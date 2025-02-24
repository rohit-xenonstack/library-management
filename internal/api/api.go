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

	"github.com/gin-contrib/cors"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:4173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := &API{
		Router:  router,
		Config:  cfg,
		Handler: h,
	}
	api.SetupRouter()

	return api
}

func (api *API) SetupRouter() {
	baseRoute := api.Router.Group("/api")
	{
		baseRoute.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		authRoutes := baseRoute.Group("/auth")
		{
			authRoutes.POST("/login", api.Handler.AuthHandler.Login)
			authRoutes.POST("/register", api.Handler.AuthHandler.UserSignup)
			authRoutes.POST("/refresh", api.Handler.AuthHandler.RefreshToken)
		}

		protectedRoutes := baseRoute.Group("/protected")
		protectedRoutes.Use(middleware.JWTAuth())
		{
			protectedRoutes.GET("/me", api.Handler.AuthHandler.UserDetails)
			ownerRoutes := protectedRoutes.Group("/owner")
			ownerRoutes.Use(middleware.RequirePrivilege(util.OwnerRole))
			{
				ownerRoutes.POST("/create-library", api.Handler.OwnerHandler.CreateLibrary)
				ownerRoutes.POST("/onboard-admin", api.Handler.OwnerHandler.CreateAdmin)
				ownerRoutes.GET("/libraries", api.Handler.OwnerHandler.GetLibraries)
				ownerRoutes.POST("/admins", api.Handler.OwnerHandler.GetAdmins)
			}
			adminRoutes := protectedRoutes.Group("/admin")
			adminRoutes.Use(middleware.RequirePrivilege(util.AdminRole))
			{
				adminRoutes.POST("/add-book", api.Handler.AdminHandler.AddBook)
				adminRoutes.POST("/remove-book", api.Handler.AdminHandler.RemoveBook)
				adminRoutes.PATCH("/update-book", api.Handler.AdminHandler.UpdateBook)
				adminRoutes.GET("/issue-requests", api.Handler.AdminHandler.ListIssueRequests)
				adminRoutes.POST("/approve-issue-request", api.Handler.AdminHandler.ApproveIssueRequest)
				adminRoutes.POST("/reject-issue-request", api.Handler.AdminHandler.RejectIssueRequest)
				adminRoutes.POST("/books", api.Handler.AdminHandler.SearchBook)
				adminRoutes.GET("/books/:isbn", api.Handler.AdminHandler.SearchBookByISBN)

			}
			readerRoutes := protectedRoutes.Group("/reader")
			readerRoutes.Use(middleware.RequirePrivilege(util.ReaderRole))
			{
				readerRoutes.GET("/latest/:isbn", api.Handler.ReaderHandler.GetLatestAvailability)
				readerRoutes.POST("/books", api.Handler.ReaderHandler.SearchBook)
				readerRoutes.POST("/request-issue", api.Handler.ReaderHandler.RaiseIssueRequest)
			}
		}
	}
}

func (api *API) Run() error {
	srv := &http.Server{
		Addr:    api.Config.Server.Port,
		Handler: api.Router.Handler(),
	}

	go func() {
		log.Print("Starting Server in 5 seconds...")
		time.Sleep(1 * time.Second)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutting down Server in 5 seconds...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	<-ctx.Done()
	return nil
}
