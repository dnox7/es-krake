package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	domainScope "pech/es-krake/internal/domain/shared/scope"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/pkg/log"
)

type productRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewProductRepository(pg *db.PostgreSQL) domainRepo.ProductRepository {
	return &productRepository{
		logger: log.With("repo", "product_repo"),
		pg:     pg,
	}
}

// Create implements repository.ProductRepository.
func (p *productRepository) Create(ctx context.Context, attributes map[string]interface{}) (entity.Product, error) {
	panic("unimplemented")
}

// CreateWithTx implements repository.ProductRepository.
func (p *productRepository) CreateWithTx(ctx context.Context, attributes map[string]interface{}) (entity.Product, error) {
	panic("unimplemented")
}

// FindByConditions implements repository.ProductRepository.
func (p *productRepository) FindByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...domainScope.Base) ([]entity.Product, error) {
	panic("unimplemented")
}

// TakeByConditions implements repository.ProductRepository.
func (p *productRepository) TakeByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...domainScope.Base) (entity.Product, error) {
	panic("unimplemented")
}

// UpdateWithTx implements repository.ProductRepository.
func (p *productRepository) UpdateWithTx(ctx context.Context, product entity.Product, attributesToUpdate map[string]interface{}) (entity.Product, error) {
	panic("unimplemented")
}
