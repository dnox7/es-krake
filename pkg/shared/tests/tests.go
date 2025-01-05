package tests

import (
	"database/sql"
	"database/sql/driver"
	"net/http/httptest"
	"pech/es-krake/internal/services/main-service/pkg/container"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MockedApplication struct {
	Test         *testing.T
	Logger       *logrus.Logger
	DB           *sql.DB
	Gorm         *gorm.DB
	GoMock       *gomock.Controller
	SqlMock      sqlmock.Sqlmock
	Gin          *gin.Engine
	Server       *httptest.Server
	Repositories *container.RepositoryContainers
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyInputTime struct{}

func (a AnyInputTime) Match(v driver.Value) bool {
    inputTime := v.(string)
    _,  err := 
}
