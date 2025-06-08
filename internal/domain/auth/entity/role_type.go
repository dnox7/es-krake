package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const RoleTypeTableName = "role_types"

type RoleType struct {
	ID   int    `gorm:"column:id;primaryKey;type:smallint;autoIncrement;not null" json:"id"`
	Name string `gorm:"column:name;type:varchar(50);not null"                     json:"name"`
	model.BaseModel
}

func (RoleType) TableName() string {
	return RoleTypeTableName
}
