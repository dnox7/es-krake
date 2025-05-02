package entity

import (
	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

type AttributeType struct {
	ID   int    `gorm:"column:id;primaryKey;type:smallint;autoIncrement;not null" json:"id"`
	Name string `gorm:"column:name;type:varchar(20);not null"                     json:"name"`
	model.BaseModel
}
