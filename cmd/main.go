package main

import (
	"fmt"
	"food/api"
	"food/config"
	"food/pkg/logger"
	"food/service"
	"net/http"
	"time"

	postgres "food/storage/postgres"
	"food/storage/redis"

	"github.com/gin-gonic/gin"
)

func KeepAlive(cfg *config.Config) {
	for {
		_, err := http.Get(fmt.Sprintf("http://localhost%s/ping", cfg.HTTPPort))
		if err != nil {
			fmt.Println("Error while sending ping:", err)
		} else {
			fmt.Println("Ping sent successfully")
		}
		time.Sleep(1 * time.Minute)
	}
}

func main() {
	cfg := config.Load()

	var loggerLevel = new(string)

	*loggerLevel = logger.LevelDebug

	switch cfg.Environment {
	case config.DebugMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		*loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger("app", *loggerLevel)
	defer func() {
		err := logger.Cleanup(log)
		if err != nil {
			return
		}
	}()

	pgconn, err := postgres.NewConnectionPostgres(&cfg)
	if err != nil {
		panic("postgres no connection: " + err.Error())
	}
	defer pgconn.CloseDB()

	newRedis := redis.New(cfg)
	services := service.New(pgconn, log, newRedis)

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	api.NewApi(r, &cfg, pgconn, log, services)

	go KeepAlive(&cfg)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	fmt.Println("Listening server", cfg.PostgresHost+cfg.HTTPPort)
	err = r.Run(cfg.HTTPPort)
	if err != nil {
		panic(err)
	}

}
