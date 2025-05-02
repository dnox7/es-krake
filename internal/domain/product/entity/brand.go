package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

type Brand struct {
	ID                 int     `gorm:"column:id;type:bigint;primaryKey;not null;autoIncrement" json:"id"`
	Name               string  `gorm:"column:name;type:varchar(100);not null"                  json:"name"`
	IsActive           bool    `gorm:"column:is_active;type:bool;not null;default:true"        json:"is_active"`
	Description        *string `gorm:"column:description;type:text"                            json:"description"`
	ThumbnailImagePath *string `gorm:"column:thumbnail_image_path;type:varchar(255)"           json:"thumbnail_image_path"`
	WebsitePath        *string `gorm:"column:website_path;type:varchar(255)"                   json:"website_path"`
	model.BaseModel
}
