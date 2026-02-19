import { useState, useEffect } from "react";
import { toast } from "sonner";
import {
	getNodes,
	createNode,
	updateNode,
	deleteNode,
	nodeAction,
	type Node,
	type CreateNodeRequest,
	type UpdateNodeRequest,
} from "@/lib/api/nodes";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Server, Play, Square, RotateCw, Trash2, Plus, Edit } from "lucide-react";

export default function Nodes() {
	const [nodes, setNodes] = useState<Node[]>([]);
	const [loading, setLoading] = useState(true);
	const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
	const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
	const [selectedNode, setSelectedNode] = useState<Node | null>(null);
	const [formData, setFormData] = useState<CreateNodeRequest>({
		name: "",
		image: "",
		port: 8080,
		internal_port: 8080,
		cpu_limit: 1.0,
		memory_limit: 512,
		restart_policy: "unless-stopped",
	});

	useEffect(() => {
		loadNodes();
	}, []);

	const loadNodes = async () => {
		try {
			setLoading(true);
			const data = await getNodes();
			setNodes(data || []); // Ensure it's always an array
		} catch (error) {
			console.error("Failed to load nodes:", error);
			toast.error("Failed to load nodes");
			setNodes([]); // Set empty array on error
		} finally {
			setLoading(false);
		}
	};

	const handleCreateNode = async () => {
		try {
			await createNode(formData);
			toast.success("Node created successfully");
			setIsCreateDialogOpen(false);
			resetForm();
			loadNodes();
		} catch (error: any) {
			console.error("Failed to create node:", error);
			toast.error(error.response?.data?.error || "Failed to create node");
		}
	};

	const handleUpdateNode = async () => {
		if (!selectedNode) return;

		try {
			const updateData: UpdateNodeRequest = {
				name: formData.name,
				cpu_limit: formData.cpu_limit,
				memory_limit: formData.memory_limit,
				restart_policy: formData.restart_policy,
			};
			await updateNode(selectedNode.id, updateData);
			toast.success("Node updated successfully");
			setIsEditDialogOpen(false);
			setSelectedNode(null);
			resetForm();
			loadNodes();
		} catch (error: any) {
			console.error("Failed to update node:", error);
			toast.error(error.response?.data?.error || "Failed to update node");
		}
	};

	const handleDeleteNode = async (id: number) => {
		if (!confirm("Are you sure you want to delete this node? This will stop and remove the container.")) {
			return;
		}

		try {
			await deleteNode(id);
			toast.success("Node deleted successfully");
			loadNodes();
		} catch (error: any) {
			console.error("Failed to delete node:", error);
			toast.error(error.response?.data?.error || "Failed to delete node");
		}
	};

	const handleNodeAction = async (id: number, action: "start" | "stop" | "restart") => {
		try {
			await nodeAction(id, action);
			toast.success(`Node ${action}ed successfully`);
			loadNodes();
		} catch (error: any) {
			console.error(`Failed to ${action} node:`, error);
			toast.error(error.response?.data?.error || `Failed to ${action} node`);
		}
	};

	const openEditDialog = (node: Node) => {
		setSelectedNode(node);
		setFormData({
			name: node.name,
			image: node.image,
			port: node.port,
			internal_port: node.internal_port,
			cpu_limit: node.cpu_limit,
			memory_limit: node.memory_limit,
			restart_policy: node.restart_policy,
		});
		setIsEditDialogOpen(true);
	};

	const resetForm = () => {
		setFormData({
			name: "",
			image: "",
			port: 8080,
			internal_port: 8080,
			cpu_limit: 1.0,
			memory_limit: 512,
			restart_policy: "unless-stopped",
		});
	};

	const getStatusColor = (status: string) => {
		switch (status) {
			case "running":
				return "bg-green-500";
			case "stopped":
				return "bg-red-500";
			case "created":
				return "bg-blue-500";
			case "error":
				return "bg-yellow-500";
			default:
				return "bg-gray-500";
		}
	};

	return (
		<div className="p-6 space-y-6">
			<div className="flex justify-between items-center">
				<div>
					<h1 className="text-3xl font-bold">Nodes Management</h1>
					<p className="text-muted-foreground">Manage your Docker containers</p>
				</div>
				<Button onClick={() => setIsCreateDialogOpen(true)}>
					<Plus className="mr-2 h-4 w-4" />
					Create Node
				</Button>
			</div>

			{loading ? (
				<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{[1, 2, 3].map((i) => (
						<Card key={i} className="animate-pulse">
							<CardHeader>
								<div className="h-6 bg-muted rounded w-3/4"></div>
								<div className="h-4 bg-muted rounded w-1/2 mt-2"></div>
							</CardHeader>
							<CardContent>
								<div className="h-20 bg-muted rounded"></div>
							</CardContent>
						</Card>
					))}
				</div>
			) : !nodes || nodes.length === 0 ? (
				<Card>
					<CardContent className="flex flex-col items-center justify-center py-12">
						<Server className="h-12 w-12 text-muted-foreground mb-4" />
						<p className="text-lg font-medium">No nodes yet</p>
						<p className="text-sm text-muted-foreground mb-4">
							Create your first node to get started
						</p>
						<Button onClick={() => setIsCreateDialogOpen(true)}>
							<Plus className="mr-2 h-4 w-4" />
							Create Node
						</Button>
					</CardContent>
				</Card>
			) : (
				<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{nodes.map((node) => (
						<Card key={node.id}>
							<CardHeader>
								<div className="flex items-start justify-between">
									<div className="flex-1">
										<CardTitle className="flex items-center gap-2">
											<Server className="h-5 w-5" />
											{node.name}
										</CardTitle>
										<CardDescription className="mt-1">
											{node.image}
										</CardDescription>
									</div>
									<Badge className={getStatusColor(node.status)}>
										{node.status}
									</Badge>
								</div>
							</CardHeader>
							<CardContent className="space-y-4">
								<div className="space-y-2 text-sm">
									<div className="flex justify-between">
										<span className="text-muted-foreground">Port:</span>
										<span className="font-medium">
											{node.port} → {node.internal_port}
										</span>
									</div>
									<div className="flex justify-between">
										<span className="text-muted-foreground">CPU:</span>
										<span className="font-medium">{node.cpu_limit} cores</span>
									</div>
									<div className="flex justify-between">
										<span className="text-muted-foreground">Memory:</span>
										<span className="font-medium">{node.memory_limit} MB</span>
									</div>
									<div className="flex justify-between">
										<span className="text-muted-foreground">Restart:</span>
										<span className="font-medium">{node.restart_policy}</span>
									</div>
								</div>

								<div className="flex gap-2">
									{node.status === "stopped" || node.status === "created" ? (
										<Button
											size="sm"
											variant="default"
											className="flex-1"
											onClick={() => handleNodeAction(node.id, "start")}
										>
											<Play className="mr-1 h-3 w-3" />
											Start
										</Button>
									) : (
										<Button
											size="sm"
											variant="secondary"
											className="flex-1"
											onClick={() => handleNodeAction(node.id, "stop")}
										>
											<Square className="mr-1 h-3 w-3" />
											Stop
										</Button>
									)}
									<Button
										size="sm"
										variant="outline"
										onClick={() => handleNodeAction(node.id, "restart")}
										disabled={node.status === "stopped" || node.status === "created"}
									>
										<RotateCw className="h-3 w-3" />
									</Button>
									<Button
										size="sm"
										variant="outline"
										onClick={() => openEditDialog(node)}
									>
										<Edit className="h-3 w-3" />
									</Button>
									<Button
										size="sm"
										variant="destructive"
										onClick={() => handleDeleteNode(node.id)}
									>
										<Trash2 className="h-3 w-3" />
									</Button>
								</div>
							</CardContent>
						</Card>
					))}
				</div>
			)}

			{/* Create Node Dialog */}
			<Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
				<DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
					<DialogHeader>
						<DialogTitle>Create New Node</DialogTitle>
						<DialogDescription>
							Create a new Docker container node
						</DialogDescription>
					</DialogHeader>
					<div className="grid gap-4 py-4">
						<div className="grid gap-2">
							<label htmlFor="name" className="text-sm font-medium">
								Node Name *
							</label>
							<Input
								id="name"
								placeholder="my-node"
								value={formData.name}
								onChange={(e) =>
									setFormData({ ...formData, name: e.target.value })
								}
							/>
						</div>
						<div className="grid gap-2">
							<label htmlFor="image" className="text-sm font-medium">
								Docker Image *
							</label>
							<Input
								id="image"
								placeholder="nginx:latest"
								value={formData.image}
								onChange={(e) =>
									setFormData({ ...formData, image: e.target.value })
								}
							/>
						</div>
						<div className="grid grid-cols-2 gap-4">
							<div className="grid gap-2">
								<label htmlFor="port" className="text-sm font-medium">
									Host Port *
								</label>
								<Input
									id="port"
									type="number"
									placeholder="8080"
									value={formData.port}
									onChange={(e) =>
										setFormData({
											...formData,
											port: parseInt(e.target.value),
										})
									}
								/>
							</div>
							<div className="grid gap-2">
								<label htmlFor="internal_port" className="text-sm font-medium">
									Container Port *
								</label>
								<Input
									id="internal_port"
									type="number"
									placeholder="80"
									value={formData.internal_port}
									onChange={(e) =>
										setFormData({
											...formData,
											internal_port: parseInt(e.target.value),
										})
									}
								/>
							</div>
						</div>
						<div className="grid grid-cols-2 gap-4">
							<div className="grid gap-2">
								<label htmlFor="cpu_limit" className="text-sm font-medium">
									CPU Limit (cores)
								</label>
								<Input
									id="cpu_limit"
									type="number"
									step="0.1"
									placeholder="1.0"
									value={formData.cpu_limit}
									onChange={(e) =>
										setFormData({
											...formData,
											cpu_limit: parseFloat(e.target.value),
										})
									}
								/>
							</div>
							<div className="grid gap-2">
								<label htmlFor="memory_limit" className="text-sm font-medium">
									Memory Limit (MB)
								</label>
								<Input
									id="memory_limit"
									type="number"
									placeholder="512"
									value={formData.memory_limit}
									onChange={(e) =>
										setFormData({
											...formData,
											memory_limit: parseInt(e.target.value),
										})
									}
								/>
							</div>
						</div>
						<div className="grid gap-2">
							<label htmlFor="restart_policy" className="text-sm font-medium">
								Restart Policy
							</label>
							<select
								id="restart_policy"
								className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
								value={formData.restart_policy}
								onChange={(e) =>
									setFormData({ ...formData, restart_policy: e.target.value })
								}
							>
								<option value="no">No</option>
								<option value="on-failure">On Failure</option>
								<option value="always">Always</option>
								<option value="unless-stopped">Unless Stopped</option>
							</select>
						</div>
					</div>
					<DialogFooter>
						<Button
							variant="outline"
							onClick={() => {
								setIsCreateDialogOpen(false);
								resetForm();
							}}
						>
							Cancel
						</Button>
						<Button onClick={handleCreateNode}>Create Node</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>

			{/* Edit Node Dialog */}
			<Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
				<DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
					<DialogHeader>
						<DialogTitle>Edit Node</DialogTitle>
						<DialogDescription>
							Update node configuration
						</DialogDescription>
					</DialogHeader>
					<div className="grid gap-4 py-4">
						<div className="grid gap-2">
							<label htmlFor="edit_name" className="text-sm font-medium">
								Node Name
							</label>
							<Input
								id="edit_name"
								value={formData.name}
								onChange={(e) =>
									setFormData({ ...formData, name: e.target.value })
								}
							/>
						</div>
						<div className="grid grid-cols-2 gap-4">
							<div className="grid gap-2">
								<label htmlFor="edit_cpu_limit" className="text-sm font-medium">
									CPU Limit (cores)
								</label>
								<Input
									id="edit_cpu_limit"
									type="number"
									step="0.1"
									value={formData.cpu_limit}
									onChange={(e) =>
										setFormData({
											...formData,
											cpu_limit: parseFloat(e.target.value),
										})
									}
								/>
							</div>
							<div className="grid gap-2">
								<label htmlFor="edit_memory_limit" className="text-sm font-medium">
									Memory Limit (MB)
								</label>
								<Input
									id="edit_memory_limit"
									type="number"
									value={formData.memory_limit}
									onChange={(e) =>
										setFormData({
											...formData,
											memory_limit: parseInt(e.target.value),
										})
									}
								/>
							</div>
						</div>
						<div className="grid gap-2">
							<label htmlFor="edit_restart_policy" className="text-sm font-medium">
								Restart Policy
							</label>
							<select
								id="edit_restart_policy"
								className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
								value={formData.restart_policy}
								onChange={(e) =>
									setFormData({ ...formData, restart_policy: e.target.value })
								}
							>
								<option value="no">No</option>
								<option value="on-failure">On Failure</option>
								<option value="always">Always</option>
								<option value="unless-stopped">Unless Stopped</option>
							</select>
						</div>
					</div>
					<DialogFooter>
						<Button
							variant="outline"
							onClick={() => {
								setIsEditDialogOpen(false);
								setSelectedNode(null);
								resetForm();
							}}
						>
							Cancel
						</Button>
						<Button onClick={handleUpdateNode}>Update Node</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</div>
	);
}
