package entity

import "pech/es-krake/internal/domain"

type Category struct {
	ID               int         `json:"id" db:"id"`                             // Unique identifier for the category
	ParentID         *int        `json:"parent_id" db:"parent_id"`               // ID of the parent category (nil if it's a top-level category)
	Children         *[]Category `json:"children"`                               // List of child categories (used for hierarchical structures)
	Name             string      `json:"name" db:"name"`                         // Name of the category
	Slug             string      `json:"slug" db:"slug"`                         // URL-friendly identifier for the category (used in routing)
	Description      string      `json:"description" db:"description"`           // Detailed description of the category
	IsPublished      bool        `json:"is_publised" db:"is_publised"`           // Indicates if the category is published and visible
	DisplayOrder     int         `json:"display_order" db:"display_order"`       // Order in which the category is displayed (lower values appear first)
	MetaDescription  *string     `json:"meta_description" db:"meta_description"` // SEO meta description for the category, used in search optimization
	ImageURL         string      `json:"image_url" db:"image_url"`               // URL of the image associated with the category
	domain.BaseModel             // Base model containing common fields like created_at, updated_at.
}

func (c *Category) AddChild(children ...Category) {
	if c.Children == nil {
		c.Children = &[]Category{}
	}
	*c.Children = append(*c.Children, children...)
}
