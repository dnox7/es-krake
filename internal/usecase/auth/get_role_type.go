package usecase

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/pkg/log"
)

func (au AuthUsecase) GetAllRoleTypes(ctx context.Context) ([]entity.RoleType, error) {
	return au.deps.RoleTypeRepo.FindByConditions(ctx, nil, nil)
}

func (au AuthUsecase) GetRoleTypeByID(ctx context.Context, ID int) (entity.RoleType, error) {
	log.Debug(ctx, "into usecase")
	return au.deps.RoleTypeRepo.TakeByConditions(ctx, map[string]interface{}{
		"id": ID,
	}, nil)
}
