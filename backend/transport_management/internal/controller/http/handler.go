package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/dto"
	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/usecase"
)

type Handler struct {
	createBulkUseCase *usecase.CreateBulkRequestsUseCase
	listActiveUseCase *usecase.ListActiveRequestsUseCase
	updateStatusCase  *usecase.UpdateStatusUseCase
	logger            *slog.Logger
}

func NewHandler(
	createBulkUseCase *usecase.CreateBulkRequestsUseCase,
	listActiveUseCase *usecase.ListActiveRequestsUseCase,
	updateStatusCase *usecase.UpdateStatusUseCase,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		createBulkUseCase: createBulkUseCase,
		listActiveUseCase: listActiveUseCase,
		updateStatusCase:  updateStatusCase,
		logger:            logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", h.handleHealth)
	mux.HandleFunc("POST /v1/requests/bulk", h.handleCreateBulk)
	mux.HandleFunc("GET /v1/requests/active", h.handleListActive)
	mux.HandleFunc("PATCH /v1/requests/{id}/status", h.handleUpdateStatus)
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleCreateBulk(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateRequestsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	created, err := h.createBulkUseCase.Execute(r.Context(), dto.ToDrafts(request.Requests))
	if err != nil {
		h.logger.Error("failed to create bulk requests", "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, dto.CreateRequestsResponse{
		Requests: dto.FromEntityRequests(created),
	})
}

func (h *Handler) handleListActive(w http.ResponseWriter, r *http.Request) {
	requests, err := h.listActiveUseCase.Execute(r.Context())
	if err != nil {
		h.logger.Error("failed to list active requests", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list active requests")
		return
	}

	writeJSON(w, http.StatusOK, dto.ActiveRequestsResponse{
		Requests: dto.FromEntityRequests(requests),
	})
}

func (h *Handler) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing request id")
		return
	}

	var request dto.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	status, err := dto.ToEntityStatus(request.Status)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.updateStatusCase.Execute(r.Context(), id, status)
	if err != nil {
		h.logger.Error("failed to update request status", "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, contracts.TransportRequest{
		ID:               updated.ID,
		RouteID:          updated.RouteID,
		OfficeFromID:     updated.OfficeFromID,
		Timestamp:        updated.Timestamp,
		RequiredVehicles: updated.RequiredVehicles,
		PredictedVolume:  updated.PredictedVolume,
		Status:           contracts.TransportRequestStatus(updated.Status),
		Reason:           updated.Reason,
		CreatedAt:        updated.CreatedAt,
		UpdatedAt:        updated.UpdatedAt,
	})
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, contracts.ErrorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(value)
}
