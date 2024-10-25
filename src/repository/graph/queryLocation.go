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
)

func (db *Db) GetLocations(options domain.LocationQueryOptions, ctx context.Context) ([]domain.LocationDTO, error) {
	query, bindVars := buildLocationQuery(options)
	return db.RunLocationQuery(query, bindVars, ctx)
}

func (db *Db) RunLocationQuery(query string, bindVars map[string]interface{}, ctx context.Context) ([]domain.LocationDTO, error) {
	cursor, err := db.Database.Query(ctx, query, &arangodb.QueryOptions{BindVars: bindVars})
	if err != nil {
		return nil, t.Errorf("query string invalid: %w", err)
	}
	defer cursor.Close()

	var result []domain.LocationDTO
	for {
		var doc domain.LocationDTO
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

func BuildImportIdQuery(source string, id string) (string, map[string]interface{}) {
	query := `
        FOR location IN locations
        FILTER location.information.importedFrom == @importedFrom AND location.information.importedId == @importedId
        RETURN location
    `
	bindVars := map[string]interface{}{
		"importedFrom": source,
		"importedId":   id,
	}
	return query, bindVars
}

func buildLocationQuery(options domain.LocationQueryOptions) (string, map[string]interface{}) {
	includeSet := options.Include
	var query string
	bindVars := map[string]interface{}{
		"lat":         options.Lat,
		"lng":         options.Lng,
		"maxDistance": options.MaxDistance,
	}
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
		query += "FOR location IN `locations-descriptions`\n"
		query += fmt.Sprintf(`  SEARCH ANALYZER(TOKENS(@text, "text_%s") ALL == location.descriptions.%s.text, "text_%s")`, lang, lang, lang)
		bindVars["text"] = options.Text
	} else {
		query += "FOR location IN locations\n"
	}
	if options.Type != "" {
		query += "\n  FILTER location.type == @type"
		bindVars["type"] = options.Type
	}
	query += "\n  LET distance = GEO_DISTANCE([@lat, @lng], [location.lat, location.lng])"
	query += "\n  FILTER distance <= @maxDistance"
	query += "\n  SORT distance"
	if options.Skip > 0 || options.Limit > 0 {
		if options.Limit == 0 {
			options.Limit = math.MaxInt
		}
		query += "\n  LIMIT @skip, @limit"
		bindVars["skip"] = options.Skip
		bindVars["limit"] = options.Limit
	}

	unsetLocation := buildUnsetParts(includeSet, "")
	unsetLocation = appendUnsetPart(unsetLocation, includeSet, "descriptions", "descriptions")
	unsetLocation = appendUnsetPart(unsetLocation, includeSet, "photos", "photos")
	unsetLocation = appendUnsetPart(unsetLocation, includeSet, "comments", "comments")
	unsetLocationStr := buildUnsetString("location", unsetLocation)

	query += "\n  RETURN MERGE(" + unsetLocationStr + ", { distance: distance })"

	return query, bindVars
}
