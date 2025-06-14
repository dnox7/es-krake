package aws

import (
	"context"
	"fmt"

	"github.com/aws/smithy-go/logging"
	"github.com/dpe27/es-krake/pkg/log"
)

type awsLogger struct {
	logger *log.Logger
}

func (l awsLogger) WithContext(ctx context.Context) logging.Logger {
	return &awsLogger{
		logger: l.logger.WithContext(ctx),
	}
}

func (l awsLogger) Logf(classification logging.Classification, format string, v ...interface{}) {
	if classification == logging.Warn {
		l.logger.Warn(context.Background(), fmt.Sprintf(format, v...))
	}
	if classification == logging.Debug {
		l.logger.Debug(context.Background(), fmt.Sprintf(format, v...))
	}
}
