package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const ProductTableName = "products"

type Product struct {
	ID                   int              `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement"     json:"id"`
	Name                 string           `gorm:"column:name;type:varchar(255);not null"                      json:"name"`
	SKU                  string           `gorm:"column:sku;type:varchar(50);not null"                        json:"sku"`
	Description          *string          `gorm:"column:description;type:varchar(255)"                        json:"description"`
	Price                float64          `gorm:"column:price;type:float;not null"                            json:"price"`
	HasOptions           bool             `gorm:"column:has_options;type:bool;not null;default:false"         json:"has_options"`
	IsAllowedToOrder     bool             `gorm:"column:is_allowed_to_order;type:bool;not null;default:false" json:"is_allowed_to_order"`
	IsPublished          bool             `gorm:"column:is_published;type:bool;not null;default:false"        json:"is_published"`
	StockTrackingEnabled bool             `gorm:"column:stock_tracking_enabled;not null;default:true"         json:"stock_tracking_enabled"`
	StockQuantity        int64            `gorm:"column:stock_quantity;type:bigint;not null;default:0"        json:"stock_quantity"`
	ThumbnailPath        *string          `gorm:"column:thumbnail_path;type:varchar(255)"                     json:"thumbnail_path"`
	BrandID              int              `gorm:"column:brand_id;type:bigint;not null" json:"brand_id"`
	Brand                *Brand           `gorm:"foreignKey:BrandID"`
	Options              []*ProductOption `gorm:"-"`
	Categories           []*Category      `gorm:"-"`
	model.BaseModel
}

func (Product) TableName() string {
	return ProductTableName
}
