<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { apiClient } from '$lib/api/client';
	import { toastStore } from '$lib/stores/toast';
	import Button from '../ui/Button.svelte';
	import IconButton from '../ui/IconButton.svelte';
	import { Copy, RefreshCw, Trash, X } from 'lucide-svelte';
	
	export let webhook: any = null;
	export let isLoading = false;

	let webhookRequests: any[] = [];
	let selectedRequest: any = null;
	let isLoadingRequests = false;
	let searchQuery = '';
	let filterMethod = 'ALL';
	let filterStatus = 'ALL';
	let showRequestDetails = false;
	let autoRefresh = true;
	let refreshInterval: any;

	// Reactive filtered requests
	$: filteredRequests = webhookRequests.filter(request => {
		const matchesSearch = !searchQuery || 
			request.headers?.['user-agent']?.toLowerCase().includes(searchQuery.toLowerCase()) ||
			request.headers?.['content-type']?.toLowerCase().includes(searchQuery.toLowerCase()) ||
			request.body?.toLowerCase().includes(searchQuery.toLowerCase()) ||
			request.query_params?.toString().toLowerCase().includes(searchQuery.toLowerCase());
		
		const matchesMethod = filterMethod === 'ALL' || request.method === filterMethod;
		const matchesStatus = filterStatus === 'ALL' || request.status_code.toString().startsWith(filterStatus);
		
		return matchesSearch && matchesMethod && matchesStatus;
	});

	onMount(() => {
		if (webhook) {
			loadWebhookRequests();
			if (autoRefresh) {
				startAutoRefresh();
			}
		}
	});

	onDestroy(() => {
		if (refreshInterval) {
			clearInterval(refreshInterval);
		}
	});

	function startAutoRefresh() {
		refreshInterval = setInterval(() => {
			if (autoRefresh && webhook) {
				loadWebhookRequests(false); // Silent refresh
			}
		}, 3000); // Refresh every 3 seconds
	}

	function stopAutoRefresh() {
		if (refreshInterval) {
			clearInterval(refreshInterval);
			refreshInterval = null;
		}
	}

	function toggleAutoRefresh() {
		autoRefresh = !autoRefresh;
		if (autoRefresh) {
			startAutoRefresh();
		} else {
			stopAutoRefresh();
		}
	}

	async function loadWebhookRequests(showLoading = true) {
		if (!webhook) return;
		
		if (showLoading) isLoadingRequests = true;
		
		try {
			const response = await apiClient.getWebhookRequests(webhook.id);
			webhookRequests = response.requests || response || [];
		} catch (error: any) {
			console.error('Error loading webhook requests:', error);
			if (showLoading) {
				toastStore.error('Failed to load webhook requests');
			}
		} finally {
			if (showLoading) isLoadingRequests = false;
		}
	}

	function selectRequest(request: any) {
		selectedRequest = request;
		showRequestDetails = true;
	}

	function closeRequestDetails() {
		showRequestDetails = false;
		selectedRequest = null;
	}

	function copyWebhookUrl() {
		if (webhook?.url) {
			navigator.clipboard.writeText(webhook.url);
			toastStore.success('Webhook URL copied to clipboard!');
		}
	}

	function formatTimestamp(timestamp: string) {
		return new Date(timestamp).toLocaleString();
	}

	function getStatusColor(statusCode: number) {
		if (statusCode >= 200 && statusCode < 300) return 'text-green-400';
		if (statusCode >= 300 && statusCode < 400) return 'text-yellow-400';
		if (statusCode >= 400) return 'text-red-400';
		return 'text-gray-400';
	}

	function getMethodColor(method: string) {
		switch (method?.toLowerCase()) {
			case 'get': return 'bg-green-700 text-green-300';
			case 'post': return 'bg-blue-700 text-blue-300';
			case 'put': return 'bg-yellow-700 text-yellow-300';
			case 'patch': return 'bg-orange-700 text-orange-300';
			case 'delete': return 'bg-red-700 text-red-300';
			default: return 'bg-gray-700 text-gray-300';
		}
	}

	function formatJSON(obj: any) {
		try {
			return JSON.stringify(obj, null, 2);
		} catch {
			return obj?.toString() || '';
		}
	}

	function clearAllRequests() {
		if (confirm('Are you sure you want to clear all webhook requests? This action cannot be undone.')) {
			webhookRequests = [];
			selectedRequest = null;
			showRequestDetails = false;
			toastStore.success('All requests cleared');
		}
	}
</script>

