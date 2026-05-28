import { useCallback, useRef, useState } from "react";
import { Upload, X, AlertCircle, CheckCircle2, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useAssetUpload, useAssetsUpload, useAssetValidation } from "@/hooks/useAssets";
import { toast } from "sonner";
import type { UploadAssetResponse } from "@/lib/api/assets";

interface AssetUploaderProps {
  onUploadSuccess?: () => void;
  acceptedTypes?: string;
  multiple?: boolean;
  className?: string;
}

export function AssetUploader({
  onUploadSuccess,
  acceptedTypes = "image/*,.pdf,.doc,.docx,.xls,.xlsx,.txt,.mp4,.webm",
  multiple = false,
  className = "",
}: AssetUploaderProps) {
  const singleUpload = useAssetUpload();
  const multiUpload = useAssetsUpload();
  const { validateFile } = useAssetValidation();
  const inputRef = useRef<HTMLInputElement>(null);
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [uploadErrors, setUploadErrors] = useState<Record<string, string>>({});

  const { isUploading, uploadProgress, error, clearError } = multiple
    ? multiUpload
    : singleUpload;

  const handleFileSelect = useCallback(
    (files: FileList | null) => {
      if (!files) return;

      const newFiles = Array.from(files);
      const errors: Record<string, string> = {};
      const validFiles: File[] = [];

      newFiles.forEach((file) => {
        const validationError = validateFile(file);
        if (validationError) {
          errors[file.name] = validationError;
        } else {
          validFiles.push(file);
        }
      });

      setUploadErrors(errors);

      if (Object.keys(errors).length > 0) {
        toast.error("Some files failed validation");
      }

      if (!multiple && validFiles.length > 0) {
        // For single file upload, only keep the first valid file
        setSelectedFiles([validFiles[0]]);
      } else {
        setSelectedFiles(validFiles);
      }
    },
    [validateFile, multiple],
  );

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleFileSelect(e.target.files);
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    handleFileSelect(e.dataTransfer.files);
  };

  const handleUpload = async () => {
    if (selectedFiles.length === 0) return;

    try {
      clearError();

      if (!multiple && selectedFiles.length === 1) {
        const file = selectedFiles[0];
        await singleUpload.uploadFile(file);
        toast.success("File uploaded successfully");
        onUploadSuccess?.();
        setSelectedFiles([]);
        if (inputRef.current) {
          inputRef.current.value = "";
        }
      } else if (multiple && selectedFiles.length > 0) {
        const results: UploadAssetResponse[] = await multiUpload.uploadMultipleFiles(selectedFiles);
        if (results.length > 0) {
          toast.success(`Successfully uploaded ${results.length} file(s)`);
          onUploadSuccess?.();
          setSelectedFiles([]);
          if (inputRef.current) {
            inputRef.current.value = "";
          }
        }
      }
    } catch (err) {
      toast.error(error || "Failed to upload file(s)");
    }
  };

  const removeFile = (index: number) => {
    setSelectedFiles((prev) => prev.filter((_, i) => i !== index));
  };

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Error Alert */}
      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription className="flex items-center justify-between">
            <span>{error}</span>
            <button
              onClick={clearError}
              className="text-sm font-medium underline hover:no-underline"
            >
              Dismiss
            </button>
          </AlertDescription>
        </Alert>
      )}

      {/* Validation Errors */}
      {Object.entries(uploadErrors).length > 0 && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            <div className="space-y-1">
              <p className="font-medium">Validation errors:</p>
              <ul className="ml-2 list-disc space-y-1">
                {Object.entries(uploadErrors).map(([filename, error]) => (
                  <li key={filename} className="text-sm">
                    <strong>{filename}:</strong> {error}
                  </li>
                ))}
              </ul>
            </div>
          </AlertDescription>
        </Alert>
      )}

      {/* Drop Zone */}
      <div
        onDragOver={handleDragOver}
        onDrop={handleDrop}
        className="relative border-2 border-dashed border-muted-foreground/30 rounded-lg p-6 text-center hover:border-muted-foreground/50 transition-colors cursor-pointer"
      >
        <Input
          ref={inputRef}
          type="file"
          multiple={multiple}
          accept={acceptedTypes}
          onChange={handleInputChange}
          disabled={isUploading}
          className="hidden"
        />

        <div
          onClick={() => inputRef.current?.click()}
          className="flex flex-col items-center justify-center gap-2"
        >
          <Upload className="h-8 w-8 text-muted-foreground" />
          <div className="space-y-1">
            <p className="text-sm font-medium text-foreground">
              Drop files here or click to browse
            </p>
            <p className="text-xs text-muted-foreground">
              {multiple ? "Upload multiple files" : "Upload a single file"} (Max 10MB each)
            </p>
          </div>
        </div>
      </div>

      {/* Selected Files List */}
      {selectedFiles.length > 0 && (
        <div className="space-y-2">
          <p className="text-sm font-medium">Selected files:</p>
          <div className="space-y-2">
            {selectedFiles.map((file, index) => (
              <div
                key={`${file.name}-${index}`}
                className="flex items-center justify-between rounded-lg border border-border bg-muted/50 p-3"
              >
                <div className="flex items-center gap-2 flex-1 min-w-0">
                  {isUploading && uploadProgress > 0 && (
                    <Loader2 className="h-4 w-4 animate-spin text-primary shrink-0" />
                  )}
                  {!isUploading && <CheckCircle2 className="h-4 w-4 text-green-500 shrink-0" />}
                  <div className="min-w-0">
                    <p className="text-sm font-medium truncate">{file.name}</p>
                    <p className="text-xs text-muted-foreground">
                      {(file.size / 1024 / 1024).toFixed(2)}MB
                    </p>
                  </div>
                </div>
                {!isUploading && (
                  <button
                    onClick={() => removeFile(index)}
                    className="p-1 hover:bg-destructive/10 rounded-md transition-colors"
                    title="Remove file"
                  >
                    <X className="h-4 w-4 text-destructive" />
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Upload Progress */}
      {isUploading && uploadProgress > 0 && (
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <p className="text-sm font-medium">Uploading...</p>
            <p className="text-sm text-muted-foreground">{Math.round(uploadProgress)}%</p>
          </div>
          <div className="w-full bg-secondary rounded-full h-2">
            <div
              className="bg-primary h-2 rounded-full transition-all duration-300"
              style={{ width: `${uploadProgress}%` }}
            />
          </div>
        </div>
      )}

      {/* Action Buttons */}
      {selectedFiles.length > 0 && (
        <div className="flex gap-2 justify-end">
          <Button
            variant="outline"
            onClick={() => {
              setSelectedFiles([]);
              if (inputRef.current) {
                inputRef.current.value = "";
              }
            }}
            disabled={isUploading}
          >
            Clear
          </Button>
          <Button onClick={handleUpload} disabled={isUploading || selectedFiles.length === 0}>
            {isUploading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {isUploading ? "Uploading..." : `Upload ${multiple ? "Files" : "File"}`}
          </Button>
        </div>
      )}
    </div>
  );
}
