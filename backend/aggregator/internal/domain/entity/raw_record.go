package entity

import "time"

type RawRecord struct {
	RouteID      int64
	OfficeFromID int64
	Timestamp    time.Time
	Status1      int64
	Status2      int64
	Status3      int64
	Status4      int64
	Status5      int64
	Status6      int64
	Status7      int64
	Status8      int64
	Target2H     float64
}
