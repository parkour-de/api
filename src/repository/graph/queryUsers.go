package graph

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"pkv/api/src/domain"
)

func (db *Db) GetAllUsers(ctx context.Context) ([]domain.User, error) {

	query := "FOR doc IN users RETURN doc"
	cursor, err := db.Database.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("query string invalid: %w", err)
	}
	defer cursor.Close()

	var result []domain.User
	for {
		var doc domain.User
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("obtaining documents failed: %w", err)
		}
		result = append(result, doc)
	}

	return result, nil
}
func (db *Db) GetAssumableUsers(key string, ctx context.Context) ([]domain.User, error) {

	// user{_key: @key} -[administers*0..]->user
	query := "FOR doc IN users\n"
	query += "  FILTER doc._key == @key\n"
	query += "  FOR v, e IN 0..99 OUTBOUND doc edges\n"
	query += "    PRUNE e != null && e.label != \"administers\""
	query += "    FILTER e == null || e.label == \"administers\" RETURN DISTINCT v"

	cursor, err := db.Database.Query(ctx, query, map[string]interface{}{"key": key})
	if err != nil {
		return nil, fmt.Errorf("query string invalid: %w", err)
	}
	defer cursor.Close()

	var result []domain.User
	for {
		var doc domain.User
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("obtaining documents failed: %w", err)
		}
		result = append(result, doc)
	}

	return result, nil
}
