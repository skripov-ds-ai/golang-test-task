package main

import (
	"fmt"
	"golang-test-task/internal/database"
	"golang-test-task/internal/facade"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/gorilla/mux"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @Title Ad Submitting Service
// @Version 0.1
// @Description A sample server for creating and getting ads

// @Contact.name Denis Skripov
// @Contact.email nizhikebinesi@gmail.com

// @Host localhost:8888
// @BasePath /api/v0.1

// App is wrapper to simplify app creating
type App struct {
	Router *mux.Router
	logger *zap.Logger
}

// NewApp creates an app
func NewApp(version string, client *database.Client, v *validator.Validate, logger *zap.Logger) *App {
	logic := facade.NewHandlerFacade(client, v, logger)

	getAdHandler, _ := logic.GetHandler("get_ad")
	listAdsHandler, _ := logic.GetHandler("list_ads")
	createAdHandler, _ := logic.GetHandler("create_ad")

	r := mux.NewRouter()

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/vars", http.DefaultServeMux)

	apiRouter := r.PathPrefix(fmt.Sprintf("/api/v%s", version)).Subrouter()
	apiRouter.HandleFunc("/ads", listAdsHandler).Methods("GET")
	apiRouter.HandleFunc("/ads", createAdHandler).Methods("POST")
	apiRouter.HandleFunc("/ads/{id:[0-9]+}", getAdHandler).Methods("GET")

	a := &App{Router: r, logger: logger}
	return a
}

// Run is need to run an app
func (a *App) Run() {
	err := http.ListenAndServe(":3000", a.Router)
	if err != nil {
		a.logger.Panic("not nil serving", zap.Error(err))
	}
}

func main() {
	runtime.GOMAXPROCS(4)

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build(zap.Hooks(func(entry zapcore.Entry) error {
		if entry.Level == zapcore.DebugLevel ||
			entry.Level == zapcore.WarnLevel ||
			entry.Level == zapcore.ErrorLevel ||
			entry.Level == zapcore.PanicLevel ||
			entry.Level == zapcore.DPanicLevel {
			defer sentry.Flush(2 * time.Second)
			sentry.CaptureMessage(fmt.Sprintf("%s, Line No: %d :: %s", entry.Caller.File, entry.Caller.Line, entry.Message))
		}
		return nil
	}))
	defer func() {
		_ = logger.Sync()
	}()

	// https://github.com/shopspring/decimal/issues/21
	decimal.MarshalJSONWithoutQuotes = true

	// TODO: sync it with git tags
	apiVersion := os.Getenv("API_VERSION")

	dsn := os.Getenv("DB_DSN")

	// TODO: add zap to sentry - https://github.com/TheZeroSlave/zapsentry
	sentryDSN := os.Getenv("SENTRY_DSN")
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDSN,
		Release:          fmt.Sprintf("golang-test-task@%s", apiVersion),
		Debug:            true,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		logger.Panic("sentry does not init", zap.Error(err))
	}

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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&database.AdItem{}, &database.ImageURL{})
	if err != nil {
		logger.Panic("failed to automigrate", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panic("failed to get DB from gorm", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(time.Hour)

	client := database.NewClient(db)
	app := NewApp(apiVersion, client, v, logger)
	app.Run()
}
