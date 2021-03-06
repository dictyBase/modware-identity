package arangodb

import (
	"log"
	"os"
	"strconv"
	"testing"

	driver "github.com/arangodb/go-driver"
	"github.com/dictyBase/arangomanager/testarango"
	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	"github.com/dictyBase/go-genproto/dictybaseapis/identity"
	"github.com/dictyBase/modware-identity/storage"
)

var gta *testarango.TestArango
var ahost, aport, auser, apass, adb string

func newIdentity(email string) *identity.NewIdentityAttributes {
	return &identity.NewIdentityAttributes{
		Identifier: email,
		Provider:   "google",
		UserId:     19,
	}
}

func TestMain(m *testing.M) {
	ta, err := testarango.NewTestArangoFromEnv(true)
	if err != nil {
		log.Fatalf("unable to construct new TestArango instance %s", err)
	}
	gta = ta
	dbh, err := ta.DB(ta.Database)
	if err != nil {
		log.Fatalf("unable to get database %s", err)
	}
	auser = gta.User
	apass = gta.Pass
	ahost = gta.Host
	aport = strconv.Itoa(gta.Port)
	adb = gta.Database
	coll := "test-collection"
	_, err = dbh.CreateCollection(coll, &driver.CreateCollectionOptions{})
	if err != nil {
		if err := dbh.Drop(); err != nil {
			log.Printf("error in dropping database %s", err)
		}
		log.Fatalf("unable to create collection %s %s", coll, err)
	}
	code := m.Run()
	if err := dbh.Drop(); err != nil {
		log.Printf("error in dropping database %s", err)
	}
	os.Exit(code)
}

func TestCreateIdentity(t *testing.T) {
	ds, err := NewDataSource(auser, apass, adb, ahost, aport)
	if err != nil {
		t.Fatalf("cannot connect to datasource %s", err)
	}
	defer func() {
		err := ds.ClearIdentities()
		if err != nil {
			t.Fatalf("error in clearing identities %s", err)
		}
	}()
	res, err := ds.CreateIdentity(&identity.NewIdentityAttributes{
		Identifier: "hello@gmail.com",
		Provider:   "google",
		UserId:     20,
	})
	if err != nil {
		t.Fatalf("could not create new identity %s", err)
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
	defer func() {
		err := ds.ClearIdentities()
		if err != nil {
			t.Fatalf("error in clearing identities %s", err)
		}
	}()
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
	defer func() {
		err := ds.ClearIdentities()
		if err != nil {
			t.Fatalf("error in clearing identities %s", err)
		}
	}()
	var allRes []storage.Result
	for _, e := range []string{
		"bitnitu@gmail.com",
		"jamba@gmail.com",
	} {
		r, err := ds.CreateIdentity(newIdentity(e))
		if err != nil {
			t.Fatalf("could not create new identity with %s", e)
		}
		allRes = append(allRes, r)
	}
	nres, err := ds.GetIdentity(&jsonapi.IdRequest{Id: allRes[0].GetId()})
	if err != nil {
		t.Fatalf("cannot retrieve identity with id %d %s", allRes[0].GetId(), err)
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

func TestGetIdentityWithAttr(t *testing.T) {
	ds, err := NewDataSource(auser, apass, adb, ahost, aport)
	if err != nil {
		t.Fatalf("cannot connect to datasource %s", err)
	}
	defer func() {
		err := ds.ClearIdentities()
		if err != nil {
			t.Fatalf("error in clearing identities %s", err)
		}
	}()
	var allRes []storage.Result
	for _, e := range []string{
		"bitnitu@gmail.com",
		"jamba@gmail.com",
	} {
		r, err := ds.CreateIdentity(newIdentity(e))
		if err != nil {
			t.Fatalf("could not create new identity with %s", e)
		}
		allRes = append(allRes, r)
	}
	nres, err := ds.GetIdentityWithAttr(
		&jsonapi.IdRequest{Id: allRes[0].GetId()},
		[]string{"identifier"},
	)
	if err != nil {
		t.Fatalf("cannot retrieve identifier with id %d %s", allRes[0].GetId(), err)
	}
	if nres.NotFound() == true {
		t.Fatal("error in retrieving result")
	}
	if allRes[0].GetId() != nres.GetId() {
		t.Fatalf("expected id %d does not match %d", allRes[0].GetId(), nres.GetId())
	}
	attr1 := allRes[0].GetAttributes()
	nattr1 := nres.GetAttributes()
	if attr1.Identifier != nattr1.Identifier {
		t.Fatalf("expected identifier %s does not match %s", attr1.Identifier, nattr1.Identifier)
	}
	if len(nattr1.Provider) != 0 {
		t.Fatalf("expected empty provider does not match %s\n", nattr1.Provider)
	}
}

func TestDeleteIdentity(t *testing.T) {
	ds, err := NewDataSource(auser, apass, adb, ahost, aport)
	if err != nil {
		t.Fatalf("cannot connect to datasource %s", err)
	}
	defer func() {
		err := ds.ClearIdentities()
		if err != nil {
			t.Fatalf("error in clearing identities %s", err)
		}
	}()
	res, err := ds.CreateIdentity(&identity.NewIdentityAttributes{
		Identifier: "hello@gmail.com",
		Provider:   "google",
		UserId:     20,
	})
	if err != nil {
		t.Fatalf("could not create new identity %s", err)
	}
	done, err := ds.DeleteIdentity(&jsonapi.IdRequest{Id: res.GetId()})
	if err != nil {
		t.Fatalf("cannot delete identity with id %d %s", res.GetId(), err)
	}
	if done == false {
		t.Fatal("cannot delete identity")
	}
}
