import { useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import type { Comment } from "@/lib/api/comments";
import { CommentForm } from "./CommentForm";
import { Button } from "@/components/ui/button";
import { Trash2, MessageCircleReply, MoreVertical } from "lucide-react";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

// Simple date formatter
const formatDistanceToNow = (date: Date) => {
	const now = new Date();
	const diffMs = now.getTime() - date.getTime();
	const diffSecs = Math.floor(diffMs / 1000);
	const diffMins = Math.floor(diffSecs / 60);
	const diffHours = Math.floor(diffMins / 60);
	const diffDays = Math.floor(diffHours / 24);

	if (diffSecs < 60) return "just now";
	if (diffMins < 60) return `${diffMins}m ago`;
	if (diffHours < 24) return `${diffHours}h ago`;
	if (diffDays < 7) return `${diffDays}d ago`;

	// Fallback to formatted date
	return date.toLocaleDateString();
};

interface CommentListProps {
	comments: Comment[];
	onAddReply: (data: any) => Promise<void>;
	onEditComment: (id: string, data: any) => Promise<void>;
	onDeleteComment: (id: string) => Promise<void>;
	isSubmitting?: boolean;
	isLoading?: boolean;
}

interface CommentItemProps {
	comment: Comment;
	postId: number;
	depth?: number;
	onAddReply: (data: any) => Promise<void>;
	onEditComment: (id: string, data: any) => Promise<void>;
	onDeleteComment: (id: string) => Promise<void>;
	isSubmitting?: boolean;
}

function CommentItem({
	comment,
	postId,
	depth = 0,
	onAddReply,
	onEditComment,
	onDeleteComment,
	isSubmitting = false,
}: CommentItemProps) {
	const { user } = useAuth();
	const [isReplying, setIsReplying] = useState(false);
	const [isEditing, setIsEditing] = useState(false);
	const isOwner = user?.id === comment.user_id;
	const isDeleted = comment.deleted_at != null && comment.deleted_at !== undefined;

	const handleReplySubmit = async (data: any) => {
		await onAddReply(data);
		setIsReplying(false);
	};

	const handleEditSubmit = async (data: any) => {
		await onEditComment(comment.id, data);
		setIsEditing(false);
	};

	const handleDelete = async () => {
		if (confirm("Are you sure you want to delete this comment?")) {
			await onDeleteComment(comment.id);
		}
	};

	const paddingClass = depth > 0 ? `pl-${Math.min(depth * 4, 12)}` : "";

	if (isDeleted) {
		return (
			<div
				className={`py-4 ${paddingClass} border-l-2 border-border/30 ${
					depth > 0 ? "ml-4" : ""
				}`}
			>
				<p className="text-sm text-muted-foreground italic">This comment has been deleted.</p>
			</div>
		);
	}

	return (
		<div
			className={`py-4 border-l-2 border-border/30 ${paddingClass} ${
				depth > 0 ? "ml-4" : ""
			}`}
		>
			{/* Comment Header */}
			<div className="flex items-start justify-between mb-3">
				<div className="flex items-center gap-3 flex-1">
					<div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
						<span className="text-xs font-semibold text-primary">
							{comment.user_name?.charAt(0).toUpperCase() || "U"}
						</span>
					</div>
					<div className="flex-1 min-w-0">
						<p className="text-sm font-semibold text-foreground truncate">
							{comment.user_name}
						</p>
						<p className="text-xs text-muted-foreground">
							{formatDistanceToNow(new Date(comment.created_at * 1000))}
							{comment.updated_at !== comment.created_at && " (edited)"}
						</p>
					</div>
				</div>

				{/* Actions Menu */}
				{isOwner && (
					<DropdownMenu>
						<DropdownMenuTrigger asChild>
							<Button
								variant="ghost"
								size="sm"
								className="h-8 w-8 p-0"
							>
								<MoreVertical className="h-4 w-4" />
							</Button>
						</DropdownMenuTrigger>
						<DropdownMenuContent align="end" className="w-48">
							<DropdownMenuItem onClick={() => setIsEditing(true)}>
								Edit
							</DropdownMenuItem>
							<DropdownMenuItem
								onClick={handleDelete}
								className="text-destructive focus:text-destructive"
							>
								<Trash2 className="h-4 w-4 mr-2" />
								Delete
							</DropdownMenuItem>
						</DropdownMenuContent>
					</DropdownMenu>
				)}
			</div>

			{/* Comment Content */}
			{isEditing ? (
				<CommentForm
					postId={postId}
					onSubmit={handleEditSubmit}
					onCancel={() => setIsEditing(false)}
					isSubmitting={isSubmitting}
				/>
			) : (
				<>
					<div className="text-sm text-foreground mb-4 leading-relaxed break-words">
						{comment.content}
					</div>

					{/* Reply Button */}
					{depth < 2 && (
						<div className="flex gap-2">
							<Button
								variant="ghost"
								size="sm"
								className="h-8 px-2 text-xs gap-1"
								onClick={() => setIsReplying(!isReplying)}
							>
								<MessageCircleReply className="h-3 w-3" />
								Reply
							</Button>
						</div>
					)}
				</>
			)}

			{/* Reply Form */}
			{isReplying && !isEditing && (
				<div className="mt-4 p-4 bg-muted/30 rounded-lg border border-border/30">
					<CommentForm
						postId={postId}
						isReply
						parentId={comment.id}
						onSubmit={handleReplySubmit}
						onCancel={() => setIsReplying(false)}
						isSubmitting={isSubmitting}
					/>
				</div>
			)}

			{/* Nested Replies */}
			{comment.replies && comment.replies.length > 0 && (
				<div className="mt-4 space-y-0">
					{comment.replies.map((reply) => (
						<CommentItem
							key={reply.id}
							comment={reply}
							postId={postId}
							depth={depth + 1}
							onAddReply={onAddReply}
							onEditComment={onEditComment}
							onDeleteComment={onDeleteComment}
							isSubmitting={isSubmitting}
						/>
					))}
				</div>
			)}
		</div>
	);
}

export function CommentList({
	comments,
	onAddReply,
	onEditComment,
	onDeleteComment,
	isSubmitting = false,
	isLoading = false,
}: CommentListProps) {
	const postId = comments[0]?.post_id || 0;

	if (isLoading) {
		return (
			<div className="space-y-4">
				{[1, 2, 3].map((i) => (
					<div key={i} className="h-24 bg-muted/30 rounded-lg animate-pulse" />
				))}
			</div>
		);
	}

	if (comments.length === 0) {
		return (
			<div className="py-12 text-center">
				<MessageCircleReply className="h-12 w-12 mx-auto text-muted-foreground mb-4 opacity-50" />
				<p className="text-muted-foreground">No comments yet. Be the first to comment!</p>
			</div>
		);
	}

	return (
		<div className="space-y-0">
			{comments.map((comment) => (
				<CommentItem
					key={comment.id}
					comment={comment}
					postId={postId}
					onAddReply={onAddReply}
					onEditComment={onEditComment}
					onDeleteComment={onDeleteComment}
					isSubmitting={isSubmitting}
				/>
			))}
		</div>
	);
}
