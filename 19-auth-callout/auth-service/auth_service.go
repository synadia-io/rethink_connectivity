package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nats.go/micro"
	"github.com/nats-io/nkeys"
)

type WorkspaceUser struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	PhotoURL string `json:"photoURL"`
}

type GoogleClaims struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
}

type AuthService struct {
	issuerKeyPair nkeys.KeyPair
	verifier      *oidc.IDTokenVerifier
	clientID      string
	workspace     jetstream.KeyValue

	username string
	password string
}

func NewAuthService(issuer nkeys.KeyPair, clientID string, username string, password string) *AuthService {
	return &AuthService{
		issuerKeyPair: issuer,
		clientID:      clientID,
		username:      username,
		password:      password,
	}
}

func (a *AuthService) Run(ctx context.Context) error {
	// Initialize OIDC Provider
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return err
	}
	a.verifier = provider.Verifier(&oidc.Config{ClientID: ClientID})

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL, nats.UserInfo("auth", "auth"))
	if err != nil {
		return err
	}

	// get the workspace kv
	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}

	kv, err := js.KeyValue(ctx, "chat_workspace")
	if err != nil {
		return err
	}
	a.workspace = kv

	// Add our Auth Service
	_, err = micro.AddService(nc, micro.Config{
		Name:        "auth",
		Version:     "0.0.1",
		Description: "Handle authorization of Google JWTs for chat applications",
		Endpoint: &micro.EndpointConfig{
			Subject: "$SYS.REQ.USER.AUTH",
			Handler: a,
		},
	})
	if err != nil {
		return err
	}

	log.Println("listening on $SYS.REQ.USER.AUTH")
	runtime.Goexit()
	return nil
}

func (a *AuthService) Handle(r micro.Request) {
	log.Println("Received Request")

	rc, err := jwt.DecodeAuthorizationRequestClaims(string(r.Data()))
	if err != nil {
		log.Println("Error", err)
		r.Error("500", err.Error(), nil)
	}

	userNkey := rc.UserNkey
	serverId := rc.Server.ID
	claims := jwt.NewUserClaims(rc.UserNkey)
	claims.Audience = "CHAT"

	// this gives me a backdoor with the CLI. Don't do this in production!
	if rc.ConnectOptions.Username == "cli" && rc.ConnectOptions.Password == "my-password" {
		claims.Name = rc.ConnectOptions.Username
		claims.Permissions = jwt.Permissions{}

		token, err := ValidateAndSign(claims, a.issuerKeyPair)
		a.Respond(r, userNkey, serverId, token, err)
		return
	} else {
		// Try to get a google JWT from the token field
		googleJWT := rc.ConnectOptions.Token
		gclaims, err := a.VerifyGoogleJWT(context.Background(), googleJWT)
		if err != nil {
			log.Println("error", err)
			a.Respond(r, userNkey, serverId, "", err)
		}

		log.Printf("google claims: %+v", gclaims)

		// Add user to workspace kv
		err = a.AddUserToWorkspace(gclaims)
		if err != nil {
			log.Println("error", err)
			a.Respond(r, userNkey, serverId, "", err)
		}

		// Assign JWT permissions based off the
		a.AssignPermissions(gclaims, claims)

		token, err := ValidateAndSign(claims, a.issuerKeyPair)
		a.Respond(r, userNkey, serverId, token, err)
	}
}

func (a *AuthService) Respond(req micro.Request, userNKey, serverId, userJWT string, err error) {
	rc := jwt.NewAuthorizationResponseClaims(userNKey)
	rc.Audience = serverId
	rc.Jwt = userJWT
	if err != nil {
		rc.Error = err.Error()
	}

	token, err := rc.Encode(a.issuerKeyPair)
	if err != nil {
		log.Println("error encoding response jwt:", err)
	}

	req.Respond([]byte(token))
}

func ValidateAndSign(claims *jwt.UserClaims, kp nkeys.KeyPair) (string, error) {
	// Validate the claims.
	vr := jwt.CreateValidationResults()
	claims.Validate(vr)
	if len(vr.Errors()) > 0 {
		return "", errors.Join(vr.Errors()...)
	}

	// Sign it with the issuer key since this is non-operator mode.
	return claims.Encode(kp)
}

func (a *AuthService) VerifyGoogleJWT(ctx context.Context, token string) (*GoogleClaims, error) {
	claims := &GoogleClaims{}

	idToken, err := a.verifier.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	if err := idToken.Claims(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func (a *AuthService) AddUserToWorkspace(gclaims *GoogleClaims) error {
	user := WorkspaceUser{
		Id:       base64.StdEncoding.EncodeToString([]byte(gclaims.Email)),
		Name:     gclaims.Name,
		Email:    gclaims.Email,
		PhotoURL: gclaims.Picture,
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = a.workspace.Put(context.Background(), fmt.Sprintf("users.%s", user.Id), data)
	return err
}

func (a *AuthService) AssignPermissions(gclaims *GoogleClaims, uc *jwt.UserClaims) {
	userId := base64.StdEncoding.EncodeToString([]byte(gclaims.Email))
	uc.Name = userId

	uc.Permissions = jwt.Permissions{
		Pub: jwt.Permission{
			Allow: jwt.StringList{
				"$JS.API.INFO", // General JS Info

				// Chat permisions
				fmt.Sprintf("chat.*.%s", userId),            // Publishing chat messages for this user id
				"$JS.API.STREAM.INFO.chat_messages",         // Getting info on chat_messages stream
				"$JS.API.CONSUMER.CREATE.chat_messages.>",   // Creating consumers on chat_messages stream
				"$JS.API.CONSUMER.MSG.NEXT.chat_messages.>", // Creating consumers on chat_messages stream

				// Workspace permissions
				"$JS.API.DIRECT.GET.KV_chat_workspace.>",        // Gets from workspace KV
				"$JS.API.STREAM.INFO.KV_chat_workspace",         // Info about workspace KV
				"$JS.API.CONSUMER.CREATE.KV_chat_workspace.*.>", // Creating consumers/watchers on workspace KV
			},
		},
	}
}
