package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"runtime"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/nats-io/nkeys"
)

// Warning! You will want this in a secret store
// of some kind instead of the source code like this
const NKeySeed string = "SAAIRUPUPZ4CZZX4EYX2MF6A3KN7FGI3AQPEVF3HI2JXVNV6DJGSTZGDFE"

// Use your own google app client ID here
const ClientID = "476732082534-gpp6r1t67lirjfddjckm81vr6c2kikr1.apps.googleusercontent.com"

func main() {
	nc, err := nats.Connect(nats.DefaultURL, nats.UserInfo("auth", "auth"))
	if err != nil {
		log.Fatal(err)
	}

	kp, err := nkeys.FromSeed([]byte(NKeySeed))
	if err != nil {
		log.Fatal(err)
	}

	v, err := NewGoogleVerifier(context.Background(), ClientID)
	if err != nil {
		log.Fatal(err)
	}

	workspace, err := NewWorkspaceKV(nc, "chat_workspace")
	if err != nil {
		log.Fatal(err)
	}

	auth := NewAuthService(kp, func(req *jwt.AuthorizationRequestClaims) (*jwt.UserClaims, error) {
		log.Println("Received Request")

		claims := jwt.NewUserClaims(req.UserNkey)
		claims.Audience = "CHAT"

		if req.ConnectOptions.Username == "cli" && req.ConnectOptions.Password == "my-password" {
			// Return claims with no permissions, backdoor method
			return claims, nil
		}

		gClaims, err := v.VerifyGoogleJWT(req.ConnectOptions.Token)
		if err != nil {
			return nil, err
		}

		userId := base64.StdEncoding.EncodeToString([]byte(gClaims.Email))

		// Add user to workspace
		err = workspace.AddUser(&WorkspaceUser{
			Id:       userId,
			Name:     gClaims.Name,
			Email:    gClaims.Email,
			PhotoURL: gClaims.Picture,
		})
		if err != nil {
			return nil, err
		}

		// Assign Permissions
		claims.Name = userId
		claims.Permissions = jwt.Permissions{
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

		return claims, nil
	})

	_, err = micro.AddService(nc, micro.Config{
		Name:        "auth",
		Version:     "0.0.1",
		Description: "Handle authorization of google jwts for chat applications",
		Endpoint: &micro.EndpointConfig{
			Subject: "$SYS.REQ.USER.AUTH",
			Handler: auth,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on $SYS.REQ.USER.AUTH...")
	// TODO: do a real graceful shutdown here
	runtime.Goexit()
}
