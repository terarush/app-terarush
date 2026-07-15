import { Calendar, Eye, User, ArrowRight } from "lucide-react";
import { getImageUrl, type Blog } from "@/lib/api/blogs";

interface BlogCardProps {
	blog: Blog;
	onClick?: () => void;
}

export function BlogCard({ blog, onClick }: BlogCardProps) {
	return (
		<article
			onClick={onClick}
			className="group relative rounded-2xl overflow-hidden bg-card border border-border/50 hover:border-primary/30 shadow-sm hover:shadow-lg transition-all duration-300 cursor-pointer hover:-translate-y-1 h-full flex flex-col"
		>
			{/* Image Container */}
			<div className="relative h-48 overflow-hidden bg-muted">
				{blog.image && (
					<img
						src={getImageUrl(blog.image)}
						alt={blog.title}
						className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
					/>
				)}
				<div className="absolute inset-0 bg-gradient-to-t from-background via-transparent to-transparent opacity-60" />

				{/* Category Badge */}
				{blog.category && (
					<div className="absolute top-4 left-4">
						<span className="px-3 py-1 text-xs font-semibold rounded-full bg-primary/90 text-primary-foreground backdrop-blur-sm">
							{blog.category}
						</span>
					</div>
				)}
			</div>

			{/* Content */}
			<div className="p-6 relative z-10 flex-1 flex flex-col justify-between">
				<div>
					<h3 className="text-lg font-bold mb-2 line-clamp-2 group-hover:text-primary transition-colors">
						{blog.title}
					</h3>
					<p className="text-sm text-muted-foreground line-clamp-2 mb-4">
						{blog.excerpt}
					</p>
				</div>

				{/* Meta Info */}
				<div className="space-y-3 pt-4 border-t border-border/30">
					<div className="flex items-center gap-4 text-xs text-muted-foreground flex-wrap">
						<div className="flex items-center gap-1">
							<User className="h-3.5 w-3.5" />
							<span className="truncate">{blog.author}</span>
						</div>
						<div className="flex items-center gap-1">
							<Calendar className="h-3.5 w-3.5" />
							<span>
								{new Date(blog.created_at).toLocaleDateString('en-US', {
									month: 'short',
									day: 'numeric',
								})}
							</span>
						</div>
						<div className="flex items-center gap-1 ml-auto">
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
			</div>
		</article>
	);
}
