import { useState, useCallback } from "react";
import {
	getCommentsByPost,
	createComment,
	createReply,
	updateComment,
	deleteComment,
	type Comment,
	type CreateCommentRequest,
	type CreateReplyRequest,
	type UpdateCommentRequest,
} from "@/lib/api/comments";
import { toast } from "sonner";

interface UseCommentsOptions {
	postId: number;
	onSuccess?: (action: "create" | "reply" | "update" | "delete") => void;
	onError?: (error: Error) => void;
}

/**
 * Hook to manage blog comments with CRUD operations
 */
export const useComments = ({ postId, onSuccess, onError }: UseCommentsOptions) => {
	const [comments, setComments] = useState<Comment[]>([]);
	const [isLoading, setIsLoading] = useState(false);
	const [isSubmitting, setIsSubmitting] = useState(false);
	const [error, setError] = useState<Error | null>(null);

	// Fetch comments for the post
	const fetchComments = useCallback(async () => {
		try {
			setIsLoading(true);
			setError(null);
			const response = await getCommentsByPost(postId);
			setComments(response.comments || []);
		} catch (err) {
			const error = err instanceof Error ? err : new Error("Failed to fetch comments");
			setError(error);
			onError?.(error);
		} finally {
			setIsLoading(false);
		}
	}, [postId, onError]);

	// Add a new comment
	const addComment = useCallback(
		async (data: CreateCommentRequest) => {
			try {
				setIsSubmitting(true);
				setError(null);
				const newComment = await createComment(data);
				setComments((prev) => [newComment, ...prev]);
				toast.success("Comment posted successfully");
				onSuccess?.("create");
				return newComment;
			} catch (err) {
				const error = err instanceof Error ? err : new Error("Failed to post comment");
				setError(error);
				toast.error("Failed to post comment");
				onError?.(error);
				throw error;
			} finally {
				setIsSubmitting(false);
			}
		},
		[onSuccess, onError],
	);

	// Add a reply to a comment
	const addReply = useCallback(
		async (data: CreateReplyRequest) => {
			try {
				setIsSubmitting(true);
				setError(null);
				const reply = await createReply(data);

				// Add reply to the parent comment
				setComments((prev) =>
					prev.map((comment) => {
						if (comment.id === data.parent_id) {
							return {
								...comment,
								replies: [...(comment.replies || []), reply],
							};
						}
						return comment;
					}),
				);

				toast.success("Reply posted successfully");
				onSuccess?.("reply");
				return reply;
			} catch (err) {
				const error = err instanceof Error ? err : new Error("Failed to post reply");
				setError(error);
				toast.error("Failed to post reply");
				onError?.(error);
				throw error;
			} finally {
				setIsSubmitting(false);
			}
		},
		[onSuccess, onError],
	);

	// Update a comment
	const editComment = useCallback(
		async (id: string, data: UpdateCommentRequest) => {
			try {
				setIsSubmitting(true);
				setError(null);
				const updatedComment = await updateComment(id, data);

				// Update comment in list
				setComments((prev) =>
					prev.map((comment) => (comment.id === id ? updatedComment : comment)),
				);

				toast.success("Comment updated successfully");
				onSuccess?.("update");
				return updatedComment;
			} catch (err) {
				const error = err instanceof Error ? err : new Error("Failed to update comment");
				setError(error);
				toast.error("Failed to update comment");
				onError?.(error);
				throw error;
			} finally {
				setIsSubmitting(false);
			}
		},
		[onSuccess, onError],
	);

	// Delete a comment
	const removeComment = useCallback(
		async (id: string) => {
			try {
				setIsSubmitting(true);
				setError(null);
				await deleteComment(id);

				// Remove comment from list
				setComments((prev) => prev.filter((comment) => comment.id !== id));

				toast.success("Comment deleted successfully");
				onSuccess?.("delete");
			} catch (err) {
				const error = err instanceof Error ? err : new Error("Failed to delete comment");
				setError(error);
				toast.error("Failed to delete comment");
				onError?.(error);
				throw error;
			} finally {
				setIsSubmitting(false);
			}
		},
		[onSuccess, onError],
	);

	return {
		comments,
		isLoading,
		isSubmitting,
		error,
		fetchComments,
		addComment,
		addReply,
		editComment,
		removeComment,
	};
};
