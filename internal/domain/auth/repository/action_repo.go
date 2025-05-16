package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type ActionType string

const (
	ActionTableName = "actions"

	ActionCreate ActionType = "create"
	ActionRead   ActionType = "read"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"

	ActionActivate   ActionType = "activate"
	ActionDeactivate ActionType = "deactivate"

	ActionUpload   ActionType = "upload"
	ActionDownload ActionType = "download"

	ActionExport ActionType = "export"
	ActionImport ActionType = "import"

	ActionView    ActionType = "view"
	ActionSubmit  ActionType = "submit"
	ActionReject  ActionType = "reject"
	ActionCancel  ActionType = "cancel"
	ActionRestore ActionType = "restore"
)

type ActionRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.Action, error)
	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.Action, error)
}
