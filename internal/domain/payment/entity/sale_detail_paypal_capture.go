package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const SaleDetailPaypalCaptureTableName = "sale_detail_paypal_captures"

type SaleDetailPaypalCapture struct {
	ID                    int     `gorm:"column:id;type:bigint;primaryKey;autoIncrement;not null;unique" json:"id"`
	SaleDetailPaypalID    int     `gorm:"column:sale_detail_paypal_id;type:bigint;not null" json:"sale_detail_paypal_id"`
	PaymentsCaptureID     string  `gorm:"column:payments_capture_id;type:varchar(50);unique;not null" json:"payments_capture_id"`
	PaymentsCaptureStatus string  `gorm:"column:payments_capture_status;type:varchar(50);not null" json:"payments_capture_status"`
	LastCaptureStatus     *string `gorm:"column:last_capture_status;type:varchar(50)" json:"last_capture_status"`
	RefundedAt            *string `gorm:"column:refunded_at;type:timestamp" json:"refunded_at"`
	RefundedPaypalDebugID *string `gorm:"column:refunded_paypal_debug_id;type:varchar(255);unique" json:"refunded_paypal_debug_id"`
	model.BaseModel
}

func (SaleDetailPaypalCapture) TableName() string {
	return SaleDetailPaypalCaptureTableName
}
