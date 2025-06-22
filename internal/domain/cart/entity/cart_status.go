package entity

import (
	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

const CartStatusTableName = "cart_statuses"

type CartStatus struct {
	ID     int    `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	Status string `gorm:"column:status;type:varchar(255);not null" json:"status"`
	model.BaseModel
}

func (CartStatus) TableName() string {
	return CartStatusTableName
}
