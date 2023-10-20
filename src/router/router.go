package router

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"pkv/api/src/api"
	"pkv/api/src/domain"
	"pkv/api/src/endpoints/authentication"
	"pkv/api/src/endpoints/crud"
	"pkv/api/src/endpoints/query"
	"pkv/api/src/endpoints/user"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/graph"
	user2 "pkv/api/src/service/user"
)

func NewServer(configPath string, test bool) *http.Server {
	db, config, err := graph.Init(configPath, test)
	if err != nil {
		log.Fatal(err)
	}
	dpv.ConfigInstance = config

	userService := user2.NewService(db)
	authenticationHandler := authentication.NewHandler(db, userService)
	queryHandler := query.NewHandler(db)
	userHandler := user.NewHandler(db, userService)
	trainingCrudHandler := crud.NewHandler[*domain.Training](db, db.Trainings)
	locationCrudHandler := crud.NewHandler[*domain.Location](db, db.Locations)
	userCrudHandler := crud.NewHandler[*domain.User](db, db.Users)
	pageCrudHandler := crud.NewHandler[*domain.Page](db, db.Pages)

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

	r.POST("/api/admin/trainings", trainingCrudHandler.Create)
	r.GET("/api/admin/trainings/:key", trainingCrudHandler.Read)
	r.PUT("/api/admin/trainings", trainingCrudHandler.Update)
	r.DELETE("/api/admin/trainings/:key", trainingCrudHandler.Delete)

	r.POST("/api/admin/locations", locationCrudHandler.Create)
	r.GET("/api/admin/locations/:key", locationCrudHandler.Read)
	r.PUT("/api/admin/locations", locationCrudHandler.Update)
	r.DELETE("/api/admin/locations/:key", locationCrudHandler.Delete)

	r.POST("/api/admin/users", userCrudHandler.Create)
	r.GET("/api/admin/users/:key", userCrudHandler.Read)
	r.PUT("/api/admin/users", userCrudHandler.Update)
	r.DELETE("/api/admin/users/:key", userCrudHandler.Delete)

	r.POST("/api/admin/pages", pageCrudHandler.Create)
	r.GET("/api/admin/pages/:key", pageCrudHandler.Read)
	r.PUT("/api/admin/pages", pageCrudHandler.Update)
	r.DELETE("/api/admin/pages/:key", pageCrudHandler.Delete)

	r.GET("/api/facebook", authenticationHandler.Facebook)
	r.GET("/api/trainings", queryHandler.GetTrainings)
	r.GET("/api/trainings/:key", queryHandler.GetTraining)
	r.GET("/api/pages", queryHandler.GetPages)
	r.GET("/api/pages/:key", queryHandler.GetPage)
	r.GET("/api/locations", queryHandler.GetLocations)
	r.GET("/api/locations/:key", queryHandler.GetLocation)
	r.GET("/api/users", queryHandler.GetUsers)
	r.GET("/api/users/:key", queryHandler.GetUser)
	r.POST("/api/users/:key", userHandler.Create)
	r.GET("/api/users/:key/exists", userHandler.Exists)
	r.POST("/api/users/:key/claim", userHandler.Claim)
	r.GET("/api/users/:key/facebook", userHandler.LinkFacebook)
	r.GET("/api/users/:key/password", userHandler.LinkPassword)
	r.GET("/api/users/:key/totp", userHandler.RequestTOTP)
	r.POST("/api/users/:key/totp", userHandler.EnableTOTP)
	r.GET("/api/users/:key/email", userHandler.RequestEmail)
	r.GET("/api/users/:key/email/:login", userHandler.EnableEmail)

	r.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Printf("panic: %+v", err)
		api.Error(w, r, fmt.Errorf("Whoops! It seems we've stumbled upon a glitch here. In the meantime, consider this a chance to take a breather."), http.StatusInternalServerError)
	}
	r.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.Error(w, r, fmt.Errorf("Oops, your %v move is impressive, but this method doesn't match the route's rhythm. Let's stick to the right Parkour technique – we've got OPTIONS waiting for you, not this wild %v dance!", r.Method, r.Method), http.StatusMethodNotAllowed)
	})
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.Error(w, r, fmt.Errorf("Oops, you're performing a daring stunt! But this route seems to be off our servers. Maybe let's stick to known paths for now and avoid tumbling into the broken API!"), http.StatusNotFound)
	})

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
