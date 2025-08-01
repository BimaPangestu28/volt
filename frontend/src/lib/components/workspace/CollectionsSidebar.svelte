<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Collection } from '$lib/stores/workspace';
	import { toastStore } from '$lib/stores/toast';
	import { apiClient } from '$lib/api/client';
	import RenameCollectionModal from './RenameCollectionModal.svelte';
	import IconButton from '../ui/IconButton.svelte';
	import Button from '../ui/Button.svelte';
	import { ChevronRight, Star, MoreVertical, Edit, Copy, Trash, Plus, Download, X, Zap, Clock } from 'lucide-svelte';

	export let collections: Collection[] = [];
	export let selectedCollection: Collection | null = null;
	export let selectedRequest: any = null;
	export let requests: any[] = [];
	export let isLoading = false;
	export let activeTab: 'collections' | 'environments' | 'webhooks' = 'collections';

	// Local state for UI
	let showOptionsMenu: string | null = null;
	let showRenameModal = false;
	let collectionToRename: Collection | null = null;
	let isRenaming = false;
	
	// Search functionality
	let searchQuery = '';
	let filteredCollections: Collection[] = [];
	let filteredRequests: any[] = [];

	const dispatch = createEventDispatcher();

	// Reactive statements for filtering
	$: {
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase();
			
			// Filter collections
			filteredCollections = collections.filter(collection =>
				collection.name.toLowerCase().includes(query) ||
				(collection.description && collection.description.toLowerCase().includes(query))
			);
			
			// Filter requests
			filteredRequests = requests.filter(request =>
				request.name.toLowerCase().includes(query) ||
				(request.url && request.url.toLowerCase().includes(query)) ||
				(request.method && request.method.toLowerCase().includes(query))
			);
		} else {
			filteredCollections = collections;
			filteredRequests = requests;
		}
	}

	function clearSearch() {
		searchQuery = '';
	}

	function handleSelectCollection(collection: Collection) {
		dispatch('selectCollection', collection);
	}

	function handleSelectRequest(request: any) {
		dispatch('selectRequest', request);
	}

	function handleNewRequest() {
		dispatch('newRequest');
	}

	function handleCreateCollection() {
		dispatch('createCollection');
	}

	async function handleToggleFavorite(collection: Collection, event: Event) {
		event.stopPropagation(); // Prevent collection toggle
		
		try {
			const result = await apiClient.toggleCollectionFavorite(collection.id);
			
			// Update collection in local state
			collection.is_favorited = result.is_favorited;
			collections = collections; // Trigger reactivity
			
			toastStore.success(result.message);
			
			// Dispatch event to refresh collections data
			dispatch('refreshCollections');
		} catch (error: any) {
			const errorMessage = error.response?.data?.error || 'Failed to update favorite status';
			toastStore.error(errorMessage);
		}
	}

	function handleCollectionOptions(collection: Collection, event: Event) {
		event.stopPropagation(); // Prevent collection toggle
		
		// Toggle options menu for this specific collection
		if (showOptionsMenu === collection.id) {
			showOptionsMenu = null;
		} else {
			showOptionsMenu = collection.id;
		}
	}

	function handleRenameCollection(collection: Collection) {
		showOptionsMenu = null;
		collectionToRename = collection;
		showRenameModal = true;
	}

	async function confirmRename(event: any) {
		if (!collectionToRename) return;
		
		const newName = event.detail.newName;
		isRenaming = true;
		
		try {
			await apiClient.updateCollection(collectionToRename.id, {
				name: newName,
				description: collectionToRename.description
			});
			
			toastStore.success(`Collection renamed to "${newName}"`);
			showRenameModal = false;
			collectionToRename = null;
			
			// Refresh collections
			dispatch('refreshCollections');
		} catch (error: any) {
			const errorMessage = error.response?.data?.error || 'Failed to rename collection';
			toastStore.error(errorMessage);
		} finally {
			isRenaming = false;
		}
	}

	function closeRenameModal() {
		showRenameModal = false;
		collectionToRename = null;
		isRenaming = false;
	}

	async function handleDuplicateCollection(collection: Collection) {
		showOptionsMenu = null;
		
		try {
			const duplicateName = `${collection.name} Copy`;
			await apiClient.createCollection(collection.workspace_id, duplicateName, collection.description);
			
			toastStore.success(`Collection "${duplicateName}" created`);
			
			// Refresh collections
			dispatch('refreshCollections');
		} catch (error: any) {
			const errorMessage = error.response?.data?.error || 'Failed to duplicate collection';
			toastStore.error(errorMessage);
		}
	}

	async function handleDeleteCollection(collection: Collection) {
		showOptionsMenu = null;
		const confirmed = confirm(`Are you sure you want to delete "${collection.name}"? This action cannot be undone.`);
		if (!confirmed) return;
		
		try {
			await apiClient.deleteCollection(collection.id);
			
			toastStore.success(`Collection "${collection.name}" deleted`);
			
			// Refresh collections
			dispatch('refreshCollections');
		} catch (error: any) {
			const errorMessage = error.response?.data?.error || 'Failed to delete collection';
			toastStore.error(errorMessage);
		}
	}

	// Close options menu when clicking outside
	function handleClickOutside(event: MouseEvent) {
		if (showOptionsMenu && !event.target?.closest('.options-menu')) {
			showOptionsMenu = null;
		}
	}
