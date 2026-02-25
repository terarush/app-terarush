import { useState, useEffect } from "react";
import { getProducts, type Product } from "@/lib/api/products";
import { ProductCard } from "@/components/fragments/ProductCard";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";

export function ProductList() {
	const [products, setProducts] = useState<Product[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [page, setPage] = useState(1);
	const [total, setTotal] = useState(0);
	const pageSize = 9;

	useEffect(() => {
		loadProducts();
	}, [page]);

	const loadProducts = async () => {
		try {
			setLoading(true);
			setError(null);
			const response = await getProducts({
				page,
				page_size: pageSize,
				is_active: true, // Only show active products to users
			});
			setProducts(response.products || []);
			setTotal(response.total || 0);
		} catch (err) {
			console.error("Error loading products:", err);
			setError("Failed to load products. Please try again later.");
		} finally {
			setLoading(false);
		}
	};

	const totalPages = Math.ceil(total / pageSize);

	if (loading && products.length === 0) {
		return (
			<div className="container mx-auto px-4 py-8">
				<div className="flex items-center justify-center min-h-[400px]">
					<Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div className="container mx-auto px-4 py-8">
				<div className="text-center py-12">
					<p className="text-destructive mb-4">{error}</p>
					<Button onClick={loadProducts}>Try Again</Button>
				</div>
			</div>
		);
	}

	return (
		<div className="container mx-auto px-4 py-8">
			<div className="mb-8">
				<h1 className="text-4xl font-bold mb-2">Our Products</h1>
				<p className="text-muted-foreground text-lg">
					Choose the perfect container plan for your needs
				</p>
			</div>

			{products.length === 0 ? (
				<div className="text-center py-12">
					<p className="text-muted-foreground text-lg">
						No products available at the moment.
					</p>
				</div>
			) : (
				<>
					<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
						{products.map((product) => (
							<ProductCard key={product.id} product={product} />
						))}
					</div>

					{totalPages > 1 && (
						<div className="flex items-center justify-center gap-2 mt-8">
							<Button
								variant="outline"
								onClick={() => setPage((p) => Math.max(1, p - 1))}
								disabled={page === 1 || loading}
							>
								Previous
							</Button>
							<span className="text-sm text-muted-foreground px-4">
								Page {page} of {totalPages}
							</span>
							<Button
								variant="outline"
								onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
								disabled={page === totalPages || loading}
							>
								Next
							</Button>
						</div>
					)}
				</>
			)}
		</div>
	);
}
