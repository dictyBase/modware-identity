package message

import (
	"context"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	"github.com/dictyBase/go-genproto/dictybaseapis/pubsub"
)

type Request interface {
	Request(string, *jsonapi.IdRequest, time.Duration) (*pubsub.Reply, error)
	RequestWithContext(context.Context, string, *jsonapi.IdRequest) (*pubsub.Reply, error)
}
