package di

import (
	"context"

	"github.com/anitta/eguchi-wedding-bot/pkg/application/server"
	"github.com/anitta/eguchi-wedding-bot/pkg/infrastructure/config"
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
	return server.NewController()
}

func (di *DI) ReadConfig() (*config.Environment, error) {
	return config.Get()
}
