<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { apiClient } from '$lib/api/client';
	import { toastStore } from '$lib/stores/toast';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import JsonEditor from '$lib/components/JsonEditor.svelte';

	export let webhook: any = null;
	export let isLoading = false;
	export let isSaving = false;

	const dispatch = createEventDispatcher();

	let activeTab = 'overview';
	let webhookRequests: any[] = [];
	let requestsLoading = false;
	let requestsPage = 1;
	let hasMoreRequests = false;

	// Test webhook state
	let testMethod = 'POST';
	let testHeaders: { key: string; value: string }[] = [{ key: '', value: '' }];
	let testBody = '{\n  "message": "Hello from Volt webhook test",\n  "timestamp": "' + new Date().toISOString() + '"\n}';
	let testQuery: { key: string; value: string }[] = [{ key: '', value: '' }];
	let testResponse: any = null;
	let isTesting = false;

	// Config editing state
	let editConfig: any = null;
	let showConfigEditor = false;

	$: if (webhook && activeTab === 'requests') {
		loadWebhookRequests();
	}

	async function loadWebhookRequests() {
		if (!webhook?.id) return;
		
		requestsLoading = true;
		try {
			const response = await apiClient.getWebhookRequests(webhook.id, requestsPage);
			webhookRequests = requestsPage === 1 ? response.requests : [...webhookRequests, ...response.requests];
			hasMoreRequests = response.has_more;
		} catch (error: any) {
			toastStore.error('Failed to load webhook requests');
		} finally {
			requestsLoading = false;
		}
	}

	async function loadMoreRequests() {
		if (hasMoreRequests && !requestsLoading) {
			requestsPage++;
			await loadWebhookRequests();
		}
	}

	async function testWebhook() {
		if (!webhook?.url) return;

		isTesting = true;
		testResponse = null;

		try {
			// Build headers object
			const headers: Record<string, string> = {};
			testHeaders.forEach(({ key, value }) => {
				if (key.trim() && value.trim()) {
					headers[key] = value;
				}
			});

			// Build query string
			const queryParams = new URLSearchParams();
			testQuery.forEach(({ key, value }) => {
				if (key.trim() && value.trim()) {
					queryParams.append(key, value);
				}
			});

			let testUrl = webhook.url;
			if (queryParams.toString()) {
				testUrl += '?' + queryParams.toString();
			}

			const startTime = Date.now();
			const response = await apiClient.testWebhook(testUrl, testMethod, headers, testBody);
			const endTime = Date.now();

			testResponse = {
				...response,
				response_time: endTime - startTime
			};

			if (response.success) {
				toastStore.success('Webhook test completed successfully');
			} else {
				toastStore.error('Webhook test failed: ' + (response.error || 'Unknown error'));
			}

			// Refresh requests if on requests tab
			if (activeTab === 'requests') {
				requestsPage = 1;
				await loadWebhookRequests();
			}
		} catch (error: any) {
			testResponse = {
				success: false,
				status: 0,
				statusText: 'Network Error',
				headers: {},
				body: error.message,
				error: error.message,
				timestamp: new Date().toISOString(),
				response_time: 0
			};
			toastStore.error('Failed to test webhook: ' + error.message);
		} finally {
			isTesting = false;
		}
	}

	function addTestHeader() {
		testHeaders = [...testHeaders, { key: '', value: '' }];
	}

	function removeTestHeader(index: number) {
		testHeaders = testHeaders.filter((_, i) => i !== index);
	}

	function addTestQuery() {
		testQuery = [...testQuery, { key: '', value: '' }];
	}

	function removeTestQuery(index: number) {
		testQuery = testQuery.filter((_, i) => i !== index);
	}

	function copyWebhookUrl() {
		if (webhook?.url) {
			navigator.clipboard.writeText(webhook.url);
			toastStore.success('Webhook URL copied to clipboard');
		}
	}

	function editWebhookConfig() {
		editConfig = JSON.parse(JSON.stringify(webhook.config || {}));
		showConfigEditor = true;
	}

	async function saveWebhookConfig() {
		if (!webhook?.id || !editConfig) return;

		try {
			await apiClient.updateWebhook(webhook.id, { config: editConfig });
			toastStore.success('Webhook configuration updated');
			showConfigEditor = false;
			dispatch('webhookUpdated');
		} catch (error: any) {
			toastStore.error('Failed to update webhook configuration');
		}
	}

	async function toggleWebhookStatus() {
		if (!webhook?.id) return;

		try {
			const newStatus = webhook.status === 'active' ? 'inactive' : 'active';
			await apiClient.updateWebhook(webhook.id, { status: newStatus });
			toastStore.success(`Webhook ${newStatus}`);
			dispatch('webhookUpdated');
		} catch (error: any) {
			toastStore.error('Failed to update webhook status');
		}
	}

	async function deleteWebhook() {
		if (!webhook?.id) return;

		const confirmed = confirm(`Are you sure you want to delete webhook "${webhook.name}"? This action cannot be undone.`);
		if (!confirmed) return;

		try {
			await apiClient.deleteWebhook(webhook.id);
			toastStore.success('Webhook deleted successfully');
			dispatch('webhookDeleted');
		} catch (error: any) {
			toastStore.error('Failed to delete webhook');
		}
	}

	function formatBytes(bytes: number) {
		if (bytes === 0) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}

	function formatTimestamp(timestamp: string) {
		return new Date(timestamp).toLocaleString();
	}

	function getStatusColor(status: string) {
		switch (status) {
			case 'active':
				return 'text-green-400 bg-green-900/20';
			case 'inactive':
				return 'text-gray-400 bg-gray-900/20';
			default:
				return 'text-yellow-400 bg-yellow-900/20';
		}
	}

	function getMethodColor(method: string) {
		switch (method.toUpperCase()) {
			case 'GET':
				return 'text-green-300 bg-green-900/20';
			case 'POST':
				return 'text-blue-300 bg-blue-900/20';
			case 'PUT':
				return 'text-yellow-300 bg-yellow-900/20';
			case 'DELETE':
				return 'text-red-300 bg-red-900/20';
			case 'PATCH':
				return 'text-purple-300 bg-purple-900/20';
			default:
				return 'text-gray-300 bg-gray-900/20';
		}
	}

	// Reset requests when webhook changes
	$: if (webhook) {
		webhookRequests = [];
		requestsPage = 1;
		hasMoreRequests = false;
	}
