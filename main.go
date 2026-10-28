package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/m1ll3r1337/worker-service/cmd"
	msentry "github.com/m1ll3r1337/worker-service/internal/app/monitor/sentry"
	"github.com/m1ll3r1337/worker-service/internal/pkg/constant"
)

func main() {
	app := &cli.App{
		Name:    constant.AppName,
		Version: constant.Version,
		Usage:   "Сервис обработки событий order.created",
		Commands: []*cli.Command{
			cmd.WebServer(),
			cmd.SubscribeOrderCreated(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-json",
				Usage: "Человеко-читаемый формат логов вместо JSON",
			},
		},
	}

	defer msentry.Flush()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
