package graph

import (
	"fmt"
	"github.com/arangodb/go-driver"
	"pkv/api/src/domain"
)

type EntityManager[T Entity] struct {
	Collection  driver.Collection
	Constructor func() T
}

type Entity interface {
	GetID() string
	SetID(id string)
}

func (im *EntityManager[T]) Create(item T) error {
	meta, err := im.Collection.CreateDocument(nil, item)
	if err != nil {
		return fmt.Errorf("could not create item: %w", err)
	}
	item.SetID(meta.ID.String())
	return nil
}

func (im *EntityManager[T]) Read(id string) (T, error) {
	item := im.Constructor()
	meta, err := im.Collection.ReadDocument(nil, id, item)
	if err != nil {
		return item, fmt.Errorf("could not read item with id %v: %w", id, err)
	}
	item.SetID(meta.ID.String())
	return item, nil
}

func (im *EntityManager[T]) Update(item T) error {
	_, err := im.Collection.UpdateDocument(nil, item.GetID(), item)
	if err != nil {
		return fmt.Errorf("could not update item with id %v: %w", item.GetID(), err)
	}
	return nil
}

func (im *EntityManager[T]) Delete(item T) error {
	_, err := im.Collection.RemoveDocument(nil, item.GetID())
	if err != nil {
		return fmt.Errorf("could not delete item with id %v: %w", item.GetID(), err)
	}
	return nil
}

func (db *Db) ConnectTrainingLocation(training *domain.Training, location *domain.Location) error {
	if _, err := db.Edges.CreateDocument(nil, domain.Edge{
		From:  training.ID,
		To:    location.ID,
		Label: "happens_at",
	}); err != nil {
		return fmt.Errorf("could not connect training %s to location %s: %w", training.ID, location.ID, err)
	}
	return nil
}

func (db *Db) ConnectUserTraining(user domain.User, training domain.Training) error {
	if _, err := db.Edges.CreateDocument(nil, domain.Edge{
		From:  user.ID,
		To:    training.ID,
		Label: "organises",
	}); err != nil {
		return fmt.Errorf("could not connect user %s to training %s: %w", user.ID, training.ID, err)
	}
	return nil
}
