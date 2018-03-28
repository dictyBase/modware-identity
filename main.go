package main

import (
	"os"

	"github.com/dictyBase/modware-identity/commands"
	"github.com/dictyBase/modware-identity/validate"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "modware-identity"
	app.Usage = "cli for modware-identity microservice"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-format",
			Usage: "format of the logging out, either of json or text.",
			Value: "json",
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "log level for the application",
			Value: "error",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "start-server",
			Usage:  "starts the modware-identity microservice with HTTP and grpc backends",
			Action: commands.RunServer,
			Before: validate.ValidateServerArgs,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "arangodb-pass, pass",
					EnvVar: "ARANGODB_PASS",
					Usage:  "arangodb database password",
				},
				cli.StringFlag{
					Name:   "arangodb-database, db",
					EnvVar: "ARANGODB_DATABASE",
					Usage:  "arangodb database name",
				},
				cli.StringFlag{
					Name:   "arangodb-user, user",
					EnvVar: "ARANGODB_USER",
					Usage:  "arangodb database user",
				},
				cli.StringFlag{
					Name:   "arangodb-host, host",
					Value:  "arangodb-backend",
					EnvVar: "ARANGODB_BACKEND_SERVICE_HOST",
					Usage:  "arangodb database host",
				},
				cli.StringFlag{
					Name:   "arangodb-port",
					EnvVar: "ARAGODB_BACKEND_SERVICE_PORT",
					Usage:  "arangodb database port",
					Value:  "8529",
				},
				cli.StringFlag{
					Name:   "identity-api-http-host",
					EnvVar: "IDENTITY_API_HTTP_HOST",
					Usage:  "public hostname serving the http api, by default the default port will be appended to http://localhost",
				},
				cli.StringFlag{
					Name:   "nats-host",
					EnvVar: "NATS_SERVICE_HOST",
					Usage:  "nats messaging server host",
				},
				cli.StringFlag{
					Name:   "nats-port",
					EnvVar: "NATS_SERVICE_PORT",
					Usage:  "nats messaging server port",
				},
				cli.StringFlag{
					Name:  "port",
					Usage: "tcp port at which the server will be available",
					Value: "9560",
				},
			},
		},
	}
	app.Run(os.Args)
}
