package handler

import (
	"fmt"
	"io"

	"github.com/dpe27/es-krake/internal/interfaces/batch/dto"
	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/gin-gonic/gin"
)

func (h *BatchHander) DeadLetter(c *gin.Context) {
	event := dto.BatchEvent{}
	if err := c.ShouldBindHeader(&event); err != nil {
		h.logger.Error(c, "Error while parsing the dead-letter event body", "error", err)
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Warn(c, "failed to read request body", "error", err)
	}

	prevErrCode := c.Request.Header.Get("X-Event-Error-Code")
	prevErrData := c.Request.Header.Get("X-Event-Error-Data")

	h.logger.
		With("event", fmt.Sprintf("%#v", event)).
		With("previousErrorCode", prevErrCode).
		With("previousErrorData", prevErrData).
		With("body", string(body)).
		Error(c, "Event ended-up as a dead-letter")

		// messageToSlack := fmt.Errorf(
		// 	"the batch failed too many times and was terminated.\nPrevious response code: %v\nPrevious response body:\n%v",
		// 	previousErrorCode,
		// 	previousErrorData,
		// )
		// slackErr := h.slackService.NotifyBatchError("Dead-letter batch", messageToSlack, event, body)
		// if slackErr != nil {
		// 	// This should not return an error code to knative since
		// 	// we don't want to retry the batch in this case
		// 	log.Error(c, "Error while notifying the error to Slack", "error", slackErr)
		// }
	nethttp.SetNoContentResponse(c)
}