{#if webhook}
	<div class="flex flex-col h-full bg-gray-900">
		<!-- Header -->
		<div class="flex items-center justify-between p-4 border-b border-gray-700 bg-gray-800">
			<div class="flex-1">
				<h1 class="text-xl font-bold text-white mb-2">{webhook.name}</h1>
				<div class="flex items-center space-x-4">
					<div class="flex items-center bg-gray-700/50 rounded-lg px-3 py-2">
						<span class="text-sm text-gray-300 mr-2">URL:</span>
						<code class="text-sm font-mono text-blue-400 bg-gray-800 px-2 py-1 rounded select-all">{webhook.url}</code>
						<IconButton
							on:click={copyWebhookUrl}
							size="sm"
							class="ml-2 p-1 text-gray-400 hover:text-white transition-colors"
							title="Copy URL"
						>
							<Copy class="w-4 h-4" />
						</IconButton>
					</div>
					<div class="flex items-center space-x-2">
						<span class="w-2 h-2 bg-green-400 rounded-full"></span>
						<span class="text-sm text-green-400">Active</span>
					</div>
				</div>
			</div>
			
			<!-- Controls -->
			<div class="flex items-center space-x-2">
				<Button
					on:click={toggleAutoRefresh}
					variant="secondary"
					size="sm"
					class="flex items-center space-x-2 px-3 py-1.5 text-sm bg-gray-700 hover:bg-gray-600 text-white rounded-md transition-colors {autoRefresh ? 'ring-1 ring-purple-500 bg-purple-600/20' : ''}"
				>
					<RefreshCw class="w-4 h-4" />
					<span>{autoRefresh ? 'Auto ON' : 'Auto OFF'}</span>
				</Button>
				
				<IconButton
					on:click={() => loadWebhookRequests()}
					size="sm"
					class="p-1.5 bg-purple-600 hover:bg-purple-700 text-white rounded-md transition-colors"
					title="Refresh requests"
				>
					<RefreshCw class="w-4 h-4" />
				</IconButton>
				
				<IconButton
					on:click={clearAllRequests}
					size="sm"
					class="p-1.5 bg-red-600 hover:bg-red-700 text-white rounded-md transition-colors"
					title="Clear all requests"
				>
					<Trash class="w-4 h-4" />
				</IconButton>
			</div>
		</div>

		<!-- Filters -->
		<div class="flex items-center space-x-3 p-3 bg-gray-800/50 border-b border-gray-700">
			<div class="flex-1 relative">
				<svg class="absolute left-3 top-2.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
				</svg>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search requests..."
					class="w-full bg-gray-700 border border-gray-600 text-white placeholder-gray-400 pl-10 pr-3 py-2 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
				/>
			</div>
			
			<select
				bind:value={filterMethod}
				class="bg-gray-700 border border-gray-600 text-white text-sm px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
			>
				<option value="ALL">All Methods</option>
				<option value="GET">GET</option>
				<option value="POST">POST</option>
				<option value="PUT">PUT</option>
				<option value="PATCH">PATCH</option>
				<option value="DELETE">DELETE</option>
			</select>
			
			<select
				bind:value={filterStatus}
				class="bg-gray-700 border border-gray-600 text-white text-sm px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
			>
				<option value="ALL">All Status</option>
				<option value="2">2xx Success</option>
				<option value="3">3xx Redirect</option>
				<option value="4">4xx Client Error</option>
				<option value="5">5xx Server Error</option>
			</select>
			
			<div class="text-xs text-gray-400 whitespace-nowrap">
				{filteredRequests.length} of {webhookRequests.length}
			</div>
		</div>

		<!-- Content -->
		<div class="flex flex-1 overflow-hidden">
			<!-- Requests List -->
			<div class="w-1/2 bg-gray-800 border-r border-gray-700 flex flex-col">
				<div class="p-3 border-b border-gray-700 bg-gray-800/50">
					<h2 class="text-base font-semibold text-white">Incoming Requests</h2>
				</div>
				
				<div class="flex-1 overflow-y-auto">
					{#if isLoadingRequests}
						<div class="flex items-center justify-center h-32">
							<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
						</div>
					{:else if filteredRequests.length === 0}
						<div class="flex flex-col items-center justify-center h-64 text-gray-400 p-6">
							<svg class="w-12 h-12 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path>
							</svg>
							<p class="text-base font-medium mb-1">No requests yet</p>
							<p class="text-sm text-center max-w-sm text-gray-500">
								Send HTTP requests to your webhook URL to see them appear here.
							</p>
						</div>
					{:else}
						{#each filteredRequests as request, index (request.id || index)}
							<div
								class="border-b border-gray-700 p-3 hover:bg-gray-750 cursor-pointer transition-colors {selectedRequest?.id === request.id ? 'bg-gray-750 border-l-2 border-l-purple-500' : ''}"
								on:click={() => selectRequest(request)}
								role="button"
								tabindex="0"
								on:keydown={(e) => e.key === 'Enter' && selectRequest(request)}
							>
								<div class="flex items-center justify-between mb-2">
									<div class="flex items-center space-x-3">
										<span class="text-xs font-mono px-2 py-1 rounded {getMethodColor(request.method)}">
											{request.method || 'GET'}
										</span>
										<span class="text-sm font-mono {getStatusColor(request.status_code || 200)}">
											{request.status_code || 200}
										</span>
									</div>
									<span class="text-xs text-gray-400">
										{formatTimestamp(request.created_at || new Date().toISOString())}
									</span>
								</div>
								
								<div class="text-sm text-gray-300 mb-1">
									{request.headers?.['user-agent'] || 'Unknown User Agent'}
								</div>
								
								{#if request.headers?.['content-type']}
									<div class="text-xs text-gray-400">
										{request.headers['content-type']}
									</div>
								{/if}
							</div>
						{/each}
					{/if}
				</div>
			</div>

			<!-- Request Details -->
			<div class="w-1/2 bg-gray-900 flex flex-col">
				{#if selectedRequest}
					<div class="p-3 border-b border-gray-700 bg-gray-800/50">
						<div class="flex items-center justify-between">
							<h2 class="text-base font-semibold text-white">Request Details</h2>
							<IconButton
								on:click={closeRequestDetails}
								size="sm"
								class="p-1 text-gray-400 hover:text-white transition-colors"
							>
								<X class="w-4 h-4" />
							</IconButton>
						</div>
					</div>
					
					<div class="flex-1 overflow-y-auto p-3 space-y-4">
						<!-- Request Summary -->
						<div>
							<h3 class="text-sm font-semibold text-gray-300 mb-2">Request Summary</h3>
							<div class="bg-gray-800/50 rounded-lg p-3 space-y-2">
								<div class="flex items-center space-x-4">
									<span class="text-xs font-mono px-2 py-1 rounded {getMethodColor(selectedRequest.method)}">
										{selectedRequest.method || 'GET'}
									</span>
									<span class="text-sm font-mono {getStatusColor(selectedRequest.status_code || 200)}">
										{selectedRequest.status_code || 200}
									</span>
									<span class="text-sm text-gray-400">
										{formatTimestamp(selectedRequest.created_at || new Date().toISOString())}
									</span>
								</div>
								{#if selectedRequest.query_params}
									<div class="text-sm text-gray-300">
										<span class="font-medium">Query:</span> {selectedRequest.query_params}
									</div>
								{/if}
							</div>
						</div>

						<!-- Headers -->
						<div>
							<h3 class="text-sm font-semibold text-gray-300 mb-2">Headers</h3>
							<div class="bg-gray-800/50 rounded-lg p-3">
								<pre class="text-xs text-gray-300 font-mono whitespace-pre-wrap">{formatJSON(selectedRequest.headers || {})}</pre>
							</div>
						</div>

						<!-- Body -->
						{#if selectedRequest.body}
							<div>
								<h3 class="text-sm font-semibold text-gray-300 mb-2">Body</h3>
								<div class="bg-gray-800/50 rounded-lg p-3">
									<pre class="text-xs text-gray-300 font-mono whitespace-pre-wrap">{selectedRequest.body}</pre>
								</div>
							</div>
						{/if}

						<!-- Response -->
						{#if selectedRequest.response}
							<div>
								<h3 class="text-sm font-semibold text-gray-300 mb-2">Response</h3>
								<div class="bg-gray-800/50 rounded-lg p-3">
									<pre class="text-xs text-gray-300 font-mono whitespace-pre-wrap">{formatJSON(selectedRequest.response)}</pre>
								</div>
							</div>
						{/if}
					</div>
				{:else}
					<div class="flex-1 flex items-center justify-center text-gray-400 p-6">
						<div class="text-center">
							<svg class="w-12 h-12 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
							</svg>
							<p class="text-base font-medium mb-1">Select a request</p>
							<p class="text-sm text-gray-500">Click on a request from the list to view its details</p>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>
{:else}
	<div class="flex-1 flex items-center justify-center text-gray-400">
		<div class="text-center">
			<svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path>
			</svg>
			<p class="text-lg font-medium">No webhook selected</p>
			<p class="text-sm">Select or create a webhook to start inspecting requests</p>
		</div>
	</div>
{/if}