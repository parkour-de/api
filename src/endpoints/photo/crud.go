package photo

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

func photosFromRequest(r *http.Request) ([]string, error) {
	var items []string
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&items); err != nil {
		return items, fmt.Errorf("decoding request body failed: %w", err)
	}
	return items, nil
}

func (h *PhotoEntityHandler[T]) UpdatePhotos(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	files, err := photosFromRequest(r)
	if err != nil {
		api.Error(w, r, fmt.Errorf("updating user photos failed: %w", err), 400)
		return
	}
	entity, err := h.em.Read(urlParams.ByName("key"), r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("updating user photos failed: %w", err), 400)
		return
	}
	photos, err := h.service.Update(entity.GetPhotos(), files, r.Context())
	oldPhotos := entity.GetPhotos()
	entity.SetPhotos(photos)
	err = h.em.Update(entity, r.Context())
	if err != nil {
		var oldPhotosFiles []string
		for _, photo := range oldPhotos {
			oldPhotosFiles = append(oldPhotosFiles, photo.Src)
		}
		_, err2 := h.service.Update(entity.GetPhotos(), oldPhotosFiles, r.Context())
		if err2 != nil {
			api.Error(w, r, fmt.Errorf("saving updated user photos failed, additionally an error occured while rolling back file changes: %w, %v", err, err2), 400)
			return
		} else {
			api.Error(w, r, fmt.Errorf("saving updated user photos failed, changes to files have been rolled back: %w", err), 400)
		}
		return
	}
	api.SuccessJson(w, r, entity)
}
