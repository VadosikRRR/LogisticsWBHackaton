package dto

import (
	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type IngestRecordsRequest = contracts.RawRecordBatch
type GetWindowResponse = contracts.RawRecordBatch

type StatsResponse struct {
	TotalRecords int `json:"total_records"`
}

func ToEntityRecords(records []contracts.RawRecord) []entity.RawRecord {
	out := make([]entity.RawRecord, 0, len(records))
	for _, record := range records {
		out = append(out, entity.RawRecord{
			RouteID:      record.RouteID,
			OfficeFromID: record.OfficeFromID,
			Timestamp:    record.Timestamp,
			Status1:      record.Status1,
			Status2:      record.Status2,
			Status3:      record.Status3,
			Status4:      record.Status4,
			Status5:      record.Status5,
			Status6:      record.Status6,
			Status7:      record.Status7,
			Status8:      record.Status8,
			Target2H:     record.Target2H,
		})
	}
	return out
}

func FromEntityRecords(records []entity.RawRecord) []contracts.RawRecord {
	out := make([]contracts.RawRecord, 0, len(records))
	for _, record := range records {
		out = append(out, contracts.RawRecord{
			RouteID:      record.RouteID,
			OfficeFromID: record.OfficeFromID,
			Timestamp:    record.Timestamp,
			Status1:      record.Status1,
			Status2:      record.Status2,
			Status3:      record.Status3,
			Status4:      record.Status4,
			Status5:      record.Status5,
			Status6:      record.Status6,
			Status7:      record.Status7,
			Status8:      record.Status8,
			Target2H:     record.Target2H,
		})
	}
	return out
}
