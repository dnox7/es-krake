package entity

import "pech/es-krake/internal/domain/utils"

type AttributeType interface{}

type CategoryAttribute struct {
	ID           int
	CategoryID   int
	Name         string
	Type         AttributeType
	Value        AttributeType
	DefaultValue *AttributeType
	Description  *string
	utils.BaseModel
}
