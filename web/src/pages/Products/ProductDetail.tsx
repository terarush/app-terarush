import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getProductById, type Product } from "@/lib/api/products";
import { checkout } from "@/lib/api/transactions";
import { useAuth } from "@/hooks/useAuth";
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
import {
	Loader2,
	ArrowLeft,
	Cpu,
	HardDrive,
	Gauge,
	Network,
} from "lucide-react";

export function ProductDetail() {
	const { id } = useParams<{ id: string }>();
	const navigate = useNavigate();
	const { isAuthenticated } = useAuth();
	const [product, setProduct] = useState<Product | null>(null);
	const [loading, setLoading] = useState(true);
	const [checkoutLoading, setCheckoutLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [quantity, setQuantity] = useState(1);
	const [selectedVersion, setSelectedVersion] = useState<string>("");

	useEffect(() => {
		if (id) {
			loadProduct();
		}
	}, [id]);

	useEffect(() => {
		// Set default version when product loads
		if (product?.image_tags && product.image_tags.length > 0) {
			setSelectedVersion(product.image_tags[0]);
		}
	}, [product]);

	const loadProduct = async () => {
		try {
			setLoading(true);
			setError(null);
			const data = await getProductById(Number(id));
			setProduct(data);
		} catch (err: any) {
			console.error("Error loading product:", err);
			setError(
				err.response?.data?.message ||
					"Failed to load product details.",
			);
		} finally {
			setLoading(false);
		}
	};

	const handleCheckout = async () => {
		if (!isAuthenticated) {
			navigate("/login", { state: { from: `/products/${id}` } });
			return;
		}

		if (!product) return;

		try {
			setCheckoutLoading(true);
			const response = await checkout({
				product_id: product.id,
				quantity,
				image_tag: selectedVersion || "latest",
		});

		// Transaction created successfully, navigate to transactions page
		navigate("/dashboard/transactions", {
			state: {
				message:
					"Transaction created! Your order is being processed.",
				orderId: response.transaction?.order_id,
			},
		});
		} catch (err: any) {
			console.error("Error during checkout:", err);
			alert(
				err.response?.data?.message ||
					"Failed to process checkout. Please try again.",
			);
			setCheckoutLoading(false);
		}
	};

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

	if (loading) {
		return (
			<div className="container mx-auto px-4 py-8">
				<div className="flex items-center justify-center min-h-[400px]">
					<Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
				</div>
			</div>
		);
	}

	if (error || !product) {
		return (
			<div className="container mx-auto px-4 py-8">
				<div className="text-center py-12">
					<p className="text-destructive mb-4">
						{error || "Product not found"}
					</p>
					<Button onClick={() => navigate("/products")}>
						Back to Products
					</Button>
				</div>
			</div>
		);
	}

	const totalPrice = product.price * quantity;
	const isOutOfStock = product.stock !== -1 && product.stock < quantity;

	return (
		<div className="container mx-auto px-4 py-8">
			<Button
				variant="ghost"
				onClick={() => navigate("/products")}
				className="mb-6"
			>
				<ArrowLeft className="mr-2 h-4 w-4" />
				Back to Products
			</Button>

			<div className="grid md:grid-cols-2 gap-8">
				<Card>
					<CardHeader>
						<div className="flex items-start justify-between">
							<div>
								<CardTitle className="text-3xl">
									{product.name}
								</CardTitle>
								<CardDescription className="mt-2 text-base">
									{product.description}
								</CardDescription>
							</div>
							{!product.is_active && (
								<Badge variant="secondary">Inactive</Badge>
							)}
						</div>
					</CardHeader>

					<CardContent>
						<div className="space-y-6">
							{/* Docker Image Information */}
							{product.docker_image && (
								<div className="pb-4 border-b">
									<p className="text-sm text-muted-foreground mb-1">
										Docker Image
									</p>
									<p className="text-lg font-semibold">
										{product.docker_image}
									</p>
								</div>
							)}

							<div className="grid grid-cols-2 gap-6">
								<div className="flex items-start gap-3">
									<div className="p-2 bg-primary/10 rounded-lg">
										<Cpu className="h-5 w-5 text-primary" />
									</div>
									<div>
										<p className="text-sm text-muted-foreground">
											CPU Cores
										</p>
										<p className="text-lg font-semibold">
											{product.cpu_cores} Cores
										</p>
									</div>
								</div>

								<div className="flex items-start gap-3">
									<div className="p-2 bg-primary/10 rounded-lg">
										<Gauge className="h-5 w-5 text-primary" />
									</div>
									<div>
										<p className="text-sm text-muted-foreground">
											RAM
										</p>
										<p className="text-lg font-semibold">
											{formatStorage(product.ram_mb)}
										</p>
									</div>
								</div>

								<div className="flex items-start gap-3">
									<div className="p-2 bg-primary/10 rounded-lg">
										<HardDrive className="h-5 w-5 text-primary" />
									</div>
									<div>
										<p className="text-sm text-muted-foreground">
											Storage
										</p>
										<p className="text-lg font-semibold">
											{product.storage_gb} GB
										</p>
									</div>
								</div>

								<div className="flex items-start gap-3">
									<div className="p-2 bg-primary/10 rounded-lg">
										<Network className="h-5 w-5 text-primary" />
									</div>
									<div>
										<p className="text-sm text-muted-foreground">
											Bandwidth
										</p>
										<p className="text-lg font-semibold">
											{product.bandwidth_gb} GB
										</p>
									</div>
								</div>
							</div>

							<div className="pt-4 border-t">
								<p className="text-sm text-muted-foreground mb-1">
									Availability
								</p>
								<p className="text-lg font-semibold">
									{product.stock === -1
										? "Unlimited Stock"
										: `${product.stock} units available`}
								</p>
							</div>
						</div>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle className="text-2xl">
							Order Summary
						</CardTitle>
					</CardHeader>

					<CardContent>
						<div className="space-y-6">
							{/* Version Selector */}
							{product.image_tags && product.image_tags.length > 0 && (
								<div>
									<label
										htmlFor="version"
										className="text-sm font-medium block mb-2"
									>
										Select Docker Image Version
									</label>
									<select
										id="version"
										value={selectedVersion}
										onChange={(e) =>
											setSelectedVersion(e.target.value)
										}
										className="w-full border rounded-md px-3 py-2 bg-background"
										disabled={checkoutLoading}
									>
										{product.image_tags.map((tag) => (
											<option key={tag} value={tag}>
												{tag}
											</option>
										))}
									</select>
									<p className="text-sm text-muted-foreground mt-1">
										Full image: {product.docker_image}:
										{selectedVersion}
									</p>
								</div>
							)}

							<div>
								<label
									htmlFor="quantity"
									className="text-sm font-medium block mb-2"
								>
									Quantity
								</label>
								<div className="flex items-center gap-2">
									<Button
										variant="outline"
										size="sm"
										onClick={() =>
											setQuantity((q) =>
												Math.max(1, q - 1),
											)
										}
										disabled={checkoutLoading}
									>
										-
									</Button>
									<input
										id="quantity"
										type="number"
										min="1"
										max={
											product.stock === -1
												? 100
												: product.stock
										}
										value={quantity}
										onChange={(e) =>
											setQuantity(
												Math.max(
													1,
													parseInt(e.target.value) ||
														1,
												),
											)
										}
										className="w-20 text-center border rounded-md px-3 py-2"
										disabled={checkoutLoading}
									/>
									<Button
										variant="outline"
										size="sm"
										onClick={() =>
											setQuantity((q) => q + 1)
										}
										disabled={
											checkoutLoading ||
											(product.stock !== -1 &&
												quantity >= product.stock)
										}
									>
										+
									</Button>
								</div>
								{isOutOfStock && (
									<p className="text-sm text-destructive mt-2">
										Not enough stock available
									</p>
								)}
							</div>

							<div className="space-y-2 pt-4 border-t">
								<div className="flex justify-between">
									<span className="text-muted-foreground">
										Price per unit
									</span>
									<span className="font-medium">
										{formatPrice(product.price)}
									</span>
								</div>
								<div className="flex justify-between">
									<span className="text-muted-foreground">
										Quantity
									</span>
									<span className="font-medium">
										{quantity}
									</span>
								</div>
								<div className="flex justify-between text-lg font-bold pt-2 border-t">
									<span>Total</span>
									<span>{formatPrice(totalPrice)}</span>
								</div>
								<p className="text-sm text-muted-foreground">
									per month
								</p>
							</div>
						</div>
					</CardContent>

					<CardFooter>
						<Button
							className="w-full"
							size="lg"
							onClick={handleCheckout}
							disabled={
								!product.is_active ||
								isOutOfStock ||
								checkoutLoading ||
								quantity < 1
							}
						>
							{checkoutLoading ? (
								<>
									<Loader2 className="mr-2 h-4 w-4 animate-spin" />
									Processing...
								</>
							) : isAuthenticated ? (
								"Proceed to Payment"
							) : (
								"Login to Purchase"
							)}
						</Button>
					</CardFooter>
				</Card>
			</div>
		</div>
	);
}
