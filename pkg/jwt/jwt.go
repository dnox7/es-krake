package jwt

import (
	"fmt"
	"pech/es-krake/config"
	"pech/es-krake/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

const (
	pubKeyPrefix = "-----BEGIN PUBLIC KEY-----"
	pubKeySuffix = "-----END PUBLIC KEY-----"
)

func Verify(token, keyStr string, cfg *config.Config) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	pubKey := fmt.Sprintf("%s\n%s\n%s", pubKeyPrefix, keyStr, pubKeySuffix)

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
