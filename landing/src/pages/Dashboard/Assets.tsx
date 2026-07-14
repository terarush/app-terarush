import { useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { AssetUploader } from "@/components/AssetUploader";
import { AssetGallery } from "@/components/AssetGallery";
import { apiClient } from "@/lib/api/client";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { Asset } from "@/lib/api/assets";

export function Assets() {
	const queryClient = useQueryClient();
	const [displayMode, setDisplayMode] = useState<"grid" | "list">("grid");
	const [uploadTab, setUploadTab] = useState("single");

	const {
		data: assets = [],
		isLoading: loading,
		refetch,
	} = useQuery<Asset[]>({
		queryKey: ["assets"],
		queryFn: async () => {
			const response = await apiClient.get("/admin/assets");
			return response.data?.assets || [];
		},
	});

	const handleUploadSuccess = () => {
		queryClient.invalidateQueries({ queryKey: ["assets"] });
	};

	const handleAssetDeleted = () => {
		queryClient.invalidateQueries({ queryKey: ["assets"] });
	};

	return (
		<div className="space-y-6">
			{/* Header */}
			<div className="flex items-center justify-between">
				<div>
					<h1 className="text-3xl font-bold tracking-tight">
						Assets
					</h1>
					<p className="text-muted-foreground mt-2">
						Upload and manage your media files, documents, and other
						assets
					</p>
				</div>
			</div>

			<Tabs defaultValue="upload" className="space-y-4">
				<TabsList>
					<TabsTrigger value="upload">Upload</TabsTrigger>
					<TabsTrigger value="gallery">
						Gallery {assets.length > 0 && `(${assets.length})`}
					</TabsTrigger>
				</TabsList>

				{/* Upload Tab */}
				<TabsContent value="upload" className="space-y-4">
					<div className="rounded-lg border border-border bg-card p-6">
						<Tabs
							value={uploadTab}
							onValueChange={setUploadTab}
							className="space-y-4"
						>
							<TabsList>
								<TabsTrigger value="single">
									Single File
								</TabsTrigger>
								<TabsTrigger value="multiple">
									Multiple Files
								</TabsTrigger>
							</TabsList>

							{/* Single File Upload */}
							<TabsContent value="single" className="space-y-4">
								<div className="space-y-2">
									<h3 className="font-semibold">
										Upload Single Asset
									</h3>
									<p className="text-sm text-muted-foreground">
										Upload one file at a time. Maximum file
										size: 10MB
									</p>
								</div>
								<AssetUploader
									multiple={false}
									onUploadSuccess={handleUploadSuccess}
									acceptedTypes="image/*,.pdf,.doc,.docx,.xls,.xlsx,.txt,.mp4,.webm"
								/>
							</TabsContent>

							{/* Multiple Files Upload */}
							<TabsContent value="multiple" className="space-y-4">
								<div className="space-y-2">
									<h3 className="font-semibold">
										Upload Multiple Assets
									</h3>
									<p className="text-sm text-muted-foreground">
										Upload up to 10 files at once. Maximum
										10MB per file
									</p>
								</div>
								<AssetUploader
									multiple={true}
									onUploadSuccess={handleUploadSuccess}
									acceptedTypes="image/*,.pdf,.doc,.docx,.xls,.xlsx,.txt,.mp4,.webm"
								/>
							</TabsContent>
						</Tabs>
					</div>

					{/* Supported Formats */}
					<div className="rounded-lg border border-border bg-card p-4 space-y-3">
						<h4 className="font-semibold text-sm">
							Supported Formats
						</h4>
						<div className="grid grid-cols-2 gap-3 text-sm">
							<div>
								<p className="font-medium text-foreground">
									Images
								</p>
								<p className="text-muted-foreground">
									JPG, PNG, GIF, WebP, SVG
								</p>
							</div>
							<div>
								<p className="font-medium text-foreground">
									Documents
								</p>
								<p className="text-muted-foreground">
									PDF, TXT, DOC, DOCX
								</p>
							</div>
							<div>
								<p className="font-medium text-foreground">
									Spreadsheets
								</p>
								<p className="text-muted-foreground">
									XLS, XLSX
								</p>
							</div>
							<div>
								<p className="font-medium text-foreground">
									Videos
								</p>
								<p className="text-muted-foreground">
									MP4, WebM
								</p>
							</div>
						</div>
					</div>
				</TabsContent>

				{/* Gallery Tab */}
				<TabsContent value="gallery" className="space-y-4">
					<div className="flex items-center justify-between gap-4">
						<div>
							<h3 className="font-semibold">Asset Gallery</h3>
							<p className="text-sm text-muted-foreground">
								{assets.length} asset
								{assets.length !== 1 ? "s" : ""} uploaded
							</p>
						</div>
						<div className="flex gap-2">
							<Button
								variant={
									displayMode === "grid"
										? "default"
										: "outline"
								}
								size="sm"
								onClick={() => setDisplayMode("grid")}
							>
								Grid
							</Button>
							<Button
								variant={
									displayMode === "list"
										? "default"
										: "outline"
								}
								size="sm"
								onClick={() => setDisplayMode("list")}
							>
								List
							</Button>
							<Button
								variant="outline"
								size="sm"
								onClick={() => refetch()}
								disabled={loading}
							>
								{loading ? "Loading..." : "Refresh"}
							</Button>
						</div>
					</div>

					{loading ? (
						<div className="flex items-center justify-center py-12 rounded-lg border border-border bg-card">
							<p className="text-muted-foreground">
								Loading assets...
							</p>
						</div>
					) : (
						<AssetGallery
							assets={assets}
							onAssetDeleted={handleAssetDeleted}
							displayMode={displayMode}
							className="rounded-lg border border-border bg-card p-4"
						/>
					)}
				</TabsContent>
			</Tabs>
		</div>
	);
}
