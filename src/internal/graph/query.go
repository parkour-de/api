package graph

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"pkv/api/src/domain"
	"strings"
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

func (db *Db) GetTrainings(options domain.TrainingQueryOptions, ctx context.Context) ([]domain.TrainingDTO, error) {
	query, bindVars := buildTrainingQuery(options)
	cursor, err := db.Database.Query(ctx, query, bindVars)
	if err != nil {
		return nil, fmt.Errorf("query string invalid: %w", err)
	}
	defer cursor.Close()

	var result []domain.TrainingDTO
	for {
		var doc domain.TrainingDTO
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

func buildTrainingQuery(options domain.TrainingQueryOptions) (string, map[string]interface{}) {
	includeSet := options.Include

	unsetOrganiser := buildUnsetParts(includeSet, "organiser_")
	unsetLocation := buildUnsetParts(includeSet, "location_")
	unsetTraining := buildUnsetParts(includeSet, "")
	unsetTraining = appendUnsetPart(unsetTraining, includeSet, "cycles", "cycles")

	query := "FOR training IN trainings\n"
	bindVars := make(map[string]interface{})
	if options.Weekday != 0 {
		query += "  FILTER training.cycles[? ANY FILTER CURRENT.weekday == @weekday]\n"
		bindVars["weekday"] = options.Weekday
	}
	locationStr := buildUnsetString("location", unsetLocation)
	query += "  LET location = FIRST( FOR location, e IN OUTBOUND training edges FILTER e.label == \"happens_at\" RETURN " + locationStr + " )\n"
	if options.City != "" {
		query += "  FILTER location.city == @city\n"
		bindVars["city"] = options.City
	}
	if options.LocationKey != "" {
		query += "  FILTER @locationKey == location._key\n"
		bindVars["locationKey"] = options.LocationKey
	}
	organiserStr := buildUnsetString("organiser", unsetOrganiser)
	query += "  LET organisers = (FOR organiser, e IN 1..1 INBOUND training edges FILTER e.label == \"organises\" RETURN " + organiserStr + ")\n"
	if options.OrganiserKey != "" {
		query += "  FILTER @organiserKey IN organisers[*]._key\n"
		bindVars["organiserKey"] = options.OrganiserKey
	}

	trainingStr := buildUnsetString("training", unsetTraining)
	query += "  RETURN MERGE(" + trainingStr + ", {"
	var sections []string
	if _, ok := includeSet["location"]; ok {
		sections = append(sections, "location: location")
	}
	sections = append(sections, "locationKey: location._key")
	if _, ok := includeSet["organisers"]; ok {
		sections = append(sections, "organisers: organisers")
	}
	sections = append(sections, "organiserKeys: organisers[*]._key")
	if len(sections) > 0 {
		query += "    " + strings.Join(sections, ",\n    ") + "\n"
	}
	query += "  })"
	if options.Limit > 0 {
		query += "\n  LIMIT @limit"
		bindVars["limit"] = options.Limit
	}
	if options.Skip > 0 {
		query += "\n  OFFSET @skip"
		bindVars["skip"] = options.Skip
	}
	return query, bindVars
}

func buildUnsetParts(includeSet map[string]struct{}, prefix string) []string {
	var unsetParts []string
	if _, ok := includeSet[prefix+"photos"]; !ok {
		unsetParts = append(unsetParts, `"photos"`)
	}
	if _, ok := includeSet[prefix+"comments"]; !ok {
		unsetParts = append(unsetParts, `"comments"`)
	}
	return unsetParts
}

func appendUnsetPart(list []string, includeSet map[string]struct{}, key string, field string) []string {
	if _, ok := includeSet[key]; !ok {
		list = append(list, `"`+field+`"`)
	}
	return list
}

func buildUnsetString(sectionStr string, unsetParts []string) string {
	if len(unsetParts) > 0 {
		return "UNSET(" + sectionStr + ", " + strings.Join(unsetParts, ", ") + ")"
	}
	return sectionStr
}
