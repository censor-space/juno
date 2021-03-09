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
	e.POST("/v1/post_question", controller.PostQuestion)
	e.POST("/v1/clear_current_question", controller.UpdateClearCurrentQuestion)
	e.GET("/v1/user_score", controller.GetUserScore)
	e.POST("/v1/line_callback", controller.CallbackFromLine)

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
