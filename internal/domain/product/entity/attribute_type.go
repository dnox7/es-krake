package entity

import "pech/es-krake/internal/domain"

type AttributeType struct {
	ID   int    `gorm:"column:id;primaryKey;type:smallint;autoIncrement;not null" json:"id"`
	Name string `gorm:"column:name;type:varchar(20);not null"                     json:"name"`
	domain.BaseModel
}
