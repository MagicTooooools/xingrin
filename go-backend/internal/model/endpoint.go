package model

import (
	"time"

	"github.com/lib/pq"
)

// Endpoint represents an endpoint asset
type Endpoint struct {
	ID                int            `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID          int            `gorm:"column:target_id;not null;index:idx_endpoint_target;uniqueIndex:unique_endpoint_url_target,priority:2" json:"targetId"`
	URL               string         `gorm:"column:url;type:text;index:idx_endpoint_url;uniqueIndex:unique_endpoint_url_target,priority:1" json:"url"`
	Host              string         `gorm:"column:host;size:253;index:idx_endpoint_host" json:"host"`
	Location          string         `gorm:"column:location;type:text" json:"location"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime;index:idx_endpoint_created_at" json:"createdAt"`
	Title             string         `gorm:"column:title;type:text;index:idx_endpoint_title" json:"title"`
	Webserver         string         `gorm:"column:webserver;type:text" json:"webserver"`
	ResponseBody      string         `gorm:"column:response_body;type:text" json:"responseBody"`
	ContentType       string         `gorm:"column:content_type;type:text" json:"contentType"`
	Tech              pq.StringArray `gorm:"column:tech;type:varchar(100)[]" json:"tech"`
	StatusCode        *int           `gorm:"column:status_code;index:idx_endpoint_status_code" json:"statusCode"`
	ContentLength     *int           `gorm:"column:content_length" json:"contentLength"`
	Vhost             *bool          `gorm:"column:vhost" json:"vhost"`
	MatchedGFPatterns pq.StringArray `gorm:"column:matched_gf_patterns;type:varchar(100)[]" json:"matchedGfPatterns"`
	ResponseHeaders   string         `gorm:"column:response_headers;type:text" json:"responseHeaders"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for Endpoint
func (Endpoint) TableName() string {
	return "endpoint"
}
