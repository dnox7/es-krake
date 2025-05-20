package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

type Role struct {
	ID              int              `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	RoleTypeID      int              `gorm:"column:role_type_id;type:smallint;not null;default:1" json:"role_type_id"`
	Name            string           `gorm:"column:name;type:varchar(50);not null" json:"name"`
	DisplayOrder    int              `gorm:"column:displaty_order;type:int;not null" json:"displaty_order"`
	RolePermissions []RolePermission `gorm:"foreignKey:RoleID"`
	model.BaseModelWithDeleted
}
