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

	async updateCollection(collectionId: string, data: {name: string; description: string}) {
		const response = await this.client.put(`/collections/${collectionId}`, data);
		return response.data;
	}

	async deleteCollection(collectionId: string) {
		const response = await this.client.delete(`/collections/${collectionId}`);
		return response.data;
	}

	async toggleCollectionFavorite(collectionId: string) {
		const response = await this.client.post(`/collections/${collectionId}/favorite`);
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

	// Dashboard Analytics
	async getDashboardStats() {
		const response = await this.client.get('/dashboard/stats');
		return response.data;
	}

	async getApiHealth() {
		const response = await this.client.get('/dashboard/api-health');
		return response.data;
	}

	async getRecentActivity() {
		const response = await this.client.get('/dashboard/activity');
		return response.data;
	}

	async getTeamMembers() {
		const response = await this.client.get('/dashboard/team');
		return response.data;
	}

	async getNotifications() {
		const response = await this.client.get('/notifications');
		return response.data;
	}

	async markNotificationAsRead(notificationId: string) {
		const response = await this.client.patch(`/notifications/${notificationId}/read`);
		return response.data;
	}

	async markAllNotificationsAsRead() {
		const response = await this.client.patch('/notifications/mark-all-read');
		return response.data;
	}

	// Search functionality
	async searchWorkspace(workspaceId: string, query: string) {
		const response = await this.client.get(`/workspaces/${workspaceId}/search`, {
			params: { q: query }
		});
		return response.data;
	}

	// Workspace members
	async generateInvitationLink(workspaceId: string, role: string = 'member', expiresInHours: number = 168) {
		const response = await this.client.post(`/workspaces/${workspaceId}/invite-link`, {
			role,
			expiresInHours
		});
		return response.data;
	}

	async getWorkspaceMembers(workspaceId: string) {
		const response = await this.client.get(`/workspaces/${workspaceId}/members`);
		return response.data;
	}

	async removeWorkspaceMember(workspaceId: string, memberId: string) {
		const response = await this.client.delete(`/workspaces/${workspaceId}/members/${memberId}`);
		return response.data;
	}

	async updateWorkspaceMemberRole(workspaceId: string, memberId: string, role: string) {
		const response = await this.client.patch(`/workspaces/${workspaceId}/members/${memberId}`, {
			role
		});
		return response.data;
	}

	// Environment endpoints
	async getEnvironments(workspaceId: string) {
		const response = await this.client.get(`/workspace-environments/${workspaceId}`);
		return response.data;
	}

	async createEnvironment(workspaceId: string, name: string, variables: any[] = []) {
		const response = await this.client.post(`/workspace-environments/${workspaceId}`, {
			name,
			variables
		});
		return response.data;
	}

	async getEnvironment(environmentId: string) {
		const response = await this.client.get(`/environments/${environmentId}`);
		return response.data;
	}

	async updateEnvironment(environmentId: string, name: string, variables: any[]) {
		const response = await this.client.put(`/environments/${environmentId}`, {
			name,
			variables
		});
		return response.data;
	}

	async deleteEnvironment(environmentId: string) {
		const response = await this.client.delete(`/environments/${environmentId}`);
		return response.data;
	}

	// Webhook endpoints
	async getWebhooks(workspaceId: string, page = 1, limit = 10) {
		const response = await this.client.get(`/workspace-webhooks/${workspaceId}?page=${page}&limit=${limit}`);
		return response.data;
	}

	async createWebhook(workspaceId: string, data: any) {
		const response = await this.client.post(`/workspace-webhooks/${workspaceId}`, data);
		return response.data;
	}

	async getWebhook(webhookId: string) {
		const response = await this.client.get(`/webhooks/${webhookId}`);
		return response.data;
	}

	async updateWebhook(webhookId: string, data: any) {
		const response = await this.client.put(`/webhooks/${webhookId}`, data);
		return response.data;
	}

	async deleteWebhook(webhookId: string) {
		const response = await this.client.delete(`/webhooks/${webhookId}`);
		return response.data;
	}

	async getWebhookRequests(webhookId: string, page = 1, limit = 20) {
		const response = await this.client.get(`/webhooks/${webhookId}/requests?page=${page}&limit=${limit}`);
		return response.data;
	}

	async testWebhook(url: string, method = 'POST', headers: Record<string, string> = {}, body = '') {
		try {
			const response = await fetch(url, {
				method,
				headers: {
					'Content-Type': 'application/json',
					...headers
				},
				body: method !== 'GET' ? body : undefined
			});
			
			const responseText = await response.text();
			return {
				success: true,
				status: response.status,
				statusText: response.statusText,
				headers: Object.fromEntries(response.headers.entries()),
				body: responseText,
				timestamp: new Date().toISOString()
			};
		} catch (error: any) {
			return {
				success: false,
				status: 0,
				statusText: 'Network Error',
				headers: {},
				body: error.message,
				error: error.message,
				timestamp: new Date().toISOString()
			};
		}
	}
}

export const apiClient = new ApiClient();