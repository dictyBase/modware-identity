package message

import (
	"context"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/identity"
	"github.com/dictyBase/go-genproto/dictybaseapis/pubsub"
)

type IdentityClient interface {
	Get(int64) (*identity.Identity, error)
	GetByIdentity(*pubsub.IdentityReq) (*identity.Identity, error)
	Delete(int64) (bool, error)
	Exist(int64) (bool, error)
	ExistIdentity(*pubsub.IdentityReq) (bool, error)
}

type Request interface {
	UserRequest(string, *pubsub.IdRequest, time.Duration) (*pubsub.UserReply, error)
	UserRequestWithContext(context.Context, string, *pubsub.IdRequest) (*pubsub.UserReply, error)
	IsActive() bool
}

type ReplyFn func(string, IdentityClient, *pubsub.IdentityReq) *pubsub.IdentityReply

type Reply interface {
	Publish(string, *pubsub.IdentityReply)
	Start(string, IdentityClient, ReplyFn) error
	Stop() error
}
