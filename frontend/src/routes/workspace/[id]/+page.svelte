<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { workspaceStore } from '$lib/stores/workspace';
	import { apiClient } from '$lib/api/client';
	import { toastStore } from '$lib/stores/toast';
	
	// Components
	import WorkspaceHeader from '$lib/components/workspace/WorkspaceHeader.svelte';
	import CollectionsSidebar from '$lib/components/workspace/CollectionsSidebar.svelte';
	import EnvironmentsSidebar from '$lib/components/workspace/EnvironmentsSidebar.svelte';
	import EnvironmentManager from '$lib/components/workspace/EnvironmentManager.svelte';
	import RequestBuilder from '$lib/components/workspace/RequestBuilder.svelte';
	import ResponsePanel from '$lib/components/workspace/ResponsePanel.svelte';
	import CreateCollectionModal from '$lib/components/workspace/CreateCollectionModal.svelte';
	import CreateEnvironmentModal from '$lib/components/workspace/CreateEnvironmentModal.svelte';

	let workspaceId = '';
	let isAuthenticated = false;
	let currentWorkspace: any = null;
	let collections: any[] = [];
	let selectedCollection: any = null;
	let selectedRequest: any = null;
	let requests: any[] = [];
	let isLoading = true;
	let error = '';
	
	// Search state
	let searchResults: any[] = [];
	let isSearching = false;
	
	// Request builder state
	let requestName = 'New Request';
	let requestMethod = 'GET';
	let requestUrl = '';
	let requestHeaders = [{ key: '', value: '' }];
	let requestBody = '';
	let requestBodyType = 'raw';
	let authType = 'none';
	let authToken = '';
	let authUsername = '';
	let authPassword = '';
	let showAuthToken = false;
	let showAuthPassword = false;
	
	// Response state
	let response: any = null;
	let isExecuting = false;
	let activeResponseTab = 'body';
	let activeTab = 'Body';

	// UI state
	let showNewCollectionModal = false;
	let newCollectionName = '';
	let newCollectionDescription = '';
	let isCreatingCollection = false;

	// Environment state
	let environments: any[] = [];
	let selectedEnvironment: any = null;
	let isLoadingEnvironments = false;
	let showNewEnvironmentModal = false;
	let newEnvironmentName = '';
	let isCreatingEnvironment = false;
	let isSavingEnvironment = false;
	
	// Sidebar tab state
	let activeSidebarTab: 'collections' | 'environments' = 'collections';

	// Resizable divider state
	let responseHeight = 320; // Default height in pixels
	let isDragging = false;
	let startY = 0;
	let startHeight = 0;
	
	// Save state tracking
	let saveCompletedCounter = 0;
	
	// Tab system for multiple requests
	interface RequestTab {
		id: string;
		name: string;
		method: string;
		url: string;
		headers: { key: string; value: string }[];
		body: string;
		bodyType: string;
		authType: string;
		authToken: string;
		authUsername: string;
		authPassword: string;
		isUnsaved: boolean;
		requestId?: string; // For saved requests
		response?: any; // Store response for each tab
	}
	
	let requestTabs: RequestTab[] = [];
	let activeTabId: string = '';
	let tabCounter = 1;
	
	// Clean up duplicate tabs
	function cleanupTabs() {
		// Remove duplicate tabs for same request
		const seen = new Set();
		requestTabs = requestTabs.filter(tab => {
			if (tab.requestId) {
				if (seen.has(tab.requestId)) {
					return false; // Remove duplicate
				}
				seen.add(tab.requestId);
			}
			return true;
		});
		
		// If active tab was removed, switch to first available tab
		if (activeTabId && !requestTabs.find(t => t.id === activeTabId)) {
			if (requestTabs.length > 0) {
				switchToTab(requestTabs[0].id);
			} else {
				activeTabId = '';
			}
		}
	}

	// Subscribe to stores
	authStore.subscribe(state => {
		isAuthenticated = state.isAuthenticated;
	});

	workspaceStore.subscribe(state => {
		currentWorkspace = state.currentWorkspace;
		collections = state.collections || [];
		isLoading = state.isLoading;
		error = state.error || '';
		console.log('Workspace store updated:', { collections: collections.length, isLoading, error });
	});

	// Get workspace ID from URL
	page.subscribe(({ params }) => {
		workspaceId = params.id || '';
		console.log('Workspace ID set:', workspaceId);
	});

	onMount(async () => {
		console.log('onMount called, isAuthenticated:', isAuthenticated, 'workspaceId:', workspaceId);
		
		// Add window resize listener
		window.addEventListener('resize', handleWindowResize);
		
		let attempts = 0;
		const checkAndLoad = async () => {
			attempts++;
			console.log(`Attempt ${attempts}: auth=${isAuthenticated}, workspaceId=${workspaceId}`);
			
			if (!isAuthenticated && attempts > 10) {
				goto('/login');
				return;
			}
			
			if (isAuthenticated && workspaceId) {
				console.log('Both conditions met, loading data...');
				await loadWorkspace();
				await loadCollections();
				return;
			}
			
			setTimeout(checkAndLoad, 100);
		};
		
		checkAndLoad();
		
		// Cleanup function
		return () => {
			window.removeEventListener('resize', handleWindowResize);
		};
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
			
			const cachedCollections = workspaceStore.getCachedCollections(workspaceId);
			
			if (cachedCollections) {
				console.log('Using cached collections:', cachedCollections);
				workspaceStore.setCollections(cachedCollections);
				
				const isStale = !workspaceStore.isCacheValid(workspaceId);
				if (isStale) {
					console.log('Cache is stale, refreshing in background...');
					workspaceStore.setRefreshing(workspaceId, true);
					
					try {
						const freshData = await apiClient.getCollections(workspaceId);
						console.log('Background refresh completed:', freshData);
						workspaceStore.setCachedCollections(workspaceId, freshData);
					} catch (refreshErr: any) {
						console.error('Background refresh failed:', refreshErr);
					} finally {
						workspaceStore.setRefreshing(workspaceId, false);
					}
				}
			} else {
				workspaceStore.setLoading(true);
				const data = await apiClient.getCollections(workspaceId);
				console.log('Collections received:', data);
				workspaceStore.setCachedCollections(workspaceId, data);
				workspaceStore.setLoading(false);
			}
		} catch (err: any) {
			console.error('Error loading collections:', err);
			workspaceStore.setError(err.response?.data?.error || 'Failed to load collections');
			workspaceStore.setLoading(false);
		}
	}

	async function createCollection() {
		if (!newCollectionName.trim()) return;

		isCreatingCollection = true;
		try {
			const collection = await apiClient.createCollection(workspaceId, newCollectionName, newCollectionDescription);
			
			workspaceStore.setLoading(true);
			const freshData = await apiClient.getCollections(workspaceId);
			workspaceStore.setCachedCollections(workspaceId, freshData);
			workspaceStore.setLoading(false);
			
			toastStore.success(`Collection "${newCollectionName}" created successfully!`);
			
			newCollectionName = '';
			newCollectionDescription = '';
			showNewCollectionModal = false;
		} catch (err: any) {
			const errorMessage = err.response?.data?.error || 'Failed to create collection';
			error = errorMessage;
			toastStore.error(errorMessage);
			workspaceStore.setLoading(false);
		} finally {
			isCreatingCollection = false;
		}
	}

	async function selectCollection(event: any) {
		const collection = event.detail;
		
		// Toggle collection - if already selected, close it
		if (selectedCollection?.id === collection.id) {
			selectedCollection = null;
			selectedRequest = null;
			response = null;
			requests = [];
			return;
		}
		
		// Open new collection
		selectedCollection = collection;
		selectedRequest = null;
		response = null;
		
		try {
			const data = await apiClient.getRequests(collection.id);
			requests = data;
			
			// Clean up any duplicate tabs
			cleanupTabs();
		} catch (err: any) {
			error = err.response?.data?.error || 'Failed to load requests';
		}
	}

	function selectRequest(event: any) {
		const request = event.detail;
		createTabFromRequest(request);
	}

	async function executeRequest() {
		if (!requestUrl.trim()) {
			toastStore.error('Please enter a request URL');
			return;
		}

		isExecuting = true;
		response = null;
		
		try {
			const headers: Record<string, string> = {};
			requestHeaders.forEach(({ key, value }) => {
				if (key.trim() && value.trim()) {
					headers[key] = value;
				}
			});

			const auth: any = {
				type: authType,
				credentials: {} as Record<string, string>
			};

			if (authType === 'bearer' && authToken) {
				auth.credentials.token = authToken;
				headers['Authorization'] = `Bearer ${authToken}`;
			} else if (authType === 'basic' && authUsername) {
				auth.credentials.username = authUsername;
				auth.credentials.password = authPassword;
				const basicAuth = btoa(`${authUsername}:${authPassword}`);
				headers['Authorization'] = `Basic ${basicAuth}`;
			}

			const requestData = {
				name: requestName || 'Untitled Request',
				method: requestMethod,
				url: requestUrl,
				headers,
				body: requestMethod !== 'GET' ? {
					type: requestBodyType,
					content: requestBody
				} : undefined,
				auth
			};

			const startTime = Date.now();
			const result = await apiClient.executeRawRequest(requestData);
			const endTime = Date.now();
			
			response = {
				...result,
				time: endTime - startTime
			};
			
			// Update response in current tab
			if (activeTabId) {
				const tabIndex = requestTabs.findIndex(t => t.id === activeTabId);
				if (tabIndex >= 0) {
					requestTabs[tabIndex].response = response;
					requestTabs = [...requestTabs];
				}
			}
			
			toastStore.success('Request executed successfully!');
		} catch (err: any) {
			const errorMessage = err.response?.data?.error || err.message || 'Failed to execute request';
			error = errorMessage;
			toastStore.error(errorMessage);
			
			if (err.response) {
				const endTime = Date.now();
				response = {
					status: err.response.status,
					statusText: err.response.statusText,
					headers: err.response.headers || {},
					body: err.response.data || err.message,
					time: endTime - (response?.startTime || Date.now())
				};
			}
		} finally {
			isExecuting = false;
		}
	}

	async function saveRequest() {
		if (!selectedCollection) return;

		const headers: Record<string, string> = {};
		requestHeaders.forEach(({ key, value }) => {
			if (key.trim() && value.trim()) {
				headers[key] = value;
			}
		});

		const auth: any = {
			type: authType,
			credentials: {} as Record<string, string>
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
				await apiClient.updateRequest(selectedRequest.id, requestData);
				const updatedRequests = requests.map(r => 
					r.id === selectedRequest.id ? { ...r, ...requestData } : r
				);
				requests = updatedRequests;
				selectedRequest = { ...selectedRequest, ...requestData };
			} else {
				const newRequest = await apiClient.createRequest(selectedCollection.id, requestData);
				requests = [...requests, newRequest];
				selectedRequest = newRequest;
			}
			
			// Update tab to mark as saved and update requestId if it's a new request
			if (activeTabId) {
				const tabIndex = requestTabs.findIndex(t => t.id === activeTabId);
				if (tabIndex >= 0) {
					const updatedTab = {
						...requestTabs[tabIndex],
						isUnsaved: false,
						name: requestName
					};
					
					// If this was a new request that got saved, update the requestId
					if (!requestTabs[tabIndex].requestId && selectedRequest) {
						updatedTab.requestId = selectedRequest.id;
					}
					
					requestTabs[tabIndex] = updatedTab;
					requestTabs = [...requestTabs];
				}
			}
			
			// Increment counter to signal save completion
			saveCompletedCounter++;
		} catch (err: any) {
			error = err.response?.data?.error || 'Failed to save request';
			throw err; // Re-throw for keyboard shortcut error handling
		}
	}

	// Tab management functions
	function createNewTab() {
		const newTab: RequestTab = {
			id: `tab-${Date.now()}-${tabCounter++}`,
			name: `New Request ${tabCounter - 1}`,
			method: 'GET',
			url: '',
			headers: [{ key: '', value: '' }],
			body: '',
			bodyType: 'raw',
			authType: 'none',
			authToken: '',
			authUsername: '',
			authPassword: '',
			isUnsaved: false,
			response: null
		};
		
		requestTabs = [...requestTabs, newTab];
		switchToTab(newTab.id);
		return newTab;
	}
	
	function createTabFromRequest(request: any) {
		// If tab for this request already exists, switch to it
		const existingTab = requestTabs.find(tab => tab.requestId === request.id);
		if (existingTab) {
			switchToTab(existingTab.id);
			return existingTab;
		}
		
		const newTab: RequestTab = {
			id: `tab-${Date.now()}-${request.id}`,
			name: request.name,
			method: request.method,
			url: request.url,
			headers: Object.entries(request.headers || {}).map(([key, value]) => ({ key, value: String(value) })),
			body: request.body?.content || '',
			bodyType: request.body?.type || 'raw',
			authType: request.auth?.type || 'none',
			authToken: request.auth?.credentials?.token || '',
			authUsername: request.auth?.credentials?.username || '',
			authPassword: request.auth?.credentials?.password || '',
			isUnsaved: false,
			requestId: request.id,
			response: null
		};
		
		requestTabs = [...requestTabs, newTab];
		switchToTab(newTab.id);
		return newTab;
	}
	
	function switchToTab(tabId: string) {
		activeTabId = tabId;
		const tab = requestTabs.find(t => t.id === tabId);
		if (tab) {
			// Update current form state
			requestName = tab.name;
			requestMethod = tab.method;
			requestUrl = tab.url;
			requestHeaders = tab.headers.length > 0 ? tab.headers : [{ key: '', value: '' }];
			requestBody = tab.body;
			requestBodyType = tab.bodyType;
			authType = tab.authType;
			authToken = tab.authToken;
			authUsername = tab.authUsername;
			authPassword = tab.authPassword;
			response = tab.response;
			
			// Update selected request if it's a saved request
			if (tab.requestId) {
				selectedRequest = requests.find(r => r.id === tab.requestId) || null;
			} else {
				selectedRequest = null;
			}
		}
	}
	
	function closeTab(tabId: string) {
		const tabIndex = requestTabs.findIndex(t => t.id === tabId);
		if (tabIndex === -1) return;
		
		const tab = requestTabs[tabIndex];
		
		// If tab has unsaved changes, show confirmation (simplified for now)
		if (tab.isUnsaved) {
			const confirm = window.confirm(`"${tab.name}" has unsaved changes. Close anyway?`);
			if (!confirm) return;
		}
		
		// Remove tab
		requestTabs = requestTabs.filter(t => t.id !== tabId);
		
		// Switch to another tab or create new one
		if (activeTabId === tabId) {
			if (requestTabs.length > 0) {
				// Switch to previous tab or next available tab
				const newActiveIndex = Math.max(0, tabIndex - 1);
				switchToTab(requestTabs[newActiveIndex].id);
			} else {
				// No tabs left, create a new one
				createNewTab();
			}
		}
	}
	
	function updateCurrentTab() {
		if (activeTabId) {
			const tabIndex = requestTabs.findIndex(t => t.id === activeTabId);
			if (tabIndex >= 0) {
				const currentTab = requestTabs[tabIndex];
				const updatedTab = {
					...currentTab,
					name: requestName,
					method: requestMethod,
					url: requestUrl,
					headers: requestHeaders,
					body: requestBody,
					bodyType: requestBodyType,
					authType: authType,
					authToken: authToken,
					authUsername: authUsername,
					authPassword: authPassword,
					isUnsaved: true // Mark as unsaved when modified
				};
				
				requestTabs[tabIndex] = updatedTab;
				requestTabs = [...requestTabs]; // Trigger reactivity
			}
		}
	}
	
	// Watch for changes and update current tab
	let lastFormState = '';
	$: {
		const currentFormState = JSON.stringify({
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
		});
		
		if (activeTabId && lastFormState && currentFormState !== lastFormState) {
			updateCurrentTab();
		}
		lastFormState = currentFormState;
	}

	function newRequest() {
		// Always create a new tab when user explicitly clicks "Add request"
		// This ensures user gets a fresh tab for new request
		createNewTab();
	}

	// Resizable divider functions
	function handleMouseDown(event: MouseEvent) {
		isDragging = true;
		startY = event.clientY;
		startHeight = responseHeight;
		
		document.addEventListener('mousemove', handleMouseMove);
		document.addEventListener('mouseup', handleMouseUp);
		document.body.style.cursor = 'ns-resize';
		document.body.style.userSelect = 'none';
	}

	function handleMouseMove(event: MouseEvent) {
		if (!isDragging) return;
		
		const deltaY = startY - event.clientY; // Inverted because we want dragging up to increase height
		const availableHeight = window.innerHeight - 120; // Account for header and some margin
		const maxResponseHeight = Math.min(600, availableHeight * 0.7); // Max 70% of available height or 600px
		const minResponseHeight = Math.max(200, availableHeight * 0.2); // Min 20% of available height or 200px
		const minRequestHeight = 300; // Minimum height for request builder
		
		// Calculate new height ensuring request builder gets minimum space
		const maxAllowedResponseHeight = availableHeight - minRequestHeight;
		const newHeight = Math.max(
			minResponseHeight, 
			Math.min(maxResponseHeight, Math.min(maxAllowedResponseHeight, startHeight + deltaY))
		);
		
		responseHeight = newHeight;
	}

	function handleMouseUp() {
		isDragging = false;
		document.removeEventListener('mousemove', handleMouseMove);
		document.removeEventListener('mouseup', handleMouseUp);
		document.body.style.cursor = '';
		document.body.style.userSelect = '';
	}

	// Handle window resize to adjust response height limits
	function handleWindowResize() {
		const availableHeight = window.innerHeight - 120;
		const maxResponseHeight = Math.min(600, availableHeight * 0.7);
		const minResponseHeight = Math.max(200, availableHeight * 0.2);
		const minRequestHeight = 300;
		const maxAllowedResponseHeight = availableHeight - minRequestHeight;
		
		// Adjust responseHeight if it exceeds new limits
		responseHeight = Math.max(
			minResponseHeight, 
			Math.min(maxResponseHeight, Math.min(maxAllowedResponseHeight, responseHeight))
		);
	}

	// Handle keyboard shortcuts
	async function handleKeyDown(event: KeyboardEvent) {
		// Check for Ctrl+S (Windows/Linux) or Cmd+S (Mac)
		if ((event.ctrlKey || event.metaKey) && event.key === 's') {
			event.preventDefault(); // Prevent browser's default save dialog
			
			// Only save if we have a selected collection
			if (selectedCollection) {
				try {
					await saveRequest();
					toastStore.success('Request saved!');
				} catch (err) {
					toastStore.error('Failed to save request');
				}
			} else {
				toastStore.error('Please select a collection first');
			}
		}
	}

	// Environment functions
	async function loadEnvironments() {
		if (!workspaceId) return;
		
		isLoadingEnvironments = true;
		try {
			const data = await apiClient.getEnvironments(workspaceId);
			environments = data;
		} catch (err: any) {
			console.error('Error loading environments:', err);
			toastStore.error('Failed to load environments');
		} finally {
			isLoadingEnvironments = false;
		}
	}

	async function createEnvironment() {
		if (!newEnvironmentName.trim()) return;

		isCreatingEnvironment = true;
		try {
			const environment = await apiClient.createEnvironment(workspaceId, newEnvironmentName);
			environments = [...environments, environment];
			toastStore.success(`Environment "${newEnvironmentName}" created successfully!`);
			
			newEnvironmentName = '';
			showNewEnvironmentModal = false;
		} catch (err: any) {
			const errorMessage = err.response?.data?.error || 'Failed to create environment';
			toastStore.error(errorMessage);
		} finally {
			isCreatingEnvironment = false;
		}
	}

	function selectEnvironment(event: any) {
		selectedEnvironment = event.detail;
	}

	async function saveEnvironment(event: any) {
		const environment = event.detail;
		
		isSavingEnvironment = true;
		try {
			const updatedEnvironment = await apiClient.updateEnvironment(
				environment.id,
				environment.name,
				environment.variables
			);
			
			// Update environment in list
			const index = environments.findIndex(env => env.id === environment.id);
			if (index >= 0) {
				environments[index] = updatedEnvironment;
				environments = [...environments];
			}
			
			// Update selected environment
			if (selectedEnvironment?.id === environment.id) {
				selectedEnvironment = updatedEnvironment;
			}
			
			toastStore.success('Environment saved successfully!');
		} catch (err: any) {
			const errorMessage = err.response?.data?.error || 'Failed to save environment';
			toastStore.error(errorMessage);
		} finally {
			isSavingEnvironment = false;
		}
	}

	async function deleteEnvironment(event: any) {
		const environment = event.detail;
		
		const confirmed = confirm(`Delete environment "${environment.name}"? This action cannot be undone.`);
		if (!confirmed) return;
		
		try {
			await apiClient.deleteEnvironment(environment.id);
			environments = environments.filter(env => env.id !== environment.id);
			
			if (selectedEnvironment?.id === environment.id) {
				selectedEnvironment = null;
			}
			
			toastStore.success(`Environment "${environment.name}" deleted successfully!`);
		} catch (err: any) {
			const errorMessage = err.response?.data?.error || 'Failed to delete environment';
			toastStore.error(errorMessage);
		}
	}

	function handleSwitchTab(event: any) {
		activeSidebarTab = event.detail;
		
		// Load environments when switching to environments tab
		if (activeSidebarTab === 'environments' && environments.length === 0) {
			loadEnvironments();
		}
	}

	// Search handler
	async function handleSearch(event: any) {
		const { query, workspaceId } = event.detail;
		
		if (!query || !workspaceId) {
			searchResults = [];
			return;
		}

		isSearching = true;
		try {
			const results = await apiClient.searchWorkspace(workspaceId, query);
			searchResults = results;
		} catch (err: any) {
			console.error('Search error:', err);
			searchResults = [];
			toastStore.error('Search failed');
		} finally {
			isSearching = false;
		}
	}

	function handleSearchResultClick(event: any) {
		const result = event.detail;
		
		if (result.type === 'collection') {
			// Select the collection
			const collection = collections.find(c => c.id === result.id);
			if (collection) {
				selectCollection({ detail: collection });
			}
		} else if (result.type === 'request') {
			// Find and select the collection first, then the request
			const collection = collections.find(c => 
				c.requests && c.requests.includes(result.id)
			);
			if (collection) {
				selectCollection({ detail: collection });
				// Wait a bit for requests to load, then select the request
				setTimeout(() => {
					const request = requests.find(r => r.id === result.id);
					if (request) {
						selectRequest({ detail: request });
					}
				}, 100);
			}
		}
	}
