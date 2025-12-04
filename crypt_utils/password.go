package crypt_utils

import (
	"github.com/go-crypt/crypt/algorithm"
	"github.com/go-crypt/crypt/algorithm/argon2"
)

func HashPassword(password string) (string, error) {
	var (
		hasher *argon2.Hasher
		err    error
		digest algorithm.Digest
	)

	hasher, err = argon2.New(
		argon2.WithProfileRFC9106LowMemory(),
	)
	if err != nil {
		return "", err
	}

	digest, err = hasher.Hash(password)
	if err != nil {
		return "", err
	}

	return digest.Encode(), nil
}
