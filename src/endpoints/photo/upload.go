package photo

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"pkv/api/src/api"
)

type SomeStruct struct{}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		api.Error(w, r, fmt.Errorf("parsing multipart form failed: %v", err), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		api.Error(w, r, fmt.Errorf("getting uploaded file failed: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := header.Filename

	data, err := io.ReadAll(file)
	if err != nil {
		api.Error(w, r, fmt.Errorf("reading uploaded file failed: %v", err), http.StatusInternalServerError)
		return
	}

	photo, err := h.service.Upload(data, filename, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, photo)
}
