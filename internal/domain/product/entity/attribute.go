package entity

import "pech/es-krake/internal/domain/utils"

type AttributeType interface{}

type Attribute struct {
	ID              int           // Unique identifier for the attribute definition
	Name            string        // Name of the attribute (e.g., "Color", "Size")
	Description     *string       // Optional description for the attribute
	Type            AttributeType // Data type of the attribute value (e.g., string, int, date, boolean)
	IsRequired      bool          // Whether this attribute is required for associated entities
	DisplayOrder    int           // The order in which the attribute is displayed in a UI
	utils.BaseModel               // Common fields like created_at, updated_at
}
