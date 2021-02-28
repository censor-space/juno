package di

import (
	"context"
	"log"

	"github.com/anitta/eguchi-wedding-bot/pkg/application/server"
	"github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/config"
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
	return server.NewController(di.LineBot())
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

func (di *DI) ReadConfig() (*config.Environment, error) {
	return config.Get()
}
