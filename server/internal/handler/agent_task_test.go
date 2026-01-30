package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yyhuni/lunafox/server/internal/dto"
	"github.com/yyhuni/lunafox/server/internal/service"
)

type fakeTaskService struct {
	task         *dto.TaskAssignment
	err          error
	status       string
	errorMessage string
	taskID       int
	agentID      int
}

func (f *fakeTaskService) PullTask(ctx context.Context, agentID int) (*dto.TaskAssignment, error) {
	f.agentID = agentID
	return f.task, f.err
}

func (f *fakeTaskService) UpdateStatus(ctx context.Context, agentID, taskID int, status, errorMessage string) error {
	f.agentID = agentID
	f.taskID = taskID
	f.status = status
	f.errorMessage = errorMessage
	return f.err
}

func TestAgentTaskPullNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/agent/tasks/pull", nil)
	c.Set("agentID", 1)

	handler := NewAgentTaskHandler(&fakeTaskService{})
	handler.PullTask(c)

	if c.Writer.Status() != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", c.Writer.Status())
	}
}

func TestAgentTaskPullOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/agent/tasks/pull", nil)
	c.Set("agentID", 1)

	handler := NewAgentTaskHandler(&fakeTaskService{task: &dto.TaskAssignment{TaskID: 1}})
	handler.PullTask(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAgentTaskPullUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/agent/tasks/pull", nil)

	handler := NewAgentTaskHandler(&fakeTaskService{})
	handler.PullTask(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAgentTaskPullInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/agent/tasks/pull", nil)
	c.Set("agentID", 1)

	handler := NewAgentTaskHandler(&fakeTaskService{err: assertError("boom")})
	handler.PullTask(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestAgentTaskUpdateStatusBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/api/agent/tasks/abc/status", nil)
	c.Set("agentID", 1)

	handler := NewAgentTaskHandler(&fakeTaskService{})
	handler.UpdateTaskStatus(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAgentTaskUpdateStatusUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/api/agent/tasks/1/status", nil)
	c.Params = gin.Params{{Key: "taskId", Value: "1"}}

	handler := NewAgentTaskHandler(&fakeTaskService{})
	handler.UpdateTaskStatus(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAgentTaskUpdateStatusFromJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := strings.NewReader(`{"status":"failed","errorMessage":"boom"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/agent/tasks/7/status", body)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{{Key: "taskId", Value: "7"}}
	c.Set("agentID", 42)

	svc := &fakeTaskService{}
	handler := NewAgentTaskHandler(svc)
	handler.UpdateTaskStatus(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if svc.status != "failed" || svc.errorMessage != "boom" {
		t.Fatalf("unexpected payload for update")
	}
	if svc.taskID != 7 || svc.agentID != 42 {
		t.Fatalf("unexpected ids for update")
	}
}

func TestAgentTaskUpdateStatusBindingError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := strings.NewReader(`{}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/agent/tasks/2/status", body)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{{Key: "taskId", Value: "2"}}
	c.Set("agentID", 1)

	handler := NewAgentTaskHandler(&fakeTaskService{})
	handler.UpdateTaskStatus(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAgentTaskUpdateStatusErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/api/agent/tasks/1/status", nil)
	c.Set("agentID", 1)
	c.Params = gin.Params{{Key: "taskId", Value: "1"}}
	c.Set("validatedStatus", "completed")
	c.Set("validatedErrorMessage", "")

	handler := NewAgentTaskHandler(&fakeTaskService{err: service.ErrScanTaskNotFound})
	handler.UpdateTaskStatus(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestAgentTaskUpdateStatusForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/api/agent/tasks/1/status", nil)
	c.Set("agentID", 1)
	c.Params = gin.Params{{Key: "taskId", Value: "1"}}
	c.Set("validatedStatus", "failed")

	handler := NewAgentTaskHandler(&fakeTaskService{err: service.ErrScanTaskNotOwned})
	handler.UpdateTaskStatus(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestAgentTaskUpdateStatusInvalidTransition(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/api/agent/tasks/1/status", nil)
	c.Set("agentID", 1)
	c.Params = gin.Params{{Key: "taskId", Value: "1"}}
	c.Set("validatedStatus", "failed")

	handler := NewAgentTaskHandler(&fakeTaskService{err: service.ErrScanTaskInvalidTransition})
	handler.UpdateTaskStatus(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAgentTaskUpdateStatusInvalidUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/api/agent/tasks/1/status", nil)
	c.Set("agentID", 1)
	c.Params = gin.Params{{Key: "taskId", Value: "1"}}
	c.Set("validatedStatus", "failed")

	handler := NewAgentTaskHandler(&fakeTaskService{err: service.ErrScanTaskInvalidUpdate})
	handler.UpdateTaskStatus(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAgentTaskUpdateStatusInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/api/agent/tasks/1/status", nil)
	c.Set("agentID", 1)
	c.Params = gin.Params{{Key: "taskId", Value: "1"}}
	c.Set("validatedStatus", "failed")

	handler := NewAgentTaskHandler(&fakeTaskService{err: assertError("boom")})
	handler.UpdateTaskStatus(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

type assertError string

func (e assertError) Error() string {
	return string(e)
}
