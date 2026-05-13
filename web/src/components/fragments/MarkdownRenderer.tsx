import { marked } from "marked";
import "./MarkdownRenderer.css";

interface MarkdownRendererProps {
	content: string;
}

export function MarkdownRenderer({ content }: MarkdownRendererProps) {
	let htmlContent = "";
	try {
		htmlContent = marked.parse(content) as string;
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

