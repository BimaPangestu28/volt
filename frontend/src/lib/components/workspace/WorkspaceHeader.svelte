<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Workspace } from '$lib/stores/workspace';
	import InviteMemberModal from './InviteMemberModal.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import { Plus, Settings, X, Search, Zap, UserPlus } from 'lucide-svelte';

	export let currentWorkspace: Workspace | null = null;
	export let searchResults: any[] = [];

	const dispatch = createEventDispatcher();
	let showInviteModal = false;
	let searchQuery = '';
	let showSearchResults = false;
	let isSearching = false;

	function handleCreateCollection() {
		dispatch('createCollection');
	}

	function handleInvite() {
		showInviteModal = true;
	}

	function handleLinkGenerated(event: any) {
		const { link, role, expiresInHours } = event.detail;
		dispatch('linkGenerated', { link, role, expiresInHours });
	}

	// Search functionality
	let searchTimeout: any;
	
	function handleSearchInput() {
		clearTimeout(searchTimeout);
		
		if (searchQuery.trim().length < 2) {
			showSearchResults = false;
			return;
		}

		searchTimeout = setTimeout(() => {
			performSearch();
		}, 300); // Debounce search
	}

	async function performSearch() {
		if (!searchQuery.trim() || !currentWorkspace?.id) return;

		isSearching = true;
		showSearchResults = true;

		try {
			// Dispatch search event to parent component
			dispatch('search', { 
				query: searchQuery.trim(),
				workspaceId: currentWorkspace.id 
			});
		} catch (error) {
			console.error('Search error:', error);
		} finally {
			isSearching = false;
		}
	}

	function clearSearch() {
		searchQuery = '';
		showSearchResults = false;
		clearTimeout(searchTimeout);
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			clearSearch();
		}
	}

	// Handle keyboard shortcut (Cmd+K)
	function handleGlobalKeyDown(event: KeyboardEvent) {
		if ((event.metaKey || event.ctrlKey) && event.key === 'k') {
			event.preventDefault();
			const searchInput = document.querySelector('.global-search-input') as HTMLInputElement;
			if (searchInput) {
				searchInput.focus();
			}
		}
	}
</script>

<svelte:window on:keydown={handleGlobalKeyDown} />

