package repository

import (
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/db"
)

type RepositoryContainer struct {
	AttributeRepository             domainRepo.IAttributeRepository
	AttributeTypeRepository         domainRepo.IAttributeTypeRepository
	CategoryRepository              domainRepo.ICategoryRepository
	ProductAttributeValueRepository domainRepo.IProductAttributeValueRepository
	ProductOptionRepository         domainRepo.IProductOptionRepository
	ProductRepository               domainRepo.IProductRepository
}

func NewRepositoryContainer(pg *db.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		AttributeRepository:             NewAttributeRepository(pg),
		AttributeTypeRepository:         NewAttributeTypeRepository(pg),
		CategoryRepository:              NewCategoryRepository(pg),
		ProductAttributeValueRepository: NewProductAttributeValueRepository(pg),
		ProductOptionRepository:         NewProductOptionRepository(pg),
		ProductRepository:               NewProductRepository(pg),
	}
}
