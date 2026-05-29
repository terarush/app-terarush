import { apiClient } from "./client";

export interface Comment {
	id: string;
	post_id: number;
	user_id: number;
	user_name: string;
	content: string;
	parent_id?: string | null;
	created_at: number; // Unix timestamp (seconds) from backend int64
	updated_at: number; // Unix timestamp (seconds) from backend int64
	deleted_at?: number | null;
	replies?: Comment[];
}

export interface CommentListResponse {
	comments?: Comment[];
	total?: number;
	page?: number;
	page_size?: number;
}

export interface CreateCommentRequest {
	content: string;
	post_id: number;
}

export interface CreateReplyRequest extends CreateCommentRequest {
	parent_id: string;
}

export interface UpdateCommentRequest {
	content: string;
	post_id: number;
}

// Get comments by blog post (public endpoint)
export const getCommentsByPost = async (
	postId: number,
	params?: {
		page?: number;
		page_size?: number;
	},
): Promise<CommentListResponse> => {
	const response = await apiClient.get(`/blogs/${postId}/comments`, {
		params,
	});

	// Handle both array response and paginated response
	if (Array.isArray(response.data)) {
		return {
			comments: response.data,
			total: response.data.length,
			page: 1,
			page_size: response.data.length,
		};
	}

	return response.data;
};

// Create comment (authenticated)
export const createComment = async (
	data: CreateCommentRequest,
): Promise<Comment> => {
	const response = await apiClient.post("/comments", data);
	return response.data;
};

// Create reply (authenticated)
export const createReply = async (
	data: CreateReplyRequest,
): Promise<Comment> => {
	const response = await apiClient.post("/comments/reply", data);
	return response.data;
};

// Update comment (authenticated)
export const updateComment = async (
	id: string,
	data: UpdateCommentRequest,
): Promise<Comment> => {
	const response = await apiClient.put(`/comments/${id}`, data);
	return response.data;
};

// Delete comment (authenticated)
export const deleteComment = async (id: string): Promise<void> => {
	await apiClient.delete(`/comments/${id}`);
};
