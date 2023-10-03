package query

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
)

// GetTrainings handles the GET /api/trainings endpoint.
func (h *Handler) GetTrainings(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	// Extract the include parameter from the URL query
	query := r.URL.Query()
	weekday, err := api.ParseInt(query.Get("weekday"))
	if err != nil {
		api.Error(w, r, fmt.Errorf("invalid weekday: %w", err), 400)
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
	queryOptions := domain.TrainingQueryOptions{
		City:         query.Get("city"),
		Weekday:      weekday,
		OrganiserKey: query.Get("organiser"),
		LocationKey:  query.Get("location"),
		Include:      api.MakeSet(query.Get("include")),
		Skip:         skip,
		Limit:        limit,
	}

	trainings, err := h.db.GetTrainings(queryOptions, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("querying trainings failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, trainings)
}
