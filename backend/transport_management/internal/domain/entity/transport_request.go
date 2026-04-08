package entity

import (
	"errors"
	"fmt"
	"time"
)

type Status string

const (
	StatusCreated    Status = "CREATED"
	StatusSent       Status = "SENT"
	StatusConfirmed  Status = "CONFIRMED"
	StatusInProgress Status = "IN_PROGRESS"
	StatusCompleted  Status = "COMPLETED"
	StatusCancelled  Status = "CANCELLED"
)

type TransportRequestDraft struct {
	RouteID          int64
	OfficeFromID     int64
	Timestamp        time.Time
	RequiredVehicles int
	PredictedVolume  float64
	Reason           string
}

type TransportRequest struct {
	ID               string
	RouteID          int64
	OfficeFromID     int64
	Timestamp        time.Time
	RequiredVehicles int
	PredictedVolume  float64
	Reason           string
	Status           Status
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (d TransportRequestDraft) Validate() error {
	if d.RouteID <= 0 {
		return errors.New("route_id must be positive")
	}
	if d.OfficeFromID <= 0 {
		return errors.New("office_from_id must be positive")
	}
	if d.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if d.RequiredVehicles <= 0 {
		return errors.New("required_vehicles must be positive")
	}
	return nil
}

func (r TransportRequest) IsActive() bool {
	switch r.Status {
	case StatusCreated, StatusSent, StatusConfirmed, StatusInProgress:
		return true
	default:
		return false
	}
}

func BuildRequestKey(routeID, officeFromID int64, timestamp time.Time) string {
	return fmt.Sprintf("%d:%d:%s", routeID, officeFromID, timestamp.UTC().Format(time.RFC3339))
}
