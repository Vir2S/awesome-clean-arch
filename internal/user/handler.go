package user

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
	usersURL = "/user"
	userURL  = "/user/:id"
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
	router.GET(usersURL, h.GetUsersList)
	router.GET(userURL, h.GetUser)
}

func (h *handler) GetUsersList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	userList, err := h.repository.FindAll(context.TODO())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	userListJSON, err := json.Marshal(userList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userListJSON)
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	userID := params.ByName("id")

	user, err := h.repository.FindOne(context.TODO(), userID)
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

	userJSON, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}
