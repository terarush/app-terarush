import { useState, useEffect, useRef } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { getBlogBySlug, getImageUrl, type Blog } from "@/lib/api/blogs";
import {
	getCommentsByPost,
	createComment,
	createReply,
	updateComment,
	deleteComment,
	type Comment,
	type CommentListResponse,
} from "@/lib/api/comments";
import { Button } from "@/components/ui/button";
import {
	ArrowLeft,
	Calendar,
	User,
	Eye,
	Link as LinkIcon,
	Twitter,
	Facebook,
	ChevronRight,
	Bookmark,
	MessageCircle,
	Tag,
} from "lucide-react";
import { MarkdownRenderer } from "@/components/fragments/MarkdownRenderer";
import { CommentForm } from "@/components/fragments/CommentForm";
import { motion, useScroll, useSpring } from "framer-motion";
import { toast } from "sonner";
import Navbar from "@/components/navbar";
import Footer from "@/components/footer";
import { CommentList } from "@/components/fragments/CommentList";

export function BlogDetail() {
	const { slug } = useParams<{ slug: string }>();
	const navigate = useNavigate();
	const [blog, setBlog] = useState<Blog | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [headings, setHeadings] = useState<
		Array<{ id: string; text: string; level: number }>
	>([]);

	// Comment state
	const [comments, setComments] = useState<Comment[]>([]);
	const [commentsLoading, setCommentsLoading] = useState(false);
	const [commentSubmitting, setCommentSubmitting] = useState(false);
	const commentSectionRef = useRef<HTMLDivElement>(null);

	const { scrollYProgress } = useScroll();
	const scaleX = useSpring(scrollYProgress, {
		stiffness: 100,
		damping: 30,
		restDelta: 0.001,
	});

	useEffect(() => {
		if (slug) {
			loadBlog();
		}
		window.scrollTo(0, 0);
	}, [slug]);

	const loadBlog = async () => {
		try {
			setLoading(true);
			setError(null);
			const data = await getBlogBySlug(slug!);
			setBlog(data);
			setHeadings(extractHeadings(data.content));
			// Load comments after blog is loaded
			await loadComments(data.id);
		} catch (err) {
			console.error("Error loading blog:", err);
			setError("Failed to load blog post. Please try again later.");
		} finally {
			setLoading(false);
		}
	};

	const loadComments = async (postId: number) => {
		try {
			setCommentsLoading(true);
			const data: CommentListResponse = await getCommentsByPost(postId);
			setComments(data.comments || []);
		} catch (err) {
			console.error("Error loading comments:", err);
		} finally {
			setCommentsLoading(false);
		}
	};

	const handleAddComment = async (data: any) => {
		if (!blog) return;
		try {
			setCommentSubmitting(true);
			const newComment = await createComment({
				content: data.content,
				post_id: blog.id,
			});
			setComments((prev) => [newComment, ...prev]);
			toast.success("Comment posted!");
		} catch (err) {
			console.error("Error creating comment:", err);
			toast.error("Failed to post comment. Please try again.");
		} finally {
			setCommentSubmitting(false);
		}
	};

	const handleAddReply = async (data: any) => {
		if (!blog) return;
		try {
			setCommentSubmitting(true);
			const newReply = await createReply({
				content: data.content,
				post_id: blog.id,
				parent_id: data.parent_id,
			});
			// Re-fetch to get nested replies correctly
			await loadComments(blog.id);
			toast.success("Reply posted!");
			return newReply;
		} catch (err) {
			console.error("Error posting reply:", err);
			toast.error("Failed to post reply. Please try again.");
		} finally {
			setCommentSubmitting(false);
		}
	};

	const handleEditComment = async (id: string, data: any) => {
		if (!blog) return;
		try {
			setCommentSubmitting(true);
			const updated = await updateComment(id, {
				content: data.content,
				post_id: blog.id,
			});
			setComments((prev) =>
				prev.map((c) =>
					c.id === id
						? {
								...c,
								content: updated.content,
								updated_at: updated.updated_at,
							}
						: c,
				),
			);
			toast.success("Comment updated!");
		} catch (err) {
			console.error("Error updating comment:", err);
			toast.error("Failed to update comment. Please try again.");
		} finally {
			setCommentSubmitting(false);
		}
	};

	const handleDeleteComment = async (id: string) => {
		try {
			setCommentSubmitting(true);
			await deleteComment(id);
			setComments((prev) => prev.filter((c) => c.id !== id));
			toast.success("Comment deleted.");
		} catch (err) {
			console.error("Error deleting comment:", err);
			toast.error("Failed to delete comment. Please try again.");
		} finally {
			setCommentSubmitting(false);
		}
	};

	const estimateReadingTime = (content: string) => {
		const wordsPerMinute = 200;
		const words = content.trim().split(/\s+/).length;
		return Math.ceil(words / wordsPerMinute);
	};

	const extractHeadings = (markdown: string) => {
		const headingRegex = /^(#{1,6})\s+(.+)$/gm;
		const headingsList: Array<{ id: string; text: string; level: number }> =
			[];
		let match;

		while ((match = headingRegex.exec(markdown)) !== null) {
			const level = match[1].length;
			const text = match[2];
			const id = text
				.toLowerCase()
				.replace(/[^\w\s-]/g, "")
				.replace(/\s+/g, "-")
				.replace(/-+/g, "-");

			headingsList.push({ id, text, level });
		}

		return headingsList;
	};

	const scrollToHeading = (id: string) => {
		const element = document.getElementById(id);
		if (element) {
			element.scrollIntoView({ behavior: "smooth" });
		}
	};

	const copyToClipboard = () => {
		navigator.clipboard.writeText(window.location.href);
		toast.success("Link copied to clipboard!");
	};

	if (loading) {
		return (
			<div className="min-h-screen bg-background flex flex-col">
				<Navbar />
				<div className="flex-1 flex flex-col items-center justify-center pt-24 space-y-4">
					<div className="w-16 h-16 border-4 border-primary/20 border-t-primary rounded-full animate-spin" />
					<p className="text-muted-foreground animate-pulse font-medium">
						Preparing your article...
					</p>
				</div>
				<Footer />
			</div>
		);
	}

	if (error || !blog) {
		return (
			<div className="min-h-screen bg-background flex flex-col">
				<Navbar />
				<main className="flex-1 pt-32">
					<div className="container mx-auto px-4 text-center">
						<div className="inline-flex items-center justify-center w-20 h-20 rounded-xl bg-destructive/10 mb-6 text-destructive">
							<ArrowLeft className="h-10 w-10" />
						</div>
						<h1 className="text-3xl font-bold mb-4">
							{error || "Blog Post Not Found"}
						</h1>
						<p className="text-muted-foreground mb-8 max-w-md mx-auto">
							The article you're looking for might have been moved
							or deleted.
						</p>
						<Button
							onClick={() => navigate("/blog")}
							className="rounded-xl px-8 py-6"
						>
							<ArrowLeft className="mr-2 h-4 w-4" />
							Browse Other Articles
						</Button>
					</div>
				</main>
				<Footer />
			</div>
		);
	}

	const formattedDate = new Date(blog.created_at).toLocaleDateString(
		"en-US",
		{
			year: "numeric",
			month: "long",
			day: "numeric",
		},
	);

	const readingTime = estimateReadingTime(blog.content);

	return (
		<div className="min-h-screen bg-background flex flex-col selection:bg-primary/30 selection:text-primary">
			<Navbar />

			{/* Reading Progress Bar */}
			<motion.div
				className="fixed top-0 left-0 right-0 h-1.5 bg-primary z-[60] origin-left"
				style={{ scaleX }}
			/>

			<main className="flex-1">
				{/* Hero Section */}
				<div className="relative w-full h-[50vh] sm:h-[60vh] md:h-[70vh] min-h-[300px] overflow-hidden">
					{blog.image ? (
						<motion.div
							initial={{ scale: 1.1, opacity: 0 }}
							animate={{ scale: 1, opacity: 1 }}
							transition={{ duration: 1.5 }}
							className="w-full h-full"
						>
							<img
								src={getImageUrl(blog.image)}
								alt={blog.title}
								className="w-full h-full object-cover"
							/>
						</motion.div>
					) : (
						<div className="w-full h-full bg-gradient-to-br from-primary/20 via-accent/10 to-background" />
					)}

					{/* Overlays */}
					<div className="absolute inset-0 bg-gradient-to-t from-background via-background/40 to-transparent" />
					<div className="absolute inset-0 bg-black/20" />

					{/* Content Overlay */}
					<div className="absolute inset-0 flex flex-col justify-end">
						<div className="container mx-auto px-4 pb-12 sm:pb-20 md:pb-32">
							<motion.div
								initial={{ opacity: 0, y: 30 }}
								animate={{ opacity: 1, y: 0 }}
								transition={{ delay: 0.5, duration: 0.8 }}
								className="max-w-4xl"
							>
								{/* Breadcrumbs */}
								<nav className="hidden sm:flex items-center gap-2 text-white/70 text-xs sm:text-sm font-medium mb-4 sm:mb-6">
									<Link
										to="/"
										className="hover:text-white transition-colors"
									>
										Home
									</Link>
									<ChevronRight className="h-3 w-3 sm:h-4 sm:w-4" />
									<Link
										to="/blog"
										className="hover:text-white transition-colors"
									>
										Blog
									</Link>
									<ChevronRight className="h-3 w-3 sm:h-4 sm:w-4" />
									<span className="text-white truncate max-w-[200px]">
										{blog.category || "Article"}
									</span>
								</nav>

								<div className="flex flex-wrap gap-2 sm:gap-3 mb-4 sm:mb-6">
									{blog.category && (
										<span className="px-3 sm:px-4 py-1 sm:py-1.5 rounded-lg sm:rounded-xl bg-primary text-primary-foreground text-xs font-bold uppercase tracking-widest shadow-lg">
											{blog.category}
										</span>
									)}
									<span className="px-3 sm:px-4 py-1 sm:py-1.5 rounded-lg sm:rounded-xl bg-white/10 backdrop-blur-md text-white text-xs font-bold uppercase tracking-widest border border-white/20">
										{readingTime} min read
									</span>
								</div>

								<h1 className="text-2xl sm:text-4xl md:text-6xl lg:text-7xl font-bold text-white mb-4 sm:mb-8 leading-[1.1] tracking-tight drop-shadow-2xl">
									{blog.title}
								</h1>

								<div className="hidden sm:flex flex-wrap items-center gap-6 sm:gap-8 text-white/90">
									<div className="flex items-center gap-3">
										<div className="w-10 sm:w-12 h-10 sm:h-12 rounded-full bg-white/20 backdrop-blur-md flex items-center justify-center border border-white/30 overflow-hidden">
											{blog.user?.avatar ? (
												<img
													src={blog.user.avatar}
													alt={blog.user.name}
													className="w-full h-full object-cover"
												/>
											) : (
												<User className="h-5 sm:h-6 w-5 sm:w-6" />
											)}
										</div>
										<div>
											<p className="text-xs sm:text-sm font-bold">
												{blog.author}
											</p>
											<p className="text-xs text-white/60 line-clamp-1">
												{blog.user?.bio
													? blog.user.bio.split(
															"\n",
														)[0]
													: "Contributor"}
											</p>
										</div>
									</div>
									<div className="flex items-center gap-4 sm:gap-6">
										<div className="flex items-center gap-2">
											<Calendar className="h-4 w-4 text-primary" />
											<span className="text-xs sm:text-sm font-medium">
												{formattedDate}
											</span>
										</div>
										<div className="flex items-center gap-2">
											<Eye className="h-4 w-4 text-primary" />
											<span className="text-xs sm:text-sm font-medium">
												{blog.view_count} views
											</span>
										</div>
									</div>
								</div>
							</motion.div>
						</div>
					</div>
				</div>

				{/* Article Body */}
				<div className="container mx-auto px-4 -mt-12 sm:-mt-16 relative z-20 pb-12 sm:pb-20">
					<div className="grid grid-cols-1 lg:grid-cols-12 gap-6 sm:gap-12">
						{/* Sidebar / Social Share */}
						<aside className="hidden lg:block lg:col-span-1">
							<div className="sticky top-32 flex flex-col items-center gap-6">
								<button
									onClick={copyToClipboard}
									className="p-3 rounded-xl bg-card border border-border/50 hover:border-primary/50 hover:text-primary transition-all shadow-sm"
									title="Copy link"
								>
									<LinkIcon className="h-5 w-5" />
								</button>
								<button className="p-3 rounded-xl bg-card border border-border/50 hover:border-[#1DA1F2] hover:text-[#1DA1F2] transition-all shadow-sm">
									<Twitter className="h-5 w-5" />
								</button>
								<button className="p-3 rounded-xl bg-card border border-border/50 hover:border-[#4267B2] hover:text-[#4267B2] transition-all shadow-sm">
									<Facebook className="h-5 w-5" />
								</button>
								<div className="h-20 w-px bg-border/50" />
								<button className="p-3 rounded-xl bg-card border border-border/50 hover:text-primary transition-all shadow-sm">
									<Bookmark className="h-5 w-5" />
								</button>
							</div>
						</aside>

						{/* Main Content */}
						<div className="lg:col-span-8">
							<article className="bg-card border border-border/40 rounded-lg sm:rounded-xl p-6 sm:p-8 md:p-16 shadow-xl shadow-primary/5">
								{/* Excerpt/Introduction */}
								{blog.excerpt && (
									<div className="mb-8 sm:mb-12">
										<p className="text-lg sm:text-2xl md:text-3xl text-foreground/80 leading-relaxed font-bold border-l-4 border-primary pl-6 sm:pl-8 py-2">
											{blog.excerpt}
										</p>
									</div>
								)}

								{/* Markdown Content */}
								<div
									className="prose prose-invert max-w-none
									prose-headings:font-bold prose-headings:tracking-tight
									prose-h1:text-2xl sm:prose-h1:text-3xl md:prose-h1:text-4xl
									prose-h2:text-xl sm:prose-h2:text-2xl md:prose-h2:text-3xl prose-h2:mt-8 sm:prose-h2:mt-12 md:prose-h2:mt-16 prose-h2:mb-4 sm:prose-h2:mb-6 md:prose-h2:mb-8
									prose-h3:text-lg sm:prose-h3:text-xl
									prose-p:text-sm sm:prose-p:text-base text-muted-foreground prose-p:leading-[1.8] prose-p:mb-4 sm:prose-p:mb-6 md:prose-p:mb-8
									prose-strong:text-foreground prose-strong:font-bold
									prose-a:text-primary prose-a:no-underline hover:prose-a:underline
									prose-img:rounded-lg sm:prose-img:rounded-xl prose-img:shadow-2xl
									prose-code:text-xs sm:prose-code:text-sm prose-code:text-primary prose-code:bg-primary/5 prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded
									prose-pre:bg-muted/30 prose-pre:border prose-pre:border-border/50 prose-pre:rounded-xl prose-pre:text-xs sm:prose-pre:text-sm prose-pre:p-4
									prose-blockquote:border-primary prose-blockquote:bg-primary/5 prose-blockquote:rounded-xl prose-blockquote:py-1 prose-blockquote:text-sm sm:prose-blockquote:text-base prose-blockquote:pl-4 sm:prose-blockquote:pl-6
									prose-li:text-sm sm:prose-li:text-base prose-li:text-muted-foreground
								"
								>
									<MarkdownRenderer content={blog.content} />
								</div>

								{/* Tags */}
								{blog.tags && (
									<div className="mt-12 sm:mt-16 pt-8 sm:pt-10 border-t border-border/40">
										<div className="flex items-center gap-3 mb-4 sm:mb-6">
											<Tag className="h-4 sm:h-5 w-4 sm:w-5 text-primary" />
											<h3 className="text-base sm:text-lg font-bold">
												Related Topics
											</h3>
										</div>
										<div className="flex flex-wrap gap-2">
											{blog.tags
												.split(",")
												.map((tag, i) => (
													<span
														key={i}
														className="px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg sm:rounded-xl bg-muted/50 border border-border/50 text-xs sm:text-sm font-medium hover:bg-primary/10 hover:border-primary/30 cursor-pointer transition-all"
													>
														#{tag.trim()}
													</span>
												))}
										</div>
									</div>
								)}

								{/* Footer Actions */}
								<div className="mt-12 sm:mt-16 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-6 py-6 sm:py-8 border-t border-border/40">
									<div className="w-full sm:w-auto">
										<span className="text-xs font-bold text-muted-foreground uppercase tracking-widest block mb-3 sm:mb-0">
											Share Article
										</span>
										<div className="flex gap-2">
											<Button
												variant="outline"
												size="icon"
												className="rounded-lg sm:rounded-xl h-9 sm:h-10 w-9 sm:w-10"
											>
												<Twitter className="h-3.5 sm:h-4 w-3.5 sm:w-4" />
											</Button>
											<Button
												variant="outline"
												size="icon"
												className="rounded-lg sm:rounded-xl h-9 sm:h-10 w-9 sm:w-10"
											>
												<Facebook className="h-3.5 sm:h-4 w-3.5 sm:w-4" />
											</Button>
											<Button
												variant="outline"
												size="icon"
												className="rounded-lg sm:rounded-xl h-9 sm:h-10 w-9 sm:w-10"
												onClick={copyToClipboard}
											>
												<LinkIcon className="h-3.5 sm:h-4 w-3.5 sm:w-4" />
											</Button>
										</div>
									</div>
									<div className="w-full sm:w-auto flex items-center gap-2 sm:gap-4">
										<Button
											variant="ghost"
											className="rounded-lg sm:rounded-xl gap-2 text-xs sm:text-sm flex-1 sm:flex-none"
											onClick={() =>
												commentSectionRef.current?.scrollIntoView(
													{ behavior: "smooth" },
												)
											}
										>
											<MessageCircle className="h-4 sm:h-5 w-4 sm:w-5" />
											<span className="hidden sm:inline">
												Comments
											</span>
											{comments.length > 0 && (
												<span className="ml-1 text-xs bg-primary/10 text-primary rounded-full px-1.5 py-0.5 font-bold">
													{comments.length}
												</span>
											)}
										</Button>
										<Button
											variant="ghost"
											className="rounded-lg sm:rounded-xl gap-2 text-xs sm:text-sm flex-1 sm:flex-none"
										>
											<Bookmark className="h-4 sm:h-5 w-4 sm:w-5" />
											<span className="hidden sm:inline">
												Save
											</span>
										</Button>
									</div>
								</div>
							</article>
						</div>

						{/* Right Sidebar - Newsletter/Related */}
						<aside className="lg:col-span-3 space-y-6 sm:space-y-8">
							<div className="p-6 sm:p-8 rounded-lg sm:rounded-xl bg-card border border-border/40 shadow-lg">
								<h3 className="text-lg sm:text-xl font-bold mb-3 sm:mb-4">
									Stay Updated
								</h3>
								<p className="text-xs sm:text-sm text-muted-foreground mb-4 sm:mb-6">
									Get the latest articles and insights
									delivered to your inbox.
								</p>
								<div className="space-y-2 sm:space-y-3">
									<input
										type="email"
										placeholder="your@email.com"
										className="w-full px-4 py-2.5 sm:py-3 rounded-lg sm:rounded-xl bg-muted/50 border border-border/50 focus:border-primary/50 outline-none text-xs sm:text-sm transition-all"
									/>
									<Button className="w-full rounded-lg sm:rounded-xl py-5 sm:py-6 font-bold text-xs sm:text-sm">
										Subscribe
									</Button>
								</div>
							</div>

							{headings.length > 0 && (
								<div className="space-y-4 sm:space-y-6 p-6 sm:p-8 rounded-lg sm:rounded-xl bg-card border border-border/40">
									<h3 className="text-lg sm:text-xl font-bold px-0 sm:px-2">
										Table of Contents
									</h3>
									<nav className="space-y-2 sm:space-y-3 px-0 sm:px-4 border-l-2 border-border/30">
										{headings.map((heading) => (
											<button
												key={heading.id}
												onClick={() =>
													scrollToHeading(heading.id)
												}
												className={`block text-xs sm:text-sm font-medium transition-colors hover:text-primary text-left ${
													heading.level === 2
														? "text-primary"
														: "text-muted-foreground hover:text-foreground"
												}`}
												style={{
													paddingLeft: `${(heading.level - 2) * 12}px`,
												}}
											>
												{heading.text}
											</button>
										))}
									</nav>
								</div>
							)}
						</aside>
					</div>
				</div>
				{/* Comments Section - full width inside container */}
				<div
					ref={commentSectionRef}
					id="comments"
					className="mt-8 sm:mt-12"
				>
					<div className="bg-card border border-border/40 rounded-lg sm:rounded-xl p-6 sm:p-8 shadow-xl shadow-primary/5">
						<div className="flex items-center gap-3 mb-6">
							<MessageCircle className="h-5 w-5 text-primary" />
							<h2 className="text-xl font-bold">
								Comments
								{comments.length > 0 && (
									<span className="ml-2 text-sm font-normal text-muted-foreground">
										({comments.length})
									</span>
								)}
							</h2>
						</div>

						{/* New Comment Form */}
						<div className="mb-8 pb-8 border-b border-border/40">
							<CommentForm
								postId={blog.id}
								onSubmit={handleAddComment}
								isSubmitting={commentSubmitting}
							/>
						</div>

						{/* Comment List */}
						<CommentList
							comments={comments}
							onAddReply={handleAddReply}
							onEditComment={handleEditComment}
							onDeleteComment={handleDeleteComment}
							isSubmitting={commentSubmitting}
							isLoading={commentsLoading}
						/>
					</div>
				</div>
			</main>

			<Footer />
		</div>
	);
}
