package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"egot-tracker/internal/service"
	"egot-tracker/pkg/response"

	"github.com/jackc/pgx/v5/pgtype"
)

type OscarHandler struct {
	service *service.OscarService
}

func NewOscarHandler(service *service.OscarService) *OscarHandler {
	return &OscarHandler{service: service}
}

// GetCeremony handles GET /api/oscar-race/{year}
func (h *OscarHandler) GetCeremony(w http.ResponseWriter, r *http.Request) {
	// Extract year from path
	yearStr := r.PathValue("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1929 || year > 2100 {
		response.Error(w, http.StatusBadRequest, "invalid year")
		return
	}

	ceremony, err := h.service.GetCeremony(r.Context(), year)
	if errors.Is(err, service.ErrCeremonyNotFound) {
		response.Error(w, http.StatusNotFound, "ceremony not found")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, ceremony)
}

// GetYears handles GET /api/oscar-race/years
func (h *OscarHandler) GetYears(w http.ResponseWriter, r *http.Request) {
	years, err := h.service.GetAllYears(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if years == nil {
		years = []int{}
	}

	response.JSON(w, http.StatusOK, years)
}

// SetWinner handles PUT /api/oscar-race/{year}/category/{categoryId}/winner/{nomineeId}
func (h *OscarHandler) SetWinner(w http.ResponseWriter, r *http.Request) {
	// Extract nomineeId from path
	nomineeIdStr := r.PathValue("nomineeId")

	// Parse UUID
	var nomineeID pgtype.UUID
	if err := nomineeID.Scan(nomineeIdStr); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid nominee ID")
		return
	}

	// Validate UUID is set
	if !nomineeID.Valid {
		response.Error(w, http.StatusBadRequest, "nominee ID is required")
		return
	}

	// Set the winner
	if err := h.service.SetWinner(r.Context(), nomineeID); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			response.Error(w, http.StatusNotFound, "nominee not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "winner set"})
}
