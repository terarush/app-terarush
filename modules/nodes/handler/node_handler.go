package handler

import (
	"fmt"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/utils"
	"go-modular/internal/pkg/validator"
	"go-modular/modules/nodes/domain/service"
	"go-modular/modules/nodes/dto/request"
	"go-modular/modules/nodes/dto/response"
	"strconv"

	"github.com/labstack/echo"
)

type NodeHandler struct {
	nodeService service.NodeService
	log         *logger.Logger
	validator   *validator.CustomValidator
	response    *utils.Response
}

func NewNodeHandler(log *logger.Logger, validator *validator.CustomValidator, nodeService service.NodeService) *NodeHandler {
	return &NodeHandler{
		nodeService: nodeService,
		log:         log,
		validator:   validator,
		response:    &utils.Response{},
	}
}

// getUserIDFromContext extracts user ID from JWT claims in context
func getUserIDFromContext(c echo.Context) (uint, error) {
	claims, ok := c.Get("user").(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("user not found in context")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in claims")
	}

	return uint(userIDFloat), nil
}

// CreateNode creates a new node
// POST /api/v1/nodes
func (h *NodeHandler) CreateNode(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.log.Error("Failed to get user ID from context", "error", err)
		return h.response.UnauthorizedResponse(c, "Unauthorized")
	}

	var req request.CreateNodeRequest
	if err := c.Bind(&req); err != nil {
		h.log.Error("Failed to bind request", "error", err)
		return h.response.BadRequestResponse(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		h.log.Error("Validation failed", "error", err)
		return h.response.ValidationErrorResponse(c, err.Error())
	}

	node, err := h.nodeService.CreateNode(userID, &req)
	if err != nil {
		h.log.Error("Failed to create node", "error", err)
		return h.response.InternalServerErrorResponse(c, err.Error())
	}

	h.log.Info("Node created successfully", "node_id", node.ID, "user_id", userID)
	return h.response.CreatedResponse(c, response.ToNodeResponse(node), "Node created successfully")
}

// GetNodes returns all nodes for the current user
// GET /api/v1/nodes
func (h *NodeHandler) GetNodes(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.log.Error("Failed to get user ID from context", "error", err)
		return h.response.UnauthorizedResponse(c, "Unauthorized")
	}

	nodes, err := h.nodeService.GetNodesByUserID(userID)
	if err != nil {
		h.log.Error("Failed to get nodes", "error", err)
		return h.response.InternalServerErrorResponse(c, "Failed to get nodes")
	}

	h.log.Info("Nodes retrieved successfully", "user_id", userID, "count", len(nodes))
	return h.response.SuccessResponse(c, response.ToNodeResponseList(nodes), "Nodes retrieved successfully")
}

// GetNode returns a single node by ID
// GET /api/v1/nodes/:id
func (h *NodeHandler) GetNode(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.log.Error("Failed to get user ID from context", "error", err)
		return h.response.UnauthorizedResponse(c, "Unauthorized")
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.log.Error("Invalid node ID", "error", err)
		return h.response.BadRequestResponse(c, "Invalid node ID")
	}

	node, err := h.nodeService.GetNodeByID(uint(id), userID)
	if err != nil {
		h.log.Error("Failed to get node", "error", err)
		return h.response.NotFoundResponse(c, "Node not found")
	}

	h.log.Info("Node retrieved successfully", "node_id", id, "user_id", userID)
	return h.response.SuccessResponse(c, response.ToNodeResponse(node), "Node retrieved successfully")
}

// UpdateNode updates a node
// PUT /api/v1/nodes/:id
func (h *NodeHandler) UpdateNode(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.log.Error("Failed to get user ID from context", "error", err)
		return h.response.UnauthorizedResponse(c, "Unauthorized")
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.log.Error("Invalid node ID", "error", err)
		return h.response.BadRequestResponse(c, "Invalid node ID")
	}

	var req request.UpdateNodeRequest
	if err := c.Bind(&req); err != nil {
		h.log.Error("Failed to bind request", "error", err)
		return h.response.BadRequestResponse(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		h.log.Error("Validation failed", "error", err)
		return h.response.ValidationErrorResponse(c, err.Error())
	}

	node, err := h.nodeService.UpdateNode(uint(id), userID, &req)
	if err != nil {
		h.log.Error("Failed to update node", "error", err)
		return h.response.InternalServerErrorResponse(c, err.Error())
	}

	h.log.Info("Node updated successfully", "node_id", id, "user_id", userID)
	return h.response.SuccessResponse(c, response.ToNodeResponse(node), "Node updated successfully")
}

// DeleteNode deletes a node
// DELETE /api/v1/nodes/:id
func (h *NodeHandler) DeleteNode(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.log.Error("Failed to get user ID from context", "error", err)
		return h.response.UnauthorizedResponse(c, "Unauthorized")
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.log.Error("Invalid node ID", "error", err)
		return h.response.BadRequestResponse(c, "Invalid node ID")
	}

	if err := h.nodeService.DeleteNode(uint(id), userID); err != nil {
		h.log.Error("Failed to delete node", "error", err)
		return h.response.InternalServerErrorResponse(c, err.Error())
	}

	h.log.Info("Node deleted successfully", "node_id", id, "user_id", userID)
	return h.response.SuccessResponse(c, nil, "Node deleted successfully")
}

// NodeAction handles node actions (start, stop, restart)
// POST /api/v1/nodes/:id/action
func (h *NodeHandler) NodeAction(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.log.Error("Failed to get user ID from context", "error", err)
		return h.response.UnauthorizedResponse(c, "Unauthorized")
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.log.Error("Invalid node ID", "error", err)
		return h.response.BadRequestResponse(c, "Invalid node ID")
	}

	var req request.NodeActionRequest
	if err := c.Bind(&req); err != nil {
		h.log.Error("Failed to bind request", "error", err)
		return h.response.BadRequestResponse(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		h.log.Error("Validation failed", "error", err)
		return h.response.ValidationErrorResponse(c, err.Error())
	}

	var actionErr error
	switch req.Action {
	case "start":
		actionErr = h.nodeService.StartNode(uint(id), userID)
	case "stop":
		actionErr = h.nodeService.StopNode(uint(id), userID)
	case "restart":
		actionErr = h.nodeService.RestartNode(uint(id), userID)
	default:
		return h.response.BadRequestResponse(c, "Invalid action")
	}

	if actionErr != nil {
		h.log.Error("Failed to perform action", "action", req.Action, "error", actionErr)
		return h.response.InternalServerErrorResponse(c, actionErr.Error())
	}

	h.log.Info("Node action completed successfully", "node_id", id, "action", req.Action, "user_id", userID)
	return h.response.SuccessResponse(c, nil, fmt.Sprintf("Node %s successfully", req.Action))
}
