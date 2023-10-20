package graph

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"pkv/api/src/domain"
	"pkv/api/src/repository/dpv"
)

func (db *Db) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return db.GetUsers(buildAllUsersQuery, ctx)
}

func (db *Db) GetFilteredUsers(options domain.UserQueryOptions, ctx context.Context) ([]domain.User, error) {
	return db.GetUsers(func() (string, map[string]interface{}) { return buildUserQuery(options) }, ctx)
}

func (db *Db) GetAdministeredUsers(key string, ctx context.Context) ([]domain.User, error) {
	return db.GetUsers(func() (string, map[string]interface{}) { return buildAdministeredUsersQuery(key) }, ctx)
}

func (db *Db) GetAdministrators(key string, ctx context.Context) ([]domain.User, error) {
	return db.GetUsers(func() (string, map[string]interface{}) { return buildAdministratorsQuery(key) }, ctx)
}

func (db *Db) GetUsers(queryBuilder QueryBuilder, ctx context.Context) ([]domain.User, error) {

	query, bindVars := queryBuilder()
	cursor, err := db.Database.Query(ctx, query, bindVars)
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

func buildAllUsersQuery() (string, map[string]interface{}) {
	query := "FOR doc IN users RETURN doc"
	var bindVars map[string]interface{}
	return query, bindVars
}

func buildUserQuery(options domain.UserQueryOptions) (string, map[string]interface{}) {
	includeSet := options.Include
	bindVars := make(map[string]interface{})
	var query string
	if options.Text != "" {
		lang := options.Language
		valid := false
		for _, language := range dpv.ConfigInstance.Settings.Languages {
			if language.Key == lang {
				valid = true
				break
			}
		}
		if !valid {
			lang = "en"
		}
		query += "FOR user IN `users-descriptions`\n"
		query += fmt.Sprintf(`  SEARCH ANALYZER(TOKENS(@text, "text_%s") ALL == user.descriptions.%s.text, "text_%s")`, lang, lang, lang)
		bindVars["text"] = options.Text
	} else {
		query += "FOR user IN users\n"
	}
	if options.Key != "" {
		query += "  FILTER user._key == @key\n"
		bindVars["key"] = options.Key
	}
	if options.Name != "" {
		query += "  FILTER user.name == @name\n"
		bindVars["name"] = options.Name
	}
	if options.Type != "" {
		query += "  FILTER user.type == @type\n"
		bindVars["type"] = options.Type
	}

	unsetUser := buildUnsetParts(includeSet, "")
	userStr := buildUnsetString("user", unsetUser)
	query += "  RETURN " + userStr
	return query, bindVars
}

func buildAdministeredUsersQuery(key string) (string, map[string]interface{}) {
	// user{_key: @key}-[administers*0..]->user
	query := "FOR doc IN users\n"
	query += "  FILTER doc._key == @key\n"
	query += "  FOR v, e IN 0..99 OUTBOUND doc edges\n"
	query += "    PRUNE e != null && e.label != \"administers\""
	query += "    FILTER e == null || e.label == \"administers\" RETURN DISTINCT v"
	return query, map[string]interface{}{"key": key}
}

func buildAdministratorsQuery(key string) (string, map[string]interface{}) {
	// user{_key: @key}<-[administers*0..]-user
	query := "FOR doc IN users\n"
	query += "  FILTER doc._key == @key\n"
	query += "  FOR v, e IN 0..99 INBOUND doc edges\n"
	query += "    PRUNE e != null && e.label != \"administers\""
	query += "    FILTER e == null || e.label == \"administers\" RETURN DISTINCT v"
	return query, map[string]interface{}{"key": key}
}
