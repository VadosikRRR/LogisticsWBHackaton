package contracts

import "time"

// RawRecord mirrors one warehouse state snapshot for a route and timestamp.
type RawRecord struct {
	RouteID      int64     `json:"route_id"`
	OfficeFromID int64     `json:"office_from_id"`
	Timestamp    time.Time `json:"timestamp"`
	Status1      int64     `json:"status_1"`
	Status2      int64     `json:"status_2"`
	Status3      int64     `json:"status_3"`
	Status4      int64     `json:"status_4"`
	Status5      int64     `json:"status_5"`
	Status6      int64     `json:"status_6"`
	Status7      int64     `json:"status_7"`
	Status8      int64     `json:"status_8"`
	Target2H     float64   `json:"target_2h"`
}

type RawRecordBatch struct {
	Records []RawRecord `json:"records"`
}
