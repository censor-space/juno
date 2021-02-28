package server

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type Controller interface {
	HealthCheck(ctx *gin.Context)
	Metrics() gin.HandlerFunc
    PostQuestion(ctx *gin.Context)
}

type controller struct {}

func NewController() Controller {
    return &controller{}
}


func (c *controller) HealthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}

func (c *controller) Metrics() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func (c *controller) PostQuestion(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}
