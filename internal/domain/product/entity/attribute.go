package entity

import "pech/es-krake/internal/domain"

const attributeTableName = "attributes"

type Attribute struct {
	ID              int     `json:"id"`                                       // Unique identifier for the attribute definition
	Name            string  `json:"name"`                                     // Name of the attribute (e.g., "Color", "Size")
	Description     *string `json:"description"`                              // Optional description for the attribute
	AttributeTypeID int     `json:"attribute_type_id" db:"attribute_type_id"` // Data type of the attribute value (e.g., string, int, date, boolean)
	IsRequired      bool    `json:"is_required" db:"is_required"`             // Whether this attribute is required for associated entities
	DisplayOrder    int     `json:"display_order" db:"display_order"`         // The order in which the attribute is displayed in a UI
	domain.BaseModel
}

func (a *Attribute) TableName() string {
	return attributeTableName
}
