package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

type BaseLogger interface {
	logrus.FieldLogger
	WithContext(ctx context.Context) *logrus.Entry
}
