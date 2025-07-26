<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { workspaceStore } from '$lib/stores/workspace';
	import { apiClient } from '$lib/api/client';

	let workspaceId = '';
	let isAuthenticated = false;
	let currentWorkspace = null;
	let collections = [];
	let selectedCollection = null;
	let selectedRequest = null;
	let requests = [];
	let isLoading = true;
	let error = '';
	
	// Request builder state
	let requestName = 'New Request';
	let requestMethod = 'GET';
	let requestUrl = '';
	let requestHeaders = [{ key: '', value: '' }];
	let requestBody = '';
	let requestBodyType = 'json';
	let authType = 'none';
	let authToken = '';
	let authUsername = '';
	let authPassword = '';
	let showAuthToken = false;
	let showAuthPassword = false;
	
	// Response state
	let response = null;
	let isExecuting = false;

	// UI state
	let showNewCollectionModal = false;
	let newCollectionName = '';
	let newCollectionDescription = '';
	let isCreatingCollection = false;

	// Subscribe to stores
	authStore.subscribe(state => {
		isAuthenticated = state.isAuthenticated;
	});

	workspaceStore.subscribe(state => {
		currentWorkspace = state.currentWorkspace;
		collections = state.collections;
		isLoading = state.isLoading;
		error = state.error || '';
	});

	// Get workspace ID from URL
	page.subscribe(({ params }) => {
		workspaceId = params.id;
	});

	onMount(async () => {
		// Wait for auth store to initialize
		const unsubscribe = authStore.subscribe(async (state) => {
			if (state.isAuthenticated === false && !state.user) {
				goto('/login');
				return;
			}
			
			if (state.isAuthenticated && state.user && workspaceId) {
				await loadWorkspace();
				await loadCollections();
				unsubscribe();
			}
		});
	});

	async function loadWorkspace() {
		try {
			const workspace = await apiClient.getWorkspace(workspaceId);
			workspaceStore.setCurrentWorkspace(workspace);
		} catch (err: any) {
			error = err.response?.data?.error || 'Failed to load workspace';
		}
	}

	async function loadCollections() {
		try {
			console.log('Loading collections for workspace:', workspaceId);
			workspaceStore.setLoading(true);
			const data = await apiClient.getCollections(workspaceId);
			console.log('Collections received:', data);
			workspaceStore.setCollections(data);
		} catch (err: any) {
			console.error('Error loading collections:', err);
			workspaceStore.setError(err.response?.data?.error || 'Failed to load collections');
		} finally {
			workspaceStore.setLoading(false);
			console.log('Collections loading finished');
		}
	}

	async function createCollection() {
		if (!newCollectionName.trim()) return;

		isCreatingCollection = true;
		try {
			const collection = await apiClient.createCollection(workspaceId, newCollectionName, newCollectionDescription);
			workspaceStore.addCollection(collection);
			
			newCollectionName = '';
			newCollectionDescription = '';
			showNewCollectionModal = false;
		} catch (err: any) {
			error = err.response?.data?.error || 'Failed to create collection';
		} finally {
			isCreatingCollection = false;
		}
	}

	async function selectCollection(collection: any) {
		selectedCollection = collection;
		selectedRequest = null;
		response = null;
		
		try {
			const data = await apiClient.getRequests(collection.id);
			requests = data;
		} catch (err: any) {
			error = err.response?.data?.error || 'Failed to load requests';
		}
	}

	function selectRequest(request: any) {
		selectedRequest = request;
		
		// Populate request builder with selected request data
		requestName = request.name;
		requestMethod = request.method;
		requestUrl = request.url;
		requestBody = request.body?.content || '';
		requestBodyType = request.body?.type || 'json';
		authType = request.auth?.type || 'none';
		
		// Convert headers object to array
		requestHeaders = Object.entries(request.headers || {}).map(([key, value]) => ({ key, value }));
		if (requestHeaders.length === 0) {
			requestHeaders = [{ key: '', value: '' }];
		}
		
		// Set auth credentials
		if (request.auth?.credentials) {
			authToken = request.auth.credentials.token || '';
			authUsername = request.auth.credentials.username || '';
			authPassword = request.auth.credentials.password || '';
		}
	}

	function addHeader() {
		requestHeaders = [...requestHeaders, { key: '', value: '' }];
	}

	function removeHeader(index: number) {
		requestHeaders = requestHeaders.filter((_, i) => i !== index);
	}

	async function executeRequest() {
		if (!selectedRequest) {
			await saveRequest();
		}
		
		if (!selectedRequest) return;

		isExecuting = true;
		response = null;
		
		try {
			const result = await apiClient.executeRequest(selectedRequest.id);
			response = result;
		} catch (err: any) {
			error = err.response?.data?.error || 'Failed to execute request';
		} finally {
			isExecuting = false;
		}
	}

	async function saveRequest() {
		if (!selectedCollection) return;

		const headers = {};
		requestHeaders.forEach(({ key, value }) => {
			if (key.trim() && value.trim()) {
				headers[key] = value;
			}
		});

		const auth = {
			type: authType,
			credentials: {}
		};

		if (authType === 'bearer' && authToken) {
			auth.credentials.token = authToken;
		} else if (authType === 'basic' && authUsername) {
			auth.credentials.username = authUsername;
			auth.credentials.password = authPassword;
		}

		const requestData = {
			name: requestName,
			method: requestMethod,
			url: requestUrl,
			headers,
			body: {
				type: requestBodyType,
				content: requestBody
			},
			auth,
			tests: []
		};

		try {
			if (selectedRequest) {
				// Update existing request
				await apiClient.updateRequest(selectedRequest.id, requestData);
				const updatedRequests = requests.map(r => 
					r.id === selectedRequest.id ? { ...r, ...requestData } : r
				);
				requests = updatedRequests;
				selectedRequest = { ...selectedRequest, ...requestData };
			} else {
				// Create new request
				const newRequest = await apiClient.createRequest(selectedCollection.id, requestData);
				requests = [...requests, newRequest];
				selectedRequest = newRequest;
			}
		} catch (err: any) {
			error = err.response?.data?.error || 'Failed to save request';
		}
	}

	function newRequest() {
		selectedRequest = null;
		requestName = 'New Request';
		requestMethod = 'GET';
		requestUrl = '';
		requestHeaders = [{ key: '', value: '' }];
		requestBody = '';
		requestBodyType = 'json';
		authType = 'none';
		authToken = '';
		authUsername = '';
		authPassword = '';
		response = null;
	}

	function formatResponseBody(body: string) {
		try {
			return JSON.stringify(JSON.parse(body), null, 2);
		} catch {
			return body;
		}
	}
