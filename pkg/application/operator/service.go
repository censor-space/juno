package operator

import (
	"log"

	firebasesdk "github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/firebase"
    "github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/line"
)

type Operator interface {
    ThinkingTime() error
    CalculateScore() error
}

type operator struct {
    ThinkingTimeFunc func()
    FirebaseApp firebasesdk.FirebaseApp
    LineBot     line.LineBot
}

func NewOperator(thinkingTimeFunc func(), firebaseApp firebasesdk.FirebaseApp, lineBot line.LineBot) Operator {
   return &operator{
        ThinkingTimeFunc: thinkingTimeFunc,
        FirebaseApp: firebaseApp,
        LineBot: lineBot,
   }
}


func (o *operator) ThinkingTime() error {
    log.Println("o-i")
    o.ThinkingTimeFunc()
    log.Println("ocha")
    q, err := o.FirebaseApp.GetCurrentQuestionTitle()
    if err != nil {
        return err
    }
    if q == "" {
        return nil
    }
    err = o.FirebaseApp.SetCurrentQuestionTitle("")
    if err != nil {
        return err
    }
    err = o.LineBot.PostMessage("回答は締め切りました。")
    if err != nil {
        return err
    }
    return nil
}

func (o *operator) CalculateScore() error {
    return nil
}
