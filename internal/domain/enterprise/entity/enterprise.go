package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

type Enterprise struct {
	ID int `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	model.BaseModel
}
