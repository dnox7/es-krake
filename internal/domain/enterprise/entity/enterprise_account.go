package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const EnterpriseAccountTableName = "enterprise_accounts"

type EnterpriseAccount struct {
	ID           int     `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	KcUserID     string  `gorm:"column:kc_user_id;type:varchar(255);not null;unique"     json:"kc_user_id"`
	RoleID       int     `gorm:"column:role_id;type:bigint;not null"                     json:"role_id"`
	EnterpriseID int     `gorm:"column:enterprise_id;type:bigint;not null"               json:"enterprise_id"`
	HasPassword  bool    `gorm:"column:has_password;type:tinyint(1);not null;default:1"  json:"has_password"`
	Notes        *string `gorm:"column:notes;type:varchar(512)"                          json:"notes"`
	model.BaseModelWithDeleted
}

func (EnterpriseAccount) TableName() string {
	return EnterpriseAccountTableName
}
