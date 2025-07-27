import axios from 'axios';
import type { AxiosInstance, AxiosResponse } from 'axios';
import { authStore } from '../stores/auth';
import { get } from 'svelte/store';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

class ApiClient {
	private client: AxiosInstance;

	constructor() {
		this.client = axios.create({
			baseURL: `${API_URL}/api`,
			headers: {
				'Content-Type': 'application/json'
			},
			withCredentials: true // Include cookies in requests
		});

		// Response interceptor for error handling
		this.client.interceptors.response.use(
			(response) => response,
			(error) => {
				if (error.response?.status === 401) {
					authStore.logout();
				}
				return Promise.reject(error);
			}
		);
	}

	// Auth endpoints
	async register(username: string, email: string, password: string) {
		const response = await this.client.post('/auth/register', {
			username,
			email,
			password
		});
		return response.data;
	}

	async login(email: string, password: string) {
		const response = await this.client.post('/auth/login', {
			email,
			password
		});
		return response.data;
	}

	async logout() {
		try {
			await this.client.delete('/auth/logout');
		} catch (error) {
			// Ignore error, logout locally anyway
		}
		authStore.logout();
	}

	// Workspaces
	async getWorkspaces() {
		const response = await this.client.get('/workspaces');
		return response.data;
	}

	async createWorkspace(name: string, description?: string) {
		const response = await this.client.post('/workspaces', {
			name,
			description
		});
		return response.data;
	}

	async getWorkspace(id: string) {
		const response = await this.client.get(`/workspace-detail/${id}`);
		return response.data;
	}

	// Collections
	async getCollections(workspaceId: string) {
		const response = await this.client.get(`/workspace-collections/${workspaceId}`);
		return response.data;
	}

	async createCollection(workspaceId: string, name: string, description?: string) {
		const response = await this.client.post(`/workspace-collections/${workspaceId}`, {
			name,
			description
		});
		return response.data;
	}

	// Requests
	async getRequests(collectionId: string) {
		const response = await this.client.get(`/collection-requests/${collectionId}`);
		return response.data;
	}

	async createRequest(collectionId: string, requestData: any) {
		const response = await this.client.post(`/collection-requests/${collectionId}`, requestData);
		return response.data;
	}

	async getRequest(requestId: string) {
		const response = await this.client.get(`/requests/${requestId}`);
		return response.data;
	}

	async updateRequest(requestId: string, requestData: any) {
		const response = await this.client.put(`/requests/${requestId}`, requestData);
		return response.data;
	}

	async deleteRequest(requestId: string) {
		const response = await this.client.delete(`/requests/${requestId}`);
		return response.data;
	}

	async executeRequest(requestId: string) {
		const response = await this.client.post(`/requests/${requestId}/execute`);
		return response.data;
	}

	async executeRawRequest(requestData: any) {
		const response = await this.client.post('/execute-request', requestData);
		return response.data;
	}
}

export const apiClient = new ApiClient();