package operator

import (
	"fmt"
	"log"

	"github.com/anitta/eguchi-wedding-bot/pkg/domain/quiz"
	firebasesdk "github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/firebase"
	"github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/line"
)

type Operator interface {
    ThinkingTime() error
    CalculateScore(question []string) ([]quiz.UserResult, error)
    CalculateScoreOfQuestion(questions []string) ([]quiz.QuetionResult, error)
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

func (o *operator) CalculateScore(question []string) ([]quiz.UserResult, error) {
    results := map[string]int64{}
    for _, q := range question {
        ans, err := o.FirebaseApp.GetAnswerByQuestion(q)
        if ans == "" {
            continue
        }
        if err != nil {
            return nil, err
        }
        log.Println(fmt.Sprintf("answer: %s", ans))
        userids, err := o.FirebaseApp.GetUserByAnswerChoice(q, ans)
        log.Printf("equals %#v",userids)
        if err != nil {
            return nil, err
        }
        for _, userid := range userids {
            _, ok := results[userid]
            if !ok {
                results[userid] = 1
            } else {
                results[userid] += 1
            }
        }
        userids, err = o.FirebaseApp.GetUserNotEqualAnswerChoice(q, ans)
        if err != nil {
            return nil, err
        }
        log.Printf("not equals %#v",userids)
        for _, userid := range userids {
            _, ok := results[userid]
            if !ok {
                results[userid] = 0
            }
        }
    }
    var userResults []quiz.UserResult
    for userid, score := range results {
        username, err := o.LineBot.GetUserNameByUserID(userid)
        if err != nil {
            return nil, err
        }
        userResults = append(userResults, quiz.UserResult{
            UserId: userid,
            Name: username,
            Score: score,
        })
    }
    return userResults, nil
}


func (o *operator) CalculateScoreOfQuestion(questions []string) ([]quiz.QuetionResult, error) {
    var userResults []quiz.QuetionResult
    for _, question := range questions {
        for _, choice := range []string{"1", "2", "3"} {
            userids, err := o.FirebaseApp.GetUserByAnswerChoice(question, choice)
            if err != nil {
                return nil, err
            }
            userResults = append(userResults, quiz.QuetionResult{
                QuetionTitle: question,
                Choice: choice,
                Score: int64(len(userids)),
            })
        }
    }
    return userResults, nil
}
