package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const AttributeTableName = "attributes"

type Attribute struct {
	ID          int     `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement;unique" json:"id"`
	Name        string  `gorm:"column:name;type:varchar(10);not null"                          json:"name"`
	Description *string `gorm:"column:description;type:text"                                   json:"description"`
	model.BaseModel
}

func (Attribute) TableName() string {
	return AttributeTableName
}
