package entity

import (
	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

const ProductCategoryTableName = "product_categories"

type ProductCategory struct {
	ID                int       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement"     json:"id"`
	ProductID         int       `gorm:"column:product_id;type:bigint;not null"                      json:"product_id"`
	CategoryID        int       `gorm:"column:category_id;type:bigint;not null"                     json:"category_id"`
	DisplayOrder      int       `gorm:"column:display_order;type:int;not null"                      json:"display_order"`
	IsFeaturedProduct bool      `gorm:"column:is_featured_product;type:bool;not null;default:false" json:"is_featured_product"`
	Product           *Product  `gorm:"foreignKey:ProductID"`
	Category          *Category `gorm:"foreignKey:CategoryID"`
	model.BaseModelWithDeleted
}

func (ProductCategory) TableName() string {
	return ProductCategoryTableName
}
