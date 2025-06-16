package usecase

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
)

func (au *AuthUsecase) GetOperationsByPermissionID(ctx context.Context, permissionID int) ([]entity.AccessOperation, error) {
	return au.deps.AccessOperationService.GetOperationsByPermissionID(ctx, permissionID)
}
