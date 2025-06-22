package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const PaymentMethodTableName = "payment_methods"

type PaymentMethod struct {
	ID             int      `gorm:"column:id;type:bigint;primaryKey;autoIncrement;not null;unique" json:"id"`
	Code           string   `gorm:"column:code;type:varchar(255);not null;unique" json:"code"`
	Name           string   `gorm:"column:name;type:varchar(255);not null" json:"name"`
	ShortName      string   `gorm:"column:short_name;type:varchar(255);not null" json:"short_name"`
	CommissionRate *float64 `gorm:"column:commission_rate;type:decimal(10,2)" json:"commission_rate"`
	model.BaseModel
}

func (PaymentMethod) TableName() string {
	return PaymentMethodTableName
}
