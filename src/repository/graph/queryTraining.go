package graph

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/arangodb/shared"
	"math"
	"pkv/api/src/domain"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/t"
	"strings"
)

func (db *Db) GetAllTrainings(ctx context.Context) ([]domain.TrainingDTO, error) {
	return db.GetTrainings(buildAllTrainingsQuery, ctx)
}

func (db *Db) GetFilteredTrainings(options domain.TrainingQueryOptions, ctx context.Context) ([]domain.TrainingDTO, error) {
	return db.GetTrainings(func() (string, map[string]interface{}) { return buildTrainingQuery(options) }, ctx)
}

func (db *Db) GetTrainings(queryBuilder QueryBuilder, ctx context.Context) ([]domain.TrainingDTO, error) {
	query, bindVars := queryBuilder()
	cursor, err := db.Database.Query(ctx, query, &arangodb.QueryOptions{BindVars: bindVars})
	if err != nil {
		return nil, t.Errorf("query string invalid: %w", err)
	}
	defer cursor.Close()

	var result []domain.TrainingDTO
	for {
		var doc domain.TrainingDTO
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

func buildAllTrainingsQuery() (string, map[string]interface{}) {
	query := "FOR doc IN trainings RETURN doc"
	var bindVars map[string]interface{}
	return query, bindVars
}

func buildTrainingQuery(options domain.TrainingQueryOptions) (string, map[string]interface{}) {
	includeSet := options.Include
	var query string
	bindVars := make(map[string]interface{})
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
		query += "FOR training IN `trainings-descriptions`\n"
		query += fmt.Sprintf(`  SEARCH ANALYZER(TOKENS(@text, "text_%s") ALL == training.descriptions.%s.text, "text_%s")`, lang, lang, lang)
		bindVars["text"] = options.Text
	} else {
		query += "FOR training IN trainings\n"
	}
	if options.Weekday != 0 {
		query += "  FILTER training.cycles[? ANY FILTER CURRENT.weekday == @weekday]\n"
		bindVars["weekday"] = options.Weekday
	}
	unsetLocation := buildUnsetParts(includeSet, "location_")
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
	unsetOrganiser := buildUnsetParts(includeSet, "organiser_")
	organiserStr := buildUnsetString("organiser", unsetOrganiser)
	query += "  LET organisers = (FOR organiser, e IN 1..1 INBOUND training edges FILTER e.label == \"organises\" RETURN " + organiserStr + ")\n"
	if options.OrganiserKey != "" {
		query += "  FILTER @organiserKey IN organisers[*]._key\n"
		bindVars["organiserKey"] = options.OrganiserKey
	}
	if options.Skip > 0 || options.Limit > 0 {
		if options.Limit == 0 {
			options.Limit = math.MaxInt
		}
		query += "\n  LIMIT @skip, @limit"
		bindVars["skip"] = options.Skip
		bindVars["limit"] = options.Limit
	}
	unsetTraining := buildUnsetParts(includeSet, "")
	unsetTraining = appendUnsetPart(unsetTraining, includeSet, "cycles", "cycles")
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

	return query, bindVars
}
