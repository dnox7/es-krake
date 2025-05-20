package scope

import (
	"fmt"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/pkg/utils"

	"gorm.io/gorm"
)

type gormScope struct {
	scopes []func(*gorm.DB) *gorm.DB
}

func GormScope() *gormScope {
	return &gormScope{}
}

func (g *gormScope) GetSpec() interface{} {
	return g.scopes
}

func (g *gormScope) Where(query string, args ...interface{}) *gormScope {
	g.scopes = append(g.scopes, func(d *gorm.DB) *gorm.DB {
		return d.Where(query, args...)
	})
	return g
}

func (g *gormScope) Join(query string, args ...interface{}) *gormScope {
	g.scopes = append(g.scopes, func(d *gorm.DB) *gorm.DB {
		return d.Joins(query, args...)
	})
	return g
}

func (g *gormScope) Order(query string) *gormScope {
	g.scopes = append(g.scopes, func(d *gorm.DB) *gorm.DB {
		return d.Order(query)
	})
	return g
}

func (g *gormScope) Limit(limit int) *gormScope {
	g.scopes = append(g.scopes, func(d *gorm.DB) *gorm.DB {
		return d.Limit(limit)
	})
	return g
}

func (g *gormScope) Paginate(pageData map[string]int) *gormScope {
	g.scopes = append(g.scopes, func(d *gorm.DB) *gorm.DB {
		pageNo := utils.DefaultPageNo
		if valPage, ok := pageData["pageNo"]; ok && valPage > 0 {
			pageNo = valPage
		}

		pageSize := utils.DefaultPageSize
		if valPageSize, ok := pageData["pageSize"]; ok && valPageSize > 0 {
			pageSize = valPageSize
		}

		offset := (pageNo - 1) * pageSize
		return d.Offset(offset).Limit(pageSize)
	})
	return g
}

func (g *gormScope) Preload(query string, args ...interface{}) *gormScope {
	g.scopes = append(g.scopes, func(d *gorm.DB) *gorm.DB {
		return d.Preload(query, args...)
	})
	return g
}

func ToGormScopes(spec specification.Base) ([]func(*gorm.DB) *gorm.DB, error) {
	gormScopes, ok := spec.GetSpec().([]func(*gorm.DB) *gorm.DB)
	if !ok {
		return nil, fmt.Errorf(utils.ErrorGetSpec)
	}
	return gormScopes, nil
}
