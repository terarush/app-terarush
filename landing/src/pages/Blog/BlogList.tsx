import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { getBlogs } from "@/lib/api/blogs";
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
import { useInfiniteQuery } from "@tanstack/react-query";

export function BlogList() {
	const [searchQuery, setSearchQuery] = useState("");
	const navigate = useNavigate();
	const pageSize = 10;

	const {
		data,
		error,
		fetchNextPage,
		hasNextPage,
		isFetchingNextPage,
		status,
		refetch,
	} = useInfiniteQuery({
		queryKey: ["blogs", searchQuery],
		queryFn: ({ pageParam = 1 }) =>
			getBlogs({
				page: pageParam,
				page_size: pageSize,
				search: searchQuery || undefined,
			}),
		initialPageParam: 1,
		getNextPageParam: (lastPage, allPages) => {
			const loadedCount = allPages.reduce(
				(acc, curr) => acc + (curr.blogs?.length || 0),
				0,
			);
			return loadedCount < (lastPage.total || 0) ? allPages.length + 1 : undefined;
		},
	});

	const blogs = data?.pages.flatMap((page) => page.blogs || []) || [];
	const loading = status === "pending";
	const loadingMore = isFetchingNextPage;

	// Server-side search is used. Keep client-side list as-is for rendering.
	const filteredBlogs = blogs;

	// Only show load more if there are more blogs to load
	const hasMore = hasNextPage;

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
					{status === "error" && (
						<div className="mb-8">
							<ErrorMessage
								message={error instanceof Error ? error.message : "Failed to load blogs. Please try again later."}
								onRetry={() => refetch()}
							/>
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
								onClick={() => fetchNextPage()}
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
