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

func (db *Db) GetLoginsByProvider(provider string, sub string, ctx context.Context) ([]domain.Login, error) {
	query := "FOR login IN logins\n"
	query += "  FILTER login.provider == @provider\n"
	query += "  FILTER login.subject == @sub\n"
	query += "  RETURN login"
	cursor, err := db.Database.Query(ctx, query, map[string]interface{}{"provider": provider, "sub": sub})
	if err != nil {
		return nil, fmt.Errorf("query string invalid: %w", err)
	}
	defer cursor.Close()

	var result []domain.Login
	for {
		var doc domain.Login
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

func (db *Db) GetLoginUsers(key string, ctx context.Context) ([]domain.User, error) {
	query := "FOR login IN logins\n"
	query += "  FILTER login._key == @key\n"
	query += "  FOR user, e IN 1..1 OUTBOUND login edges\n"
	query += "    FILTER e.label == \"authenticates\"\n"
	query += "    RETURN user"
	cursor, err := db.Database.Query(ctx, query, map[string]interface{}{"key": key})
	if err != nil {
		return nil, fmt.Errorf("query string invalid: %w", err)
	}
	defer cursor.Close()

	var result []domain.User
	for {
		var doc domain.User
		_, err = cursor.ReadDocument(ctx, &result)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("obtaining documents failed: %w", err)
		}
		result = append(result, doc)
	}
	return result, nil
}
