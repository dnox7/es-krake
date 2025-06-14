package jobs

import (
	"context"

	"github.com/dpe27/es-krake/internal/interfaces/batch/dto"
)

type Batch interface {
	Run(ctx context.Context, event dto.BatchEvent) error
}
