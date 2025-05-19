package service

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/auth/repository"
	domainService "github.com/dpe27/es-krake/internal/domain/auth/service"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"

	"golang.org/x/sync/errgroup"
)

type accessOperationService struct {
	accessReqRepo repository.AccessRequirementRepository
	logger        *log.Logger
}

func NewAccessOperationService(
	accessReqRepo repository.AccessRequirementRepository,
) domainService.AccessOperationService {
	return &accessOperationService{
		accessReqRepo,
		log.With("service", "access_operation_service"),
	}
}

// GetOperationsWithAccessReqCode implements service.AccessOperationService.
func (a *accessOperationService) GetOperationsWithAccessReqCode(
	ctx context.Context,
	code string,
) ([]entity.AccessOperation, error) {
	return nil, nil
}

// HasRequiredOperations implements service.AccessOperationService.
func (a *accessOperationService) HasRequiredOperations(
	perms []entity.Permission,
	requiredOps []entity.AccessOperation,
) (bool, error) {
	var (
		g         errgroup.Group
		permCodes []string
		opCodes   []string
	)

	g.Go(func() error {
		var err error
		permCodes, err = entity.MapPermissionsToCodes(perms)
		return err
	})

	g.Go(func() error {
		var err error
		opCodes, err = entity.MapOperationsToCodes(requiredOps)
		return err
	})

	if err := g.Wait(); err != nil {
		return false, err
	}

	return utils.IsSubSet(opCodes, permCodes), nil
}
