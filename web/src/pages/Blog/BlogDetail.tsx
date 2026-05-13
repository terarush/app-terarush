import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getBlogBySlug, type Blog } from "@/lib/api/blogs";
import { Button } from "@/components/ui/button";
import { Loader2, ArrowLeft, Calendar, User, Eye } from "lucide-react";
import { MarkdownRenderer } from "@/components/fragments/MarkdownRenderer";

export function BlogDetail() {
	const { slug } = useParams<{ slug: string }>();
	const navigate = useNavigate();
	const [blog, setBlog] = useState<Blog | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		if (slug) {
			loadBlog();
		}
	}, [slug]);

	const loadBlog = async () => {
		try {
			setLoading(true);
			setError(null);
			const data = await getBlogBySlug(slug!);
			setBlog(data);
		} catch (err) {
			console.error("Error loading blog:", err);
			setError("Failed to load blog post. Please try again later.");
		} finally {
			setLoading(false);
		}
	};

	if (loading) {
		return (
			<div className="container mx-auto px-4 py-8">
				<div className="flex items-center justify-center min-h-[400px]">
					<Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
				</div>
			</div>
		);
	}

	if (error || !blog) {
		return (
			<div className="container mx-auto px-4 py-8">
				<Button
					variant="ghost"
					onClick={() => navigate("/blog")}
					className="mb-4"
				>
					<ArrowLeft className="h-4 w-4 mr-2" />
					Back to Blog
				</Button>
				<div className="text-center py-12">
					<p className="text-destructive mb-4">{error || "Blog not found"}</p>
					<Button onClick={() => navigate("/blog")}>Return to Blog</Button>
				</div>
			</div>
		);
	}

	const formattedDate = new Date(blog.created_at).toLocaleDateString("en-US", {
		year: "numeric",
		month: "long",
		day: "numeric",
	});

	return (
		<div className="min-h-screen bg-background">
			<div className="container mx-auto px-4 py-8">
				<Button
					variant="ghost"
					onClick={() => navigate("/blog")}
					className="mb-8"
				>
					<ArrowLeft className="h-4 w-4 mr-2" />
					Back to Blog
				</Button>

				<article className="max-w-3xl mx-auto">
					{/* Header */}
					<header className="mb-8">
						<div className="flex gap-2 mb-4">
							{blog.category && (
								<span className="inline-block px-3 py-1 text-xs font-semibold rounded-full bg-primary/10 text-primary">
									{blog.category}
								</span>
							)}
							{blog.tags && (
								blog.tags.split(",").map((tag, index) => (
									<span
										key={index}
										className="inline-block px-3 py-1 text-xs font-semibold rounded-full bg-secondary/10 text-secondary-foreground"
									>
										{tag.trim()}
									</span>
								))
							)}
						</div>

						<h1 className="text-4xl md:text-5xl font-bold mb-4">
							{blog.title}
						</h1>

						<div className="flex flex-wrap items-center gap-6 text-sm text-muted-foreground border-b pb-4">
							<div className="flex items-center gap-2">
								<User className="h-4 w-4" />
								<span className="font-medium">{blog.author}</span>
							</div>
							<div className="flex items-center gap-2">
								<Calendar className="h-4 w-4" />
								<span>{formattedDate}</span>
							</div>
							<div className="flex items-center gap-2">
								<Eye className="h-4 w-4" />
								<span>{blog.view_count} views</span>
							</div>
						</div>
					</header>

					{/* Featured Image */}
					{blog.image && (
						<div className="mb-8 rounded-lg overflow-hidden">
							<img
								src={blog.image}
								alt={blog.title}
								className="w-full h-96 object-cover"
							/>
						</div>
					)}

					{/* Excerpt */}
					{blog.excerpt && (
						<p className="text-lg text-muted-foreground italic mb-8 p-4 border-l-4 border-primary">
							{blog.excerpt}
						</p>
					)}

				{/* Content */}
				<div className="mb-8">
					<MarkdownRenderer content={blog.content} />
				</div>

					{/* Footer */}
					<div className="border-t pt-8 mt-12">
						<div className="bg-muted p-6 rounded-lg">
							<div className="flex gap-4">
								<div>
									<h3 className="font-semibold mb-2">About the Author</h3>
									<p className="text-sm text-muted-foreground">
										{blog.author} - Published on {formattedDate}
									</p>
								</div>
							</div>
						</div>
					</div>

					{/* Back to Blog Button */}
					<div className="mt-12">
						<Button
							onClick={() => navigate("/blog")}
							className="w-full sm:w-auto"
						>
							Back to Blog
						</Button>
					</div>
				</article>
			</div>
		</div>
	);
}
