package firebasesdk

import (
	"context"

	firebase "firebase.google.com/go"
	"github.com/anitta/eguchi-wedding-bot/pkg/domain/quiz"
	"google.golang.org/api/option"
)

type FirebaseApp interface {
	SetCurrentQuestionTitle(questionTitle string) error
    GetCurrentQuestionTitle() (string, error)
    SetUserAnswer(userID, questionTitle string, userAnswer quiz.Answer) error
    SetQuestion(question quiz.Question) error
}

type firebaseApp struct {
	App *firebase.App
    Ctx context.Context
}

func NewFirebaseApp(ctx context.Context, credentialsFilePath, databaseURL string) (FirebaseApp, error) {
	opt := option.WithCredentialsFile(credentialsFilePath)
	config := &firebase.Config{DatabaseURL: databaseURL}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, err
	}
	return &firebaseApp{
		App: app,
        Ctx: ctx,
	}, nil
}

func (fa *firebaseApp) SetCurrentQuestionTitle(questionTitle string) error {
	client, err := fa.App.Firestore(fa.Ctx)
	if err != nil {
		return err
	}
	_, err = client.Collection("question").Doc("current").Set(fa.Ctx, map[string]interface{}{
        "title": questionTitle,
    })
	if err != nil {
		return err
	}
	return nil
}


func (fa *firebaseApp) GetCurrentQuestionTitle() (string, error) {
	client, err := fa.App.Firestore(fa.Ctx)
	if err != nil {
		return "", err
	}
    doc, err := client.Collection("question").Doc("current").Get(fa.Ctx)
    if err != nil {
        return "", err
    }
    title := doc.Data()
	return title["title"].(string) , nil
}


func (fa *firebaseApp) SetUserAnswer(userID, questionTitle string, userAnswer quiz.Answer) error {
	client, err := fa.App.Firestore(fa.Ctx)
	if err != nil {
		return err
	}
    _, err = client.Collection("user-answer").Doc(questionTitle).Collection("userid").Doc(userID).Set(fa.Ctx, userAnswer)
	if err != nil {
		return err
	}
	return nil
}

func (fa *firebaseApp) SetQuestion(question quiz.Question) error {
	client, err := fa.App.Firestore(fa.Ctx)
	if err != nil {
		return err
	}
	_, err = client.Collection("question").Doc(question.Title).Set(fa.Ctx, question)
	if err != nil {
		return err
	}
	return nil
}
