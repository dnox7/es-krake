package repository

import (
	"pech/es-krake/internal/infrastructure/database"
	"pech/es-krake/pkg/log"
)

type CategoryRepository struct {
	logger *log.Logger
	pg     *database.PostgreSQL
}
