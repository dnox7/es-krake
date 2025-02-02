package entity

import (
	"pech/es-krake/internal/domain"
)

type ProductCategory struct {
	ID                int  `json:"id" db:"id"`                                   // Unique identifier for the product-category relationship
	ProductID         int  `json:"product_id" db:"product_id"`                   // The ID of the product that is part of the category
	CategoryID        int  `json:"category_id" db:"category_id"`                 // The ID of the category to which the product belongs
	DisplayOrder      int  `json:"display_order" db:"display_order"`             // The order in which the product appears within the category
	IsFeaturedProduct bool `json:"is_featured_product" db:"is_featured_product"` // Flag indicating if the product is featured within the category
	domain.BaseModel
}
