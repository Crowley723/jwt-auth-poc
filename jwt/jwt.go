package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"os"
)

// CreateSigningKeys generates am ecdsa keypair and writes them to disk.
func CreateSigningKeys() error {
	jwtSigningKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("error generating jwt signing key: %v", err)
	}

	jwtSigningKeyBytes, err := x509.MarshalPKCS8PrivateKey(jwtSigningKey)
	if err != nil {
		return fmt.Errorf("error marshalling jwt private key: %v", err)
	}

	jwtSigningPublicKeyBytes, err := x509.MarshalPKIXPublicKey(&jwtSigningKey.PublicKey)
	if err != nil {
		return fmt.Errorf("error marshalling jwt public key: %v", err)
	}

	err = os.MkdirAll(getCertsDirPath(), 0700)
	if err != nil {
		return fmt.Errorf("unable to create certificate directory: %v", err)
	}

	err = writePrivateKeyFile(GetJWTPrivateKeyPath(), jwtSigningKeyBytes)
	if err != nil {
		return fmt.Errorf("error writing jwt private key file: %v", err)
	}

	err = writePublicKeyFile(GetJWTPublicKeyPath(), jwtSigningPublicKeyBytes)
	if err != nil {
		return fmt.Errorf("error writing jwt public key file: %v", err)
	}

	return nil
}
