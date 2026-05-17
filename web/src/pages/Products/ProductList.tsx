import { useState, useEffect } from "react";
import { getProducts, type Product } from "@/lib/api/products";
import Navbar from "@/components/navbar";
import Footer from "@/components/footer";
import { LoadingSpinner } from "@/components/elements/loading-spinner";
import { ErrorMessage } from "@/components/elements/error-message";
import { Pagination } from "@/components/elements/pagination";
import { ProductCard } from "@/components/fragments/ProductCard";
import {
	PageLayout,
	PageContent,
	HeroSection,
	ContentSection,
} from "@/components/layouts/page-layout";

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
			<PageLayout>
				<Navbar />
				<PageContent className="flex items-center justify-center">
					<LoadingSpinner size="lg" text="Loading products..." />
				</PageContent>
				<Footer />
			</PageLayout>
		);
	}

	return (
		<PageLayout>
			<Navbar />

			<PageContent>
				{/* Hero Section */}
				<HeroSection>
					<div className="max-w-3xl mx-auto text-center">
						<h1 className="text-4xl md:text-5xl font-bold mb-4">
							Our Products
						</h1>
						<p className="text-lg text-muted-foreground leading-relaxed">
							Choose the perfect container plan for your needs
						</p>
					</div>
				</HeroSection>

				{/* Products Grid */}
				<ContentSection>
					{error && (
						<div className="mb-8">
							<ErrorMessage message={error} onRetry={loadProducts} />
						</div>
					)}

					{products.length === 0 ? (
						<div className="text-center py-16">
							<p className="text-muted-foreground text-lg">
								No products available at the moment.
							</p>
						</div>
					) : (
						<>
							<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
								{products.map((product) => (
									<ProductCard key={product.id} product={product} />
								))}
							</div>

							{/* Pagination */}
							{totalPages > 1 && (
								<div className="flex justify-center pt-8 border-t border-border/50">
									<Pagination
										currentPage={page}
										totalPages={totalPages}
										onPageChange={setPage}
										isLoading={loading}
									/>
								</div>
							)}
						</>
					)}
				</ContentSection>
			</PageContent>

			<Footer />
		</PageLayout>
	);
}
