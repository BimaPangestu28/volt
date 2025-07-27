<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Workspace } from '$lib/stores/workspace';

	export let currentWorkspace: Workspace | null = null;

	const dispatch = createEventDispatcher();

	function handleCreateCollection() {
		dispatch('createCollection');
	}
</script>

<header class="bg-gray-800 border-b border-gray-700 flex-shrink-0">
	<div class="flex justify-between items-center h-16 px-6">
		<!-- Left side - Logo and workspace -->
		<div class="flex items-center space-x-4 min-w-0 flex-1">
			<a href="/dashboard" class="flex items-center space-x-2 electric-gradient bg-clip-text text-transparent hover:opacity-80 transition-opacity flex-shrink-0">
				<svg class="w-5 h-5 text-blue-500" fill="currentColor" viewBox="0 0 24 24">
					<path d="M13 10V3L4 14h7v7l9-11h-7z"/>
				</svg>
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
				<svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
				</svg>
				<input
					type="text"
					placeholder="Search Volt"
					class="w-full pl-10 pr-16 py-2 bg-gray-700 border border-gray-600 rounded-md text-xs text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:bg-gray-650 transition-colors"
				/>
				<div class="absolute right-3 top-1/2 transform -translate-y-1/2">
					<kbd class="px-1.5 py-0.5 text-xs font-medium text-gray-400 bg-gray-800 border border-gray-600 rounded text-center min-w-[20px]">⌘K</kbd>
				</div>
			</div>
		</div>
		
		<!-- Right side - Actions -->
		<div class="flex items-center space-x-3 justify-end flex-1">
			<!-- Invite button -->
			<button class="px-3 py-2 text-xs font-medium text-gray-300 hover:text-white border border-gray-600 rounded-md hover:border-gray-500 hover:bg-gray-700 transition-all duration-200">
				Invite
			</button>
			
			<!-- New Collection button -->
			<button
				on:click={handleCreateCollection}
				class="electric-gradient text-white px-4 py-2 text-xs rounded-md font-medium hover:opacity-90 transition-opacity flex items-center space-x-2"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
				</svg>
				<span>New</span>
			</button>
			
			<!-- Settings button -->
			<button class="p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 transition-all duration-200">
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
				</svg>
			</button>
		</div>
	</div>
</header>

<style>
	.electric-gradient {
		background: linear-gradient(135deg, var(--electric-blue), var(--deep-blue));
	}
	
	.bg-clip-text {
		background-clip: text;
		-webkit-background-clip: text;
	}
</style>