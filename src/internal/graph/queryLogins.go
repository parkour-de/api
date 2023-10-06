package graph

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"pkv/api/src/domain"
)

func (db *Db) GetLoginsForUser(key string, ctx context.Context) ([]domain.Login, error) {
	query := "FOR user IN users\n"
	query += "  FILTER user._key == @key\n"
	query += "  LET logins = (FOR login, e IN 1..1 INBOUND user edges FILTER e.label == \"authenticates\"\n"
	query += "    RETURN UNSET(login, \"subject\"))\n"
	query += "  RETURN logins"
	cursor, err := db.Database.Query(ctx, query, map[string]interface{}{"key": key})
	if err != nil {
		return nil, fmt.Errorf("query string invalid: %w", err)
	}
	defer cursor.Close()

	var result []domain.Login
	for {
		var doc []domain.Login
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("obtaining documents failed: %w", err)
		}
		result = append(result, doc...)
	}
	return result, nil
}
