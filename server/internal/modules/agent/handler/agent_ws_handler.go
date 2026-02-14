package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	agentapp "github.com/yyhuni/lunafox/server/internal/modules/agent/application"
	agentdomain "github.com/yyhuni/lunafox/server/internal/modules/agent/domain"
	"github.com/yyhuni/lunafox/server/internal/modules/agent/dto"
	"github.com/yyhuni/lunafox/server/internal/pkg"
	ws "github.com/yyhuni/lunafox/server/internal/websocket"
	"go.uber.org/zap"
)

// AgentWebSocketHandler handles WebSocket connections from agents.
type AgentWebSocketHandler struct {
	hub            *ws.Hub
	runtimeService *agentapp.AgentRuntimeService
	upgrader       websocket.Upgrader
}

// NewAgentWebSocketHandler creates a new AgentWebSocketHandler.
func NewAgentWebSocketHandler(hub *ws.Hub, runtimeService *agentapp.AgentRuntimeService) *AgentWebSocketHandler {
	return &AgentWebSocketHandler{
		hub:            hub,
		runtimeService: runtimeService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Handle upgrades HTTP to WebSocket and starts the client loops.
// GET /api/agent/ws
func (h *AgentWebSocketHandler) Handle(c *gin.Context) {
	agent, ok := contextAgent(c)
	if !ok || agent == nil {
		dto.Unauthorized(c, "Invalid agent context")
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		pkg.Error("Failed to upgrade WebSocket", zap.Error(err))
		return
	}

	client := &ws.Client{AgentID: agent.ID, Conn: conn, Send: make(chan []byte, 256), Hub: h.hub}
	h.hub.Register(client)

	if err := h.runtimeService.OnConnected(c.Request.Context(), agent, getForwardedIP(c)); err != nil {
		pkg.Warn("Failed to update agent connection state", zap.Error(err))
	}
	pkg.Info("Agent websocket connected", zap.Int("agent_id", agent.ID), zap.String("name", agent.Name), zap.String("ip", agent.IPAddress))

	go h.writePump(client)
	h.readPump(c.Request.Context(), client, agent)
}

func (h *AgentWebSocketHandler) readPump(ctx context.Context, client *ws.Client, agent *agentdomain.Agent) {
	defer func() {
		h.hub.Unregister(client)
		_ = client.Conn.Close()
		if err := h.runtimeService.OnDisconnected(context.Background(), agent.ID); err != nil {
			pkg.Warn("Failed to mark agent offline on disconnect", zap.Error(err), zap.Int("agent_id", agent.ID))
		}
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				pkg.Warn("Agent websocket closed unexpectedly", zap.Int("agent_id", agent.ID), zap.Error(err))
			} else {
				pkg.Info("Agent websocket disconnected", zap.Int("agent_id", agent.ID), zap.Error(err))
			}
			return
		}
		if err := h.runtimeService.HandleMessage(ctx, agent, message); err != nil {
			pkg.Warn("Failed to handle WebSocket message", zap.Error(err), zap.Int("agent_id", agent.ID))
		}
	}
}

func (h *AgentWebSocketHandler) writePump(client *ws.Client) {
	defer func() {
		_ = client.Conn.Close()
	}()
	for msg := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func contextAgent(c *gin.Context) (*agentdomain.Agent, bool) {
	agentVal, ok := c.Get("agent")
	if !ok {
		return nil, false
	}
	agent, ok := agentVal.(*agentdomain.Agent)
	if !ok || agent == nil {
		return nil, false
	}
	return agent, true
}
