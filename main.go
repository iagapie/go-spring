package main

import (
	"github.com/iagapie/go-spring/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var cfgFiles = []string{
	"./configs/app",
	"./configs/cms",
	"./configs/jwt",
}

func main() {
	app := cli.NewApp()
	app.Name = "Spring CMS"
	app.Version = "1.0.0"
	app.Commands = []*cli.Command{
		cmd.Web,
	}

	defaultFlags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   cli.NewStringSlice(cfgFiles...),
			Usage:   "Custom configuration files",
		},
	}

	app.Flags = append(app.Flags, cmd.Web.Flags...)
	app.Flags = append(app.Flags, defaultFlags...)
	app.Action = cmd.Web.Action

	for i := range app.Commands {
		setFlagsOnSubcommands(app.Commands[i], defaultFlags)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setFlagsOnSubcommands(command *cli.Command, defaultFlags []cli.Flag) {
	command.Flags = append(command.Flags, defaultFlags...)
	for i := range command.Subcommands {
		setFlagsOnSubcommands(command.Subcommands[i], defaultFlags)
	}
}
