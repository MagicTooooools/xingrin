package dto

import "github.com/gin-gonic/gin"

// BindJSON binds JSON request body and handles validation errors automatically.
// Returns true if binding succeeded, false if failed (response already sent).
func BindJSON(c *gin.Context, obj any) bool {
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
