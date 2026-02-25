import { useState, useEffect } from "react";
import {
	getProducts,
	createProduct,
	updateProduct,
	deleteProduct,
	type Product,
	type CreateProductRequest,
} from "@/lib/api/products";
import { ProductCard } from "@/components/fragments/ProductCard";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Loader2, Plus } from "lucide-react";
import { toast } from "sonner";

export function AdminProducts() {
	const [products, setProducts] = useState<Product[]>([]);
	const [loading, setLoading] = useState(true);
	const [dialogOpen, setDialogOpen] = useState(false);
	const [editingProduct, setEditingProduct] = useState<Product | null>(null);
	const [formLoading, setFormLoading] = useState(false);
	const [page, setPage] = useState(1);
	const [total, setTotal] = useState(0);
	const pageSize = 12;

	const [formData, setFormData] = useState<CreateProductRequest>({
		name: "",
		description: "",
		price: 0,
		cpu_cores: 1,
		ram_mb: 512,
		storage_gb: 10,
		bandwidth_gb: 100,
		is_active: true,
		stock: -1,
		image_url: "",
	});

	useEffect(() => {
		loadProducts();
	}, [page]);

	const loadProducts = async () => {
		try {
			setLoading(true);
			const response = await getProducts({
				page,
				page_size: pageSize,
			});
			setProducts(response.products || []);
			setTotal(response.total || 0);
		} catch (err: any) {
			console.error("Error loading products:", err);
			toast.error("Failed to load products");
		} finally {
			setLoading(false);
		}
	};

	const handleOpenDialog = (product?: Product) => {
		if (product) {
			setEditingProduct(product);
			setFormData({
				name: product.name,
				description: product.description,
				price: product.price,
				cpu_cores: product.cpu_cores,
				ram_mb: product.ram_mb,
				storage_gb: product.storage_gb,
				bandwidth_gb: product.bandwidth_gb,
				is_active: product.is_active,
				stock: product.stock,
				image_url: product.image_url,
			});
		} else {
			setEditingProduct(null);
			setFormData({
				name: "",
				description: "",
				price: 0,
				cpu_cores: 1,
				ram_mb: 512,
				storage_gb: 10,
				bandwidth_gb: 100,
				is_active: true,
				stock: -1,
				image_url: "",
			});
		}
		setDialogOpen(true);
	};

	const handleCloseDialog = () => {
		setDialogOpen(false);
		setEditingProduct(null);
	};

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setFormLoading(true);

		try {
			if (editingProduct) {
				await updateProduct(editingProduct.id, formData);
				toast.success("Product updated successfully");
			} else {
				await createProduct(formData);
				toast.success("Product created successfully");
			}
			handleCloseDialog();
			await loadProducts();
		} catch (err: any) {
			console.error("Error saving product:", err);
			toast.error(
				err.response?.data?.message || "Failed to save product",
			);
		} finally {
			setFormLoading(false);
		}
	};

	const handleDelete = async (product: Product) => {
		if (
			!confirm(
				`Are you sure you want to delete "${product.name}"? This action cannot be undone.`,
			)
		) {
			return;
		}

		try {
			await deleteProduct(product.id);
			toast.success("Product deleted successfully");
			await loadProducts();
		} catch (err: any) {
			console.error("Error deleting product:", err);
			toast.error(
				err.response?.data?.message || "Failed to delete product",
			);
		}
	};

	const totalPages = Math.ceil(total / pageSize);

	if (loading && products.length === 0) {
		return (
			<div className="flex items-center justify-center min-h-[400px]">
				<Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
			</div>
		);
	}

	return (
		<div className="space-y-6">
			<div className="flex items-center justify-between">
				<div>
					<h1 className="text-3xl font-bold">Products Management</h1>
					<p className="text-muted-foreground">
						Manage container products for Node.js, Python, and Ubuntu
					</p>
				</div>
				<Button onClick={() => handleOpenDialog()}>
					<Plus className="mr-2 h-4 w-4" />
					Add Product
				</Button>
			</div>

			{products.length === 0 && !loading ? (
				<div className="text-center py-12 border rounded-lg">
					<p className="text-muted-foreground text-lg mb-4">
						No products found. Create your first product!
					</p>
					<Button onClick={() => handleOpenDialog()}>
						<Plus className="mr-2 h-4 w-4" />
						Add Product
					</Button>
				</div>
			) : (
				<>
					<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
						{products.map((product) => (
							<ProductCard
								key={product.id}
								product={product}
								showAdminActions={true}
								onEdit={handleOpenDialog}
								onDelete={handleDelete}
							/>
						))}
					</div>

					{totalPages > 1 && (
						<div className="flex items-center justify-center gap-2">
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

			<Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
				<DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
					<DialogHeader>
						<DialogTitle>
							{editingProduct ? "Edit Product" : "Create New Product"}
						</DialogTitle>
						<DialogDescription>
							Configure container specifications for bot hosting (Node.js,
							Python, Ubuntu)
						</DialogDescription>
					</DialogHeader>

					<form onSubmit={handleSubmit} className="space-y-4">
						<div className="grid grid-cols-2 gap-4">
							<div className="col-span-2">
								<Label htmlFor="name">Product Name *</Label>
								<Input
									id="name"
									placeholder="e.g., Node.js Bot Starter"
									value={formData.name}
									onChange={(e) =>
										setFormData({ ...formData, name: e.target.value })
									}
									required
								/>
							</div>

							<div className="col-span-2">
								<Label htmlFor="description">Description *</Label>
								<Textarea
									id="description"
									placeholder="Perfect for running Discord, Telegram, or WhatsApp bots..."
									value={formData.description}
									onChange={(e) =>
										setFormData({ ...formData, description: e.target.value })
									}
									rows={3}
									required
								/>
							</div>

							<div>
								<Label htmlFor="price">Price (IDR/month) *</Label>
								<Input
									id="price"
									type="number"
									min="0"
									step="1000"
									placeholder="50000"
									value={formData.price}
									onChange={(e) =>
										setFormData({
											...formData,
											price: Number(e.target.value),
										})
									}
									required
								/>
							</div>

							<div>
								<Label htmlFor="cpu_cores">CPU Cores *</Label>
								<Input
									id="cpu_cores"
									type="number"
									min="1"
									max="32"
									placeholder="1"
									value={formData.cpu_cores}
									onChange={(e) =>
										setFormData({
											...formData,
											cpu_cores: Number(e.target.value),
										})
									}
									required
								/>
							</div>

							<div>
								<Label htmlFor="ram_mb">RAM (MB) *</Label>
								<Input
									id="ram_mb"
									type="number"
									min="256"
									step="256"
									placeholder="512"
									value={formData.ram_mb}
									onChange={(e) =>
										setFormData({
											...formData,
											ram_mb: Number(e.target.value),
										})
									}
									required
								/>
							</div>

							<div>
								<Label htmlFor="storage_gb">Storage (GB) *</Label>
								<Input
									id="storage_gb"
									type="number"
									min="5"
									placeholder="10"
									value={formData.storage_gb}
									onChange={(e) =>
										setFormData({
											...formData,
											storage_gb: Number(e.target.value),
										})
									}
									required
								/>
							</div>

							<div>
								<Label htmlFor="bandwidth_gb">Bandwidth (GB/month) *</Label>
								<Input
									id="bandwidth_gb"
									type="number"
									min="10"
									placeholder="100"
									value={formData.bandwidth_gb}
									onChange={(e) =>
										setFormData({
											...formData,
											bandwidth_gb: Number(e.target.value),
										})
									}
									required
								/>
							</div>

							<div>
								<Label htmlFor="stock">
									Stock (-1 for unlimited)
								</Label>
								<Input
									id="stock"
									type="number"
									min="-1"
									placeholder="-1"
									value={formData.stock}
									onChange={(e) =>
										setFormData({
											...formData,
											stock: Number(e.target.value),
										})
									}
								/>
							</div>

							<div className="col-span-2">
								<Label htmlFor="image_url">Image URL (optional)</Label>
								<Input
									id="image_url"
									type="url"
									placeholder="https://example.com/image.png"
									value={formData.image_url}
									onChange={(e) =>
										setFormData({ ...formData, image_url: e.target.value })
									}
								/>
							</div>

							<div className="col-span-2 flex items-center space-x-2">
								<Switch
									id="is_active"
									checked={formData.is_active}
									onCheckedChange={(checked) =>
										setFormData({ ...formData, is_active: checked })
									}
								/>
								<Label htmlFor="is_active" className="cursor-pointer">
									Active (visible to users)
								</Label>
							</div>
						</div>

						<DialogFooter>
							<Button
								type="button"
								variant="outline"
								onClick={handleCloseDialog}
								disabled={formLoading}
							>
								Cancel
							</Button>
							<Button type="submit" disabled={formLoading}>
								{formLoading ? (
									<>
										<Loader2 className="mr-2 h-4 w-4 animate-spin" />
										Saving...
									</>
								) : editingProduct ? (
									"Update Product"
								) : (
									"Create Product"
								)}
							</Button>
						</DialogFooter>
					</form>
				</DialogContent>
			</Dialog>
		</div>
	);
}
