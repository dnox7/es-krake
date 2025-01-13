package log

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestInitialize(t *testing.T) {
	type expect struct {
		level    slog.Level
		ctxValue loggerCtxKey
	}
	tests := []struct {
		name    string
		isDebug bool
		expect  expect
	}{
		{"can initialize logger with debug mode", true, expect{slog.LevelDebug, loggerCtxKey{}}},
		{"can initialize logger with info mode", false, expect{slog.LevelInfo, loggerCtxKey{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			buf := &bytes.Buffer{}
			ctx := Initialize(context.Background(), buf, tt.isDebug, []string{})
			if level := slog.Default().Enabled(ctx, tt.expect.level); !level {
				t.Errorf("expected log level to be %v, got %v", tt.expect.level, level)
			}
			if _, ok := ctx.Value(tt.expect.ctxValue).(*sync.Map); !ok {
				t.Errorf("expected context value to be of type sync.Map")
			}
		})
	}
}

type msgOutput struct {
	Message        string
	Time           time.Time
	Level          slog.Level
	AdditionalInfo string
	TestValues     []string
}

func TestLog(t *testing.T) {
	tests := []struct {
		name    string
		isDebug bool
		logFunc func(ctx context.Context, msg string, args ...any)
		expect  msgOutput
	}{
		{"can log info message", false, func(ctx context.Context, msg string, args ...any) { Info(ctx, msg) }, msgOutput{"info message", time.Time{}, slog.LevelInfo, "request_id", nil}},
		{"can log debug message", true, func(ctx context.Context, msg string, args ...any) { Debug(ctx, msg) }, msgOutput{"debug message", time.Time{}, slog.LevelDebug, "request_id", nil}},
		{"can log warn message", false, func(ctx context.Context, msg string, args ...any) { Warn(ctx, msg) }, msgOutput{"warn message", time.Time{}, slog.LevelWarn, "request_id", nil}},
		{"can log error message", false, func(ctx context.Context, msg string, args ...any) { Error(ctx, msg) }, msgOutput{"error message", time.Time{}, slog.LevelError, "request_id", nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			buf := &bytes.Buffer{}
			var msgOutput msgOutput
			ctx := Initialize(context.Background(), buf, tt.isDebug, []string{})

			key := reflect.ValueOf(&tt.expect).Elem().Type().Field(3).Name
			ctx = CtxWithValue(ctx, key, tt.expect.AdditionalInfo)
			tt.logFunc(ctx, tt.expect.Message)

			if err := json.Unmarshal(buf.Bytes(), &msgOutput); err != nil {
				t.Errorf("failed to unmarshal log output: %v", err)
			}
			if msgOutput.Message != tt.expect.Message {
				t.Errorf("expected message to be %s, got %s", tt.expect.Message, msgOutput.Message)
			}
			if msgOutput.Level != tt.expect.Level {
				t.Errorf("expected level to be %v, got %v", tt.expect.Level, msgOutput.Level)
			}
			if msgOutput.AdditionalInfo != tt.expect.AdditionalInfo {
				t.Errorf("expected value to be %s, got %s", tt.expect.AdditionalInfo, msgOutput.AdditionalInfo)
			}
		})
	}
}

func TestWithValue(t *testing.T) {
	type expect struct {
		key string
		val any
	}
	tests := []struct {
		name   string
		key    string
		val    any
		expect expect
	}{
		{"can add key-value pair to context", "key", "value", expect{"key", "value"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			ctx = CtxWithValue(ctx, tt.key, tt.val)
			if v, ok := ctx.Value(logMapCtxKey).(*sync.Map); ok {
				v.Range(func(key, value any) bool {
					if key, ok := key.(string); ok {
						if key != tt.expect.key {
							t.Errorf("expected key to be %s, got %s", tt.expect.key, key)
						}
						if value != tt.expect.val {
							t.Errorf("expected value to be %s, got %s", tt.expect.val, value)
						}
					}
					return true
				})
			}
		})
	}
}

type groupMsgOutput struct {
	Message string
	Time    time.Time
	Level   slog.Level
	key1    string
	values  []any
}

func TestGroup(t *testing.T) {
	tests := []struct {
		name   string
		expect groupMsgOutput
	}{
		{"key-value group logs nested key value pair", groupMsgOutput{"info message", time.Time{}, slog.LevelInfo, "key1", []any{"value1", "value2"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			buf := &bytes.Buffer{}
			var msgOutput map[string]any
			ctx := Initialize(context.Background(), buf, true, []string{})
			Info(ctx, tt.expect.Message, Group(tt.expect.key1, tt.expect.values...))

			if err := json.Unmarshal(buf.Bytes(), &msgOutput); err != nil {
				t.Errorf("failed to unmarshal log output: %v", err)
			}
			if msgOutput["message"] != tt.expect.Message {
				t.Errorf("expected message to be %s, got %s", tt.expect.Message, msgOutput["Message"])
			}
			if msgOutput["level"] != tt.expect.Level.String() {
				t.Errorf("expected level to be %v, got %v", tt.expect.Level, msgOutput["level"])
			}
			if !reflect.DeepEqual(msgOutput[tt.expect.key1], map[string]any{tt.expect.values[0].(string): tt.expect.values[1]}) {
				t.Errorf("expected value to be %v, got %v", tt.expect.values, msgOutput[tt.expect.key1])
			}
		})
	}
}
func TestWith(t *testing.T) {
	tests := []struct {
		name    string
		logFunc func(l *Logger, ctx context.Context, msg string, args ...any)
		expect  msgOutput
	}{
		{"can add key-value pair to info logger", func(l *Logger, ctx context.Context, msg string, args ...any) { l.Info(ctx, msg, args...) }, msgOutput{"info message", time.Time{}, slog.LevelInfo, "request_id", []string{"value1", "value2"}}},
		{"can add key-value pair to debug logger", func(l *Logger, ctx context.Context, msg string, args ...any) { l.Debug(ctx, msg, args...) }, msgOutput{"debug message", time.Time{}, slog.LevelDebug, "request_id", []string{"value1", "value2"}}},
		{"can add key-value pair to warn logger", func(l *Logger, ctx context.Context, msg string, args ...any) { l.Warn(ctx, msg, args...) }, msgOutput{"warn message", time.Time{}, slog.LevelWarn, "request_id", []string{"value1", "value2"}}},
		{"can add key-value pair to error logger", func(l *Logger, ctx context.Context, msg string, args ...any) { l.Error(ctx, msg, args...) }, msgOutput{"error message", time.Time{}, slog.LevelError, "request_id", []string{"value1", "value2"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			buf := &bytes.Buffer{}
			var msgOutput msgOutput
			ctx := Initialize(context.Background(), buf, true, []string{})

			key := reflect.ValueOf(&tt.expect).Elem().Type().Field(3).Name
			logger := With(ctx, key, tt.expect.AdditionalInfo)
			key2 := reflect.ValueOf(&tt.expect).Elem().Type().Field(4).Name
			l := logger.With(key2, tt.expect.TestValues)
			tt.logFunc(l, ctx, tt.expect.Message)
			if err := json.Unmarshal(buf.Bytes(), &msgOutput); err != nil {
				t.Errorf("failed to unmarshal log output: %v", err)
			}
			if msgOutput.Message != tt.expect.Message {
				t.Errorf("expected message to be %s, got %s", tt.expect.Message, msgOutput.Message)
			}
			if msgOutput.Level != tt.expect.Level {
				t.Errorf("expected level to be %v, got %v", tt.expect.Level, msgOutput.Level)
			}
			if msgOutput.AdditionalInfo != tt.expect.AdditionalInfo {
				t.Errorf("expected value to be %s, got %s", tt.expect.AdditionalInfo, msgOutput.AdditionalInfo)
			}
			if !reflect.DeepEqual(msgOutput.TestValues, tt.expect.TestValues) {
				t.Errorf("expected value to be %v, got %v", tt.expect.TestValues, msgOutput.TestValues)
			}
		})
	}
}

func TestHandle(t *testing.T) {
	tests := []struct {
		name   string
		expect msgOutput
	}{
		{"ctx can pass value in function", msgOutput{"info message", time.Time{}, slog.LevelInfo, "request_id", []string{"value1", "value2"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			buf := &bytes.Buffer{}
			output1, output2 := msgOutput{}, msgOutput{}
			ctx := Initialize(context.Background(), buf, true, []string{})

			key := reflect.ValueOf(&tt.expect).Elem().Type().Field(3).Name
			ctx = CtxWithValue(ctx, key, tt.expect.AdditionalInfo)
			Info(ctx, tt.expect.Message)
			func(ctx context.Context) {
				key2 := reflect.ValueOf(&tt.expect).Elem().Type().Field(4).Name
				ctx = CtxWithValue(ctx, key2, tt.expect.TestValues)
				Info(ctx, "test message")
			}(ctx)

			result := strings.Split(buf.String(), "\n")
			if err := json.Unmarshal([]byte(result[0]), &output1); err != nil {
				t.Errorf("failed to unmarshal log output: %v", err)
			}
			if output1.AdditionalInfo != tt.expect.AdditionalInfo {
				t.Errorf("expected value to be %s, got %s", tt.expect.AdditionalInfo, output1.AdditionalInfo)
			}
			if output1.TestValues != nil {
				t.Errorf("expected value to be nil, got %v", output1.TestValues)
			}
			if err := json.Unmarshal([]byte(result[1]), &output2); err != nil {
				t.Errorf("failed to unmarshal log output: %v", err)
			}
			if output2.AdditionalInfo != tt.expect.AdditionalInfo {
				t.Errorf("expected value to be %s, got %s", tt.expect.AdditionalInfo, output2.AdditionalInfo)
			}
			if !reflect.DeepEqual(output2.TestValues, tt.expect.TestValues) {
				t.Errorf("expected value to be nil, got %v", output2.TestValues)
			}
		})
	}
}
