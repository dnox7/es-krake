package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const CartItemTableName = "cart_items"

type CartItem struct {
	ID        int  `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	CartID    int  `gorm:"column:cart_id;type:bigint;not null" json:"cart_id"`
	ProductID int  `gorm:"column:product_id;type:bigint;not null" json:"product_id"`
	Quantity  int  `gorm:"column:quantity;type:int;not null" json:"quantity"`
	Price     int  `gorm:"column:price;type:int;not null" json:"price"`
	Selected  bool `gorm:"column:selected;type:boolean;not null" json:"selected"`
	model.BaseModelWithDeleted
}

func (CartItem) TableName() string {
	return CartItemTableName
}
