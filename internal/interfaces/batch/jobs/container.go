package jobs

import "github.com/dpe27/es-krake/internal/usecase"

type BatchContainer map[string]Batch

func NewBatchContainer(
	usecases *usecase.UsecasesContainer,
) BatchContainer {
	return BatchContainer{
		"esk.batch.simple": newSimpleBatch(),
	}
}
