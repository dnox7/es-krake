package scope

import (
	"fmt"
	"pech/es-krake/internal/domain/shared/scope"
	"pech/es-krake/pkg/utils"

	"gorm.io/gorm"
)

type gormScope struct {
	sc func(*gorm.DB) *gorm.DB
}

func (g *gormScope) GetScope() interface{} {
	return g.sc
}

func Where(query string, args ...interface{}) scope.Base {
	return &gormScope{sc: func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}}
}

func Join(query string, args ...interface{}) scope.Base {
	return &gormScope{sc: func(db *gorm.DB) *gorm.DB {
		return db.Joins(query, args...)
	}}
}

func Order(query string) scope.Base {
	return &gormScope{sc: func(db *gorm.DB) *gorm.DB {
		return db.Order(query)
	}}
}

func Limit(limit int) scope.Base {
	return &gormScope{sc: func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}}
}

func Paginate(pageData map[string]int) scope.Base {
	return &gormScope{sc: func(db *gorm.DB) *gorm.DB {
		pageNo := utils.DefaultPageNo
		if valPage, ok := pageData["pageNo"]; ok && valPage > 0 {
			pageNo = valPage
		}

		pageSize := utils.DefaultPageSize
		if valPageSize, ok := pageData["pageSize"]; ok && valPageSize > 0 {
			pageSize = valPageSize
		}

		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}}
}

func Preload(query string, args ...interface{}) scope.Base {
	return &gormScope{sc: func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, args...)
	}}
}

func ToGormScopes(scopes ...scope.Base) ([]func(*gorm.DB) *gorm.DB, error) {
	gormScopes := make([]func(*gorm.DB) *gorm.DB, 0)
	for _, s := range scopes {
		gs, ok := s.GetScope().(func(*gorm.DB) *gorm.DB)
		if !ok {
			return nil, fmt.Errorf(utils.ErrorGetTx)
		}
		gormScopes = append(gormScopes, gs)
	}
	return gormScopes, nil
}
