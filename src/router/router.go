package router

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"pkv/api/src/api"
	"pkv/api/src/domain"
	"pkv/api/src/endpoints/accounting"
	"pkv/api/src/endpoints/authentication"
	"pkv/api/src/endpoints/crud"
	"pkv/api/src/endpoints/location"
	"pkv/api/src/endpoints/photo"
	"pkv/api/src/endpoints/query"
	"pkv/api/src/endpoints/server"
	"pkv/api/src/endpoints/user"
	"pkv/api/src/endpoints/verband"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/graph"
	accountingService "pkv/api/src/service/accounting"
	photoService "pkv/api/src/service/photo"
	serverService "pkv/api/src/service/server"
	userService "pkv/api/src/service/user"
	verbandService "pkv/api/src/service/verband"
	"time"
)

func NewServer(configPath string, test bool) *http.Server {
	attempts := 0
	if !test {
		attempts = 5
	}
	db, config, err := graph.Init(configPath, test)
	for attempt := 0; attempt < attempts; attempt++ {
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * time.Duration(attempt*attempt)) // 0, 1, 4, 9, 16 seconds
			db, config, err = graph.Init(configPath, test)
			continue
		}
		break
	}
	if err != nil {
		log.Fatal(err)
	}
	dpv.ConfigInstance = config

	trainingCrudHandler := crud.NewHandler[*domain.Training](db, db.Trainings)
	locationCrudHandler := crud.NewHandler[*domain.Location](db, db.Locations)
	userCrudHandler := crud.NewHandler[*domain.User](db, db.Users)
	pageCrudHandler := crud.NewHandler[*domain.Page](db, db.Pages)

	userService := userService.NewService(db)
	authenticationHandler := authentication.NewHandler(db, userService)
	queryHandler := query.NewHandler(db)
	userHandler := user.NewHandler(db, userService)
	userPhotoHandler := photo.NewPhotoEntityHandler[*domain.User](photoService.NewService(), db.Users)

	photoService := photoService.NewService()

	serverHandler := server.NewHandler(serverService.NewService())
	photoHandler := photo.NewHandler(photoService)
	locationHandler := location.NewHandler(db, photoService, db.Locations)

	accountingHandler := accounting.NewHandler(accountingService.NewService())
	verbandHandler := verband.NewHandler(verbandService.NewService())

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

	r.GET("/api/version", Version)

	r.POST("/api/admin/training", trainingCrudHandler.Create)
	r.GET("/api/admin/training/:key", trainingCrudHandler.Read)
	r.PUT("/api/admin/training", trainingCrudHandler.Update)
	r.DELETE("/api/admin/training/:key", trainingCrudHandler.Delete)

	r.POST("/api/admin/location", locationCrudHandler.Create)
	r.GET("/api/admin/location/:key", locationCrudHandler.Read)
	r.PUT("/api/admin/location", locationCrudHandler.Update)
	r.DELETE("/api/admin/location/:key", locationCrudHandler.Delete)

	r.POST("/api/admin/user", userCrudHandler.Create)
	r.GET("/api/admin/user/:key", userCrudHandler.Read)
	r.PUT("/api/admin/user", userCrudHandler.Update)
	r.DELETE("/api/admin/user/:key", userCrudHandler.Delete)

	r.POST("/api/admin/page", pageCrudHandler.Create)
	r.GET("/api/admin/page/:key", pageCrudHandler.Read)
	r.PUT("/api/admin/page", pageCrudHandler.Update)
	r.DELETE("/api/admin/page/:key", pageCrudHandler.Delete)

	r.GET("/api/login/facebook", authenticationHandler.Facebook)

	r.POST("/api/locations/import/pkorg", locationHandler.ImportPkOrgSpot)

	r.GET("/api/training", queryHandler.GetTrainings)
	r.GET("/api/training/:key", queryHandler.GetTraining)
	r.GET("/api/page", queryHandler.GetPages)
	r.GET("/api/page/:key", queryHandler.GetPage)
	r.GET("/api/location", queryHandler.GetLocations)
	r.GET("/api/location/:key", queryHandler.GetLocation)
	r.GET("/api/user", queryHandler.GetUsers)
	r.GET("/api/user/:key", queryHandler.GetUser)
	r.POST("/api/user/:key", userHandler.Create)
	r.GET("/api/user/:key/exists", userHandler.Exists)
	r.POST("/api/user/:key/claim", userHandler.Claim)
	r.GET("/api/user/:key/facebook", userHandler.LinkFacebook)
	r.GET("/api/user/:key/password", userHandler.LinkPassword)
	r.GET("/api/user/:key/totp", userHandler.RequestTOTP)
	r.POST("/api/user/:key/totp", userHandler.EnableTOTP)
	r.GET("/api/user/:key/email", userHandler.RequestEmail)
	r.GET("/api/user/:key/email/:login", userHandler.EnableEmail)
	r.POST("/api/user/:key/photos", userPhotoHandler.UpdatePhotos)

	r.POST("/api/user/:key/comment", userHandler.AddComment)
	r.PUT("/api/user/:key/comment", userHandler.EditComment)
	r.DELETE("/api/user/:key/comment", userHandler.DeleteComment)

	r.POST("/api/server/mail", serverHandler.ChangeMailPassword)
	r.POST("/api/server/minecraft/whitelist", serverHandler.AddUsernameToWhitelist)

	r.POST("/api/photo/upload", photoHandler.Upload)

	r.GET("/api/account/:key", accountingHandler.GetBalanceSheetCSV)
	r.POST("/api/account/:key", accountingHandler.AddStatements)

	r.GET("/api/verband/vereine", verbandHandler.GetVereine)
	r.GET("/api/verband/bundeslaender", verbandHandler.GetBundeslaender)

	r.POST("/api/verband/mitmachen", verbandHandler.Mitmachen)

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

func Version(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	// the only endpoint that does not use JSON-formatted response, i.e. no quotes around version string
	api.Success(w, r, []byte(dpv.ConfigInstance.Settings.Version))
}
