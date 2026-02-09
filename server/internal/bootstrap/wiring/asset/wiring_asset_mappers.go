package assetwiring

import (
	assetdomain "github.com/yyhuni/lunafox/server/internal/modules/asset/domain"
	assetmodel "github.com/yyhuni/lunafox/server/internal/modules/asset/repository/persistence"
)

func assetModelWebsiteToDomain(item *assetmodel.Website) *assetdomain.Website {
	if item == nil {
		return nil
	}
	return &assetdomain.Website{
		ID:              item.ID,
		TargetID:        item.TargetID,
		URL:             item.URL,
		Host:            item.Host,
		Location:        item.Location,
		CreatedAt:       item.CreatedAt,
		Title:           item.Title,
		Webserver:       item.Webserver,
		ResponseBody:    item.ResponseBody,
		ContentType:     item.ContentType,
		Tech:            item.Tech,
		StatusCode:      item.StatusCode,
		ContentLength:   item.ContentLength,
		Vhost:           item.Vhost,
		ResponseHeaders: item.ResponseHeaders,
	}
}

func assetDomainWebsiteToModel(item *assetdomain.Website) *assetmodel.Website {
	if item == nil {
		return nil
	}
	return &assetmodel.Website{
		ID:              item.ID,
		TargetID:        item.TargetID,
		URL:             item.URL,
		Host:            item.Host,
		Location:        item.Location,
		CreatedAt:       item.CreatedAt,
		Title:           item.Title,
		Webserver:       item.Webserver,
		ResponseBody:    item.ResponseBody,
		ContentType:     item.ContentType,
		Tech:            item.Tech,
		StatusCode:      item.StatusCode,
		ContentLength:   item.ContentLength,
		Vhost:           item.Vhost,
		ResponseHeaders: item.ResponseHeaders,
	}
}

func assetDomainWebsiteListToModel(items []assetdomain.Website) []assetmodel.Website {
	results := make([]assetmodel.Website, 0, len(items))
	for index := range items {
		results = append(results, *assetDomainWebsiteToModel(&items[index]))
	}
	return results
}

func assetModelEndpointToDomain(item *assetmodel.Endpoint) *assetdomain.Endpoint {
	if item == nil {
		return nil
	}
	return &assetdomain.Endpoint{
		ID:                item.ID,
		TargetID:          item.TargetID,
		URL:               item.URL,
		Host:              item.Host,
		Location:          item.Location,
		CreatedAt:         item.CreatedAt,
		Title:             item.Title,
		Webserver:         item.Webserver,
		ResponseBody:      item.ResponseBody,
		ContentType:       item.ContentType,
		Tech:              item.Tech,
		StatusCode:        item.StatusCode,
		ContentLength:     item.ContentLength,
		Vhost:             item.Vhost,
		MatchedGFPatterns: item.MatchedGFPatterns,
		ResponseHeaders:   item.ResponseHeaders,
	}
}

func assetDomainEndpointToModel(item *assetdomain.Endpoint) *assetmodel.Endpoint {
	if item == nil {
		return nil
	}
	return &assetmodel.Endpoint{
		ID:                item.ID,
		TargetID:          item.TargetID,
		URL:               item.URL,
		Host:              item.Host,
		Location:          item.Location,
		CreatedAt:         item.CreatedAt,
		Title:             item.Title,
		Webserver:         item.Webserver,
		ResponseBody:      item.ResponseBody,
		ContentType:       item.ContentType,
		Tech:              item.Tech,
		StatusCode:        item.StatusCode,
		ContentLength:     item.ContentLength,
		Vhost:             item.Vhost,
		MatchedGFPatterns: item.MatchedGFPatterns,
		ResponseHeaders:   item.ResponseHeaders,
	}
}

func assetDomainEndpointListToModel(items []assetdomain.Endpoint) []assetmodel.Endpoint {
	results := make([]assetmodel.Endpoint, 0, len(items))
	for index := range items {
		results = append(results, *assetDomainEndpointToModel(&items[index]))
	}
	return results
}

