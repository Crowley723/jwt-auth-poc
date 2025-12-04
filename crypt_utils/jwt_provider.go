package crypt_utils

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

type JWTProvider interface {
	Sign(claims jwt.Claims) (string, error)
	Validate(token string) (*jwt.Claims, error)
}

type ecdsaJWTProvider struct {
	signer    jose.Signer
	publicKey *ecdsa.PublicKey
}

func NewECDSAJWTProvider(privateKey *ecdsa.PrivateKey) (JWTProvider, error) {
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.ES256,
			Key:       privateKey,
		},
		&jose.SignerOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECDSA JWT signer: %w", err)
	}

	return &ecdsaJWTProvider{
		signer:    signer,
		publicKey: &privateKey.PublicKey,
	}, nil
}

func (p *ecdsaJWTProvider) Sign(claims jwt.Claims) (string, error) {
	token, err := jwt.Signed(p.signer).Claims(claims).Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return token, nil
}

func (p *ecdsaJWTProvider) Validate(token string) (*jwt.Claims, error) {
	parsed, err := jwt.ParseSigned(token, []jose.SignatureAlgorithm{jose.ES256})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	var claims jwt.Claims
	if err := parsed.Claims(p.publicKey, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if err := claims.Validate(jwt.Expected{
		Time: time.Now(),
	}); err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return &claims, nil
}
