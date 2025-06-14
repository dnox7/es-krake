package jobs

import (
	"context"

	"github.com/dpe27/es-krake/internal/interfaces/batch/dto"
)

type simpleBatch struct{}

func newSimpleBatch() *simpleBatch {
	return &simpleBatch{}
}

func (batch *simpleBatch) Run(ctx context.Context, event dto.BatchEvent) error {
	return nil
}
