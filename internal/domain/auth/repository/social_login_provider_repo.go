package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type (
	SocialLoginProviderType int
	KeycloakSnsAliasType    string
)

const (
	GoogleProvider SocialLoginProviderType = iota + 1

	GoogleKeycloakSnsAlias KeycloakSnsAliasType = "google"
)

var KeycloakSnsAliasMap = map[int]string{
	int(GoogleProvider): string(GoogleKeycloakSnsAlias),
}

type SocialLoginProviderRepository interface {
	FindByCondition(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.SocialLoginProvider, error)
}
