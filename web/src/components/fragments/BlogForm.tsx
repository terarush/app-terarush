import { useState } from "react";
import { useForm } from "react-hook-form";
import {
	createBlog,
	updateBlog,
	uploadBlogImage,
	type Blog,
	type CreateBlogRequest,
} from "@/lib/api/blogs";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/components/ui/form";
import { Badge } from "@/components/ui/badge";
import { Loader2, Upload, X } from "lucide-react";
import { toast } from "sonner";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

interface BlogFormProps {
	blog?: Blog;
	onClose: (saved?: boolean) => void;
}

export function BlogForm({ blog, onClose }: BlogFormProps) {
	const [loading, setLoading] = useState(false);
	const [uploading, setUploading] = useState(false);
	const [previewUrl, setPreviewUrl] = useState<string>(blog?.image || "");
	const form = useForm<CreateBlogRequest>({
		defaultValues: blog
			? {
					title: blog.title,
					slug: blog.slug,
					content: blog.content,
					excerpt: blog.excerpt,
					author: blog.author,
					category: blog.category,
					tags: blog.tags,
					image: blog.image,
					is_published: blog.is_published,
				}
			: {
					title: "",
					slug: "",
					content: "",
					excerpt: "",
					author: "",
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
			setPreviewUrl(imageUrl);
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
		setPreviewUrl(url);
	};

	const clearImage = () => {
		form.setValue("image", "");
		setPreviewUrl("");
	};

	return (
		<Dialog open={true} onOpenChange={(open) => !open && onClose()}>
			<DialogContent className="max-w-2xl scrollbar-hide max-h-[90vh] overflow-y-auto">
				<DialogHeader>
					<DialogTitle>
						{blog ? "Edit Blog Post" : "Create New Blog Post"}
					</DialogTitle>
					<DialogDescription>
						{blog
							? "Update your blog post details"
							: "Fill in the form to create a new blog post"}
					</DialogDescription>
				</DialogHeader>

				<Form {...form}>
					<form
						onSubmit={form.handleSubmit(onSubmit)}
						className="space-y-4"
					>
						<FormField
							control={form.control}
							name="title"
							rules={{
								required: "Title is required",
								maxLength: {
									value: 255,
									message:
										"Title must be less than 255 characters",
								},
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>Title</FormLabel>
									<FormControl>
										<Input
											placeholder="Enter blog post title"
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
							rules={{
								required: "Slug is required",
								maxLength: {
									value: 255,
									message:
										"Slug must be less than 255 characters",
								},
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>Slug</FormLabel>
									<FormControl>
										<Input
											placeholder="blog-post-slug"
											{...field}
										/>
									</FormControl>
									<FormDescription>
										URL-friendly version of the title
									</FormDescription>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="excerpt"
							rules={{
								maxLength: {
									value: 500,
									message:
										"Excerpt must be less than 500 characters",
								},
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>Excerpt</FormLabel>
									<FormControl>
										<Textarea
											placeholder="Brief summary of the post"
											{...field}
											rows={2}
										/>
									</FormControl>
									<FormDescription>
										{field.value?.length || 0}/500 characters
									</FormDescription>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="content"
							rules={{
								required: "Content is required",
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>Content</FormLabel>
									<FormControl>
										<Textarea
											placeholder="Write your blog post content here..."
											{...field}
											rows={6}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<div className="grid grid-cols-2 gap-4">
							<FormField
								control={form.control}
								name="author"
								rules={{
									required: "Author is required",
									maxLength: {
										value: 100,
										message:
											"Author must be less than 100 characters",
									},
								}}
								render={({ field }) => (
									<FormItem>
										<FormLabel>Author</FormLabel>
										<FormControl>
											<Input
												placeholder="Author name"
												{...field}
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								)}
							/>

							<FormField
								control={form.control}
								name="category"
								rules={{
									maxLength: {
										value: 100,
										message:
											"Category must be less than 100 characters",
									},
								}}
								render={({ field }) => (
									<FormItem>
										<FormLabel>Category</FormLabel>
										<FormControl>
											<Input
												placeholder="e.g. Web Development"
												{...field}
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								)}
							/>
						</div>

						<FormField
							control={form.control}
							name="tags"
							rules={{
								maxLength: {
									value: 255,
									message:
										"Tags must be less than 255 characters",
								},
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>Tags</FormLabel>
									<FormControl>
										<Input
											placeholder="Comma-separated tags (e.g. react, javascript, web)"
											{...field}
										/>
									</FormControl>
									<FormDescription>
										Separate multiple tags with commas
									</FormDescription>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="image"
							rules={{
								maxLength: {
									value: 255,
									message:
										"Image URL must be less than 255 characters",
								},
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>Featured Image</FormLabel>
									<Tabs defaultValue="url" className="w-full">
										<TabsList className="grid w-full grid-cols-2">
											<TabsTrigger value="url">From URL</TabsTrigger>
											<TabsTrigger value="upload">Upload</TabsTrigger>
										</TabsList>
										<TabsContent value="url" className="space-y-3">
											<FormControl>
												<Input
													placeholder="https://example.com/image.jpg"
													value={field.value}
													onChange={(e) => {
														handleUrlInput(e.target.value);
													}}
												/>
											</FormControl>
											<FormDescription>
												Enter the URL of your featured image
											</FormDescription>
										</TabsContent>
										<TabsContent value="upload" className="space-y-3">
											<div className="relative">
												<label className="flex flex-col items-center justify-center w-full border-2 border-dashed border-gray-300 rounded-lg cursor-pointer hover:bg-gray-50 p-6">
													<div className="flex flex-col items-center justify-center">
														<Upload className="h-8 w-8 text-gray-400 mb-2" />
														<p className="text-sm text-gray-500">
															{uploading
																? "Uploading..."
																: "Click to upload or drag and drop"}
														</p>
														<p className="text-xs text-gray-400">
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
									{previewUrl && (
										<div className="relative mt-4">
											<img
												src={previewUrl}
												alt="Preview"
												className="w-full h-48 object-cover rounded-lg"
											/>
											<button
												type="button"
												onClick={clearImage}
												className="absolute top-2 right-2 bg-red-500 text-white p-1 rounded-full hover:bg-red-600"
											>
												<X className="h-4 w-4" />
											</button>
										</div>
									)}
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="is_published"
							render={({ field }) => (
								<FormItem className="flex items-center justify-between p-3 border rounded-lg">
									<div>
										<FormLabel>Publish</FormLabel>
										<FormDescription>
											Make this post visible to the public
										</FormDescription>
									</div>
									<FormControl>
										<div
											className="cursor-pointer"
											onClick={() =>
												field.onChange(!field.value)
											}
										>
											<Badge
												variant={
													field.value ? "default" : "secondary"
												}
											>
												{field.value
													? "Published"
													: "Draft"}
											</Badge>
										</div>
									</FormControl>
								</FormItem>
							)}
						/>

						<div className="flex justify-end gap-2 pt-4">
							<Button
								type="button"
								variant="outline"
								onClick={() => onClose()}
								disabled={loading}
							>
								Cancel
							</Button>
							<Button type="submit" disabled={loading}>
								{loading && (
									<Loader2 className="h-4 w-4 mr-2 animate-spin" />
								)}
								{blog ? "Update" : "Create"} Blog Post
							</Button>
						</div>
					</form>
				</Form>
			</DialogContent>
		</Dialog>
	);
}
