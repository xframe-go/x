package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/validate"
)

type EchoServer struct {
	*echo.Echo
	cfg Config
}

func NewEcho(cfg Config) *EchoServer {
	e := echo.New()

	e.Validator = &validate.FormRequestValidator{}

	return &EchoServer{
		Echo: e,
		cfg:  cfg,
	}
}

func (s *EchoServer) Resource(group *echo.Group, name string, handler contracts.ResourceHandler) {
	group.GET("/"+name, handler.List)
	group.GET("/"+name+"/_batch", handler.BatchList)
	group.GET("/"+name+"/:id", handler.Show)
	group.POST("/"+name, handler.Create)
	group.PUT("/"+name+"/:id", handler.Update)
	group.DELETE("/"+name+"/:id", handler.Destroy)
}

func (s *EchoServer) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         fmt.Sprintf(":%d", s.cfg.Port),
		GracefulTimeout: 5 * time.Second,
	}
	if err := sc.Start(ctx, s.Echo); err != nil {
		s.Echo.Logger.Error("failed to start server", "error", err)
	}

	return sc.Start(ctx, s.Echo)
}
