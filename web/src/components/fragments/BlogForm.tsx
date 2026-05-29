import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAuth } from "@/contexts/AuthContext";
import {
	createBlog,
	updateBlog,
	uploadBlogImage,
	getImageUrl,
	type Blog,
	type CreateBlogRequest,
} from "@/lib/api/blogs";
import { createBlogSchema, type CreateBlogFormData } from "@/validations/blog.validation";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/components/ui/form";
import { Badge } from "@/components/ui/badge";
import { Loader2, Upload, X, FileText, Image as ImageIcon, Zap } from "lucide-react";
import { toast } from "sonner";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { MarkdownEditor } from "./MarkdownEditor/index";
import { Textarea } from "../ui/textarea";

interface BlogFormProps {
	blog?: Blog;
	onClose: (saved?: boolean) => void;
}

export function BlogForm({ blog, onClose }: BlogFormProps) {
	const { user } = useAuth();
	const [loading, setLoading] = useState(false);
	const [uploading, setUploading] = useState(false);
	const [previewUrl, setPreviewUrl] = useState<string>(blog?.image ? getImageUrl(blog.image) : "");
	const [imageError, setImageError] = useState(false);
	const form = useForm<CreateBlogFormData>({
		resolver: zodResolver(createBlogSchema),
		defaultValues: blog
			? {
				title: blog.title,
				slug: blog.slug,
				content: blog.content,
				excerpt: blog.excerpt,
				category: blog.category,
				tags: blog.tags,
				image: blog.image || "",
				is_published: blog.is_published,
			}
			: {
				title: "",
				slug: "",
				content: "",
				excerpt: "",
				category: "",
				tags: "",
				image: "",
				is_published: false,
			},
	});

	const onSubmit = async (data: CreateBlogRequest) => {
		try {
			setLoading(true);
			if (blog?.id) {
				await updateBlog(blog.id, data);
				toast.success("Blog post updated successfully");
			} else {
				await createBlog(data);
				toast.success("Blog post created successfully");
			}
			onClose(true);
		} catch (err) {
			console.error("Error saving blog:", err);
			toast.error("Failed to save blog post");
		} finally {
			setLoading(false);
		}
	};

	const generateSlug = (title: string) => {
		return title
			.toLowerCase()
			.replace(/[^\w\s-]/g, "")
			.replace(/\s+/g, "-")
			.replace(/-+/g, "-");
	};

	const handleTitleChange = (title: string) => {
		form.setValue("title", title);
		if (!blog) {
			form.setValue("slug", generateSlug(title));
		}
	};

	const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
		const file = e.target.files?.[0];
		if (!file) return;

		try {
			setUploading(true);
			const response = await uploadBlogImage(file);
			const imageUrl = response.url;

			form.setValue("image", imageUrl);
			setPreviewUrl(getImageUrl(imageUrl));
			setImageError(false);
			toast.success("Image uploaded successfully");
		} catch (error) {
			console.error("Error uploading image:", error);
			toast.error("Failed to upload image");
		} finally {
			setUploading(false);
		}
	};

	const handleUrlInput = (url: string) => {
		form.setValue("image", url);
		setPreviewUrl(getImageUrl(url));
		setImageError(false);
	};

	const clearImage = () => {
		form.setValue("image", "");
		setPreviewUrl("");
		setImageError(false);
	};

	return (
		<Dialog open={true} onOpenChange={(open) => !open && onClose()}>
			<DialogContent className="max-w-6xl w-full p-0 flex flex-col max-h-[95vh] rounded-2xl overflow-hidden bg-background border-border/50 shadow-2xl">
				{/* Header */}
				<div className="p-8 border-b border-border/50 bg-gradient-to-r from-background via-background to-primary/5">
					<DialogHeader>
						<div className="flex items-center gap-3 mb-2">
							<div className="p-2 rounded-lg bg-primary/10">
								<FileText className="h-5 w-5 text-primary" />
							</div>
							<DialogTitle className="text-2xl font-bold">
								{blog ? "Edit Blog Post" : "Create New Blog Post"}
							</DialogTitle>
						</div>
						<DialogDescription className="text-base mt-2">
							{blog
								? "Update your blog post details and content"
								: "Fill in the form to create a new blog post"}
						</DialogDescription>
					</DialogHeader>
				</div>

				{/* Form Content */}
				<div className="flex-1 overflow-y-auto p-8">
					<Form {...form}>
						<form
							onSubmit={form.handleSubmit(onSubmit)}
							id="blog-form"
							className="space-y-8"
						>
							{/* Title Section */}
							<div className="space-y-4">
								<div className="flex items-center gap-2">
									<h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">Basic Info</h3>
									<div className="flex-1 h-px bg-border/50" />
								</div>
								<FormField
									control={form.control}
									name="title"
									render={({ field }) => (
										<FormItem>
											<FormLabel className="text-base font-semibold">Title</FormLabel>
											<FormControl>
												<Input
													placeholder="Enter an engaging blog post title"
													className="h-11 bg-muted/50 border-border/50 rounded-lg"
													{...field}
													onChange={(e) => {
														field.onChange(e);
														handleTitleChange(e.target.value);
													}}
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>

								<FormField
									control={form.control}
									name="slug"
									render={({ field }) => (
										<FormItem>
											<FormLabel className="text-base font-semibold">Slug (URL)</FormLabel>
											<FormControl>
												<Input
													placeholder="auto-generated-from-title"
													className="h-11 bg-muted/50 border-border/50 rounded-lg"
													{...field}
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>

							{/* Author & Category */}
							<div className="space-y-4">
								<div className="flex items-center gap-2">
									<h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">Metadata</h3>
									<div className="flex-1 h-px bg-border/50" />
								</div>
								<div className="grid grid-cols-2 gap-4 flex items-center align-center">
									<FormItem>
										<FormLabel className="text-base font-semibold">Author</FormLabel>
										<div className="mt-2 px-4 py-2.5 bg-primary/5 rounded-lg border border-primary/20 text-sm font-medium flex items-center gap-2 h-11">
											<Zap className="h-4 w-4 text-primary" />
											{user?.name || "Not authenticated"}
										</div>
									</FormItem>

									<FormField
										control={form.control}
										name="category"
										render={({ field }) => (
											<FormItem>
												<FormLabel className="text-base font-semibold">Category</FormLabel>
												<FormControl>
													<Input
														placeholder="e.g. Web Development"
														className="h-11 bg-muted/50 border-border/50 rounded-lg mt-2 py-4"
														{...field}
													/>
												</FormControl>
												<FormMessage />
											</FormItem>
										)}
									/>
								</div>
							</div>

							{/* Excerpt */}
							<div className="space-y-4">
								<FormField
									control={form.control}
									name="excerpt"
									render={({ field }) => (
										<FormItem>
											<FormLabel className="text-base font-semibold">Excerpt</FormLabel>
											<FormControl>
												<Textarea
													placeholder="A brief summary that will appear in blog listings"
													className="rounded-lg bg-muted/50 border-border/50"
													{...field}
													rows={3}
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>

							{/* Tags */}
							<div className="space-y-4">
								<FormField
									control={form.control}
									name="tags"
									render={({ field }) => (
										<FormItem>
											<FormLabel className="text-base font-semibold">Tags</FormLabel>
											<FormControl>
												<Input
													placeholder="Comma-separated (e.g. react, javascript, web)"
													className="h-11 bg-muted/50 border-border/50 rounded-lg"
													{...field}
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>

							{/* Content Editor */}
							<div className="space-y-4">
								<div className="flex items-center gap-2">
									<h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">Content</h3>
									<div className="flex-1 h-px bg-border/50" />
								</div>
								<FormField
									control={form.control}
									name="content"
									render={({ field }) => (
										<FormItem>
											<FormControl>
												<MarkdownEditor
													value={field.value}
													onChange={field.onChange}
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>

							{/* Featured Image */}
							<div className="space-y-4">
								<div className="flex items-center gap-2">
									<h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">Featured Image</h3>
									<div className="flex-1 h-px bg-border/50" />
								</div>
								<FormField
									control={form.control}
									name="image"
									render={({ field }) => (
										<FormItem>
											<Tabs
												defaultValue="url"
												className="w-full"
											>
												<TabsList className="grid w-full grid-cols-2 bg-muted/50 p-1">
													<TabsTrigger value="url" className="rounded-md data-[state=active]:bg-background">
														<ImageIcon className="h-4 w-4 mr-2" />
														From URL
													</TabsTrigger>
													<TabsTrigger value="upload" className="rounded-md data-[state=active]:bg-background">
														<Upload className="h-4 w-4 mr-2" />
														Upload
													</TabsTrigger>
												</TabsList>
												<TabsContent
													value="url"
													className="space-y-3 mt-4"
												>
													<FormControl>
														<Input
															placeholder="https://example.com/image.jpg"
															className="h-11 bg-muted/50 border-border/50 rounded-lg"
															value={field.value}
															onChange={(e) => {
																handleUrlInput(
																	e.target.value
																);
															}}
														/>
													</FormControl>
												</TabsContent>
												<TabsContent
													value="upload"
													className="space-y-3 mt-4"
												>
													<div className="relative">
														<label className="flex flex-col items-center justify-center w-full border-2 border-dashed border-border/50 rounded-xl cursor-pointer hover:border-primary/50 hover:bg-primary/5 transition-all p-8 group">
															<div className="flex flex-col items-center justify-center text-center">
																<div className="p-3 rounded-lg bg-primary/10 group-hover:bg-primary/20 mb-3 transition-colors">
																	<Upload className="h-6 w-6 text-primary group-hover:text-primary/80" />
																</div>
																<p className="text-sm font-semibold text-foreground mb-1">
																	{uploading
																		? "Uploading..."
																		: "Click to upload or drag and drop"}
																</p>
																<p className="text-xs text-muted-foreground">
																	PNG, JPG, GIF, WebP up to 5MB
																</p>
															</div>
															<input
																type="file"
																className="hidden"
																accept="image/*"
																onChange={handleImageUpload}
																disabled={uploading}
															/>
														</label>
													</div>
												</TabsContent>
											</Tabs>
											{previewUrl && !imageError && (
												<div className="relative mt-6 rounded-xl overflow-hidden bg-gradient-to-br from-primary/10 via-accent/5 to-background p-1 border border-border/50">
													<div className="relative aspect-video bg-muted rounded-lg overflow-hidden">
														<img
															src={previewUrl}
															alt="Preview"
															className="w-full h-full object-cover"
															onError={() => setImageError(true)}
															onLoad={() => setImageError(false)}
														/>
														<button
															type="button"
															onClick={clearImage}
															className="absolute top-3 right-3 bg-destructive hover:bg-destructive/90 text-white p-2 rounded-lg transition-all shadow-lg hover:shadow-xl"
															title="Remove image"
														>
															<X className="h-4 w-4" />
														</button>
													</div>
												</div>
											)}
											{imageError && previewUrl && (
												<div className="mt-4 p-4 bg-destructive/10 border border-destructive/30 rounded-lg text-destructive text-sm">
													<p className="font-medium mb-2">Failed to load image</p>
													<p className="text-xs mb-3">Please check the URL and try again.</p>
													<button
														type="button"
														onClick={clearImage}
														className="text-destructive hover:text-destructive/80 underline text-xs font-medium"
													>
														Clear
													</button>
												</div>
											)}
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>

							{/* Publish Status */}
							<div className="space-y-4 p-4 bg-gradient-to-r from-primary/5 via-accent/5 to-background rounded-xl border border-border/50">
								<FormField
									control={form.control}
									name="is_published"
									render={({ field }) => (
										<FormItem className="flex items-center justify-between">
											<div className="flex-1">
												<FormLabel className="text-base font-semibold cursor-pointer">Publish Status</FormLabel>
												<p className="text-xs text-muted-foreground mt-1">
													{field.value ? "This post is live and visible to readers" : "This post is in draft mode"}
												</p>
											</div>
											<FormControl>
												<button
													type="button"
													onClick={() => field.onChange(!field.value)}
													className="focus:outline-none"
												>
													<Badge
														variant={field.value ? "default" : "secondary"}
														className="ml-4 cursor-pointer"
													>
														{field.value ? "Published" : "Draft"}
													</Badge>
												</button>
											</FormControl>
										</FormItem>
									)}
								/>
							</div>
						</form>
					</Form>
				</div>

				{/* Footer Actions */}
				<div className="p-6 border-t border-border/50 bg-muted/30 flex justify-end gap-3">
					<Button
						type="button"
						variant="outline"
						onClick={() => onClose()}
						disabled={loading || uploading}
						className="rounded-lg h-11 px-6"
					>
						Cancel
					</Button>
					<Button
						type="submit"
						form="blog-form"
						disabled={loading || uploading}
						className="rounded-lg h-11 px-6 gap-2"
					>
						{loading && (
							<Loader2 className="h-4 w-4 animate-spin" />
						)}
						{blog ? "Update" : "Create"} Blog Post
					</Button>
				</div>
			</DialogContent>
		</Dialog>
	);
}
