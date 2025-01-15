package entity

type Product struct {
	ID                    int     // Unique identifier for the product
	Name                  string  // Name of the product
	SKU                   string  // Stock Keeping Unit, unique code for inventory tracking
	Description           string  // Detailed description of the product
	Price                 float32 // Price of the product
	HasOptions            bool    // Indicates if the product has options (e.g., size, color)
	IsAllowedToOrder      bool    // Determines if the product can be ordered
	IsPublished           bool    // Indicates if the product is published and visible to customers
	IsFeatured            bool    // Marks the product as a featured item for promotions
	IsVisibleIndividually bool    // Specifies if the product is displayed as an individual item
	StockTrackingEnabled  bool    // Determines if stock tracking is enabled for the product
	StockQuantity         int64   // Quantity of the product available in stock
	TaxClassID            int64   // Identifier for the tax class applied to the product
	MetaTitle             string  // SEO title for the product, used in search engine optimization
	MetaKeyword           string  // SEO keywords for the product, used for search optimization
}
