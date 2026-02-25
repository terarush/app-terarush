package service

import (
	"context"
	"errors"
	"fmt"
	"go-modular/internal/pkg/config"
	productRepo "go-modular/modules/products/domain/repository"
	transactionEntity "go-modular/modules/transactions/domain/entity"
	"go-modular/modules/transactions/domain/repository"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInvalidStatus       = errors.New("invalid transaction status")
	ErrProductNotAvailable = errors.New("product not available")
)

type TransactionService struct {
	repo        repository.TransactionRepository
	productRepo productRepo.ProductRepository
	snapClient  snap.Client
}

// NewTransactionService creates a new transaction service
func NewTransactionService(repo repository.TransactionRepository, productRepo productRepo.ProductRepository) *TransactionService {
	// Initialize Midtrans Snap client
	var snapClient snap.Client
	snapClient.New(
		config.GetString("MIDTRANS_SERVER_KEY"),
		midtrans.Sandbox, // Change to midtrans.Production for production
	)

	return &TransactionService{
		repo:        repo,
		productRepo: productRepo,
		snapClient:  snapClient,
	}
}

// CreateTransaction creates a new transaction and gets Snap token from Midtrans
func (s *TransactionService) CreateTransaction(ctx context.Context, userID uint, productID uint, quantity int, userEmail, userName string) (*transactionEntity.Transaction, error) {
	// Get product details
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Check product availability
	if !product.IsAvailable() {
		return nil, ErrProductNotAvailable
	}

	// Check stock
	if product.Stock != -1 && product.Stock < quantity {
		return nil, errors.New("insufficient stock")
	}

	// Calculate amount
	amount := product.Price * float64(quantity)

	// Generate unique order ID
	orderID := fmt.Sprintf("TRS-%d-%d", time.Now().Unix(), userID)

	// Create transaction
	transaction := &transactionEntity.Transaction{
		OrderID:   orderID,
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Amount:    amount,
		Status:    transactionEntity.StatusPending,
	}

	// Set expiration (24 hours from now)
	expiredAt := time.Now().Add(24 * time.Hour)
	transaction.ExpiredAt = &expiredAt

	// Create Midtrans Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: userEmail,
			FName: userName,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    fmt.Sprintf("PROD-%d", productID),
				Name:  product.Name,
				Price: int64(product.Price),
				Qty:   int32(quantity),
			},
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeGopay,
			snap.PaymentTypeShopeepay,
			snap.PaymentTypeBCAVA,
			snap.PaymentTypeBNIVA,
			snap.PaymentTypeBRIVA,
			snap.PaymentTypePermataVA,
		},
	}

	// Get Snap token from Midtrans
	snapResp, err := s.snapClient.CreateTransaction(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create Midtrans transaction: %w", err)
	}

	// Store Snap token and redirect URL
	transaction.SnapToken = snapResp.Token
	transaction.PaymentURL = snapResp.RedirectURL

	// Save transaction to database
	err = s.repo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Decrement product stock if not unlimited
	if product.Stock != -1 {
		err = s.productRepo.UpdateStock(ctx, productID, quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}
	}

	return transaction, nil
}

// GetTransactionByID gets a transaction by ID
func (s *TransactionService) GetTransactionByID(ctx context.Context, id uint) (*transactionEntity.Transaction, error) {
	transaction, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	return transaction, nil
}

// GetTransactionByOrderID gets a transaction by order ID
func (s *TransactionService) GetTransactionByOrderID(ctx context.Context, orderID string) (*transactionEntity.Transaction, error) {
	transaction, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	return transaction, nil
}

// GetUserTransactions gets all transactions for a user
func (s *TransactionService) GetUserTransactions(ctx context.Context, userID uint) ([]*transactionEntity.Transaction, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// GetAllTransactions gets all transactions (admin only)
func (s *TransactionService) GetAllTransactions(ctx context.Context) ([]*transactionEntity.Transaction, error) {
	return s.repo.FindAll(ctx)
}

// HandleMidtransNotification handles Midtrans payment notification/webhook
func (s *TransactionService) HandleMidtransNotification(ctx context.Context, notification map[string]interface{}) error {
	orderID, ok := notification["order_id"].(string)
	if !ok {
		return errors.New("invalid order_id in notification")
	}

	transactionStatus, ok := notification["transaction_status"].(string)
	if !ok {
		return errors.New("invalid transaction_status in notification")
	}

	fraudStatus := ""
	if val, ok := notification["fraud_status"].(string); ok {
		fraudStatus = val
	}

	// Get transaction from database
	transaction, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	// Update transaction based on Midtrans status
	var newStatus transactionEntity.TransactionStatus

	switch transactionStatus {
	case "capture":
		if fraudStatus == "challenge" {
			newStatus = transactionEntity.StatusProcessing
		} else if fraudStatus == "accept" {
			newStatus = transactionEntity.StatusSuccess
			now := time.Now()
			transaction.PaidAt = &now
		}
	case "settlement":
		newStatus = transactionEntity.StatusSuccess
		now := time.Now()
		transaction.PaidAt = &now
	case "pending":
		newStatus = transactionEntity.StatusPending
	case "deny", "cancel":
		newStatus = transactionEntity.StatusFailed
	case "expire":
		newStatus = transactionEntity.StatusExpired
	default:
		newStatus = transactionEntity.StatusPending
	}

	// Update transaction
	transaction.Status = newStatus
	transaction.MidtransStatus = transactionStatus
	if midtransID, ok := notification["transaction_id"].(string); ok {
		transaction.MidtransTransactionID = midtransID
	}
	if paymentType, ok := notification["payment_type"].(string); ok {
		transaction.PaymentMethod = paymentType
	}

	return s.repo.Update(ctx, transaction)
}

// CancelTransaction cancels a pending transaction
func (s *TransactionService) CancelTransaction(ctx context.Context, orderID string, userID uint) error {
	transaction, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return ErrTransactionNotFound
	}

	// Check ownership
	if transaction.UserID != userID {
		return errors.New("unauthorized")
	}

	// Check if can be cancelled
	if transaction.Status != transactionEntity.StatusPending {
		return ErrInvalidStatus
	}

	transaction.Status = transactionEntity.StatusCancelled
	return s.repo.Update(ctx, transaction)
}
