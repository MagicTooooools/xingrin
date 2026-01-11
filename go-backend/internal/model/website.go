package model

import (
	"time"

	"github.com/lib/pq"
)

// Website represents a website asset
type Website struct {
	ID              int            `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID        int            `gorm:"column:target_id;not null;index:idx_website_target;uniqueIndex:unique_website_url_target,priority:2" json:"targetId"`
	URL             string         `gorm:"column:url;type:text;index:idx_website_url;uniqueIndex:unique_website_url_target,priority:1" json:"url"`
	Host            string         `gorm:"column:host;size:253;index:idx_website_host" json:"host"`
	Location        string         `gorm:"column:location;type:text" json:"location"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime;index:idx_website_created_at" json:"createdAt"`
	Title           string         `gorm:"column:title;type:text;index:idx_website_title" json:"title"`
	Webserver       string         `gorm:"column:webserver;type:text" json:"webserver"`
	ResponseBody    string         `gorm:"column:response_body;type:text" json:"responseBody"`
	ContentType     string         `gorm:"column:content_type;type:text" json:"contentType"`
	Tech            pq.StringArray `gorm:"column:tech;type:varchar(100)[]" json:"tech"`
	StatusCode      *int           `gorm:"column:status_code;index:idx_website_status_code" json:"statusCode"`
	ContentLength   *int           `gorm:"column:content_length" json:"contentLength"`
	Vhost           *bool          `gorm:"column:vhost" json:"vhost"`
	ResponseHeaders string         `gorm:"column:response_headers;type:text" json:"responseHeaders"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for Website
func (Website) TableName() string {
	return "website"
}
