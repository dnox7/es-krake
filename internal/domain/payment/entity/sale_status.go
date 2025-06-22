package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const SaleStatusTableName = "sale_statuses"

type SaleStatus struct {
	ID   int    `gorm:"column:id;type:bigint;primaryKey;autoIncrement;not null;unique" json:"id"`
	Name string `gorm:"column:name;type:varchar(20);not null" json:"name"`
	model.BaseModel
}

func (SaleStatus) TableName() string {
	return SaleStatusTableName
}
