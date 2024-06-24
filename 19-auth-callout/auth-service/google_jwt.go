package main

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
)

type GoogleVerifier struct {
	*oidc.IDTokenVerifier
}

func NewGoogleVerifier(ctx context.Context, clientID string) (*GoogleVerifier, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}
	return &GoogleVerifier{provider.Verifier(&oidc.Config{ClientID: ClientID})}, nil
}

type GoogleClaims struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
}

// Verifies that the ID token was signed by google and is valid.
// Returns the claims embedded with the token
func (v *GoogleVerifier) VerifyGoogleJWT(token string) (*GoogleClaims, error) {
	claims := &GoogleClaims{}

	idToken, err := v.Verify(context.Background(), token)
	if err != nil {
		return nil, err
	}

	if err := idToken.Claims(claims); err != nil {
		return nil, err
	}
	//TODO:Also check for expiry here

	return claims, nil
}
