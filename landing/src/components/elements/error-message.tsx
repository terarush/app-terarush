import { AlertCircle, RotateCcw } from "lucide-react";
import { Button } from "@/components/ui/button";

interface ErrorMessageProps {
	message: string;
	onRetry?: () => void;
	variant?: "default" | "destructive" | "warning";
}

export function ErrorMessage({ message, onRetry, variant = "default" }: ErrorMessageProps) {
	const variantStyles = {
		default: "bg-destructive/10 border-destructive/20 text-destructive",
		destructive: "bg-red-500/10 border-red-500/20 text-red-600",
		warning: "bg-amber-500/10 border-amber-500/20 text-amber-600",
	};

	return (
		<div className={`p-6 rounded-xl border flex items-start gap-4 ${variantStyles[variant]}`}>
			<AlertCircle className="h-5 w-5 flex-shrink-0 mt-0.5" />
			<div className="flex-1">
				<p className="font-medium">{message}</p>
				{onRetry && (
					<Button
						onClick={onRetry}
						variant="outline"
						size="sm"
						className="mt-3 gap-2"
					>
						<RotateCcw className="h-4 w-4" />
						Try Again
					</Button>
				)}
			</div>
		</div>
	);
}
