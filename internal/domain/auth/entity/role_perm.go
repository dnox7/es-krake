package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

type RolePermission struct {
	ID           int         `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	RoleID       int         `gorm:"column:role_id;type:bigint;not null" json:"role_id"`
	PermissionID int         `gorm:"column:permission_id;type:bigint;not null" json:"permission_id"`
	Permission   *Permission `gorm:"foreignKey:PermissionID"`
	model.BaseModelWithDeleted
}
