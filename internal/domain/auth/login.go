package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func checkPassword(plainPassword string, hashFromDb string) bool {
	// ComparePasswordAndHash performs a constant-time comparison between a
	// plain-text password and Argon2id hash, using the parameters and salt
	// contained in the hash. It returns true if they match, otherwise it returns
	// false.
	match, err := argon2id.ComparePasswordAndHash(plainPassword, hashFromDb)
	if err != nil {
		log.Fatal(err)
	}

	return match
}
