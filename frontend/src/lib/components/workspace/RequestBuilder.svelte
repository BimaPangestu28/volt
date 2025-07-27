<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

	export let selectedRequest: any = null;
	export let requestName = 'New Request';
	export let requestMethod = 'GET';
	export let requestUrl = '';
	export let requestHeaders = [{ key: '', value: '' }];
	export let requestBody = '';
	export let requestBodyType = 'raw';
	export let authType = 'none';
	export let authToken = '';
	export let authUsername = '';
	export let authPassword = '';
	export let showAuthToken = false;
	export let showAuthPassword = false;
	export let isExecuting = false;
	export let activeTab = 'Body';
	export let saveCompletedCounter = 0;

	// Track if request has unsaved changes
	let hasUnsavedChanges = false;
	let originalRequestData: any = null;

	// Watch for changes to determine if request is modified
	$: {
		if (selectedRequest && originalRequestData) {
			const currentData = {
				name: requestName,
				method: requestMethod,
				url: requestUrl,
				headers: requestHeaders,
				body: requestBody,
				bodyType: requestBodyType,
				authType: authType,
				authToken: authToken,
				authUsername: authUsername,
				authPassword: authPassword
			};
			
			hasUnsavedChanges = JSON.stringify(currentData) !== JSON.stringify(originalRequestData);
		} else if (!selectedRequest && (requestName !== 'New Request' || requestUrl !== '' || requestBody !== '')) {
			hasUnsavedChanges = true;
		} else {
			hasUnsavedChanges = false;
		}
	}

	// Reset unsaved changes when save is completed
	$: if (saveCompletedCounter > 0) {
		hasUnsavedChanges = false;
		// Update original data after save
		if (selectedRequest) {
			originalRequestData = {
				name: requestName,
				method: requestMethod,
				url: requestUrl,
				headers: requestHeaders,
				body: requestBody,
				bodyType: requestBodyType,
				authType: authType,
				authToken: authToken,
				authUsername: authUsername,
				authPassword: authPassword
			};
		}
	}

	// Store original data when request is loaded
	$: if (selectedRequest) {
		originalRequestData = {
			name: selectedRequest.name || requestName,
			method: selectedRequest.method || requestMethod,
			url: selectedRequest.url || requestUrl,
			headers: selectedRequest.headers ? Object.entries(selectedRequest.headers).map(([key, value]) => ({ key, value: String(value) })) : requestHeaders,
			body: selectedRequest.body?.content || requestBody,
			bodyType: selectedRequest.body?.type || requestBodyType,
			authType: selectedRequest.auth?.type || authType,
			authToken: selectedRequest.auth?.credentials?.token || authToken,
			authUsername: selectedRequest.auth?.credentials?.username || authUsername,
			authPassword: selectedRequest.auth?.credentials?.password || authPassword
		};
	}

	const dispatch = createEventDispatcher();

	function handleExecuteRequest() {
		dispatch('executeRequest');
	}

	function handleSaveRequest() {
		dispatch('saveRequest');
		// Reset unsaved changes after successful save
		hasUnsavedChanges = false;
		// Update original data to current state
		if (selectedRequest) {
			originalRequestData = {
				name: requestName,
				method: requestMethod,
				url: requestUrl,
				headers: requestHeaders,
				body: requestBody,
				bodyType: requestBodyType,
				authType: authType,
				authToken: authToken,
				authUsername: authUsername,
				authPassword: authPassword
			};
		}
	}

	function handleAddHeader() {
		requestHeaders = [...requestHeaders, { key: '', value: '' }];
	}

	function handleRemoveHeader(index: number) {
		requestHeaders = requestHeaders.filter((_, i) => i !== index);
	}
</script>

