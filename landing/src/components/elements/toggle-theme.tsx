import { Moon, Sun } from "lucide-react";
import { useTheme } from "@/components/theme-provider";
import { Button } from "@/components/ui/button";

type ThemeToggleProps = {
	variant?: "default" | "rounded";
};

export function ThemeToggle({ variant = "default" }: ThemeToggleProps) {
	const { theme, setTheme } = useTheme();

	const toggleTheme = () => {
		setTheme(theme === "dark" ? "light" : "dark");
	};

	return (
		<Button
			variant={variant === "rounded" ? "ghost" : "outline"}
			size={variant === "rounded" ? "icon" : "default"}
			aria-label="Toggle theme"
			onClick={toggleTheme}
			className={variant === "rounded" ? "rounded-full border p-2" : ""}
		>
			<Sun
				className={`h-[1.2rem] w-[1.2rem] transition-all ${
					theme === "dark"
						? "scale-0 rotate-90 absolute"
						: "scale-100 rotate-0"
				}`}
			/>
			<Moon
				className={`h-[1.2rem] w-[1.2rem] transition-all ${
					theme === "dark"
						? "scale-100 rotate-0"
						: "scale-0 -rotate-90 absolute"
				}`}
			/>
		</Button>
	);
}