</script>

{#if isLoading}
	<div class="flex-1 flex items-center justify-center">
		<div class="text-center">
			<LoadingSpinner size="lg" />
			<p class="text-gray-400 text-sm mt-2">Loading webhook...</p>
		</div>
	</div>
{:else if !webhook}
	<div class="flex-1 flex items-center justify-center">
		<div class="text-center max-w-md px-6">
			<svg class="mx-auto h-16 w-16 text-gray-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path>
			</svg>
			<h3 class="text-xl font-medium text-white mb-2">Select a Webhook</h3>
			<p class="text-gray-400">Choose a webhook from the sidebar to view its details, test it, and monitor incoming requests.</p>
		</div>
	</div>
{:else}
	<div class="flex-1 flex flex-col bg-gray-900">
		<!-- Header -->
		<div class="border-b border-gray-700 bg-gray-800">
			<div class="p-4">
				<div class="flex items-center justify-between mb-4">
					<div class="flex items-center space-x-3">
						<div class="w-10 h-10 rounded-lg bg-purple-600/20 flex items-center justify-center">
							<svg class="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path>
							</svg>
						</div>
						<div>
							<h1 class="text-xl font-semibold text-white">{webhook.name}</h1>
							{#if webhook.description}
								<p class="text-gray-400 text-sm">{webhook.description}</p>
							{/if}
						</div>
					</div>
					<div class="flex items-center space-x-2">
						<span class="px-2 py-1 rounded-full text-xs font-medium {getStatusColor(webhook.status)}">
							{webhook.status}
						</span>
						<button
							on:click={toggleWebhookStatus}
							class="px-3 py-1.5 text-xs font-medium rounded {
								webhook.status === 'active' 
									? 'bg-red-600 hover:bg-red-700 text-white' 
									: 'bg-green-600 hover:bg-green-700 text-white'
							} transition-colors"
						>
							{webhook.status === 'active' ? 'Deactivate' : 'Activate'}
						</button>
						<button
							on:click={deleteWebhook}
							class="px-3 py-1.5 bg-red-600 hover:bg-red-700 text-white text-xs font-medium rounded transition-colors"
						>
							Delete
						</button>
					</div>
				</div>

				<!-- Stats -->
				<div class="grid grid-cols-4 gap-4 mb-4">
					<div class="bg-gray-700/50 rounded-lg p-3">
						<div class="text-2xl font-bold text-white">{webhook.stats?.total_requests || 0}</div>
						<div class="text-xs text-gray-400">Total Requests</div>
					</div>
					<div class="bg-gray-700/50 rounded-lg p-3">
						<div class="text-2xl font-bold text-green-400">{webhook.stats?.success_requests || 0}</div>
						<div class="text-xs text-gray-400">Success</div>
					</div>
					<div class="bg-gray-700/50 rounded-lg p-3">
						<div class="text-2xl font-bold text-red-400">{webhook.stats?.failed_requests || 0}</div>
						<div class="text-xs text-gray-400">Failed</div>
					</div>
					<div class="bg-gray-700/50 rounded-lg p-3">
						<div class="text-2xl font-bold text-blue-400">{webhook.stats?.average_response || 0}ms</div>
						<div class="text-xs text-gray-400">Avg Response</div>
					</div>
				</div>

				<!-- Tabs -->
				<nav class="flex space-x-4">
					{#each [
						{ id: 'overview', name: 'Overview', icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z' },
						{ id: 'test', name: 'Test', icon: 'M13 10V3L4 14h7v7l9-11h-7z' },
						{ id: 'requests', name: 'Requests', icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z' },
						{ id: 'config', name: 'Configuration', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z' }
					] as tab}
						<button
							on:click={() => activeTab = tab.id}
							class="flex items-center px-3 py-2 text-sm font-medium rounded-lg transition-colors {
								activeTab === tab.id 
									? 'bg-purple-600 text-white' 
									: 'text-gray-300 hover:text-white hover:bg-gray-700'
							}"
						>
							<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={tab.icon}></path>
							</svg>
							{tab.name}
						</button>
					{/each}
				</nav>
			</div>
		</div>

		<!-- Tab Content -->
		<div class="flex-1 overflow-y-auto p-4">
			{#if activeTab === 'overview'}
				<!-- Overview Tab -->
				<div class="space-y-6">
					<!-- Webhook URL -->
					<div class="bg-gray-800 rounded-lg p-4">
						<h3 class="text-lg font-medium text-white mb-3">Webhook URL</h3>
						<div class="flex items-center space-x-2">
							<input
								type="text"
								value={webhook.url}
								readonly
								class="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white font-mono"
							/>
							<button
								on:click={copyWebhookUrl}
								class="px-3 py-2 bg-purple-600 hover:bg-purple-700 text-white text-sm rounded transition-colors"
							>
								Copy
							</button>
						</div>
						<p class="text-xs text-gray-400 mt-2">Send HTTP requests to this URL to test your webhook</p>
					</div>

					<!-- Methods and Configuration -->
					<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
						<div class="bg-gray-800 rounded-lg p-4">
							<h3 class="text-lg font-medium text-white mb-3">Allowed Methods</h3>
							<div class="flex flex-wrap gap-2">
								{#each webhook.config?.methods || ['POST'] as method}
									<span class="px-2 py-1 rounded text-xs font-medium {getMethodColor(method)}">
										{method}
									</span>
								{/each}
							</div>
						</div>

						<div class="bg-gray-800 rounded-lg p-4">
							<h3 class="text-lg font-medium text-white mb-3">Response Configuration</h3>
							<div class="space-y-2 text-sm">
								<div class="flex justify-between">
									<span class="text-gray-400">Status Code:</span>
									<span class="text-white">{webhook.config?.response_status || 200}</span>
								</div>
								<div class="flex justify-between">
									<span class="text-gray-400">Content Type:</span>
									<span class="text-white">{webhook.config?.content_type || 'application/json'}</span>
								</div>
								<div class="flex justify-between">
									<span class="text-gray-400">Response Delay:</span>
									<span class="text-white">{webhook.config?.response_delay || 0}ms</span>
								</div>
								<div class="flex justify-between">
									<span class="text-gray-400">CORS Enabled:</span>
									<span class="text-white">{webhook.config?.cors_enabled ? 'Yes' : 'No'}</span>
								</div>
							</div>
						</div>
					</div>

					<!-- Authentication -->
					{#if webhook.config?.authentication?.enabled}
						<div class="bg-gray-800 rounded-lg p-4">
							<h3 class="text-lg font-medium text-white mb-3">Authentication</h3>
							<div class="space-y-2 text-sm">
								<div class="flex justify-between">
									<span class="text-gray-400">Type:</span>
									<span class="text-white capitalize">{webhook.config.authentication.type || 'None'}</span>
								</div>
								{#if webhook.config.authentication.type === 'basic'}
									<div class="flex justify-between">
										<span class="text-gray-400">Username:</span>
										<span class="text-white">{webhook.config.authentication.username}</span>
									</div>
								{/if}
							</div>
						</div>
					{/if}
				</div>

			{:else if activeTab === 'test'}
				<!-- Test Tab -->
				<div class="max-w-4xl space-y-6">
					<div class="bg-gray-800 rounded-lg p-4">
						<h3 class="text-lg font-medium text-white mb-4">Test Webhook</h3>

						<!-- Method Selection -->
						<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
							<div>
								<label class="block text-sm font-medium text-gray-300 mb-2">HTTP Method</label>
								<select bind:value={testMethod} class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white">
									{#each webhook.config?.methods || ['POST'] as method}
										<option value={method}>{method}</option>
									{/each}
								</select>
							</div>
						</div>

						<!-- Headers -->
						<div class="mb-6">
							<div class="flex items-center justify-between mb-2">
								<label class="block text-sm font-medium text-gray-300">Headers</label>
								<button
									on:click={addTestHeader}
									class="text-xs text-purple-400 hover:text-purple-300"
								>
									Add Header
								</button>
							</div>
							<div class="space-y-2">
								{#each testHeaders as header, index}
									<div class="flex space-x-2">
										<input
											type="text"
											bind:value={header.key}
											placeholder="Header name"
											class="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white"
										/>
										<input
											type="text"
											bind:value={header.value}
											placeholder="Header value"
											class="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white"
										/>
										<button
											on:click={() => removeTestHeader(index)}
											class="px-2 py-2 text-red-400 hover:text-red-300"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
											</svg>
										</button>
									</div>
								{/each}
							</div>
						</div>

						<!-- Query Parameters -->
						<div class="mb-6">
							<div class="flex items-center justify-between mb-2">
								<label class="block text-sm font-medium text-gray-300">Query Parameters</label>
								<button
									on:click={addTestQuery}
									class="text-xs text-purple-400 hover:text-purple-300"
								>
									Add Parameter
								</button>
							</div>
							<div class="space-y-2">
								{#each testQuery as param, index}
									<div class="flex space-x-2">
										<input
											type="text"
											bind:value={param.key}
											placeholder="Parameter name"
											class="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white"
										/>
										<input
											type="text"
											bind:value={param.value}
											placeholder="Parameter value"
											class="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white"
										/>
										<button
											on:click={() => removeTestQuery(index)}
											class="px-2 py-2 text-red-400 hover:text-red-300"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
											</svg>
										</button>
									</div>
								{/each}
							</div>
						</div>

						<!-- Request Body -->
						{#if testMethod !== 'GET'}
							<div class="mb-6">
								<label class="block text-sm font-medium text-gray-300 mb-2">Request Body</label>
								<JsonEditor
									bind:value={testBody}
									height="200px"
									placeholder="Enter request body (JSON)"
								/>
							</div>
						{/if}

						<!-- Test Button -->
						<button
							on:click={testWebhook}
							disabled={isTesting}
							class="w-full py-3 bg-purple-600 hover:bg-purple-700 disabled:opacity-50 text-white font-medium rounded transition-colors flex items-center justify-center"
						>
							{#if isTesting}
								<LoadingSpinner size="sm" color="white" />
								<span class="ml-2">Testing...</span>
							{:else}
								<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
								</svg>
								Test Webhook
							{/if}
						</button>
					</div>

					<!-- Test Response -->
					{#if testResponse}
						<div class="bg-gray-800 rounded-lg p-4">
							<h3 class="text-lg font-medium text-white mb-4">Response</h3>
							
							<div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
								<div class="text-center">
									<div class="text-2xl font-bold {testResponse.success ? 'text-green-400' : 'text-red-400'}">
										{testResponse.status}
									</div>
									<div class="text-xs text-gray-400">Status</div>
								</div>
								<div class="text-center">
									<div class="text-2xl font-bold text-blue-400">{testResponse.response_time || 0}ms</div>
									<div class="text-xs text-gray-400">Time</div>
								</div>
								<div class="text-center">
									<div class="text-2xl font-bold text-yellow-400">{formatBytes(testResponse.body?.length || 0)}</div>
									<div class="text-xs text-gray-400">Size</div>
								</div>
								<div class="text-center">
									<div class="text-2xl font-bold {testResponse.success ? 'text-green-400' : 'text-red-400'}">
										{testResponse.success ? '✓' : '✗'}
									</div>
									<div class="text-xs text-gray-400">Success</div>
								</div>
							</div>

							{#if testResponse.headers && Object.keys(testResponse.headers).length > 0}
								<div class="mb-4">
									<h4 class="text-sm font-medium text-gray-300 mb-2">Response Headers</h4>
									<div class="bg-gray-900 rounded p-3 font-mono text-xs">
										{#each Object.entries(testResponse.headers) as [key, value]}
											<div class="flex">
												<span class="text-purple-400 mr-2">{key}:</span>
												<span class="text-gray-300">{value}</span>
											</div>
										{/each}
									</div>
								</div>
							{/if}

							<div>
								<h4 class="text-sm font-medium text-gray-300 mb-2">Response Body</h4>
								<div class="bg-gray-900 rounded p-3 font-mono text-xs text-gray-300 max-h-96 overflow-y-auto">
									{testResponse.body || 'No response body'}
								</div>
							</div>
						</div>
					{/if}
				</div>

			{:else if activeTab === 'requests'}
				<!-- Requests Tab -->
				<div class="space-y-4">
					<div class="flex items-center justify-between">
						<h3 class="text-lg font-medium text-white">Request History</h3>
						<button
							on:click={() => { requestsPage = 1; loadWebhookRequests(); }}
							class="px-3 py-1.5 bg-purple-600 hover:bg-purple-700 text-white text-sm rounded transition-colors"
						>
							Refresh
						</button>
					</div>

					{#if requestsLoading && webhookRequests.length === 0}
						<div class="flex items-center justify-center py-12">
							<LoadingSpinner size="lg" />
						</div>
					{:else if webhookRequests.length === 0}
						<div class="text-center py-12">
							<svg class="mx-auto h-12 w-12 text-gray-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
							</svg>
							<p class="text-gray-400">No requests received yet</p>
							<p class="text-gray-500 text-sm mt-1">Start sending requests to your webhook URL to see them here</p>
						</div>
					{:else}
						<div class="space-y-2">
							{#each webhookRequests as request}
								<div class="bg-gray-800 rounded-lg p-4">
									<div class="flex items-center justify-between mb-3">
										<div class="flex items-center space-x-3">
											<span class="px-2 py-1 rounded text-xs font-medium {getMethodColor(request.method)}">
												{request.method}
											</span>
											<span class="text-white font-mono text-sm">{request.path}</span>
											<span class="text-gray-400 text-sm">{formatTimestamp(request.timestamp)}</span>
										</div>
										<div class="flex items-center space-x-2 text-sm">
											<span class="text-gray-400">{request.response_time}ms</span>
											<span class="text-gray-400">{formatBytes(request.size || 0)}</span>
											<span class="text-gray-400">{request.ip}</span>
										</div>
									</div>

									{#if request.query && Object.keys(request.query).length > 0}
										<div class="mb-3">
											<h5 class="text-xs font-medium text-gray-400 mb-1">Query Parameters</h5>
											<div class="bg-gray-900 rounded p-2 font-mono text-xs">
												{#each Object.entries(request.query) as [key, value]}
													<div><span class="text-purple-400">{key}:</span> <span class="text-gray-300">{value}</span></div>
												{/each}
											</div>
										</div>
									{/if}

									{#if request.headers && Object.keys(request.headers).length > 0}
										<div class="mb-3">
											<h5 class="text-xs font-medium text-gray-400 mb-1">Headers</h5>
											<div class="bg-gray-900 rounded p-2 font-mono text-xs max-h-32 overflow-y-auto">
												{#each Object.entries(request.headers) as [key, value]}
													<div><span class="text-purple-400">{key}:</span> <span class="text-gray-300">{value}</span></div>
												{/each}
											</div>
										</div>
									{/if}

									{#if request.body}
										<div>
											<h5 class="text-xs font-medium text-gray-400 mb-1">Request Body</h5>
											<div class="bg-gray-900 rounded p-2 font-mono text-xs text-gray-300 max-h-32 overflow-y-auto">
												{request.body}
											</div>
										</div>
									{/if}
								</div>
							{/each}

							{#if hasMoreRequests}
								<button
									on:click={loadMoreRequests}
									disabled={requestsLoading}
									class="w-full py-2 text-purple-400 hover:text-purple-300 text-sm font-medium"
								>
									{#if requestsLoading}
										Loading more...
									{:else}
										Load More Requests
									{/if}
								</button>
							{/if}
						</div>
					{/if}
				</div>

			{:else if activeTab === 'config'}
				<!-- Configuration Tab -->
				<div class="max-w-2xl space-y-6">
					<div class="bg-gray-800 rounded-lg p-4">
						<div class="flex items-center justify-between mb-4">
							<h3 class="text-lg font-medium text-white">Webhook Configuration</h3>
							<button
								on:click={editWebhookConfig}
								class="px-3 py-1.5 bg-purple-600 hover:bg-purple-700 text-white text-sm rounded transition-colors"
							>
								Edit Configuration
							</button>
						</div>

						<div class="bg-gray-900 rounded p-4 font-mono text-sm">
							<pre class="text-gray-300 whitespace-pre-wrap">{JSON.stringify(webhook.config || {}, null, 2)}</pre>
						</div>
					</div>

					<!-- Quick Settings -->
					<div class="bg-gray-800 rounded-lg p-4">
						<h3 class="text-lg font-medium text-white mb-4">Quick Settings</h3>
						<div class="space-y-4">
							<div class="flex items-center justify-between">
								<div>
									<div class="text-white font-medium">Request Logging</div>
									<div class="text-gray-400 text-sm">Log all incoming requests for debugging</div>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input type="checkbox" checked={webhook.config?.log_requests} class="sr-only peer">
									<div class="w-11 h-6 bg-gray-600 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-600"></div>
								</label>
							</div>

							<div class="flex items-center justify-between">
								<div>
									<div class="text-white font-medium">CORS Enabled</div>
									<div class="text-gray-400 text-sm">Allow cross-origin requests</div>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input type="checkbox" checked={webhook.config?.cors_enabled} class="sr-only peer">
									<div class="w-11 h-6 bg-gray-600 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-600"></div>
								</label>
							</div>
						</div>
					</div>
				</div>
			{/if}
		</div>
	</div>
{/if}

<!-- Configuration Editor Modal -->
{#if showConfigEditor}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-gray-800 rounded-lg p-6 w-full max-w-2xl max-h-[80vh] overflow-y-auto">
			<div class="flex items-center justify-between mb-4">
				<h3 class="text-lg font-medium text-white">Edit Configuration</h3>
				<button
					on:click={() => showConfigEditor = false}
					class="text-gray-400 hover:text-white"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>

			<JsonEditor
				bind:value={editConfig}
				height="400px"
				placeholder="Enter webhook configuration (JSON)"
			/>

			<div class="flex justify-end space-x-2 mt-4">
				<button
					on:click={() => showConfigEditor = false}
					class="px-4 py-2 bg-gray-600 hover:bg-gray-700 text-white rounded transition-colors"
				>
					Cancel
				</button>
				<button
					on:click={saveWebhookConfig}
					class="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded transition-colors"
				>
					Save Configuration
				</button>
			</div>
		</div>
	</div>
{/if}