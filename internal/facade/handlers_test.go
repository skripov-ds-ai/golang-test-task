package facade

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type HandlerFacadeTestSuite struct {
	suite.Suite
	mock *sqlmock.Sqlmock
}

func (suite *HandlerFacadeTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.mock = &mock
	defer func() {
		_ = db.Close()
	}()
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when init gorm", err)
	}
	if gormDB == nil {
		suite.T().Fatalf("gormDB is null")
	}
}
