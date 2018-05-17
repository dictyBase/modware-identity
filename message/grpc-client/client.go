package client

import (
	"context"

	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	"github.com/dictyBase/go-genproto/dictybaseapis/identity"
	"github.com/dictyBase/go-genproto/dictybaseapis/pubsub"
	"github.com/dictyBase/modware-identity/message"
	"google.golang.org/grpc"
)

type grpcIdentityClient struct {
	client identity.IdentityServiceClient
}

func NewIdentityClient(conn *grpc.ClientConn) message.IdentityClient {
	return &grpcIdentityClient{
		client: identity.NewIdentityServiceClient(conn),
	}
}

func (g *grpcIdentityClient) Get(id int64) (*identity.Identity, error) {
	return g.client.GetIdentity(context.Background(), &jsonapi.IdRequest{Id: id})
}

func (g *grpcIdentityClient) GetByIdentity(r *pubsub.IdentityReq) (*identity.Identity, error) {
	return g.client.GetIdentityFromProvider(
		context.Background(),
		&identity.IdentityProviderReq{
			Provider:   r.Provider,
			Identifier: r.Identifier,
		})
}

func (g *grpcIdentityClient) Delete(id int64) (bool, error) {
	_, err := g.client.DeleteIdentity(context.Background(), &jsonapi.IdRequest{Id: id})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g *grpcIdentityClient) Exist(id int64) (bool, error) {
	identity, err := g.Get(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g *grpcIdentityClient) ExistIdentity(r *pubsub.IdentityReq) (bool, error) {
	res, err := g.client.ExistProviderIdentity(
		context.Background(),
		&identity.IdentityProviderReq{
			Provider:   r.Provider,
			Identifier: r.Identifier,
		})
	return res.Exist, err
}
