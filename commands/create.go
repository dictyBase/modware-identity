package commands

import (
	"context"
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	"github.com/dictyBase/go-genproto/dictybaseapis/identity"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v1"
)

func CreateIdentity(c *cli.Context) error {
	uconn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", c.String("user-grpc-host"), c.String("user-grpc-port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("cannot connect to grpc server for user microservice %s", err),
			2,
		)
	}
	idconn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", c.String("identity-grpc-host"), c.String("identity-grpc-port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("cannot connect to grpc server for identity microservice %s", err),
			2,
		)
	}
	uclient := user.NewUserServiceClient(uconn)
	res, err := uclient.GetUserByEmail(context.Background(), &jsonapi.GetEmailRequest{Email: c.String("email")})
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in retrieving user %s", c.Int64("email"), err),
			2,
		)
	}
	idclient := identity.NewIdentityServiceClient(idconn)
	idn, err := idclient.CreateIdentity(
		context.Background(),
		&identity.CreateIdentityReq{
			Data: &identity.CreateIdentityReq_Data{
				Type: "identities",
				Attributes: &identity.NewIdentityAttributes{
					Identifier: c.String("identifier"),
					Provider:   c.String("provider"),
					UserId:     res.Data.Id,
				},
			},
		})
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in creating identity %s", err),
			2,
		)
	}
	logger := getLogger(c)
	logger.Infof(
		"created identity %d with identifier %s provider %s and user %s",
		idn.Data.Id,
		c.String("identifier"),
		c.String("provider"),
		c.String("email"),
	)
	return nil
}
