package service

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
)

type AccessOperationService interface {
	GetOperationsWithAccessReqCode(ctx context.Context, code string) ([]entity.AccessOperation, error)

	HasRequiredOperations(perms []entity.Permission, requiredOps []entity.AccessOperation) (bool, error)
}
