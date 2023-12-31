package graph

import (
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"log"
	"math/rand"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/security"
	"strings"
	"time"
)

func Connect(config *dpv.Config, useRoot bool) (driver.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{fmt.Sprintf("http://%s:%d", config.DB.Host, config.DB.Port)},
	})
	if err != nil {
		return nil, fmt.Errorf("could not establish a http connection to ArangoDB: %w", err)
	}

	var auth driver.Authentication
	if useRoot {
		auth = driver.BasicAuthentication("root", config.DB.Root)
	} else {
		driver.BasicAuthentication(config.DB.User, config.DB.Pass)
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: auth,
	})
	if err != nil {
		return nil, fmt.Errorf("could not connect to the ArangoDB database: %w", err)
	}
	return c, nil
}

func DropTestDatabases(c driver.ClientDatabases) error {
	dbs, err := c.Databases(nil)
	if err != nil {
		return fmt.Errorf("could not list databases: %w", err)
	}
	for _, db := range dbs {
		if strings.HasPrefix(db.Name(), "test-") {
			err = db.Remove(nil)
			if err != nil {
				return fmt.Errorf("could not remove database: %w", err)
			}
		}
	}
	return nil
}

func GetOrCreateDatabase(c driver.ClientDatabases, dbname string, config *dpv.Config) (driver.Database, error) {
	var db driver.Database
	if ok, err := c.DatabaseExists(nil, dbname); !ok || err != nil {
		if err != nil {
			return nil, fmt.Errorf("failed to look for database: %w", err)
		}
		trueBool := true
		if db, err = c.CreateDatabase(nil, dbname, &driver.CreateDatabaseOptions{Users: []driver.CreateDatabaseUserOptions{
			{config.DB.User, config.DB.Pass, &trueBool, nil},
		}}); err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	} else {
		if db, err = c.Database(nil, dbname); err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}
	}
	return db, nil
}

func GetOrCreateCollection(db driver.Database, name string, edges bool) (driver.Collection, error) {
	if ok, err := db.CollectionExists(nil, name); !ok || err != nil {
		if err != nil {
			return nil, fmt.Errorf("could not check if collection exists: %w", err)
		}
		if edges {
			return db.CreateCollection(nil, name, &driver.CreateCollectionOptions{Type: driver.CollectionTypeEdge})
		} else {
			return db.CreateCollection(nil, name, &driver.CreateCollectionOptions{
				ComputedValues: []driver.ComputedValue{
					{
						Name:       "created",
						Expression: "RETURN DATE_ISO8601(DATE_NOW())",
						Overwrite:  true,
						ComputeOn:  []driver.ComputeOn{driver.ComputeOnInsert},
					},
					{
						Name:       "modified",
						Expression: "RETURN DATE_ISO8601(DATE_NOW())",
						Overwrite:  true,
						ComputeOn:  []driver.ComputeOn{driver.ComputeOnInsert, driver.ComputeOnReplace, driver.ComputeOnUpdate},
					},
				},
				Type: driver.CollectionTypeDocument,
			})
		}
	} else {
		return db.Collection(nil, name)
	}
}

var fields map[string]driver.ArangoSearchElementProperties

func FieldsForAllLanguages(config *dpv.Config) map[string]driver.ArangoSearchElementProperties {
	if fields == nil {
		fields = make(map[string]driver.ArangoSearchElementProperties)
		for _, language := range config.Settings.Languages {
			fields[language.Key] = driver.ArangoSearchElementProperties{
				Fields: map[string]driver.ArangoSearchElementProperties{
					"text": {
						Analyzers: []string{"text_" + language.Key},
					},
				},
			}
		}
	}
	return fields
}

func CreateViewIfNotExists(db driver.Database, config *dpv.Config, name string) error {
	ok, err := db.ViewExists(nil, name+"-descriptions")
	if err != nil {
		return fmt.Errorf("could not check if view for collection %v exists: %w", name, err)
	}
	if !ok {
		_, err := db.CreateArangoSearchView(nil, name+"-descriptions", &driver.ArangoSearchViewProperties{
			Links: map[string]driver.ArangoSearchElementProperties{
				"users": {
					Fields: map[string]driver.ArangoSearchElementProperties{
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

func NewEntityManager[T Entity](db driver.Database, name string, edges bool, constructor func() T) (EntityManager[T], error) {
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
		collection, err := db.Database.Collection(nil, "users")
		if err != nil {
			return nil, nil, fmt.Errorf("could not get users collection: %w", err)
		}
		count, err := collection.Count(nil)
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
