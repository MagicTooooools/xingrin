package model

import (
	"time"
)

// Session represents a Django django_session compatible model
type Session struct {
	SessionKey  string    `gorm:"column:session_key;size:40;primaryKey" json:"sessionKey"`
	SessionData string    `gorm:"column:session_data;type:text" json:"sessionData"`
	ExpireDate  time.Time `gorm:"column:expire_date;index" json:"expireDate"`
}

// TableName returns the table name for Session (Django django_session)
func (Session) TableName() string {
	return "django_session"
}
