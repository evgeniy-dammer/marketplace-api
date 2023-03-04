package main

import (
	"context"
	rediscache "github.com/evgeniy-dammer/emenu-api/pkg/store/redis"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/evgeniy-dammer/emenu-api/internal/config"
	deliveryHttp "github.com/evgeniy-dammer/emenu-api/internal/delivery/http"
	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	postgresStorage "github.com/evgeniy-dammer/emenu-api/internal/repository/storage/postgres"
	redisCache "github.com/evgeniy-dammer/emenu-api/internal/repository/storage/redis"
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
	"github.com/evgeniy-dammer/emenu-api/pkg/store/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// set log format as JSON
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// initializing config
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	// loading env variables
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
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
		logrus.Fatalf("failed to initialize database: %s", err)
	}

	defer func(database *sqlx.DB) {
		err = database.Close()
		if err != nil {
			logrus.Fatalf("failed to close database connection: %s", err)
		}
	}(database)

	var cache *cache.Cache

	if viper.GetBool("cache.mode") {
		// establishing cache connection
		cache, err = rediscache.NewRedisCache(redis.Options{
			Addr:     net.JoinHostPort(viper.GetString("cache.host"), viper.GetString("cache.port")),
			Password: "", // os.Getenv("REDIS_PASSWORD"),
			DB:       viper.GetInt("cache.database"),
		})
		if err != nil {
			logrus.Fatalf("failed to initialize cache: %s", err)
		}
	} else {
		logrus.Info("cache is turned off")
	}

	// repositories
	repoStorage := postgresStorage.New(database)
	repoCache := redisCache.New(cache)

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
	srv := new(domain.Server)

	srvConfig := domain.ServerConfig{
		Port:           viper.GetString("server.port"),
		Handler:        deliveryHTTP.InitRoutes(viper.GetString("router.mode")),
		ReadTimeout:    viper.GetInt("server.read_timeout"),
		WriteTimeout:   viper.GetInt("server.write_timeout"),
		IdleTimeout:    viper.GetInt("server.idle_timeout"),
		MaxHeaderBytes: viper.GetInt("server.max_header_bytes"),
	}

	go func() {
		if err = srv.Run(srvConfig); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Println("Application started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("Application shutdown...")

	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error ocured on server shutdown: %s", err.Error())
	}
}
