package client

import (
	"context"

	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type MockMLClient struct{}

func NewMockMLClient() *MockMLClient {
	return &MockMLClient{}
}

func (c *MockMLClient) Predict(_ context.Context, records []contracts.RawRecord) (contracts.PredictResponse, error) {
	predictions := make([]contracts.Prediction, 0, len(records))
	for _, record := range records {
		volume := float64(
			record.Status1+record.Status2+record.Status3+record.Status4+
				record.Status5+record.Status6+record.Status7+record.Status8,
		) / 4.0
		predictions = append(predictions, contracts.Prediction{
			RouteID:      record.RouteID,
			OfficeFromID: record.OfficeFromID,
			Timestamp:    record.Timestamp,
			Target2H:     volume,
		})
	}

	return contracts.PredictResponse{
		Predictions:  predictions,
		ModelVersion: "mock-ml-v1",
	}, nil
}
