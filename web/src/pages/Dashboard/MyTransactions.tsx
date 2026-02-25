import { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";
import {
	getMyTransactions,
	cancelTransaction,
	type Transaction,
	type TransactionStatus,
	getStatusColor,
	getStatusText,
} from "@/lib/api/transactions";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Loader2, Package, Calendar, CreditCard } from "lucide-react";

export function MyTransactions() {
	const location = useLocation();
	const [transactions, setTransactions] = useState<Transaction[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [page, setPage] = useState(1);
	const [total, setTotal] = useState(0);
	const [successMessage, setSuccessMessage] = useState<string | null>(null);
	const pageSize = 10;

	useEffect(() => {
		// Check for success message from location state
		if (location.state?.message) {
			setSuccessMessage(location.state.message);
			// Clear the message after 5 seconds
			setTimeout(() => setSuccessMessage(null), 5000);
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
			// Refresh the list
			await loadTransactions();
		} catch (err: any) {
			console.error("Error cancelling transaction:", err);
			alert(
				err.response?.data?.message ||
					"Failed to cancel transaction. Please try again.",
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

	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString("en-US", {
			year: "numeric",
			month: "long",
			day: "numeric",
			hour: "2-digit",
			minute: "2-digit",
		});
	};

	const canCancel = (transaction: Transaction) => {
		return transaction.status === "pending" || transaction.status === "processing";
	};

	const totalPages = Math.ceil(total / pageSize);

	if (loading && transactions.length === 0) {
		return (
			<div className="flex items-center justify-center min-h-[400px]">
				<Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
			</div>
		);
	}

	return (
		<div className="space-y-6">
			<div>
				<h1 className="text-3xl font-bold">My Transactions</h1>
				<p className="text-muted-foreground">
					View and manage your purchase history
				</p>
			</div>

			{successMessage && (
				<div className="bg-green-100 dark:bg-green-900 border border-green-200 dark:border-green-800 text-green-800 dark:text-green-200 px-4 py-3 rounded-lg">
					{successMessage}
				</div>
			)}

			{error && (
				<div className="bg-destructive/10 border border-destructive/20 text-destructive px-4 py-3 rounded-lg">
					{error}
				</div>
			)}

			{transactions.length === 0 && !loading ? (
				<Card>
					<CardContent className="py-12 text-center">
						<Package className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
						<p className="text-muted-foreground text-lg">
							No transactions found.
						</p>
						<Button className="mt-4" onClick={() => (window.location.href = "/products")}>
							Browse Products
						</Button>
					</CardContent>
				</Card>
			) : (
				<>
					<div className="space-y-4">
						{transactions.map((transaction) => (
							<Card key={transaction.id}>
								<CardHeader>
									<div className="flex items-start justify-between">
										<div>
											<CardTitle className="text-lg">
												{transaction.product?.name || "Product"}
											</CardTitle>
											<CardDescription>
												Order #{transaction.order_id}
											</CardDescription>
										</div>
										<Badge className={getStatusColor(transaction.status)}>
											{getStatusText(transaction.status)}
										</Badge>
									</div>
								</CardHeader>

								<CardContent>
									<div className="grid md:grid-cols-2 gap-6">
										<div className="space-y-3">
											<div className="flex items-center gap-2 text-sm">
												<Calendar className="h-4 w-4 text-muted-foreground" />
												<span className="text-muted-foreground">Ordered:</span>
												<span>{formatDate(transaction.created_at)}</span>
											</div>

											<div className="flex items-center gap-2 text-sm">
												<Package className="h-4 w-4 text-muted-foreground" />
												<span className="text-muted-foreground">Quantity:</span>
												<span>{transaction.quantity}</span>
											</div>

											{transaction.payment_method && (
												<div className="flex items-center gap-2 text-sm">
													<CreditCard className="h-4 w-4 text-muted-foreground" />
													<span className="text-muted-foreground">Payment:</span>
													<span>{transaction.payment_method}</span>
												</div>
											)}
										</div>

										<div className="space-y-3">
											<div className="flex justify-between items-center">
												<span className="text-muted-foreground">Amount</span>
												<span className="text-xl font-bold">
													{formatPrice(transaction.amount)}
												</span>
											</div>

											{transaction.expired_at && transaction.status === "pending" && (
												<div className="text-sm">
													<span className="text-muted-foreground">Expires:</span>
													<span className="ml-2">
														{formatDate(transaction.expired_at)}
													</span>
												</div>
											)}

											{transaction.paid_at && (
												<div className="text-sm">
													<span className="text-muted-foreground">Paid:</span>
													<span className="ml-2">
														{formatDate(transaction.paid_at)}
													</span>
												</div>
											)}
										</div>
									</div>

									{canCancel(transaction) && (
										<div className="mt-4 pt-4 border-t flex gap-2">
											{transaction.payment_url && (
												<Button
													variant="default"
													size="sm"
													onClick={() =>
														window.open(transaction.payment_url, "_blank")
													}
												>
													Continue Payment
												</Button>
											)}
											<Button
												variant="outline"
												size="sm"
												onClick={() =>
													handleCancelTransaction(transaction.order_id)
												}
											>
												Cancel Order
											</Button>
										</div>
									)}
								</CardContent>
							</Card>
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
		</div>
	);
}
