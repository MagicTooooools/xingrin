# Implementation Tasks: WebSocket Agent System

**Feature Branch**: `001-websocket-agent`
**Generated**: 2026-01-22
**Total Tasks**: 87

## Overview

This document breaks down the WebSocket Agent System implementation into executable tasks organized by user story. Each phase represents a complete, independently testable increment of functionality.

## Implementation Strategy

**MVP Scope**: Phase 3 (User Story 1 - Agent Deployment and Connection)
- Delivers core value: Agent can connect to Server and maintain connection
- Enables early testing and validation
- Foundation for all subsequent features

**Incremental Delivery**:
1. Setup + Foundational → US1 (MVP) → US2 → US3 → US4 → US5 → Polish
2. Each user story phase is independently testable
3. Parallel execution opportunities marked with [P]

## Phase 1: Setup (Project Initialization)

**Goal**: Initialize project structure and dependencies for both Agent and Server components.

### Tasks

- [ ] T001 Create agent/ module directory structure in /Users/yangyang/Desktop/orbit/agent/
- [ ] T002 Initialize Go module for agent in agent/go.mod with module name github.com/yyhuni/orbit/agent
- [ ] T003 Create agent subdirectories: cmd/agent/, internal/config/, internal/websocket/, internal/task/, internal/docker/, internal/metrics/
- [ ] T004 [P] Add gorilla/websocket dependency to agent/go.mod (go get github.com/gorilla/websocket)
- [ ] T005 [P] Add docker/docker SDK dependency to agent/go.mod (go get github.com/docker/docker)
- [ ] T006 [P] Add gopsutil v3 dependency to agent/go.mod (go get github.com/shirou/gopsutil/v3)
- [ ] T007 Create server extensions directory structure in server/internal/handler/agent.go
- [ ] T008 [P] Add gin-gonic/gin dependency to server/go.mod if not present
- [ ] T009 [P] Add gorilla/websocket dependency to server/go.mod if not present

## Phase 2: Foundational (Blocking Prerequisites)

**Goal**: Implement database schema, base models, and core infrastructure needed by all user stories.

**Independent Test**: Database migrations run successfully, base models can be instantiated, WebSocket infrastructure accepts connections.

### Tasks

- [ ] T010 Create database migration for agent table in server/migrations/YYYYMMDD_HHMMSS_create_agent_table.sql (use timestamp format: 20260122_143000)
- [ ] T011 Create database migration for scan_task table in server/migrations/YYYYMMDD_HHMMSS_create_scan_task_table.sql (includes version field, no worker_image field)
- [ ] T012 Create database indexes migration in server/migrations/YYYYMMDD_HHMMSS_create_indexes.sql
- [ ] T013 [P] Implement Agent model in server/internal/model/agent.go with GORM tags
- [ ] T014 [P] Implement ScanTask model in server/internal/model/scan_task.go with GORM tags (includes version field, no worker_image field)
- [ ] T015 Implement Agent repository interface in server/internal/repository/agent.go
- [ ] T016 Implement ScanTask repository interface in server/internal/repository/scan_task.go
- [ ] T017 [P] Implement Redis client wrapper in server/internal/cache/redis.go
- [ ] T018 [P] Implement heartbeat cache operations (set/get/delete) in server/internal/cache/heartbeat.go
- [ ] T019 Create WebSocket hub for connection management in server/internal/websocket/hub.go
- [ ] T020 Implement WebSocket authentication middleware in server/internal/middleware/agent_auth.go

## Phase 3: User Story 1 - Agent Deployment and Connection (P1) 🎯 MVP

**Story Goal**: As an operations person, I can quickly deploy Agent on remote VPS and connect to Server to start receiving tasks.

**Independent Test**:
1. Create Agent via Web UI → receive API key
2. Run installation command on remote machine → Agent shows online status
3. Disconnect network → Agent automatically reconnects within 120s
4. Provide wrong API key → connection rejected with auth error

### Tasks

