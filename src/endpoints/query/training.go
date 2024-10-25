package query

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
	"pkv/api/src/repository/t"
)

// GetTrainings handles the GET /api/trainings endpoint.
func (h *Handler) GetTrainings(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	query := r.URL.Query()
	weekday, err := api.ParseInt(query.Get("weekday"))
	if err != nil {
		api.Error(w, r, t.Errorf("invalid weekday: %w", err), 400)
		return
	}
	skip, err := api.ParseInt(query.Get("skip"))
	if err != nil {
		api.Error(w, r, t.Errorf("invalid skip: %w", err), 400)
		return
	}
	limit, err := api.ParseInt(query.Get("limit"))
	if err != nil {
		api.Error(w, r, t.Errorf("invalid limit: %w", err), 400)
		return
	}
	queryOptions := domain.TrainingQueryOptions{
		City:         query.Get("city"),
		Weekday:      weekday,
		OrganiserKey: query.Get("organiser"),
		LocationKey:  query.Get("location"),
		Type:         query.Get("type"),
		Text:         query.Get("text"),
		Language:     query.Get("language"),
		Include:      api.MakeSet(query.Get("include")),
		Skip:         skip,
		Limit:        limit,
	}

	trainings, err := h.db.GetFilteredTrainings(queryOptions, r.Context())
	if err != nil {
		api.Error(w, r, t.Errorf("querying trainings failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, trainings)
}

func (h *Handler) GetTraining(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")
	item, err := h.db.Trainings.Read(key, r.Context())
	if err != nil {
		api.Error(w, r, t.Errorf("read request failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, item)
}
