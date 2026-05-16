import { apiClient } from "./client";

// Utility function to get absolute image URL
export const getImageUrl = (imagePath: string): string => {
	if (!imagePath) return "";
	
	// If it already starts with http, return as-is
	if (imagePath.startsWith("http")) {
		return imagePath;
	}
	
	// If it starts with /, it's a relative path - prepend API base
	if (imagePath.startsWith("/")) {
		const apiUrl = import.meta.env.VITE_API_URL || "";
		return `${apiUrl}${imagePath}`;
	}
	
	// Default: prepend /public/
	const apiUrl = import.meta.env.VITE_API_URL || "";
	return `${apiUrl}/${imagePath}`;
};

export interface UserInfo {
	id: number;
	name: string;
	email: string;
	avatar: string;
	bio: string;
}

export interface Blog {
	id: number;
	title: string;
	slug: string;
	content: string;
	excerpt: string;
	author: string;
	user_id: number;
	user?: UserInfo;
	category: string;
	tags: string;
	image: string;
	is_published: boolean;
	view_count: number;
	created_at: string;
	updated_at: string;
	published_at: string | null;
}

export interface BlogListResponse {
	blogs?: Blog[];
	total?: number;
	page?: number;
	page_size?: number;
}

export interface CreateBlogRequest {
	title: string;
	slug: string;
	content: string;
	excerpt: string;
	category?: string;
	tags?: string;
	image?: string;
	is_published?: boolean;
}

export interface UpdateBlogRequest extends Partial<CreateBlogRequest> {}

// Get all published blogs (public endpoint)
export const getBlogs = async (params?: {
	page?: number;
	page_size?: number;
}): Promise<BlogListResponse> => {
	const response = await apiClient.get("/blogs", { params });

	// Handle both array response and paginated response
	if (Array.isArray(response.data)) {
		return {
			blogs: response.data,
			total: response.data.length,
			page: 1,
			page_size: response.data.length,
		};
	}

	return response.data;
};

// Get blog by slug (public endpoint)
export const getBlogBySlug = async (slug: string): Promise<Blog> => {
	const response = await apiClient.get(`/blogs/${slug}`);
	return response.data;
};

// Get all blogs (admin only)
export const getAllBlogs = async (params?: {
	page?: number;
	page_size?: number;
}): Promise<BlogListResponse> => {
	const response = await apiClient.get("/admin/blogs", { params });

	// Handle both array response and paginated response
	if (Array.isArray(response.data)) {
		return {
			blogs: response.data,
			total: response.data.length,
			page: 1,
			page_size: response.data.length,
		};
	}

	return response.data;
};

// Get blog by ID (admin only)
export const getBlogById = async (id: number): Promise<Blog> => {
	const response = await apiClient.get(`/admin/blogs/${id}`);
	return response.data;
};

// Create blog (admin only)
export const createBlog = async (data: CreateBlogRequest): Promise<Blog> => {
	const response = await apiClient.post("/admin/blogs", data);
	return response.data;
};

// Update blog (admin only)
export const updateBlog = async (
	id: number,
	data: UpdateBlogRequest,
): Promise<Blog> => {
	const response = await apiClient.put(`/admin/blogs/${id}`, data);
	return response.data;
};

// Delete blog (admin only)
export const deleteBlog = async (id: number): Promise<void> => {
	await apiClient.delete(`/admin/blogs/${id}`);
};

// Upload blog image (admin only)
export interface UploadImageResponse {
	url: string;
	path: string;
}

export const uploadBlogImage = async (file: File): Promise<UploadImageResponse> => {
	const formData = new FormData();
	formData.append("image", file);

	const response = await apiClient.post("/admin/blogs/upload/image", formData, {
		headers: {
			"Content-Type": "multipart/form-data",
		},
	});

	return response.data;
};