func assetModelDirectoryToDomain(item *assetmodel.Directory) *assetdomain.Directory {
	if item == nil {
		return nil
	}
	return &assetdomain.Directory{
		ID:            item.ID,
		TargetID:      item.TargetID,
		URL:           item.URL,
		Status:        item.Status,
		ContentLength: item.ContentLength,
		ContentType:   item.ContentType,
		Duration:      item.Duration,
		CreatedAt:     item.CreatedAt,
	}
}

func assetDomainDirectoryToModel(item *assetdomain.Directory) *assetmodel.Directory {
	if item == nil {
		return nil
	}
	return &assetmodel.Directory{
		ID:            item.ID,
		TargetID:      item.TargetID,
		URL:           item.URL,
		Status:        item.Status,
		ContentLength: item.ContentLength,
		ContentType:   item.ContentType,
		Duration:      item.Duration,
		CreatedAt:     item.CreatedAt,
	}
}

func assetDomainDirectoryListToModel(items []assetdomain.Directory) []assetmodel.Directory {
	results := make([]assetmodel.Directory, 0, len(items))
	for index := range items {
		results = append(results, *assetDomainDirectoryToModel(&items[index]))
	}
	return results
}

func assetModelSubdomainToDomain(item *assetmodel.Subdomain) *assetdomain.Subdomain {
	if item == nil {
		return nil
	}
	return &assetdomain.Subdomain{
		ID:        item.ID,
		TargetID:  item.TargetID,
		Name:      item.Name,
		CreatedAt: item.CreatedAt,
	}
}

func assetDomainSubdomainToModel(item *assetdomain.Subdomain) *assetmodel.Subdomain {
	if item == nil {
		return nil
	}
	return &assetmodel.Subdomain{
		ID:        item.ID,
		TargetID:  item.TargetID,
		Name:      item.Name,
		CreatedAt: item.CreatedAt,
	}
}

func assetDomainSubdomainListToModel(items []assetdomain.Subdomain) []assetmodel.Subdomain {
	results := make([]assetmodel.Subdomain, 0, len(items))
	for index := range items {
		results = append(results, *assetDomainSubdomainToModel(&items[index]))
	}
	return results
}

func assetModelHostPortToDomain(item *assetmodel.HostPort) *assetdomain.HostPort {
	if item == nil {
		return nil
	}
	return &assetdomain.HostPort{
		ID:        item.ID,
		TargetID:  item.TargetID,
		Host:      item.Host,
		IP:        item.IP,
		Port:      item.Port,
		CreatedAt: item.CreatedAt,
	}
}

func assetDomainHostPortToModel(item *assetdomain.HostPort) *assetmodel.HostPort {
	if item == nil {
		return nil
	}
	return &assetmodel.HostPort{
		ID:        item.ID,
		TargetID:  item.TargetID,
		Host:      item.Host,
		IP:        item.IP,
		Port:      item.Port,
		CreatedAt: item.CreatedAt,
	}
}

func assetDomainHostPortListToModel(items []assetdomain.HostPort) []assetmodel.HostPort {
	results := make([]assetmodel.HostPort, 0, len(items))
	for index := range items {
		results = append(results, *assetDomainHostPortToModel(&items[index]))
	}
	return results
}

func assetModelScreenshotToDomain(item *assetmodel.Screenshot) *assetdomain.Screenshot {
	if item == nil {
		return nil
	}
	return &assetdomain.Screenshot{
		ID:         item.ID,
		TargetID:   item.TargetID,
		URL:        item.URL,
		StatusCode: item.StatusCode,
		Image:      item.Image,
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
	}
}

func assetDomainScreenshotToModel(item *assetdomain.Screenshot) *assetmodel.Screenshot {
	if item == nil {
		return nil
	}
	return &assetmodel.Screenshot{
		ID:         item.ID,
		TargetID:   item.TargetID,
		URL:        item.URL,
		StatusCode: item.StatusCode,
		Image:      item.Image,
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
	}
}

func assetDomainScreenshotListToModel(items []assetdomain.Screenshot) []assetmodel.Screenshot {
	results := make([]assetmodel.Screenshot, 0, len(items))
	for index := range items {
		results = append(results, *assetDomainScreenshotToModel(&items[index]))
	}
	return results
}
