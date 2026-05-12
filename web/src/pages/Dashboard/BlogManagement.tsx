import { useState, useEffect } from "react";
import { Plus, Edit2, Trash2, Loader2 } from "lucide-react";
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
			<div className="mb-6 flex items-center justify-between">
				<div>
					<h1 className="text-3xl font-bold">Blog Management</h1>
					<p className="text-muted-foreground mt-1">
						Create, edit, and manage your blog posts
					</p>
				</div>
				<Button onClick={handleCreate} className="gap-2">
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
				<div className="border rounded-lg overflow-x-auto">
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead>Title</TableHead>
								<TableHead>Author</TableHead>
								<TableHead>Category</TableHead>
								<TableHead>Status</TableHead>
								<TableHead>Views</TableHead>
								<TableHead>Date</TableHead>
								<TableHead className="text-right">Actions</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							{blogs.map((blog) => (
								<TableRow key={blog.id}>
									<TableCell className="font-medium max-w-xs truncate">
										{blog.title}
									</TableCell>
									<TableCell>{blog.author}</TableCell>
									<TableCell>{blog.category || "-"}</TableCell>
									<TableCell>
										<Badge
											variant={
												blog.is_published ? "default" : "secondary"
											}
										>
											{blog.is_published ? "Published" : "Draft"}
										</Badge>
									</TableCell>
									<TableCell>{blog.view_count}</TableCell>
									<TableCell>
										{new Date(blog.created_at).toLocaleDateString()}
									</TableCell>
									<TableCell className="text-right space-x-2">
										<Button
											size="sm"
											variant="ghost"
											onClick={() => handleEdit(blog)}
										>
											<Edit2 className="h-4 w-4" />
										</Button>
										<Button
											size="sm"
											variant="ghost"
											className="text-destructive hover:text-destructive"
											onClick={() => setDeleteId(blog.id)}
										>
											<Trash2 className="h-4 w-4" />
										</Button>
									</TableCell>
								</TableRow>
							))}
						</TableBody>
					</Table>
				</div>
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
