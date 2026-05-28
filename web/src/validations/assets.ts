import { z } from "zod";
import { ALLOWED_MIME_TYPES, MAX_FILE_SIZE } from "@/lib/api/assets";

// Asset file validation schema
export const assetFileSchema = z
  .instanceof(File)
  .refine((file) => file.size > 0, "File is required")
  .refine((file) => file.size <= MAX_FILE_SIZE, {
    message: `File size must be less than ${MAX_FILE_SIZE / 1024 / 1024}MB`,
  })
  .refine((file) => ALLOWED_MIME_TYPES[file.type as keyof typeof ALLOWED_MIME_TYPES], {
    message: "File type not allowed. Allowed types: images, PDFs, documents, and videos",
  });

export type AssetFileFormData = z.infer<typeof assetFileSchema>;

// Single asset upload schema
export const uploadAssetSchema = z.object({
  file: assetFileSchema,
});

export type UploadAssetFormData = z.infer<typeof uploadAssetSchema>;

// Multiple assets upload schema
export const uploadAssetsSchema = z.object({
  files: z
    .array(assetFileSchema)
    .min(1, "At least one file is required")
    .max(10, "Maximum 10 files can be uploaded at once"),
});

export type UploadAssetsFormData = z.infer<typeof uploadAssetsSchema>;

// Asset deletion confirmation schema
export const deleteAssetSchema = z.object({
  id: z.string().min(1, "Asset ID is required"),
  confirm: z.boolean().refine((val) => val === true, {
    message: "You must confirm the deletion",
  }),
});

export type DeleteAssetFormData = z.infer<typeof deleteAssetSchema>;