<header class="bg-gray-800 border-b border-gray-700 flex-shrink-0">
	<div class="flex justify-between items-center h-16 px-6">
		<!-- Left side - Logo and workspace -->
		<div class="flex items-center space-x-4 min-w-0 flex-1">
			<a href="/dashboard" class="flex items-center space-x-2 electric-gradient bg-clip-text text-transparent hover:opacity-80 transition-opacity flex-shrink-0">
				<Zap class="w-5 h-5 text-blue-500" />
				<span class="text-base font-bold">Volt</span>
			</a>
			<div class="w-px h-6 bg-gray-600"></div>
			<div class="flex items-center space-x-2 min-w-0">
				<svg class="w-4 h-4 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0H5m14 0l-2-2H7l-2 2"/>
				</svg>
				<h1 class="text-xs font-medium text-white truncate">{currentWorkspace?.name || 'Loading...'}</h1>
			</div>
		</div>
		
		<!-- Center - Search -->
		<div class="flex items-center justify-center flex-1 max-w-lg mx-8">
			<div class="relative w-full max-w-sm">
				<Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
				<input
					type="text"
					bind:value={searchQuery}
					on:input={handleSearchInput}
					on:keydown={handleKeyDown}
					placeholder="Search collections, requests..."
					class="global-search-input w-full pl-10 pr-16 py-2 bg-gray-700 border border-gray-600 rounded-md text-xs text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:bg-gray-650 transition-colors"
				/>
				
				{#if searchQuery}
					<IconButton
						size="sm"
						class="absolute right-8 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-white transition-colors"
						on:click={clearSearch}
					>
						<X class="w-4 h-4" />
					</IconButton>
				{/if}

				<div class="absolute right-3 top-1/2 transform -translate-y-1/2">
					<kbd class="px-1.5 py-0.5 text-xs font-medium text-gray-400 bg-gray-800 border border-gray-600 rounded text-center min-w-[20px]">⌘K</kbd>
				</div>

				<!-- Search Results Dropdown -->
				{#if showSearchResults}
					<div class="absolute top-full left-0 right-0 mt-1 bg-gray-800 border border-gray-700 rounded-md shadow-xl max-h-64 overflow-y-auto z-50">
						{#if isSearching}
							<div class="px-4 py-3 text-center">
								<div class="flex items-center justify-center space-x-2">
									<svg class="w-4 h-4 animate-spin text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
									</svg>
									<span class="text-sm text-gray-400">Searching...</span>
								</div>
							</div>
						{:else if searchResults.length === 0}
							<div class="px-4 py-3 text-center">
								<p class="text-sm text-gray-400">No results found for "{searchQuery}"</p>
							</div>
						{:else}
							{#each searchResults as result}
								<Button
									variant="ghost"
									size="sm"
									class="w-full px-4 py-2 text-left hover:bg-gray-700 transition-colors border-b border-gray-700 last:border-b-0 justify-start"
									on:click={() => dispatch('searchResultClick', result)}
								>
									<div class="flex items-center space-x-3">
										<div class="flex-shrink-0">
											{#if result.type === 'collection'}
												<svg class="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
												</svg>
											{:else if result.type === 'request'}
												<div class="w-4 h-4 rounded text-xs font-bold flex items-center justify-center text-white" 
													 class:bg-green-600={result.method === 'GET'}
													 class:bg-blue-600={result.method === 'POST'}
													 class:bg-yellow-600={result.method === 'PUT'}
													 class:bg-red-600={result.method === 'DELETE'}
													 class:bg-purple-600={result.method === 'PATCH'}>
													{result.method?.charAt(0) || 'R'}
												</div>
											{/if}
										</div>
										<div class="flex-1 min-w-0">
											<p class="text-sm font-medium text-white truncate">{result.name}</p>
											{#if result.description}
												<p class="text-xs text-gray-400 truncate">{result.description}</p>
											{/if}
											{#if result.type === 'request' && result.url}
												<p class="text-xs text-gray-500 truncate">{result.url}</p>
											{/if}
										</div>
										<div class="flex-shrink-0">
											<span class="text-xs text-gray-500 capitalize">{result.type}</span>
										</div>
									</div>
								</Button>
							{/each}
						{/if}
					</div>
				{/if}
			</div>
		</div>
		
		<!-- Right side - Actions -->
		<div class="flex items-center space-x-3 justify-end flex-1">
			<!-- Invite button -->
			<Button
				on:click={handleInvite}
				variant="secondary"
				size="sm"
				class="px-3 py-2 text-xs font-medium text-gray-300 hover:text-white border border-gray-600 rounded-md hover:border-gray-500 hover:bg-gray-700 transition-all duration-200 flex items-center space-x-2"
			>
				<UserPlus class="w-4 h-4" />
				<span>Invite</span>
			</Button>
			
			<!-- New Collection button -->
			<Button
				on:click={handleCreateCollection}
				variant="primary"
				size="sm"
				class="electric-gradient text-white px-4 py-2 text-xs rounded-md font-medium hover:opacity-90 transition-opacity flex items-center space-x-2"
			>
				<Plus class="w-4 h-4" />
				<span>New</span>
			</Button>
			
			<!-- Settings button -->
			<IconButton
				size="sm"
				class="p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 transition-all duration-200"
			>
				<Settings class="w-5 h-5" />
			</IconButton>
		</div>
	</div>
</header>

<!-- Invite Member Modal -->
<InviteMemberModal
	bind:show={showInviteModal}
	workspaceId={currentWorkspace?.id || ''}
	on:close={() => showInviteModal = false}
	on:linkGenerated={handleLinkGenerated}
/>

<style>
	.electric-gradient {
		background: linear-gradient(135deg, var(--electric-blue), var(--deep-blue));
	}
	
	.bg-clip-text {
		background-clip: text;
		-webkit-background-clip: text;
	}
</style>