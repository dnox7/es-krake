package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

type Action struct {
	ID   int    `gorm:"column:id;primaryKey;type:smallint;autoIncrement;not null" json:"id"`
	Name string `gorm:"column:name;type:varchar(25);not null" json:"name"`
	Code string `gorm:"column:code;type:varchar(25);not null;unique" json:"code"`
	model.BaseModel
}
