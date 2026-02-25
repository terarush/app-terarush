package request

// CreateTransactionRequest represents the request to create a transaction
type CreateTransactionRequest struct {
	ProductID uint `json:"product_id" validate:"required,gt=0"`
	Quantity  int  `json:"quantity" validate:"required,gte=1"`
}
