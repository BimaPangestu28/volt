<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Collection } from '$lib/stores/workspace';

	export let collections: Collection[] = [];
	export let selectedCollection: Collection | null = null;
	export let selectedRequest: any = null;
	export let requests: any[] = [];
	export let isLoading = false;

	const dispatch = createEventDispatcher();

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
</script>

<div class="flex bg-gray-800 border-r border-gray-700">
	<!-- Left Navigation -->
	<div class="w-20 bg-gray-900 border-r border-gray-700 flex flex-col">
		<div class="py-4 space-y-3">
			<!-- Collections Icon -->
			<div class="flex flex-col items-center px-1">
				<button class="w-10 h-10 flex items-center justify-center rounded-md bg-blue-600 text-white mb-1.5 hover:bg-blue-500 transition-colors">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
					</svg>
				</button>
				<span class="text-[9px] text-blue-400 font-medium text-center leading-tight px-0.5">Collections</span>
			</div>
			
			<!-- Environments Icon -->
			<div class="flex flex-col items-center px-1">
				<button class="w-10 h-10 flex items-center justify-center rounded-md text-gray-400 hover:bg-gray-700 hover:text-white transition-colors mb-1.5">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 919-9"></path>
					</svg>
				</button>
				<span class="text-[9px] text-gray-500 font-medium text-center leading-tight px-0.5">Environments</span>
			</div>
			
			<!-- APIs Icon -->
			<div class="flex flex-col items-center px-1">
				<button class="w-10 h-10 flex items-center justify-center rounded-md text-gray-400 hover:bg-gray-700 hover:text-white transition-colors mb-1.5">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
					</svg>
				</button>
				<span class="text-[9px] text-gray-500 font-medium text-center leading-tight px-0.5">APIs</span>
			</div>
			
			<!-- Flows Icon -->
			<div class="flex flex-col items-center px-1">
				<button class="w-10 h-10 flex items-center justify-center rounded-md text-gray-400 hover:bg-gray-700 hover:text-white transition-colors mb-1.5">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
					</svg>
				</button>
				<span class="text-[9px] text-gray-500 font-medium text-center leading-tight px-0.5">Flows</span>
			</div>
			
			<!-- History Icon -->
			<div class="flex flex-col items-center px-1">
				<button class="w-10 h-10 flex items-center justify-center rounded-md text-gray-400 hover:bg-gray-700 hover:text-white transition-colors mb-1.5">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
					</svg>
				</button>
				<span class="text-[9px] text-gray-500 font-medium text-center leading-tight px-0.5">History</span>
			</div>
		</div>
	</div>

	<!-- Collections Panel -->
	<div class="w-80 bg-gray-800 flex flex-col">
		<!-- Header with Search -->
		<div class="p-4 border-b border-gray-700">
			<div class="flex items-center justify-between mb-4">
				<h2 class="text-lg font-semibold text-white">Collections</h2>
				<div class="flex items-center space-x-1">
					<button
						class="p-1.5 rounded text-gray-400 hover:text-white hover:bg-gray-700 transition-colors"
						title="Import"
					>
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"></path>
						</svg>
					</button>
					<button
						on:click={handleCreateCollection}
						class="p-1.5 rounded text-gray-400 hover:text-white hover:bg-gray-700 transition-colors"
						title="New Collection"
					>
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
						</svg>
					</button>
					<button
						class="p-1.5 rounded text-gray-400 hover:text-white hover:bg-gray-700 transition-colors"
						title="More actions"
					>
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path>
						</svg>
					</button>
				</div>
			</div>
			
			<!-- Search Bar -->
			<div class="relative">
				<svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
				</svg>
				<input
					type="text"
					placeholder="Search collections"
					class="w-full pl-10 pr-3 py-2.5 bg-gray-700 border border-gray-600 rounded-lg text-xs text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
				/>
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
					<button
						on:click={handleCreateCollection}
						class="text-xs text-blue-400 hover:text-blue-300 font-medium"
					>
						Create your first collection
					</button>
				</div>
			{:else}
				<div>
					{#each collections || [] as collection}
						<div class="group border-b border-gray-700/50 last:border-b-0">
							<!-- Collection Header -->
							<div class="flex items-center px-4 py-3 hover:bg-gray-750 cursor-pointer transition-colors">
								<button
									class="mr-3 p-1 rounded text-gray-400 hover:text-white transition-colors"
									on:click={() => handleSelectCollection(collection)}
								>
									<svg class="w-3 h-3 transform {selectedCollection?.id === collection.id ? 'rotate-90' : ''} transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
									</svg>
								</button>
								
								<div class="flex items-center space-x-3 flex-1 min-w-0">
									<svg class="w-4 h-4 text-blue-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
									</svg>
									<span class="text-xs text-white truncate font-medium">{collection.name}</span>
								</div>
								
								<div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
									<!-- Star/Favorite -->
									<button class="p-1 rounded text-gray-400 hover:text-yellow-400 transition-colors" title="Add to favorites">
										<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.196-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path>
										</svg>
									</button>
									<!-- More Options -->
									<button class="p-1 rounded text-gray-400 hover:text-white transition-colors" title="More actions">
										<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path>
										</svg>
									</button>
								</div>
							</div>
							
							<!-- Collection Requests (when expanded) -->
							{#if selectedCollection?.id === collection.id}
								<div class="bg-gray-750/30 border-t border-gray-700/30">
									{#each requests as request}
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
									<button
										on:click={handleNewRequest}
										class="flex items-center px-4 py-2 pl-12 hover:bg-gray-700/50 cursor-pointer transition-colors w-full text-left group"
									>
										<svg class="w-3.5 h-3.5 text-gray-400 group-hover:text-blue-400 mr-2 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
										</svg>
										<span class="text-[11px] text-gray-400 group-hover:text-blue-400 transition-colors">Add request</span>
									</button>
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>