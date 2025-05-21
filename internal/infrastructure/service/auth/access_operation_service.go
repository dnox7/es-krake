package service

import (
	"context"
	"fmt"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/auth/repository"
	domainService "github.com/dpe27/es-krake/internal/domain/auth/service"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
)

type accessOperationService struct {
	accessOpRepo repository.AccessOperationRepository
	logger       *log.Logger
}

func NewAccessOperationService(
	accessOpRepo repository.AccessOperationRepository,
) domainService.AccessOperationService {
	return &accessOperationService{
		accessOpRepo,
		log.With("service", "access_operation_service"),
	}
}

// GetOperationsWithAccessReqCode implements service.AccessOperationService.
func (a *accessOperationService) GetOperationsWithAccessReqCode(
	ctx context.Context,
	code string,
) ([]entity.AccessOperation, error) {
	scopes := scope.GormScope().
		Join(fmt.Sprintf(
			"INNER JOIN %s AS aro ON aro.access_operation_id = %s.id",
			repository.AccessRequirementOperationTableName,
			repository.AccessOperationsTableName,
		)).
		Join(fmt.Sprintf(
			"INNER JOIN %s AS ar ON ar.id = aro.access_requirement_id",
			repository.AccessRequirementTableName,
		)).
		Where("ar.code = ?", code).
		Preload("Action").
		Preload("Resource")
	return a.accessOpRepo.FindByConditions(ctx, nil, scopes)
}
