package entity

import (
	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

const CategoryTableName = "categories"

type Category struct {
	ID              int         `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" json:"id"`
	Name            string      `gorm:"column:name;type:varchar(50);not null"                   json:"name"`
	Slug            string      `gorm:"column:slug;type:varchar(100);not null"                  json:"slug"`
	Description     *string     `gorm:"column:description;type:varchar(255)"                    json:"description"`
	IsPublished     bool        `gorm:"column:is_published;not null;default:false"              json:"is_publised"`
	DisplayOrder    int         `gorm:"column:display_order;type:smallint;not null"             json:"display_order"`
	MetaDescription *string     `gorm:"column:meta_description;type:text"                       json:"meta_description"`
	ThumbnailPath   *string     `gorm:"column:thumbnail_path;type:varchar(255)"                 json:"thumbnail_path"`
	ParentID        *int        `gorm:"column:parent_id;type:bigint"                            json:"parent_id"`
	Parent          *Category   `gorm:"foreignKey:ParentID"`
	Children        []*Category `gorm:"-:all"`
	model.BaseModel
}

func (Category) TableName() string {
	return CategoryTableName
}
