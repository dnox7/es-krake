package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const SaleDetailPaypalTableName = "sale_detail_paypals"

type SaleDetailPaypal struct {
	ID             int     `gorm:"column:id;type:bigint;primaryKey;autoIncrement;not null;unique" json:"id"`
	SaleID         int     `gorm:"column:sale_id;type:bigint;not null" json:"sale_id"`
	ResponseID     *string `gorm:"column:response_id;type:varchar(255);unique" json:"response_id"`
	ResponseStatus string  `gorm:"column:response_status;type:varchar(50)" json:"response_status"`
	PayerID        *string `gorm:"column:payer_id;type:varchar(255)" json:"payer_id"`
	PaypalDebugID  string  `gorm:"column:paypal_debug_id;type:varchar(255);unique;not null" json:"paypal_debug_id"`
	model.BaseModel
}

func (SaleDetailPaypal) TableName() string {
	return SaleDetailPaypalTableName
}
