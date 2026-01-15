package dto

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// BindJSON binds JSON request body and handles validation errors automatically.
// Returns true if binding succeeded, false if failed (response already sent).
func BindJSON(c *gin.Context, obj any) bool {
	// Validate Content-Type
	contentType := c.GetHeader("Content-Type")
	if contentType == "" || !strings.HasPrefix(contentType, "application/json") {
		Error(c, http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return false
	}

	if err := c.ShouldBindJSON(obj); err != nil {
		if HandleBindingError(c, err) {
			return false
		}
		BadRequest(c, "Invalid request body")
		return false
	}
	return true
}

// BindQuery binds query parameters and handles validation errors automatically.
// Returns true if binding succeeded, false if failed (response already sent).
func BindQuery(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		if HandleBindingError(c, err) {
			return false
		}
		BadRequest(c, "Invalid query parameters")
		return false
	}
	return true
}

// BindURI binds URI parameters and handles validation errors automatically.
// Returns true if binding succeeded, false if failed (response already sent).
func BindURI(c *gin.Context, obj any) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		if HandleBindingError(c, err) {
			return false
		}
		BadRequest(c, "Invalid URI parameters")
		return false
	}
	return true
}
