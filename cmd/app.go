package main

import (
	"github.com/gictorbit/sikabiz/db/userdb"
	"github.com/gictorbit/sikabiz/seeder"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "sikabiz",
		Usage: "sikabiz cmd",
		Commands: []*cli.Command{
			{
				Name:  "seeder",
				Usage: "runs seeder",
				Action: func(ctx *cli.Context) error {
					logger, err := zap.NewProduction()
					if err != nil {
						return err
					}
					filePath := os.Getenv("JSON_FILE_NAME")
					userdbConn, err := userdb.ConnectToUserDB(os.Getenv("USERDB_URL"))
					if err != nil {
						return err
					}
					userDatabase := userdb.NewUserDB(userdbConn)
					s := seeder.NewSeeder(userDatabase, logger, filePath)
					s.RunSeeder(ctx.Context)
					return nil
				},
			},
			{
				Name:  "api",
				Usage: "runs api",
				Action: func(ctx *cli.Context) error {

					return nil
				},
			},
		},
	}

	if e := app.Run(os.Args); e != nil {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("create new logger failed:%v\n", err)
		}
		logger.Error("failed to run app", zap.Error(e))
	}
}
