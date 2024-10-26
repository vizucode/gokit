package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/vizucode/gokit/logger"
)

// OptionFunc setter rest options
type OptionFunc func(*option)

// option an instance of rest options
type option struct {
	cors         fiber.Handler
	httpPort     string
	httpHost     string
	engineOption func(app *fiber.App)
	log          *logrus.Logger

	// it's recomended to set error handling, default is fiber.DefaultErrorHandler
	errorHandler fiber.ErrorHandler
}

// defaultOption default options for rest
func defaultOption() option {
	return option{
		httpPort: "8080",
		log:      logger.Logrus(),
		cors: func(c *fiber.Ctx) error {
			return c.Next()
		},
		errorHandler: fiber.DefaultErrorHandler,
	}
}

// SetHTTPPort set http port
func SetHTTPPort(httpPort int) OptionFunc {
	return func(o *option) {
		o.httpPort = fmt.Sprintf("%d", httpPort)
	}
}

// SetCors set cors options
func SetCors(cors fiber.Handler) OptionFunc {
	return func(o *option) {
		o.cors = cors
	}
}

// SetEngineOption set rest engine
func SetEngineOption(app func(*fiber.App)) OptionFunc {
	return func(o *option) {
		o.engineOption = app
	}
}

func SetHTTPHost(httpHost string) OptionFunc {
	return func(o *option) {
		o.httpHost = httpHost
	}
}

func SetErrorHandler(errorHandler fiber.ErrorHandler) OptionFunc {
	return func(o *option) {
		o.errorHandler = errorHandler
	}
}
