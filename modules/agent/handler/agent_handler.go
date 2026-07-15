package handler

import (
	"fmt"
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/agent/domain/entity"
	"go-modular/modules/agent/domain/service"
	"go-modular/modules/agent/dto/request"
	"go-modular/modules/agent/dto/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AgentHandler struct {
	svc *service.AgentService
	log *logger.Logger
	evt *bus.EventBus
	r   *utils.Response
}

func NewAgentHandler(log *logger.Logger, evt *bus.EventBus, svc *service.AgentService) *AgentHandler {
	return &AgentHandler{svc: svc, log: log, evt: evt, r: &utils.Response{}}
}

func getUserID(c echo.Context) (uint, error) {
	claims, ok := c.Get("user").(map[string]any)
	if !ok {
		return 0, fmt.Errorf("user not found in context")
	}
	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID")
	}
	return uint(id), nil
}

// ─── Agents ──────────────────────────────────────────────────────────────────

// ListAgents lists all agents
// @Summary List agents
// @Description List all AI agents (admin only)
// @Tags Agent
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Agents list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents [get]
// @Security Bearer
func (h *AgentHandler) ListAgents(c echo.Context) error {
	agents, err := h.svc.GetAgents(c.Request().Context(), nil)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.AgentFromEntities(agents), "")
}

// GetAgent gets an agent by ID
// @Summary Get agent
// @Description Get a specific AI agent by ID (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Agent ID"
// @Success 200 {object} utils.SuccessResponse "Agent details"
// @Failure 404 {object} utils.ErrorResponse "Agent not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/{id} [get]
// @Security Bearer
func (h *AgentHandler) GetAgent(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid agent ID")
	}
	agent, err := h.svc.GetAgentByID(c.Request().Context(), uint(id))
	if err != nil {
		if err == service.ErrAgentNotFound {
			return h.r.NotFoundResponse(c, "agent not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.AgentFromEntity(agent), "")
}

// CreateAgent creates a new agent
// @Summary Create agent
// @Description Create a new AI agent configuration (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param request body request.CreateAgentRequest true "Agent creation data"
// @Success 201 {object} utils.CreatedResponse "Agent created"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents [post]
// @Security Bearer
func (h *AgentHandler) CreateAgent(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(request.CreateAgentRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}

	userID, _ := getUserID(c)
	agent := entity.NewAgent(req.Name, req.Description, req.Model, req.SystemPrompt, "")
	if req.MaxConcurrent > 0 {
		agent.MaxConcurrent = req.MaxConcurrent
	}
	if req.Timeout > 0 {
		agent.Timeout = req.Timeout
	}
	if req.RetryLimit > 0 {
		agent.RetryLimit = req.RetryLimit
	}
	agent.CreatedBy = userID

	if err := h.svc.CreateAgent(ctx, agent); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	h.evt.Publish(bus.Event{Type: "agent.created", Payload: agent})
	return h.r.CreatedResponse(c, response.AgentFromEntity(agent), "")
}

// UpdateAgent updates an agent
// @Summary Update agent
// @Description Update an existing AI agent configuration (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param id path int true "Agent ID"
// @Param request body request.UpdateAgentRequest true "Agent update data"
// @Success 200 {object} utils.SuccessResponse "Agent updated"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 404 {object} utils.ErrorResponse "Agent not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/{id} [put]
// @Security Bearer
func (h *AgentHandler) UpdateAgent(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid agent ID")
	}

	req := new(request.UpdateAgentRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}

	agent, err := h.svc.GetAgentByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "agent not found")
	}
	if req.Name != "" {
		agent.Name = req.Name
	}
	if req.Description != "" {
		agent.Description = req.Description
	}
	if req.Model != "" {
		agent.Model = req.Model
	}
	if req.SystemPrompt != "" {
		agent.SystemPrompt = req.SystemPrompt
	}
	if req.Status != "" {
		agent.Status = req.Status
	}
	if req.MaxConcurrent > 0 {
		agent.MaxConcurrent = req.MaxConcurrent
	}
	if req.Timeout > 0 {
		agent.Timeout = req.Timeout
	}
	if req.RetryLimit > 0 {
		agent.RetryLimit = req.RetryLimit
	}

	if err := h.svc.UpdateAgent(ctx, agent); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.AgentFromEntity(agent), "")
}

