package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const BatchLogTypeTableName = "batch_log_types"

type BatchLogType struct {
	ID   int    `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	Type string `gorm:"column:type;type:varchar(255);not null;unique"           json:"type"`
	Name string `gorm:"column:name;type:varchar(255)"                           json:"name"`
	model.BaseModel
}

func (BatchLogType) TableName() string {
	return BatchLogTypeTableName
}