<div class="flex-1 flex flex-col bg-gray-900">
	<!-- Request Tab Bar -->
	<div class="bg-gray-800 border-b border-gray-700">
		<div class="flex items-center">
			<!-- Request Tab -->
			<div class="flex items-center px-3 py-1.5 border-b-2 border-blue-500 bg-gray-750 min-w-0">
				<div class="flex items-center space-x-2 min-w-0">
					<span class="text-[9px] font-mono px-1 py-0.5 rounded font-semibold uppercase {
						requestMethod === 'GET' ? 'bg-green-700 text-green-300' : 
						requestMethod === 'POST' ? 'bg-blue-700 text-blue-300' : 
						requestMethod === 'PUT' ? 'bg-yellow-700 text-yellow-300' : 
						requestMethod === 'DELETE' ? 'bg-red-700 text-red-300' :
						'bg-purple-700 text-purple-300'
					}">
						{requestMethod}
					</span>
					<div class="flex items-center space-x-1 min-w-0 flex-1">
						<input
							type="text"
							bind:value={requestName}
							placeholder="Untitled Request"
							class="bg-transparent text-white text-[11px] border-none outline-none placeholder-gray-400 min-w-0 flex-1 {hasUnsavedChanges ? 'italic' : 'font-normal'}"
						/>
						{#if hasUnsavedChanges}
							<span class="w-1.5 h-1.5 bg-orange-500 rounded-full flex-shrink-0" title="Unsaved changes"></span>
						{/if}
					</div>
					<button class="text-gray-400 hover:text-white p-0.5 rounded transition-colors">
						<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
						</svg>
					</button>
				</div>
			</div>
			
			<!-- Add Tab Button -->
			<button class="p-2 text-gray-400 hover:text-white hover:bg-gray-700 transition-colors border-b-2 border-transparent hover:border-gray-500">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
				</svg>
			</button>
		</div>
	</div>

	<!-- Request Builder -->
	<div class="flex-1 flex flex-col">
			<!-- URL Bar -->
			<div class="p-4 bg-gray-800 border-b border-gray-700">
				<div class="flex items-center space-x-3">
					<select bind:value={requestMethod} class="px-3 py-2.5 border border-gray-600 rounded text-xs font-semibold bg-gray-700 text-white outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 min-w-[80px] {
						requestMethod === 'GET' ? 'text-green-300' :
						requestMethod === 'POST' ? 'text-blue-300' :
						requestMethod === 'PUT' ? 'text-yellow-300' :
						requestMethod === 'DELETE' ? 'text-red-300' :
						'text-purple-300'
					}">
						<option value="GET" class="bg-gray-700">GET</option>
						<option value="POST" class="bg-gray-700">POST</option>
						<option value="PUT" class="bg-gray-700">PUT</option>
						<option value="DELETE" class="bg-gray-700">DELETE</option>
						<option value="PATCH" class="bg-gray-700">PATCH</option>
					</select>
					
					<input
						type="url"
						bind:value={requestUrl}
						placeholder="https://api.example.com/users"
						class="flex-1 px-4 py-2.5 border border-gray-600 rounded text-xs bg-gray-700 text-white placeholder-gray-400 outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-mono"
					/>
					
					<button
						on:click={handleExecuteRequest}
						disabled={isExecuting}
						class="electric-gradient text-white px-6 py-2.5 text-xs font-semibold rounded hover:opacity-90 transition-opacity disabled:opacity-50 flex items-center min-w-[60px] justify-center"
					>
						{#if isExecuting}
							<LoadingSpinner size="sm" color="white" />
						{:else}
							Send
						{/if}
					</button>
					
					<div class="flex items-center space-x-2">
						<button
							on:click={handleSaveRequest}
							class="px-4 py-2.5 text-xs font-medium text-gray-300 border border-gray-600 rounded hover:bg-gray-700 hover:text-white transition-colors"
						>
							Save
						</button>
						
						<button class="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded transition-colors" title="More options">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path>
							</svg>
						</button>
					</div>
				</div>
			</div>

			<!-- Request Tabs -->
			<div class="bg-gray-800 border-b border-gray-700">
				<nav class="flex" aria-label="Tabs">
					{#each ['Params', 'Authorization', 'Headers (8)', 'Body', 'Scripts', 'Settings'] as tab}
						<button 
							on:click={() => activeTab = tab.split(' ')[0]}
							class="px-3 py-2 text-xs font-normal border-b-2 transition-colors"
							class:border-blue-500={activeTab === tab.split(' ')[0]}
							class:text-blue-400={activeTab === tab.split(' ')[0]}
							class:border-transparent={activeTab !== tab.split(' ')[0]}
							class:text-gray-400={activeTab !== tab.split(' ')[0]}
							class:hover:text-gray-200={activeTab !== tab.split(' ')[0]}
							class:hover:border-gray-500={activeTab !== tab.split(' ')[0]}
						>
							{tab}
						</button>
					{/each}
				</nav>
			</div>

			<div class="flex-1 overflow-y-auto bg-gray-800">
				<div class="p-4">
				<!-- Params Tab -->
				{#if activeTab === 'Params'}
					<div class="space-y-3">
						<div class="flex items-center justify-between">
							<h3 class="text-xs font-medium text-gray-300">Query Params</h3>
							<div class="flex items-center space-x-2">
								<button class="text-[10px] text-gray-400 hover:text-gray-200 font-medium">⋯ Bulk Edit</button>
							</div>
						</div>
						
						<!-- Table Header -->
						<div class="grid grid-cols-12 gap-2 text-[10px] font-medium text-gray-400 border-b border-gray-700 pb-2">
							<div class="col-span-4">Key</div>
							<div class="col-span-4">Value</div>
							<div class="col-span-4">Description</div>
						</div>
						
						<!-- Table Rows -->
						<div class="space-y-1">
							<div class="grid grid-cols-12 gap-2 items-center py-1">
								<div class="col-span-4">
									<input type="text" placeholder="Key" class="w-full px-2 py-1 bg-gray-700 border border-gray-600 rounded text-[10px] text-white placeholder-gray-500">
								</div>
								<div class="col-span-4">
									<input type="text" placeholder="Value" class="w-full px-2 py-1 bg-gray-700 border border-gray-600 rounded text-[10px] text-white placeholder-gray-500">
								</div>
								<div class="col-span-4">
									<input type="text" placeholder="Description" class="w-full px-2 py-1 bg-gray-700 border border-gray-600 rounded text-[10px] text-white placeholder-gray-500">
								</div>
							</div>
							<div class="grid grid-cols-12 gap-2 items-center py-1">
								<div class="col-span-4">
									<input type="text" placeholder="Key" class="w-full px-2 py-1 bg-gray-700 border border-gray-600 rounded text-[10px] text-white placeholder-gray-500">
								</div>
								<div class="col-span-4">
									<input type="text" placeholder="Value" class="w-full px-2 py-1 bg-gray-700 border border-gray-600 rounded text-[10px] text-white placeholder-gray-500">
								</div>
								<div class="col-span-4">
									<input type="text" placeholder="Description" class="w-full px-2 py-1 bg-gray-700 border border-gray-600 rounded text-[10px] text-white placeholder-gray-500">
								</div>
							</div>
						</div>
					</div>
				
				<!-- Authorization Tab -->
				{:else if activeTab === 'Authorization'}
					<div class="space-y-4">
						<div class="flex items-center justify-between">
							<h3 class="text-sm font-medium text-white">Authorization</h3>
						</div>
						<div class="space-y-4">
							<div>
								<label class="block text-xs font-medium text-gray-300 mb-2">Authorization Type</label>
								<select bind:value={authType} class="w-full max-w-sm px-3 py-2 border border-gray-600 rounded-md text-sm bg-gray-700 text-white outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
									<option value="none" class="bg-gray-700 text-white">No Authorization</option>
									<option value="bearer" class="bg-gray-700 text-white">Bearer Token</option>
									<option value="basic" class="bg-gray-700 text-white">Basic Auth</option>
									<option value="apikey" class="bg-gray-700 text-white">API Key</option>
								</select>
							</div>
						
						{#if authType === 'bearer'}
							<div class="bg-gray-700 p-4 rounded-md border border-gray-600">
								<label class="block text-xs font-medium text-gray-300 mb-2">Bearer Token</label>
								<div class="relative">
									{#if showAuthToken}
										<input
											type="text"
											bind:value={authToken}
											placeholder="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
											class="w-full px-3 py-2 pr-10 border border-gray-600 rounded-md text-sm bg-gray-800 text-white placeholder-gray-400 outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-mono"
										/>
									{:else}
										<input
											type="password"
											bind:value={authToken}
											placeholder="Enter your bearer token..."
											class="w-full px-3 py-2 pr-10 border border-gray-600 rounded-md text-sm bg-gray-800 text-white placeholder-gray-400 outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
										/>
									{/if}
									<button
										type="button"
										class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-200 transition-colors"
										on:click={() => showAuthToken = !showAuthToken}
									>
										{#if showAuthToken}
											<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L8.464 8.464M14.12 14.12l1.415 1.415M14.12 14.12L9.878 9.878m4.242 4.242L8.464 8.464m5.656 5.656l1.415 1.415"></path>
											</svg>
										{:else}
											<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
											</svg>
										{/if}
									</button>
								</div>
								<p class="text-xs text-gray-500 mt-1">The bearer token will be included in the Authorization header</p>
							</div>
						{:else if authType === 'basic'}
							<div class="bg-gray-800 p-4 rounded-md space-y-3">
								<div>
									<label class="block text-xs font-medium text-gray-300 mb-2">Username</label>
									<input
										type="text"
										bind:value={authUsername}
										placeholder="Enter username"
										class="w-full px-3 py-2 border border-gray-600 rounded-md text-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
									/>
								</div>
								<div>
									<label class="block text-xs font-medium text-gray-300 mb-2">Password</label>
									<div class="relative">
										{#if showAuthPassword}
											<input
												type="text"
												bind:value={authPassword}
												placeholder="Enter password"
												class="w-full px-3 py-2 pr-10 border border-gray-600 rounded-md text-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
											/>
										{:else}
											<input
												type="password"
												bind:value={authPassword}
												placeholder="Enter password"
												class="w-full px-3 py-2 pr-10 border border-gray-600 rounded-md text-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
											/>
										{/if}
										<button
											type="button"
											class="absolute inset-y-0 right-0 pr-3 flex items-center"
											on:click={() => showAuthPassword = !showAuthPassword}
										>
											{#if showAuthPassword}
												<svg class="h-4 w-4 text-gray-500 hover:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L8.464 8.464M14.12 14.12l1.415 1.415M14.12 14.12L9.878 9.878m4.242 4.242L8.464 8.464m5.656 5.656l1.415 1.415"></path>
												</svg>
											{:else}
												<svg class="h-4 w-4 text-gray-500 hover:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
												</svg>
											{/if}
										</button>
									</div>
								</div>
							</div>
						{:else if authType === 'none'}
							<div class="bg-gray-700 rounded-md border border-gray-600">
								<div class="text-center py-8 text-gray-400">
									<svg class="mx-auto h-8 w-8 text-gray-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
									</svg>
									<p class="text-sm font-medium mb-1">No authorization required</p>
									<p class="text-xs text-gray-500">This request will be sent without any authorization</p>
								</div>
							</div>
						{/if}
						</div>
					</div>
				
				<!-- Headers Tab -->
				{:else if activeTab === 'Headers'}
					<div class="space-y-4">
						<div class="flex items-center justify-between">
							<h3 class="text-sm font-medium text-white">Headers</h3>
							<button
								on:click={handleAddHeader}
								class="text-xs text-blue-400 hover:text-blue-300 flex items-center font-medium"
							>
								<svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
								</svg>
								Add Header
							</button>
						</div>
						<div class="space-y-2">
							{#each requestHeaders as header, index}
								<div class="flex space-x-2 items-center bg-gray-800 p-3 rounded-md">
									<input
										type="text"
										bind:value={header.key}
										placeholder="Header name"
										class="flex-1 px-3 py-2 border border-gray-600 rounded text-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
									/>
									<span class="text-gray-500">:</span>
									<input
										type="text"
										bind:value={header.value}
										placeholder="Header value"
										class="flex-1 px-3 py-2 border border-gray-600 rounded text-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
									/>
									<button
										on:click={() => handleRemoveHeader(index)}
										class="p-2 text-gray-500 hover:text-red-400 rounded"
									>
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
										</svg>
									</button>
								</div>
							{/each}
							{#if requestHeaders.length === 0}
								<div class="text-center py-8 text-gray-400">
									<svg class="mx-auto h-8 w-8 text-gray-500 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"></path>
									</svg>
									<p class="text-sm">No headers added</p>
									<p class="text-xs text-gray-500 mt-1">Add custom headers to your request</p>
								</div>
							{/if}
						</div>
					</div>
				
				<!-- Body Tab -->
				{:else if activeTab === 'Body'}
					<div class="space-y-4">
						{#if requestMethod === 'GET'}
							<div class="bg-gray-700 rounded-md border border-gray-600">
								<div class="text-center py-8 text-gray-400">
									<svg class="mx-auto h-8 w-8 text-gray-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
									</svg>
									<p class="text-sm font-medium mb-1">This request does not have a body</p>
									<p class="text-xs text-gray-500">GET requests typically don't include a body</p>
								</div>
							</div>
						{:else}
							<!-- Body Type Selection - Postman style -->
							<div class="flex items-center space-x-4 mb-4">
								<label class="flex items-center space-x-1.5 cursor-pointer">
									<input type="radio" bind:group={requestBodyType} value="none" class="w-3.5 h-3.5 text-blue-500 bg-gray-800 border-gray-600">
									<span class="text-xs text-gray-300">none</span>
								</label>
								<label class="flex items-center space-x-1.5 cursor-pointer">
									<input type="radio" bind:group={requestBodyType} value="form-data" class="w-3.5 h-3.5 text-blue-500 bg-gray-800 border-gray-600">
									<span class="text-xs text-gray-300">form-data</span>
								</label>
								<label class="flex items-center space-x-1.5 cursor-pointer">
									<input type="radio" bind:group={requestBodyType} value="x-www-form-urlencoded" class="w-3.5 h-3.5 text-blue-500 bg-gray-800 border-gray-600">
									<span class="text-xs text-gray-300">x-www-form-urlencoded</span>
								</label>
								<label class="flex items-center space-x-1.5 cursor-pointer">
									<input type="radio" bind:group={requestBodyType} value="raw" class="w-3.5 h-3.5 text-blue-500 bg-gray-800 border-gray-600">
									<span class="text-xs text-gray-300">raw</span>
								</label>
								<label class="flex items-center space-x-1.5 cursor-pointer">
									<input type="radio" bind:group={requestBodyType} value="binary" class="w-3.5 h-3.5 text-blue-500 bg-gray-800 border-gray-600">
									<span class="text-xs text-gray-300">binary</span>
								</label>
								<label class="flex items-center space-x-1.5 cursor-pointer">
									<input type="radio" bind:group={requestBodyType} value="GraphQL" class="w-3.5 h-3.5 text-blue-500 bg-gray-800 border-gray-600">
									<span class="text-xs text-gray-300">GraphQL</span>
								</label>
								
								{#if requestBodyType === 'raw'}
									<select class="ml-4 px-2 py-1 border border-gray-600 rounded text-xs bg-gray-700 text-white">
										<option value="json">JSON</option>
										<option value="text">Text</option>
										<option value="javascript">JavaScript</option>
										<option value="html">HTML</option>
										<option value="xml">XML</option>
									</select>
								{/if}
							</div>
							
							{#if requestBodyType === 'none'}
								<div class="bg-gray-700/30 rounded border border-gray-600/50 mt-4">
									<div class="text-center py-12 text-gray-400">
										<p class="text-xs text-gray-500">This request does not have a body</p>
									</div>
								</div>
							{:else if requestBodyType === 'raw' || requestBodyType === 'json'}
								<div class="mt-4">
									<textarea
										bind:value={requestBody}
										placeholder={requestBodyType === 'json' ? '{\n  "campaignId": "[Test] Survey Tes #3",\n  "campaignTitle": "[Test] Survey IPSOS #3",\n  "customerId": "123456",\n  "phone": "628237186356",\n  "templateId": "test_surveyipsos",\n  "arguments": [\n    "Bima"\n  ],\n  "media": {\n    "header": {\n      "type": "image",\n      "url": "https://d1sin90wg7b4g2.cloudfront.net/survey-image-2025-04-22-at-14-11-50.jpg"\n    }\n  }\n}' : 'Enter raw text content...'}
										rows="18"
										class="w-full px-3 py-3 border border-gray-600 rounded text-xs font-mono bg-gray-700 text-white placeholder-gray-500 resize-none leading-relaxed"
									></textarea>
								</div>
							{:else if requestBodyType === 'form-data' || requestBodyType === 'x-www-form-urlencoded'}
								<div class="mt-4">
									<div class="bg-gray-700/30 rounded border border-gray-600/50 p-6">
										<div class="text-center py-8 text-gray-400">
											<p class="text-xs text-gray-500 mb-1">No form data added</p>
											<p class="text-xs text-gray-600">Add key-value pairs for your form data</p>
										</div>
									</div>
								</div>
							{:else}
								<div class="bg-gray-700/30 rounded border border-gray-600/50 mt-4">
									<div class="text-center py-12 text-gray-400">
										<p class="text-xs text-gray-500 mb-1">{requestBodyType} body type</p>
										<p class="text-xs text-gray-600">Configure your {requestBodyType} body content</p>
									</div>
								</div>
							{/if}
						{/if}
					</div>
				
				<!-- Scripts Tab -->
				{:else if activeTab === 'Scripts'}
					<div class="space-y-4">
						<div class="flex items-center justify-between">
							<h3 class="text-sm font-medium text-white">Scripts</h3>
						</div>
						<div class="bg-gray-700 rounded-md border border-gray-600">
							<div class="text-center py-8 text-gray-400">
								<svg class="mx-auto h-8 w-8 text-gray-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path>
								</svg>
								<p class="text-sm font-medium mb-1">No scripts configured</p>
								<p class="text-xs text-gray-500">Add pre-request scripts to customize your request behavior</p>
							</div>
						</div>
					</div>
				
				<!-- Tests Tab -->
				{:else if activeTab === 'Tests'}
					<div class="space-y-4">
						<div class="flex items-center justify-between">
							<h3 class="text-sm font-medium text-white">Tests</h3>
						</div>
						<div class="bg-gray-700 rounded-md border border-gray-600">
							<div class="text-center py-8 text-gray-400">
								<svg class="mx-auto h-8 w-8 text-gray-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
								</svg>
								<p class="text-sm font-medium mb-1">No tests configured</p>
								<p class="text-xs text-gray-500">Add JavaScript tests to validate your API responses</p>
							</div>
						</div>
					</div>
				
				<!-- Settings Tab -->
				{:else if activeTab === 'Settings'}
					<div class="space-y-4">
						<div class="flex items-center justify-between">
							<h3 class="text-sm font-medium text-white">Settings</h3>
						</div>
						<div class="bg-gray-700 rounded-md border border-gray-600">
							<div class="text-center py-8 text-gray-400">
								<svg class="mx-auto h-8 w-8 text-gray-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
								</svg>
								<p class="text-sm font-medium mb-1">Request Settings</p>
								<p class="text-xs text-gray-500">Configure request timeout, redirects, and other settings</p>
							</div>
						</div>
					</div>
				{/if}
				</div>
		</div>
	</div>
</div>

<style>
	/* Remove all outline styles globally */
	:global(*) {
		outline: none !important;
		box-shadow: none !important;
	}
	
	:global(*:focus) {
		outline: none !important;
		box-shadow: none !important;
	}
	
	:global(button),
	:global(input),
	:global(textarea),
	:global(select),
	:global(button:focus),
	:global(input:focus),
	:global(textarea:focus),
	:global(select:focus) {
		outline: none !important;
		box-shadow: none !important;
		border-outline: none !important;
	}
	
	:global(input[type="text"]),
	:global(input[type="password"]),
	:global(input[type="url"]),
	:global(input[type="email"]) {
		outline: none !important;
		box-shadow: none !important;
		-webkit-appearance: none !important;
		-moz-appearance: none !important;
		appearance: none !important;
	}
	
	:global(select) {
		outline: none !important;
		box-shadow: none !important;
		-webkit-appearance: none !important;
		-moz-appearance: none !important;
		appearance: none !important;
	}
</style>