- [ ] T021 [US1] Implement Config struct in agent/internal/config/config.go with ServerURL, APIKey, MaxTasks, thresholds
- [ ] T022 [US1] Implement config loading from environment variables in agent/internal/config/loader.go
- [ ] T023 [US1] Implement WebSocket client with exponential backoff in agent/internal/websocket/client.go
- [ ] T024 [US1] Implement connection authentication in agent/internal/websocket/auth.go (support header and query param)
- [ ] T025 [US1] Implement reconnection logic with backoff strategy (1s, 2s, 4s, 8s, max 60s) in agent/internal/websocket/reconnect.go
- [ ] T026 [US1] Implement main Agent entry point in agent/cmd/agent/main.go
- [ ] T027 [P] [US1] Implement Server WebSocket endpoint handler in server/internal/handler/agent_ws.go at /api/agents/ws
- [ ] T028 [P] [US1] Implement Agent registration on first connection in server/internal/service/agent_service.go
- [ ] T029 [P] [US1] Implement Agent status update to online in server/internal/repository/agent.go
- [ ] T030 [US1] Implement WebSocket message router in server/internal/websocket/router.go
- [ ] T031 [US1] Create Agent creation API endpoint in server/internal/handler/agent.go POST /api/agents
- [ ] T032 [US1] Implement API key generation (8 char hex string, 4 bytes random) in server/internal/service/agent_service.go
- [ ] T033 [US1] Implement Agent list API endpoint in server/internal/handler/agent.go GET /api/agents

## Phase 4: User Story 2 - Task Execution and Status Tracking (P1)

**Story Goal**: As a system admin, I want Agent to automatically pull tasks and execute scans, with real-time status updates for monitoring progress.

**Independent Test**:
1. Create scan task → Agent pulls and starts Worker container
2. Worker completes (exit code 0) → task status updates to completed
3. Worker fails (exit code ≠ 0) → task status updates to failed with error log
4. Cancel task in Web UI → Agent stops Worker and updates status to cancelled

### Tasks

- [ ] T034 [US2] Implement task pull HTTP client in agent/internal/task/client.go for POST /api/agent/tasks/pull
- [ ] T035 [US2] Implement task status update client in agent/internal/task/client.go for PATCH /api/agent/tasks/{taskId}/status
- [ ] T036 [US2] Implement Docker client wrapper in agent/internal/docker/client.go
- [ ] T037 [US2] Implement Worker container launcher in agent/internal/docker/runner.go (constructs image name as yyhuni/orbit-worker:v{version} from task version field, passes environment variables)
- [ ] T038 [US2] Implement container exit code monitoring in agent/internal/docker/monitor.go
- [ ] T039 [US2] Implement container log reader (last 100 lines, 4KB truncation) in agent/internal/docker/logs.go
- [ ] T040 [US2] Implement container cleanup logic in agent/internal/docker/cleanup.go (manual, not --rm)
- [ ] T041 [US2] Implement task executor orchestration in agent/internal/task/executor.go
- [ ] T042 [US2] Implement Worker timeout mechanism (default 7 days) in agent/internal/task/timeout.go
- [ ] T043 [P] [US2] Implement task pull endpoint in server/internal/handler/agent_task.go POST /api/agent/tasks/pull
- [ ] T044 [P] [US2] Implement task assignment with FOR UPDATE SKIP LOCKED in server/internal/repository/scan_task.go
- [ ] T045 [P] [US2] Implement task status update endpoint in server/internal/handler/agent_task.go PATCH /api/agent/tasks/{taskId}/status
- [ ] T046 [P] [US2] Implement status update validation in server/internal/service/scan_task_service.go (FR-022: ownership check - agent_id must match; FR-023: idempotency - duplicate status returns 200; FR-024: state transition validation - only allow pending→running, running→completed/failed/cancelled)
- [ ] T047 [P] [US2] Implement scan.status synchronization with scan_task.status in server/internal/service/scan_service.go
- [ ] T048 [US2] Implement task_available WebSocket notification in server/internal/websocket/notifier.go
- [ ] T049 [US2] Implement task_cancel WebSocket message handler in agent/internal/websocket/handlers.go
- [ ] T050 [US2] Implement scan_task creation on scan creation in server/internal/service/scan_service.go (reads VERSION file and sets version field)
- [ ] T051 [US2] Implement Agent pull strategy with backoff (5s/10s/30s, max 60s) and dynamic pull interval based on load (<50%: 1s, 50-80%: 3s, >80%: 10s) in agent/internal/task/puller.go

