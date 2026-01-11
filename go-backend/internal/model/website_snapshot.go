package model

import (
	"time"

	"github.com/lib/pq"
)

// WebsiteSnapshot represents a website snapshot
type WebsiteSnapshot struct {
	ID              int            `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID          int            `gorm:"column:scan_id;not null;index:idx_website_snap_scan;uniqueIndex:unique_website_per_scan_snapshot,priority:1" json:"scanId"`
	URL             string         `gorm:"column:url;type:text;index:idx_website_snap_url;uniqueIndex:unique_website_per_scan_snapshot,priority:2" json:"url"`
	Host            string         `gorm:"column:host;size:253;index:idx_website_snap_host" json:"host"`
	Title           string         `gorm:"column:title;type:text;index:idx_website_snap_title" json:"title"`
	StatusCode      *int           `gorm:"column:status_code" json:"statusCode"`
	ContentLength   *int64         `gorm:"column:content_length" json:"contentLength"`
	Location        string         `gorm:"column:location;type:text" json:"location"`
	Webserver       string         `gorm:"column:webserver;type:text" json:"webserver"`
	ContentType     string         `gorm:"column:content_type;type:text" json:"contentType"`
	Tech            pq.StringArray `gorm:"column:tech;type:varchar(100)[]" json:"tech"`
	ResponseBody    string         `gorm:"column:response_body;type:text" json:"responseBody"`
	Vhost           *bool          `gorm:"column:vhost" json:"vhost"`
	ResponseHeaders string         `gorm:"column:response_headers;type:text" json:"responseHeaders"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime;index:idx_website_snap_created_at" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for WebsiteSnapshot
func (WebsiteSnapshot) TableName() string {
	return "website_snapshot"
}
