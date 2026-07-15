import { useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import type { Comment } from "@/lib/api/comments";
import { CommentForm } from "./CommentForm";
import { Button } from "@/components/ui/button";
import {
	Trash2,
	MessageCircleReply,
	MoreVertical,
	Edit3,
} from "lucide-react";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { motion } from "framer-motion";

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
	return date.toLocaleDateString("en-US", { month: "short", day: "numeric" });
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
	const [isDeleted, setIsDeleted] = useState(
		comment.deleted_at != null && comment.deleted_at !== undefined,
	);
	const isOwner = user?.id === comment.user_id;

	const handleReplySubmit = async (data: any) => {
		await onAddReply(data);
		setIsReplying(false);
	};

	const handleEditSubmit = async (data: any) => {
		await onEditComment(comment.id, data);
		setIsEditing(false);
	};

	const handleDelete = async () => {
		await onDeleteComment(comment.id);
		setIsDeleted(true);
	};

	const formatDate = (timestamp: number) => {
		const date = new Date(timestamp * 1000);
		return formatDistanceToNow(date);
	};

	const isEdited = comment.updated_at !== comment.created_at;

	if (isDeleted) {
		return (
			<div className={`relative py-4 ${depth > 0 ? "ml-8 pl-5 border-l border-border/15" : ""}`}>
				<div className="flex items-center gap-3 py-2 px-4 rounded-lg bg-muted/15 border border-border/10">
					<div className="w-7 h-7 rounded-full bg-muted/20 flex items-center justify-center">
						<div className="h-2.5 w-2.5 rounded-full bg-muted-foreground/20" />
					</div>
					<p className="text-sm text-muted-foreground/60 italic">This comment has been deleted.</p>
				</div>
			</div>
		);
	}

	const avatarEl = (
		<div className="w-9 h-9 rounded-full bg-gradient-to-br from-primary/8 to-accent/8 ring-1 ring-primary/5 overflow-hidden shrink-0 flex items-center justify-center">
			{comment.user_avatar ? (
				<img src={comment.user_avatar} alt={comment.user_name} className="w-full h-full object-cover" />
			) : (
				<span className="text-xs font-semibold text-primary/60">
					{comment.user_name?.charAt(0).toUpperCase() || "U"}
				</span>
			)}
		</div>
	);

	return (
		<motion.div
			initial={{ opacity: 0, y: 8 }}
			animate={{ opacity: 1, y: 0 }}
			className={`relative ${depth > 0 ? "ml-8 pl-5 border-l border-primary/8" : ""}`}
		>
			<div className="group">
				<div className="flex items-start gap-3">
					{avatarEl}

					<div className="flex-1 min-w-0">
						{/* Header */}
						<div className="flex items-center gap-2 mb-0.5">
							<span className="text-sm font-semibold text-foreground">{comment.user_name}</span>
							<span className="text-xs text-muted-foreground/50">{formatDate(comment.created_at)}</span>
							{isEdited && <span className="text-[10px] text-muted-foreground/30">· edited</span>}
						</div>

						{/* Body */}
						{isEditing ? (
							<div className="mt-2 mb-3">
								<CommentForm
									postId={postId}
									onSubmit={handleEditSubmit}
									onCancel={() => setIsEditing(false)}
									isSubmitting={isSubmitting}
									initialContent={comment.content}
								/>
							</div>
						) : (
							<>
								<p className="text-sm text-foreground/85 leading-relaxed mt-1 break-words">
									{comment.content}
								</p>

								<div className="flex items-center gap-3 mt-2">
									<button
										onClick={() => setIsReplying(!isReplying)}
										className="text-xs text-muted-foreground/50 hover:text-primary transition-colors font-medium flex items-center gap-1 py-0.5"
									>
										<MessageCircleReply className="h-3 w-3" />
										Reply
									</button>
								</div>
							</>
						)}

						{/* Reply form */}
						{isReplying && !isEditing && (
							<motion.div
								initial={{ opacity: 0, y: -6 }}
								animate={{ opacity: 1, y: 0 }}
								className="mt-3 mb-2 p-4 rounded-xl bg-muted/10 border border-border/20"
							>
								<CommentForm
									postId={postId}
									isReply
									parentId={comment.id}
									onSubmit={handleReplySubmit}
									onCancel={() => setIsReplying(false)}
									isSubmitting={isSubmitting}
								/>
							</motion.div>
						)}
					</div>

					{/* Actions */}
					{isOwner && !isEditing && (
						<div className="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity duration-150">
							<DropdownMenu>
								<DropdownMenuTrigger asChild>
									<Button variant="ghost" size="sm" className="h-7 w-7 p-0 rounded-lg">
										<MoreVertical className="h-3.5 w-3.5 text-muted-foreground/40" />
									</Button>
								</DropdownMenuTrigger>
								<DropdownMenuContent align="end" className="w-36 rounded-lg">
									<DropdownMenuItem onClick={() => setIsEditing(true)} className="gap-2 text-xs">
										<Edit3 className="h-3.5 w-3.5" />
										Edit
									</DropdownMenuItem>
									<DropdownMenuSeparator />
									<DropdownMenuItem onClick={handleDelete} className="gap-2 text-xs text-destructive focus:text-destructive">
										<Trash2 className="h-3.5 w-3.5" />
										Delete
									</DropdownMenuItem>
								</DropdownMenuContent>
							</DropdownMenu>
						</div>
					)}
				</div>
			</div>

			{/* Nested replies */}
			{comment.replies && comment.replies.length > 0 && (
				<div className="mt-1">
					{comment.replies.map((reply: any) => (
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
		</motion.div>
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
					<div key={i} className="flex items-start gap-3 animate-pulse">
						<div className="w-9 h-9 rounded-full bg-muted/20 shrink-0" />
						<div className="flex-1 space-y-2 py-1">
							<div className="h-3 w-24 bg-muted/20 rounded" />
							<div className="h-4 w-full bg-muted/10 rounded" />
							<div className="h-4 w-2/3 bg-muted/10 rounded" />
						</div>
					</div>
				))}
			</div>
		);
	}

	if (comments.length === 0) {
		return (
			<div className="py-16 text-center">
				<div className="w-12 h-12 rounded-xl bg-muted/20 flex items-center justify-center mx-auto mb-3">
					<MessageCircleReply className="h-6 w-6 text-muted-foreground/30" />
				</div>
				<p className="text-sm text-muted-foreground font-medium">No comments yet</p>
				<p className="text-xs text-muted-foreground/50 mt-1">Be the first to share your thoughts</p>
			</div>
		);
	}

	return (
		<div className="space-y-2">
			{comments.map((comment: any) => (
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
