import { marked } from "marked";
import "./MarkdownRenderer.css";

interface MarkdownRendererProps {
	content: string;
}

export function MarkdownRenderer({ content }: MarkdownRendererProps) {
	let htmlContent = "";
	try {
		htmlContent = marked.parse(content) as string;

		// Add IDs to headings for table of contents linking
		const parser = new DOMParser();
		const doc = parser.parseFromString(htmlContent, "text/html");
		const headings = doc.querySelectorAll("h1, h2, h3, h4, h5, h6");

		headings.forEach((heading) => {
			const text = heading.textContent || "";
			const id = text
				.toLowerCase()
				.replace(/[^\w\s-]/g, "")
				.replace(/\s+/g, "-")
				.replace(/-+/g, "-");
			heading.id = id;
		});

		htmlContent = doc.body.innerHTML;
	} catch (error) {
		return (
			<div className="text-destructive p-4 bg-destructive/10 rounded">
				Error rendering content
			</div>
		);
	}

	return (
		<div className="markdown-content max-w-none">
			<div
				className="prose dark:prose-invert max-w-none"
				dangerouslySetInnerHTML={{ __html: htmlContent }}
			/>
		</div>
	);
}

