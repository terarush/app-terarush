import { useState, useEffect } from "react";
import {
	getProducts,
	createProduct,
	updateProduct,
	deleteProduct,
	type Product,
	type CreateProductRequest,
} from "@/lib/api/products";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { Loader2, Plus, Pencil, Trash2 } from "lucide-react";
import { toast } from "sonner";

export function Products() {
	const [products, setProducts] = useState<Product[]>([]);
	const [loading, setLoading] = useState(true);
	const [dialogOpen, setDialogOpen] = useState(false);
	const [editingProduct, setEditingProduct] = useState<Product | null>(null);
	const [formLoading, setFormLoading] = useState(false);

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
		docker_image: "",
		image_tags: [],
	});

	useEffect(() => {
		loadProducts();
	}, []);

	const loadProducts = async () => {
		try {
			setLoading(true);
			const response = await getProducts();
			setProducts(response.products || []);
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
				image_url: product.image_url || "",
				docker_image: product.docker_image || "",
				image_tags: [],
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
				docker_image: "",
				image_tags: [],
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
			loadProducts();
		} catch (err: any) {
			console.error("Error saving product:", err);
			toast.error(
				err.response?.data?.error?.message || "Failed to save product",
			);
		} finally {
			setFormLoading(false);
		}
	};

	const handleDelete = async (id: number) => {
		if (!confirm("Are you sure you want to delete this product?")) return;

		try {
			await deleteProduct(id);
			toast.success("Product deleted successfully");
			loadProducts();
		} catch (err: any) {
			console.error("Error deleting product:", err);
			toast.error("Failed to delete product");
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
			<div className="flex items-center justify-center min-h-[400px]">
				<Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
			</div>
		);
	}

	return (
		<>
			<div className="flex items-center justify-between mb-6">
				<div>
					<h1 className="text-3xl font-bold">Products</h1>
					<p className="text-muted-foreground mt-1">
						Manage container products for Node.js, Python, and Ubuntu
					</p>
				</div>
				<Button onClick={() => handleOpenDialog()}>
					<Plus className="mr-2 h-4 w-4" />
					Add Product
				</Button>
			</div>

			{products.length === 0 ? (
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
				<div className="border rounded-lg">
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead>Name</TableHead>
								<TableHead>Description</TableHead>
								<TableHead>Price</TableHead>
								<TableHead>CPU</TableHead>
								<TableHead>RAM</TableHead>
								<TableHead>Storage</TableHead>
								<TableHead>Bandwidth</TableHead>
								<TableHead>Stock</TableHead>
								<TableHead>Status</TableHead>
								<TableHead className="text-right">Actions</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							{products.map((product) => (
								<TableRow key={product.id}>
									<TableCell className="font-medium">
										{product.name}
									</TableCell>
									<TableCell className="max-w-xs truncate">
										{product.description}
									</TableCell>
									<TableCell>{formatPrice(product.price)}</TableCell>
									<TableCell>{product.cpu_cores} Cores</TableCell>
									<TableCell>{formatStorage(product.ram_mb)}</TableCell>
									<TableCell>{product.storage_gb} GB</TableCell>
									<TableCell>{product.bandwidth_gb} GB</TableCell>
									<TableCell>
										{product.stock === -1 ? "Unlimited" : product.stock}
									</TableCell>
									<TableCell>
										<Badge variant={product.is_active ? "default" : "secondary"}>
											{product.is_active ? "Active" : "Inactive"}
										</Badge>
									</TableCell>
									<TableCell className="text-right">
										<div className="flex items-center justify-end gap-2">
											<Button
												variant="ghost"
												size="icon"
												onClick={() => handleOpenDialog(product)}
											>
												<Pencil className="h-4 w-4" />
											</Button>
											<Button
												variant="ghost"
												size="icon"
												onClick={() => handleDelete(product.id)}
											>
												<Trash2 className="h-4 w-4 text-destructive" />
											</Button>
										</div>
									</TableCell>
								</TableRow>
							))}
						</TableBody>
					</Table>
				</div>
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
										setFormData({
											...formData,
											name: e.target.value,
										})
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
										setFormData({
											...formData,
											description: e.target.value,
										})
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
									min="1"
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
								<Label htmlFor="bandwidth_gb">Bandwidth (GB) *</Label>
								<Input
									id="bandwidth_gb"
									type="number"
									min="1"
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
								<Label htmlFor="image_url">Image URL</Label>
								<Input
									id="image_url"
									type="url"
									placeholder="https://example.com/image.png"
									value={formData.image_url}
									onChange={(e) =>
										setFormData({
											...formData,
											image_url: e.target.value,
										})
									}
								/>
							</div>

							<div className="col-span-2 flex items-center space-x-2">
								<Switch
									id="is_active"
									checked={formData.is_active}
									onCheckedChange={(checked) =>
										setFormData({
											...formData,
											is_active: checked,
										})
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
		</>
	);
}
