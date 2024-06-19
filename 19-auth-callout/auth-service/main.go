package main

import (
	"context"
	"log"

	"github.com/nats-io/nkeys"
)

// Warning! You will want this in a secret store
// of some kind instead of the source code like this
const NKeySeed string = "SAAIRUPUPZ4CZZX4EYX2MF6A3KN7FGI3AQPEVF3HI2JXVNV6DJGSTZGDFE"

// Use your own google app client ID here
const ClientID = "476732082534-gpp6r1t67lirjfddjckm81vr6c2kikr1.apps.googleusercontent.com"

func main() {
	kp, err := nkeys.FromSeed([]byte(NKeySeed))
	if err != nil {
		log.Fatal(err)
	}

	auth := NewAuthService(kp, ClientID, "cli", "my-password")
	err = auth.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
