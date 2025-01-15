package entity

import "pech/es-krake/internal/domain/utils"

type ProductOption struct {
	ID              int     // Unique identifier for the product option
	ProductID       int     // The ID of the product that this option belongs to
	Name            string  // The name of the option (e.g., "Size", "Color")
	Description     *string // Optional description of the option (e.g., "Select the size of the product")
	utils.BaseModel         // Common fields like created_at, updated_at for auditing and versioning
}
