package entity

import "pech/es-krake/internal/domain"

type ProductOption struct {
	ID          int     `json:"id" db:"id"`                   // Unique identifier for the product option
	ProductID   int     `json:"product_id" db:"product_id"`   // The ID of the product that this option belongs to
	Name        string  `json:"name" db:"name"`               // The name of the option
	Description *string `json:"description" db:"description"` // Optional description of the option (e.g., "Select the size of the product")
	Attributes  []*ProductAttributeValue
	domain.BaseModel
}
