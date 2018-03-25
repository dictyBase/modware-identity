package nats

import (
	"context"
	"time"

	gnats "github.com/nats-io/go-nats"

	"github.com/dictyBase/apihelpers/pubsub"
	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	"github.com/nats-io/go-nats/encoders/protobuf"
)

type natsRequest struct {
	econn *gnats.Conn
}

func NewRequest(url string, options ...gnats.Option) (pubsub.Request, error) {
	nc, err := gnats.Connect(url, options...)
	if err != nil {
		return &natsRequest{}, err
	}
	ec, err := gnats.NewEncodedConn(c, protobuf.PROTOBUF_ENCODER)
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
