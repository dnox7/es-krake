package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"github.com/dpe27/es-krake/internal/interfaces/batch/dto"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/wraperror"
	"github.com/gin-gonic/gin"
)

func (h *BatchHander) Batch(c *gin.Context) {
	event := dto.BatchEvent{}
	if err := c.ShouldBindHeader(&event); err != nil {
		nethttp.SetGenericErrorResponse(c, err, h.debug)
		return
	}

	c.Set("eventID", event.ID)
	c.Set("eventType", event.Type)
	h.logger.Info(c, fmt.Sprintf("start run batch %s", event.Type))

	if c.GetHeader(nethttp.HeaderContentType) != nethttp.MIMEApplicationJSON {
		nethttp.SetBadRequestResponse(
			c,
			"Expected 'application/json' for Content-Type, got '"+c.GetHeader(nethttp.HeaderContentType)+"'",
			nil, nil,
		)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		nethttp.SetGenericErrorResponse(c, err, h.debug)
		return
	}

	err = json.Unmarshal(body, &event.Data)
	if err != nil {
		nethttp.SetGenericErrorResponse(c, err, h.debug)
		return
	}

	batch, ok := h.batches[event.Type]
	if !ok {
		msg := fmt.Sprintf(
			"The event type '%v' has no defined handler on this server.",
			event.Type,
		)
		h.logger.Error(c, msg)
		nethttp.SetNotFoundResponse(c, msg, nil, nil)
	}

	batchLog, err := h.usecases.LogBatchStarting(c, event)
	if err != nil {
		h.logger.Error(c, "Error while logging batch start", "type", event.Type, "error", err)
		nethttp.SetGenericErrorResponse(c, err, h.debug)
		return
	}
	event.BatchLogID = batchLog.ID

	var batchErr error
	func() {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				log.Error(c, "Panic occurred", "stack_trace", fmt.Sprintf("%v\n%v", panicErr, string(debug.Stack())))
				batchErr = fmt.Errorf("%v\n%v", panicErr, string(debug.Stack()))
			}
		}()
		batchErr = batch.Run(c, event)
	}()

	batchLog, loggingErr := h.usecases.LogBatchEnded(c, batchLog, batchErr)
	if loggingErr != nil {
		h.logger.Error(
			c, fmt.Sprintf(
				"Error while saving the batch result in the batch_logs table for ID = %v",
				batchLog.ID,
			),
			"error", loggingErr,
		)
	}

	if batchErr != nil {
		var returnStatus int
		if batchError := (&wraperror.BatchError{}); errors.As(batchErr, &batchError) && !batchError.Retry {
			returnStatus = http.StatusNoContent

			// slackErr := h.slackService.NotifyBatchError("Error in batch", batchErr, event, body)
			// if slackErr != nil {
			// 	// This should not return an error code to knative since
			// 	// we don't want to retry the batch in this case
			// 	log.Error(c, "Error while notifying the error to Slack", "error", slackErr)
			// }
		} else {
			returnStatus = http.StatusInternalServerError
			c.Header("X-Event-Id", event.ID)
			c.Header("X-Event-Type", event.Type)
		}

		h.logger.Error(c, fmt.Sprintf("Error while running batch %v", event.Type), "error", batchErr)
		nethttp.SetErrorResponseWithStatus(c, returnStatus,
			"Error while processing the batch", batchErr.Error(), map[string]interface{}{
				"input": event,
			},
		)
		return
	}
	nethttp.SetNoContentResponse(c)
}
