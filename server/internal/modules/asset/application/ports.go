package application

import (
	"database/sql"

	assetdomain "github.com/yyhuni/lunafox/server/internal/modules/asset/domain"
)

type AssetTargetLookup interface {
	GetActiveByID(id int) (*assetdomain.TargetRef, error)
}

type WebsiteTargetLookup = AssetTargetLookup

type WebsiteQueryStore interface {
	FindByTargetID(targetID int, page, pageSize int, filter string) ([]assetdomain.Website, int64, error)
	StreamByTargetID(targetID int) (*sql.Rows, error)
	CountByTargetID(targetID int) (int64, error)
	ScanRow(rows *sql.Rows) (*assetdomain.Website, error)
}

type WebsiteCommandStore interface {
	GetByID(id int) (*assetdomain.Website, error)
	BulkCreate(websites []assetdomain.Website) (int, error)
	Delete(id int) error
	BulkDelete(ids []int) (int64, error)
	BulkUpsert(websites []assetdomain.Website) (int64, error)
}

type WebsiteStore interface {
	WebsiteQueryStore
	WebsiteCommandStore
}

type EndpointTargetLookup = AssetTargetLookup

type EndpointQueryStore interface {
	FindByTargetID(targetID int, page, pageSize int, filter string) ([]assetdomain.Endpoint, int64, error)
	GetByID(id int) (*assetdomain.Endpoint, error)
	StreamByTargetID(targetID int) (*sql.Rows, error)
	CountByTargetID(targetID int) (int64, error)
	ScanRow(rows *sql.Rows) (*assetdomain.Endpoint, error)
}

type EndpointCommandStore interface {
	GetByID(id int) (*assetdomain.Endpoint, error)
	BulkCreate(endpoints []assetdomain.Endpoint) (int, error)
	Delete(id int) error
	BulkDelete(ids []int) (int64, error)
	BulkUpsert(endpoints []assetdomain.Endpoint) (int64, error)
}

type EndpointStore interface {
	EndpointQueryStore
	EndpointCommandStore
}

type DirectoryTargetLookup = AssetTargetLookup

type DirectoryQueryStore interface {
	FindByTargetID(targetID int, page, pageSize int, filter string) ([]assetdomain.Directory, int64, error)
	StreamByTargetID(targetID int) (*sql.Rows, error)
	CountByTargetID(targetID int) (int64, error)
	ScanRow(rows *sql.Rows) (*assetdomain.Directory, error)
}

type DirectoryCommandStore interface {
	BulkCreate(directories []assetdomain.Directory) (int, error)
	BulkDelete(ids []int) (int64, error)
	BulkUpsert(directories []assetdomain.Directory) (int64, error)
}

type DirectoryStore interface {
	DirectoryQueryStore
	DirectoryCommandStore
}

type SubdomainTargetLookup = AssetTargetLookup

type SubdomainQueryStore interface {
	FindByTargetID(targetID int, page, pageSize int, filter string) ([]assetdomain.Subdomain, int64, error)
	StreamByTargetID(targetID int) (*sql.Rows, error)
	CountByTargetID(targetID int) (int64, error)
	ScanRow(rows *sql.Rows) (*assetdomain.Subdomain, error)
}

type SubdomainCommandStore interface {
	BulkCreate(subdomains []assetdomain.Subdomain) (int, error)
	BulkDelete(ids []int) (int64, error)
}

type SubdomainStore interface {
	SubdomainQueryStore
	SubdomainCommandStore
}

type HostPortTargetLookup = AssetTargetLookup

type HostPortQueryStore interface {
	GetIPAggregation(targetID int, page, pageSize int, filter string) ([]assetdomain.IPAggregationRow, int64, error)
	GetHostsAndPortsByIP(targetID int, ip string, filter string) ([]string, []int, error)
	StreamByTargetID(targetID int) (*sql.Rows, error)
	StreamByTargetIDAndIPs(targetID int, ips []string) (*sql.Rows, error)
	CountByTargetID(targetID int) (int64, error)
	ScanRow(rows *sql.Rows) (*assetdomain.HostPort, error)
}

type HostPortCommandStore interface {
	BulkUpsert(mappings []assetdomain.HostPort) (int64, error)
	DeleteByIPs(ips []string) (int64, error)
}

type HostPortStore interface {
	HostPortQueryStore
	HostPortCommandStore
}

type ScreenshotTargetLookup = AssetTargetLookup

type ScreenshotQueryStore interface {
	FindByTargetID(targetID int, page, pageSize int, filter string) ([]assetdomain.Screenshot, int64, error)
	GetByID(id int) (*assetdomain.Screenshot, error)
}

type ScreenshotCommandStore interface {
	BulkDelete(ids []int) (int64, error)
	BulkUpsert(screenshots []assetdomain.Screenshot) (int64, error)
}

type ScreenshotStore interface {
	ScreenshotQueryStore
	ScreenshotCommandStore
}
