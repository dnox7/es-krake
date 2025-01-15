package entity

import "pech/es-krake/internal/domain/utils"

type ProductAttributeValue struct {
	ID              int           // Unique identifier for the product attribute value
	ProductID       int           // The ID of the product this attribute value is associated with
	AttributeID     int           // The ID of the attribute that this value corresponds to (e.g., "Color", "Size")
	Value           AttributeType // The actual value of the attribute (e.g., "Red", "Large")
	utils.BaseModel               // Common fields like created_at, updated_at, for auditing and versioning
}
