package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)

	handler.ServeHTTP(rr, req)

	expectedContentType := "application/json"
	if rr.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %v want %v",
			rr.Header().Get("Content-Type"), expectedContentType)
	}

	expectedBody := `{"name":"World","url":"/"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestServer(t *testing.T) {
	port := os.Getenv("PORT")
	os.Setenv("PORT", "8081")
	defer os.Setenv("PORT", port)

	server := NewServer()
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Error(err)
		}
	}()
	defer server.Close()

	// Wait 50 milliseconds for server to start listening to requests
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://localhost:8081/")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	expectedContentType := "application/json"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %v want %v",
			resp.Header.Get("Content-Type"), expectedContentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	expectedBody := `{"name":"World","url":"/"}`
	if string(body) != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(body), expectedBody)
	}
}
