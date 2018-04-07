package arangodb

import (
	"context"
	"log"
	"os"
	"testing"

	driver "github.com/arangodb/go-driver"
	"github.com/dictyBase/apihelpers/aphdocker"
	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	"github.com/dictyBase/go-genproto/dictybaseapis/identity"
	"github.com/dictyBase/modware-identity/storage"
)

var ahost, aport, auser, apass, adb string
var coll driver.Collection

func newIdentity(email string) *identity.NewIdentityAttributes {
	return &identity.NewIdentityAttributes{
		Identifier: email,
		Provider:   "google",
		UserId:     19,
	}
}

func TestMain(m *testing.M) {
	adocker, err := aphdocker.NewArangoDocker()
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	aresource, err := adocker.Run()
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	client, err := adocker.RetryConnection()
	if err != nil {
		log.Fatalf("unable to get client connection %s", err)
	}
	adb = aphdocker.RandString(6)
	dbh, err := client.CreateDatabase(context.Background(), adb, &driver.CreateDatabaseOptions{})
	if err != nil {
		log.Fatalf("could not create arangodb database %s %s\n", adb, err)
	}

	coll, err = dbh.CreateCollection(context.Background(), collection, &driver.CreateCollectionOptions{})
	if err != nil {
		log.Fatalf("could not create arangodb collection %s", collection)
	}
	auser = adocker.GetUser()
	apass = adocker.GetPassword()
	ahost = adocker.GetIP()
	aport = adocker.GetPort()
	code := m.Run()
	if err = adocker.Purge(aresource); err != nil {
		log.Fatalf("unable to remove arangodb container %s\n", err)
	}
	os.Exit(code)
}

func TestCreateIdentity(t *testing.T) {
	ds, err := NewDataSource(auser, apass, adb, ahost, aport)
	if err != nil {
		t.Fatalf("cannot connect to datasource %s", err)
	}
	defer coll.Truncate(context.Background())
	res, err := ds.CreateIdentity(&identity.NewIdentityAttributes{
		Identifier: "hello@gmail.com",
		Provider:   "google",
		UserId:     20,
	})
	if err != nil {
		t.Fatal("could not create new identity")
	}
	if res.GetId() <= 1 {
		t.Fatalf("expected id does not match %d", res.GetId())
	}
	attr := res.GetAttributes()
	if attr.UserId != 20 {
		t.Fatalf("expected user id does not match %d", attr.UserId)
	}
	if attr.Identifier != "hello@gmail.com" {
		t.Fatalf("expected identifier does not match %s", attr.Identifier)
	}
}

func TestHasIdentity(t *testing.T) {
	ds, err := NewDataSource(auser, apass, adb, ahost, aport)
	if err != nil {
		t.Fatalf("cannot connect to datasource %s", err)
	}
	defer coll.Truncate(context.Background())
	res, err := ds.CreateIdentity(&identity.NewIdentityAttributes{
		Identifier: "janto@gmail.com",
		Provider:   "google",
		UserId:     25,
	})
	if err != nil {
		t.Fatal("could not create new identity")
	}
	found, err := ds.HasIdentity(&jsonapi.IdRequest{Id: res.GetId()})
	if err != nil {
		t.Fatalf("error in finding id %d %s", res.GetId(), err)
	}
	if !found {
		t.Fatalf("could not find id %d in storage", res.GetId())
	}
	ifound, err := ds.HasProviderIdentity(
		&identity.IdentityProviderReq{
			Identifier: "janto@gmail.com",
			Provider:   "google",
		})
	if err != nil {
		t.Fatalf("error in finding identity with identifier and provider %s", err)
	}
	if !ifound {
		t.Fatal("could not find identity with identifier and provider in storage")
	}
}

func TestGetIdentity(t *testing.T) {
	ds, err := NewDataSource(auser, apass, adb, ahost, aport)
	if err != nil {
		t.Fatalf("cannot connect to datasource %s", err)
	}
	defer coll.Truncate(context.Background())

	var allRes []storage.Result
	for _, e := range []string{
		"bitnitu@gmail.com",
		"jamba@gmail.com",
	} {
		r, err := ds.CreateIdentity(newIdentity(e))
		if err != nil {
			t.Fatal("could not create new identity with %s", e)
		}
		allRes = append(allRes, r)
	}
	nres, err := ds.GetIdentity(&jsonapi.IdRequest{Id: allRes[0].GetId()})
	if err != nil {
		t.Fatalf("cannot retrieve identity with id %s %s", allRes[0].GetId(), err)
	}
	attr1 := allRes[0].GetAttributes()
	nattr1 := nres.GetAttributes()
	if attr1.Identifier != nattr1.Identifier {
		t.Fatalf("expected identifier %s does not match %s", attr1.Identifier, nattr1.Identifier)
	}
	pres, err := ds.GetProviderIdentity(
		&identity.IdentityProviderReq{
			Identifier: "jamba@gmail.com",
			Provider:   "google",
		})
	if err != nil {
		t.Fatalf("cannot retrieve identity with identity and provider %s", err)
	}
	attr2 := allRes[1].GetAttributes()
	nattr2 := pres.GetAttributes()
	if attr2.Identifier != nattr2.Identifier {
		t.Fatalf("expected identifier %s does not match %s", attr2.Identifier, nattr2.Identifier)
	}
}
