package graph

import (
	"fmt"
	"github.com/arangodb/go-driver"
	"pkv/api/src/domain"
)

/*
// DB is an interface that defines the methods required for a database connection
type DB interface {
	ConnectUserTraining(user domain.User, training domain.Training, ctx context.Context) error

	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetAllPages(ctx context.Context) ([]domain.Page, error)
	GetTrainings(options domain.TrainingQueryOptions, ctx context.Context) ([]domain.TrainingDTO, error)
}*/

type Db struct {
	Database       driver.Database
	Trainings      EntityManager[*domain.Training]
	Locations      EntityManager[*domain.Location]
	Users          EntityManager[*domain.User]
	Logins         EntityManager[*domain.Login]
	Pages          EntityManager[*domain.Page]
	Edges          driver.Collection
	LocationsIndex driver.Index
}

func NewDB(database driver.Database) (*Db, error) {
	trainings, err := NewEntityManager[*domain.Training](database, "trainings", false, func() *domain.Training { return new(domain.Training) })
	if err != nil {
		return nil, err
	}
	locations, err := NewEntityManager[*domain.Location](database, "locations", false, func() *domain.Location { return new(domain.Location) })
	if err != nil {
		return nil, err
	}
	users, err := NewEntityManager[*domain.User](database, "users", false, func() *domain.User { return new(domain.User) })
	if err != nil {
		return nil, err
	}
	logins, err := NewEntityManager[*domain.Login](database, "logins", false, func() *domain.Login { return new(domain.Login) })
	if err != nil {
		return nil, err
	}
	pages, err := NewEntityManager[*domain.Page](database, "pages", false, func() *domain.Page { return new(domain.Page) })
	if err != nil {
		return nil, err
	}
	edges, err := GetOrCreateCollection(database, "edges", true)
	if err != nil {
		return nil, fmt.Errorf("could not get or create edges collection: %w", err)
	}
	locationsIndex, _, err := locations.Collection.EnsureGeoIndex(nil, []string{"lat", "lng"}, nil)
	if err != nil {
		return nil, fmt.Errorf("could not ensure geo index for locations: %w", err)
	}
	return &Db{
		database,
		trainings,
		locations,
		users,
		logins,
		pages,
		edges,
		locationsIndex,
	}, nil
}
