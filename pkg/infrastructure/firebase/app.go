package firebasesdk

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type FirebaseApp interface {
	SetCurrentQuestionTitle(questionTitle string) error
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
        log.Println("ocha")
		return nil, err
	}
	return &firebaseApp{
		App: app,
        Ctx: ctx,
	}, nil
}

func (fa *firebaseApp) SetCurrentQuestionTitle(questionTitle string) error {
	client, err := fa.App.Database(fa.Ctx)
	if err != nil {
		return err
	}
	ref := client.NewRef("server/current")
	questionRef := ref.Child("question")

	err = questionRef.Set(fa.Ctx, map[string]string{
		"title": questionTitle,
	})
	if err != nil {
		return err
	}
	return nil
}
