package api

import (
	"context"
	"github.com/duyquang6/git-watchdog/internal/configuration"
	"github.com/duyquang6/git-watchdog/internal/rabbitmq"
	"github.com/duyquang6/git-watchdog/internal/repository"
	"github.com/streadway/amqp"
	"net/http"

	_ "github.com/duyquang6/git-watchdog/docs"
	repoControllerPkg "github.com/duyquang6/git-watchdog/internal/controller/repository"
	"github.com/duyquang6/git-watchdog/internal/database"
	"github.com/duyquang6/git-watchdog/internal/middleware"
	"github.com/duyquang6/git-watchdog/internal/service"
	"github.com/duyquang6/git-watchdog/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func (s *httpServer) setupDependencyAndRouter(ctx context.Context, r *gin.Engine,
	db database.DBFactory, amqpClient *amqp.Connection) {
	appConfig := configuration.FromContext(ctx)
	repoRepo := repository.NewRepoRepository()
	scanRepo := repository.NewScanRepository()

	channel := rabbitmq.NewQueueChannelFromConnection(appConfig.RabbitMQConfig(), amqpClient)

	repoService := service.NewRepositoryService(db, repoRepo, scanRepo, channel)

	repoController := repoControllerPkg.NewController(repoService)

	initRoute(ctx, r, repoController)
}

// @title GitWatchdog API
// @version 1.0
// @description GitWatchdog API Spec
// @termsOfService http://swagger.io/terms/

// @contact.name Quang Nguyen
// @contact.url https://duyquang6.github.io
// @contact.email nguyenduyquang06@gmail.com

// @BasePath /api/v1
func initRoute(ctx context.Context, r *gin.Engine, repoController *repoControllerPkg.Controller) {
	r.Use(middleware.PopulateRequestID())
	r.Use(middleware.PopulateLogger(logging.FromContext(ctx)))
	apiV1 := r.Group("/api/v1")
	{
		sub := apiV1.Group("/repositories")
		{
			sub.POST("", repoController.HandleCreateRepository())
			sub.GET("/:id", repoController.HandleGetOneRepository())
			sub.DELETE("/:id", repoController.HandleDelete())
			sub.PUT("/:id", repoController.HandleUpdateRepository())

			sub.POST("/:id/scans", repoController.HandleIssueScan())
			sub.GET("/:id/scans", repoController.HandleListScan())
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
