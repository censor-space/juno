package line

import (
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
    PostMessageToUserID(userid, messageText string) error
    GetUserNameByUserID(userid string) (string, error)
}

type lineBot struct {
	Client *linebotsdk.Client
}

func NewLineBot(channelSecret, channelToken string) (LineBot, error) {
	client, err := linebotsdk.New(channelSecret, channelToken)
	if err != nil {
		return nil, err
	}
	return &lineBot{
		Client: client,
	}, nil
}

func (lb *lineBot) PostQuiz(question quiz.Question) error {
	log.Println("Post Quiz")

	message := createQuizTemplateMessage(question)

    log.Printf("%#v", message)
	// append some message to messages
	_, err := lb.Client.BroadcastMessage(message).Do()
	if err != nil {
		// Do something when some bad happened
		log.Println("Do something when some bad happened.")
        log.Println(err)
		return err
	}

	return nil
}

func createQuizTemplateMessage(question quiz.Question) *linebotsdk.TemplateMessage {
	template := linebotsdk.NewButtonsTemplate(
		question.ImageURL,
		question.Title,
		question.Text,
        linebotsdk.NewMessageAction("1", question.Choice1),
		linebotsdk.NewMessageAction("2", question.Choice2),
		linebotsdk.NewMessageAction("3", question.Choice3),
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


func (lb *lineBot) PostMessageToUserID(userid, messageText string) error {
    _, err := lb.Client.PushMessage(userid, linebotsdk.NewTextMessage(messageText)).Do()
    return err
}

func (lb*lineBot) GetUserNameByUserID(userid string) (string, error) {
    profile, err := lb.Client.GetProfile(userid).Do()
    if err != nil {
        return "", err
    }
    return profile.DisplayName, nil
}
