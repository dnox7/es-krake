package entity

import "pech/es-krake/internal/domain"

type Category struct {
	ID               int         `json:"id" db:"id"`                             // Unique identifier for the category
	Name             string      `json:"name"`                                   // Name of the category
	Slug             string      `json:"slug"`                                   // URL-friendly identifier for the category (used in routing)
	Description      string      `json:"description"`                            // Detailed description of the category
	IsPublished      bool        `json:"is_publised" db:"is_publised"`           // Indicates if the category is published and visible
	DisplayOrder     int         `json:"display_order" db:"display_order"`       // Order in which the category is displayed (lower values appear first)
	MetaDescription  *string     `json:"meta_description" db:"meta_description"` // SEO meta description for the category, used in search optimization
	ThumbnailURL     *string     `json:"thumbnail_url" db:"thumbnail_url"`       // URL of the image associated with the category
	Children         []*Category // List of child categories (used for hierarchical structures)
	Parents          []*Category // List of parent categories
	domain.BaseModel             // Base model containing common fields like created_at, updated_at.
}

type CategoryParent struct {
	ChildID  int `json:"child_id" db:"child_id"`
	ParentID int `json:"parent_id" db:"child_id"`
	Child    *Category
	Parent   *Category
	domain.BaseModel
}
