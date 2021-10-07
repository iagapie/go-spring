package cmd

import (
	"github.com/iagapie/go-spring/modules/backend/user"
	userdb "github.com/iagapie/go-spring/modules/backend/user/db"
	"github.com/iagapie/go-spring/modules/sys/config"
	"github.com/iagapie/go-spring/modules/sys/helper"
	"github.com/iagapie/go-spring/modules/sys/logger"
	"github.com/iagapie/go-spring/modules/sys/password"
	"github.com/iagapie/go-spring/modules/sys/postgresdb"
	"github.com/urfave/cli/v2"
)

type globalData struct {
	cfg         config.Cfg
	log         *logger.Logger
	encoder     password.Encoder
	db          *postgresdb.Database
	userStorage user.Storage
	userService user.Service
}

func initGlobalData(ctx *cli.Context) (*globalData, error) {
	var cfg config.Cfg
	if err := helper.ReadConfig(&cfg, ctx.StringSlice("config")...); err != nil {
		return nil, err
	}

	log := logger.New(logger.WithDebug(cfg.App.Debug))
	log.Infoln("config and logger initialized")

	log.Infoln("password encoder initializing")
	encoder := password.NewDefault()

	log.Infoln("database initializing")
	postgres, err := postgresdb.New(cfg.DB, log.Entry)
	if err != nil {
		return nil, err
	}

	log.Infoln("auto migrate")
	if err = postgres.AutoMigrate(&user.User{}); err != nil {
		postgres.Close()
		return nil, err
	}

	log.Infoln("user storage initializing")
	userStorage := userdb.NewStorage(postgres, log.Entry)

	log.Infoln("user service initializing")
	userService := user.NewService(userStorage, encoder)

	return &globalData{
		cfg:         cfg,
		log:         log,
		encoder:     encoder,
		db:          postgres,
		userStorage: userStorage,
		userService: userService,
	}, nil
}
