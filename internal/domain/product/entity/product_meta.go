package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductMeta struct {
	ID                     primitive.ObjectID     `bson:"_id,omitempty"                      json:"id"`
	ProductID              int                    `bson:"product_id"                         json:"product_id"`
	DisplayName            string                 `bson:"display_name"                       json:"display_name"`
	DescriptionHTML        string                 `bson:"description_html"                   json:"description_html"`
	TechnicalSpec          map[string]interface{} `bson:"technical_spec,omitempty"           json:"technical_spec"`
	IllustrationImagePaths []string               `bson:"illustration_image_paths,omitempty" json:"illustration_image_paths"`
	CreatedAt              time.Time              `bson:"created_at,omitempty"               json:"created_at"`
	UpdatedAt              time.Time              `bson:"updated_at,omitempty"               json:"updated_at"`
}

func (p *ProductMeta) PrepareForInsert() {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *ProductMeta) PrepareForUpdate() {
	p.UpdatedAt = time.Now()
}
