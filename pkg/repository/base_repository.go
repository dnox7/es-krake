package repository

import (
	"pech/es-krake/pkg/shared/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BaseRepository struct {
	Logger *logrus.Logger
}

func NewBaseRepository(logger *logrus.Logger) *BaseRepository {
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
