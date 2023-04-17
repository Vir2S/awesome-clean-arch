package profile

import (
	"awesome-clean-arch/internal/handlers"
	"awesome-clean-arch/pkg/logging"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	profilesURL      = "/profile"
	profileURL       = "/profile/:username"
	createProfileURL = "/profile/create"
)

var _ handlers.Handler = &handler{}

type ErrorResponse struct {
	Error string `json:"error"`
}

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
	router.POST(createProfileURL, h.CreateProfile)
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
	w.Header().Set("Content-Type", "application/json")

	username := params.ByName("username")

	profile, err := h.repository.FindOne(context.TODO(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			errorResponse := ErrorResponse{Error: "User not found"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	profileJSON, err := json.Marshal(profile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(profileJSON)
}

func (h *handler) CreateProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var profile Profile
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	id, err := h.repository.Create(context.TODO(), profile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"id": id}
	json.NewEncoder(w).Encode(response)
}
