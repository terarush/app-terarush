import React, { createContext, useContext, useState, useCallback } from "react";
import type { ReactNode } from "react";
import {
  uploadAsset,
  deleteAsset,
  uploadAssets,
  type UploadAssetResponse,
} from "@/lib/api/assets";

interface AssetContextType {
  isUploading: boolean;
  isDeleting: boolean;
  uploadProgress: number;
  error: string | null;
  uploadFile: (file: File) => Promise<UploadAssetResponse>;
  uploadMultipleFiles: (files: File[]) => Promise<UploadAssetResponse[]>;
  removeAsset: (id: string) => Promise<void>;
  clearError: () => void;
}

const AssetContext = createContext<AssetContextType | undefined>(undefined);

interface AssetProviderProps {
  children: ReactNode;
}

export const AssetProvider: React.FC<AssetProviderProps> = ({ children }) => {
  const [isUploading, setIsUploading] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);

  const uploadFile = useCallback(async (file: File): Promise<UploadAssetResponse> => {
    setIsUploading(true);
    setError(null);
    setUploadProgress(0);

    try {
      // Simulate upload progress
      const progressInterval = setInterval(() => {
        setUploadProgress((prev) => {
          if (prev >= 90) {
            clearInterval(progressInterval);
            return prev;
          }
          return prev + Math.random() * 30;
        });
      }, 200);

      const result = await uploadAsset(file);

      clearInterval(progressInterval);
      setUploadProgress(100);

      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to upload file";
      setError(errorMessage);
      throw err;
    } finally {
      setIsUploading(false);
      // Reset progress after a delay
      setTimeout(() => setUploadProgress(0), 1000);
    }
  }, []);

  const uploadMultipleFiles = useCallback(
    async (files: File[]): Promise<UploadAssetResponse[]> => {
      setIsUploading(true);
      setError(null);
      setUploadProgress(0);

      try {
        const results = await uploadAssets(files);

        if (results.length === 0) {
          setError("Failed to upload any files. Please check file types and sizes.");
        }

        setUploadProgress(100);
        return results;
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : "Failed to upload files";
        setError(errorMessage);
        throw err;
      } finally {
        setIsUploading(false);
        // Reset progress after a delay
        setTimeout(() => setUploadProgress(0), 1000);
      }
    },
    [],
  );

  const removeAsset = useCallback(async (id: string): Promise<void> => {
    setIsDeleting(true);
    setError(null);

    try {
      await deleteAsset(id);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to delete asset";
      setError(errorMessage);
      throw err;
    } finally {
      setIsDeleting(false);
    }
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  const value: AssetContextType = {
    isUploading,
    isDeleting,
    uploadProgress,
    error,
    uploadFile,
    uploadMultipleFiles,
    removeAsset,
    clearError,
  };

  return (
    <AssetContext.Provider value={value}>
      {children}
    </AssetContext.Provider>
  );
};

export const useAssetsContext = () => {
  const context = useContext(AssetContext);
  if (context === undefined) {
    throw new Error("useAssetsContext must be used within an AssetProvider");
  }
  return context;
};
