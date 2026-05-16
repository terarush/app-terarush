import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { getBlogs, getImageUrl, type Blog } from "@/lib/api/blogs";
import { Button } from "@/components/ui/button";
import { Loader2, Calendar, User, Eye, ArrowRight, Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import Navbar from "@/components/navbar";
import Footer from "@/components/footer";

export function BlogList() {
	const [blogs, setBlog] = useState<Blog[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [page, setPage] = useState(1);
	const [total, setTotal] = useState(0);
	const [searchQuery, setSearchQuery] = useState("");
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

	const filteredBlogs = blogs.filter((blog) =>
		blog.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
		blog.excerpt.toLowerCase().includes(searchQuery.toLowerCase())
	);

	const totalPages = Math.ceil(total / pageSize);

	if (loading && blogs.length === 0) {
		return (
			<div className="min-h-screen bg-background flex flex-col">
				<Navbar />
				<div className="flex-1 flex items-center justify-center pt-24">
					<Loader2 className="h-8 w-8 animate-spin text-primary" />
				</div>
				<Footer />
			</div>
		);
	}

	return (
		<div className="min-h-screen bg-background flex flex-col">
			<Navbar />
			
			<main className="flex-1 pt-24">
				{/* Hero Section */}
				<div className="relative overflow-hidden border-b border-border">
					<div className="container mx-auto px-4 py-16 md:py-24">
						<div className="max-w-3xl mx-auto text-center">
							<h1 className="text-5xl md:text-6xl font-bold mb-6">
								Our Blog
							</h1>
							<p className="text-xl text-muted-foreground mb-8 leading-relaxed">
								Explore our latest insights on web development, best practices, and technology trends.
								Stay updated with in-depth articles from our team.
							</p>

							{/* Search Bar */}
							<div className="relative max-w-md mx-auto">
								<Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
								<Input
									placeholder="Search articles..."
									value={searchQuery}
									onChange={(e) => setSearchQuery(e.target.value)}
									className="pl-12 py-3 rounded-xl bg-background/50 backdrop-blur-sm border border-border/50 hover:border-primary/30 transition-colors"
								/>
							</div>
						</div>
					</div>
				</div>

				{/* Blog Grid */}
				<div className="container mx-auto px-4 py-16">
					{error && (
						<div className="mb-8 p-6 rounded-xl bg-destructive/10 border border-destructive/20 text-destructive">
							<p className="font-medium">{error}</p>
							<Button 
								onClick={loadBlogs} 
								variant="outline" 
								className="mt-4"
							>
								Try Again
							</Button>
						</div>
					)}

					{filteredBlogs.length === 0 ? (
						<div className="text-center py-16">
							<p className="text-muted-foreground text-lg mb-4">
								{searchQuery ? "No articles found matching your search." : "No blog posts available at the moment."}
							</p>
						</div>
					) : (
						<>
							<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
								{filteredBlogs.map((blog) => (
									<article
										key={blog.id}
										className="group relative rounded-2xl overflow-hidden bg-card border border-border/50 hover:border-primary/30 shadow-sm hover:shadow-lg transition-all duration-300 cursor-pointer hover:scale-105 transform"
										onClick={() => navigate(`/blog/${blog.slug}`)}
									>
										{/* Image Container */}
										<div className="relative h-56 overflow-hidden bg-muted">
											{blog.image && (
												<img
													src={getImageUrl(blog.image)}
													alt={blog.title}
													className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
												/>
											)}
											<div className="absolute inset-0 bg-gradient-to-t from-background via-transparent to-transparent opacity-60" />

											{/* Tags */}
											<div className="absolute top-4 left-4 flex gap-2 z-10">
												{blog.category && (
													<span className="px-3 py-1 text-xs font-semibold rounded-full bg-primary/90 text-primary-foreground backdrop-blur-sm">
														{blog.category}
													</span>
												)}
											</div>
										</div>

										{/* Content */}
										<div className="p-6 relative z-10">
											<h2 className="text-xl font-bold mb-3 line-clamp-2 group-hover:text-primary transition-colors">
												{blog.title}
											</h2>

											<p className="text-muted-foreground text-sm mb-4 line-clamp-2">
												{blog.excerpt}
											</p>

											{/* Meta Info */}
											<div className="flex items-center gap-4 text-xs text-muted-foreground mb-4 pb-4 border-t border-border/30">
												<div className="flex items-center gap-1 pt-4">
													<User className="h-3.5 w-3.5" />
													<span className="truncate">{blog.author}</span>
												</div>
												<div className="flex items-center gap-1 pt-4">
													<Calendar className="h-3.5 w-3.5" />
													<span>
														{new Date(blog.created_at).toLocaleDateString('en-US', {
															month: 'short',
															day: 'numeric'
														})}
													</span>
												</div>
												<div className="flex items-center gap-1 pt-4 ml-auto">
													<Eye className="h-3.5 w-3.5" />
													<span>{blog.view_count}</span>
												</div>
											</div>

											{/* Read More Link */}
											<div className="flex items-center gap-2 text-sm font-medium text-primary group-hover:gap-3 transition-all">
												Read More
												<ArrowRight className="h-4 w-4" />
											</div>
										</div>
									</article>
								))}
							</div>

							{/* Pagination */}
							{totalPages > 1 && (
								<div className="flex items-center justify-center gap-4 py-8">
									<Button
										variant="outline"
										onClick={() => setPage((p) => Math.max(1, p - 1))}
										disabled={page === 1 || loading}
										className="rounded-xl"
									>
										Previous
									</Button>
									<div className="flex items-center gap-2">
										{Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
											const pageNum = i + 1;
											return (
												<Button
													key={pageNum}
													variant={page === pageNum ? "default" : "outline"}
													onClick={() => setPage(pageNum)}
													className="rounded-lg h-10 w-10"
													disabled={loading}
												>
													{pageNum}
												</Button>
											);
										})}
									</div>
									<Button
										variant="outline"
										onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
										disabled={page === totalPages || loading}
										className="rounded-xl"
									>
										Next
									</Button>
								</div>
							)}
						</>
					)}
				</div>
			</main>

			<Footer />
		</div>
	);
}
