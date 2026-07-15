import { apiClient } from "./client";

/**
 * Asset interface matching the backend entity
 */
export interface Asset {
	id: string;
	url: string;
	path: string;
	file_name: string;
	mime_type: string;
	created_at: string;
	updated_at: string;
}

/**
 * Upload response from the backend
 */
export interface UploadAssetResponse {
	id: string;
	url: string;
	path: string;
}

/**
 * Asset listing response
 */
export interface AssetListResponse {
	assets?: Asset[];
	total?: number;
	page?: number;
	page_size?: number;
}

/**
 * Allowed file types for upload
 */
export const ALLOWED_MIME_TYPES = {
	// Images
	"image/jpeg": true,
	"image/png": true,
	"image/gif": true,
	"image/webp": true,
	"image/svg+xml": true,
	// Documents
	"application/pdf": true,
	"text/plain": true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	// Spreadsheets
	"application/vnd.ms-excel": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
	// Videos
	"video/mp4": true,
	"video/webm": true,
};

/**
 * Max file size in bytes (10MB)
 */
export const MAX_FILE_SIZE = 10 * 1024 * 1024;

/**
 * Get the absolute URL for an asset
 * @param assetPath - The asset path or URL
 * @returns Absolute URL
 */
export const getAssetUrl = (assetPath: string): string => {
	if (!assetPath) return "";

	// If it already starts with http, return as-is
	if (assetPath.startsWith("http")) {
		return assetPath;
	}

	// If it starts with /, it's a relative path - prepend API base
	if (assetPath.startsWith("/")) {
		const apiUrl = import.meta.env.VITE_API_URL || "";
		return `${apiUrl}${assetPath}`;
	}

	// Default: prepend /public/
	const apiUrl = import.meta.env.VITE_API_URL || "";
	return `${apiUrl}/${assetPath}`;
};

/**
 * Validate file before upload
 * @param file - File to validate
 * @returns Error message if invalid, null if valid
 */
export const validateAssetFile = (file: File): string | null => {
	// Check file size
	if (file.size > MAX_FILE_SIZE) {
		return `File size must be less than ${MAX_FILE_SIZE / 1024 / 1024}MB`;
	}

	// Check MIME type
	if (!ALLOWED_MIME_TYPES[file.type as keyof typeof ALLOWED_MIME_TYPES]) {
		return "File type not allowed. Allowed types: images, PDFs, documents, and videos";
	}

	return null;
};

/**
 * Upload an asset file (admin only)
 * @param file - File to upload
 * @returns Upload response with URL and ID
 */
export const uploadAsset = async (file: File): Promise<UploadAssetResponse> => {
	// Validate file
	const validationError = validateAssetFile(file);
	if (validationError) {
		throw new Error(validationError);
	}

	const formData = new FormData();
	formData.append("file", file);

	const response = await apiClient.post("/admin/assets/upload", formData, {
		headers: {
			"Content-Type": "multipart/form-data",
		},
	});

	return response.data;
};

/**
 * Delete an asset by ID (admin only)
 * @param id - Asset ID to delete
 */
export const deleteAsset = async (id: string): Promise<void> => {
	await apiClient.delete(`/admin/assets/${id}`);
};

/**
 * Get an asset file by filename (public endpoint)
 * @param filename - Asset filename
 * @returns Asset file blob
 */
export const getAsset = async (filename: string): Promise<Blob> => {
	const response = await apiClient.get(`/assets/${filename}`, {
		responseType: "blob",
	});
	return response.data;
};

/**
 * Get asset URL for embedding
 * @param filename - Asset filename
 * @returns Absolute URL for the asset
 */
export const getAssetFileUrl = (filename: string): string => {
	const apiUrl = import.meta.env.VITE_API_URL || "";
	return `${apiUrl}/api/v1/assets/${filename}`;
};

/**
 * Upload multiple asset files
 * @param files - Array of files to upload
 * @returns Array of upload responses
 */
export const uploadAssets = async (
	files: File[],
): Promise<UploadAssetResponse[]> => {
	const results: UploadAssetResponse[] = [];

	for (const file of files) {
		try {
			const result = await uploadAsset(file);
			results.push(result);
		} catch (error) {
			console.error(`Failed to upload ${file.name}:`, error);
			// Continue with next file even if one fails
		}
	}

	return results;
};

/**
 * Check if a file is an image
 * @param mimeType - MIME type of the file
 * @returns True if file is an image
 */
export const isImageType = (mimeType: string): boolean => {
	return mimeType.startsWith("image/");
};

/**
 * Check if a file is a video
 * @param mimeType - MIME type of the file
 * @returns True if file is a video
 */
export const isVideoType = (mimeType: string): boolean => {
	return mimeType.startsWith("video/");
};

/**
 * Check if a file is a document
 * @param mimeType - MIME type of the file
 * @returns True if file is a document
 */
export const isDocumentType = (mimeType: string): boolean => {
	const documentTypes = [
		"application/pdf",
		"text/plain",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	];
	return documentTypes.includes(mimeType);
};
