package router

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"pkv/api/src/domain"
	"pkv/api/src/endpoints/authentication"
	"pkv/api/src/endpoints/crud"
	"pkv/api/src/endpoints/query"
	"pkv/api/src/internal/dpv"
	"pkv/api/src/internal/graph"
)

func NewServer(configPath string) *http.Server {
	var err error
	dpv.ConfigInstance, err = dpv.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	c, err := graph.Connect(dpv.ConfigInstance, true)
	if err != nil {
		log.Fatal(err)
	}
	database, err := graph.GetOrCreateDatabase(c, "dpv", dpv.ConfigInstance)
	if err != nil {
		log.Fatal(err)
	}

	db, err := graph.NewDB(database)
	authenticationHandler := authentication.NewHandler(db)
	queryHandler := query.NewHandler(db)
	trainingCrudHandler := crud.NewHandler[*domain.Training](db, db.Trainings, "")
	locationCrudHandler := crud.NewHandler[*domain.Location](db, db.Locations, "")
	userCrudHandler := crud.NewHandler[*domain.User](db, db.Users, "")
	pageCrudHandler := crud.NewHandler[*domain.Page](db, db.Pages, "")

	r := httprouter.New()

	r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		}

		w.WriteHeader(http.StatusNoContent)
	})

	r.GET("/api/facebook", authenticationHandler.Facebook)
	r.GET("/api/users", queryHandler.GetUsers)
	r.GET("/api/trainings", queryHandler.GetTrainings)
	r.GET("/api/pages", queryHandler.GetPages)
	r.POST("/api/trainings", trainingCrudHandler.Create)
	r.GET("/api/trainings/:id", trainingCrudHandler.Read)
	r.PUT("/api/trainings", trainingCrudHandler.Update)
	r.DELETE("/api/trainings/:id", trainingCrudHandler.Delete)
	r.POST("/api/locations", locationCrudHandler.Create)
	r.GET("/api/locations/:id", locationCrudHandler.Read)
	r.PUT("/api/locations", locationCrudHandler.Update)
	r.DELETE("/api/locations/:id", locationCrudHandler.Delete)
	r.POST("/api/users", userCrudHandler.Create)
	r.GET("/api/users/:id", userCrudHandler.Read)
	r.PUT("/api/users", userCrudHandler.Update)
	r.DELETE("/api/users/:id", userCrudHandler.Delete)
	r.POST("/api/pages", pageCrudHandler.Create)
	r.GET("/api/pages/:id", pageCrudHandler.Read)
	r.PUT("/api/pages", pageCrudHandler.Update)
	r.DELETE("/api/pages/:id", pageCrudHandler.Delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	addr := "localhost:" + port
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
