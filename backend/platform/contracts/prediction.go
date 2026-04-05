package contracts

import "time"

type PredictRequest struct {
	Records []RawRecord `json:"records"`
}

type Prediction struct {
	RouteID      int64     `json:"route_id"`
	OfficeFromID int64     `json:"office_from_id"`
	Timestamp    time.Time `json:"timestamp"`
	Target2H     float64   `json:"target_2h"`
}

type PredictResponse struct {
	Predictions  []Prediction `json:"predictions"`
	ModelVersion string       `json:"model_version,omitempty"`
}
