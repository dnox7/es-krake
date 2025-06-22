package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const CartTableName = "carts"

type Cart struct {
	ID        int        `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	BuyerID   int        `gorm:"column:buyer_id;type:bigint;not null" json:"buyer_id"`
	StatusID  int        `gorm:"column:status_id;type:bigint;not null" json:"status_id"`
	CartItems []CartItem `gorm:"foreignKey:CartID;references:ID" json:"cart_items"`
	model.BaseModel
}

func (Cart) TableName() string {
	return CartTableName
}
