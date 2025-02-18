package api

import (
	"context"
	"library-management/backend/internal/config"
	"library-management/backend/pkg/token"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type API struct {
	Router   *gin.Engine
	Config   *config.Config
	Token    token.Token
	Database *gorm.DB
}

func NewAPI(cfg *config.Config, db *gorm.DB) (*API, error) {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	token, err := token.NewJWToken(cfg.JWT.SecretKey)
	if err != nil {
		return nil, err
	}

	API := &API{
		Config:   cfg,
		Database: db,
		Token:    token,
	}

	API.SetupRouter()

	return API, nil
}

func (api *API) SetupRouter() {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	// Add all routes here
	api.Router = router
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
