package server

import (
	"net/http"

	"github.com/anitta/eguchi-wedding-bot/pkg/domain/quiz"
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
	var jsonQuestion quiz.Question
	err := ctx.ShouldBindJSON(&jsonQuestion)
	if err != nil {
		ctx.String(http.StatusBadRequest, "400 Bad Request")
		return
	}

	err = c.LineBot.PostQuiz(jsonQuestion)
	if err != nil {
		ctx.String(http.StatusBadRequest, "400 Bad Request")
		return
	}
	ctx.String(http.StatusOK, "OK")
	return
}
