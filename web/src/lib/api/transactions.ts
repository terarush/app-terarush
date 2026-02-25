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

// Get status badge color based on transaction status
export const getStatusColor = (status: TransactionStatus): string => {
	switch (status) {
		case "success":
			return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300";
		case "pending":
			return "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300";
		case "processing":
			return "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300";
		case "failed":
			return "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300";
		case "cancelled":
			return "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300";
		case "expired":
			return "bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300";
		default:
			return "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300";
	}
};

// Format status text
export const getStatusText = (status: TransactionStatus): string => {
	return status.charAt(0).toUpperCase() + status.slice(1);
};
