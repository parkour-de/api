package accounting

import (
	"bufio"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/service/accounting"
	"regexp"
)

type Handler struct {
	service *accounting.Service
}

func NewHandler(service *accounting.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddStatements(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")
	defer r.Body.Close()

	s := h.service
	bs, err := s.LoadFromJson(dpv.ConfigInstance.Server.Account)
	if err != nil {
		api.Error(w, r, fmt.Errorf("could not open accounting file: %w", err), 400)
		return
	}

	if key != bs.Key {
		api.Error(w, r, fmt.Errorf("wrong key provided"), 403)
		return
	}

	scanner := bufio.NewScanner(r.Body)
	var messages []string
	for scanner.Scan() {
		line := scanner.Text()
		messages = append(messages, line)
	}

	if err := scanner.Err(); err != nil {
		api.Error(w, r, fmt.Errorf("reading request body failed: %w", err), 400)
		return
	}

	dateRegex := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}`)

	for line, msg := range messages {
		if !dateRegex.MatchString(msg) {
			continue
		}
		err := s.UpdateBalanceSheet(&bs, msg)
		if err != nil {
			api.Error(w, r, fmt.Errorf("updating balance sheet failed, error on line %d: %w", line, err), 400)
			return
		}
	}

	err = s.SaveToJson(bs, dpv.ConfigInstance.Server.Account)
	if err != nil {
		api.Error(w, r, fmt.Errorf("could not save accounting file: %w", err), 500)
		return
	}

	api.SuccessJson(w, r, "ok")
}

func (h *Handler) GetBalanceSheetCSV(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")
	s := h.service
	bs, err := s.LoadFromJson(dpv.ConfigInstance.Server.Account)
	if err != nil {
		api.Error(w, r, fmt.Errorf("could not open accounting file: %w", err), 400)
		return
	}

	if key != bs.Key {
		api.Error(w, r, fmt.Errorf("wrong key provided"), 403)
		return
	}

	csv, err := h.service.ExportToCSV(bs)
	if err != nil {
		api.Error(w, r, fmt.Errorf("could not get balance sheet: %w", err), 500)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	api.Success(w, r, []byte(csv))
}
