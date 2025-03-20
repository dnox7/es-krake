package repository

import (
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/pkg/log"
)

type categoryRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewCategoryRepository(pg *db.PostgreSQL) domainRepo.ICategoryRepository {
	return &cctegoryRepository{
		logger: log.With("repo", "category_repo"),
		pg:     pg,
	}
}
