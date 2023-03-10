package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/config"
	deliveryHttp "github.com/evgeniy-dammer/emenu-api/internal/delivery/http"
	postgresStorage "github.com/evgeniy-dammer/emenu-api/internal/repository/storage/postgres"
	redisStorage "github.com/evgeniy-dammer/emenu-api/internal/repository/storage/redis"
	useCaseAuthentication "github.com/evgeniy-dammer/emenu-api/internal/usecase/authentication"
	useCaseCategory "github.com/evgeniy-dammer/emenu-api/internal/usecase/category"
	useCaseComment "github.com/evgeniy-dammer/emenu-api/internal/usecase/comment"
	useCaseFavorite "github.com/evgeniy-dammer/emenu-api/internal/usecase/favorite"
	useCaseImage "github.com/evgeniy-dammer/emenu-api/internal/usecase/image"
	useCaseItem "github.com/evgeniy-dammer/emenu-api/internal/usecase/item"
	useCaseOrder "github.com/evgeniy-dammer/emenu-api/internal/usecase/order"
	useCaseOrganization "github.com/evgeniy-dammer/emenu-api/internal/usecase/organization"
	useCaseRule "github.com/evgeniy-dammer/emenu-api/internal/usecase/rule"
	useCaseSpecification "github.com/evgeniy-dammer/emenu-api/internal/usecase/specification"
	useCaseTable "github.com/evgeniy-dammer/emenu-api/internal/usecase/table"
	useCaseUser "github.com/evgeniy-dammer/emenu-api/internal/usecase/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/evgeniy-dammer/emenu-api/pkg/server"
	"github.com/evgeniy-dammer/emenu-api/pkg/store/postgres"
	redisStore "github.com/evgeniy-dammer/emenu-api/pkg/store/redis"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// initializing config
	if err := config.InitConfig(); err != nil {
		logger.Logger.Fatal("config initialization failed", zap.String("error", err.Error()))
	}

	// loading env variables
	if err := godotenv.Load(); err != nil {
		logger.Logger.Fatal("env variables loading failed", zap.String("error", err.Error()))
	}

	err := logger.InitLogger()
	if err != nil {
		logger.Logger.Fatal("logger initialization failed", zap.String("error", err.Error()))
	}

	// establishing database connection
	database, adapter, err := postgres.NewPostgresDB(postgres.DBConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("database.dbname"),
		SSLMode:  viper.GetString("database.sslmode"),
	})
	if err != nil {
		logger.Logger.Fatal("database initialization failed", zap.String("error", err.Error()))
	}

	defer func(database *sqlx.DB) {
		err = database.Close()
		if err != nil {
			logger.Logger.Fatal("failed to close database connection", zap.String("error", err.Error()))
		}
	}(database)

	var rcache *cache.Cache

	if viper.GetBool("cache.mode") {
		// establishing cache connection
		rcache, err = redisStore.NewRedisCache(redis.Options{
			Addr:     net.JoinHostPort(viper.GetString("cache.host"), viper.GetString("cache.port")),
			Password: "", // os.Getenv("REDIS_PASSWORD"),
			DB:       viper.GetInt("cache.database"),
		})
		if err != nil {
			logger.Logger.Fatal("cache initialization failed", zap.String("error", err.Error()))
		}
	} else {
		logger.Logger.Info("cache is turned off")
	}

	if viper.GetBool("service.tracing") {
		closer, err := tracing.New()
		if err != nil {
			panic(err)
		}
		defer func() {
			if err = closer.Close(); err != nil {
				log.Error(err)
			}
		}()
	} else {
		logger.Logger.Info("tracing is turned off")
	}

	// repositories
	repoStorage := postgresStorage.New(
		database,
		postgresStorage.Options{Timeout: time.Duration(viper.GetInt("database.timeout")) * time.Second},
	)

	repoCache := redisStorage.New(
		rcache,
		redisStorage.Options{Timeout: time.Duration(viper.GetInt("database.timeout")) * time.Second},
	)

	// use cases
	ucAuthentication := useCaseAuthentication.New(repoStorage, repoCache)
	ucUser := useCaseUser.New(repoStorage, repoCache)
	ucOrganization := useCaseOrganization.New(repoStorage, repoCache)
	ucCategory := useCaseCategory.New(repoStorage, repoCache)
	ucItem := useCaseItem.New(repoStorage, repoCache)
	ucTable := useCaseTable.New(repoStorage, repoCache)
	ucOrder := useCaseOrder.New(repoStorage, repoCache)
	ucImage := useCaseImage.New(repoStorage, repoCache)
	ucComment := useCaseComment.New(repoStorage, repoCache)
	ucSpecification := useCaseSpecification.New(repoStorage, repoCache)
	ucFavorite := useCaseFavorite.New(repoStorage, repoCache)
	ucRule := useCaseRule.New(repoStorage, repoCache)

	// deliveries
	deliveryHTTP := deliveryHttp.New(
		ucAuthentication,
		ucUser,
		ucOrganization,
		ucCategory,
		ucItem,
		ucTable,
		ucOrder,
		ucImage,
		ucComment,
		ucSpecification,
		ucFavorite,
		ucRule,
		adapter,
	)

	// create new server
	srv := new(server.Server)

	srvConfig := server.Config{
		Port:           viper.GetString("server.port"),
		Handler:        deliveryHTTP.InitRoutes(viper.GetString("router.mode")),
		ReadTimeout:    viper.GetInt("server.read_timeout"),
		WriteTimeout:   viper.GetInt("server.write_timeout"),
		IdleTimeout:    viper.GetInt("server.idle_timeout"),
		MaxHeaderBytes: viper.GetInt("server.max_header_bytes"),
	}

	go func() {
		if err = srv.Run(srvConfig); err != nil {
			logger.Logger.Fatal("server starting failed", zap.String("error", err.Error()))
		}
	}()

	logger.Logger.Info("application started", zap.String("port", viper.GetString("server.port")))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Logger.Info("application shutdown")

	if err = srv.Shutdown(context.Background()); err != nil {
		logger.Logger.Fatal("error occurred on server shutdown", zap.String("error", err.Error()))
	}
}
