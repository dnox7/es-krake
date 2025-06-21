package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const PlatformAccountEnterpriseAccessTableName = "platform_account_enterprise_accesses"

type PlatformAccountEnterpriseAccess struct {
	ID                int  `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	PlatformAccountID int  `gorm:"column:platform_account_id;type:bigint;not null" json:"platform_account_id"`
	Enabled           bool `gorm:"column:enabled;type:boolean;not null;default:true" json:"enabled"`
	EnterpriseID      int  `gorm:"column:enterprise_id;type:bigint;not null" json:"enterprise_id"`
	model.BaseModelWithDeleted
}

func (PlatformAccountEnterpriseAccess) TableName() string {
	return PlatformAccountEnterpriseAccessTableName
}
