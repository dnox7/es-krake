package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

type Product struct {
	ID                    int              `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement"     json:"id"`
	Name                  string           `gorm:"column:name;type:varchar(250);not null"                      json:"name"`
	SKU                   string           `gorm:"column:sku;type:varchar(50);not null"                        json:"sku"`
	Description           string           `gorm:"column:description;type:text"                                json:"description"`
	Price                 float64          `gorm:"column:price;type:float;not null"                            json:"price"`
	HasOptions            bool             `gorm:"column:has_options;type:bool;not null;default:false"         json:"has_options"`
	IsAllowedToOrder      bool             `gorm:"column:is_allowed_to_order;type:bool;not null;default:false" json:"is_allowed_to_order"`
	IsPublished           bool             `gorm:"column:is_published;type:bool;not null;default:false"        json:"is_published"`
	IsFeatured            bool             `gorm:"column:is_featured;type:bool;not null;default:false"         json:"is_featured"`
	IsVisibleIndividually bool             `gorm:"column:is_visible_individually;not null;default:false"       json:"is_visible_individually"`
	StockTrackingEnabled  bool             `gorm:"column:stock_tracking_enabled;not null;default:false"        json:"stock_tracking_enabled"`
	StockQuantity         int64            `gorm:"column:stock_quantity;type:bigint;not null"                  json:"stock_quantity"`
	TaxClassID            int64            `gorm:"column:tax_class_id;type:smallint"                           json:"tax_class_id"`
	Options               []*ProductOption `gorm:"-"`
	Categories            []*Category      `gorm:"-"`
	model.BaseModel
}
