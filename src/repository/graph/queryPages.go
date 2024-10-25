package graph

import (
	"context"
	"github.com/arangodb/go-driver/v2/arangodb/shared"
	"pkv/api/src/domain"
	"pkv/api/src/repository/t"
)

func (db *Db) GetAllPages(ctx context.Context) ([]domain.Page, error) {

	query := "FOR doc IN pages RETURN doc"
	cursor, err := db.Database.Query(ctx, query, nil)
	if err != nil {
		return nil, t.Errorf("query string invalid: %w", err)
	}
	defer cursor.Close()

	var result []domain.Page
	for {
		var doc domain.Page
		_, err := cursor.ReadDocument(ctx, &doc)
		if shared.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, t.Errorf("obtaining documents failed: %w", err)
		}
		result = append(result, doc)
	}

	return result, nil
}
