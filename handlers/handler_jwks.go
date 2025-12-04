package handlers

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log/slog"

	"jwt-auth-poc/jwt"
	"jwt-auth-poc/middlewares"
	"net/http"
	"os"

	"github.com/go-jose/go-jose/v4"
)

// handleJWKSPublicKeyGET returns the JWKS with the JWT public key
func handleJWKSPublicKeyGET(ctx *middlewares.AppContext) {
	bytes, err := os.ReadFile(jwt.GetJWTPublicKeyPath())
	if err != nil {
		ctx.Logger.Error("failed to read jwt public key", "err", err)
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		ctx.Logger.Error("failed to decode jwt public key")
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	switch block.Type {
	case "PRIVATE KEY", "EC PRIVATE KEY", "RSA PRIVATE KEY":
		ctx.Logger.Error("decoded private key, canceling request")
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	case "PUBLIC KEY":
	default:
		ctx.Logger.Error("unknown jwt public key type", "type", block.Type)
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		ctx.Logger.Error("failed to parse jwt public key", "err", err)
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	var alg string
	switch key := pubKey.(type) {
	case *rsa.PublicKey:
		alg = string(jose.RS256)
	case *ecdsa.PublicKey:
		switch key.Params().BitSize {
		case 256:
			alg = string(jose.ES256)
		case 384:
			alg = string(jose.ES384)
		case 521:
			alg = string(jose.ES512)
		default:
			ctx.Logger.Error("unsupported EC curve size", "bits", key.Params().BitSize)
			ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
			return
		}
	case ed25519.PublicKey:
		alg = string(jose.EdDSA)
	default:
		ctx.Logger.Error("unsupported public key type")
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	keyID := generateKeyID(ctx.Logger, pubKey)

	jwk := jose.JSONWebKey{
		Key:       pubKey,
		KeyID:     keyID,
		Algorithm: alg,
		Use:       "sig",
	}

	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{jwk},
	}

	ctx.Response.Header().Set("Cache-Control", "public, max-age=3600")
	ctx.WriteJSON(http.StatusOK, jwks)
}

// generateKeyID creates a unique identifier for the key
func generateKeyID(logger *slog.Logger, key interface{}) string {
	jwk := jose.JSONWebKey{Key: key}
	thumbprint, err := jwk.Thumbprint(crypto.SHA256)

	if err != nil {
		logger.Warn("failed to generate key thumbprint, using default", "err", err)
		return "default"
	}

	return base64.RawURLEncoding.EncodeToString(thumbprint)
}
