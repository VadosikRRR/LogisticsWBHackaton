package http

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/artem/logisticswbhackaton/backend/api/internal/dto"
	"github.com/artem/logisticswbhackaton/backend/api/internal/usecase"
	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type Handler struct {
	ingestUseCase    *usecase.IngestRecordsUseCase
	triggerDispatch  *usecase.TriggerDispatchUseCase
	getDashboardCase *usecase.GetDashboardOverviewUseCase
	logger           *slog.Logger
}

func NewHandler(
	ingestUseCase *usecase.IngestRecordsUseCase,
	triggerDispatch *usecase.TriggerDispatchUseCase,
	getDashboardCase *usecase.GetDashboardOverviewUseCase,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		ingestUseCase:    ingestUseCase,
		triggerDispatch:  triggerDispatch,
		getDashboardCase: getDashboardCase,
		logger:           logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", h.handleHealth)
	mux.HandleFunc("POST /v1/warehouse/records", h.handleIngestWarehouse)
	mux.HandleFunc("POST /v1/dispatch/run", h.handleRunDispatch)
	mux.HandleFunc("GET /v1/dashboard/overview", h.handleDashboardOverview)
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleIngestWarehouse(w http.ResponseWriter, r *http.Request) {
	var request dto.IngestWarehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := h.ingestUseCase.Execute(r.Context(), request); err != nil {
		h.logger.Error("failed to ingest records through gateway", "error", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]int{"accepted_records": len(request.Records)})
}

func (h *Handler) handleRunDispatch(w http.ResponseWriter, r *http.Request) {
	request, err := decodeRunDispatchRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	plan, err := h.triggerDispatch.Execute(r.Context(), request)
	if err != nil {
		h.logger.Error("failed to run dispatch from gateway", "error", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, contracts.RunDispatchResponse{Plan: plan})
}

func (h *Handler) handleDashboardOverview(w http.ResponseWriter, r *http.Request) {
	overview, err := h.getDashboardCase.Execute(r.Context())
	if err != nil {
		h.logger.Error("failed to build dashboard overview", "error", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, dto.FromEntityDashboard(overview))
}

func decodeRunDispatchRequest(r *http.Request) (dto.RunDispatchRequest, error) {
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
