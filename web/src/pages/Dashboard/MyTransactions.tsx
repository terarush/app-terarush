import { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";
import {
	getMyTransactions,
	cancelTransaction,
	type Transaction,
	getStatusColor,
	getStatusText,
} from "@/lib/api/transactions";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Loader2, XCircle } from "lucide-react";
import { toast } from "sonner";

export function MyTransactions() {
	const location = useLocation();
	const [transactions, setTransactions] = useState<Transaction[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [page, setPage] = useState(1);
	const [total, setTotal] = useState(0);
	const pageSize = 10;

	useEffect(() => {
		// Check for success message from location state
		if (location.state?.message) {
			toast.success(location.state.message);
		}
	}, [location]);

	useEffect(() => {
		loadTransactions();
	}, [page]);

	const loadTransactions = async () => {
		try {
			setLoading(true);
			setError(null);
			const response = await getMyTransactions({
				page,
				page_size: pageSize,
			});
			setTransactions(response.transactions || []);
			setTotal(response.total || 0);
		} catch (err: any) {
			console.error("Error loading transactions:", err);
			setError(
				err.response?.data?.message ||
					"Failed to load transactions. Please try again.",
			);
		} finally {
			setLoading(false);
		}
	};

	const handleCancelTransaction = async (orderId: string) => {
		if (!confirm("Are you sure you want to cancel this transaction?")) {
			return;
		}

		try {
			await cancelTransaction(orderId);
			toast.success("Transaction cancelled successfully");
			await loadTransactions();
		} catch (err: any) {
			console.error("Error cancelling transaction:", err);
			toast.error(
				err.response?.data?.message || "Failed to cancel transaction",
			);
		}
	};

	const formatPrice = (price: number) => {
		return new Intl.NumberFormat("id-ID", {
			style: "currency",
			currency: "IDR",
			minimumFractionDigits: 0,
		}).format(price);
	};

	const formatDate = (date: string) => {
		return new Date(date).toLocaleDateString("id-ID", {
			year: "numeric",
			month: "short",
			day: "numeric",
			hour: "2-digit",
			minute: "2-digit",
		});
	};

	const totalPages = Math.ceil(total / pageSize);

	if (loading && transactions.length === 0) {
		return (
			<div className="flex items-center justify-center min-h-[400px]">
				<Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
			</div>
		);
	}

	if (error) {
		return (
			<div className="text-center py-12">
				<p className="text-destructive mb-4">{error}</p>
				<Button onClick={loadTransactions}>Try Again</Button>
			</div>
		);
	}

	return (
		<>
			<div className="mb-6">
				<h1 className="text-3xl font-bold">My Transactions</h1>
				<p className="text-muted-foreground mt-1">
					View and manage your purchase history
				</p>
			</div>

			{transactions.length === 0 ? (
				<div className="text-center py-12 border rounded-lg">
					<p className="text-muted-foreground text-lg">
						No transactions found. Start by purchasing a product!
					</p>
				</div>
			) : (
				<>
					<div className="border rounded-lg">
						<Table>
							<TableHeader>
								<TableRow>
									<TableHead>Order ID</TableHead>
									<TableHead>Amount</TableHead>
									<TableHead>Quantity</TableHead>
									<TableHead>Status</TableHead>
									<TableHead>Payment Method</TableHead>
									<TableHead>Date</TableHead>
									<TableHead className="text-right">Actions</TableHead>
								</TableRow>
							</TableHeader>
							<TableBody>
								{transactions.map((transaction) => (
									<TableRow key={transaction.id}>
										<TableCell className="font-medium">
											{transaction.order_id}
										</TableCell>
										<TableCell>{formatPrice(transaction.amount)}</TableCell>
										<TableCell>{transaction.quantity}</TableCell>
										<TableCell>
											<Badge variant={getStatusColor(transaction.status)}>
												{getStatusText(transaction.status)}
											</Badge>
										</TableCell>
										<TableCell>
											{transaction.payment_method || "-"}
										</TableCell>
										<TableCell>
											{formatDate(transaction.created_at)}
										</TableCell>
										<TableCell className="text-right">
											{transaction.status === "pending" && (
												<Button
													variant="ghost"
													size="sm"
													onClick={() =>
														handleCancelTransaction(transaction.order_id)
													}
												>
													<XCircle className="h-4 w-4 mr-2" />
													Cancel
												</Button>
											)}
										</TableCell>
									</TableRow>
								))}
							</TableBody>
						</Table>
					</div>

					{totalPages > 1 && (
						<div className="flex items-center justify-center gap-2 mt-6">
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
		</>
	);
}
