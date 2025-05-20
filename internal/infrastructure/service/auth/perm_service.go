package service

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/auth/repository"
	domainService "github.com/dpe27/es-krake/internal/domain/auth/service"
	"github.com/dpe27/es-krake/pkg/log"
)

type permService struct {
	permRepo repository.PermissionRepository
	logger   *log.Logger
}

func NewPermissionService(
	permRepo repository.PermissionRepository,
) domainService.PermissionService {
	return &permService{
		permRepo,
		log.With("service", "permission_serivce"),
	}
}

// GetPermissionsWithKcUserID implements service.PermissionService.
func (p *permService) GetPermissionsWithKcUserID(
	ctx context.Context,
	kcUserID int,
) ([]entity.Permission, error) {
	panic("unimplemented")
}
