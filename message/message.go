package message

import (
	"context"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/pubsub"
)

type Request interface {
	Request(string, *pubsub.IdRequest, time.Duration) (*pubsub.UserReply, error)
	RequestWithContext(context.Context, string, *pubsub.IdRequest) (*pubsub.UserReply, error)
}
