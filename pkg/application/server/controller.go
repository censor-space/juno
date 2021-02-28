package server

import (
	"net/http"

	"github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/line"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Controller interface {
	HealthCheck(ctx *gin.Context)
	Metrics() gin.HandlerFunc
	PostQuestion(ctx *gin.Context)
}

type controller struct {
	LineBot line.LineBot
}

func NewController(lineBot line.LineBot) Controller {
	return &controller{
		LineBot: lineBot,
	}
}

func (c *controller) HealthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}

func (c *controller) Metrics() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func (c *controller) PostQuestion(ctx *gin.Context) {
	err := c.LineBot.PostQuiz()
	if err != nil {
		panic(err)
	}
	ctx.String(http.StatusOK, "OK")
}
