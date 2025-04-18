package entity

import "pech/es-krake/internal/domain"

type Category struct {
	ID              int         `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" json:"id"`
	Name            string      `gorm:"column:name;type:varchar(50);not null"                   json:"name"`
	Slug            string      `gorm:"column:slug;type:varchar(50);not null"                   json:"slug"`
	Description     *string     `gorm:"column:description;type:text"                            json:"description"`
	IsPublished     bool        `gorm:"column:is_publised;not null;default:false"               json:"is_publised"`
	DisplayOrder    *int        `gorm:"column:display_order;type:smallint;"                     json:"display_order"`
	MetaDescription *string     `gorm:"column:meta_description;type:text"                       json:"meta_description"`
	ThumbnailURL    *string     `gorm:"column:thumbnail_url;type:varchar"                       json:"thumbnail_url"`
	Children        []*Category `gorm:"-:all"`
	Parents         []*Category `gorm:"-:all"`
	domain.BaseModel
}

type CategoryParent struct {
	ChildID  int `json:"child_id"`
	ParentID int `json:"parent_id"`
	Child    *Category
	Parent   *Category
	domain.BaseModel
}
