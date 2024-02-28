package graph

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/connection"
	"log"
	"math/rand"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/security"
	"strings"
	"time"
)

func Connect(config *dpv.Config, useRoot bool) (arangodb.Client, error) {
	var auth connection.Authentication
	if useRoot {
		auth = connection.NewBasicAuth("root", config.DB.Root)
	} else {
		connection.NewBasicAuth(config.DB.User, config.DB.Pass)
	}
	conn := connection.NewHttpConnection(connection.HttpConfiguration{
		Authentication: auth,
		Endpoint:       connection.NewRoundRobinEndpoints([]string{fmt.Sprintf("http://%s:%d", config.DB.Host, config.DB.Port)}),
	})

	c := arangodb.NewClient(conn)

	_, err := c.Version(context.Background())

	return c, err
}

func DropTestDatabases(c arangodb.Client) error {
	dbs, err := c.Databases(context.Background())
	if err != nil {
		return fmt.Errorf("could not list databases: %w", err)
	}
	for _, db := range dbs {
		if strings.HasPrefix(db.Name(), "test-") {
			err = db.Remove(context.Background())
			if err != nil {
				return fmt.Errorf("could not remove database: %w", err)
			}
		}
	}
	return nil
}

func GetOrCreateDatabase(c arangodb.Client, dbname string, config *dpv.Config) (arangodb.Database, error) {
	var db arangodb.Database
	if ok, err := c.DatabaseExists(context.Background(), dbname); !ok || err != nil {
		// fix arangodb implementation bug
		if err.Error() == "database not found" {
			ok = false
			err = nil
		}
		if err != nil {
			return nil, fmt.Errorf("failed to look for database: %w", err)
		}
		trueBool := true
		if db, err = c.CreateDatabase(context.Background(), dbname, &arangodb.CreateDatabaseOptions{Users: []arangodb.CreateDatabaseUserOptions{
			{config.DB.User, config.DB.Pass, &trueBool, nil},
		}}); err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	} else {
		if db, err = c.Database(context.Background(), dbname); err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}
	}
	return db, nil
}

func GetOrCreateCollection(db arangodb.Database, name string, edges bool) (arangodb.Collection, error) {
	if ok, err := db.CollectionExists(context.Background(), name); !ok || err != nil {
		// fix arangodb implementation bug
		if err.Error() == "collection or view not found" {
			ok = false
			err = nil
		}
		if err != nil {
			return nil, fmt.Errorf("could not check if collection exists: %w", err)
		}
		if edges {
			return db.CreateCollection(context.Background(), name, &arangodb.CreateCollectionProperties{Type: arangodb.CollectionTypeEdge})
		} else {
			return db.CreateCollection(context.Background(), name, &arangodb.CreateCollectionProperties{
				ComputedValues: []arangodb.ComputedValue{
					{
						Name:       "created",
						Expression: "RETURN DATE_ISO8601(DATE_NOW())",
						Overwrite:  true,
						ComputeOn:  []arangodb.ComputeOn{arangodb.ComputeOnInsert},
					},
					{
						Name:       "modified",
						Expression: "RETURN DATE_ISO8601(DATE_NOW())",
						Overwrite:  true,
						ComputeOn:  []arangodb.ComputeOn{arangodb.ComputeOnInsert, arangodb.ComputeOnReplace, arangodb.ComputeOnUpdate},
					},
				},
				Type: arangodb.CollectionTypeDocument,
			})
		}
	} else {
		return db.Collection(context.Background(), name)
	}
}

var fields map[string]arangodb.ArangoSearchElementProperties

func FieldsForAllLanguages(config *dpv.Config) map[string]arangodb.ArangoSearchElementProperties {
	if fields == nil {
		fields = make(map[string]arangodb.ArangoSearchElementProperties)
		for _, language := range config.Settings.Languages {
			fields[language.Key] = arangodb.ArangoSearchElementProperties{
				Fields: map[string]arangodb.ArangoSearchElementProperties{
					"text": {
						Analyzers: []string{"text_" + language.Key},
					},
				},
			}
		}
	}
	return fields
}

func CreateViewIfNotExists(db arangodb.Database, config *dpv.Config, name string) error {
	ok, err := db.ViewExists(context.Background(), name+"-descriptions")
	if err != nil {
		return fmt.Errorf("could not check if view for collection %v exists: %w", name, err)
	}
	if !ok {
		_, err := db.CreateArangoSearchView(context.Background(), name+"-descriptions", &arangodb.ArangoSearchViewProperties{
			Links: map[string]arangodb.ArangoSearchElementProperties{
				"users": {
					Fields: map[string]arangodb.ArangoSearchElementProperties{
						"descriptions": {
							Fields: FieldsForAllLanguages(config),
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("could not create view for collection %v: %w", name, err)
		}
	}
	return nil
}

func NewEntityManager[T Entity](db arangodb.Database, name string, edges bool, constructor func() T) (EntityManager[T], error) {
	collection, err := GetOrCreateCollection(db, name, edges)
	if err != nil {
		return EntityManager[T]{}, fmt.Errorf("could not get or create %s collection: %w", name, err)
	}
	return EntityManager[T]{collection, constructor}, nil
}

func Init(configPath string, test bool) (*Db, *dpv.Config, error) {
	var err error
	config, err := dpv.NewConfig(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("could not initialise config instance: %w", err)
	}
	c, err := Connect(config, true)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to database server: %w", err)
	}
	dbname := "dpv"
	if test {
		dbname = "test-" + dbname + "-" + security.HashToken(fmt.Sprintf("%s-%x", time.Now().String(), rand.Int()))[0:8]
		log.Printf("Using database %s\n", dbname)
	}
	database, err := GetOrCreateDatabase(c, dbname, config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not use database: %w", err)
	}
	db, err := NewDB(database, config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not initialise database: %w", err)
	}
	if !test {
		collection, err := db.Database.Collection(context.Background(), "users")
		if err != nil {
			return nil, nil, fmt.Errorf("could not get users collection: %w", err)
		}
		count, err := collection.Count(context.Background())
		if err != nil {
			return nil, nil, fmt.Errorf("could not count users: %w", err)
		}
		if count == 0 {
			log.Println("Creating sample data")
			SampleData(db)
		}
	}
	return db, config, err
}
