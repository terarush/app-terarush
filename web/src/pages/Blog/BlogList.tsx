import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { getBlogs, type Blog } from "@/lib/api/blogs";
import Navbar from "@/components/navbar";
import Footer from "@/components/footer";
import { LoadingSpinner } from "@/components/elements/loading-spinner";
import { ErrorMessage } from "@/components/elements/error-message";
import { SearchInput } from "@/components/elements/search-input";
import { LoadMoreButton } from "@/components/elements/load-more-button";
import { BlogCard } from "@/components/fragments/blog-card";
import {
	PageLayout,
	PageContent,
	HeroSection,
	ContentSection,
} from "@/components/layouts/page-layout";

export function BlogList() {
	const [blogs, setBlogs] = useState<Blog[]>([]);
	const [loading, setLoading] = useState(true);
	const [loadingMore, setLoadingMore] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [page, setPage] = useState(1);
	const [total, setTotal] = useState(0);
	const [searchQuery, setSearchQuery] = useState("");
	const navigate = useNavigate();
	const pageSize = 10;

	useEffect(() => {
		if (page === 1) {
			loadBlogs();
		} else {
			loadMoreBlogs();
		}
	}, [page]);

	const loadBlogs = async () => {
		try {
			setLoading(true);
			setError(null);
			const response = await getBlogs({
				page: 1,
				page_size: pageSize,
			});
			setBlogs(response.blogs || []);
			setTotal(response.total || 0);
			setPage(1);
		} catch (err) {
			console.error("Error loading blogs:", err);
			setError("Failed to load blogs. Please try again later.");
		} finally {
			setLoading(false);
		}
	};

	const loadMoreBlogs = async () => {
		try {
			setLoadingMore(true);
			setError(null);
			const response = await getBlogs({
				page,
				page_size: pageSize,
			});
			setBlogs((prevBlogs) => [...prevBlogs, ...(response.blogs || [])]);
			setTotal(response.total || 0);
		} catch (err) {
			console.error("Error loading more blogs:", err);
			setError("Failed to load more blogs. Please try again later.");
		} finally {
			setLoadingMore(false);
		}
	};

	const filteredBlogs = blogs.filter((blog) =>
		blog.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
		blog.excerpt.toLowerCase().includes(searchQuery.toLowerCase())
	);

	// Only show load more if not searching, and there are more blogs to load
	const hasMore = !searchQuery && blogs.length < total;

	if (loading && blogs.length === 0) {
		return (
			<PageLayout>
				<Navbar />
				<PageContent className="flex items-center justify-center">
					<LoadingSpinner size="lg" text="Loading blogs..." />
				</PageContent>
				<Footer />
			</PageLayout>
		);
	}

	return (
		<PageLayout>
			<Navbar />

			<PageContent>
				{/* Hero Section */}
				<HeroSection>
					<div className="max-w-3xl mx-auto text-center">
						<h1 className="text-4xl md:text-5xl font-bold mb-4">
							Blog & Insights
						</h1>
						<p className="text-lg text-muted-foreground mb-8 leading-relaxed">
							Explore our latest insights on web development, best practices, and technology trends.
						</p>

						{/* Search Bar */}
						<div className="max-w-md mx-auto">
							<SearchInput
								placeholder="Search articles..."
								value={searchQuery}
								onChange={setSearchQuery}
							/>
						</div>
					</div>
				</HeroSection>

				{/* Blog Grid */}
				<ContentSection>
					{error && (
						<div className="mb-8">
							<ErrorMessage message={error} onRetry={loadBlogs} />
						</div>
					)}

					{filteredBlogs.length === 0 ? (
						<div className="text-center py-16">
							<p className="text-muted-foreground text-lg">
								{searchQuery
									? "No articles found matching your search."
									: "No blog posts available at the moment."}
							</p>
						</div>
					) : (
						<>
							<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
								{filteredBlogs.map((blog) => (
									<BlogCard
										key={blog.id}
										blog={blog}
										onClick={() => navigate(`/blog/${blog.slug}`)}
									/>
								))}
							</div>

							{/* Load More Button */}
							<LoadMoreButton
								onClick={() => setPage((prev) => prev + 1)}
								isLoading={loadingMore}
								hasMore={hasMore}
								itemCount={blogs.length}
								pageSize={pageSize}
							/>
						</>
					)}
				</ContentSection>
			</PageContent>

			<Footer />
		</PageLayout>
	);
}
