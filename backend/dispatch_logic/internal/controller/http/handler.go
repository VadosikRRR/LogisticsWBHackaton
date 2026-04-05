package http

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/dto"
	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/usecase"
	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type Handler struct {
	runDispatchUseCase *usecase.RunDispatchCycleUseCase
	getLatestUseCase   *usecase.GetLatestPlanUseCase
	logger             *slog.Logger
}

func NewHandler(
	runDispatchUseCase *usecase.RunDispatchCycleUseCase,
	getLatestUseCase *usecase.GetLatestPlanUseCase,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		runDispatchUseCase: runDispatchUseCase,
		getLatestUseCase:   getLatestUseCase,
		logger:             logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", h.handleHealth)
	mux.HandleFunc("POST /v1/dispatch/run", h.handleRunDispatch)
	mux.HandleFunc("GET /v1/dispatch/plans/latest", h.handleGetLatestPlan)
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleRunDispatch(w http.ResponseWriter, r *http.Request) {
	request, err := decodeRunRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	plan, err := h.runDispatchUseCase.Execute(r.Context(), request.WindowLimit)
	if err != nil {
		h.logger.Error("dispatch run failed", "error", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, dto.RunDispatchResponse{
		Plan: dto.FromEntityPlan(plan),
	})
}

func (h *Handler) handleGetLatestPlan(w http.ResponseWriter, r *http.Request) {
	plan, err := h.getLatestUseCase.Execute(r.Context())
	if err != nil {
		h.logger.Error("failed to get latest plan", "error", err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, dto.LatestPlanResponse{
		Plan: dto.FromEntityPlan(plan),
	})
}

func decodeRunRequest(r *http.Request) (dto.RunDispatchRequest, error) {
	if r.Body == nil {
		return dto.RunDispatchRequest{}, nil
	}
	var request dto.RunDispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		if errors.Is(err, io.EOF) {
			return dto.RunDispatchRequest{}, nil
		}
		return dto.RunDispatchRequest{}, errors.New("invalid JSON body")
	}
	return request, nil
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, contracts.ErrorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(value)
}
