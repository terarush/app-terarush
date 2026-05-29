import { z } from "zod";

// Create blog schema
export const createBlogSchema = z.object({
  title: z
    .string()
    .min(1, "Title is required")
    .min(5, "Title must be at least 5 characters")
    .max(255, "Title must not exceed 255 characters"),
  slug: z
    .string()
    .min(1, "Slug is required")
    .min(3, "Slug must be at least 3 characters")
    .max(255, "Slug must not exceed 255 characters")
    .regex(/^[a-z0-9]+(?:-[a-z0-9]+)*$/, "Slug must contain only lowercase letters, numbers, and hyphens"),
  content: z
    .string()
    .min(1, "Content is required")
    .min(10, "Content must be at least 10 characters"),
  excerpt: z
    .string()
    .max(500, "Excerpt must not exceed 500 characters")
    .optional()
    .default(""),
  category: z
    .string()
    .min(1, "Category is required")
    .max(100, "Category must not exceed 100 characters"),
  tags: z
    .string()
    .max(255, "Tags must not exceed 255 characters")
    .optional()
    .default(""),
  image: z
    .string()
    .url("Invalid image URL")
    .optional()
    .default(""),
  is_published: z.boolean().optional().default(false),
});

export type CreateBlogFormData = z.infer<typeof createBlogSchema>;

// Update blog schema
export const updateBlogSchema = createBlogSchema;

export type UpdateBlogFormData = z.infer<typeof updateBlogSchema>;

// Blog filter/search schema
export const blogFilterSchema = z.object({
  search: z.string().optional().default(""),
  page: z.number().int().positive().optional().default(1),
  page_size: z.number().int().min(1).max(100).optional().default(10),
  category: z.string().optional().default(""),
});

export type BlogFilterFormData = z.infer<typeof blogFilterSchema>;

// Delete blog confirmation schema
export const deleteBlogSchema = z.object({
  id: z.number().int().positive("Blog ID must be a positive number"),
  confirm: z.boolean().refine((val) => val === true, {
    message: "You must confirm the deletion",
  }),
});

export type DeleteBlogFormData = z.infer<typeof deleteBlogSchema>;
