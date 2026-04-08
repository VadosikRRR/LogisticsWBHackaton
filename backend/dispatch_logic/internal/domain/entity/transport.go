package entity

import (
	"fmt"
	"time"
)

type ActiveTransportRequest struct {
	RouteID          int64
	OfficeFromID     int64
	Timestamp        time.Time
	RequiredVehicles int
}

type PlannedTransportRequest struct {
	RouteID          int64
	OfficeFromID     int64
	Timestamp        time.Time
	RequiredVehicles int
	PredictedVolume  float64
	Reason           string
}

func BuildActiveRequestKey(routeID, officeFromID int64, timestamp time.Time) string {
	return BuildTripletKey(routeID, officeFromID, timestamp)
}

func BuildTripletKey(routeID, officeFromID int64, timestamp time.Time) string {
	return timestamp.UTC().Format(time.RFC3339) + "|" + int64ToString(routeID) + "|" + int64ToString(officeFromID)
}

func int64ToString(v int64) string {
	return fmt.Sprintf("%d", v)
}
