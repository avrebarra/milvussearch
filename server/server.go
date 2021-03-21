package server

import (
	"net/http"

	"github.com/avrebarra/milvus-dating/component/matchmakerservice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.linkaja.com/gopkg/validator"
)

type Config struct {
	Service matchmakerservice.Service `validate:"required"`
}

type Server struct {
	config Config
	router http.Handler
}

func New(cfg Config) (*Server, error) {
	if err := validator.Validate(cfg); err != nil {
		return nil, err
	}

	s := Server{
		config: cfg,
		router: echo.New(),
	}

	// setup routes
	if err := s.setupRoutes(); err != nil {
		return nil, err
	}

	return &s, nil
}

// setupRoutes defines mappings of server paths to handler's functions
func (s *Server) setupRoutes() (err error) {
	// ---
	// setup handlers
	handlerCommon := HandlerCommon{Config: s.config}

	// ---
	// setup router
	router := s.router.(*echo.Echo)
	router.Validator = &CustomValidator{}

	router.GET("/", handlerCommon.PingV1())

	router.POST("/user", handlerCommon.PingV1())
	router.GET("/user/:id", handlerCommon.PingV1())
	router.POST("/find_match", handlerCommon.PingV1())

	// middlewares and add-ons
	router.Pre(middleware.RemoveTrailingSlash())
	router.Pre(middleware.RequestID())

	return
}

func (s *Server) GetHandler() http.Handler {
	return s.router
}

func (s *Server) Close() error {
	return nil
}

// ***

type CustomValidator struct{}

func (cv *CustomValidator) Validate(i interface{}) error {
	return validator.Validate(i)
}

func GetHandler() []string {
	return []string{
		"POST /user",
		"GET  /user/:id",
		"POST /find_match",
	}
}
