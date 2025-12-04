package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"jwt-auth-poc/crypt_utils"
	"jwt-auth-poc/db"
	"jwt-auth-poc/middlewares"
	"strconv"
	"time"

	"github.com/go-jose/go-jose/v4/jwt"
)

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func GenerateRefreshToken() (token, hash string, err error) {
	b := make([]byte, 64)
	_, err = rand.Read(b)
	if err != nil {
		return "", "", err
	}

	token = hex.EncodeToString(b)

	return token, HashToken(token), nil
}

func GenerateAccessToken(ctx *middlewares.AppContext, userDetails *db.User) (string, error) {
	var claims = jwt.Claims{
		Subject:  strconv.Itoa(userDetails.ID),
		Expiry:   jwt.NewNumericDate(time.Now().Add(crypt_utils.ConstAccessTokenValidityPeriod)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Issuer:   "http://localhost",
	}

	token, err := ctx.JWTProvider.Sign(claims)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return token, nil
}
