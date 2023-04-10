package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evgeniy-dammer/marketplace-api/internal/config"
	deliveryHttp "github.com/evgeniy-dammer/marketplace-api/internal/delivery/http"
	postgresStorage "github.com/evgeniy-dammer/marketplace-api/internal/repository/storage/postgres"
	redisStorage "github.com/evgeniy-dammer/marketplace-api/internal/repository/storage/redis"
	useCaseAuthentication "github.com/evgeniy-dammer/marketplace-api/internal/usecase/authentication"
	useCaseAuthorization "github.com/evgeniy-dammer/marketplace-api/internal/usecase/authorization"
	useCaseCategory "github.com/evgeniy-dammer/marketplace-api/internal/usecase/category"
	useCaseComment "github.com/evgeniy-dammer/marketplace-api/internal/usecase/comment"
	useCaseFavorite "github.com/evgeniy-dammer/marketplace-api/internal/usecase/favorite"
	useCaseImage "github.com/evgeniy-dammer/marketplace-api/internal/usecase/image"
	useCaseItem "github.com/evgeniy-dammer/marketplace-api/internal/usecase/item"
	useCaseMessage "github.com/evgeniy-dammer/marketplace-api/internal/usecase/message"
	useCaseOrder "github.com/evgeniy-dammer/marketplace-api/internal/usecase/order"
	useCaseOrganization "github.com/evgeniy-dammer/marketplace-api/internal/usecase/organization"
	useCaseRule "github.com/evgeniy-dammer/marketplace-api/internal/usecase/rule"
	useCaseSpecification "github.com/evgeniy-dammer/marketplace-api/internal/usecase/specification"
	useCaseTable "github.com/evgeniy-dammer/marketplace-api/internal/usecase/table"
	useCaseUser "github.com/evgeniy-dammer/marketplace-api/internal/usecase/user"
	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/evgeniy-dammer/marketplace-api/pkg/server"
	"github.com/evgeniy-dammer/marketplace-api/pkg/store/postgres"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	vault "github.com/hashicorp/vault/api"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Run() {
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

	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = viper.GetString("vault.address")

	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		logger.Logger.Fatal("unable to connect to vault", zap.String("error", err.Error()))
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	marketplace, err := client.Logical().Read("cubbyhole/marketplace")
	if err != nil {
		logger.Logger.Fatal("unable to read data from vault", zap.String("error", err.Error()))
	}

	viper.Set("JWT_KEY", marketplace.Data["JWT_KEY"])

	dbMasterPassword, exists := marketplace.Data["DB_MASTER_PASSWORD"].(string)
	if !exists {
		logger.Logger.Fatal("unable to get master database password")
	}

	dbSlavePassword, exists := marketplace.Data["DB_SLAVE_PASSWORD"].(string)
	if !exists {
		logger.Logger.Fatal("unable to get slave database password")
	}

	// service settings
	isCacheOn := viper.GetBool("service.cache")
	isTracingOn := viper.GetBool("service.tracing")
	routerMode := viper.GetString("service.router")

	// establishing database connection
	master, adapterMaster, err := postgres.NewPostgresDB(postgres.DBConfig{
		Host:     viper.GetString("database.master.host"),
		Port:     viper.GetString("database.master.port"),
		Username: viper.GetString("database.master.username"),
		Password: dbMasterPassword,
		DBName:   viper.GetString("database.master.dbname"),
		SSLMode:  viper.GetString("database.master.sslmode"),
	})
	if err != nil {
		logger.Logger.Fatal("master database initialization failed", zap.String("error", err.Error()))
	}

	defer func(database *sqlx.DB) {
		err = database.Close()
		if err != nil {
			logger.Logger.Fatal("failed to close master database connection", zap.String("error", err.Error()))
		}
	}(master)

	slave, adapterSlave, err := postgres.NewPostgresDB(postgres.DBConfig{
		Host:     viper.GetString("database.slave.host"),
		Port:     viper.GetString("database.slave.port"),
		Username: viper.GetString("database.slave.username"),
		Password: dbSlavePassword,
		DBName:   viper.GetString("database.slave.dbname"),
		SSLMode:  viper.GetString("database.slave.sslmode"),
	})
	if err != nil {
		logger.Logger.Fatal("slave database initialization failed", zap.String("error", err.Error()))
	}

	defer func(database *sqlx.DB) {
		err = database.Close()
		if err != nil {
			logger.Logger.Fatal("failed to close slave database connection", zap.String("error", err.Error()))
		}
	}(slave)

	var redisClient *redis.Client

	if isCacheOn {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     net.JoinHostPort(viper.GetString("cache.host"), viper.GetString("cache.port")),
			Password: "", // marketplace.Data["REDIS_PASSWORD"].(string)
			DB:       viper.GetInt("cache.database"),
		})

		defer func(redisClient *redis.Client) {
			err = redisClient.Close()
			if err != nil {
				logger.Logger.Fatal("unable to close redis client", zap.String("error", err.Error()))
			}
		}(redisClient)

		_, err = redisClient.Ping(context.TODO()).Result()

		if err != nil {
			logger.Logger.Fatal("cache initialization failed", zap.String("error", err.Error()))
		}
	} else {
		logger.Logger.Info("cache is turned off")
	}

	if isTracingOn {
		tracerProvider, err := tracing.InitTracer(viper.GetString("tracing.url"), viper.GetString("service.name"))
		if err != nil {
			logger.Logger.Fatal("unable to start tracer provider", zap.String("error", err.Error()))
		}

		defer func() {
			if err := tracerProvider.Shutdown(context.Background()); err != nil {
				logger.Logger.Fatal("unable to shutdown tracer provider", zap.String("error", err.Error()))
			}
		}()
	} else {
		logger.Logger.Info("tracing is turned off")
	}

	// repositories
	repoStorage := postgresStorage.New(
		master,
		slave,
		postgresStorage.Options{Timeout: time.Duration(viper.GetInt("database.timeout")) * time.Second},
		isTracingOn,
	)

	repoCache := redisStorage.New(
		redisClient,
		redisStorage.Options{
			Timeout: time.Duration(viper.GetInt("cache.timeout")) * time.Second,
			TTL:     time.Duration(viper.GetInt("cache.ttl")) * time.Minute,
		},
		isTracingOn,
	)

	// use cases
	ucAuthentication := useCaseAuthentication.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucAuthorization := useCaseAuthorization.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucUser := useCaseUser.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucOrganization := useCaseOrganization.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucCategory := useCaseCategory.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucItem := useCaseItem.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucTable := useCaseTable.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucOrder := useCaseOrder.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucImage := useCaseImage.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucComment := useCaseComment.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucSpecification := useCaseSpecification.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucFavorite := useCaseFavorite.New(repoStorage, repoCache, isTracingOn)
	ucRule := useCaseRule.New(repoStorage, repoCache, isTracingOn, isCacheOn)
	ucMessage := useCaseMessage.New(repoStorage, repoCache, isTracingOn, isCacheOn)

	// deliveries
	deliveryHTTP := deliveryHttp.New(
		ucAuthentication,
		ucAuthorization,
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
		ucMessage,
		adapterMaster,
		adapterSlave,
		isTracingOn,
	)

	// create new server
	srv := new(server.Server)

	srvConfig := server.Config{
		Port:           viper.GetString("server.port"),
		Handler:        deliveryHTTP.InitRoutes(routerMode),
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
