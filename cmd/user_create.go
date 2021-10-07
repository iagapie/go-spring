package cmd

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/iagapie/go-spring/modules/backend/user"
	"github.com/urfave/cli/v2"
)

var UserCreate = &cli.Command{
	Name:   "user:create",
	Usage:  "Create backend user",
	Action: runUserCreate,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "User name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "email",
			Aliases:  []string{"e"},
			Usage:    "User email",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Usage:    "User password",
			Required: true,
		},
	},
}

func runUserCreate(ctx *cli.Context) error {
	data, err := initData(ctx)
	if err != nil {
		return err
	}
	defer data.db.Close()

	dto := user.CreateUserDTO{
		Name:           ctx.String("name"),
		Email:          ctx.String("email"),
		Password:       ctx.String("password"),
		RepeatPassword: ctx.String("password"),
	}

	if err = validator.New().Struct(&dto); err != nil {
		return err
	}

	id, err := data.userService.Create(context.Background(), dto)
	if err != nil {
		return err
	}

	data.log.Infof("user %s <%s> was created with UUID %s", dto.Name, dto.Email, id)
	return nil
}
