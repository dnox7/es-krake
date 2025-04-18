package utils

import (
	"gorm.io/gorm"
)

type Scope func(*gorm.DB) *gorm.DB

func WhereScope(query string, args ...interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

func JoinScope(query string, args ...interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(query, args...)
	}
}

func OrderScope(query string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(query)
	}
}

func LimitScope(limit int) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func PaginateScope(pageData map[string]int) Scope {
	return func(db *gorm.DB) *gorm.DB {
		pageNo := DefaultPageNo
		if valPage, ok := pageData["pageNo"]; ok && valPage > 0 {
			pageNo = valPage
		}

		pageSize := DefaultPageSize
		if valPageSize, ok := pageData["pageSize"]; ok && valPageSize > 0 {
			pageSize = valPageSize
		}

		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func PreloadScope(query string, args ...interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, args...)
	}
}

func ToGormScope(scopes ...Scope) []func(*gorm.DB) *gorm.DB {
	gormScopes := make([]func(*gorm.DB) *gorm.DB, len(scopes))
	for i, s := range scopes {
		gormScopes[i] = s
	}
	return gormScopes
}