## Phase 5: User Story 3 - Load Monitoring and Smart Scheduling (P2)

**Story Goal**: As a system admin, I want Agent to monitor its load and intelligently decide whether to accept new tasks to avoid system overload.

**Independent Test**:
1. Agent CPU exceeds 85% → Agent waits before pulling new tasks
2. Agent reaches max concurrent tasks (5) → Agent waits until task completes
3. Agent sends heartbeat → Server records CPU, memory, disk, task count
4. Agent heartbeat timeout (>120s) → Server marks Agent offline

### Tasks

- [ ] T052 [US3] Implement system metrics collector using gopsutil in agent/internal/metrics/collector.go
- [ ] T053 [US3] Implement cgroup-aware metrics for containerized Agent in agent/internal/metrics/cgroup.go
- [ ] T054 [US3] Implement heartbeat message builder in agent/internal/websocket/heartbeat.go
- [ ] T055 [US3] Implement heartbeat sender (every 5 seconds) in agent/internal/websocket/sender.go
- [ ] T056 [US3] Implement load check before task pull in agent/internal/task/scheduler.go
- [ ] T057 [US3] Implement concurrent task counter in agent/internal/task/counter.go
- [ ] T058 [P] [US3] Implement heartbeat WebSocket message handler in server/internal/websocket/handlers.go
- [ ] T059 [P] [US3] Integrate heartbeat handler with cache layer (call T018 operations) in server/internal/websocket/handlers.go
- [ ] T060 [P] [US3] Implement Agent offline detection background job in server/internal/job/agent_monitor.go (runs every minute)
- [ ] T061 [P] [US3] Implement task recovery for offline Agents in server/internal/job/task_recovery.go
- [ ] T062 [US3] Implement heartbeat API endpoint for Web UI in server/internal/handler/agent.go GET /api/agents/{id}/heartbeat

## Phase 6: User Story 4 - Dynamic Configuration Updates (P2)

**Story Goal**: As a system admin, I can dynamically adjust Agent configuration parameters in Web UI to optimize performance without restarting Agent.

**Independent Test**:
1. Modify max tasks in Web UI → Agent receives and applies new config immediately
2. Update load thresholds → Agent uses new thresholds for load checks

### Tasks

- [ ] T063 [US4] Implement config_update WebSocket message handler in agent/internal/websocket/handlers.go
- [ ] T064 [US4] Implement dynamic config update in agent/internal/config/updater.go
- [ ] T065 [P] [US4] Implement Agent config update API endpoint in server/internal/handler/agent.go PATCH /api/agents/{id}/config
- [ ] T066 [P] [US4] Implement config_update WebSocket message sender in server/internal/websocket/notifier.go

## Phase 7: User Story 5 - Auto-Update (P3)

**Story Goal**: As an operations person, I want Agent to automatically update to the latest version without manual intervention.

**Independent Test**:
1. Agent reports old version → Server sends update_required message
2. Agent receives update → pulls new image, starts new container, exits old process
3. Image pull fails → Agent logs error and continues running current version

### Tasks

- [ ] T067 [US5] Implement update_required WebSocket message handler in agent/internal/websocket/handlers.go
- [ ] T068 [US5] Implement self-update logic in agent/internal/update/updater.go (pull image, start container, exit)
- [ ] T069 [US5] Implement graceful shutdown waiting for tasks in agent/internal/task/shutdown.go
- [ ] T070 [P] [US5] Implement version check in server/internal/service/agent_service.go
- [ ] T071 [P] [US5] Implement update_required message sender in server/internal/websocket/notifier.go

## Phase 8: Polish & Cross-Cutting Concerns

**Goal**: Add production-ready features, error handling, logging, and documentation.

### Tasks

