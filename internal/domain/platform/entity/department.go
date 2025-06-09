package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

type Department struct {
	ID   int    `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement;unique" json:"id"`
	Code string `gorm:"column:code;type:varchar(50);not null;unique"                   json:"code"`
	Name string `gorm:"column:name;type:varchar(255);not null"                         json:"name"`
	model.BaseModel
}
