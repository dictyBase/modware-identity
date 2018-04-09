package message

import (
	"context"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/pubsub"
)

type IdentityClient interface {
	Get(int64) (*idenity.Identity, error)
	Delete(int64) (bool, error)
	Exist(int64) (bool, error)
}

type Request interface {
	UserRequest(string, *pubsub.IdRequest, time.Duration) (*pubsub.UserReply, error)
	UserRequestWithContext(context.Context, string, *pubsub.IdRequest) (*pubsub.UserReply, error)
}

type ReplyFn func(string, IdentityClient, *pubsub.IdRequest) *pubsub.IdentityReply

type Reply interface {
	Publish(string, *pubsub.IdentityReply)
	Start(string, IdentityClient, ReplyFn) error
	Stop() error
}
