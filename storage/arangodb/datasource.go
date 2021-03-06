package arangodb

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dictyBase/apihelpers/aphcollection"
	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/dictyBase/go-genproto/dictybaseapis/identity"
	"github.com/dictyBase/modware-identity/storage"
)

const (
	collection = "identity"
)

type identityDoc struct {
	Identifier string    `json:"identifier"`
	Provider   string    `json:"provider"`
	UserId     int64     `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type arangoSource struct {
	database   driver.Database
	collection driver.Collection
}

func NewTLSDataSource(user, pass, database, host, port string) (storage.DataSource, error) {
	var ds *arangoSource
	conn, err := http.NewConnection(
		http.ConnectionConfig{
			Endpoints: []string{
				fmt.Sprintf("https://%s:%s", host, port),
			},
			TLSConfig: &tls.Config{InsecureSkipVerify: true},
		})
	if err != nil {
		return ds, fmt.Errorf("could not connect %s", err)
	}
	client, err := driver.NewClient(
		driver.ClientConfig{
			Connection: conn,
			Authentication: driver.BasicAuthentication(
				user,
				pass,
			),
		})
	if err != nil {
		return ds, fmt.Errorf("could not get a client instance %s", err)
	}
	db, err := getDatabase(database, client)
	if err != nil {
		return ds, err
	}
	c, err := getCollection(db)
	if err != nil {
		return ds, err
	}
	return &arangoSource{
		database:   db,
		collection: c,
	}, nil
}

func NewDataSource(user, pass, database, host, port string) (storage.DataSource, error) {
	var ds *arangoSource
	conn, err := http.NewConnection(
		http.ConnectionConfig{
			Endpoints: []string{
				fmt.Sprintf("http://%s:%s", host, port),
			},
		})
	if err != nil {
		return ds, fmt.Errorf("could not connect %s", err)
	}
	client, err := driver.NewClient(
		driver.ClientConfig{
			Connection: conn,
			Authentication: driver.BasicAuthentication(
				user,
				pass,
			),
		})
	if err != nil {
		return ds, fmt.Errorf("could not get a client instance %s", err)
	}
	db, err := getDatabase(database, client)
	if err != nil {
		return ds, err
	}
	c, err := getCollection(db)
	if err != nil {
		return ds, err
	}
	return &arangoSource{
		database:   db,
		collection: c,
	}, nil
}

func (ds *arangoSource) GetIdentity(r *jsonapi.IdRequest) (storage.Result, error) {
	doc := &identityDoc{}
	idStr := strconv.FormatInt(r.Id, 10)
	_, err := ds.collection.ReadDocument(context.Background(), idStr, doc)
	if err != nil {
		if driver.IsNotFound(err) {
			return &arangoResult{
				noResult: true,
			}, nil
		}
		return &arangoResult{
			noResult: true,
		}, err
	}
	return &arangoResult{
		doc:      doc,
		id:       r.Id,
		noResult: false,
	}, nil
}

func (ds *arangoSource) GetIdentityWithAttr(r *jsonapi.IdRequest, fields []string) (storage.Result, error) {
	bindParams := aphcollection.MapIdx(
		fields,
		func(s string, i int) string { return fmt.Sprintf("@attr%d", i) },
	)
	query := fmt.Sprintf(`
				FOR d in @@collection
				FILTER d._key == @id
				RETURN KEEP(d,["_id","_key","_rev", %s])`,
		strings.Join(bindParams, ","),
	)
	bindVars := map[string]interface{}{
		"@collection": collection,
		"id":          strconv.FormatInt(r.Id, 10),
	}
	for i, v := range bindParams {
		bindVars[strings.Replace(v, "@", "", 1)] = fields[i]
	}
	cursor, err := ds.database.Query(nil, query, bindVars)
	if err != nil {
		return &arangoResult{
			noResult: true,
		}, err
	}
	defer cursor.Close()
	doc := &identityDoc{}
	meta, err := cursor.ReadDocument(context.Background(), doc)
	if err != nil {
		if driver.IsNotFound(err) {
			return &arangoResult{
				noResult: true,
			}, nil
		}
		return &arangoResult{
			noResult: true,
		}, err
	}
	id, err := strconv.ParseInt(meta.Key, 10, 64)
	if err != nil {
		return &arangoResult{
			noResult: true,
		}, err
	}
	return &arangoResult{
		doc:      doc,
		id:       id,
		noResult: false,
	}, nil
}

func (ds *arangoSource) GetProviderIdentity(r *identity.IdentityProviderReq) (storage.Result, error) {
	query := `FOR d in @@collection
				FILTER d.identifier == @identifier
				AND d.provider == @provider
				RETURN d`
	bindVars := map[string]interface{}{
		"@collection": collection,
		"identifier":  r.Identifier,
		"provider":    r.Provider,
	}
	cursor, err := ds.database.Query(nil, query, bindVars)
	if err != nil {
		return &arangoResult{
			noResult: true,
		}, err
	}
	defer cursor.Close()
	doc := &identityDoc{}
	meta, err := cursor.ReadDocument(nil, doc)
	if err != nil {
		if driver.IsNotFound(err) {
			return &arangoResult{
				noResult: true,
			}, nil
		}
		return &arangoResult{
			noResult: true,
		}, err
	}
	id, err := strconv.ParseInt(meta.Key, 10, 64)
	if err != nil {
		return &arangoResult{
			noResult: true,
		}, err
	}
	return &arangoResult{
		doc:      doc,
		id:       id,
		noResult: false,
	}, nil
}

func (ds *arangoSource) CreateIdentity(attr *identity.NewIdentityAttributes) (storage.Result, error) {
	insert := `INSERT {
					identifier: @identifier,
					provider: @provider,
					user_id: @user_id,
					created_at: DATE_ISO8601(DATE_NOW()),
					updated_at: DATE_ISO8601(DATE_NOW())
			   } IN @@collection RETURN NEW`
	bindVars := map[string]interface{}{
		"@collection": collection,
		"identifier":  attr.Identifier,
		"provider":    attr.Provider,
		"user_id":     attr.UserId,
	}
	cursor, err := ds.database.Query(context.Background(), insert, bindVars)
	if err != nil {
		return &arangoResult{
			noResult: true,
		}, err
	}
	defer cursor.Close()
	doc := &identityDoc{}
	meta, err := cursor.ReadDocument(context.Background(), doc)
	if err != nil {
		return &arangoResult{
			noResult: true,
		}, err
	}
	id, err := strconv.ParseInt(meta.Key, 10, 64)
	if err != nil {
		return &arangoResult{
			noResult: true,
		}, err
	}
	return &arangoResult{
		doc:      doc,
		id:       id,
		noResult: false,
	}, nil
}

func (ds *arangoSource) DeleteIdentity(r *jsonapi.IdRequest) (bool, error) {
	_, err := ds.collection.RemoveDocument(nil, strconv.FormatInt(r.Id, 10))
	if err != nil {
		if driver.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ds *arangoSource) HasProviderIdentity(r *identity.IdentityProviderReq) (bool, error) {
	ctx := driver.WithQueryCount(context.Background())
	query := `FOR d in @@collection
				FILTER d.identifier == @identifier
				AND d.provider == @provider
				RETURN d`
	bindVars := map[string]interface{}{
		"@collection": collection,
		"identifier":  r.Identifier,
		"provider":    r.Provider,
	}
	cursor, err := ds.database.Query(ctx, query, bindVars)
	if err != nil {
		return false, err
	}
	defer cursor.Close()
	if cursor.Count() > 0 {
		return true, nil
	}
	return false, nil
}

func (ds *arangoSource) HasIdentity(r *jsonapi.IdRequest) (bool, error) {
	return ds.collection.DocumentExists(context.Background(), strconv.FormatInt(r.Id, 10))
}

// ClearIdentities clears all identities from the repository datasource
func (ds *arangoSource) ClearIdentities() error {
	if err := ds.collection.Truncate(context.Background()); err != nil {
		return err
	}
	return nil
}
