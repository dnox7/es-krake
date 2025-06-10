package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/product/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	AttributeRepo            domainRepo.AttributeRepository
	BrandRepo                domainRepo.BrandRepository
	CategoryRepo             domainRepo.CategoryRepository
	OptionAttributeValueRepo domainRepo.OptionAttributeValueRepository
	ProductOptionRepo        domainRepo.ProductOptionRepository
	ProductRepo              domainRepo.ProductRepository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		AttributeRepo:            NewAttributeRepository(pg),
		BrandRepo:                NewBrandRepository(pg),
		CategoryRepo:             NewCategoryRepository(pg),
		OptionAttributeValueRepo: NewOptionAttributeValueRepository(pg),
		ProductOptionRepo:        NewProductOptionRepository(pg),
		ProductRepo:              NewProductRepository(pg),
	}
}
