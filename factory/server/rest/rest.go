package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/hellofresh/health-go/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vizucode/gokit/factory"
	"github.com/vizucode/gokit/logger"
	"github.com/vizucode/gokit/types"
	"github.com/vizucode/gokit/utils/timezone"
)

// rest an instance of rest handler
type rest struct {
	serverEngine *fiber.App
	service      factory.ServiceFactory
	opt          option
	tz           *time.Location
}

// New creates new handler for rest server
func New(svc factory.ServiceFactory, opts ...OptionFunc) factory.ApplicationFactory {
	var (
		tz = timezone.JakartaTz()

		fiberConfig = fiber.Config{
			AppName:           svc.Name(),
			Prefork:           false,
			ReduceMemoryUsage: true,
		}
	)

	// init an instance rest handler
	srv := &rest{
		tz:           tz,
		opt:          defaultOption(),
		service:      svc,
		serverEngine: fiber.New(fiberConfig),
	}

	for _, o := range opts {
		o(&srv.opt)
	}

	// set custom fiber error handling
	fiberConfig.ErrorHandler = srv.opt.errorHandler
	srv.serverEngine = fiber.New(fiberConfig)

	// add cors middleware
	srv.serverEngine.Use(srv.opt.cors)
	// start handler for health-check
	h, _ := health.New()
	lg := srv.serverEngine.Group("/live")
	lg.Get("/status", adaptor.HTTPHandler(h.Handler()))
	// metrics for prometheus
	mg := srv.serverEngine.Group("/metrics")
	mg.Get("", adaptor.HTTPHandler(promhttp.Handler()))

	// root path for http handler
	rootPath := srv.serverEngine.Group("")
	rootPath.Use(srv.restTraceLogger) // implement http logging

	// print all routes
	for _, route := range srv.serverEngine.GetRoutes(true) {
		if strings.EqualFold(route.Method, http.MethodHead) {
			continue
		}

		logger.Blue(fmt.Sprintf(`[REST-API-ROUTE] (method): %-6s (route): %s`, `"`+route.Method+`"`, `"`+route.Path+`"`))
	}

	return srv
}

func (r *rest) Serve() {
	err := r.serverEngine.Listen(r.opt.httpHost + ":" + r.opt.httpPort)

	switch e := err.(type) {
	case *net.OpError:
		panic(fmt.Errorf("rest server: %s", e))
	}
}

func (r *rest) Shutdown(_ context.Context) {
	defer logger.RedBold("Stopping REST Server")
	_ = r.serverEngine.Shutdown()
}

func (r *rest) Name() string {
	return types.REST.String()
}
