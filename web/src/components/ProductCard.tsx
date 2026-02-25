import { Link } from "react-router-dom";
import type { Product } from "@/lib/api/products";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

interface ProductCardProps {
	product: Product;
	showAdminActions?: boolean;
	onEdit?: (product: Product) => void;
	onDelete?: (product: Product) => void;
}

export function ProductCard({
	product,
	showAdminActions = false,
	onEdit,
	onDelete,
}: ProductCardProps) {
	const formatPrice = (price: number) => {
		return new Intl.NumberFormat("id-ID", {
			style: "currency",
			currency: "IDR",
			minimumFractionDigits: 0,
		}).format(price);
	};

	const formatStorage = (mb: number) => {
		if (mb >= 1024) {
			return `${(mb / 1024).toFixed(0)} GB`;
		}
		return `${mb} MB`;
	};

	return (
		<Card className="h-full flex flex-col">
			<CardHeader>
				<div className="flex items-start justify-between">
					<div className="flex-1">
						<CardTitle className="text-xl">{product.name}</CardTitle>
						<CardDescription className="mt-2">
							{product.description}
						</CardDescription>
					</div>
					{!product.is_active && (
						<Badge variant="secondary" className="ml-2">
							Inactive
						</Badge>
					)}
				</div>
			</CardHeader>

			<CardContent className="flex-1">
				<div className="space-y-3">
					<div className="grid grid-cols-2 gap-4 text-sm">
						<div>
							<p className="text-muted-foreground">CPU Cores</p>
							<p className="font-medium">{product.cpu_cores} Cores</p>
						</div>
						<div>
							<p className="text-muted-foreground">RAM</p>
							<p className="font-medium">{formatStorage(product.ram_mb)}</p>
						</div>
						<div>
							<p className="text-muted-foreground">Storage</p>
							<p className="font-medium">{product.storage_gb} GB</p>
						</div>
						<div>
							<p className="text-muted-foreground">Bandwidth</p>
							<p className="font-medium">{product.bandwidth_gb} GB</p>
						</div>
					</div>

					<div className="pt-2 border-t">
						<p className="text-muted-foreground text-sm">Stock</p>
						<p className="font-medium">
							{product.stock === -1 ? "Unlimited" : `${product.stock} available`}
						</p>
					</div>

					<div className="pt-2">
						<p className="text-2xl font-bold">{formatPrice(product.price)}</p>
						<p className="text-muted-foreground text-sm">per month</p>
					</div>
				</div>
			</CardContent>

			<CardFooter className="flex-col gap-2">
				{showAdminActions ? (
					<div className="flex gap-2 w-full">
						<Button
							variant="outline"
							className="flex-1"
							onClick={() => onEdit?.(product)}
						>
							Edit
						</Button>
						<Button
							variant="destructive"
							className="flex-1"
							onClick={() => onDelete?.(product)}
						>
							Delete
						</Button>
					</div>
				) : (
					<Button asChild className="w-full">
						<Link to={`/products/${product.id}`}>View Details</Link>
					</Button>
				)}
			</CardFooter>
		</Card>
	);
}
