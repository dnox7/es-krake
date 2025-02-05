package entity

import "pech/es-krake/internal/domain"

type Product struct {
	ID                    int     `json:"id"`                                                   // Unique identifier for the product
	Name                  string  `json:"name"`                                                 // Name of the product
	SKU                   string  `json:"sku"`                                                  // Stock Keeping Unit, unique code for inventory tracking
	Description           string  `json:"description"`                                          // Detailed description of the product
	Price                 float32 `json:"price"`                                                // Price of the product
	HasOptions            bool    `json:"has_options" db:"has_options`                          // Indicates if the product has options (e.g., size, color)
	IsAllowedToOrder      bool    `json:"is_allowed_to_order" db:"is_allowed_to_order"`         // Determines if the product can be ordered
	IsPublished           bool    `json:"is_published" db:"is_published"`                       // Indicates if the product is published and visible to customers
	IsFeatured            bool    `json:"is_featured" db:"is_featured"`                         // Marks the product as a featured item for promotions
	IsVisibleIndividually bool    `json:"is_visible_individually" db:"is_visible_individually"` // Specifies if the product is displayed as an individual item
	StockTrackingEnabled  bool    `json:"stock_tracking_enabled" db:"stock_tracking_enabled"`   // Determines if stock tracking is enabled for the product
	StockQuantity         int64   `json:"stock_quantity" db:"stock_quantity"`                   // Quantity of the product available in stock
	TaxClassID            int64   `json:"tax_class_id" db:"tax_class_id"`                       // Identifier for the tax class applied to the product
	MetaTitle             *string `json:"meta_title" db:"meta_title"`                           // SEO title for the product, used in search engine optimization
	MetaKeyword           *string `json:"meta_keyword" db:"meta_keyword"`                       // SEO keywords for the product, used for search optimization
	domain.BaseModel
}
