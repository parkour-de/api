package graph

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"math"
	"pkv/api/src/domain"
)

func (db *Db) GetLocations(options domain.LocationQueryOptions, ctx context.Context) ([]domain.LocationDTO, error) {
	query, bindVars := buildLocationQuery(options)
	cursor, err := db.Database.Query(ctx, query, bindVars)
	if err != nil {
		return nil, fmt.Errorf("query string invalid: %w", err)
	}
	defer cursor.Close()

	var result []domain.LocationDTO
	for {
		var doc domain.LocationDTO
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

func buildLocationQuery(options domain.LocationQueryOptions) (string, map[string]interface{}) {
	includeSet := options.Include
	query := "FOR location IN locations"
	bindVars := map[string]interface{}{
		"lat":         options.Lat,
		"lng":         options.Lng,
		"maxDistance": options.MaxDistance,
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
