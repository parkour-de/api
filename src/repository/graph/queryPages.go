package graph

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"pkv/api/src/domain"
)

func (db *Db) GetAllPages(ctx context.Context) ([]domain.Page, error) {

	query := "FOR doc IN pages RETURN doc"
	cursor, err := db.Database.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("query string invalid: %w", err)
	}
	defer cursor.Close()

	var result []domain.Page
	for {
		var doc domain.Page
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
