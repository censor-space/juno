package line

import (
    "fmt"
	"log"
    "net/http"

	"github.com/anitta/eguchi-wedding-bot/pkg/domain/quiz"
	linebotsdk "github.com/line/line-bot-sdk-go/linebot"
)

type LineBot interface {
	PostQuiz(question quiz.Question) error
    ParseLineEventRequest(req *http.Request) ([]*linebotsdk.Event, error)
    PostReplyMessage(eplyToken, messageText string) error
    PostMessage(messageText string) error
}

type lineBot struct {
	Client *linebotsdk.Client
}

func NewLineBot(channelSecret, channelToken string) (LineBot, error) {
	client, err := linebotsdk.New(channelSecret, channelToken)
	if err != nil {
		log.Println("NewLineBot err")
		return nil, err
	}
	return &lineBot{
		Client: client,
	}, nil
}

func (lb *lineBot) PostQuiz(question quiz.Question) error {
	log.Println("Post Quiz")

	message := createQuizTemplateMessage(question)

	// append some message to messages
	_, err := lb.Client.BroadcastMessage(message).Do()
	if err != nil {
		// Do something when some bad happened
		log.Println("Do something when some bad happened.")
		return err
	}

	return nil
}

func createQuizTemplateMessage(question quiz.Question) *linebotsdk.TemplateMessage {
	template := linebotsdk.NewButtonsTemplate(
		question.ImageURL,
		question.Title,
		question.Text,
        linebotsdk.NewMessageAction(fmt.Sprintf("A: %s",question.Choice1), question.Choice1),
		linebotsdk.NewMessageAction(fmt.Sprintf("B: %s",question.Choice2), question.Choice2),
		linebotsdk.NewMessageAction(fmt.Sprintf("C: %s",question.Choice3), question.Choice3),
		linebotsdk.NewMessageAction(fmt.Sprintf("D: %s",question.Choice4), question.Choice4),
	)

	return linebotsdk.NewTemplateMessage(question.NotificationMessage, template)
}

func (lb *lineBot) ParseLineEventRequest(req *http.Request) ([]*linebotsdk.Event, error) {
    return lb.Client.ParseRequest(req)
}

func (lb *lineBot) PostReplyMessage(replyToken, messageText string) error {
    _, err := lb.Client.ReplyMessage(replyToken, linebotsdk.NewTextMessage(messageText)).Do();
    return err
}


func (lb *lineBot) PostMessage(messageText string) error {
    _, err := lb.Client.BroadcastMessage(linebotsdk.NewTextMessage(messageText)).Do()
    return err
}
