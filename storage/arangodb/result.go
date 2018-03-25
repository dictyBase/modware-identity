package arangodb

import (
	"github.com/dictyBase/apihelpers/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/identity"
)

type arangoResult struct {
	noResult bool
	doc      *identityDoc
	id       int64
}

func (a *arangoResult) NotFound() bool {
	return a.noResult
}

func (a *arangoResult) GetId() int64 {
	return a.id
}

func (a *arangoResult) GetAttributes() *identity.IdentityAttributes {
	return a.docToAttributes(a.doc)
}

func (a *arangoResult) docToAttributes(doc *identityDoc) *identity.IdentityAttributes {
	return &identity.IdentityAttributes{
		Identifier: doc.Identifier,
		Provider:   doc.Provider,
		UserId:     doc.UserId,
		CreatedAt:  aphgrpc.TimestampProto(doc.CreatedAt),
		UpdatedAt:  aphgrpc.TimestampProto(doc.UpdatedAt),
	}
}
