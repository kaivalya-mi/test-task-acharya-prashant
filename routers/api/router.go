package api

import (
	"context"
	"fmt"
	"net/http"
	v1Ctl "test-task/controllers/v1"
	v1Service "test-task/services/v1"
	"test-task/shared/config"
	"test-task/shared/log"
	"test-task/shared/utils/middleware"
	"test-task/validator"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// IRoutes is
type IRoutes interface {
	Setup()
	Run()
	Close(ctx context.Context) error
}

// Routes is
type Routes struct {
	router     *gin.Engine
	server     *http.Server
	config     config.IConfig
	authCtl    *v1Ctl.AuthCtl
	middleware middleware.IMiddleware
}

// NewRouter is
func NewRouter(config config.IConfig) IRoutes {
	validation := validator.NewAPIValidatorService()
	authSrv := v1Service.NewAuthService()
	middlewareSrv := middleware.NewMiddlewareService(config)

	authCtl := v1Ctl.AuthController(validation, authSrv, middlewareSrv)

	router := gin.Default()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.App().Port),
		Handler: router,
	}

	return &Routes{
		router,
		server,
		config,
		authCtl,
		middlewareSrv,
	}
}

func (rt *Routes) Run() {
	log.GetLog().Info("", "service listen on "+rt.config.App().Port)

	err := rt.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.GetLog().Fatal("", "listen: %s\n", err)
	}
}

func (rt *Routes) Close(ctx context.Context) error {
	if rt.server != nil {
		return rt.server.Shutdown(ctx)
	}
	return nil
}

func (rt *Routes) Setup() {
	router := rt.router
	auth := rt.authCtl
	middleware := rt.middleware

	rt.setupCors()
	rt.setupDefaultEndpoints()

	app := router.Group("/api/v1")

	app.POST("/sign-up", auth.SignUp)
	app.POST("/sign-in", auth.SignIn)
	app.POST("/refresh-token", auth.RefreshToken)

	//protected route
	app.GET("/user-profile", middleware.AuthHandler(), auth.GetProfile)
	app.POST("/sign-out", middleware.AuthHandler(), auth.SignOut)

}

func (rt *Routes) setupCors() {
	rt.router.Use(cors.New(cors.Config{
		ExposeHeaders:   []string{"Data-Length"},
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
		AllowAllOrigins: true,
		MaxAge:          12 * time.Hour,
	}))
}

func (rt *Routes) setupDefaultEndpoints() {
	rt.router.GET("/ping", func(c *gin.Context) {
		var msg string
		if rt.config.Env() == "production" {
			msg = fmt.Sprintf("Pong! I am %s. Version is %s.", rt.config.AppRegion(), rt.config.AppVersion())
		} else {
			msg = "pong"
		}
		c.JSON(200, gin.H{"message": msg})
	})
}
