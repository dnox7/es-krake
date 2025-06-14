package handler

import (
	"github.com/dpe27/es-krake/internal/infrastructure/notify"
	"github.com/dpe27/es-krake/internal/interfaces/batch/jobs"
	batchLogUC "github.com/dpe27/es-krake/internal/usecase/batchlog"
	"github.com/dpe27/es-krake/pkg/log"
)

type BatchHander struct {
	logger   *log.Logger
	batches  jobs.BatchContainer
	usecases batchLogUC.BatchLogUsecase
	notifier notify.DiscordNotifier
	debug    bool
}

func NewBatchHandler(
	batches jobs.BatchContainer,
	usecases batchLogUC.BatchLogUsecase,
	notifier notify.DiscordNotifier,
	debug bool,
) *BatchHander {
	return &BatchHander{
		logger:   log.With("object", "batch_http_handler"),
		batches:  batches,
		notifier: notifier,
		debug:    debug,
	}
}
