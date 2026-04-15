import { apiClient } from "./client";
import type { Product } from "./products";

export type TransactionStatus =
	| "pending"
	| "processing"
	| "success"
	| "failed"
	| "cancelled"
	| "expired";

export interface Transaction {
	id: number;
	order_id: string;
	user_id: number;
	product_id: number;
	product?: Product;
	quantity: number;
	amount: number;
	status: TransactionStatus;
	payment_method?: string;
	payment_url?: string;
	snap_token?: string;
	midtrans_transaction_id?: string;
	midtrans_status?: string;
	paid_at?: string;
	expired_at: string;
	metadata?: Record<string, any>;
	created_at: string;
	updated_at: string;
}

export interface TransactionListResponse {
	transactions: Transaction[];
	total: number;
	page: number;
	page_size: number;
}

export interface CheckoutRequest {
	product_id: number;
	quantity: number;
	image_tag: string;
}

export interface CheckoutResponse {
	transaction: Transaction;
	snap_token: string;
	payment_url: string;
}

// Create transaction / Checkout (authenticated)
export const checkout = async (
	data: CheckoutRequest,
): Promise<CheckoutResponse> => {
	const response = await apiClient.post("/transactions", data);
	return response.data;
};

// Get user's transactions (authenticated)
export const getMyTransactions = async (params?: {
	page?: number;
	page_size?: number;
	status?: TransactionStatus;
}): Promise<TransactionListResponse> => {
	const response = await apiClient.get("/transactions/my", { params });
	return response.data;
};

// Get transaction by ID (authenticated)
export const getTransactionById = async (id: number): Promise<Transaction> => {
	const response = await apiClient.get(`/transactions/${id}`);
	return response.data;
};

// Get all transactions (admin only)
export const getAllTransactions = async (params?: {
	page?: number;
	page_size?: number;
	status?: TransactionStatus;
	user_id?: number;
}): Promise<TransactionListResponse> => {
	const response = await apiClient.get("/transactions", { params });
	return response.data;
};

// Cancel transaction (authenticated)
export const cancelTransaction = async (orderId: string): Promise<void> => {
	await apiClient.post(`/transactions/${orderId}/cancel`);
};

// Get status badge variant based on transaction status
export const getStatusColor = (
	status: TransactionStatus,
): "default" | "destructive" | "outline" | "secondary" => {
	switch (status) {
		case "success":
			return "default";
		case "pending":
			return "secondary";
		case "processing":
			return "outline";
		case "failed":
			return "destructive";
		case "cancelled":
			return "secondary";
		case "expired":
			return "destructive";
		default:
			return "default";
	}
};

// Format status text
export const getStatusText = (status: TransactionStatus): string => {
	return status.charAt(0).toUpperCase() + status.slice(1);
};
