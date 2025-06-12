package entity

import (
	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

const OptionAttributeValueTableName = "option_attribute_values"

type OptionAttributeValue struct {
	ID              int        `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	AttributeID     int        `gorm:"column:attribute_id;type:bigint;not null"                json:"attribute_id"`
	ProductOptionID int        `gorm:"column:product_option_id;type:bigint;not null"           json:"product_option_id"`
	Value           string     `gorm:"column:value;type:varchar(50);not null"                  json:"value"`
	Attribute       *Attribute `gorm:"<-:false;foreignKey:AttributeID"`
	model.BaseModelWithDeleted
}

func (OptionAttributeValue) TableName() string {
	return OptionAttributeValueTableName
}
