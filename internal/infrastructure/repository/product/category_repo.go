package repository

import (
	"context"
	"pech/es-krake/internal/domain"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/pkg/log"
)

type categoryRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewCategoryRepository(pg *db.PostgreSQL) domainRepo.CategoryRepository {
	return &categoryRepository{
		logger: log.With("repo", "category_repo"),
		pg:     pg,
	}
}

// CreateWithTx implements repository.CategoryRepository.
func (c *categoryRepository) CreateWithTx(ctx context.Context, attributes map[string]interface{}) (entity.Category, error) {
	panic("unimplemented")
}

// FindByConditions implements repository.CategoryRepository.
func (c *categoryRepository) FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Category, error) {
	panic("unimplemented")
}

// FindByConditionsWithScope implements repository.CategoryRepository.
func (c *categoryRepository) FindByConditionsWithScope(ctx context.Context, conditions map[string]interface{}, scopes ...domain.Scope) ([]entity.Category, error) {
	panic("unimplemented")
}

// TakeByConditions implements repository.CategoryRepository.
func (c *categoryRepository) TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Category, error) {
	panic("unimplemented")
}

// TakeByID implements repository.CategoryRepository.
func (c *categoryRepository) TakeByID(ctx context.Context, ID int) (entity.Category, error) {
	panic("unimplemented")
}

// UpdateWithTx implements repository.CategoryRepository.
func (c *categoryRepository) UpdateWithTx(ctx context.Context, category entity.Category, attributesToUpdate map[string]interface{}) (entity.Category, error) {
	panic("unimplemented")
}
