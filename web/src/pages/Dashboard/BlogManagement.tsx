import { useState, useEffect } from "react";
import { Plus, Edit2, Trash2, Loader2, Eye, Calendar, Tag, Search, ChevronLeft, ChevronRight } from "lucide-react";
import {
	getAllBlogs,
	deleteBlog,
	type Blog,
} from "@/lib/api/blogs";
import { Button } from "@/components/ui/button";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { BlogForm } from "@/components/fragments/BlogForm";

export function BlogManagement() {
	const [blogs, setBlogs] = useState<Blog[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [selectedBlog, setSelectedBlog] = useState<Blog | null>(null);
	const [isFormOpen, setIsFormOpen] = useState(false);
	const [deleteId, setDeleteId] = useState<number | null>(null);
	const [isDeleting, setIsDeleting] = useState(false);
	const [searchQuery, setSearchQuery] = useState("");
	const [currentPage, setCurrentPage] = useState(1);
	const itemsPerPage = 10;

	useEffect(() => {
		loadBlogs();
	}, []);

	const loadBlogs = async () => {
		try {
			setLoading(true);
			setError(null);
			const response = await getAllBlogs();
			setBlogs(response.blogs || []);
		} catch (err) {
			console.error("Error loading blogs:", err);
			setError("Failed to load blogs");
			toast.error("Failed to load blogs");
		} finally {
			setLoading(false);
		}
	};

	// Filter blogs based on search query
	const filteredBlogs = blogs.filter((blog) => {
		const searchLower = searchQuery.toLowerCase();
		return (
			blog.title.toLowerCase().includes(searchLower) ||
			blog.excerpt?.toLowerCase().includes(searchLower) ||
			blog.category?.toLowerCase().includes(searchLower) ||
			blog.author.toLowerCase().includes(searchLower)
		);
	});

	// Pagination logic
	const totalPages = Math.ceil(filteredBlogs.length / itemsPerPage);
	const startIndex = (currentPage - 1) * itemsPerPage;
	const endIndex = startIndex + itemsPerPage;
	const paginatedBlogs = filteredBlogs.slice(startIndex, endIndex);

	// Reset to first page when search query changes
	useEffect(() => {
		setCurrentPage(1);
	}, [searchQuery]);

	const handleCreate = () => {
		setSelectedBlog(null);
		setIsFormOpen(true);
	};

	const handleEdit = (blog: Blog) => {
		setSelectedBlog(blog);
		setIsFormOpen(true);
	};

	const handleFormClose = (saved?: boolean) => {
		setIsFormOpen(false);
		setSelectedBlog(null);
		if (saved) {
			loadBlogs();
		}
	};

	const handleDelete = async (id: number) => {
		try {
			setIsDeleting(true);
			await deleteBlog(id);
			toast.success("Blog deleted successfully");
			setDeleteId(null);
			loadBlogs();
		} catch (err) {
			console.error("Error deleting blog:", err);
			toast.error("Failed to delete blog");
		} finally {
			setIsDeleting(false);
		}
	};

	if (loading) {
		return (
			<div className="flex items-center justify-center min-h-96">
				<Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
			</div>
		);
	}

	if (error) {
		return (
			<div className="text-center py-12">
				<p className="text-destructive mb-4">{error}</p>
				<Button onClick={loadBlogs}>Try Again</Button>
			</div>
		);
	}

	return (
		<>
			<div className="mb-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
				<div>
					<h1 className="text-3xl font-bold">Blog Management</h1>
					<p className="text-muted-foreground mt-1">
						Create, edit, and manage your blog posts
					</p>
				</div>
				<Button onClick={handleCreate} className="gap-2 w-full md:w-auto">
					<Plus className="h-4 w-4" />
					New Blog Post
				</Button>
			</div>

			{isFormOpen && (
				<BlogForm
					blog={selectedBlog || undefined}
					onClose={handleFormClose}
				/>
			)}

			{blogs.length === 0 ? (
				<div className="text-center py-12 border rounded-lg">
					<p className="text-muted-foreground mb-4">
						No blog posts yet. Create your first one!
					</p>
					<Button onClick={handleCreate} variant="outline">
						Create Blog Post
					</Button>
				</div>
			) : (
				<>
					{/* Search Bar */}
					<div className="mb-6 flex items-center gap-2">
						<div className="relative flex-1 max-w-md">
							<Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
							<input
								type="text"
								placeholder="Search by title, excerpt, category, or author..."
								value={searchQuery}
								onChange={(e) => setSearchQuery(e.target.value)}
								className="w-full pl-10 pr-4 py-2 border border-border/50 rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary"
							/>
						</div>
						{searchQuery && (
							<Button
								variant="ghost"
								size="sm"
								onClick={() => setSearchQuery("")}
								className="text-muted-foreground hover:text-foreground"
							>
								Clear
							</Button>
						)}
					</div>

					{/* Results info */}
					<div className="mb-4 text-sm text-muted-foreground">
						Showing {paginatedBlogs.length === 0 ? 0 : startIndex + 1} to {Math.min(endIndex, filteredBlogs.length)} of {filteredBlogs.length} {filteredBlogs.length === 1 ? 'post' : 'posts'}
						{searchQuery && <span> (filtered from {blogs.length} total)</span>}
					</div>

					{paginatedBlogs.length === 0 ? (
						<div className="text-center py-12 border rounded-lg">
							<p className="text-muted-foreground">
								No blog posts found matching your search criteria.
							</p>
						</div>
					) : (
						<>
							<div className="border border-border/50 rounded-xl overflow-hidden bg-card shadow-sm">
								<Table>
									<TableHeader className="bg-muted/50 border-b border-border/50">
										<TableRow className="hover:bg-transparent">
											<TableHead className="font-semibold text-foreground w-1/4">Article</TableHead>
											<TableHead className="font-semibold text-foreground">Category</TableHead>
											<TableHead className="font-semibold text-foreground">Author</TableHead>
											<TableHead className="font-semibold text-foreground">Stats</TableHead>
											<TableHead className="font-semibold text-foreground">Status</TableHead>
											<TableHead className="font-semibold text-foreground">Date</TableHead>
											<TableHead className="text-right font-semibold text-foreground">Actions</TableHead>
										</TableRow>
									</TableHeader>
									<TableBody>
										{paginatedBlogs.map((blog) => (
											<TableRow key={blog.id} className="border-border/50 hover:bg-muted/40">
												<TableCell className="py-4">
													<div className="space-y-2">
														<p className="font-semibold text-sm leading-snug line-clamp-2 hover:text-primary transition-colors cursor-pointer">
															{blog.title}
														</p>
														<p className="text-xs text-muted-foreground line-clamp-2">
															{blog.excerpt || "No excerpt provided"}
														</p>
													</div>
												</TableCell>
												<TableCell className="py-4">
													<div className="flex items-center gap-1">
														<Tag className="h-3.5 w-3.5 text-muted-foreground" />
														<span className="text-sm font-medium">
															{blog.category || <span className="text-muted-foreground">-</span>}
														</span>
													</div>
												</TableCell>
												<TableCell className="py-4">
													<div className="flex items-center gap-2">
														<div className="w-6 h-6 rounded-full bg-primary/10 flex items-center justify-center">
															<span className="text-xs font-bold text-primary">
																{blog.author.charAt(0).toUpperCase()}
															</span>
														</div>
														<span className="text-sm">{blog.author}</span>
													</div>
												</TableCell>
												<TableCell className="py-4">
													<div className="space-y-1.5">
														<div className="flex items-center gap-2 text-sm">
															<Eye className="h-3.5 w-3.5 text-amber-500" />
															<span className="font-medium">{blog.view_count}</span>
															<span className="text-xs text-muted-foreground">views</span>
														</div>
													</div>
												</TableCell>
												<TableCell className="py-4">
													<Badge
														variant={
															blog.is_published ? "default" : "secondary"
														}
														className="font-medium"
													>
														{blog.is_published ? "Published" : "Draft"}
													</Badge>
												</TableCell>
												<TableCell className="py-4">
													<div className="flex items-center gap-2 text-sm">
														<Calendar className="h-3.5 w-3.5 text-muted-foreground" />
														<span>{new Date(blog.created_at).toLocaleDateString('en-US', {
															month: 'short',
															day: 'numeric',
															year: 'numeric'
														})}</span>
													</div>
												</TableCell>
												<TableCell className="text-right py-4">
													<div className="flex items-center justify-end gap-1">
														<Button
															size="sm"
															variant="ghost"
															className="h-8 w-8 p-0 hover:bg-primary/10 hover:text-primary"
															onClick={() => handleEdit(blog)}
															title="Edit blog"
														>
															<Edit2 className="h-4 w-4" />
														</Button>
														<Button
															size="sm"
															variant="ghost"
															className="h-8 w-8 p-0 hover:bg-destructive/10 hover:text-destructive"
															onClick={() => setDeleteId(blog.id)}
															title="Delete blog"
														>
															<Trash2 className="h-4 w-4" />
														</Button>
													</div>
												</TableCell>
											</TableRow>
										))}
									</TableBody>
								</Table>
							</div>

							{/* Pagination */}
							{totalPages > 1 && (
								<div className="mt-6 flex items-center justify-between">
									<div className="text-sm text-muted-foreground">
										Page {currentPage} of {totalPages}
									</div>
									<div className="flex items-center gap-2">
										<Button
											variant="outline"
											size="sm"
											onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
											disabled={currentPage === 1}
											className="gap-1"
										>
											<ChevronLeft className="h-4 w-4" />
											Previous
										</Button>
										<div className="flex gap-1">
											{Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
												<Button
													key={page}
													variant={currentPage === page ? "default" : "outline"}
													size="sm"
													onClick={() => setCurrentPage(page)}
													className="h-8 w-8 p-0"
												>
													{page}
												</Button>
											))}
										</div>
										<Button
											variant="outline"
											size="sm"
											onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
											disabled={currentPage === totalPages}
											className="gap-1"
										>
											Next
											<ChevronRight className="h-4 w-4" />
										</Button>
									</div>
								</div>
							)}
						</>
					)}
				</>
			)}

			<AlertDialog open={deleteId !== null} onOpenChange={(open) => !open && setDeleteId(null)}>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle>Delete Blog Post</AlertDialogTitle>
						<AlertDialogDescription>
							Are you sure you want to delete this blog post? This action cannot be undone.
						</AlertDialogDescription>
					</AlertDialogHeader>
					<div className="flex justify-end gap-2">
						<AlertDialogCancel>Cancel</AlertDialogCancel>
						<AlertDialogAction
							onClick={() => deleteId && handleDelete(deleteId)}
							disabled={isDeleting}
							className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
						>
							{isDeleting ? (
								<Loader2 className="h-4 w-4 animate-spin mr-2" />
							) : null}
							Delete
						</AlertDialogAction>
					</div>
				</AlertDialogContent>
			</AlertDialog>
		</>
	);
}
