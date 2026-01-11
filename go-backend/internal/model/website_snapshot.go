package model

import (
	"time"

	"github.com/lib/pq"
)

// WebsiteSnapshot represents a website discovered in a scan
type WebsiteSnapshot struct {
	ID              int            `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID          int            `gorm:"column:scan_id;not null" json:"scanId"`
	URL             string         `gorm:"column:url;type:text" json:"url"`
	Host            string         `gorm:"column:host;size:253" json:"host"`
	Title           string         `gorm:"column:title;type:text" json:"title"`
	StatusCode      *int           `gorm:"column:status_code" json:"statusCode"`
	ContentLength   *int64         `gorm:"column:content_length" json:"contentLength"`
	Location        string         `gorm:"column:location;type:text" json:"location"`
	Webserver       string         `gorm:"column:webserver;type:text" json:"webserver"`
	ContentType     string         `gorm:"column:content_type;type:text" json:"contentType"`
	Tech            pq.StringArray `gorm:"column:tech;type:varchar(100)[]" json:"tech"`
	ResponseBody    string         `gorm:"column:response_body;type:text" json:"responseBody"`
	Vhost           *bool          `gorm:"column:vhost" json:"vhost"`
	ResponseHeaders string         `gorm:"column:response_headers;type:text" json:"responseHeaders"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for WebsiteSnapshot
func (WebsiteSnapshot) TableName() string {
	return "website_snapshot"
}
