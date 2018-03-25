package storage

import (
	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	"github.com/dictyBase/go-genproto/dictybaseapis/identity"
)

type DataSource interface {
	HasProviderIdentity(*identity.IdentityProviderReq) (bool, error)
	HasIdentity(*jsonapi.IdRequest) (bool, error)
	GetIdentity(*jsonapi.IdRequest) (Result, error)
	GetProviderIdentity(*identity.IdentityProviderReq) (Result, error)
	CreateIdentity(*identity.NewIdentityAttributes) (Result, error)
	DeleteIdentity(*jsonapi.IdRequest) (bool, error)
}

type Result interface {
	NotFound() bool
	GetId() int64
	GetAttributes() *identity.IdentityAttributes
}
