package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"egot-tracker/internal/models"
	"egot-tracker/internal/service"
	"egot-tracker/pkg/response"
)

type CelebrityHandler struct {
	service *service.CelebrityService
}

func NewCelebrityHandler(service *service.CelebrityService) *CelebrityHandler {
	return &CelebrityHandler{service: service}
}

func (h *CelebrityHandler) Search(w http.ResponseWriter, r *http.Request) {
	// Validate query parameter
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		response.Error(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	// Call service layer
	result, err := h.service.SearchCelebrity(r.Context(), query)
	if errors.Is(err, service.ErrCelebrityNotFound) {
		response.Error(w, http.StatusNotFound, "celebrity not found")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *CelebrityHandler) Autocomplete(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		response.JSON(w, http.StatusOK, []models.Celebrity{})
		return
	}

	results, err := h.service.Autocomplete(r.Context(), query, 10)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if results == nil {
		results = []models.Celebrity{}
	}

	response.JSON(w, http.StatusOK, results)
}

// CloseToEGOT handles GET /api/celebrity/close-to-egot
func (h *CelebrityHandler) CloseToEGOT(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	results, err := h.service.GetCloseToEGOT(r.Context(), limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if results == nil {
		results = []models.CelebrityWithEGOTProgress{}
	}

	response.JSON(w, http.StatusOK, results)
}
