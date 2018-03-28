package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/pubsub"
	"github.com/dictyBase/modware-identity/message"
	gnats "github.com/nats-io/go-nats"

	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	"github.com/nats-io/go-nats/encoders/protobuf"
)

type natsRequest struct {
	econn *gnats.EncodedConn
}

func NewRequest(host, port string, options ...gnats.Option) (message.Request, error) {
	nc, err := gnats.Connect(fmt.Sprintf("nats://%s:%s", host, port), options...)
	if err != nil {
		return &natsRequest{}, err
	}
	ec, err := gnats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		return &natsRequest{}, err
	}
	return &natsRequest{econn: ec}, nil
}

func (n *natsRequest) Request(subj string, r *jsonapi.IdRequest, timeout time.Duration) (*pubsub.Reply, error) {
	reply := &pubsub.Reply{}
	err := n.econn.Request(subj, r, reply, timeout)
	return reply, err
}

func (n *natsRequest) RequestWithContext(ctx context.Context, subj string, r *jsonapi.IdRequest) (*pubsub.Reply, error) {
	reply := &pubsub.Reply{}
	err := n.econn.RequestWithContext(ctx, subj, r, reply)
	return reply, err
}
