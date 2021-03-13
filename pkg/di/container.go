package di

import (
	"context"
	"log"
	"time"

	"github.com/anitta/eguchi-wedding-bot/pkg/application/operator"
	"github.com/anitta/eguchi-wedding-bot/pkg/application/server"
	"github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/config"
	firebasesdk "github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/firebase"
	"github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/line"
)

type DI struct{}

func Start(ctx context.Context) error {
	di := DI{}
	env, err := di.ReadConfig()
	if err != nil {
		panic(err)
	}
	return server.Run(ctx, env.Port, di.Controller())
}

func (di *DI) Controller() server.Controller {
	return server.NewController(di.LineBot(), di.FirebaseApp(), di.Operator())
}

func (di *DI) LineBot() line.LineBot {
	env, err := di.ReadConfig()
	if err != nil {
		panic(err)
	}
	bot, err := line.NewLineBot(env.LineChannelSecret, env.LineChannelToken)
	if err != nil {
		log.Println("LineBot()")
		panic(err)
	}
	return bot
}

func (di *DI) FirebaseApp() firebasesdk.FirebaseApp {
	env, err := di.ReadConfig()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	app, err := firebasesdk.NewFirebaseApp(ctx, env.FirebaseCredentialsFilePath, env.FirebaseDatabaseURL)
	if err != nil {
		log.Println("FirebaseApp()")
		panic(err)
	}
	return app
}

func (di *DI) Operator() operator.Operator {
	env, err := di.ReadConfig()
	if err != nil {
        log.Println("Operator()")
		panic(err)
	}
    log.Println(env.ThinkingTimeSec)
    return operator.NewOperator(func(){
        time.Sleep((time.Duration(1*env.ThinkingTimeSec)* time.Second))
    },
    di.FirebaseApp(),
    di.LineBot())
}

func (di *DI) ReadConfig() (*config.Environment, error) {
	return config.Get()
}
