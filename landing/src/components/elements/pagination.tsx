import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";

interface PaginationProps {
	currentPage: number;
	totalPages: number;
	onPageChange: (page: number) => void;
	isLoading?: boolean;
	maxVisiblePages?: number;
}

export function Pagination({
	currentPage,
	totalPages,
	onPageChange,
	isLoading = false,
	maxVisiblePages = 5,
}: PaginationProps) {
	const pageNumbers = Array.from({ length: Math.min(totalPages, maxVisiblePages) }, (_, i) => i + 1);

	return (
		<div className="flex items-center justify-between gap-4">
			<div className="text-sm text-muted-foreground">
				Page {currentPage} of {totalPages}
			</div>
			<div className="flex items-center gap-2">
				<Button
					variant="outline"
					size="sm"
					onClick={() => onPageChange(Math.max(1, currentPage - 1))}
					disabled={currentPage === 1 || isLoading}
					className="gap-1"
				>
					<ChevronLeft className="h-4 w-4" />
					Previous
				</Button>
				<div className="flex gap-1">
					{pageNumbers.map((page) => (
						<Button
							key={page}
							variant={currentPage === page ? "default" : "outline"}
							size="sm"
							onClick={() => onPageChange(page)}
							disabled={isLoading}
							className="h-9 w-9 p-0"
						>
							{page}
						</Button>
					))}
				</div>
				<Button
					variant="outline"
					size="sm"
					onClick={() => onPageChange(Math.min(totalPages, currentPage + 1))}
					disabled={currentPage === totalPages || isLoading}
					className="gap-1"
				>
					Next
					<ChevronRight className="h-4 w-4" />
				</Button>
			</div>
		</div>
	);
}
