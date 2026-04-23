import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { getBlogs, type Blog } from "@/lib/api/blogs";
import { Button } from "@/components/ui/button";
import { Loader2, Calendar, User, Eye } from "lucide-react";

export function BlogList() {
	const [blogs, setBlog] = useState<Blog[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [page, setPage] = useState(1);
	const [total, setTotal] = useState(0);
	const navigate = useNavigate();
	const pageSize = 6;

	useEffect(() => {
		loadBlogs();
	}, [page]);

	const loadBlogs = async () => {
		try {
			setLoading(true);
			setError(null);
			const response = await getBlogs({
				page,
				page_size: pageSize,
			});
			setBlog(response.blogs || []);
			setTotal(response.total || 0);
		} catch (err) {
			console.error("Error loading blogs:", err);
			setError("Failed to load blogs. Please try again later.");
		} finally {
			setLoading(false);
		}
	};

	const totalPages = Math.ceil(total / pageSize);

	if (loading && blogs.length === 0) {
		return (
			<div className="container mx-auto px-4 py-8">
				<div className="flex items-center justify-center min-h-[400px]">
					<Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div className="container mx-auto px-4 py-8">
				<div className="text-center py-12">
					<p className="text-destructive mb-4">{error}</p>
					<Button onClick={loadBlogs}>Try Again</Button>
				</div>
			</div>
		);
	}

	return (
		<div className="container mx-auto px-4 py-12">
			<div className="mb-12">
				<h1 className="text-4xl font-bold mb-4">Blog</h1>
				<p className="text-muted-foreground text-lg max-w-2xl">
					Explore our latest articles on web development, app development, and
					best practices for building scalable applications.
				</p>
			</div>

			{blogs.length === 0 ? (
				<div className="text-center py-12">
					<p className="text-muted-foreground text-lg">
						No blog posts available at the moment.
					</p>
				</div>
			) : (
				<>
					<div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-12">
						{blogs.map((blog) => (
							<article
								key={blog.id}
								className="border rounded-lg overflow-hidden hover:shadow-lg transition-shadow cursor-pointer bg-card"
								onClick={() => navigate(`/blog/${blog.slug}`)}
							>
								{blog.image && (
									<img
										src={blog.image}
										alt={blog.title}
										className="w-full h-48 object-cover"
									/>
								)}
								<div className="p-6">
									<div className="flex gap-2 mb-3">
										{blog.category && (
											<span className="inline-block px-3 py-1 text-xs font-semibold rounded-full bg-primary/10 text-primary">
												{blog.category}
											</span>
										)}
										{blog.tags && (
											<span className="inline-block px-3 py-1 text-xs font-semibold rounded-full bg-secondary/10 text-secondary-foreground">
												{blog.tags.split(",")[0]}
											</span>
										)}
									</div>

									<h2 className="text-2xl font-bold mb-3 hover:text-primary transition-colors">
										{blog.title}
									</h2>

									<p className="text-muted-foreground mb-4 line-clamp-2">
										{blog.excerpt}
									</p>

									<div className="flex items-center gap-4 text-sm text-muted-foreground border-t pt-4">
										<div className="flex items-center gap-1">
											<User className="h-4 w-4" />
											<span>{blog.author}</span>
										</div>
										<div className="flex items-center gap-1">
											<Calendar className="h-4 w-4" />
											<span>
												{new Date(blog.created_at).toLocaleDateString()}
											</span>
										</div>
										<div className="flex items-center gap-1 ml-auto">
											<Eye className="h-4 w-4" />
											<span>{blog.view_count}</span>
										</div>
									</div>
								</div>
							</article>
						))}
					</div>

					{totalPages > 1 && (
						<div className="flex items-center justify-center gap-2">
							<Button
								variant="outline"
								onClick={() => setPage((p) => Math.max(1, p - 1))}
								disabled={page === 1 || loading}
							>
								Previous
							</Button>
							<span className="text-sm text-muted-foreground px-4">
								Page {page} of {totalPages}
							</span>
							<Button
								variant="outline"
								onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
								disabled={page === totalPages || loading}
							>
								Next
							</Button>
						</div>
					)}
				</>
			)}
		</div>
	);
}
