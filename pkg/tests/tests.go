package tests

//
// import (
// 	"bytes"
// 	"database/sql"
// 	"database/sql/driver"
// 	"encoding/json"
// 	"io"
// 	"log/slog"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/http/httputil"
// 	"os"
// 	"pech/es-krake/internal/services/main-service/pkg/container"
// 	"pech/es-krake/pkg/infrastructure"
// 	"pech/es-krake/pkg/log"
// 	"pech/es-krake/pkg/utils"
// 	"strings"
// 	"testing"
// 	"time"
//
// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// 	"github.com/h2non/gock"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// 	gormLog "gorm.io/gorm/logger"
// )
//
// type MockedApplication struct {
// 	Test         *testing.T
// 	Logger       *log.Logger
// 	DB           *sql.DB
// 	Gorm         *gorm.DB
// 	GoMock       *gomock.Controller
// 	SqlMock      sqlmock.Sqlmock
// 	Gin          *gin.Engine
// 	Server       *httptest.Server
// 	Repositories *container.RepositoryContainers
// }
//
// type AnyTime struct{}
//
// func (a AnyTime) Match(v driver.Value) bool {
// 	_, ok := v.(time.Time)
// 	return ok
// }
//
// type AnyInputTime struct{}
//
// func (a AnyInputTime) Match(v driver.Value) bool {
// 	inputTime := v.(string)
// 	_, err := utils.ParseDateTimeFromSQlOrISO(&inputTime)
// 	return err == nil
// }
//
// type AnyString struct{}
//
// func (a AnyString) Match(v driver.Value) bool {
// 	_, ok := v.(string)
// 	return ok
// }
//
// type serviceName string
//
// const (
// 	MainService  serviceName = "main-service"
// 	BatchService serviceName = "batch-service"
// )
//
// type gormSlogWritter struct {
// 	logger *log.Logger
// }
//
// func (w *gormSlogWritter) Printf(fmt string, args ...interface{}) {
// 	// impl func Printf
// }
//
// func NewMockedApplication(test *testing.T, service serviceName) (*MockedApplication, error) {
// 	if err := os.Setenv("PE_SERVICE_NAME", string(service)); err != nil {
// 		return nil, err
// 	}
//
// 	logger := log.With()
// 	db, dbMock, err := sqlmock.New()
// 	if err != nil {
// 		return nil, err
// 	}
// 	dbMock.MatchExpectationsInOrder(false)
//
// 	gormConf := infrastructure.GetGormConfig()
// 	gormConf.Logger = gormLog.New(
// 		&gormSlogWritter{logger: logger},
// 		gormLog.Config{
// 			LogLevel: gormLog.Warn,
// 			Colorful: true,
// 		},
// 	)
// 	gormConf.SkipDefaultTransaction = true
//
// 	gorm, err := gorm.Open(
// 		postgres.New(postgres.Config{
// 			Conn: db,
// 		}),
// 		gormConf,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	gin.DefaultWriter = io.Discard
// 	gin.DefaultErrorWriter = os.Stderr
//
// 	ginServer := infrastructure.NewServer(gorm)
// 	server := httptest.NewServer(ginServer)
// 	gock.New(server.URL).EnableNetworking().Persist()
//
// 	mockController := gomock.NewController(test)
// 	repositories, err := container.NewRepositoryContainers(gorm)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &MockedApplication{
// 		Test:         test,
// 		Logger:       logger,
// 		DB:           db,
// 		Gorm:         gorm,
// 		GoMock:       mockController,
// 		SqlMock:      dbMock,
// 		Gin:          ginServer,
// 		Server:       server,
// 		Repositories: repositories,
// 	}, nil
// }
//
// func NewMockedApplicationOrFail(test *testing.T, service serviceName) *MockedApplication {
// 	app, err := NewMockedApplication(test, service)
// 	if err != nil {
// 		app.Test.Fatalf("Error while mocking the application: '%s'", err)
// 	}
// 	return app
// }
//
// func (app *MockedApplication) Close() error {
// 	app.Server.Close()
// 	gock.DisableNetworking()
// 	app.GoMock.Finish()
// 	return app.DB.Close()
// }
//
// func (app *MockedApplication) AssertMockExpectationsFullfilled() {
// 	if err := app.SqlMock.ExpectationsWereMet(); err != nil {
// 		app.Test.Errorf("Unfullfilled expectations about the mocked database: %s", err.Error())
// 	}
// }
//
// func (app *MockedApplication) AssertResponseStatus(
// 	response *http.Response,
// 	expectedStatus int,
// ) {
// 	if response.StatusCode != expectedStatus {
// 		dump, _ := httputil.DumpResponse(response, true)
// 		app.Test.Fatalf("Expected status code %v, got %v. Response: %s", expectedStatus, response.StatusCode, dump)
// 	}
// }
//
// func (app *MockedApplication) GetBodyOrFail(response *http.Response) []byte {
// 	body, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		dump, _ := httputil.DumpResponse(response, true)
// 		app.Test.Fatalf("Error while reading response body: %s, %v", dump, err)
// 	}
// 	return body
// }
//
// func MakeSingleSqlMockRow(rows map[string]driver.Value) *sqlmock.Rows {
// 	return utils.CreateSingleMockRow(rows)
// }
//
// type MockFileInfo struct {
// 	FileName string
// 	Data     []byte
// }
//
// func (mf MockFileInfo) Name() string       { return mf.FileName }
// func (mf MockFileInfo) Size() int64        { return int64(len(mf.Data)) }
// func (mf MockFileInfo) Mode() os.FileMode  { return 0444 }
// func (mf MockFileInfo) ModTime() time.Time { return time.Time{} }
// func (mf MockFileInfo) IsDir() bool        { return false }
// func (mf MockFileInfo) Sys() interface{}   { return nil }
//
// func ParseLogEntries(buf *bytes.Buffer) []map[string]interface{} {
// 	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
// 	var entries []map[string]interface{}
// 	for _, line := range lines {
// 		if line == "" {
// 			continue
// 		}
// 		var entry map[string]interface{}
// 		if err := json.Unmarshal([]byte(line), &entry); err == nil {
// 			entries = append(entries, entry)
// 		}
// 	}
// 	return entries
// }
//
// func NewMockSlog() *bytes.Buffer {
// 	logBuffer := &bytes.Buffer{}
// 	options := &slog.HandlerOptions{
// 		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
// 			if v, ok := a.Value.Any().(time.Duration); ok {
// 				a.Value = slog.StringValue(v.String())
// 			}
// 			return a
// 		},
// 	}
// 	newLogger := slog.New(slog.NewJSONHandler(logBuffer, options))
// 	slog.SetDefault(newLogger)
// 	return logBuffer
// }
