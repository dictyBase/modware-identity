package nats

import (
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/pubsub"
	"github.com/dictyBase/modware-identity/message"
	gnats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
)

type natsReply struct {
	econn *gnats.EncodedConn
}

func NewReply(host, port string, options ...gnats.Option) (message.Reply, error) {
	nc, err := gnats.Connect(fmt.Sprintf("nats://%s:%s", host, port), options...)
	if err != nil {
		return &natsReply{}, err
	}
	ec, err := gnats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		return &natsReply{}, err
	}
	return &natsReply{econn: ec}, nil
}

func (n *natsReply) Publish(subj string, rep *pubsub.IdentityReply) {
	n.econn.Publish(subj, rep)
}

func (n *natsReply) Start(subj string, client message.IdentityClient, replyFn message.ReplyFn) error {
	_, err := n.econn.Subscribe(subj, func(s, rep string, req *pubsub.IdentityReq) {
		n.Publish(rep, replyFn(s, client, req))
	})
	if err != nil {
		return err
	}
	if err := n.econn.Flush(); err != nil {
		return err
	}
	if err := n.econn.LastError(); err != nil {
		return err
	}
	return nil
}

func (n *natsReply) Stop() error {
	n.econn.Close()
	return nil
}
