package repository

import (
	"context"
	"go-modular/modules/transactions/domain/entity"
)

// TransactionRepository defines the interface for transaction data operations
type TransactionRepository interface {
	// Create operations
	Create(ctx context.Context, transaction *entity.Transaction) error

	// Read operations
	FindByID(ctx context.Context, id uint) (*entity.Transaction, error)
	FindByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.Transaction, error)
	FindAll(ctx context.Context) ([]*entity.Transaction, error)

	// Update operations
	Update(ctx context.Context, transaction *entity.Transaction) error
	UpdateStatus(ctx context.Context, orderID string, status entity.TransactionStatus) error
}
