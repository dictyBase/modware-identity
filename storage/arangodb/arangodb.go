package arangodb

import (
	"context"
	"fmt"

	driver "github.com/arangodb/go-driver"
)

func getDatabase(database string, client driver.Client) (driver.Database, error) {
	var db driver.Database
	ok, err := client.DatabaseExists(context.Background(), database)
	if err != nil { 
		return db, fmt.Errorf("unable to check for database %s", err)
	}
	if !ok {
		return db, fmt.Errorf("database %s has to be created", database)
	}
	return client.Database(context.Background(), database)
}

func getCollection(db driver.Database) (driver.Collection, error) {
	var c driver.Collection
	ok, err := db.CollectionExists(context.Background(), collection)
	if err != nil {
		return c, fmt.Errorf("unable to check for collection %s", collection)
	}
	if !ok {
		// create collection if it doesn't exist
		_, err := db.CreateCollection(context.Background(), collection, nil)
		if err != nil {
			return c, fmt.Errorf("failed to create %s collection: %v", collection, err)
		}
	}

	return db.Collection(context.Background(), collection)
}
