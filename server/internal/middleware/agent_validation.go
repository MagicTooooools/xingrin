package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	scandto "github.com/yyhuni/lunafox/server/internal/modules/scan/dto"
)

// AgentValidationMiddleware validates input for Agent API endpoints
func AgentValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.Contains(path, "/tasks/") && c.Request.Method == "PATCH" {
			var req scandto.TaskStatusUpdateRequest
			body, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				c.Abort()
				return
			}
			if err := json.Unmarshal(body, &req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				c.Abort()
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			validStatuses := map[string]bool{
				"completed": true,
				"failed":    true,
				"cancelled": true,
			}

			if !validStatuses[req.Status] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
				c.Abort()
				return
			}

			if len(req.ErrorMessage) > 4096 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error message exceeds 4KB limit"})
				c.Abort()
				return
			}

			c.Set("validatedStatus", req.Status)
			c.Set("validatedErrorMessage", req.ErrorMessage)
		}

		c.Next()
	}
}
