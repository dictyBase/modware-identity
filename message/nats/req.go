package nats

import (
	"context"
	"time"

	"github.com/dictyBase/modware-identity/message"
	gnats "github.com/nats-io/go-nats"

	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	"github.com/nats-io/go-nats/encoders/protobuf"
)

type natsRequest struct {
	econn *gnats.Conn
}

func NewRequest(url string, options ...gnats.Option) (message.Request, error) {
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

func (n *natsRequest) Request(subj string, r *jsonapi.IdRequest, timeout time.Duration) (*message.Reply, error) {
	reply := &message.Reply{}
	err := n.econn.Request(subj, r, reply, timeout)
	return reply, err
}

func (n *natsRequest) RequestWithContext(ctx context.Context, subj string, r *jsonapi.IdRequest) (*message.Reply, error) {
	reply := &message.Reply{}
	err := n.econn.RequestWithContext(ctx, subj, r, reply)
	return reply, err
}
