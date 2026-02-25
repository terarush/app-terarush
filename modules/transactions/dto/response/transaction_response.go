package response

import (
	"go-modular/modules/transactions/domain/entity"
	"time"
)

// TransactionResponse represents the transaction response
type TransactionResponse struct {
	ID                    uint                     `json:"id"`
	OrderID               string                   `json:"order_id"`
	UserID                uint                     `json:"user_id"`
	ProductID             uint                     `json:"product_id"`
	Quantity              int                      `json:"quantity"`
	Amount                float64                  `json:"amount"`
	Status                entity.TransactionStatus `json:"status"`
	PaymentMethod         string                   `json:"payment_method,omitempty"`
	PaymentURL            string                   `json:"payment_url,omitempty"`
	SnapToken             string                   `json:"snap_token,omitempty"`
	MidtransTransactionID string                   `json:"midtrans_transaction_id,omitempty"`
	MidtransStatus        string                   `json:"midtrans_status,omitempty"`
	PaidAt                *time.Time               `json:"paid_at,omitempty"`
	ExpiredAt             *time.Time               `json:"expired_at,omitempty"`
	CreatedAt             time.Time                `json:"created_at"`
	UpdatedAt             time.Time                `json:"updated_at"`
}

// FromEntity converts entity to response
func FromEntity(transaction *entity.Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID:                    transaction.ID,
		OrderID:               transaction.OrderID,
		UserID:                transaction.UserID,
		ProductID:             transaction.ProductID,
		Quantity:              transaction.Quantity,
		Amount:                transaction.Amount,
		Status:                transaction.Status,
		PaymentMethod:         transaction.PaymentMethod,
		PaymentURL:            transaction.PaymentURL,
		SnapToken:             transaction.SnapToken,
		MidtransTransactionID: transaction.MidtransTransactionID,
		MidtransStatus:        transaction.MidtransStatus,
		PaidAt:                transaction.PaidAt,
		ExpiredAt:             transaction.ExpiredAt,
		CreatedAt:             transaction.CreatedAt,
		UpdatedAt:             transaction.UpdatedAt,
	}
}

// FromEntities converts multiple entities to responses
func FromEntities(transactions []*entity.Transaction) []*TransactionResponse {
	responses := make([]*TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = FromEntity(transaction)
	}
	return responses
}
