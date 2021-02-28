package line

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineBot interface {
	PostQuiz() error
}

type lineBot struct {
	Client *linebot.Client
}

func NewLineBot(channelSecret, channelToken string) (LineBot, error) {
	client, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		log.Println("NewLineBot err")
		return nil, err
	}
	return &lineBot{
		Client: client,
	}, nil
}

func (lb *lineBot) PostQuiz() error {
	log.Println("Post Quiz")
	return nil
}
