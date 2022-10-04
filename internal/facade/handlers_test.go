package facade

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang-test-task/internal/database"
	"golang-test-task/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http/httptest"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"
)

// AnyTime is mock for time
type AnyTime struct{}

// Match for AnyTime mocking
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type HandlerFacadeTestSuite struct {
	suite.Suite
	mock   *sqlmock.Sqlmock
	db     *sql.DB
	gormDB *gorm.DB
	v      *validator.Validate
	client *database.Client
	logger *zap.Logger
	logic  *HandlerFacade
}

func (suite *HandlerFacadeTestSuite) SetupSuite() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.mock = &mock
	suite.db = db
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 suite.db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when init gorm", err)
	}
	if gormDB == nil {
		suite.T().Fatalf("gormDB is null")
	}
	suite.gormDB = gormDB

	v := validator.New()
	_ = v.RegisterValidation("checkURL", func(fl validator.FieldLevel) bool {
		arr, ok := fl.Field().Interface().([]string)
		if !ok {
			return false
		}
		for _, a := range arr {
			_, err := url.ParseRequestURI(a)
			if err != nil {
				return false
			}
		}
		return true
	})
	suite.v = v

	suite.client = database.NewClient(suite.gormDB)

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	suite.logger = logger
}

func (suite *HandlerFacadeTestSuite) TearDownSuite() {
	_ = suite.db.Close()
}

func (suite *HandlerFacadeTestSuite) SetupTest() {
	suite.logic = NewHandlerFacade(suite.client, suite.v, suite.logger)
}

func (suite *HandlerFacadeTestSuite) TestCreateAd() {
	id := 0
	expectedResult := entities.CreateAdAnswer{Status: "success", ID: &id}
	item := entities.AdJSONItem{Title: "title", Price: decimal.NewFromInt32(0)}
	mockedRow := sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(id))
	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "ad_items" ("created_at","updated_at","deleted_at","title","description","price")
			VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id","id"`)).
		WithArgs(AnyTime{}, AnyTime{}, nil, item.Title, item.Description, item.Price).
		WillReturnRows(mockedRow)
	(*suite.mock).ExpectCommit()
	(*suite.mock).ExpectClose()

	defer func() {
		_ = suite.logger.Sync()
	}()

	h, _ := suite.logic.GetHandler("create_ad")
	//svr := httptest.NewServer(http.HandlerFunc(h))
	//defer svr.Close()
	//
	bs, err := json.Marshal(item)
	if err != nil {
		suite.T().Fatalf("err during Marshal item; item = %v ; err = %v", item, err)
	}
	req := httptest.NewRequest("POST", "/ads", bytes.NewBuffer(bs))
	res := httptest.NewRecorder()

	h(res, req)
	resp := res.Result()
	body, _ := io.ReadAll(resp.Body)
	var actualResult entities.CreateAdAnswer
	_ = json.Unmarshal(body, &actualResult)

	if !reflect.DeepEqual(actualResult, expectedResult) {
		suite.T().Fatalf("actualResult != expectedResult; actualResult = %v ; expectedResult = %v",
			actualResult, expectedResult)
	}
}

func TestHandlerFacadeTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerFacadeTestSuite))
}
