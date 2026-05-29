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
	FormLabel,
	FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Loader2, Send, X } from "lucide-react";

interface CommentFormProps {
	postId: number;
	isReply?: boolean;
	parentId?: string;
	onSubmit: (data: CreateCommentFormData | CreateReplyFormData) => Promise<void>;
	onCancel?: () => void;
	isSubmitting?: boolean;
}

export function CommentForm({
	postId,
	isReply = false,
	parentId,
	onSubmit,
	onCancel,
	isSubmitting = false,
}: CommentFormProps) {
	const { user } = useAuth();
	const [loading, setLoading] = useState(false);

	const schema = isReply ? createReplySchema : createCommentSchema;
	type FormData = typeof isReply extends true ? CreateReplyFormData : CreateCommentFormData;

	const form = useForm<FormData>({
		resolver: zodResolver(schema),
		defaultValues: isReply
			? ({
				content: "",
				post_id: postId,
				parent_id: parentId,
			} as FormData)
			: ({
				content: "",
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
				{/* User Info */}
				{user && (
					<div className="flex items-center gap-3 mb-4">
						<div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
							<span className="text-sm font-semibold text-primary">
								{user.name?.charAt(0).toUpperCase() || "U"}
							</span>
						</div>
						<div className="text-sm">
							<p className="font-semibold text-foreground">{user.name || "Anonymous"}</p>
							<p className="text-xs text-muted-foreground">{user.email}</p>
						</div>
					</div>
				)}

				{/* Content Field */}
				<FormField
					control={form.control}
					name="content"
					render={({ field }) => (
						<FormItem>
							<FormLabel className="text-sm font-medium">
								{isReply ? "Write a reply" : "Write a comment"}
							</FormLabel>
							<FormControl>
								<Textarea
									placeholder={
										isReply
											? "Share your reply..."
											: "Share your thoughts on this post..."
									}
									className="min-h-24 rounded-lg bg-muted/50 border-border/50 resize-none"
									{...field}
									disabled={isDisabled}
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>

				{/* Action Buttons */}
				<div className="flex justify-end gap-2">
					{isReply && onCancel && (
						<Button
							type="button"
							variant="outline"
							size="sm"
							onClick={onCancel}
							disabled={isDisabled}
							className="rounded-md h-9"
						>
							<X className="h-4 w-4 mr-2" />
							Cancel
						</Button>
					)}
					<Button
						type="submit"
						size="sm"
						disabled={isDisabled}
						className="rounded-md h-9 gap-2"
					>
						{loading || isSubmitting ? (
							<Loader2 className="h-4 w-4 animate-spin" />
						) : (
							<Send className="h-4 w-4" />
						)}
						{isReply ? "Post Reply" : "Post Comment"}
					</Button>
				</div>

				{!user && (
					<div className="p-3 bg-muted/50 rounded-md border border-border/50 text-sm text-muted-foreground">
						Please log in to post a comment.
					</div>
				)}
			</form>
		</Form>
	);
}
