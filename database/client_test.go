package database

import (
	"database/sql/driver"
	"golang-test-task/entities"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AnyTime is mock for time
type AnyTime struct{}

// Match for AnyTime mocking
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestClientGetAd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
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
		t.Fatalf("an error '%s' was not expected when init gorm", err)
	}
	if gormDB == nil {
		t.Fatalf("gormDB is null")
	}

	id := 0
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "ad_items"
			WHERE "ad_items"."id" = $1 AND "ad_items"."deleted_at" IS NULL
			ORDER BY "ad_items"."id" LIMIT 1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{}))

	client := NewClient(gormDB)
	_, err = client.GetAd(id)
	if err != gorm.ErrRecordNotFound {
		t.Fatalf("gorm.ErrRecordNotFound was expected; actual = %v", err)
	}
}

func TestClientListAds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
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
		t.Fatalf("an error '%s' was not expected when init gorm", err)
	}
	if gormDB == nil {
		t.Fatalf("gormDB is null")
	}

	mock.MatchExpectationsInOrder(true)

	//mock.ExpectQuery(regexp.QuoteMeta(
	//	`SELECT * FROM "image_url";`))

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT "ad_items"."id","ad_items"."title","ad_items"."price" 
			FROM "ad_items" WHERE "ad_items"."deleted_at" IS NULL ORDER BY id asc LIMIT 10`)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "image_url";`))

	client := NewClient(gormDB)
	items, err := client.ListAds(0, 10, "id", true)
	if err != nil {
		t.Fatalf("nil error was expected; actual = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("zero length of items was expected ; items = %v", items)
	}
}

func TestClientCreateAd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
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
		t.Fatalf("an error '%s' was not expected when init gorm", err)
	}
	if gormDB == nil {
		t.Fatalf("gormDB is null")
	}

	id := 0
	item := entities.AdJSONItem{}
	mockedRow := sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(id))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "ad_items" ("created_at","updated_at","deleted_at","title","description","price")
			VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id","id"`)).
		WithArgs(AnyTime{}, AnyTime{}, nil, item.Title, item.Description, item.Price).
		WillReturnRows(mockedRow)
	mock.ExpectCommit()
	mock.ExpectClose()

	client := NewClient(gormDB)
	actualID, err := client.CreateAd(item)
	if err != nil {
		t.Fatalf("err != nil ; err = %v", err)
	}
	if id != actualID {
		t.Fatalf("id != actualID ; id = %d ; actualID = %d", id, actualID)
	}
}
