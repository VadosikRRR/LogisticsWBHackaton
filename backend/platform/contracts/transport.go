package contracts

import "time"

type TransportRequestStatus string

const (
	TransportRequestStatusCreated    TransportRequestStatus = "CREATED"
	TransportRequestStatusSent       TransportRequestStatus = "SENT"
	TransportRequestStatusConfirmed  TransportRequestStatus = "CONFIRMED"
	TransportRequestStatusInProgress TransportRequestStatus = "IN_PROGRESS"
	TransportRequestStatusCompleted  TransportRequestStatus = "COMPLETED"
	TransportRequestStatusCancelled  TransportRequestStatus = "CANCELLED"
)

type CreateTransportRequest struct {
	RouteID          int64     `json:"route_id"`
	OfficeFromID     int64     `json:"office_from_id"`
	Timestamp        time.Time `json:"timestamp"`
	RequiredVehicles int       `json:"required_vehicles"`
	PredictedVolume  float64   `json:"predicted_volume"`
	Reason           string    `json:"reason,omitempty"`
}

type CreateTransportRequestsRequest struct {
	Requests []CreateTransportRequest `json:"requests"`
}

type TransportRequest struct {
	ID               string                 `json:"id"`
	RouteID          int64                  `json:"route_id"`
	OfficeFromID     int64                  `json:"office_from_id"`
	Timestamp        time.Time              `json:"timestamp"`
	RequiredVehicles int                    `json:"required_vehicles"`
	PredictedVolume  float64                `json:"predicted_volume"`
	Status           TransportRequestStatus `json:"status"`
	Reason           string                 `json:"reason,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type CreateTransportRequestsResponse struct {
	Requests []TransportRequest `json:"requests"`
}

type UpdateTransportRequestStatusRequest struct {
	Status TransportRequestStatus `json:"status"`
}
