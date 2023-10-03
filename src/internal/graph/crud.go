package graph

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"pkv/api/src/domain"
)

type EntityManager[T Entity] struct {
	Collection  driver.Collection
	Constructor func() T
}

type Entity interface {
	GetKey() string
	SetKey(key string)
}

func (im *EntityManager[T]) Create(item T, ctx context.Context) error {
	meta, err := im.Collection.CreateDocument(ctx, item)
	if err != nil {
		return fmt.Errorf("could not create item: %w", err)
	}
	item.SetKey(meta.Key)
	return nil
}

func (im *EntityManager[T]) Has(key string, ctx context.Context) (bool, error) {
	exists, err := im.Collection.DocumentExists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("could not check for item with key %v: %w", key, err)
	}
	return exists, nil
}

func (im *EntityManager[T]) Read(key string, ctx context.Context) (T, error) {
	item := im.Constructor()
	meta, err := im.Collection.ReadDocument(ctx, key, item)
	if err != nil {
		return item, fmt.Errorf("could not read item with key %v: %w", key, err)
	}
	item.SetKey(meta.Key)
	return item, nil
}

func (im *EntityManager[T]) Update(item T, ctx context.Context) error {
	_, err := im.Collection.UpdateDocument(ctx, item.GetKey(), item)
	if err != nil {
		return fmt.Errorf("could not update item with key %v: %w", item.GetKey(), err)
	}
	return nil
}

func (im *EntityManager[T]) Delete(item T, ctx context.Context) error {
	_, err := im.Collection.RemoveDocument(ctx, item.GetKey())
	if err != nil {
		return fmt.Errorf("could not delete item with key %v: %w", item.GetKey(), err)
	}
	return nil
}

func (db *Db) ConnectTrainingLocation(training *domain.Training, location *domain.Location, ctx context.Context) error {
	if _, err := db.Edges.CreateDocument(ctx, domain.Edge{
		From:  "trainings/" + training.Key,
		To:    "locations/" + location.Key,
		Label: "happens_at",
	}); err != nil {
		return fmt.Errorf("could not connect training %s to location %s: %w", training.Key, location.Key, err)
	}
	return nil
}

func (db *Db) ConnectUserTraining(user domain.User, training domain.Training, ctx context.Context) error {
	if _, err := db.Edges.CreateDocument(ctx, domain.Edge{
		From:  "users/" + user.Key,
		To:    "trainings/" + training.Key,
		Label: "organises",
	}); err != nil {
		return fmt.Errorf("could not connect user %s to training %s: %w", user.Key, training.Key, err)
	}
	return nil
}
