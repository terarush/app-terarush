import { ChevronDown } from "lucide-react";
import { Button } from "@/components/ui/button";

interface LoadMoreButtonProps {
	onClick: () => void;
	isLoading?: boolean;
	hasMore: boolean;
	itemCount: number;
	pageSize: number;
}

export function LoadMoreButton({
	onClick,
	isLoading = false,
	hasMore,
	itemCount,
	pageSize,
}: LoadMoreButtonProps) {
	if (!hasMore) {
		return null;
	}

	return (
		<div className="flex flex-col items-center gap-4 py-12">
			<p className="text-sm text-muted-foreground">
				Showing {itemCount} posts
			</p>
			<Button
				onClick={onClick}
				disabled={isLoading}
				className="gap-2 px-6 py-2"
				variant="outline"
			>
				{isLoading ? (
					<>
						<div className="h-4 w-4 animate-spin rounded-full border-2 border-muted-foreground border-t-foreground" />
						Loading...
					</>
				) : (
					<>
						Load More
						<ChevronDown className="h-4 w-4" />
					</>
				)}
			</Button>
		</div>
	);
}
