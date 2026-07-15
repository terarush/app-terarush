import { Search, X } from "lucide-react";

interface SearchInputProps {
	placeholder: string;
	value: string;
	onChange: (value: string) => void;
	onClear?: () => void;
}

export function SearchInput({ placeholder, value, onChange, onClear }: SearchInputProps) {
	return (
		<div className="relative">
			<Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
			<input
				type="text"
				placeholder={placeholder}
				value={value}
				onChange={(e) => onChange(e.target.value)}
				className="w-full pl-10 pr-10 py-2 border border-border/50 rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all"
			/>
			{value && (
				<button
					onClick={() => {
						onChange("");
						onClear?.();
					}}
					className="absolute right-2 top-1/2 -translate-y-1/2 p-1 hover:bg-muted rounded-md transition-colors"
					title="Clear search"
				>
					<X className="h-4 w-4 text-muted-foreground hover:text-foreground" />
				</button>
			)}
		</div>
	);
}