- [ ] T072 [P] Add structured logging to Agent using zerolog in agent/internal/logger/
- [ ] T073 [P] Add structured logging to Server Agent handlers using existing logger
- [ ] T074 [P] Implement API key masking in logs for both Agent and Server
- [ ] T075 [P] Add OOM score adjustment (-500 for Agent, 500 for Worker) in agent/internal/docker/runner.go
- [ ] T076 [P] Implement ping/pong WebSocket keep-alive in agent/internal/websocket/keepalive.go
- [ ] T077 [P] Add error message truncation (4KB limit) in agent/internal/docker/logs.go
- [ ] T078 [P] Create Agent Dockerfile in agent/Dockerfile
- [ ] T079 [P] Create Agent installation script in scripts/install-agent.sh (implement FR-015: support PUBLIC_URL config or infer from request headers; generate SERVER_URL for Agent; auto-convert https→wss, http→ws for WebSocket)
- [ ] T080 [P] Add Agent deployment documentation in docs/agent-deployment.md
- [ ] T081 [P] Add WebSocket protocol documentation in docs/websocket-protocol.md
- [ ] T082 Add integration test for Agent connection flow in agent/test/integration/connection_test.go
- [ ] T083 Add integration test for task execution flow in agent/test/integration/task_test.go
- [ ] T084 Add unit tests for WebSocket reconnection logic in agent/internal/websocket/reconnect_test.go
- [ ] T085 Add unit tests for task assignment with database locks in server/internal/repository/scan_task_test.go
- [ ] T086 Add unit tests for status update validation in server/internal/service/scan_task_service_test.go
- [ ] T087 Add unit tests for heartbeat processing in server/internal/websocket/handlers_test.go

## Dependencies

### User Story Completion Order

```
Setup (Phase 1) → Foundational (Phase 2)
                        ↓
                    US1 (MVP) ← Must complete first
                        ↓
                      US2 ← Depends on US1 (connection required for task execution)
                        ↓
                    ┌───┴───┐
                   US3     US4 ← Can be done in parallel after US2
                    └───┬───┘
                        ↓
                      US5 ← Depends on US1 (connection) and US3 (heartbeat/version)
                        ↓
                     Polish
```

### Critical Path

1. **Setup + Foundational** (T001-T020): Required for all subsequent work
2. **US1: Connection** (T021-T033): Blocking for all other user stories
3. **US2: Task Execution** (T034-T051): Blocking for US3, US4, US5
4. **US3 + US4** (T052-T066): Can be implemented in parallel
5. **US5: Auto-Update** (T067-T071): Requires US1 and US3
6. **Polish** (T072-T087): Can be done incrementally throughout

## Parallel Execution Opportunities

### Phase 1 (Setup)
- T004, T005, T006: Agent dependencies (parallel)
- T008, T009: Server dependencies (parallel)

### Phase 2 (Foundational)
- T013, T014: Models (parallel)
- T017, T018: Redis cache (parallel after T017)

### Phase 3 (US1)
- T027, T028, T029: Server-side connection handling (parallel)

### Phase 4 (US2)
- T043, T044, T045, T046, T047: Server-side task APIs (parallel)

### Phase 5 (US3)
- T058, T059, T060, T061: Server-side heartbeat and monitoring (parallel)

### Phase 6 (US4)
- T065, T066: Server-side config update (parallel)

### Phase 7 (US5)
- T070, T071: Server-side version check (parallel)

### Phase 8 (Polish)
- T072, T073, T074, T075, T076, T077, T078, T079, T080, T081: Documentation and logging (all parallel)
- T082-T087: Tests (can be done in parallel with implementation)

## Task Metrics

- **Total Tasks**: 87
- **Setup Phase**: 9 tasks
- **Foundational Phase**: 11 tasks
- **User Story 1 (P1)**: 13 tasks - MVP
- **User Story 2 (P1)**: 18 tasks
- **User Story 3 (P2)**: 11 tasks
- **User Story 4 (P2)**: 4 tasks
- **User Story 5 (P3)**: 5 tasks
- **Polish Phase**: 16 tasks
- **Parallelizable Tasks**: 42 tasks (48%)

## Notes

- All file paths are absolute from project root: `/Users/yangyang/Desktop/orbit/`
- Tasks marked with [P] can be executed in parallel with other [P] tasks in the same phase
- Tasks marked with [US1], [US2], etc. belong to specific user stories
- Each user story phase is independently testable
- MVP scope (Phase 3) delivers core value and enables early validation
- Tests are included in Polish phase but can be written incrementally during implementation
