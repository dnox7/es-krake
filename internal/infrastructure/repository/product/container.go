package repository

import (
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/db"
)

type RepositoryContainer struct {
	AttributeRepository            domainRepo.AttributeRepository
	AttributeTypeRepository        domainRepo.AttributeTypeRepository
	BrandRepository                domainRepo.BrandRepository
	CategoryRepository             domainRepo.CategoryRepository
	OptionAttributeValueRepository domainRepo.OptionAttributeValueRepository
	ProductOptionRepository        domainRepo.ProductOptionRepository
	ProductRepository              domainRepo.ProductRepository
}

func NewRepositoryContainer(pg *db.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		AttributeRepository:            NewAttributeRepository(pg),
		AttributeTypeRepository:        NewAttributeTypeRepository(pg),
		BrandRepository:                NewBrandRepository(pg),
		CategoryRepository:             NewCategoryRepository(pg),
		OptionAttributeValueRepository: NewOptionAttributeValueRepository(pg),
		ProductOptionRepository:        NewProductOptionRepository(pg),
		ProductRepository:              NewProductRepository(pg),
	}
}
