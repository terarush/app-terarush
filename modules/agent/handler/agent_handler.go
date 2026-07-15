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

func (h *AgentHandler) ListAgents(c echo.Context) error {
	agents, err := h.svc.GetAgents(c.Request().Context(), nil)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.AgentFromEntities(agents), "")
}

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

func (h *AgentHandler) ProcessDueSchedules(c echo.Context) error {
	tasks, err := h.svc.ProcessDueSchedules(c.Request().Context())
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, map[string]any{"tasks_created": len(tasks)}, "")
}

// ─── Dashboard ───────────────────────────────────────────────────────────────

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
