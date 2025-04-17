package entity

import "pech/es-krake/internal/domain"

type Category struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	Description     string  `json:"description"`
	IsPublished     bool    `json:"is_publised" `
	DisplayOrder    int     `json:"display_order" `
	MetaDescription *string `json:"meta_description" `
	ThumbnailURL    *string `json:"thumbnail_url"`
	Children        []*Category
	Parents         []*Category
	domain.BaseModel
}

type CategoryParent struct {
	ChildID  int `json:"child_id" `
	ParentID int `json:"parent_id"`
	Child    *Category
	Parent   *Category
	domain.BaseModel
}
