import { apiClient } from "./client";

export interface Product {
	id: number;
	name: string;
	description: string;
	price: number;
	cpu_cores: number;
	ram_mb: number;
	storage_gb: number;
	bandwidth_gb: number;
	is_active: boolean;
	stock: number;
	image_url: string;
	docker_image: string;
	image_tags: string[];
	created_at: string;
	updated_at: string;
}

export interface ProductListResponse {
	products?: Product[];
	total?: number;
	page?: number;
	page_size?: number;
}

export interface CreateProductRequest {
	name: string;
	description: string;
	price: number;
	cpu_cores: number;
	ram_mb: number;
	storage_gb: number;
	bandwidth_gb: number;
	is_active?: boolean;
	stock?: number;
	image_url?: string;
	docker_image: string;
	image_tags: string[];
}

export interface UpdateProductRequest extends Partial<CreateProductRequest> {}

// Get all products (public endpoint)
export const getProducts = async (params?: {
	page?: number;
	page_size?: number;
	is_active?: boolean;
}): Promise<ProductListResponse> => {
	const response = await apiClient.get("/products", { params });
	
	// Handle both array response (current backend) and paginated response (future)
	if (Array.isArray(response.data)) {
		return {
			products: response.data,
			total: response.data.length,
			page: 1,
			page_size: response.data.length,
		};
	}
	
	return response.data;
};

// Get product by ID
export const getProductById = async (id: number): Promise<Product> => {
	const response = await apiClient.get(`/products/${id}`);
	return response.data;
};

// Create product (admin only)
export const createProduct = async (
	data: CreateProductRequest,
): Promise<Product> => {
	const response = await apiClient.post("/products", data);
	return response.data;
};

// Update product (admin only)
export const updateProduct = async (
	id: number,
	data: UpdateProductRequest,
): Promise<Product> => {
	const response = await apiClient.put(`/products/${id}`, data);
	return response.data;
};

// Delete product (admin only)
export const deleteProduct = async (id: number): Promise<void> => {
	await apiClient.delete(`/products/${id}`);
};
