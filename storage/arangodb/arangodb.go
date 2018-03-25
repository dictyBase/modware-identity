package arangodb

import (
	"fmt"

	driver "github.com/arangodb/go-driver"
)

func getDatabase(database string, client driver.Client) (driver.Database, error) {
	var db driver.Database
	ok, err := client.DatabaseExists(nil, database)
	if err != nil {
		return db, fmt.Errorf("unable to check for database %s", err)
	}
	if !ok {
		return db, fmt.Errorf("database %s has to be created", database)
	}
	return client.Database(nil, database)
}

func getCollection(db driver.Database) (driver.Collection, error) {
	var c driver.Collection
	ok, err := db.CollectionExists(nil, collection)
	if err != nil {
		return c, fmt.Errorf("unable to check for collection %s", collection)
	}
	if !ok {
		return c, fmt.Errorf("collection %s has to be created", collection)
	}
	return db.CreateCollection(nil, collection, nil)
}
