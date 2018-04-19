package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dictyBase/go-genproto/dictybaseapis/pubsub"
	"github.com/dictyBase/modware-identity/message"
	gclient "github.com/dictyBase/modware-identity/message/grpc-client"
	"github.com/dictyBase/modware-identity/message/nats"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/urfave/cli.v1"
)

func shutdown(r message.Reply, logger *logrus.Entry) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	logger.Info("received kill signal")
	if err := r.Stop(); err != nil {
		logger.Fatalf("unable to close the subscription %s\n", err)
	}
	logger.Info("closed the connections gracefully")
}

func ReplyIdentity(subj string, c message.IdentityClient, req *pubsub.IdentityReq) *pubsub.IdentityReply {
	switch subj {
	case "IdentityService.Get":
		identity, err := c.Get(req.Id)
		if err != nil {
			st, _ := status.FromError(err)
			return &pubsub.IdentityReply{
				Status: st.Proto(),
				Exist:  false,
			}
		}
		return &pubsub.IdentityReply{
			Exist:    true,
			Identity: identity,
		}
	case "IdentityService.ExistIdentity":
		exist, err := c.ExistIdentity(req)
		if err != nil {
			st, _ := status.FromError(err)
			return &pubsub.IdentityReply{
				Status: st.Proto(),
				Exist:  exist,
			}
		}
		return &pubsub.IdentityReply{
			Exist: exist,
		}
	case "IdentityService.Exist":
		exist, err := c.Exist(req.Id)
		if err != nil {
			st, _ := status.FromError(err)
			return &pubsub.IdentityReply{
				Status: st.Proto(),
				Exist:  exist,
			}
		}
		return &pubsub.IdentityReply{
			Exist: exist,
		}
	case "IdentityService.Delete":
		deleted, err := c.Delete(req.Id)
		if err != nil {
			st, _ := status.FromError(err)
			return &pubsub.IdentityReply{
				Status: st.Proto(),
				Exist:  deleted,
			}
		}
		return &pubsub.IdentityReply{
			Exist: deleted,
		}
	case "IdentityService.GetIdentity":
		identity, err := c.GetByIdentity(req)
		if err != nil {
			st, _ := status.FromError(err)
			return &pubsub.IdentityReply{
				Status: st.Proto(),
				Exist:  false,
			}
		}
		return &pubsub.IdentityReply{
			Exist:    true,
			Identity: identity,
		}

	default:
		return &pubsub.IdentityReply{
			Status: status.Newf(codes.Internal, "subject %s is not supported", subj).Proto(),
		}
	}
}

func RunIdentityReply(c *cli.Context) error {
	reply, err := nats.NewReply(
		c.String("messaging-host"),
		c.String("messaging-port"),
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("cannot connect to reply server %s", err),
			2,
		)
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", c.String("identity-grpc-host"), c.String("identity-grpc-port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("cannot connect to grpc server for identity microservice %s", err),
			2,
		)
	}
	err = reply.Start(
		"IdentityService.*",
		gclient.NewIdentityClient(conn),
		ReplyIdentity,
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("cannot start the reply server %s", err),
			2,
		)
	}
	logger := getLogger(c)
	logger.Info("starting the identity reply messaging backend")
	shutdown(reply, logger)
	return nil
}
