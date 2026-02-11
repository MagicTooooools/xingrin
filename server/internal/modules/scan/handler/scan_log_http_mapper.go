package handler

import (
	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	"github.com/yyhuni/lunafox/server/internal/modules/scan/dto"
)

func toScanLogListQuery(query *dto.ScanLogListQuery) *scanapp.ScanLogListQuery {
	if query == nil {
		return nil
	}
	return &scanapp.ScanLogListQuery{
		AfterID: query.AfterID,
		Limit:   query.Limit,
	}
}

func toScanLogBulkCreateRequest(scanID int, req *dto.BulkCreateScanLogsRequest) *scanapp.ScanLogBulkCreateRequest {
	if req == nil {
		return &scanapp.ScanLogBulkCreateRequest{ScanID: scanID}
	}
	items := make([]scanapp.ScanLogCreateItem, 0, len(req.Logs))
	for index := range req.Logs {
		item := req.Logs[index]
		items = append(items, scanapp.ScanLogCreateItem{
			Level:   item.Level,
			Content: item.Content,
		})
	}
	return &scanapp.ScanLogBulkCreateRequest{
		ScanID: scanID,
		Items:  items,
	}
}
