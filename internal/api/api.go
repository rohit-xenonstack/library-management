package api

import (
	"context"
	"library-management/backend/internal/api/handler"
	"library-management/backend/internal/api/middleware"
	"library-management/backend/internal/config"
	"library-management/backend/internal/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type API struct {
	Router  *gin.Engine
	Config  *config.Config
	Handler *handler.Handler
}

func NewAPI(cfg *config.Config, h *handler.Handler) *API {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	api := &API{
		Router:  router,
		Config:  cfg,
		Handler: h,
	}
	api.SetupRouter()

	return api
}

func (api *API) SetupRouter() {
	api.Router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Add all routes here
	authRoutes := api.Router.Group("/auth")
	{
		authRoutes.POST("/login", api.Handler.AuthHandler.Login)
	}

	protectedRoutes := api.Router.Group("/")
	protectedRoutes.Use(middleware.JWTAuth())
	{
		ownerRoutes := protectedRoutes.Group("/owner")
		ownerRoutes.Use(middleware.RequirePrivilege(util.OwnerRole))
		{
			ownerRoutes.POST("/create-library", api.Handler.OwnerHandler.CreateLibrary)
			ownerRoutes.POST("/create-new-owner", api.Handler.OwnerHandler.CreateOwner)
			ownerRoutes.POST("/create-new-admin", api.Handler.OwnerHandler.CreateAdmin)
		}
	}
	return
}

func (api *API) Run() error {
	srv := &http.Server{
		Addr:    api.Config.Server.Port,
		Handler: api.Router.Handler(),
	}

	go func() {
		log.Print("Starting Server in 5 seconds...")
		time.Sleep(5 * time.Second)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutting down Server in 5 seconds...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	<-ctx.Done()
	return nil
}
