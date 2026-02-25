package entity

import (
	"encoding/json"
	"time"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	StatusPending    TransactionStatus = "pending"
	StatusProcessing TransactionStatus = "processing"
	StatusSuccess    TransactionStatus = "success"
	StatusFailed     TransactionStatus = "failed"
	StatusCancelled  TransactionStatus = "cancelled"
	StatusExpired    TransactionStatus = "expired"
)

// Transaction represents a payment transaction
type Transaction struct {
	ID                    uint              `gorm:"primaryKey" json:"id"`
	OrderID               string            `gorm:"size:255;unique;not null" json:"order_id"`
	UserID                uint              `gorm:"not null" json:"user_id"`
	ProductID             uint              `gorm:"not null" json:"product_id"`
	Quantity              int               `gorm:"not null;default:1" json:"quantity"`
	Amount                float64           `gorm:"type:decimal(15,2);not null" json:"amount"`
	Status                TransactionStatus `gorm:"type:enum('pending','processing','success','failed','cancelled','expired');default:'pending'" json:"status"`
	PaymentMethod         string            `gorm:"size:50" json:"payment_method,omitempty"`
	PaymentURL            string            `gorm:"type:text" json:"payment_url,omitempty"`
	SnapToken             string            `gorm:"type:text" json:"snap_token,omitempty"`
	MidtransTransactionID string            `gorm:"size:255" json:"midtrans_transaction_id,omitempty"`
	MidtransStatus        string            `gorm:"size:50" json:"midtrans_status,omitempty"`
	PaidAt                *time.Time        `json:"paid_at,omitempty"`
	ExpiredAt             *time.Time        `json:"expired_at,omitempty"`
	Metadata              json.RawMessage   `gorm:"type:json" json:"metadata,omitempty"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}

// TableName specifies the table name for Transaction
func (Transaction) TableName() string {
	return "transactions"
}

// IsCompleted checks if transaction is completed (success or failed)
func (t *Transaction) IsCompleted() bool {
	return t.Status == StatusSuccess || t.Status == StatusFailed || t.Status == StatusCancelled || t.Status == StatusExpired
}

// CanBePaid checks if transaction can be paid
func (t *Transaction) CanBePaid() bool {
	return t.Status == StatusPending && (t.ExpiredAt == nil || time.Now().Before(*t.ExpiredAt))
}
