import { z } from "zod";

// Create comment schema
export const createCommentSchema = z.object({
  content: z
    .string()
    .min(1, "Comment content is required")
    .min(2, "Comment must be at least 2 characters")
    .max(5000, "Comment must not exceed 5000 characters"),
  post_id: z
    .number()
    .int()
    .positive("Post ID must be a positive number"),
});

export type CreateCommentFormData = z.infer<typeof createCommentSchema>;

// Create reply schema
export const createReplySchema = z.object({
  content: z
    .string()
    .min(1, "Reply content is required")
    .min(2, "Reply must be at least 2 characters")
    .max(5000, "Reply must not exceed 5000 characters"),
  post_id: z
    .number()
    .int()
    .positive("Post ID must be a positive number"),
  parent_id: z
    .string()
    .min(1, "Parent comment ID is required"),
});

export type CreateReplyFormData = z.infer<typeof createReplySchema>;

// Update comment schema
export const updateCommentSchema = z.object({
  content: z
    .string()
    .min(1, "Comment content is required")
    .min(2, "Comment must be at least 2 characters")
    .max(5000, "Comment must not exceed 5000 characters"),
  post_id: z
    .number()
    .int()
    .positive("Post ID must be a positive number"),
});

export type UpdateCommentFormData = z.infer<typeof updateCommentSchema>;

// Delete comment confirmation schema
export const deleteCommentSchema = z.object({
  id: z.string().min(1, "Comment ID is required"),
  confirm: z.boolean().refine((val) => val === true, {
    message: "You must confirm the deletion",
  }),
});

export type DeleteCommentFormData = z.infer<typeof deleteCommentSchema>;

// Comments filter schema
export const commentsFilterSchema = z.object({
  post_id: z
    .number()
    .int()
    .positive("Post ID must be a positive number"),
  page: z.number().int().positive().optional().default(1),
  page_size: z.number().int().min(1).max(100).optional().default(10),
});

export type CommentsFilterFormData = z.infer<typeof commentsFilterSchema>;
