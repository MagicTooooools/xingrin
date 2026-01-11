package model

import (
	"time"

	"github.com/lib/pq"
)

// Endpoint represents a discovered endpoint
type Endpoint struct {
	ID                int            `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID          int            `gorm:"column:target_id;not null" json:"targetId"`
	URL               string         `gorm:"column:url;type:text" json:"url"`
	Host              string         `gorm:"column:host;size:253" json:"host"`
	Location          string         `gorm:"column:location;type:text" json:"location"`
	Title             string         `gorm:"column:title;type:text" json:"title"`
	Webserver         string         `gorm:"column:webserver;type:text" json:"webserver"`
	ResponseBody      string         `gorm:"column:response_body;type:text" json:"responseBody"`
	ContentType       string         `gorm:"column:content_type;type:text" json:"contentType"`
	Tech              pq.StringArray `gorm:"column:tech;type:varchar(100)[]" json:"tech"`
	StatusCode        *int           `gorm:"column:status_code" json:"statusCode"`
	ContentLength     *int           `gorm:"column:content_length" json:"contentLength"`
	Vhost             *bool          `gorm:"column:vhost" json:"vhost"`
	MatchedGFPatterns pq.StringArray `gorm:"column:matched_gf_patterns;type:varchar(100)[]" json:"matchedGfPatterns"`
	ResponseHeaders   string         `gorm:"column:response_headers;type:text" json:"responseHeaders"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for Endpoint
func (Endpoint) TableName() string {
	return "endpoint"
}