</script>

<!-- Global click listener to close options menu -->
<svelte:window on:click={handleClickOutside} />

<div class="flex bg-gray-800 border-r border-gray-700">
	<!-- Left Navigation -->
	<div class="w-20 bg-gray-900 border-r border-gray-700 flex flex-col">
		<div class="py-4 space-y-3">
			<!-- Collections Icon -->
			<div class="flex flex-col items-center px-1">
				<IconButton
					variant="{activeTab === 'collections' ? 'default' : 'ghost'}"
					size="md"
					class="w-10 h-10 mb-1.5 {activeTab === 'collections' 
						? 'bg-blue-600 text-white hover:bg-blue-500' 
						: 'text-gray-400 hover:bg-gray-700 hover:text-white'}"
					on:click={() => dispatch('switchTab', 'collections')}
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
					</svg>
				</IconButton>
				<span class="text-[9px] font-medium text-center leading-tight px-0.5 {activeTab === 'collections' ? 'text-blue-400' : 'text-gray-500'}">Collections</span>
			</div>
			
			<!-- Environments Icon -->
			<div class="flex flex-col items-center px-1">
				<IconButton
					variant="{activeTab === 'environments' ? 'default' : 'ghost'}"
					size="md"
					class="w-10 h-10 mb-1.5 {activeTab === 'environments' 
						? 'bg-blue-600 text-white hover:bg-blue-500' 
						: 'text-gray-400 hover:bg-gray-700 hover:text-white'}"
					on:click={() => dispatch('switchTab', 'environments')}
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 0 1 9-9"></path>
					</svg>
				</IconButton>
				<span class="text-[9px] font-medium text-center leading-tight px-0.5 {activeTab === 'environments' ? 'text-blue-400' : 'text-gray-500'}">Environments</span>
			</div>
			
			<!-- Webhooks Icon -->
			<div class="flex flex-col items-center px-1">
				<IconButton
					variant="ghost"
					size="md"
					class="w-10 h-10 mb-1.5 text-gray-400 hover:bg-gray-700 hover:text-white"
					on:click={() => dispatch('switchTab', 'webhooks')}
				>
					<Zap class="w-5 h-5" />
				</IconButton>
				<span class="text-[9px] text-gray-500 font-medium text-center leading-tight px-0.5">Webhooks</span>
			</div>
			
			<!-- Flows Icon -->
			<div class="flex flex-col items-center px-1">
				<IconButton
					variant="ghost"
					size="md"
					class="w-10 h-10 mb-1.5 text-gray-400 hover:bg-gray-700 hover:text-white"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
					</svg>
				</IconButton>
				<span class="text-[9px] text-gray-500 font-medium text-center leading-tight px-0.5">Flows</span>
			</div>
			
			<!-- History Icon -->
			<div class="flex flex-col items-center px-1">
				<IconButton
					variant="ghost"
					size="md"
					class="w-10 h-10 mb-1.5 text-gray-400 hover:bg-gray-700 hover:text-white"
				>
					<Clock class="w-5 h-5" />
				</IconButton>
				<span class="text-[9px] text-gray-500 font-medium text-center leading-tight px-0.5">History</span>
			</div>
		</div>
	</div>

	<!-- Content Panel -->
	<div class="w-80 bg-gray-800 flex flex-col">
		{#if activeTab === 'collections'}
			<!-- Collections Header with Search -->
			<div class="p-4 border-b border-gray-700">
				<div class="flex items-center justify-between mb-4">
					<h2 class="text-lg font-semibold text-white">Collections</h2>
				<div class="flex items-center space-x-1">
					<IconButton
						variant="ghost"
						size="sm"
						class="p-1.5 text-gray-400 hover:text-white hover:bg-gray-700"
						title="Import"
					>
						<Download class="w-4 h-4" />
					</IconButton>
					<IconButton
						variant="ghost"
						size="sm"
						class="p-1.5 text-gray-400 hover:text-white hover:bg-gray-700"
						title="New Collection"
						on:click={handleCreateCollection}
					>
						<Plus class="w-4 h-4" />
					</IconButton>
					<IconButton
						variant="ghost"
						size="sm"
						class="p-1.5 text-gray-400 hover:text-white hover:bg-gray-700"
						title="More actions"
					>
						<MoreVertical class="w-4 h-4" />
					</IconButton>
				</div>
			</div>
			
			<!-- Search Bar -->
			<div class="relative">
				<svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
				</svg>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search collections & requests..."
					class="w-full pl-10 pr-8 py-2.5 bg-gray-700 border border-gray-600 rounded-lg text-xs text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
				/>
				{#if searchQuery}
					<IconButton
						variant="ghost"
						size="xs"
						class="absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-white"
						on:click={clearSearch}
					>
						<X class="w-4 h-4" />
					</IconButton>
				{/if}
			</div>
		</div>

		<!-- Collections Tree -->
		<div class="flex-1 overflow-y-auto bg-gray-800">
			{#if isLoading}
				<div class="p-3 space-y-2">
					{#each Array(4) as _, i}
						<div class="flex items-center space-x-2 p-2 animate-pulse">
							<div class="w-4 h-4 bg-gray-600 rounded"></div>
							<div class="flex-1 h-4 bg-gray-600 rounded"></div>
							<div class="w-4 h-4 bg-gray-600 rounded"></div>
						</div>
					{/each}
				</div>
			{:else if !collections || collections.length === 0}
				<div class="text-center py-12 px-4">
					<svg class="mx-auto h-12 w-12 text-gray-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
					</svg>
					<p class="text-xs text-gray-400 mb-4">No collections yet</p>
					<Button
						variant="ghost"
						size="xs"
						class="text-xs text-blue-400 hover:text-blue-300 font-medium"
						on:click={handleCreateCollection}
					>
						Create your first collection
					</Button>
				</div>
			{:else if searchQuery && filteredCollections.length === 0}
				<div class="text-center py-8 px-4">
					<svg class="mx-auto h-8 w-8 text-gray-500 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
					</svg>
					<p class="text-xs text-gray-400">No results found for "{searchQuery}"</p>
				</div>
			{:else}
				<div>
					{#each filteredCollections || [] as collection}
						<div class="group border-b border-gray-700/50 last:border-b-0">
							<!-- Collection Header -->
							<div class="flex items-center px-4 py-3 hover:bg-gray-750 cursor-pointer transition-colors">
								<IconButton
									size="sm"
									class="mr-3 p-1 rounded text-gray-400 hover:text-white transition-colors"
									on:click={() => handleSelectCollection(collection)}
								>
									<ChevronRight class="w-3 h-3 transform {selectedCollection?.id === collection.id ? 'rotate-90' : ''} transition-transform" />
								</IconButton>
								
								<div 
									class="flex items-center space-x-3 flex-1 min-w-0 cursor-pointer"
									on:click={() => handleSelectCollection(collection)}
									on:keydown={(e) => e.key === 'Enter' && handleSelectCollection(collection)}
									role="button"
									tabindex="0"
								>
									<svg class="w-4 h-4 text-blue-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
									</svg>
									<span class="text-xs text-white truncate font-medium flex items-center">
										{collection.name}
										{#if collection.is_favorited}
											<svg class="w-3 h-3 text-yellow-400 ml-1.5 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
												<path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.196-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path>
											</svg>
										{/if}
									</span>
								</div>
								
								<div class="flex items-center space-x-1 transition-opacity {collection.is_favorited || showOptionsMenu === collection.id ? 'opacity-100' : 'opacity-0 group-hover:opacity-100'}">
									<!-- Star/Favorite -->
									<IconButton
										size="sm"
										class="p-1 rounded transition-colors {collection.is_favorited ? 'text-yellow-400 hover:text-yellow-300' : 'text-gray-400 hover:text-yellow-400'}"
										title="{collection.is_favorited ? 'Remove from favorites' : 'Add to favorites'}"
										on:click={(e) => handleToggleFavorite(collection, e)}
									>
										<Star class="w-3.5 h-3.5" fill="{collection.is_favorited ? 'currentColor' : 'none'}" />
									</IconButton>
									<!-- More Options -->
									<div class="relative options-menu">
										<IconButton
											size="sm"
											class="p-1 rounded text-gray-400 hover:text-white transition-colors"
											title="More actions"
											on:click={(e) => handleCollectionOptions(collection, e)}
										>
											<MoreVertical class="w-3.5 h-3.5" />
										</IconButton>
										
										<!-- Options Dropdown -->
										{#if showOptionsMenu === collection.id}
											<!-- svelte-ignore a11y-no-static-element-interactions -->
											<div 
												class="absolute top-full right-0 z-50 bg-gray-800 border border-gray-600 rounded-md shadow-xl min-w-[160px] py-1"
												on:click|stopPropagation
											>
												<Button
													variant="ghost"
													size="xs"
													class="w-full px-3 py-2 text-left text-xs text-gray-300 hover:bg-gray-700 hover:text-white transition-colors flex items-center justify-start"
													on:click={() => handleRenameCollection(collection)}
												>
													<Edit class="w-3.5 h-3.5 mr-2" />
													Rename
												</Button>
												<Button
													variant="ghost"
													size="xs"
													class="w-full px-3 py-2 text-left text-xs text-gray-300 hover:bg-gray-700 hover:text-white transition-colors flex items-center justify-start"
													on:click={() => handleDuplicateCollection(collection)}
												>
													<Copy class="w-3.5 h-3.5 mr-2" />
													Duplicate
												</Button>
												<div class="border-t border-gray-700 my-1"></div>
												<Button
													variant="ghost"
													size="xs"
													class="w-full px-3 py-2 text-left text-xs text-red-400 hover:bg-red-900/20 hover:text-red-300 transition-colors flex items-center justify-start"
													on:click={() => handleDeleteCollection(collection)}
												>
													<Trash class="w-3.5 h-3.5 mr-2" />
													Delete
												</Button>
											</div>
										{/if}
									</div>
								</div>
							</div>
							
							<!-- Collection Requests (when expanded) -->
							{#if selectedCollection?.id === collection.id}
								<div class="bg-gray-750/30 border-t border-gray-700/30">
									{#each filteredRequests as request}
										<div 
											class="flex items-center px-4 py-2 pl-12 hover:bg-gray-700/50 cursor-pointer transition-colors {selectedRequest?.id === request.id ? 'bg-blue-900/30 border-l-2 border-blue-500' : ''}"
											on:click={() => handleSelectRequest(request)}
											on:keydown={(e) => e.key === 'Enter' && handleSelectRequest(request)}
											role="button"
											tabindex="0"
										>
											<span class="text-[10px] font-mono px-1.5 py-0.5 rounded mr-2 font-semibold uppercase {
												request.method === 'GET' ? 'bg-green-700 text-green-300' : 
												request.method === 'POST' ? 'bg-blue-700 text-blue-300' : 
												request.method === 'PUT' ? 'bg-yellow-700 text-yellow-300' : 
												request.method === 'DELETE' ? 'bg-red-700 text-red-300' :
												'bg-purple-700 text-purple-300'
											}">
												{request.method}
											</span>
											<span class="text-[11px] text-gray-200 truncate flex-1 leading-tight">{request.name}</span>
										</div>
									{/each}
									
									<!-- Add New Request -->
									<Button
										variant="ghost"
										size="xs"
										class="flex items-center px-4 py-2 pl-12 hover:bg-gray-700/50 cursor-pointer transition-colors w-full text-left group justify-start"
										on:click={handleNewRequest}
									>
										<Plus class="w-3.5 h-3.5 text-gray-400 group-hover:text-blue-400 mr-2 transition-colors" />
										<span class="text-[11px] text-gray-400 group-hover:text-blue-400 transition-colors">Add request</span>
									</Button>
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
		{:else if activeTab === 'environments'}
			<!-- Environments content will be handled by parent component -->
			<div class="flex-1 flex items-center justify-center">
				<div class="text-center text-gray-400">
					<p class="text-sm">Environment panel will be displayed here</p>
				</div>
			</div>
		{/if}
	</div>
</div>

<!-- Rename Collection Modal -->
<RenameCollectionModal
	bind:show={showRenameModal}
	collectionName={collectionToRename?.name || ''}
	{isRenaming}
	on:rename={confirmRename}
	on:close={closeRenameModal}
/>