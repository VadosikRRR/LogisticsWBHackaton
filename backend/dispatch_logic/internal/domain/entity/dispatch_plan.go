package entity

import "time"

type DispatchPlan struct {
	GeneratedAt      time.Time
	WindowLimit      int
	RecordsScanned   int
	PredictionsCount int
	ModelVersion     string
	CreatedRequests  int
	Requests         []PlannedTransportRequest
}
