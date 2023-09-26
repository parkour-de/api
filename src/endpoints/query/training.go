package query

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
)

// GetTrainings handles the GET /api/trainings endpoint.
//
//	@Summary		Get a list of trainings
//	@Description	Returns a list of trainings.
//	@Tags			trainings
//	@Param			weekday		query	int			false	"Day of the week (1-7) or 0 to ignore"																																			example(0)
//	@Param			city		query	string		false	"City name"																																										example(Hamburg)
//	@Param			organiser	query	string		false	"Return only trainings that match provided Organiser ID"																														example(user/135)
//	@Param			location	query	string		false	"Return only trainings that match provided Location ID"																															example(location/246)
//	@Param			include		query	[]string	false	"comma-separated list of sections to include. Choose from: cycles,photos,comments,location,location_photos,location_comments,organisers,organiser_photos,organiser_comments"	example(cycles,photos,comments,location,organisers)	collectionFormat(csv)
//	@Success		200			{array}	domain.TrainingDTO
//	@Router			/api/trainings [get]
func (h *Handler) GetTrainings(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	// Extract the include parameter from the URL query
	query := r.URL.Query()
	weekday, err := api.ParseInt(query.Get("weekday"))
	if err != nil {
		api.Error(w, fmt.Errorf("invalid weekday: %w", err), 400)
		return
	}
	skip, err := api.ParseInt(query.Get("skip"))
	if err != nil {
		api.Error(w, fmt.Errorf("invalid skip: %w", err), 400)
		return
	}
	limit, err := api.ParseInt(query.Get("limit"))
	if err != nil {
		api.Error(w, fmt.Errorf("invalid limit: %w", err), 400)
		return
	}
	queryOptions := domain.TrainingQueryOptions{
		City:        query.Get("city"),
		Weekday:     weekday,
		OrganiserID: query.Get("organiser"),
		LocationID:  query.Get("location"),
		Include:     api.MakeSet(query.Get("include")),
		Skip:        skip,
		Limit:       limit,
	}

	trainings, err := h.db.GetTrainings(queryOptions)
	jsonMsg, err := json.Marshal(trainings)
	if err != nil {
		api.Error(w, fmt.Errorf("querying trainings failed: %w", err), 400)
		return
	}

	api.Success(w, jsonMsg)
}
