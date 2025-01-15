package entity

import "pech/es-krake/internal/domain/utils"

type ProductCategory struct {
	ID                int  // Unique identifier for the product-category relationship
	ProductID         int  // The ID of the product that is part of the category
	CategoryID        int  // The ID of the category to which the product belongs
	DisplayOrder      int  // The order in which the product appears within the category
	IsFeaturedProduct bool // Flag indicating if the product is featured within the category
	utils.BaseModel        // Common fields like created_at, updated_at for auditing purposes
}
