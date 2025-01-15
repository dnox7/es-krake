package entity

import "pech/es-krake/internal/domain/utils"

type Category struct {
	ID              int         // Unique identifier for the category
	ParentID        *int        // ID of the parent category (nil if it's a top-level category)
	Children        *[]Category // List of child categories (used for hierarchical structures)
	Name            string      // Name of the category
	Slug            string      // URL-friendly identifier for the category (used in routing)
	Description     string      // Detailed description of the category
	IsPublished     bool        // Indicates if the category is published and visible
	DisplayOrder    int         // Order in which the category is displayed (lower values appear first)
	MetaDescription *string     // SEO meta description for the category, used in search optimization
	ImageURL        string      // URL of the image associated with the category
	utils.BaseModel             // Base model containing common fields like created_at, updated_at.
}

func (c *Category) AddChild(children ...Category) {
	if c.Children == nil {
		c.Children = &[]Category{}
	}
	*c.Children = append(*c.Children, children...)
}
