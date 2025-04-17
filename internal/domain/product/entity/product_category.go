package entity

import (
	"pech/es-krake/internal/domain"
)

type ProductCategory struct {
	ID                int  `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement"     json:"id"`
	ProductID         int  `gorm:"column:product_id;type:bigint;not null"                      json:"product_id"`
	CategoryID        int  `gorm:"column:category_id;type:bigint;not null"                     json:"category_id"`
	DisplayOrder      int  `gorm:"column:display_order;type:int"                               json:"display_order"`
	IsFeaturedProduct bool `gorm:"column:is_featured_product;type:bool;not null;default:false" json:"is_featured_product"`
	domain.BaseModel
}
