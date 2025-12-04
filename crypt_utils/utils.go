package crypt_utils

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
	"os"
)

func GetJWTPrivateKeyPath() string {
	return fmt.Sprintf("%s/%s", constCertDir, constJWTSigningKeyName)
}

func GetJWTPublicKeyPath() string {
	return fmt.Sprintf("%s/%s", constCertDir, constJWTSigningPublicKeyName)
}

func getCertsDirPath() string {
	return constCertDir
}

func LoadECDSAPrivateKeyFromPEM() (*ecdsa.PrivateKey, error) {
	keyData, err := os.ReadFile(GetJWTPrivateKeyPath())
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		ecKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("failed to parse ECDSA private key")
		}
		return ecKey, nil
	}

	ecKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ECDSA private key: %v", err)
	}

	return ecKey, nil
}

// writePrivateKeyFile takes the bytes of a private key and writes it to a file on disk at the specified path.
func writePrivateKeyFile(filePath string, value []byte) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error opening file for writing: %v", err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			slog.Error("error closing private key file", "err", err)
		}
	}(f)

	err = pem.Encode(f, &pem.Block{Type: constPrivateKeyHeader, Bytes: value})
	if err != nil {
		return fmt.Errorf("error encoding key: %v", err)
	}

	return nil
}

// writePublicKeyFile takes the bytes of a public key and writes it to a file on disk at the specified path.
func writePublicKeyFile(filePath string, value []byte) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error opening file for writing: %v", err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			slog.Error("error closing node key file", "err", err)
		}
	}(f)

	err = pem.Encode(f, &pem.Block{Type: constPublicKeyHeader, Bytes: value})
	if err != nil {
		return fmt.Errorf("error encoding key: %v", err)
	}

	return nil
}
