package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/product/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
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

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
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
