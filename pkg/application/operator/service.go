package operator

import (
	"fmt"
	"log"
    "sort"

	"github.com/anitta/eguchi-wedding-bot/pkg/domain/quiz"
	firebasesdk "github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/firebase"
	"github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/line"
)

type Operator interface {
    ThinkingTime() error
    CalculateScore(question []string) ([]quiz.UserResult, error)
    CalculateScoreOfQuestion(questions []string) ([]quiz.QuetionResult, error)
    PostCalculateScoreToUser(questions []string) error
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

func (o *operator) PostCalculateScoreToUser(questions []string) error {
    quizResults, err := o.CalculateScore(questions)
    if err != nil {
        return err
    }
    return o.calculateRankingOfUserResult(quizResults)
}

func (o *operator) calculateRankingOfUserResult(quizUserResults []quiz.UserResult) error {
    qurs := quiz.UserResults(quizUserResults)
    log.Println(qurs)
    sort.Sort(sort.Reverse(qurs))                                 // ユーザのスコア順にソートする。同一スコアの場合、順番は考慮しない。
    log.Println(qurs)
    var currentRanking int
    var currentScore int64
    for k, userResult := range qurs {
        if k == 0 {                                 // ランキング1位の時
            currentRanking = 1                      // 現在のランキングに1位で初期化
            currentScore = userResult.Score         // ランキング1位のスコアで初期化
        } else if currentScore > userResult.Score { // ユーザのランキングが一つ前の結果より低い時
            currentScore = userResult.Score         // 現在のユーザのスコアに更新
            currentRanking = k + 1                  // ランキング更新
        }
        err := o.sendRankingAndScoreMessage(currentRanking, qurs.Len(), currentScore, userResult.UserId)
        if err != nil {
            log.Printf("score cannot send. username: %s, lineid: %s", userResult.Name, userResult.UserId)
        }
    }
    return nil
}

func (o *operator) sendRankingAndScoreMessage(ranking, answererCount int, score int64,userid string) error {
    return o.LineBot.PostMessageToUserID(userid, fmt.Sprintf("あなたの順位は %d / %d\nあなたのスコアは %d 問 正解です", ranking, answererCount, score))
}
