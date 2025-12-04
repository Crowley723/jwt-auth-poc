package crypt_utils

import "time"

const (
	constPrivateKeyHeader  = "PRIVATE KEY"
	constPublicKeyHeader   = "PUBLIC KEY"
	constCertificateHeader = "CERTIFICATE"
)

const (
	constJWTSigningKeyName       = "crypt_utils.key"
	constJWTSigningPublicKeyName = "crypt_utils.pub"
)

const (
	constCertDir = "./app/certs"
)

const (
	constRefreshTokenValidityPeriod = 30 * 24 * time.Hour //30 days
	constAccessTokenValidityPeriod  = 24 * time.Hour      //24 hours
)