</script>

<svelte:head>
	<title>{currentWorkspace?.name || 'Workspace'} - Volt</title>
</svelte:head>

<!-- Global keyboard shortcuts -->
<svelte:window on:keydown={handleKeyDown} />

<div class="h-screen flex flex-col bg-gray-900">
	<!-- Header -->
	<WorkspaceHeader 
		{currentWorkspace}
		{searchResults}
		on:createCollection={() => showNewCollectionModal = true}
		on:linkGenerated={(e) => {
			console.log('Invitation link generated:', e.detail);
			toastStore.success('Invitation link generated successfully!');
		}}
		on:search={handleSearch}
		on:searchResultClick={handleSearchResultClick}
	/>

	<div class="flex-1 flex overflow-hidden">
		<!-- Sidebar -->
		{#if activeSidebarTab === 'collections'}
			<CollectionsSidebar
				{collections}
				{selectedCollection}
				{selectedRequest}
				{requests}
				{isLoading}
				activeTab={activeSidebarTab}
				on:selectCollection={selectCollection}
				on:selectRequest={selectRequest}
				on:newRequest={newRequest}
				on:createCollection={() => showNewCollectionModal = true}
				on:refreshCollections={loadCollections}
				on:switchTab={handleSwitchTab}
			/>
		{:else if activeSidebarTab === 'environments'}
			<div class="flex bg-gray-800 border-r border-gray-700">
				<!-- Left Navigation (same as collections) -->
				<div class="w-20 bg-gray-900 border-r border-gray-700 flex flex-col">
					<div class="py-4 space-y-3">
						<!-- Collections Icon -->
						<div class="flex flex-col items-center px-1">
							<button 
								on:click={() => handleSwitchTab({ detail: 'collections' })}
								class="w-10 h-10 flex items-center justify-center rounded-md mb-1.5 transition-colors text-gray-400 hover:bg-gray-700 hover:text-white"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
								</svg>
							</button>
							<span class="text-[9px] font-medium text-center leading-tight px-0.5 text-gray-500">Collections</span>
						</div>
						
						<!-- Environments Icon -->
						<div class="flex flex-col items-center px-1">
							<button 
								class="w-10 h-10 flex items-center justify-center rounded-md mb-1.5 transition-colors bg-blue-600 text-white hover:bg-blue-500"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 919-9"></path>
								</svg>
							</button>
							<span class="text-[9px] font-medium text-center leading-tight px-0.5 text-blue-400">Environments</span>
						</div>
					</div>
				</div>

				<!-- Environments Sidebar -->
				<EnvironmentsSidebar
					{environments}
					{selectedEnvironment}
					isLoading={isLoadingEnvironments}
					on:selectEnvironment={selectEnvironment}
					on:createEnvironment={() => showNewEnvironmentModal = true}
					on:deleteEnvironment={deleteEnvironment}
				/>
			</div>
		{/if}

		<!-- Main Content -->
		<div class="flex-1 flex flex-col bg-gray-900 overflow-hidden relative">
			{#if activeSidebarTab === 'environments'}
				<!-- Environment Manager -->
				<EnvironmentManager
					environment={selectedEnvironment}
					isLoading={isLoadingEnvironments}
					isSaving={isSavingEnvironment}
					on:saveEnvironment={saveEnvironment}
				/>
			{:else if selectedCollection}
				<!-- Request Builder -->
				<div 
					class="overflow-hidden bg-gray-900" 
					style="height: calc(100vh - 64px - {response || isExecuting ? responseHeight + 4 : 0}px); min-height: 300px;"
				>
					<RequestBuilder
						{selectedRequest}
						bind:requestName
						bind:requestMethod
						bind:requestUrl
						bind:requestHeaders
						bind:requestBody
						bind:requestBodyType
						bind:authType
						bind:authToken
						bind:authUsername
						bind:authPassword
						bind:showAuthToken
						bind:showAuthPassword
						bind:activeTab
						{isExecuting}
						{saveCompletedCounter}
						{requestTabs}
						{activeTabId}
						onSwitchTab={switchToTab}
						onCloseTab={closeTab}
						onCreateTab={createNewTab}
						on:executeRequest={executeRequest}
						on:saveRequest={saveRequest}
					/>
				</div>
				
				{#if response || isExecuting}
					<!-- Resizable Divider -->
					<div 
						class="h-1 bg-gray-700 hover:bg-blue-500 cursor-ns-resize transition-colors relative group flex-shrink-0"
						on:mousedown={handleMouseDown}
						role="separator"
						tabindex="0"
					>
						<!-- Drag Handle Visual -->
						<div class="absolute inset-x-0 top-1/2 transform -translate-y-1/2 flex justify-center">
							<div class="w-8 h-0.5 bg-gray-500 group-hover:bg-blue-400 transition-colors"></div>
						</div>
					</div>
					
					<!-- Response Panel -->
					<div 
						class="border-t border-gray-700 overflow-hidden flex-shrink-0" 
						style="height: {responseHeight}px; min-height: 200px; max-height: 70vh;"
					>
						<ResponsePanel 
							{response}
							{isExecuting}
							bind:activeResponseTab
						/>
					</div>
				{/if}
			{:else}
				<!-- Welcome State -->
				<div class="flex-1 flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900">
					<div class="text-center max-w-md mx-auto px-6">
						<!-- Electric Logo -->
						<div class="relative mb-8">
							<div class="absolute inset-0 flex items-center justify-center">
								<div class="w-20 h-20 rounded-full bg-blue-500/10 blur-xl"></div>
							</div>
							<div class="relative flex items-center justify-center">
								<svg class="w-16 h-16 text-blue-500" fill="currentColor" viewBox="0 0 24 24">
									<path d="M13 10V3L4 14h7v7l9-11h-7z"/>
								</svg>
							</div>
						</div>
						
						<!-- Main Content -->
						<div class="space-y-4">
							<h2 class="text-2xl font-bold text-white">
								Welcome to <span class="electric-gradient bg-clip-text text-transparent">Volt</span>
							</h2>
							<p class="text-gray-400 text-base leading-relaxed">
								Ready to power up your API testing? Select a collection from the sidebar to start building and testing your requests.
							</p>
						</div>
						
						<!-- Feature Cards -->
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-8">
							<div class="bg-gray-800/50 border border-gray-700/50 rounded-lg p-4 hover:bg-gray-800/70 transition-colors">
								<div class="flex items-center space-x-3">
									<div class="w-8 h-8 rounded-lg bg-green-500/20 flex items-center justify-center">
										<svg class="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
										</svg>
									</div>
									<div class="text-left">
										<h4 class="text-sm font-medium text-white">Fast Testing</h4>
										<p class="text-xs text-gray-400">Lightning quick API requests</p>
									</div>
								</div>
							</div>
							
							<div class="bg-gray-800/50 border border-gray-700/50 rounded-lg p-4 hover:bg-gray-800/70 transition-colors">
								<div class="flex items-center space-x-3">
									<div class="w-8 h-8 rounded-lg bg-blue-500/20 flex items-center justify-center">
										<svg class="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
										</svg>
									</div>
									<div class="text-left">
										<h4 class="text-sm font-medium text-white">Organized</h4>
										<p class="text-xs text-gray-400">Collections keep you tidy</p>
									</div>
								</div>
							</div>
							
							<div class="bg-gray-800/50 border border-gray-700/50 rounded-lg p-4 hover:bg-gray-800/70 transition-colors">
								<div class="flex items-center space-x-3">
									<div class="w-8 h-8 rounded-lg bg-purple-500/20 flex items-center justify-center">
										<svg class="w-4 h-4 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
										</svg>
									</div>
									<div class="text-left">
										<h4 class="text-sm font-medium text-white">Collaborative</h4>
										<p class="text-xs text-gray-400">Team-friendly workspace</p>
									</div>
								</div>
							</div>
							
							<div class="bg-gray-800/50 border border-gray-700/50 rounded-lg p-4 hover:bg-gray-800/70 transition-colors">
								<div class="flex items-center space-x-3">
									<div class="w-8 h-8 rounded-lg bg-yellow-500/20 flex items-center justify-center">
										<svg class="w-4 h-4 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"></path>
										</svg>
									</div>
									<div class="text-left">
										<h4 class="text-sm font-medium text-white">Multiple Tabs</h4>
										<p class="text-xs text-gray-400">Work on multiple requests</p>
									</div>
								</div>
							</div>
						</div>
						
						<!-- Call to Action -->
						<div class="mt-8 space-y-3">
							<p class="text-sm text-gray-500">
								Get started by selecting a collection or creating a new one
							</p>
							<div class="flex items-center justify-center space-x-2 text-xs text-gray-600">
								<kbd class="px-2 py-1 bg-gray-800 border border-gray-700 rounded text-gray-400">⌘</kbd>
								<span>+</span>
								<kbd class="px-2 py-1 bg-gray-800 border border-gray-700 rounded text-gray-400">K</kbd>
								<span class="ml-2">to search</span>
							</div>
						</div>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>

<!-- Create Collection Modal -->
<CreateCollectionModal
	show={showNewCollectionModal}
	bind:collectionName={newCollectionName}
	bind:collectionDescription={newCollectionDescription}
	isCreating={isCreatingCollection}
	on:close={() => showNewCollectionModal = false}
	on:create={(e) => {
		newCollectionName = e.detail.name;
		newCollectionDescription = e.detail.description;
		createCollection();
	}}
/>

<!-- Create Environment Modal -->
<CreateEnvironmentModal
	show={showNewEnvironmentModal}
	bind:environmentName={newEnvironmentName}
	isCreating={isCreatingEnvironment}
	on:close={() => showNewEnvironmentModal = false}
	on:create={(e) => {
		newEnvironmentName = e.detail.name;
		createEnvironment();
	}}
/>

<!-- Error Toast -->
{#if error}
	<div class="fixed bottom-4 right-4 bg-red-900 border border-red-700 rounded-md p-4 max-w-sm z-50">
		<div class="flex">
			<div class="flex-shrink-0">
				<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
				</svg>
			</div>
			<div class="ml-3">
				<p class="text-sm text-red-200">{error}</p>
			</div>
			<div class="ml-auto pl-3">
				<button
					on:click={() => error = ''}
					class="text-red-400 hover:text-red-300"
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
	/* Remove all outline styles globally */
	* {
		outline: none !important;
	}
	
	*:focus {
		outline: none !important;
	}
	
	button:focus,
	input:focus,
	textarea:focus,
	select:focus {
		outline: none !important;
	}
</style>