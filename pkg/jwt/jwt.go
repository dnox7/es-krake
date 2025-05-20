package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

const (
	PemPublicKeyHeader  = "-----BEGIN PUBLIC KEY-----"
	PemPublicKeyFooter  = "-----END PUBLIC KEY-----"
	PemPrivateKeyHeader = "-----BEGIN PRIVATE KEY-----"
	PemPrivateKeyFooter = "-----END PRIVATE KEY-----"
)

func GenerateES256JWT(keyID, secret string, payload map[string]interface{}) (string, error) {
	block, _ := pem.Decode([]byte(PemPrivateKeyHeader + "\n" + secret + "\n" + PemPrivateKeyFooter))
	if block == nil || block.Type != "PRIVATE KEY" {
		return "", errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = keyID
	return token.SignedString(privateKey)
}

func Verify(token, keyStr string, cfg *config.Config) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	pubKey := PemPublicKeyHeader + "\n" + keyStr + "\n" + PemPublicKeyFooter

	keyByte, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubKey))
	if err != nil {
		if cfg.App.Env == utils.TestingEnv && (err == jwt.ErrNotRSAPublicKey || err == jwt.ErrKeyMustBePEMEncoded) {
			_, parseErr := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(keyStr), nil
			})
			return claims, parseErr
		}
		return claims, err
	}

	_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return keyByte, nil
	})
	return claims, err
}

func DecodeJWTUnverified(tokenStr string) (*jwt.Token, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	return token, nil
}

func GetKeycloakUserID(claims jwt.MapClaims) string {
	sub, _ := claims.GetSubject()
	return sub
}

func GetEmail(claims jwt.MapClaims) string {
	email, ok := claims["email"].(string)
	if !ok {
		return ""
	}
	return email
}

func VerifyExpired(claims jwt.Claims) bool {
	v := jwt.NewValidator(jwt.WithExpirationRequired())
	err := v.Validate(claims)
	return !errors.Is(err, jwt.ErrTokenExpired)
}
