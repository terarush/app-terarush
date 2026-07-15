import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAuth } from "@/contexts/AuthContext";
import {
	createCommentSchema,
	createReplySchema,
	type CreateCommentFormData,
	type CreateReplyFormData,
} from "@/validations/comments.validation";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Loader2, Send, X, User } from "lucide-react";
import { motion } from "framer-motion";

interface CommentFormProps {
	postId: number;
	isReply?: boolean;
	parentId?: string;
	onSubmit: (data: CreateCommentFormData | CreateReplyFormData) => Promise<void>;
	onCancel?: () => void;
	isSubmitting?: boolean;
	initialContent?: string;
}

export function CommentForm({
	postId,
	isReply = false,
	parentId,
	onSubmit,
	onCancel,
	isSubmitting = false,
	initialContent,
}: CommentFormProps) {
	const { user } = useAuth();
	const [loading, setLoading] = useState(false);
	const [focused, setFocused] = useState(false);

	const schema = isReply ? createReplySchema : createCommentSchema;
	type FormData = typeof isReply extends true
		? CreateReplyFormData
		: CreateCommentFormData;

	const form = useForm<FormData>({
		resolver: zodResolver(schema),
		defaultValues: isReply
			? ({
					content: initialContent || "",
					post_id: postId,
					parent_id: parentId,
				} as FormData)
			: ({
					content: initialContent || "",
					post_id: postId,
				} as FormData),
	});

	const handleSubmit = async (data: FormData) => {
		try {
			setLoading(true);
			await onSubmit(data);
			form.reset();
		} finally {
			setLoading(false);
		}
	};

	const isDisabled = loading || isSubmitting || !user;

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
				{user && (
					<motion.div
						initial={{ opacity: 0, y: -8 }}
						animate={{ opacity: 1, y: 0 }}
						className="flex items-center gap-3"
					>
						<div className="w-9 h-9 rounded-full bg-gradient-to-br from-primary/15 to-accent/15 ring-2 ring-primary/8 overflow-hidden shrink-0 flex items-center justify-center">
							{user.avatar ? (
								<img src={user.avatar} alt={user.name} className="w-full h-full object-cover" />
							) : (
								<User className="h-4 w-4 text-primary/60" />
							)}
						</div>
						<div className="text-sm leading-tight min-w-0">
							<p className="font-semibold text-foreground truncate">{user.name || "Anonymous"}</p>
							<p className="text-xs text-muted-foreground">{isReply ? "Replying..." : "Share your thoughts"}</p>
						</div>
					</motion.div>
				)}

				{isReply && onCancel && (
					<div className="flex items-center gap-3 text-xs text-muted-foreground/60">
						<span className="text-primary/70 font-medium">Reply</span>
						<div className="h-px flex-1 bg-border/40" />
					</div>
				)}

				<FormField
					control={form.control}
					name="content"
					render={({ field }) => (
						<FormItem>
							<FormControl>
								<div
									className={`relative rounded-xl border transition-all duration-200 ${
										focused
											? "border-primary/40 shadow-[0_0_0_3px_rgba(var(--primary)/0.08)]"
											: "border-border/60 hover:border-border"
									} ${isDisabled ? "opacity-60" : ""}`}
								>
									<Textarea
										placeholder={isReply ? "Write a reply..." : "Share your thoughts on this post..."}
										className="min-h-[100px] border-0 bg-muted/25 rounded-xl resize-none focus-visible:ring-0 focus-visible:ring-offset-0 p-4 text-sm leading-relaxed"
										{...field}
										disabled={isDisabled}
										onFocus={() => setFocused(true)}
										onBlur={() => setFocused(false)}
									/>
								</div>
							</FormControl>
							<FormMessage className="text-xs ml-1" />
						</FormItem>
					)}
				/>

				<div className="flex items-center justify-between gap-3">
					{!user && (
						<div className="flex-1 px-4 py-2.5 bg-muted/20 rounded-lg border border-border/30 text-sm text-muted-foreground">
							Please{" "}
							<a href="/login" className="text-primary font-semibold hover:underline">
								log in
							</a>{" "}
							to post.
						</div>
					)}

					<div className="flex items-center gap-2 ml-auto">
						{isReply && onCancel && (
							<Button
								type="button"
								variant="ghost"
								size="sm"
								onClick={onCancel}
								disabled={isDisabled}
								className="rounded-lg h-9 px-3 text-xs text-muted-foreground"
							>
								<X className="h-3.5 w-3.5 mr-1.5" />
								Cancel
							</Button>
						)}
						<Button
							type="submit"
							size="sm"
							disabled={isDisabled || !form.watch("content")?.trim()}
							className="rounded-lg h-9 px-4 text-xs font-semibold gap-1.5"
						>
							{loading || isSubmitting ? (
								<Loader2 className="h-3.5 w-3.5 animate-spin" />
							) : (
								<Send className="h-3.5 w-3.5" />
							)}
							{isReply ? "Reply" : "Comment"}
						</Button>
					</div>
				</div>
			</form>
		</Form>
	);
}
