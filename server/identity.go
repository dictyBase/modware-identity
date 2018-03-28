package server

import (
	"context"

	"github.com/dictyBase/apihelpers/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	"github.com/dictyBase/go-genproto/dictybaseapis/identity"
	"github.com/dictyBase/modware-identity/message"
	"github.com/dictyBase/modware-identity/storage"
	"github.com/golang/protobuf/ptypes/empty"
)

type IdentityService struct {
	*aphgrpc.Service
	storage storage.DataSource
	request message.Request
}

func defaultOptions() *aphgrpc.ServiceOptions {
	return &aphgrpc.ServiceOptions{
		PathPrefix: "identities",
		Resource:   "identities",
	}
}

func NewIdentityService(st storage.DataSource, r message.Request, opt ...aphgrpc.Option) *IdentityService {
	so := defaultOptions()
	for _, optfn := range opt {
		optfn(so)
	}
	srv := &aphgrpc.Service{}
	aphgrpc.AssignFieldsToStructs(so, srv)
	return &IdentityService{
		Service: srv,
		storage: st,
		request: r,
	}
}

func (s *IdentityService) GetIdentityFromProvider(ctx context.Context, r *identity.IdentityProviderReq) (*identity.Identity, error) {
	rs, err := s.storage.GetProviderIdentity(r)
	if err != nil {
		return &identity.Identity{}, aphgrpc.HandleGetError(ctx, err)
	}
	if rs.NotFound() {
		return &identity.Identity{}, aphgrpc.HandleNotFoundError(ctx, err)
	}
	return s.buildResource(rs.GetId(), rs.GetAttributes()), nil
}

func (s *IdentityService) GetIdentity(ctx context.Context, r *jsonapi.IdRequest) (*identity.Identity, error) {
	return &identity.Identity{}, nil
}

func (s *IdentityService) CreateIdentity(ctx context.Context, r *identity.CreateIdentityReq) (*identity.Identity, error) {
	return &identity.Identity{}, nil
}

func (s *IdentityService) DeleteIdentity(ctx context.Context, r *jsonapi.IdRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *IdentityService) Healthz(ctx context.Context, r jsonapi.HealthzIdRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// -- Functions that builds up the various parts of the final user resource objects
func (s *IdentityService) buildResourceData(id int64, attr *identity.IdentityAttributes) *identity.IdentityData {
	return &identity.IdentityData{
		Attributes: attr,
		Id:         id,
		Type:       s.GetResourceName(),
		Links: &jsonapi.Links{
			Self: s.GenResourceSelfLink(id),
		},
	}
}

func (s *IdentityService) buildResource(id int64, attr *identity.IdentityAttributes) *identity.Identity {
	return &identity.Identity{
		Data: s.buildResourceData(id, attr),
		Links: &jsonapi.Links{
			Self: s.GenResourceSelfLink(id),
		},
	}
}
