package server

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context, port string, controller Controller) error {
	e := gin.New()
	e.GET("/metrics", controller.Metrics())
	e.GET("/health_check", controller.HealthCheck)

	s := http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: e,
	}

	go func() {
		<-ctx.Done()
		s.Shutdown(context.Background())
	}()

	err := s.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}
