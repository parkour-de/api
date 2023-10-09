package query

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
)

// GetLocations handles the GET request to /api/locations
func (h *Handler) GetLocations(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	// Extract the include parameter from the URL query
	query := r.URL.Query()
	lat, err := api.ParseFloat(query.Get("lat"))
	if err != nil {
		api.Error(w, r, fmt.Errorf("invalid lat: %w", err), 400)
		return
	}
	lng, err := api.ParseFloat(query.Get("lng"))
	if err != nil {
		api.Error(w, r, fmt.Errorf("invalid lng: %w", err), 400)
		return
	}
	maxDistance, err := api.ParseFloat(query.Get("maxDistance"))
	if err != nil {
		api.Error(w, r, fmt.Errorf("invalid maxDistance: %w", err), 400)
		return
	}
	skip, err := api.ParseInt(query.Get("skip"))
	if err != nil {
		api.Error(w, r, fmt.Errorf("invalid skip: %w", err), 400)
		return
	}
	limit, err := api.ParseInt(query.Get("limit"))
	if err != nil {
		api.Error(w, r, fmt.Errorf("invalid limit: %w", err), 400)
		return
	}
	queryOptions := domain.LocationQueryOptions{
		Lat:         lat,
		Lng:         lng,
		MaxDistance: maxDistance,
		Type:        query.Get("type"),
		Include:     api.MakeSet(query.Get("include")),
		Skip:        skip,
		Limit:       limit,
	}

	locations, err := h.db.GetLocations(queryOptions, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("querying locations failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, locations)
}
