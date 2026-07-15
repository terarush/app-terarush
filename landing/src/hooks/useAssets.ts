import { useAssetsContext } from "@/contexts/AssetContext";
import { validateAssetFile, isImageType, isVideoType, isDocumentType } from "@/lib/api/assets";

// Re-export context hook
export { useAssetsContext } from "@/contexts/AssetContext";

/**
 * Hook to upload a single asset file
 * @returns Object with upload function and upload state
 */
export const useAssetUpload = () => {
  const { uploadFile, isUploading, uploadProgress, error, clearError } =
    useAssetsContext();

  return {
    uploadFile,
    isUploading,
    uploadProgress,
    error,
    clearError,
  };
};

/**
 * Hook to upload multiple asset files
 * @returns Object with upload function and upload state
 */
export const useAssetsUpload = () => {
  const { uploadMultipleFiles, isUploading, uploadProgress, error, clearError } =
    useAssetsContext();

  return {
    uploadMultipleFiles,
    isUploading,
    uploadProgress,
    error,
    clearError,
  };
};

/**
 * Hook to delete an asset
 * @returns Object with delete function and delete state
 */
export const useAssetDelete = () => {
  const { removeAsset, isDeleting, error, clearError } = useAssetsContext();

  return {
    deleteAsset: removeAsset,
    isDeleting,
    error,
    clearError,
  };
};

/**
 * Hook to validate a file before upload
 * @returns Object with validation function
 */
export const useAssetValidation = () => {
  return {
    validateFile: validateAssetFile,
    isImageType,
    isVideoType,
    isDocumentType,
  };
};

/**
 * Hook to get all asset hooks
 * Combines upload, delete, and validation hooks
 * @returns Object with all asset operations
 */
export const useAssets = () => {
  const upload = useAssetUpload();
  const multiUpload = useAssetsUpload();
  const deletion = useAssetDelete();
  const validation = useAssetValidation();

  return {
    upload,
    multiUpload,
    deletion,
    validation,
  };
};
