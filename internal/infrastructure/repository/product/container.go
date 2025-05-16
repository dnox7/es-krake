package repository

import (
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	AttributeRepo            domainRepo.AttributeRepository
	AttributeTypeRepo        domainRepo.AttributeTypeRepository
	BrandRepo                domainRepo.BrandRepository
	CategoryRepo             domainRepo.CategoryRepository
	OptionAttributeValueRepo domainRepo.OptionAttributeValueRepository
	ProductOptionRepo        domainRepo.ProductOptionRepository
	ProductRepo              domainRepo.ProductRepository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		AttributeRepo:            NewAttributeRepository(pg),
		AttributeTypeRepo:        NewAttributeTypeRepository(pg),
		BrandRepo:                NewBrandRepository(pg),
		CategoryRepo:             NewCategoryRepository(pg),
		OptionAttributeValueRepo: NewOptionAttributeValueRepository(pg),
		ProductOptionRepo:        NewProductOptionRepository(pg),
		ProductRepo:              NewProductRepository(pg),
	}
}
