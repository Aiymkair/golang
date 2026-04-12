package app

import (
	"Assignment-7/config"
	"Assignment-7/internal/controller/htpp/v1"
	"Assignment-7/internal/usecase"
	"Assignment-7/internal/usecase/repo"
	"Assignment-7/pkg/logger"
	"Assignment-7/pkg/postgres"
	"Assignment-7/utils"

	"github.com/gin-gonic/gin"
)

type App struct {
	router *gin.Engine
	cfg    *config.Config
	log    logger.Logger
}

func NewApp() *App {
	//загружаем конфиг
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	//логгер
	log := logger.New()

	//postgres(неть)
	_ = postgres.New()

	//репозиторий in-memory
	userRepo := repo.NewUserRepo()

	//usecase
	userUsecase := usecase.NewUserUseCase(userRepo)

	//роутер
	router := gin.Default()

	//глобальный rate limiter
	router.Use(utils.RateLimitMiddleware())

	//регистрация роутов
	v1.NewRouter(router, userUsecase, log)

	return &App{
		router: router,
		cfg:    cfg,
		log:    log,
	}
}

func (a *App) Run() error {
	a.log.Info("Starting server on port " + a.cfg.Port)
	return a.router.Run(":" + a.cfg.Port)
}
