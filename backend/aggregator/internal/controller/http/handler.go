package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/dto"
	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/usecase"
	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

const defaultWindowLimit = 14000

type Handler struct {
	ingestRecordsUseCase *usecase.IngestRecordsUseCase
	getWindowUseCase     *usecase.GetWindowUseCase
	getStatsUseCase      *usecase.GetStatsUseCase
	logger               *slog.Logger
}

func NewHandler(
	ingestRecordsUseCase *usecase.IngestRecordsUseCase,
	getWindowUseCase *usecase.GetWindowUseCase,
	getStatsUseCase *usecase.GetStatsUseCase,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		ingestRecordsUseCase: ingestRecordsUseCase,
		getWindowUseCase:     getWindowUseCase,
		getStatsUseCase:      getStatsUseCase,
		logger:               logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", h.handleHealth)
	mux.HandleFunc("POST /v1/records", h.handleIngestRecords)
	mux.HandleFunc("GET /v1/records/window", h.handleGetWindow)
	mux.HandleFunc("GET /v1/stats", h.handleGetStats)
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleIngestRecords(w http.ResponseWriter, r *http.Request) {
	var request dto.IngestRecordsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := h.ingestRecordsUseCase.Execute(r.Context(), dto.ToEntityRecords(request.Records)); err != nil {
		h.logger.Error("failed to ingest records", "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]int{"accepted_records": len(request.Records)})
}

func (h *Handler) handleGetWindow(w http.ResponseWriter, r *http.Request) {
	limit := defaultWindowLimit
	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = parsed
	}

	records, err := h.getWindowUseCase.Execute(r.Context(), limit)
	if err != nil {
		h.logger.Error("failed to fetch window", "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, dto.GetWindowResponse{Records: dto.FromEntityRecords(records)})
}

func (h *Handler) handleGetStats(w http.ResponseWriter, r *http.Request) {
	count, err := h.getStatsUseCase.Execute(r.Context())
	if err != nil {
		h.logger.Error("failed to fetch stats", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch stats")
		return
	}

	writeJSON(w, http.StatusOK, dto.StatsResponse{TotalRecords: count})
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, contracts.ErrorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(value)
}
