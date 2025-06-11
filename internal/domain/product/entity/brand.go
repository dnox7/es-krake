package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const BrandTableName = "brands"

type Brand struct {
	ID            int     `gorm:"column:id;type:bigint;primaryKey;not null;autoIncrement" json:"id"`
	Name          string  `gorm:"column:name;type:varchar(50);not null"                   json:"name"`
	IsActive      bool    `gorm:"column:is_active;type:bool;not null;default:true"        json:"is_active"`
	Description   *string `gorm:"column:description;type:text"                            json:"description"`
	ThumbnailPath *string `gorm:"column:thumbnail_path;type:varchar(255)"                 json:"thumbnail_path"`
	WebsitePath   *string `gorm:"column:website_path;type:varchar(255)"                   json:"website_path"`
	model.BaseModel
}

func (Brand) TableName() string {
	return BrandTableName
}
