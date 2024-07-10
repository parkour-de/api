package location

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
	"pkv/api/src/repository/graph"
	"pkv/api/src/service/description"
	"strconv"
	"strings"
	"time"
)

type SpotResponse struct {
	Spot                Spot                 `json:"spot"`
	Images              []Image              `json:"images"`
	SpotCategoryDetails []SpotCategoryDetail `json:"spot_category_details"`
}

type Spot struct {
	Id          int    `json:"id"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Created     int64  `json:"created"`
	Changed     int64  `json:"changed"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
	Geohash     string `json:"geohash"`
	Zoom        int    `json:"zoom"`
	P0          string `json:"p0"`
	Description string `json:"description"`
	UserId      int    `json:"user_id"`
	UrlAlias    string `json:"url_alias"`
	UserCreated string `json:"user_created"`
	UserChanged string `json:"user_changed"`
}

type Image struct {
	SpotId   int    `json:"spot_id"`
	Delta    int    `json:"delta"`
	Filename string `json:"filename"`
}

type SpotCategoryDetail struct{}

func (h *Handler) ImportPkOrgSpot(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	spotID := r.URL.Query().Get("spot")
	existing, err := h.isPkOrgLocationExisting(spotID, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("checking for existing locations failed: %w", err), 400)
		return
	}
	if existing {
		api.Error(w, r, fmt.Errorf("location already found in database"), 409)
		return
	}
	if spotID == "" {
		api.Error(w, r, fmt.Errorf("missing 'spot' query parameter"), 400)
		return
	}
	spot, filenames, err := h.readPkOrgData(spotID, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to extract information from PkOrg spot %s: %w", spotID, err), 500)
		return
	}

	photos, err := h.photoService.Update([]domain.Photo{}, filenames, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to update photos for spot %s: %w", spotID, err), 500)
		return
	}

	location := mapLocationFromPkOrgSpot(spot)
	location.Photos = domain.Photos{Photos: photos}

	err = h.em.Create(&location, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to create location for spot %s: %w", spotID, err), 500)
		return
	}

	api.SuccessJson(w, r, location.Key)
}

func (h *Handler) readPkOrgData(spotID string, ctx context.Context) (Spot, []string, error) {
	url := fmt.Sprintf("https://map.parkour.org/api/v1/spot/%s", spotID)

	spotResponse, err := h.extractPkOrgSpotInfo(url)
	if err != nil {
		return Spot{}, nil, err
	}

	filenames, err := h.extractPkOrgImages(spotResponse.Images, ctx)
	return spotResponse.Spot, filenames, err
}

func (h *Handler) isPkOrgLocationExisting(spotID string, ctx context.Context) (bool, error) {
	query, bindVars := graph.BuildImportIdQuery("pkorg", spotID)
	locations, err := h.db.RunLocationQuery(query, bindVars, ctx)
	if err != nil {
		return false, err
	}
	return len(locations) > 0, nil
}

func (h *Handler) extractPkOrgSpotInfo(url string) (SpotResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return SpotResponse{}, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SpotResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var spotResponse SpotResponse
	err = json.Unmarshal(body, &spotResponse)
	if err != nil {
		return SpotResponse{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return spotResponse, err
}

func (h *Handler) extractPkOrgImages(images []Image, ctx context.Context) ([]string, error) {
	var filenames []string
	var errors []string
	for i, img := range images {
		imageURL := fmt.Sprintf("https://map.parkour.org/images/spots/%s", img.Filename)
		resp, err := http.Get(imageURL)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to fetch image data for image %d: %w", i, err).Error())
			continue
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to read image data for image %d: %w", i, err).Error())
			continue
		}

		photo, err := h.photoService.Upload(data, img.Filename, ctx)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to upload photo for image %d: %w", i, err).Error())
			continue
		}

		filenames = append(filenames, photo.Src)
	}
	var err error
	if len(errors) > 0 {
		err = fmt.Errorf("errors occured with spot images: %v", strings.Join(errors, "; "))
	}
	return filenames, err
}

func parseFloat(value string) float64 {
	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return parsedValue
}

func mapLocationFromPkOrgSpot(spot Spot) domain.Location {
	location := domain.Location{
		Entity: domain.Entity{
			Created:  time.Unix(spot.Created, 0).UTC(),
			Modified: time.Unix(spot.Changed, 0).UTC(),
		},
		Lat:  parseFloat(spot.Lat),
		Lng:  parseFloat(spot.Lng),
		Type: spot.Type,
		Information: map[string]string{
			"importedFrom":        "pkorg",
			"importedId":          fmt.Sprintf("%d", spot.Id),
			"importedCategory":    spot.Category,
			"importedUserCreated": spot.UserCreated,
			"importedUserChanged": spot.UserChanged,
		},
		Descriptions: domain.Descriptions{
			"de": {
				Title:  spot.Title,
				Text:   spot.Description,
				Render: description.Render([]byte(spot.Description)),
			},
		},
	}
	return location
}
