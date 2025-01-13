package repository

import (
	"log"
	"pech/es-krake/pkg/shared/utils"

	"gorm.io/gorm"
)

type BaseRepository struct {
	Logger *log.Logger
}

func NewBaseRepository(logger *log.Logger) *BaseRepository {
	return &BaseRepository{Logger: logger}
}

func Paginate(pageData map[string]int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageNo := utils.DefaultPageNo
		if val, ok := pageData["PageNo"]; ok && val > 0 {
			pageNo = val
		}

		pageSize := utils.DefaultPageSize
		if val, ok := pageData["PageSize"]; ok && val > 0 {
			pageSize = val
		}

		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
