package contracts

import "time"

type RunDispatchRequest struct {
	WindowLimit int `json:"window_limit"`
}

type DispatchPlan struct {
	GeneratedAt      time.Time                `json:"generated_at"`
	WindowLimit      int                      `json:"window_limit"`
	RecordsScanned   int                      `json:"records_scanned"`
	PredictionsCount int                      `json:"predictions_count"`
	ModelVersion     string                   `json:"model_version,omitempty"`
	CreatedRequests  int                      `json:"created_requests"`
	Requests         []CreateTransportRequest `json:"requests"`
}

type RunDispatchResponse struct {
	Plan DispatchPlan `json:"plan"`
}
