package model

import (
	"time"

	"github.com/lib/pq"
)

// EndpointSnapshot represents an endpoint snapshot
type EndpointSnapshot struct {
	ID                int            `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID            int            `gorm:"column:scan_id;not null;index:idx_endpoint_snap_scan;uniqueIndex:unique_endpoint_per_scan_snapshot,priority:1" json:"scanId"`
	URL               string         `gorm:"column:url;type:text;index:idx_endpoint_snap_url;uniqueIndex:unique_endpoint_per_scan_snapshot,priority:2" json:"url"`
	Host              string         `gorm:"column:host;size:253;index:idx_endpoint_snap_host" json:"host"`
	Title             string         `gorm:"column:title;type:text;index:idx_endpoint_snap_title" json:"title"`
	StatusCode        *int           `gorm:"column:status_code;index:idx_endpoint_snap_status_code" json:"statusCode"`
	ContentLength     *int           `gorm:"column:content_length" json:"contentLength"`
	Location          string         `gorm:"column:location;type:text" json:"location"`
	Webserver         string         `gorm:"column:webserver;type:text;index:idx_endpoint_snap_webserver" json:"webserver"`
	ContentType       string         `gorm:"column:content_type;type:text" json:"contentType"`
	Tech              pq.StringArray `gorm:"column:tech;type:varchar(100)[]" json:"tech"`
	ResponseBody      string         `gorm:"column:response_body;type:text" json:"responseBody"`
	Vhost             *bool          `gorm:"column:vhost" json:"vhost"`
	MatchedGFPatterns pq.StringArray `gorm:"column:matched_gf_patterns;type:varchar(100)[]" json:"matchedGfPatterns"`
	ResponseHeaders   string         `gorm:"column:response_headers;type:text" json:"responseHeaders"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime;index:idx_endpoint_snap_created_at" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for EndpointSnapshot
func (EndpointSnapshot) TableName() string {
	return "endpoint_snapshot"
}
