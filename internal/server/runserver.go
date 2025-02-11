package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlers "github.com/bllooop/coinshop/internal/delivery/api"
	"github.com/bllooop/coinshop/internal/repository"
	"github.com/bllooop/coinshop/internal/usecase"
	logger "github.com/bllooop/coinshop/pkg/logging"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Run() {
	logger.Log.Debug().Msg("Initializing server... / Инициализация сервера...")

	if err := godotenv.Load(); err != nil {
		logger.Log.Error().Err(err).Msg("")
		logger.Log.Fatal().Msg("There was an error with env / Возникла ошибка с env")
	}
	logger.Log.Debug().Msg("Environment variables loaded successfully / Переменные окружения успешно загружены")
	dbpool, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBname:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.Log.Error().Err(err).Msg("Database connection failed / Не удалось установить соединение с базой данных")
		logger.Log.Fatal().Msg("There was an error with database / Произошла ошибка с базой данных")
	}
	logger.Log.Debug().Msg("Database connected successfully / База данных успешно подключена")

	migratePath := "./migrations"
	logger.Log.Debug().Msgf("Running database migrations from path: %s", migratePath)
	if err := repository.RunMigrate(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBname:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}, migratePath); err != nil {
		logger.Log.Error().Err(err).Msg("")
		logger.Log.Fatal().Msg("There was an error when migrating / Возникла ошибка при переносе")
	}
	if err != nil {
		logger.Log.Error().Err(err).Msg("")
		logger.Log.Fatal().Msg("There was an error with database / Произошла ошибка с базой данных")
	}
	logger.Log.Debug().Msg("Initializing repository layer / Инициализация слоя репозитория")
	repos := repository.NewRepository(dbpool)
	logger.Log.Debug().Msg("Initializing usecase layer / Инициализация usecase слоя")
	usecases := usecase.NewUsecase(repos)
	logger.Log.Debug().Msg("Initializing API handlers / Инициализация обработчиков API")
	handler := handlers.NewHandler(usecases)
	srv := new(Server)

	go func() {
		logger.Log.Info().Msg("Starting server... / Запуск сервера...")
		if err := srv.RunServer(os.Getenv("SERVERPORT"), handler.InitRoutes()); err != nil && err == http.ErrServerClosed {
			logger.Log.Info().Msg("Server was shut down gracefully / Сервер был закрыт аккуратно")
		} else {
			logger.Log.Error().Err(err).Msg("")
			logger.Log.Fatal().Msg("There was an error when starting the server / При запуске сервера произошла ошибка")
		}
	}()
	logger.Log.Info().Msg("Server is running / Сервер работает")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	logger.Log.Debug().Msg("Listening for OS termination signals / Прослушивание сигналов завершения работы ОС")
	<-quit
	logger.Log.Info().Msg("Server is shutting down / Сервер отключается")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer dbpool.Close()
	logger.Log.Debug().Msg("Closing database connection / Закрытие соединения с базой данных ")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error().Err(err).Msg("")
		logger.Log.Fatal().Msg("There was an error while shutting down the server / При выключении сервера произошла ошибка")
	}
}

func initConfig() error {
	viper.AddConfigPath("./backend/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
