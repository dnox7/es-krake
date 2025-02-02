package entity

import "pech/es-krake/internal/domain"

const attributeTypeTableName = "attribute_types"

type AttributeType struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	domain.BaseModel
}

func (at *AttributeType) TableName() string {
	return attributeTypeTableName
}
