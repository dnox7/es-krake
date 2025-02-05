package entity

import "pech/es-krake/internal/domain"

type ProductAttributeValue struct {
	ID          int    `json:"id"`                             // Unique identifier for the product attribute value
	ProductID   int    `json:"product_id" db:"product_id"`     // The ID of the product this attribute value is associated with
	AttributeID int    `json:"attribute_id" db:"attribute_id"` // The ID of the attribute that this value corresponds to (e.g., "Color", "Size")
	Value       string `json :"value"`                         // The actual value of the attribute (e.g., "Red", "Large")
	Attribute   *Attribute
	Type        *AttributeType
	domain.BaseModel
}
