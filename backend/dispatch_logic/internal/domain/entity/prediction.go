package entity

import "time"

type Prediction struct {
	RouteID      int64
	OfficeFromID int64
	Timestamp    time.Time
	Target2H     float64
}

func BuildPredictionKey(routeID, officeFromID int64, timestamp time.Time) string {
	return BuildTripletKey(routeID, officeFromID, timestamp)
}
