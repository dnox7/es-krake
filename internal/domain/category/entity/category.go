package entity

import "pech/es-krake/internal/domain/utils"

type Category struct {
	ID              int
	ParentID        *int
	Children        *[]Category
	Name            string
	Slug            string
	Description     string
	IsPublished     bool
	DisplayOrder    int
	MetaDescription *string
	ImageURL        string
	Attributes      *[]CategoryAttribute
	utils.BaseModel
}

func (c *Category) AddChild(children ...Category) {
	if c.Children == nil {
		c.Children = &[]Category{}
	}
	*c.Children = append(*c.Children, children...)
}

func (c *Category) AddAttribute(attr ...CategoryAttribute) {
	if c.Attributes == nil {
		c.Attributes = &[]CategoryAttribute{}
	}
	*c.Attributes = append(*c.Attributes, attr...)
}
