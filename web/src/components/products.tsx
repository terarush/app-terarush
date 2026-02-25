import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { getProducts, type Product } from "@/lib/api/products";
import { ProductCard } from "@/components/fragments/ProductCard";
import { Button } from "@/components/ui/button";
import { Loader2, ArrowRight } from "lucide-react";

export default function Products() {
	const [products, setProducts] = useState<Product[]>([]);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		loadProducts();
	}, []);

	const loadProducts = async () => {
		try {
			setLoading(true);
			const response = await getProducts({ is_active: true });
			// Get only first 3 products for landing page
			setProducts((response.products || []).slice(0, 3));
		} catch (err) {
			console.error("Error loading products:", err);
		} finally {
			setLoading(false);
		}
	};

	if (loading) {
		return (
			<section className="py-20 bg-background" id="products">
				<div className="container mx-auto px-4">
					<div className="flex items-center justify-center min-h-[400px]">
						<Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
					</div>
				</div>
			</section>
		);
	}

	if (products.length === 0) {
		return null; // Don't show section if no products
	}

	return (
		<section className="py-20 bg-background" id="products">
			<div className="container mx-auto px-4">
				{/* Section Header */}
				<div className="text-center mb-12">
					<h2 className="text-3xl md:text-4xl font-bold mb-4">
						Our Products
					</h2>
					<p className="text-muted-foreground text-lg max-w-2xl mx-auto">
						Choose the perfect container plan for your Node.js, Python, or
						Ubuntu bot projects
					</p>
				</div>

				{/* Products Grid */}
				<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 mb-12">
					{products.map((product) => (
						<ProductCard key={product.id} product={product} />
					))}
				</div>

				{/* View All Button */}
				<div className="text-center">
					<Button asChild size="lg">
						<Link to="/products">
							View All Products
							<ArrowRight className="ml-2 h-4 w-4" />
						</Link>
					</Button>
				</div>
			</div>
		</section>
	);
}
