package jwt

import (
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
