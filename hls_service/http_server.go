package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var instance *HttpServer
var instance_sync sync.Once

func GetHttpServer() *HttpServer {
	if instance == nil {
		instance_sync.Do(func() {
			instance = &HttpServer{
				logger: logger.NewLogger("HLS_HTTP_SERVER"),
				config: GetConfig(),
			}
		})

	}
	return instance
}

type HttpServer struct {
	fiber_app *fiber.App
	logger    *logger.Logger
	config    *Config
}

func (a *HttpServer) Init() *HttpServer {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			a.logger.Error("Internal server error: %v", err)
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			a.logger.Error("%s %s %s", ctx.IP(), ctx.OriginalURL(), ctx.BodyRaw())
			return ctx.Status(code).SendString("Internal Server Error")
		},
		BodyLimit: 5 * 1024 * 1024,
	})
	app.Static("/", "./public")
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Accept-Encoding, X-Token, X-Refresh, Ngrok-Skip-Browser-Warning, Tz-Offset",
	}))
	app.Use(healthcheck.New())
	a.fiber_app = app
	return a
}

func (a *HttpServer) Listen() error {
	return a.fiber_app.Listen(fmt.Sprintf("%s:%d", a.config.HLS_HTTP_HOST, a.config.HLS_HTTP_PORT))
}
