import { useState } from "react";
import { Trash2, Download, ExternalLink, Loader2, FileText, Music, Video } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useAssetDelete, useAssetValidation } from "@/hooks/useAssets";
import { getAssetUrl } from "@/lib/api/assets";
import { toast } from "sonner";
import type { Asset } from "@/lib/api/assets";

interface AssetGalleryProps {
  assets: Asset[];
  onAssetDeleted?: (assetId: string) => void;
  isLoading?: boolean;
  className?: string;
  displayMode?: "grid" | "list";
}

function getAssetIcon(mimeType: string) {
  if (mimeType.startsWith("image/")) {
    return null; // Will show image
  }
  if (mimeType.startsWith("video/")) {
    return <Video className="h-8 w-8 text-blue-500" />;
  }
  if (
    mimeType === "application/pdf" ||
    mimeType === "application/msword" ||
    mimeType.includes("wordprocessingml") ||
    mimeType.includes("spreadsheet")
  ) {
    return <FileText className="h-8 w-8 text-orange-500" />;
  }
  if (mimeType.startsWith("audio/")) {
    return <Music className="h-8 w-8 text-purple-500" />;
  }
  return <FileText className="h-8 w-8 text-gray-500" />;
}

function AssetGridItem({
  asset,
  onDelete,
  isDeletingId,
}: {
  asset: Asset;
  onDelete: (id: string) => void;
  isDeletingId: string | null;
}) {
  const { isImageType } = useAssetValidation();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const assetUrl = getAssetUrl(asset.url);

  return (
    <>
      <div className="group relative overflow-hidden rounded-lg border border-border bg-muted/50 hover:border-border/80 transition-colors">
        {/* Asset Preview */}
        <div className="aspect-square bg-muted flex items-center justify-center overflow-hidden">
          {isImageType(asset.mime_type) ? (
            <img
              src={assetUrl}
              alt={asset.file_name}
              className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-200"
              onError={(e) => {
                e.currentTarget.src = assetUrl;
              }}
            />
          ) : (
            <div className="flex flex-col items-center justify-center gap-2 p-4 text-center">
              {getAssetIcon(asset.mime_type)}
              <p className="text-xs text-muted-foreground line-clamp-2">{asset.file_name}</p>
            </div>
          )}
        </div>

        {/* Overlay Actions */}
        <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
          <Button
            size="icon"
            variant="ghost"
            className="bg-white/20 hover:bg-white/30 text-white"
            onClick={() => window.open(assetUrl, "_blank")}
            title="Open in new tab"
          >
            <ExternalLink className="h-4 w-4" />
          </Button>
          <Button
            size="icon"
            variant="ghost"
            className="bg-white/20 hover:bg-white/30 text-white"
            onClick={() => {
              const link = document.createElement("a");
              link.href = assetUrl;
              link.download = asset.file_name;
              link.click();
            }}
            title="Download file"
          >
            <Download className="h-4 w-4" />
          </Button>
          <Button
            size="icon"
            variant="ghost"
            className="bg-white/20 hover:bg-destructive/80 text-white"
            onClick={() => setShowDeleteDialog(true)}
            disabled={isDeletingId === asset.id}
            title="Delete file"
          >
            {isDeletingId === asset.id ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Trash2 className="h-4 w-4" />
            )}
          </Button>
        </div>

        {/* Asset Info */}
        <div className="p-2 space-y-1">
          <p className="text-xs font-medium line-clamp-1" title={asset.file_name}>
            {asset.file_name}
          </p>
          <p className="text-xs text-muted-foreground">
            {new Date(asset.created_at).toLocaleDateString()}
          </p>
        </div>
      </div>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Asset</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete <strong>{asset.file_name}</strong>? This action
              cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className="flex gap-2 justify-end">
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              onClick={() => {
                onDelete(asset.id);
                setShowDeleteDialog(false);
              }}
              disabled={isDeletingId === asset.id}
            >
              {isDeletingId === asset.id && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Delete
            </AlertDialogAction>
          </div>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

function AssetListItem({
  asset,
  onDelete,
  isDeletingId,
}: {
  asset: Asset;
  onDelete: (id: string) => void;
  isDeletingId: string | null;
}) {
  const { isImageType } = useAssetValidation();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const assetUrl = getAssetUrl(asset.url);

  return (
    <>
      <div className="flex items-center justify-between gap-4 rounded-lg border border-border bg-muted/50 p-4 hover:bg-muted transition-colors">
        <div className="flex items-center gap-3 flex-1 min-w-0">
          {isImageType(asset.mime_type) ? (
            <img
              src={assetUrl}
              alt={asset.file_name}
              className="h-12 w-12 rounded object-cover shrink-0"
              onError={(e) => {
                e.currentTarget.style.display = "none";
              }}
            />
          ) : (
            <div className="h-12 w-12 rounded bg-muted flex items-center justify-center shrink-0">
              {getAssetIcon(asset.mime_type)}
            </div>
          )}
          <div className="min-w-0">
            <p className="text-sm font-medium truncate" title={asset.file_name}>
              {asset.file_name}
            </p>
            <p className="text-xs text-muted-foreground">
              {new Date(asset.created_at).toLocaleDateString()} •{" "}
              {asset.mime_type}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2 shrink-0">
          <Button
            size="sm"
            variant="ghost"
            onClick={() => window.open(assetUrl, "_blank")}
            title="Open in new tab"
          >
            <ExternalLink className="h-4 w-4" />
          </Button>
          <Button
            size="sm"
            variant="ghost"
            onClick={() => {
              const link = document.createElement("a");
              link.href = assetUrl;
              link.download = asset.file_name;
              link.click();
            }}
            title="Download file"
          >
            <Download className="h-4 w-4" />
          </Button>
          <Button
            size="sm"
            variant="ghost"
            className="text-destructive hover:text-destructive hover:bg-destructive/10"
            onClick={() => setShowDeleteDialog(true)}
            disabled={isDeletingId === asset.id}
            title="Delete file"
          >
            {isDeletingId === asset.id ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Trash2 className="h-4 w-4" />
            )}
          </Button>
        </div>
      </div>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Asset</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete <strong>{asset.file_name}</strong>? This action
              cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className="flex gap-2 justify-end">
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              onClick={() => {
                onDelete(asset.id);
                setShowDeleteDialog(false);
              }}
              disabled={isDeletingId === asset.id}
            >
              {isDeletingId === asset.id && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Delete
            </AlertDialogAction>
          </div>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

export function AssetGallery({
  assets,
  onAssetDeleted,
  isLoading = false,
  className = "",
  displayMode = "grid",
}: AssetGalleryProps) {
  const { deleteAsset } = useAssetDelete();
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const handleDelete = async (assetId: string) => {
    try {
      setDeletingId(assetId);
      await deleteAsset(assetId);
      toast.success("Asset deleted successfully");
      onAssetDeleted?.(assetId);
    } catch (err) {
      toast.error("Failed to delete asset");
    } finally {
      setDeletingId(null);
    }
  };

  if (isLoading) {
    return (
      <div className={`flex items-center justify-center p-8 ${className}`}>
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (assets.length === 0) {
    return (
      <div className={`text-center p-8 rounded-lg border border-dashed border-muted-foreground/30 ${className}`}>
        <p className="text-muted-foreground">No assets uploaded yet</p>
      </div>
    );
  }

  return (
    <div className={className}>
      {displayMode === "grid" ? (
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
          {assets.map((asset) => (
            <AssetGridItem
              key={asset.id}
              asset={asset}
              onDelete={handleDelete}
              isDeletingId={deletingId}
            />
          ))}
        </div>
      ) : (
        <div className="space-y-2">
          {assets.map((asset) => (
            <AssetListItem
              key={asset.id}
              asset={asset}
              onDelete={handleDelete}
              isDeletingId={deletingId}
            />
          ))}
        </div>
      )}
    </div>
  );
}
