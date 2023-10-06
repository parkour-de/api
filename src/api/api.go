package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func SuccessJson(w http.ResponseWriter, r *http.Request, data interface{}) {
	jsonMsg, err := json.Marshal(data)
	if err != nil {
		Error(w, r, fmt.Errorf("serialising response failed: %w", err), 400)
		return
	} else {
		Success(w, r, jsonMsg)
	}
}

func Success(w http.ResponseWriter, r *http.Request, jsonMsg []byte) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonMsg); err != nil {
		log.Printf("Error writing response: %v", err)
	}

	log.Printf(
		"%s %s %s 200",
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
	)
}

func Error(w http.ResponseWriter, r *http.Request, err error, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		err = fmt.Errorf("nil err")
	}
	logErr := err
	errorMsgJSON, err := json.Marshal(ErrorResponse{
		err.Error(),
	})
	if err != nil {
		log.Println(err)
	} else {
		if _, err = w.Write(errorMsgJSON); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}

	log.Printf(
		"%s %s %s %d %s",
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
		code,
		logErr.Error(),
	)
}

func MakeCors(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func MakeSet(queryParam string) map[string]struct{} {
	set := make(map[string]struct{})
	if queryParam != "" {
		tokens := strings.Split(queryParam, ",")
		for _, token := range tokens {
			set[token] = struct{}{}
		}
	}
	return set
}

func ParseInt(queryValue string) (int, error) {
	if queryValue == "" {
		return 0, nil
	}
	return strconv.Atoi(queryValue)
}

func ParseFloat(queryValue string) (float64, error) {
	if queryValue == "" {
		return 0, nil
	}
	return strconv.ParseFloat(queryValue, 64)
}
