package entity

import (
	"time"

	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

const SaleTableName = "sales"

type Sale struct {
	ID              int        `gorm:"column:id;type:bigint;primaryKey;autoIncrement;not null;unique" json:"id"`
	OrderID         int        `gorm:"column:order_id;type:bigint;not null" json:"order_id"`
	StatusID        int        `gorm:"column:status_id;type:bigint;not null" json:"status_id"`
	PaymentMethodID int        `gorm:"column:payment_method_id;type:bigint;not null" json:"payment_method_id"`
	PaidAt          *time.Time `gorm:"column:paid_at;type:timestamp" mapstructure:"paid_at"`

	model.BaseModel
}

func (Sale) TableName() string {
	return SaleTableName
}
