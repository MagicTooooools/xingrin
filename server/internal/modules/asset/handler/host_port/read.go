package hostport

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
	"github.com/yyhuni/lunafox/server/internal/modules/asset/dto"
)

// List returns paginated host-ports aggregated by IP.
// GET /api/targets/:id/host-ports
func (h *HostPortHandler) List(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var query dto.HostPortListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	results, total, err := h.svc.ListByTarget(targetID, &query)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to list host-ports")
		return
	}

	for i := range results {
		if results[i].Hosts == nil {
			results[i].Hosts = []string{}
		}
		if results[i].Ports == nil {
			results[i].Ports = []int{}
		}
	}

	dto.Paginated(c, results, total, query.GetPage(), query.GetPageSize())
}
