package validate

import (
	"fmt"

	"github.com/urfave/cli"
)

func ValidateCreateIdentityArgs(c *cli.Context) error {
	for _, p := range []string{
		"identity-grpc-host",
		"user-grpc-host",
		"user-grpc-port",
		"identity-grpc-port",
		"identifier",
		"provider",
		"email",
	} {
		if len(c.String(p)) == 0 {
			return cli.NewExitError(
				fmt.Sprintf("argument %s is missing", p),
				2,
			)
		}
	}
	return nil
}

func ValidateReplyArgs(c *cli.Context) error {
	for _, p := range []string{
		"identity-grpc-host",
		"identity-grpc-port",
		"messaging-host",
		"messaging-port",
	} {
		if len(c.String(p)) == 0 {
			return cli.NewExitError(
				fmt.Sprintf("argument %s is missing", p),
				2,
			)
		}
	}
	return nil
}

func ValidateServerArgs(c *cli.Context) error {
	for _, p := range []string{
		"arangodb-pass",
		"arangodb-database",
		"arangodb-user",
		"nats-host",
		"nats-port",
	} {
		if len(c.String(p)) == 0 {
			return cli.NewExitError(
				fmt.Sprintf("argument %s is missing", p),
				2,
			)
		}
	}
	return nil
}
