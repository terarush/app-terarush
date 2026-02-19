import { apiClient } from "./client";

export interface Node {
	id: number;
	user_id: number;
	name: string;
	image: string;
	container_id: string;
	status: string;
	port: number;
	internal_port: number;
	cpu_limit: number;
	memory_limit: number;
	environment: Record<string, string>;
	volumes: Record<string, string>;
	command: string;
	restart_policy: string;
	created_at: string;
	updated_at: string;
}

export interface CreateNodeRequest {
	name: string;
	image: string;
	port: number;
	internal_port: number;
	cpu_limit?: number;
	memory_limit?: number;
	environment?: Record<string, string>;
	volumes?: Record<string, string>;
	command?: string;
	restart_policy?: string;
}

export interface UpdateNodeRequest {
	name?: string;
	cpu_limit?: number;
	memory_limit?: number;
	environment?: Record<string, string>;
	volumes?: Record<string, string>;
	command?: string;
	restart_policy?: string;
}

export interface NodeActionRequest {
	action: "start" | "stop" | "restart";
}

// Get all nodes for current user
export const getNodes = async (): Promise<Node[]> => {
	const response = await apiClient.get("/nodes");
	return response.data.data;
};

// Get a single node by ID
export const getNode = async (id: number): Promise<Node> => {
	const response = await apiClient.get(`/nodes/${id}`);
	return response.data.data;
};

// Create a new node
export const createNode = async (data: CreateNodeRequest): Promise<Node> => {
	const response = await apiClient.post("/nodes", data);
	return response.data.data;
};

// Update a node
export const updateNode = async (
	id: number,
	data: UpdateNodeRequest
): Promise<Node> => {
	const response = await apiClient.put(`/nodes/${id}`, data);
	return response.data.data;
};

// Delete a node
export const deleteNode = async (id: number): Promise<void> => {
	await apiClient.delete(`/nodes/${id}`);
};

// Perform action on node (start, stop, restart)
export const nodeAction = async (
	id: number,
	action: "start" | "stop" | "restart"
): Promise<void> => {
	await apiClient.post(`/nodes/${id}/action`, { action });
};
