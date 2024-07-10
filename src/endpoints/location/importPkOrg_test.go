package location

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"os"
	"pkv/api/src/domain"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/graph"
	"pkv/api/src/service/photo"
	"reflect"
	"testing"
	"time"
)

var changedDir = false

func setup(imgDir string, tmpDir string) {
	dpv.ConfigInstance.Server.ImgPath = imgDir + "/"
	dpv.ConfigInstance.Server.TmpPath = tmpDir + "/"
}

func createHandler(t *testing.T) Handler {
	if !changedDir {
		if err := os.Chdir("../../../"); err != nil {
			t.Fatalf("failed trying to switch folders: %s", err)
		}
		changedDir = true
	}
	db, config, err := graph.Init("./config.yml", true)
	dpv.ConfigInstance = config
	if err != nil {
		t.Fatalf("db initialisation failed: %s", err)
	}
	photoService := photo.NewService()
	handler := Handler{
		db, photoService, db.Locations,
	}
	return handler
}

func TestHandler_ImportPkOrgSpot(t *testing.T) {
	h := createHandler(t)

	imgDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(imgDir)
	tmpDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	setup(imgDir, tmpDir)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("request creation failed: %s", err)
	}
	rr := httptest.NewRecorder()
	h.ImportPkOrgSpot(rr, req, httprouter.Params{})
	if rr.Code != 400 {
		t.Errorf("should have rejected missing spot id, got %v want %v", rr.Code, 400)
	}
	req, err = http.NewRequest("GET", "/?spot=3366", nil)
	if err != nil {
		t.Fatalf("request creation failed: %s", err)
	}
	rr = httptest.NewRecorder()
	h.ImportPkOrgSpot(rr, req, httprouter.Params{})
	if rr.Code != 200 {
		t.Errorf("should have succeeded, got %v want %v", rr.Code, 200)
	}
}

func TestHandler_extractPkOrgImages(t *testing.T) {
	h := createHandler(t)

	imgDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(imgDir)
	tmpDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	setup(imgDir, tmpDir)

	t.Run("file should exist", func(t *testing.T) {
		temporaryFiles, err := h.extractPkOrgImages([]Image{
			{
				1337, 0, "2014-04-09_12.44.21_andere.jpg",
			},
		}, context.Background())
		if err != nil {
			t.Errorf("expected no error, but got %v", err)
		}
		if len(temporaryFiles) != 1 {
			t.Errorf("expected 1 image, but got %d", len(temporaryFiles))
		}
	})
	t.Run("file should not exist", func(t *testing.T) {
		_, err := h.extractPkOrgImages([]Image{
			{
				1337, 0, ".well-known",
			},
		}, context.Background())
		if err == nil {
			t.Errorf("expected error, but got %v", err)
		}
	})
}

/*
	func TestHandler_extractPkOrgSpotInfo(t *testing.T) {
		type fields struct {
			db           *graph.Db
			service      *user.Service
			photoService *photo.Service
			em           graph.EntityManager
		}
		type args struct {
			url    string
			spotID string
		}
		tests := []struct {
			name    string
			fields  fields
			args    args
			want    SpotResponse
			wantErr bool
		}{
			// TODO: Add test cases.
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := &Handler{
					db:           tt.fields.db,
					service:      tt.fields.service,
					photoService: tt.fields.photoService,
					em:           tt.fields.em,
				}
				got, err := h.extractPkOrgSpotInfo(tt.args.url, tt.args.spotID)
				if (err != nil) != tt.wantErr {
					t.Errorf("extractPkOrgSpotInfo() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("extractPkOrgSpotInfo() got = %v, want %v", got, tt.want)
				}
			})
		}
	}

	func TestHandler_readPkOrgData(t *testing.T) {
		type fields struct {
			db           *graph.Db
			service      *user.Service
			photoService *photo.Service
			em           graph.EntityManager
		}
		type args struct {
			spotID string
			ctx    context.Context
		}
		tests := []struct {
			name    string
			fields  fields
			args    args
			want    Spot
			want1   []string
			wantErr bool
		}{
			// TODO: Add test cases.
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := &Handler{
					db:           tt.fields.db,
					service:      tt.fields.service,
					photoService: tt.fields.photoService,
					em:           tt.fields.em,
				}
				got, got1, err := h.readPkOrgData(tt.args.spotID, tt.args.ctx)
				if (err != nil) != tt.wantErr {
					t.Errorf("readPkOrgData() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("readPkOrgData() got = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got1, tt.want1) {
					t.Errorf("readPkOrgData() got1 = %v, want %v", got1, tt.want1)
				}
			})
		}
	}
*/
func Test_mapLocationFromPkOrgSpot(t *testing.T) {
	tests := []struct {
		name string
		spot Spot
		want domain.Location
	}{
		{"minimum fields", Spot{}, domain.Location{
			Entity: domain.Entity{
				Created:  time.UnixMilli(0).UTC(),
				Modified: time.UnixMilli(0).UTC(),
			},
			Information: map[string]string{
				"importedCategory":    "",
				"importedFrom":        "pkorg",
				"importedId":          "0",
				"importedUserCreated": "",
				"importedUserChanged": "",
			},
			Descriptions: domain.Descriptions{
				"de": {},
			},
		}},
		{"with category and description", Spot{
			Id:          42,
			Type:        "playground",
			Category:    "bouncy castle",
			Title:       "Hello World",
			Lat:         "42",
			Lng:         "-69",
			Description: "# Heading\n<b>Bold</b><script>hide me</script>",
			UserCreated: "creator",
			UserChanged: "editor",
		}, domain.Location{
			Entity: domain.Entity{
				Created:  time.UnixMilli(0).UTC(),
				Modified: time.UnixMilli(0).UTC(),
			},
			Lat:  42,
			Lng:  -69,
			Type: "playground",
			Information: map[string]string{
				"importedCategory":    "bouncy castle",
				"importedFrom":        "pkorg",
				"importedId":          "42",
				"importedUserCreated": "creator",
				"importedUserChanged": "editor",
			},
			Descriptions: domain.Descriptions{
				"de": {
					"Hello World",
					"# Heading\n<b>Bold</b><script>hide me</script>",
					"<h1 id=\"heading\">Heading</h1>\n\n<p><b>Bold</b></p>\n",
					false,
				},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapLocationFromPkOrgSpot(tt.spot); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapLocationFromPkOrgSpot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFloat(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  float64
	}{
		{"integer number", "42", 42},
		{"decimal number", "420.69", 420.69},
		{"invalid number", "invalid", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFloat(tt.value); got != tt.want {
				t.Errorf("parseFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
