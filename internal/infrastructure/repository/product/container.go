package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/product/repository"
	mdb "github.com/dpe27/es-krake/internal/infrastructure/mongodb"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	AttributeRepo            domainRepo.AttributeRepository
	BrandRepo                domainRepo.BrandRepository
	CategoryRepo             domainRepo.CategoryRepository
	OptionAttributeValueRepo domainRepo.OptionAttributeValueRepository
	ProductOptionRepo        domainRepo.ProductOptionRepository
	ProductRepo              domainRepo.ProductRepository
	ProductCategoryRepo      domainRepo.ProductCategoryRepository
	ProductMetaRepo          domainRepo.ProductMetaRespository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL, mongo *mdb.Mongo) RepositoryContainer {
	return RepositoryContainer{
		AttributeRepo:            NewAttributeRepository(pg),
		BrandRepo:                NewBrandRepository(pg),
		CategoryRepo:             NewCategoryRepository(pg),
		OptionAttributeValueRepo: NewOptionAttributeValueRepository(pg),
		ProductOptionRepo:        NewProductOptionRepository(pg),
		ProductRepo:              NewProductRepository(pg),
		ProductCategoryRepo:      NewProductCategoryRepository(pg),
		ProductMetaRepo:          NewProductMetaRepository(mongo),
	}
}
