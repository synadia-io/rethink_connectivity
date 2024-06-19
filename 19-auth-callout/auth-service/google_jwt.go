package main

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
)

var ClientID = "476732082534-gpp6r1t67lirjfddjckm81vr6c2kikr1.apps.googleusercontent.com"
var verifier *oidc.IDTokenVerifier

type GoogleClaims struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
}

func InitGoogleOIDC(ctx context.Context) error {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return err
	}

	verifier = provider.Verifier(&oidc.Config{ClientID: ClientID})

	return nil
}

func VerifyGoogleJWT(ctx context.Context, token string) (*GoogleClaims, error) {
	claims := &GoogleClaims{}

	idToken, err := verifier.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	if err := idToken.Claims(claims); err != nil {
		return nil, err
	}

	return claims, nil
}
