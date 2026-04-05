package dto

import (
	"fmt"

	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/entity"
)

type CreateRequestsRequest = contracts.CreateTransportRequestsRequest
type CreateRequestsResponse = contracts.CreateTransportRequestsResponse

type ActiveRequestsResponse struct {
	Requests []contracts.TransportRequest `json:"requests"`
}

type UpdateStatusRequest = contracts.UpdateTransportRequestStatusRequest

func ToDrafts(requests []contracts.CreateTransportRequest) []entity.TransportRequestDraft {
	out := make([]entity.TransportRequestDraft, 0, len(requests))
	for _, request := range requests {
		out = append(out, entity.TransportRequestDraft{
			RouteID:          request.RouteID,
			OfficeFromID:     request.OfficeFromID,
			Timestamp:        request.Timestamp,
			RequiredVehicles: request.RequiredVehicles,
			PredictedVolume:  request.PredictedVolume,
			Reason:           request.Reason,
		})
	}
	return out
}

func FromEntityRequests(requests []entity.TransportRequest) []contracts.TransportRequest {
	out := make([]contracts.TransportRequest, 0, len(requests))
	for _, request := range requests {
		out = append(out, contracts.TransportRequest{
			ID:               request.ID,
			RouteID:          request.RouteID,
			OfficeFromID:     request.OfficeFromID,
			Timestamp:        request.Timestamp,
			RequiredVehicles: request.RequiredVehicles,
			PredictedVolume:  request.PredictedVolume,
			Status:           contracts.TransportRequestStatus(request.Status),
			Reason:           request.Reason,
			CreatedAt:        request.CreatedAt,
			UpdatedAt:        request.UpdatedAt,
		})
	}
	return out
}

func ToEntityStatus(status contracts.TransportRequestStatus) (entity.Status, error) {
	switch status {
	case contracts.TransportRequestStatusCreated:
		return entity.StatusCreated, nil
	case contracts.TransportRequestStatusSent:
		return entity.StatusSent, nil
	case contracts.TransportRequestStatusConfirmed:
		return entity.StatusConfirmed, nil
	case contracts.TransportRequestStatusInProgress:
		return entity.StatusInProgress, nil
	case contracts.TransportRequestStatusCompleted:
		return entity.StatusCompleted, nil
	case contracts.TransportRequestStatusCancelled:
		return entity.StatusCancelled, nil
	default:
		return "", fmt.Errorf("unknown status %q", status)
	}
}