// DeleteAgent deletes an agent
// @Summary Delete agent
// @Description Delete an AI agent by ID (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Agent ID"
// @Success 204 "No content"
// @Failure 400 {object} utils.ErrorResponse "Invalid agent ID"
// @Failure 404 {object} utils.ErrorResponse "Agent not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/{id} [delete]
// @Security Bearer
func (h *AgentHandler) DeleteAgent(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid agent ID")
	}
	if err := h.svc.DeleteAgent(c.Request().Context(), uint(id)); err != nil {
		if err == service.ErrAgentNotFound {
			return h.r.NotFoundResponse(c, "agent not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// SetAgentStatus updates agent status
// @Summary Set agent status
// @Description Set the operational status of an AI agent (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param id path int true "Agent ID"
// @Param request body request.SetAgentStatusRequest true "Status data"
// @Success 200 {object} utils.SuccessResponse "Status updated"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 404 {object} utils.ErrorResponse "Agent not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/{id}/status [patch]
// @Security Bearer
func (h *AgentHandler) SetAgentStatus(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid agent ID")
	}
	req := new(request.SetAgentStatusRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}
	if err := h.svc.SetAgentStatus(ctx, uint(id), req.Status); err != nil {
		if err == service.ErrAgentNotFound {
			return h.r.NotFoundResponse(c, "agent not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "agent status updated")
}

// ─── Tasks ───────────────────────────────────────────────────────────────────

// ListTasks lists agent tasks
// @Summary List tasks
// @Description List all agent tasks with optional filters (admin only)
// @Tags Agent
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param status query string false "Filter by status"
// @Param agent_id query string false "Filter by agent ID"
// @Success 200 {object} utils.SuccessResponse "Tasks list with pagination"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/tasks [get]
// @Security Bearer
func (h *AgentHandler) ListTasks(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	status := c.QueryParam("status")
	agentID := c.QueryParam("agent_id")

	filter := map[string]any{}
	if status != "" {
		filter["status"] = status
	}
	if agentID != "" {
		filter["assigned_to"] = agentID
	}

	tasks, total, err := h.svc.GetTasks(c.Request().Context(), filter, page, limit)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.TaskListResponse{
		Tasks: response.TaskFromEntities(tasks),
		Total: total,
		Page:  page,
		Limit: limit,
	}, "")
}

// GetTask gets a task by ID
// @Summary Get task
// @Description Get a specific agent task by ID (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} utils.SuccessResponse "Task details"
// @Failure 404 {object} utils.ErrorResponse "Task not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/tasks/{id} [get]
// @Security Bearer
func (h *AgentHandler) GetTask(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid task ID")
	}
	task, err := h.svc.GetTaskByID(c.Request().Context(), uint(id))
	if err != nil {
		if err == service.ErrTaskNotFound {
			return h.r.NotFoundResponse(c, "task not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.TaskFromEntity(task), "")
}

// CreateTask creates a new task
// @Summary Create task
// @Description Create a new agent task (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param request body request.CreateTaskRequest true "Task creation data"
// @Success 201 {object} utils.CreatedResponse "Task created"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/tasks [post]
// @Security Bearer
func (h *AgentHandler) CreateTask(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(request.CreateTaskRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}

	userID, _ := getUserID(c)
	task := entity.NewTask(req.Title, req.Description, req.Category, req.Priority, userID)
	task.Tags = req.Tags
	task.Input = req.Input
	task.MaxAttempts = req.MaxAttempts
	task.ParentTaskID = req.ParentTaskID

	if err := h.svc.CreateTask(ctx, task); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	h.evt.Publish(bus.Event{Type: "agent.task.created", Payload: task})
	return h.r.CreatedResponse(c, response.TaskFromEntity(task), "")
}

// UpdateTask updates a task
// @Summary Update task
// @Description Update an existing agent task (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body request.UpdateTaskRequest true "Task update data"
// @Success 200 {object} utils.SuccessResponse "Task updated"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 404 {object} utils.ErrorResponse "Task not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/tasks/{id} [put]
// @Security Bearer
func (h *AgentHandler) UpdateTask(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid task ID")
	}
	req := new(request.UpdateTaskRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}

	task, err := h.svc.GetTaskByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "task not found")
	}
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.Category != "" {
		task.Category = req.Category
	}
	if req.Tags != "" {
		task.Tags = req.Tags
	}
	task.Priority = req.Priority
	if req.Progress > 0 {
		task.Progress = req.Progress
	}
	if req.Input != "" {
		task.Input = req.Input
	}
	if req.MaxAttempts > 0 {
		task.MaxAttempts = req.MaxAttempts
	}

	if err := h.svc.UpdateTask(ctx, task); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.TaskFromEntity(task), "")
}

// DeleteTask deletes a task
// @Summary Delete task
// @Description Delete an agent task by ID (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Task ID"
// @Success 204 "No content"
// @Failure 400 {object} utils.ErrorResponse "Invalid task ID"
// @Failure 404 {object} utils.ErrorResponse "Task not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/tasks/{id} [delete]
// @Security Bearer
func (h *AgentHandler) DeleteTask(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid task ID")
	}
	if err := h.svc.DeleteTask(c.Request().Context(), uint(id)); err != nil {
		if err == service.ErrTaskNotFound {
			return h.r.NotFoundResponse(c, "task not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// AssignTask assigns a task to an agent
// @Summary Assign task
// @Description Assign an agent task to a specific agent (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body request.AssignTaskRequest true "Assignment data"
// @Success 200 {object} utils.SuccessResponse "Task assigned"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 404 {object} utils.ErrorResponse "Task or agent not found"
// @Failure 409 {object} utils.ErrorResponse "Agent busy or disabled"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/tasks/{id}/assign [post]
// @Security Bearer
func (h *AgentHandler) AssignTask(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid task ID")
	}
	req := new(request.AssignTaskRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}
	if err := h.svc.AssignTask(ctx, uint(id), req.AgentID); err != nil {
		if err == service.ErrTaskNotFound || err == service.ErrAgentNotFound {
			return h.r.NotFoundResponse(c, err.Error())
		}
		if err == service.ErrAgentBusy || err == service.ErrAgentDisabled {
			return h.r.ConflictResponse(c, err.Error())
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "task assigned")
}

// StartTask starts executing a task
// @Summary Start task
// @Description Start executing an agent task (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} utils.SuccessResponse "Session details"
// @Failure 404 {object} utils.ErrorResponse "Task not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/tasks/{id}/start [post]
// @Security Bearer
func (h *AgentHandler) StartTask(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid task ID")
	}
	userID, _ := getUserID(c)
	session, err := h.svc.StartTask(ctx, uint(id), userID)
	if err != nil {
		if err == service.ErrTaskNotFound {
			return h.r.NotFoundResponse(c, "task not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.SessionFromEntity(session), "")
}

// CompleteTask marks a task as complete
// @Summary Complete task
// @Description Mark an agent task as completed with output (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body request.CompleteTaskRequest true "Completion data"
// @Success 200 {object} utils.SuccessResponse "Task completed"
// @Failure 404 {object} utils.ErrorResponse "Task not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/tasks/{id}/complete [post]
// @Security Bearer
func (h *AgentHandler) CompleteTask(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid task ID")
	}
	req := new(request.CompleteTaskRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := h.svc.CompleteTask(ctx, uint(id), req.Output, req.ResultSummary); err != nil {
		if err == service.ErrTaskNotFound {
			return h.r.NotFoundResponse(c, "task not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	h.evt.Publish(bus.Event{Type: "agent.task.completed", Payload: map[string]any{"task_id": id}})
	return h.r.SuccessResponse(c, nil, "task completed")
}

// FailTask marks a task as failed
// @Summary Fail task
// @Description Mark an agent task as failed with error details (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body request.FailTaskRequest true "Failure data"
// @Success 200 {object} utils.SuccessResponse "Task failed"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 404 {object} utils.ErrorResponse "Task not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/tasks/{id}/fail [post]
// @Security Bearer
func (h *AgentHandler) FailTask(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid task ID")
	}
	req := new(request.FailTaskRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}
	if err := h.svc.FailTask(ctx, uint(id), req.ErrorMessage); err != nil {
		if err == service.ErrTaskNotFound {
			return h.r.NotFoundResponse(c, "task not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "task failed")
}

// AutoAssign automatically assigns pending tasks
// @Summary Auto-assign tasks
// @Description Automatically assign pending tasks to idle agents (admin only)
// @Tags Agent
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Assigned count"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/tasks/auto-assign [post]
// @Security Bearer
func (h *AgentHandler) AutoAssign(c echo.Context) error {
	assigned, err := h.svc.AutoAssign(c.Request().Context())
	if err != nil {
		if err == service.ErrNoIdleAgent {
			return h.r.SuccessResponse(c, map[string]any{"assigned": 0, "message": "no idle agents"}, "")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, map[string]any{"assigned": len(assigned)}, "")
}

// ─── Sessions ────────────────────────────────────────────────────────────────

// ListSessions lists agent sessions
// @Summary List sessions
// @Description List all agent execution sessions with optional filters (admin only)
// @Tags Agent
// @Produce json
// @Param agent_id query string false "Filter by agent ID"
// @Param status query string false "Filter by status"
// @Success 200 {object} utils.SuccessResponse "Sessions list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/sessions [get]
// @Security Bearer
func (h *AgentHandler) ListSessions(c echo.Context) error {
	agentID := c.QueryParam("agent_id")
	status := c.QueryParam("status")
	filter := map[string]any{}
	if agentID != "" {
		filter["agent_id"] = agentID
	}
	if status != "" {
		filter["status"] = status
	}
	sessions, err := h.svc.GetSessions(c.Request().Context(), filter)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.SessionFromEntities(sessions), "")
}

// GetSession gets a session by ID
// @Summary Get session
// @Description Get a specific agent execution session by ID (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} utils.SuccessResponse "Session details"
// @Failure 404 {object} utils.ErrorResponse "Session not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/sessions/{id} [get]
// @Security Bearer
func (h *AgentHandler) GetSession(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid session ID")
	}
	session, err := h.svc.GetSessionByID(c.Request().Context(), uint(id))
	if err != nil {
		if err == service.ErrSessionNotFound {
			return h.r.NotFoundResponse(c, "session not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.SessionFromEntity(session), "")
}

// EndSession ends a session
// @Summary End session
// @Description Terminate an active agent execution session (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} utils.SuccessResponse "Session ended"
// @Failure 404 {object} utils.ErrorResponse "Session not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/sessions/{id}/end [post]
// @Security Bearer
func (h *AgentHandler) EndSession(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid session ID")
	}
	if err := h.svc.EndSession(c.Request().Context(), uint(id)); err != nil {
		if err == service.ErrSessionNotFound {
			return h.r.NotFoundResponse(c, "session not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "session ended")
}

// ─── Capabilities ────────────────────────────────────────────────────────────

// ListCapabilities lists agent capabilities
// @Summary List capabilities
// @Description List all capabilities for a specific agent (admin only)
// @Tags Agent
// @Produce json
// @Param agentId path int true "Agent ID"
// @Success 200 {object} utils.SuccessResponse "Capabilities list"
// @Failure 400 {object} utils.ErrorResponse "Invalid agent ID"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/agents/{agentId}/capabilities [get]
// @Security Bearer
func (h *AgentHandler) ListCapabilities(c echo.Context) error {
	agentID, err := strconv.ParseUint(c.Param("agentId"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid agent ID")
	}
	caps, err := h.svc.GetCapabilities(c.Request().Context(), uint(agentID))
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.CapabilityFromEntities(caps), "")
}

// AddCapability adds a capability to an agent
// @Summary Add capability
// @Description Add a new capability to an agent (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param request body request.AddCapabilityRequest true "Capability data"
// @Success 201 {object} utils.CreatedResponse "Capability added"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/capabilities [post]
// @Security Bearer
func (h *AgentHandler) AddCapability(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(request.AddCapabilityRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}
	cap := entity.NewCapability(req.AgentID, req.Capability)
	cap.Config = req.Config
	if err := h.svc.AddCapability(ctx, cap); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.CreatedResponse(c, response.CapabilityFromEntity(cap), "")
}

// RemoveCapability removes a capability
// @Summary Remove capability
// @Description Remove a capability by ID (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Capability ID"
// @Success 204 "No content"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/capabilities/{id} [delete]
// @Security Bearer
func (h *AgentHandler) RemoveCapability(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid capability ID")
	}
	if err := h.svc.RemoveCapability(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// ToggleCapability toggles a capability
// @Summary Toggle capability
// @Description Enable or disable an agent capability (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param id path int true "Capability ID"
// @Param request body request.ToggleCapabilityRequest true "Toggle data"
// @Success 200 {object} utils.SuccessResponse "Capability updated"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 404 {object} utils.ErrorResponse "Capability not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/capabilities/{id}/toggle [patch]
// @Security Bearer
func (h *AgentHandler) ToggleCapability(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid capability ID")
	}
	req := new(request.ToggleCapabilityRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}
	if err := h.svc.ToggleCapability(ctx, uint(id), req.Enabled); err != nil {
		if err == service.ErrCapabilityNotFound {
			return h.r.NotFoundResponse(c, "capability not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "capability updated")
}

// ─── Logs ────────────────────────────────────────────────────────────────────

// ListLogs lists agent logs
// @Summary List logs
// @Description List agent execution logs with optional filters (admin only)
// @Tags Agent
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param agent_id query string false "Filter by agent ID"
// @Param task_id query string false "Filter by task ID"
// @Param level query string false "Filter by log level"
// @Success 200 {object} utils.SuccessResponse "Logs list with pagination"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/logs [get]
// @Security Bearer
func (h *AgentHandler) ListLogs(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	agentID := c.QueryParam("agent_id")
	taskID := c.QueryParam("task_id")
	level := c.QueryParam("level")

	filter := map[string]any{}
	if agentID != "" {
		filter["agent_id"] = agentID
	}
	if taskID != "" {
		filter["task_id"] = taskID
	}
	if level != "" {
		filter["level"] = level
	}

	logs, total, err := h.svc.GetLogs(c.Request().Context(), filter, page, limit)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.LogListResponse{
		Logs:  response.LogFromEntities(logs),
		Total: total,
		Page:  page,
		Limit: limit,
	}, "")
}

// ─── Templates ───────────────────────────────────────────────────────────────

// ListTemplates lists agent templates
// @Summary List templates
// @Description List all task templates with optional category filter (admin only)
// @Tags Agent
// @Produce json
// @Param category query string false "Filter by category"
// @Success 200 {object} utils.SuccessResponse "Templates list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/templates [get]
// @Security Bearer
func (h *AgentHandler) ListTemplates(c echo.Context) error {
	category := c.QueryParam("category")
	filter := map[string]any{}
	if category != "" {
		filter["category"] = category
	}
	templates, err := h.svc.GetTemplates(c.Request().Context(), filter)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.TemplateFromEntities(templates), "")
}

// GetTemplate gets a template by ID
// @Summary Get template
// @Description Get a specific task template by ID (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Template ID"
// @Success 200 {object} utils.SuccessResponse "Template details"
// @Failure 404 {object} utils.ErrorResponse "Template not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/templates/{id} [get]
// @Security Bearer
func (h *AgentHandler) GetTemplate(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid template ID")
	}
	tmpl, err := h.svc.GetTemplateByID(c.Request().Context(), uint(id))
	if err != nil {
		if err == service.ErrTemplateNotFound {
			return h.r.NotFoundResponse(c, "template not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.TemplateFromEntity(tmpl), "")
}

// CreateTemplate creates a task template
// @Summary Create template
// @Description Create a new agent task template (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param request body request.CreateTemplateRequest true "Template data"
// @Success 201 {object} utils.CreatedResponse "Template created"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/templates [post]
// @Security Bearer
func (h *AgentHandler) CreateTemplate(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(request.CreateTemplateRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}
	userID, _ := getUserID(c)
	tmpl := entity.NewTemplate(req.Name, req.Description, req.Category, req.Prompt, userID)
	tmpl.Priority = req.Priority
	tmpl.MaxAttempts = req.MaxAttempts
	tmpl.InputSchema = req.InputSchema
	tmpl.Variables = req.Variables
	tmpl.Tags = req.Tags
	tmpl.IsPublic = req.IsPublic

	if err := h.svc.CreateTemplate(ctx, tmpl); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.CreatedResponse(c, response.TemplateFromEntity(tmpl), "")
}

// DeleteTemplate deletes a template
// @Summary Delete template
// @Description Delete a task template by ID (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Template ID"
// @Success 204 "No content"
// @Failure 404 {object} utils.ErrorResponse "Template not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/templates/{id} [delete]
// @Security Bearer
func (h *AgentHandler) DeleteTemplate(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid template ID")
	}
	if err := h.svc.DeleteTemplate(c.Request().Context(), uint(id)); err != nil {
		if err == service.ErrTemplateNotFound {
			return h.r.NotFoundResponse(c, "template not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// CreateTaskFromTemplate creates a task from a template
// @Summary Create task from template
// @Description Create a new task from an existing template (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param request body request.CreateTaskFromTemplateRequest true "Template task data"
// @Success 201 {object} utils.CreatedResponse "Task created"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 404 {object} utils.ErrorResponse "Template not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/templates/{id}/create-task [post]
// @Security Bearer
func (h *AgentHandler) CreateTaskFromTemplate(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(request.CreateTaskFromTemplateRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}
	userID, _ := getUserID(c)
	task, err := h.svc.CreateTaskFromTemplate(ctx, req.TemplateID, req.Input, userID)
	if err != nil {
		if err == service.ErrTemplateNotFound {
			return h.r.NotFoundResponse(c, "template not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.CreatedResponse(c, response.TaskFromEntity(task), "")
}

// ─── Schedules ───────────────────────────────────────────────────────────────

// ListSchedules lists agent schedules
// @Summary List schedules
// @Description List all scheduled agent tasks with optional filters (admin only)
// @Tags Agent
// @Produce json
// @Param agent_id query string false "Filter by agent ID"
// @Param status query string false "Filter by status"
// @Success 200 {object} utils.SuccessResponse "Schedules list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/schedules [get]
// @Security Bearer
func (h *AgentHandler) ListSchedules(c echo.Context) error {
	agentID := c.QueryParam("agent_id")
	status := c.QueryParam("status")
	filter := map[string]any{}
	if agentID != "" {
		filter["agent_id"] = agentID
	}
	if status != "" {
		filter["status"] = status
	}
	schedules, err := h.svc.GetSchedules(c.Request().Context(), filter)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.ScheduleFromEntities(schedules), "")
}

// GetSchedule gets a schedule by ID
// @Summary Get schedule
// @Description Get a specific agent schedule by ID (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} utils.SuccessResponse "Schedule details"
// @Failure 404 {object} utils.ErrorResponse "Schedule not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/schedules/{id} [get]
// @Security Bearer
func (h *AgentHandler) GetSchedule(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid schedule ID")
	}
	sched, err := h.svc.GetScheduleByID(c.Request().Context(), uint(id))
	if err != nil {
		if err == service.ErrScheduleNotFound {
			return h.r.NotFoundResponse(c, "schedule not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.ScheduleFromEntity(sched), "")
}

// CreateSchedule creates a schedule
// @Summary Create schedule
// @Description Create a new scheduled task for an agent (admin only)
// @Tags Agent
// @Accept json
// @Produce json
// @Param request body request.CreateScheduleRequest true "Schedule data"
// @Success 201 {object} utils.CreatedResponse "Schedule created"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/schedules [post]
// @Security Bearer
func (h *AgentHandler) CreateSchedule(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(request.CreateScheduleRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}
	userID, _ := getUserID(c)
	sched := entity.NewSchedule(req.AgentID, req.Name, req.CronExpr, userID)
	sched.TemplateID = req.TemplateID
	sched.Description = req.Description
	sched.Timezone = req.Timezone
	sched.MaxRuns = req.MaxRuns
	sched.InputConfig = req.InputConfig

	if err := h.svc.CreateSchedule(ctx, sched); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.CreatedResponse(c, response.ScheduleFromEntity(sched), "")
}

// DeleteSchedule deletes a schedule
// @Summary Delete schedule
// @Description Delete a scheduled agent task by ID (admin only)
// @Tags Agent
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 204 "No content"
// @Failure 404 {object} utils.ErrorResponse "Schedule not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/schedules/{id} [delete]
// @Security Bearer
func (h *AgentHandler) DeleteSchedule(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "invalid schedule ID")
	}
	if err := h.svc.DeleteSchedule(c.Request().Context(), uint(id)); err != nil {
		if err == service.ErrScheduleNotFound {
			return h.r.NotFoundResponse(c, "schedule not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// ProcessDueSchedules processes due schedules
// @Summary Process due schedules
// @Description Process all schedules that are due for execution (admin only)
// @Tags Agent
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Tasks created count"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/schedules/process-due [post]
// @Security Bearer
func (h *AgentHandler) ProcessDueSchedules(c echo.Context) error {
	tasks, err := h.svc.ProcessDueSchedules(c.Request().Context())
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, map[string]any{"tasks_created": len(tasks)}, "")
}

// ─── Dashboard ───────────────────────────────────────────────────────────────

// GetStats gets agent dashboard stats
// @Summary Get agent stats
// @Description Get aggregate statistics for agent dashboard (admin only)
// @Tags Agent
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Stats"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/agents/stats [get]
// @Security Bearer
func (h *AgentHandler) GetStats(c echo.Context) error {
	stats, err := h.svc.GetStats(c.Request().Context())
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, stats, "")
}

// ─── Event handlers ──────────────────────────────────────────────────────────

func (h *AgentHandler) HandleAgentEvent(event bus.Event) {
	h.log.Info("Agent event received: %s", event.Type)
}

// ─── Routes ──────────────────────────────────────────────────────────────────

func (h *AgentHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	adm := e.Group(basePath+"/admin/agents", middleware.Auth, middleware.AdminOnly)

	// Agent CRUD
	adm.GET("", h.ListAgents)
	adm.POST("", h.CreateAgent)
	adm.GET("/:id", h.GetAgent)
	adm.PUT("/:id", h.UpdateAgent)
	adm.DELETE("/:id", h.DeleteAgent)
	adm.PATCH("/:id/status", h.SetAgentStatus)

	// Task management
	adm.GET("/tasks", h.ListTasks)
	adm.POST("/tasks", h.CreateTask)
	adm.GET("/tasks/:id", h.GetTask)
	adm.PUT("/tasks/:id", h.UpdateTask)
	adm.DELETE("/tasks/:id", h.DeleteTask)
	adm.POST("/tasks/:id/assign", h.AssignTask)
	adm.POST("/tasks/:id/start", h.StartTask)
	adm.POST("/tasks/:id/complete", h.CompleteTask)
	adm.POST("/tasks/:id/fail", h.FailTask)
	adm.POST("/tasks/auto-assign", h.AutoAssign)

	// Sessions
	adm.GET("/sessions", h.ListSessions)
	adm.GET("/sessions/:id", h.GetSession)
	adm.POST("/sessions/:id/end", h.EndSession)

	// Capabilities
	adm.GET("/agents/:agentId/capabilities", h.ListCapabilities)
	adm.POST("/capabilities", h.AddCapability)
	adm.DELETE("/capabilities/:id", h.RemoveCapability)
	adm.PATCH("/capabilities/:id/toggle", h.ToggleCapability)

	// Logs
	adm.GET("/logs", h.ListLogs)

	// Templates
	adm.GET("/templates", h.ListTemplates)
	adm.POST("/templates", h.CreateTemplate)
	adm.GET("/templates/:id", h.GetTemplate)
	adm.DELETE("/templates/:id", h.DeleteTemplate)
	adm.POST("/templates/:id/create-task", h.CreateTaskFromTemplate)

	// Schedules
	adm.GET("/schedules", h.ListSchedules)
	adm.POST("/schedules", h.CreateSchedule)
	adm.GET("/schedules/:id", h.GetSchedule)
	adm.DELETE("/schedules/:id", h.DeleteSchedule)
	adm.POST("/schedules/process-due", h.ProcessDueSchedules)

	// Dashboard
	adm.GET("/stats", h.GetStats)
}