</script>

<svelte:head>
	<title>{currentWorkspace?.name || 'Workspace'} - Volt</title>
</svelte:head>

<div class="h-screen flex flex-col bg-gray-50">
	<!-- Header -->
	<header class="bg-white shadow-sm border-b border-gray-200 flex-shrink-0">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between items-center h-16">
				<div class="flex items-center space-x-4">
					<a href="/dashboard" class="text-2xl font-bold electric-gradient bg-clip-text text-transparent">⚡ Volt</a>
					<span class="text-gray-300">→</span>
					<h1 class="text-lg font-medium text-gray-900">{currentWorkspace?.name || 'Loading...'}</h1>
				</div>
				
				<div class="flex items-center space-x-4">
					<button
						on:click={() => showNewCollectionModal = true}
						class="text-sm bg-blue-600 text-white px-3 py-1.5 rounded-md hover:bg-blue-700 transition-colors"
					>
						New Collection
					</button>
				</div>
			</div>
		</div>
	</header>

	<div class="flex-1 flex overflow-hidden">
		<!-- Sidebar -->
		<div class="w-80 bg-white border-r border-gray-200 flex flex-col">
			<!-- Collections -->
			<div class="flex-1 overflow-y-auto">
				<div class="p-4">
					<h2 class="text-sm font-medium text-gray-900 mb-3">Collections</h2>
					
					{#if isLoading}
						<div class="flex items-center justify-center py-8">
							<div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
						</div>
					{:else if collections.length === 0}
						<div class="text-center py-8">
							<p class="text-sm text-gray-500 mb-4">No collections yet</p>
							<button
								on:click={() => showNewCollectionModal = true}
								class="text-sm text-blue-600 hover:text-blue-700"
							>
								Create your first collection
							</button>
						</div>
					{:else}
						<div class="space-y-2">
							{#each collections as collection}
								<div 
									class="p-3 rounded-lg border cursor-pointer transition-colors {selectedCollection?.id === collection.id ? 'bg-blue-50 border-blue-200' : 'bg-gray-50 border-gray-200 hover:bg-gray-100'}"
									on:click={() => selectCollection(collection)}
									on:keydown={(e) => e.key === 'Enter' && selectCollection(collection)}
									role="button"
									tabindex="0"
								>
									<h3 class="font-medium text-sm text-gray-900">{collection.name}</h3>
									{#if collection.description}
										<p class="text-xs text-gray-500 mt-1">{collection.description}</p>
									{/if}
								</div>
							{/each}
						</div>
					{/if}
				</div>

				<!-- Requests -->
				{#if selectedCollection}
					<div class="border-t border-gray-200 p-4">
						<div class="flex items-center justify-between mb-3">
							<h3 class="text-sm font-medium text-gray-900">Requests</h3>
							<button
								on:click={newRequest}
								class="text-xs text-blue-600 hover:text-blue-700 flex items-center"
							>
								<svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
								</svg>
								New
							</button>
						</div>
						
						<div class="space-y-1">
							{#each requests as request}
								<div 
									class="p-2 rounded cursor-pointer transition-colors flex items-center {selectedRequest?.id === request.id ? 'bg-blue-100' : 'hover:bg-gray-100'}"
									on:click={() => selectRequest(request)}
									on:keydown={(e) => e.key === 'Enter' && selectRequest(request)}
									role="button"
									tabindex="0"
								>
									<span class="text-xs font-mono px-1.5 py-0.5 rounded {request.method === 'GET' ? 'bg-green-100 text-green-800' : request.method === 'POST' ? 'bg-blue-100 text-blue-800' : request.method === 'PUT' ? 'bg-yellow-100 text-yellow-800' : 'bg-red-100 text-red-800'}">{request.method}</span>
									<span class="ml-2 text-sm text-gray-900 truncate">{request.name}</span>
								</div>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		</div>

		<!-- Main Content -->
		<div class="flex-1 flex flex-col">
			{#if selectedCollection}
				<!-- Request Builder -->
				<div class="flex-1 flex">
					<!-- Request Panel -->
					<div class="flex-1 flex flex-col bg-white">
						<div class="border-b border-gray-200 p-4">
							<div class="flex items-center space-x-3">
								<input
									type="text"
									bind:value={requestName}
									placeholder="Request name"
									class="flex-1 text-lg font-medium border-none outline-none bg-transparent"
								/>
								<button
									on:click={saveRequest}
									class="text-sm bg-gray-100 hover:bg-gray-200 px-3 py-1 rounded transition-colors"
								>
									Save
								</button>
								<button
									on:click={executeRequest}
									disabled={isExecuting}
									class="electric-gradient text-white px-4 py-2 rounded font-medium hover:opacity-90 transition-opacity disabled:opacity-50 flex items-center"
								>
									{#if isExecuting}
										<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
											<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
											<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
										</svg>
										Sending...
									{:else}
										Send
									{/if}
								</button>
							</div>
						</div>

						<div class="flex-1 overflow-y-auto p-4 space-y-6">
							<!-- URL -->
							<div>
								<div class="flex space-x-2">
									<select bind:value={requestMethod} class="px-3 py-2 border border-gray-300 rounded-md text-sm font-mono">
										<option value="GET">GET</option>
										<option value="POST">POST</option>
										<option value="PUT">PUT</option>
										<option value="DELETE">DELETE</option>
										<option value="PATCH">PATCH</option>
									</select>
									<input
										type="url"
										bind:value={requestUrl}
										placeholder="https://api.example.com/endpoint"
										class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
									/>
								</div>
							</div>

							<!-- Headers -->
							<div>
								<h3 class="text-sm font-medium text-gray-700 mb-2">Headers</h3>
								<div class="space-y-2">
									{#each requestHeaders as header, index}
										<div class="flex space-x-2">
											<input
												type="text"
												bind:value={header.key}
												placeholder="Header name"
												class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm"
											/>
											<input
												type="text"
												bind:value={header.value}
												placeholder="Header value"
												class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm"
											/>
											<button
												on:click={() => removeHeader(index)}
												class="px-2 py-2 text-gray-400 hover:text-red-500"
											>
												<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
												</svg>
											</button>
										</div>
									{/each}
									<button
										on:click={addHeader}
										class="text-sm text-blue-600 hover:text-blue-700 flex items-center"
									>
										<svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
										</svg>
										Add Header
									</button>
								</div>
							</div>

							<!-- Auth -->
							<div>
								<h3 class="text-sm font-medium text-gray-700 mb-2">Authentication</h3>
								<select bind:value={authType} class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm mb-3">
									<option value="none">No Auth</option>
									<option value="bearer">Bearer Token</option>
									<option value="basic">Basic Auth</option>
								</select>
								
								{#if authType === 'bearer'}
									<div class="relative">
										{#if showAuthToken}
											<input
												type="text"
												bind:value={authToken}
												placeholder="Bearer token"
												class="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md text-sm"
											/>
										{:else}
											<input
												type="password"
												bind:value={authToken}
												placeholder="Bearer token"
												class="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md text-sm"
											/>
										{/if}
										<button
											type="button"
											class="absolute inset-y-0 right-0 pr-3 flex items-center"
											on:click={() => showAuthToken = !showAuthToken}
										>
											{#if showAuthToken}
												<!-- Eye slash icon (hide) -->
												<svg class="h-4 w-4 text-gray-400 hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L8.464 8.464M14.12 14.12l1.415 1.415M14.12 14.12L9.878 9.878m4.242 4.242L8.464 8.464m5.656 5.656l1.415 1.415"></path>
												</svg>
											{:else}
												<!-- Eye icon (show) -->
												<svg class="h-4 w-4 text-gray-400 hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
												</svg>
											{/if}
										</button>
									</div>
								{:else if authType === 'basic'}
									<div class="space-y-2">
										<input
											type="text"
											bind:value={authUsername}
											placeholder="Username"
											class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
										/>
										<div class="relative">
											{#if showAuthPassword}
												<input
													type="text"
													bind:value={authPassword}
													placeholder="Password"
													class="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md text-sm"
												/>
											{:else}
												<input
													type="password"
													bind:value={authPassword}
													placeholder="Password"
													class="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md text-sm"
												/>
											{/if}
											<button
												type="button"
												class="absolute inset-y-0 right-0 pr-3 flex items-center"
												on:click={() => showAuthPassword = !showAuthPassword}
											>
												{#if showAuthPassword}
													<!-- Eye slash icon (hide) -->
													<svg class="h-4 w-4 text-gray-400 hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L8.464 8.464M14.12 14.12l1.415 1.415M14.12 14.12L9.878 9.878m4.242 4.242L8.464 8.464m5.656 5.656l1.415 1.415"></path>
													</svg>
												{:else}
													<!-- Eye icon (show) -->
													<svg class="h-4 w-4 text-gray-400 hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
													</svg>
												{/if}
											</button>
										</div>
									</div>
								{/if}
							</div>

							<!-- Body -->
							{#if requestMethod !== 'GET'}
								<div>
									<h3 class="text-sm font-medium text-gray-700 mb-2">Body</h3>
									<select bind:value={requestBodyType} class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm mb-3">
										<option value="json">JSON</option>
										<option value="form">Form Data</option>
										<option value="raw">Raw</option>
									</select>
									<textarea
										bind:value={requestBody}
										placeholder={requestBodyType === 'json' ? '{\n  "key": "value"\n}' : 'Request body...'}
										rows="8"
										class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm font-mono"
									></textarea>
								</div>
							{/if}
						</div>
					</div>

					<!-- Response Panel -->
					{#if response}
						<div class="w-1/2 bg-gray-50 border-l border-gray-200 flex flex-col">
							<div class="border-b border-gray-200 p-4">
								<h3 class="text-lg font-medium text-gray-900">Response</h3>
								<div class="flex items-center space-x-4 mt-2 text-sm">
									<span class="px-2 py-1 rounded text-xs font-mono {response.status >= 200 && response.status < 300 ? 'bg-green-100 text-green-800' : response.status >= 400 ? 'bg-red-100 text-red-800' : 'bg-yellow-100 text-yellow-800'}">{response.status}</span>
									<span class="text-gray-600">{response.time}ms</span>
								</div>
							</div>
							
							<div class="flex-1 overflow-y-auto p-4">
								<div class="space-y-4">
									<!-- Response Headers -->
									<div>
										<h4 class="text-sm font-medium text-gray-700 mb-2">Headers</h4>
										<div class="bg-white rounded border text-xs font-mono">
											{#each Object.entries(response.headers) as [key, value]}
												<div class="px-3 py-2 border-b border-gray-100 last:border-b-0">
													<span class="text-gray-600">{key}:</span> <span class="text-gray-900">{value}</span>
												</div>
											{/each}
										</div>
									</div>
									
									<!-- Response Body -->
									<div>
										<h4 class="text-sm font-medium text-gray-700 mb-2">Body</h4>
										<pre class="bg-white rounded border p-3 text-xs font-mono overflow-x-auto">{formatResponseBody(response.body)}</pre>
									</div>
								</div>
							</div>
						</div>
					{/if}
				</div>
			{:else}
				<!-- Welcome State -->
				<div class="flex-1 flex items-center justify-center">
					<div class="text-center">
						<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
						</svg>
						<h3 class="mt-2 text-sm font-medium text-gray-900">No collection selected</h3>
						<p class="mt-1 text-sm text-gray-500">Select a collection from the sidebar to start building requests.</p>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>

<!-- Create Collection Modal -->
{#if showNewCollectionModal}
	<div class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-lg shadow-xl max-w-md w-full">
			<div class="px-6 py-4 border-b border-gray-200">
				<h3 class="text-lg font-medium text-gray-900">Create New Collection</h3>
			</div>
			
			<div class="px-6 py-4">
				<div class="space-y-4">
					<div>
						<label for="collectionName" class="block text-sm font-medium text-gray-700">
							Collection Name
						</label>
						<input
							id="collectionName"
							type="text"
							bind:value={newCollectionName}
							placeholder="User API"
							class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>
					
					<div>
						<label for="collectionDescription" class="block text-sm font-medium text-gray-700">
							Description (optional)
						</label>
						<textarea
							id="collectionDescription"
							bind:value={newCollectionDescription}
							placeholder="API endpoints for user management..."
							rows="3"
							class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						></textarea>
					</div>
				</div>
			</div>
			
			<div class="px-6 py-4 bg-gray-50 flex justify-end space-x-3 rounded-b-lg">
				<button
					on:click={() => showNewCollectionModal = false}
					disabled={isCreatingCollection}
					class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
				>
					Cancel
				</button>
				<button
					on:click={createCollection}
					disabled={isCreatingCollection || !newCollectionName.trim()}
					class="px-4 py-2 text-sm font-medium text-white electric-gradient rounded-md hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 flex items-center"
				>
					{#if isCreatingCollection}
						<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Creating...
					{:else}
						Create Collection
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Error Toast -->
{#if error}
	<div class="fixed bottom-4 right-4 bg-red-50 border border-red-200 rounded-md p-4 max-w-sm z-50">
		<div class="flex">
			<div class="flex-shrink-0">
				<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
				</svg>
			</div>
			<div class="ml-3">
				<p class="text-sm text-red-800">{error}</p>
			</div>
			<div class="ml-auto pl-3">
				<button
					on:click={() => error = ''}
					class="text-red-400 hover:text-red-600"
				>
					<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.electric-gradient {
		background: linear-gradient(135deg, var(--electric-blue), var(--deep-blue));
	}
	
	.bg-clip-text {
		background-clip: text;
		-webkit-background-clip: text;
	}
</style>