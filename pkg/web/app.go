package web

import (
	"github.com/gopusher/gateway/pkg/web/middlewares/auth"
	"github.com/gopusher/gateway/pkg/web/middlewares/errors"
	"github.com/gopusher/gateway/pkg/web/middlewares/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

type Config struct {
	Address string `mapstructure:"address" validate:"required"`
	Token   string `mapstructure:"token"`
}

func NewEngine(debug bool) *gin.Engine {
	//mode: debug | release | test
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()

	return engine
}

func NewServer(engine *gin.Engine, config *Config) *http.Server {
	engine.Use(logger.Logger(), errors.Recovery(), auth.Check(config.Token), cors.AllowAll())
	engine.NoRoute(errors.NoFound())

	//handle static file
	//engine.StaticFile("/", "public/index.html")
	//engine.StaticFile("/favicon.ico", "public/favicon.ico")
	//engine.Static("/static", "public/static")

	return &http.Server{
		Addr: config.Address,
		//Handler:        http.TimeoutHandler(engine, 60*time.Second, "request timeout"),
		Handler:      engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		//IdleTimeout:	30 * time.Second,
	}
}
