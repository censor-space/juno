package server

import (
	"log"
    "fmt"
	"net/http"

	"github.com/anitta/eguchi-wedding-bot/pkg/application/operator"
	"github.com/anitta/eguchi-wedding-bot/pkg/domain/quiz"
	firebasesdk "github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/firebase"
	"github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/line"
	"github.com/gin-gonic/gin"
	linebotsdk "github.com/line/line-bot-sdk-go/linebot"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Controller interface {
	HealthCheck(ctx *gin.Context)
	Metrics() gin.HandlerFunc
	PostQuestion(ctx *gin.Context)
    UpdateClearCurrentQuestion(ctx *gin.Context)
    GetUserScore(ctx *gin.Context)
	CallbackFromLine(ctx *gin.Context)
}

type controller struct {
	LineBot     line.LineBot
	FirebaseApp firebasesdk.FirebaseApp
    Operator operator.Operator
}

func NewController(lineBot line.LineBot, firebaseApp firebasesdk.FirebaseApp, operator operator.Operator) Controller {
	return &controller{
		LineBot: lineBot,
        FirebaseApp: firebaseApp,
        Operator: operator,
	}
}

func (c *controller) HealthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "200 Status OK")
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

    err = c.FirebaseApp.SetQuestion(jsonQuestion)
    if err != nil {
        ctx.String(http.StatusInternalServerError, "500 Internal Server Error")
        return
    }

	err = c.LineBot.PostQuiz(jsonQuestion)
	if err != nil {
        ctx.String(http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	err = c.FirebaseApp.SetCurrentQuestionTitle(jsonQuestion.Title)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	ctx.String(http.StatusOK, "200 Status OK")
    go c.Operator.ThinkingTime()
	return
}

func (c *controller) UpdateClearCurrentQuestion(ctx *gin.Context) {
    err := c.FirebaseApp.SetCurrentQuestionTitle("")
    if err != nil {
        ctx.String(http.StatusInternalServerError, "500 Internal Server Error")
        return
    }
    ctx.String(http.StatusOK, "200 Status OK")
}

func (c *controller) GetUserScore(ctx *gin.Context) {
    var json quiz.Quetions
    err := ctx.ShouldBindJSON(&json)
    if err != nil {
		ctx.String(http.StatusBadRequest, "400 Bad Request")
		return
	}
    log.Println(fmt.Sprintf("%#v",json))
    userResult, err := c.Operator.CalculateScore(json.Titles)
    if err != nil {
        ctx.String(http.StatusInternalServerError, "500 Internal Server Error")
        return
    }
    ctx.JSON(http.StatusOK, userResult)
}

func (c *controller) CallbackFromLine(ctx *gin.Context) {
    events, err := c.LineBot.ParseLineEventRequest(ctx.Request)
		if err != nil {
			if err == linebotsdk.ErrInvalidSignature {
                ctx.String(http.StatusBadRequest, "400 Bad Request")
			} else {
                ctx.String(http.StatusInternalServerError, "500 Internal Server Error")
			}
			return
		}
		for _, event := range events {
			if event.Type == linebotsdk.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebotsdk.TextMessage:
                    log.Println("linebotsdk.TextMessage")
                    title, err := c.FirebaseApp.GetCurrentQuestionTitle()
                    if err != nil {
                        log.Println("500: GetCurrentQuestionTitle()")
                        ctx.String(http.StatusInternalServerError, "500 Internal Server Error")
                        return
                    }
                    replyMessage := "現在は回答を受付しておりません。"
                    if title != "" {
                        err = c.FirebaseApp.SetUserAnswer(event.Source.UserID, title, quiz.Answer{
                            Answer: message.Text,
                            ID: event.Source.UserID,
                        })
                        if err != nil {
                            log.Println("500: SetUserAnswer()")
                            ctx.String(http.StatusInternalServerError, "500 Internal Server Error")
                            return
					    }
                        replyMessage = fmt.Sprintf("%sを選択しました。", message.Text)
                    }
                    err = c.LineBot.PostReplyMessage(event.ReplyToken, replyMessage)
                    if err != nil {
                        log.Println("500: PostReplyMessage()")
                        ctx.String(http.StatusInternalServerError, "500 Internal Server Error")
                        return
					}
                    ctx.String(http.StatusOK, "200 Status OK")
				}
			}
		}
}

