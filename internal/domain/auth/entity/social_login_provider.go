package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const SocialLoginProviderTableName = "social_login_provider"

type SocialLoginProvider struct {
	ID           int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProviderName string `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`
	model.BaseModel
}

func (SocialLoginProvider) TableName() string {
	return SocialLoginProviderTableName
}
