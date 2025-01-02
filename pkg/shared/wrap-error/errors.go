package wraperror

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

type Fields map[string]interface{}

type DetailTrace interface {
	GetTraces() string
	GetTraceSource() string
	GetStackTrace() []string
	GetFields() Fields
}

type detailTrace struct {
	tracer     TraceHistory
	stackTrace []string
	fields     Fields
}

func (e *detailTrace) GetTraces() string {
	res := ""
	if e.tracer != nil {
		res = e.tracer.Traces()
	}
	return res
}

func (e *detailTrace) GetTraceSource() string {
	res := ""
	if e.tracer != nil {
		res = e.tracer.TraceSource()
	}
	return res
}

func (e *detailTrace) GetStackTrace() []string {
	return e.stackTrace
}

func (e *detailTrace) GetFields() Fields {
	return e.fields
}

type withTrace struct {
	error
	DetailTrace
}

type TraceHistory interface {
	Traces() string
	AddTraces(traces ...interface{})
	TraceSource() string
	SetTraceSource(source string)
}

type stack []uintptr

func WithTrace(err error, fields Fields, tracer TraceHistory) error {
	if err == nil {
		return nil
	}

	stackTrace := []string{}
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(2, pcs[:])
	var st stack = pcs[0:n]

	for _, pc := range st {
		f := errors.Frame(pc)
		fStr := fmt.Sprintf("%+v", f)
		fStr = strings.Replace(fStr, "\n\t", " in ", -1)
		stackTrace = append(stackTrace, fStr)
	}

	dt := &detailTrace{
		tracer:     tracer,
		fields:     fields,
		stackTrace: stackTrace,
	}

	return &withTrace{err, dt}
}
