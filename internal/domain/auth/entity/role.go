package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const RoleTableName = "roles"

type Role struct {
	ID              int              `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	RoleTypeID      int              `gorm:"column:role_type_id;type:smallint;not null;default:1"    json:"role_type_id"`
	Name            string           `gorm:"column:name;type:varchar(50);not null"                   json:"name"`
	DisplayOrder    int              `gorm:"column:display_order;type:int;not null"                  json:"display_order"`
	RolePermissions []RolePermission `gorm:"foreignKey:RoleID"`
	model.BaseModelWithDeleted
}

func (Role) TableName() string {
	return RoleTableName
}
