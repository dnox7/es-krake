package entity

import (
	"pech/es-krake/internal/domain/shared/model"
)

type ProductAttributeValue struct {
	ID              int        `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	ProductID       int        `gorm:"column:product_id;type:bigint;not null"                  json:"product_id"`
	AttributeID     int        `gorm:"column:attribute_id;type:bigint;not null"                json:"attribute_id"`
	ProductOptionID *int       `gorm:"column:product_option_id;type:bigint"                    json:"product_option_id"`
	Value           string     `gorm:"column:value;type:varchar(50);not null"                  json:"value"`
	Attribute       *Attribute `gorm:"<-:false;foreignKey:AttributeID"`
	model.BaseModelWithDeleted
}
