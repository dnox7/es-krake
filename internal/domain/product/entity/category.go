package entity

import (
	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

const CategoryTableName = "categories"

type Category struct {
	ID              int         `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" json:"id"`
	Name            string      `gorm:"column:name;type:varchar(50);not null"                   json:"name"`
	Slug            string      `gorm:"column:slug;type:varchar(50);not null"                   json:"slug"`
	Description     *string     `gorm:"column:description;type:text"                            json:"description"`
	IsPublished     bool        `gorm:"column:is_publised;not null;default:false"               json:"is_publised"`
	DisplayOrder    *int        `gorm:"column:display_order;type:smallint;"                     json:"display_order"`
	MetaDescription *string     `gorm:"column:meta_description;type:text"                       json:"meta_description"`
	ThumbnailURL    *string     `gorm:"column:thumbnail_url;type:varchar"                       json:"thumbnail_url"`
	ParentID        *int        `gorm:"column:parent_id;type:bigint"                            json:"parent_id"`
	Children        []*Category `gorm:"-:all"`
	Parent          *Category   `gorm:"-:all"`
	model.BaseModel
}

func (Category) TableName() string {
	return CategoryTableName
}
