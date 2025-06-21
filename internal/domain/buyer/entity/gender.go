package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const GenderTableName = "genders"

type Gender struct {
	ID   int    `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	Name string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	model.BaseModel
}

func (Gender) TableName() string {
	return GenderTableName
}
