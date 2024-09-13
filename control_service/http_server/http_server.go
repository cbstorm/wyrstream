package httpserver

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cbstorm/wyrstream/control_service/configs"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var instance *HttpServer
var instance_sync sync.Once

func GetApp() *HttpServer {
	if instance == nil {
		instance_sync.Do(func() {
			instance = &HttpServer{
				logger: logger.NewLogger("HTTP_SERVER"),
				config: configs.GetConfig(),
			}
		})
	}
	return instance
}

type HttpServer struct {
	FiberApp *fiber.App
	logger   *logger.Logger
	config   *configs.Config
	// routes
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
		BodyLimit: 30 * 1024 * 1024,
	})
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Accept-Encoding, X-Token, X-Refresh, Ngrok-Skip-Browser-Warning, Tz-Offset",
	}))
	app.Use(healthcheck.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	a.FiberApp = app
	return a
}

func (a *HttpServer) LoadRoutes() *HttpServer {
	return a
}

func (a *HttpServer) Listen() error {
	return a.FiberApp.Listen(fmt.Sprintf("%s:%d", a.config.HTTP_HOST, a.config.HTTP_PORT))
}
