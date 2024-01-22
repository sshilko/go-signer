package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type HTTPServer struct {
	echo       *echo.Echo
	listenPort uint16
}

type JWTService interface {
	EchoMiddleware() (echo.MiddlewareFunc, error)
}

type HTTPHandler interface {
	GetRoutes() map[string]map[string]echo.HandlerFunc
}

// NewServer returns new web server.
func NewServer(port uint16, jwtService JWTService) *HTTPServer {

	echoServer := echo.New()
	echoServer.HideBanner = true
	echoServer.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${status} ${method} ${uri}\n",
	}))
	echoServer.Logger.SetLevel(log.INFO)

	jwtMiddleware, err := jwtService.EchoMiddleware()
	if err != nil {
		log.Fatal(err.Error())
	}
	echoServer.Use(
		middleware.Recover(),
		middleware.RequestID(),
		jwtMiddleware,
	)
	echoServer.Debug = true

	return &HTTPServer{
		echo:       echoServer,
		listenPort: port,
	}
}

// AddHandler adds new handler to server.
func (s *HTTPServer) AddHandler(h HTTPHandler) {
	routes := h.GetRoutes()
	if len(routes) == 0 {
		return
	}
	for method, v := range routes {
		for path, handler := range v {
			s.echo.Add(method, path, handler)
		}
	}
}

// StopServer stops web server.
func (s *HTTPServer) StopServer(c context.Context) error {
	return s.echo.Shutdown(c)
}

// StartServer starts web server.
func (s *HTTPServer) StartServer() error {
	return s.echo.Start(fmt.Sprintf(":%d", s.listenPort))
}
