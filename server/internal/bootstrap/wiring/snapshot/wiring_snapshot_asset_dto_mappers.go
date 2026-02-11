package snapshotwiring

import (
	assetdto "github.com/yyhuni/lunafox/server/internal/modules/asset/dto"
	snapshotapp "github.com/yyhuni/lunafox/server/internal/modules/snapshot/application"
)

func snapshotWebsiteAssetUpsertItemsToDTO(items []snapshotapp.WebsiteAssetUpsertItem) []assetdto.WebsiteUpsertItem {
	results := make([]assetdto.WebsiteUpsertItem, 0, len(items))
	for index := range items {
		item := items[index]
		results = append(results, assetdto.WebsiteUpsertItem{
			URL:             item.URL,
			Host:            item.Host,
			Location:        item.Location,
			Title:           item.Title,
			Webserver:       item.Webserver,
			ContentType:     item.ContentType,
			StatusCode:      item.StatusCode,
			ContentLength:   item.ContentLength,
			ResponseBody:    item.ResponseBody,
			Tech:            item.Tech,
			Vhost:           item.Vhost,
			ResponseHeaders: item.ResponseHeaders,
		})
	}
	return results
}

func snapshotEndpointAssetUpsertItemsToDTO(items []snapshotapp.EndpointAssetUpsertItem) []assetdto.EndpointUpsertItem {
	results := make([]assetdto.EndpointUpsertItem, 0, len(items))
	for index := range items {
		item := items[index]
		results = append(results, assetdto.EndpointUpsertItem{
			URL:             item.URL,
			Host:            item.Host,
			Location:        item.Location,
			Title:           item.Title,
			Webserver:       item.Webserver,
			ContentType:     item.ContentType,
			StatusCode:      item.StatusCode,
			ContentLength:   item.ContentLength,
			ResponseBody:    item.ResponseBody,
			Tech:            item.Tech,
			Vhost:           item.Vhost,
			ResponseHeaders: item.ResponseHeaders,
		})
	}
	return results
}

func snapshotDirectoryAssetUpsertItemsToDTO(items []snapshotapp.DirectoryAssetUpsertItem) []assetdto.DirectoryUpsertItem {
	results := make([]assetdto.DirectoryUpsertItem, 0, len(items))
	for index := range items {
		item := items[index]
		results = append(results, assetdto.DirectoryUpsertItem{
			URL:           item.URL,
			Status:        item.Status,
			ContentLength: item.ContentLength,
			ContentType:   item.ContentType,
			Duration:      item.Duration,
		})
	}
	return results
}

func snapshotHostPortAssetItemsToDTO(items []snapshotapp.HostPortAssetItem) []assetdto.HostPortItem {
	results := make([]assetdto.HostPortItem, 0, len(items))
	for index := range items {
		item := items[index]
		results = append(results, assetdto.HostPortItem{Host: item.Host, IP: item.IP, Port: item.Port})
	}
	return results
}

func snapshotScreenshotAssetRequestToDTO(req *snapshotapp.ScreenshotAssetUpsertRequest) *assetdto.BulkUpsertScreenshotRequest {
	if req == nil {
		return nil
	}
	items := make([]assetdto.ScreenshotItem, 0, len(req.Screenshots))
	for index := range req.Screenshots {
		item := req.Screenshots[index]
		items = append(items, assetdto.ScreenshotItem{URL: item.URL, StatusCode: item.StatusCode, Image: item.Image})
	}
	return &assetdto.BulkUpsertScreenshotRequest{Screenshots: items}
}
