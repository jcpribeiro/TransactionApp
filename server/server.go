package server

import (
	"os"
	"time"
	"transactionapp/api"
	"transactionapp/app"
	"transactionapp/config"

	"transactionapp/internal/cache"
	"transactionapp/internal/mongodb"
	"transactionapp/internal/validate"

	"transactionapp/store"

	"transactionapp/api/swagger"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// Server is a interface to define contract to server up
type Server interface {
	Start()
	Stop()
}

type server struct {
	echo        *echo.Echo
	startedAt   time.Time
	log         logrus.Logger
	app         *app.Container
	redis       *redis.Client
	mongoReader mongodb.MongoDB
	mongoWriter mongodb.MongoDB
	stores      *store.Container
}

func setLogLevel() logrus.Level {
	if config.GlobalConfig.ENV != "prod" {
		return logrus.DebugLevel
	}
	return logrus.ErrorLevel
}

func NewServer() Server {
	return &server{
		startedAt: time.Now(),
		log: logrus.Logger{
			Out:   os.Stderr,
			Level: setLogLevel(),
			Formatter: &logrus.JSONFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
			},
		},
	}
}

func (s *server) Start() {
	// ---- setup echo ----
	s.echo = echo.New()
	s.echo.Validator = validate.New()
	s.echo.Debug = config.GlobalConfig.ENV != "prod"
	s.echo.HideBanner = true

	// ---- setup middlewares ----
	s.echo.Use(emiddleware.Logger())
	s.echo.Use(emiddleware.BodyLimit("2M"))
	s.echo.Use(emiddleware.Recover())
	s.echo.Use(emiddleware.RequestID())
	s.echo.Use(emiddleware.Secure())

	// ---- setup Redis ----
	s.redis = redis.NewClient(&redis.Options{
		Addr:     config.GlobalConfig.Redis.URL,
		Password: config.GlobalConfig.Redis.Password,
		DB:       config.GlobalConfig.Redis.DB,
	})

	cache := cache.NewCache(s.redis, s.log)

	// ---- setup Mongodb ----
	s.mongoReader = mongodb.NewMongoDB(
		config.GlobalConfig.MongoDbReader.URL,
		config.GlobalConfig.MongoDbReader.Scheme,
		true,
	)

	s.mongoWriter = mongodb.NewMongoDB(
		config.GlobalConfig.MongoDbReader.URL,
		config.GlobalConfig.MongoDbReader.Scheme,
		true,
	)

	// ---- setup Store ----
	s.stores = store.NewStore(store.Options{
		MongodbConReader: s.mongoReader.Connect(),
		MongodbConWriter: s.mongoWriter.Connect(),
		Log:              s.log,
	})

	// ---- setup App ----
	s.app = app.NewApp(app.Options{
		Log:    s.log,
		URL:    config.GlobalConfig.FiscalData.URL,
		Stores: s.stores,
	})

	// ---- setup Api ----
	api.Register(api.Options{
		Group: s.echo.Group(""),
		Apps:  s.app,
		Cache: cache,
	})

	// ---- setup documentation ----
	swagger.Register(swagger.Options{
		Group: s.echo.Group("/swagger"),
	})

	// ---- start server ----
	s.log.Info("Start server PID: ", os.Getpid())
	if err := s.echo.Start(config.GlobalConfig.Server.Port); err != nil {
		s.log.Error("cannot starting server ", err.Error())
	}
}

func (s *server) Stop() {
	if err := s.echo.Close(); err != nil {
		s.log.Error("cannot close echo ", err.Error())
	}
}
