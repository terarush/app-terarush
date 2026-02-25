package repository

import (
	"context"
	"errors"
	"go-modular/modules/transactions/domain/entity"

	"gorm.io/gorm"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type transactionRepositoryImpl struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

// Create inserts a new transaction
func (r *transactionRepositoryImpl) Create(ctx context.Context, transaction *entity.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

// FindByID finds a transaction by ID
func (r *transactionRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.WithContext(ctx).First(&transaction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ERR_RECORD_NOT_FOUND
		}
		return nil, err
	}
	return &transaction, nil
}

// FindByOrderID finds a transaction by order ID
func (r *transactionRepositoryImpl) FindByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ERR_RECORD_NOT_FOUND
		}
		return nil, err
	}
	return &transaction, nil
}

// FindByUserID finds all transactions for a user
func (r *transactionRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// FindAll finds all transactions
func (r *transactionRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// Update updates a transaction
func (r *transactionRepositoryImpl) Update(ctx context.Context, transaction *entity.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

// UpdateStatus updates transaction status
func (r *transactionRepositoryImpl) UpdateStatus(ctx context.Context, orderID string, status entity.TransactionStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Transaction{}).
		Where("order_id = ?", orderID).
		Update("status", status).
		Error
}
