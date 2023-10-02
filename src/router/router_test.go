package router

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"pkv/api/src/domain"
	"pkv/api/src/endpoints/query"
	"pkv/api/src/internal/graph"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	port := os.Getenv("PORT")
	os.Setenv("PORT", "8081")
	defer os.Setenv("PORT", port)

	server := NewServer("../../config.yml")
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Error(err)
		}
	}()
	defer server.Close()

	// Wait 50 milliseconds for server to start listening to requests
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://localhost:8081/api/users")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	expectedContentType := "application/json"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %v want %v",
			resp.Header.Get("Content-Type"), expectedContentType)
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	/*expectedBody := `"hello world"`
	if string(body) != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(body), expectedBody)
	}*/
}

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	graphDB := graph.NewTestDB()
	queryHandler := query.NewHandler(graphDB)
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		queryHandler.GetTrainings(writer, request, httprouter.Params{})
	})

	handler.ServeHTTP(rr, req)

	expectedContentType := "application/json"
	if rr.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %v want %v",
			rr.Header().Get("Content-Type"), expectedContentType)
	}

	/*expectedBody := `"hello world"`
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}*/

	/*slice, err := graphDB.GetAllUsers()
	for _, v := range slice {
		fmt.Printf("%+v\n", v)
	}*/

	x1 := time.Now()
	slice2, err := graphDB.GetTrainings(domain.TrainingQueryOptions{
		City:    "Hamburg",
		Weekday: 5,
	}, nil)
	if err != nil {
		t.Error(err)
	}
	x2 := time.Now()
	x3 := time.Now()
	fmt.Printf("Find friday trainings in Hamburg: %d ms\n", x2.Sub(x1).Milliseconds())
	fmt.Printf("Find friday trainings in Hamburg: %d ms\n", x3.Sub(x2).Milliseconds())
	fmt.Println(len(slice2))
	for _, v := range slice2 {
		b, err := json.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(b))
		//fmt.Printf("%#v\n", v)
	}
	x4 := time.Now()
	fmt.Printf("Print trainings: %d ms\n", x4.Sub(x3).Milliseconds())
}
