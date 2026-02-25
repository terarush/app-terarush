package handler

import (
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/transactions/domain/service"
	"go-modular/modules/transactions/dto/request"
	"go-modular/modules/transactions/dto/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
	log                *logger.Logger
	r                  *utils.Response
}

func NewTransactionHandler(transactionService *service.TransactionService, log *logger.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		log:                log,
		r:                  &utils.Response{},
	}
}

// CreateTransaction godoc
// @Summary Create transaction (checkout)
// @Description Create a new transaction for product purchase
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction body request.CreateTransactionRequest true "Transaction data"
// @Success 201 {object} utils.SuccessResponse{data=response.TransactionResponse}
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	// Get user ID from context
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		// Try to get from claims map
		claims, ok := c.Get("user").(map[string]interface{})
		if !ok {
			return h.r.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		}
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return h.r.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID")
		}
		userID = uint(userIDFloat)
	}

	// Get user email and name from context
	claims, _ := c.Get("user").(map[string]interface{})
	userEmail, _ := claims["email"].(string)
	userName, _ := claims["name"].(string)

	var req request.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	transaction, err := h.transactionService.CreateTransaction(
		c.Request().Context(),
		userID,
		req.ProductID,
		req.Quantity,
		userEmail,
		userName,
	)
	if err != nil {
		h.log.Error("Failed to create transaction:", err)
		if err == service.ErrProductNotAvailable {
			return h.r.ErrorResponse(c, http.StatusBadRequest, "Product not available")
		}
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(transaction), "Transaction created successfully")
}

// GetUserTransactions godoc
// @Summary Get user transactions
// @Description Get all transactions for the authenticated user
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponse{data=[]response.TransactionResponse}
// @Router /transactions/my [get]
func (h *TransactionHandler) GetUserTransactions(c echo.Context) error {
	// Get user ID from context
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		claims, ok := c.Get("user").(map[string]interface{})
		if !ok {
			return h.r.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		}
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return h.r.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID")
		}
		userID = uint(userIDFloat)
	}

	transactions, err := h.transactionService.GetUserTransactions(c.Request().Context(), userID)
	if err != nil {
		h.log.Error("Failed to get user transactions:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to get transactions")
	}

	return h.r.SuccessResponse(c, response.FromEntities(transactions), "Transactions retrieved successfully")
}

// GetTransactionByID godoc
// @Summary Get transaction by ID
// @Description Get transaction details by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} utils.SuccessResponse{data=response.TransactionResponse}
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
	}

	transaction, err := h.transactionService.GetTransactionByID(c.Request().Context(), uint(id))
	if err != nil {
		if err == service.ErrTransactionNotFound {
			return h.r.ErrorResponse(c, http.StatusNotFound, "Transaction not found")
		}
		h.log.Error("Failed to get transaction:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to get transaction")
	}

	// Check ownership (users can only see their own transactions, admins can see all)
	userRole, _ := c.Get("role").(string)
	if userRole != "admin" {
		claims, _ := c.Get("user").(map[string]interface{})
		userIDFloat, _ := claims["user_id"].(float64)
		userID := uint(userIDFloat)

		if transaction.UserID != userID {
			return h.r.ErrorResponse(c, http.StatusForbidden, "Forbidden")
		}
	}

	return h.r.SuccessResponse(c, response.FromEntity(transaction), "Transaction retrieved successfully")
}

// GetAllTransactions godoc
// @Summary Get all transactions (Admin only)
// @Description Get all transactions in the system
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponse{data=[]response.TransactionResponse}
// @Router /transactions [get]
func (h *TransactionHandler) GetAllTransactions(c echo.Context) error {
	transactions, err := h.transactionService.GetAllTransactions(c.Request().Context())
	if err != nil {
		h.log.Error("Failed to get all transactions:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to get transactions")
	}

	return h.r.SuccessResponse(c, response.FromEntities(transactions), "Transactions retrieved successfully")
}

// MidtransWebhook godoc
// @Summary Midtrans payment webhook
// @Description Handle Midtrans payment notification
// @Tags transactions
// @Accept json
// @Produce json
// @Param notification body map[string]interface{} true "Midtrans notification"
// @Success 200 {object} utils.SuccessResponse
// @Router /transactions/webhook [post]
func (h *TransactionHandler) MidtransWebhook(c echo.Context) error {
	var notification map[string]interface{}
	if err := c.Bind(&notification); err != nil {
		h.log.Error("Failed to bind notification:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	h.log.Info("Received Midtrans notification:", notification)

	err := h.transactionService.HandleMidtransNotification(c.Request().Context(), notification)
	if err != nil {
		h.log.Error("Failed to handle notification:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to process notification")
	}

	return h.r.SuccessResponse(c, nil, "Notification processed successfully")
}

// CancelTransaction godoc
// @Summary Cancel transaction
// @Description Cancel a pending transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "Order ID"
// @Success 200 {object} utils.SuccessResponse
// @Router /transactions/{order_id}/cancel [post]
func (h *TransactionHandler) CancelTransaction(c echo.Context) error {
	orderID := c.Param("order_id")

	// Get user ID from context
	claims, ok := c.Get("user").(map[string]interface{})
	if !ok {
		return h.r.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return h.r.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID")
	}
	userID := uint(userIDFloat)

	err := h.transactionService.CancelTransaction(c.Request().Context(), orderID, userID)
	if err != nil {
		if err == service.ErrTransactionNotFound {
			return h.r.ErrorResponse(c, http.StatusNotFound, "Transaction not found")
		}
		if err.Error() == "unauthorized" {
			return h.r.ErrorResponse(c, http.StatusForbidden, "Forbidden")
		}
		h.log.Error("Failed to cancel transaction:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.r.SuccessResponse(c, nil, "Transaction cancelled successfully")
}

// RegisterRoutes registers transaction routes
func (h *TransactionHandler) RegisterRoutes(e *echo.Echo, basePath string, authMiddleware echo.MiddlewareFunc, adminMiddleware echo.MiddlewareFunc) {
	group := e.Group(basePath + "/transactions")

	// Public webhook endpoint (no auth)
	group.POST("/webhook", h.MidtransWebhook)

	// Authenticated user routes
	auth := group.Group("", authMiddleware)
	auth.POST("", h.CreateTransaction)
	auth.GET("/my", h.GetUserTransactions)
	auth.GET("/:id", h.GetTransactionByID)
	auth.POST("/:order_id/cancel", h.CancelTransaction)

	// Admin only routes
	admin := group.Group("", authMiddleware, adminMiddleware)
	admin.GET("", h.GetAllTransactions)
}
