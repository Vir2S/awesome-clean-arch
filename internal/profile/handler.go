package profile

import (
	"awesome-clean-arch/internal/handlers"
	"awesome-clean-arch/pkg/logging"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	profilesURL = "/profile"
	profileURL  = "/profile/:username"
)

var _ handlers.Handler = &handler{}

type handler struct {
	logger     *logging.Logger
	repository Repository
}

func NewHandler(logger *logging.Logger, repository Repository) handlers.Handler {
	return &handler{
		logger:     logger,
		repository: repository,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(profilesURL, h.GetProfilesList)
	router.GET(profileURL, h.GetProfile)
}

func (h *handler) GetProfilesList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	all, err := h.repository.FindAll(context.TODO())
	if err != nil {
		w.WriteHeader(400)
		return
	}

	allBytes, err := json.Marshal(all)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

	return
}

func (h *handler) GetProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("Username"))
}
