<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { apiClient } from '$lib/api/client';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import Button from '../ui/Button.svelte';
	import IconButton from '../ui/IconButton.svelte';
	import { RefreshCw, Plus, Search, Link } from 'lucide-svelte';

	export let webhooks: any[] = [];
	export let selectedWebhook: any = null;
	export let isLoading = false;

	const dispatch = createEventDispatcher();

	let searchQuery = '';
	let filteredWebhooks: any[] = [];

	// Filter webhooks based on search query
	$: {
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase();
			filteredWebhooks = webhooks.filter(webhook =>
				webhook.name.toLowerCase().includes(query) ||
				(webhook.description && webhook.description.toLowerCase().includes(query)) ||
				webhook.url.toLowerCase().includes(query)
			);
		} else {
			filteredWebhooks = webhooks;
		}
	}

	function selectWebhook(webhook: any) {
		dispatch('selectWebhook', webhook);
	}

	function createWebhook() {
		dispatch('createWebhook');
	}

	function refreshWebhooks() {
		dispatch('refreshWebhooks');
	}

	function getStatusColor(status: string) {
		switch (status) {
			case 'active':
				return 'text-green-400';
			case 'inactive':
				return 'text-gray-400';
			default:
				return 'text-yellow-400';
		}
	}

	function getMethodBadgeColor(methods: string[]) {
		if (methods.includes('POST')) return 'bg-blue-600';
		if (methods.includes('GET')) return 'bg-green-600';
		if (methods.includes('PUT')) return 'bg-yellow-600';
		if (methods.includes('DELETE')) return 'bg-red-600';
		return 'bg-gray-600';
	}

	function formatLastRequest(lastRequestAt: string) {
		if (!lastRequestAt) return 'Never';
		const date = new Date(lastRequestAt);
		const now = new Date();
		const diff = now.getTime() - date.getTime();
		
		if (diff < 60000) return 'Just now';
		if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
		if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
		return `${Math.floor(diff / 86400000)}d ago`;
	}
</script>

<div class="w-80 bg-gray-800 border-r border-gray-700 flex flex-col h-full">
	<!-- Header -->
	<div class="p-4 border-b border-gray-700">
		<div class="flex items-center justify-between mb-4">
			<h2 class="text-lg font-semibold text-white flex items-center">
				<svg class="w-5 h-5 mr-2 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path>
				</svg>
				Webhooks
			</h2>
			<div class="flex items-center space-x-2">
				<IconButton
					size="sm"
					class="p-1.5 text-gray-400 hover:text-white hover:bg-gray-700 rounded transition-colors"
					title="Refresh webhooks"
					on:click={refreshWebhooks}
				>
					<RefreshCw class="w-4 h-4" />
				</IconButton>
				<Button
					variant="primary"
					size="sm"
					class="px-3 py-1.5 bg-purple-600 hover:bg-purple-700 text-white rounded text-sm font-medium transition-colors flex items-center"
					on:click={createWebhook}
				>
					<Plus class="w-4 h-4 mr-1" />
					New
				</Button>
			</div>
		</div>

		<!-- Search -->
		<div class="relative">
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Search webhooks..."
				class="w-full px-3 py-2 pl-9 bg-gray-700 border border-gray-600 rounded text-sm text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
			/>
			<Search class="absolute left-3 top-2.5 w-4 h-4 text-gray-400" />
		</div>
	</div>

	<!-- Loading State -->
	{#if isLoading}
		<div class="flex-1 flex items-center justify-center">
			<div class="text-center">
				<LoadingSpinner size="lg" />
				<p class="text-gray-400 text-sm mt-2">Loading webhooks...</p>
			</div>
		</div>
	{:else if filteredWebhooks.length === 0 && searchQuery}
		<!-- No Search Results -->
		<div class="flex-1 flex items-center justify-center p-4">
			<div class="text-center">
				<svg class="mx-auto h-12 w-12 text-gray-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
				</svg>
				<p class="text-gray-400 text-sm">No webhooks found for "{searchQuery}"</p>
			</div>
		</div>
	{:else if filteredWebhooks.length === 0}
		<!-- Empty State -->
		<div class="flex-1 flex items-center justify-center p-4">
			<div class="text-center max-w-xs">
				<svg class="mx-auto h-16 w-16 text-gray-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path>
				</svg>
				<h3 class="text-lg font-medium text-white mb-2">No Webhooks Yet</h3>
				<p class="text-gray-400 text-sm mb-4">Create your first webhook endpoint to start testing incoming requests.</p>
				<Button
					variant="primary"
					size="sm"
					class="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded text-sm font-medium transition-colors"
					on:click={createWebhook}
				>
					Create Webhook
				</Button>
			</div>
		</div>
	{:else}
		<!-- Webhooks List -->
		<div class="flex-1 overflow-y-auto">
			<div class="p-2">
				{#each filteredWebhooks as webhook (webhook.id)}
					<Button
						variant="ghost"
						size="md"
						class="w-full p-3 mb-2 text-left rounded-lg border transition-all hover:bg-gray-700/50 {
							selectedWebhook?.id === webhook.id 
								? 'bg-purple-600/20 border-purple-500 text-white' 
								: 'bg-gray-800/50 border-gray-600 text-gray-300 hover:border-gray-500'
						} justify-start"
						on:click={() => selectWebhook(webhook)}
					>
						<div class="flex items-start justify-between mb-2">
							<div class="flex-1 min-w-0">
								<h3 class="font-medium text-sm truncate mb-1">{webhook.name}</h3>
								{#if webhook.description}
									<p class="text-xs text-gray-400 truncate">{webhook.description}</p>
								{/if}
							</div>
							<div class="flex items-center space-x-1 ml-2">
								<!-- Status indicator -->
								<div class="w-2 h-2 rounded-full {webhook.status === 'active' ? 'bg-green-400' : 'bg-gray-400'}"></div>
								<!-- Method badge -->
								<span class="text-[10px] px-1.5 py-0.5 rounded {getMethodBadgeColor(webhook.config?.methods || ['POST'])} text-white font-medium">
									{webhook.config?.methods?.[0] || 'POST'}
								</span>
							</div>
						</div>

						<div class="flex items-center justify-between text-xs">
							<div class="flex items-center space-x-3 text-gray-500">
								<span class="flex items-center">
									<svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
									</svg>
									{webhook.stats?.total_requests || 0}
								</span>
								<span class="flex items-center">
									<svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
									</svg>
									{formatLastRequest(webhook.stats?.last_request_at)}
								</span>
							</div>
							<span class="{getStatusColor(webhook.status)} capitalize">
								{webhook.status}
							</span>
						</div>

						<!-- URL preview -->
						<div class="mt-2 text-xs">
							<span class="text-gray-500">URL:</span>
							<span class="text-gray-400 font-mono ml-1 truncate block">
								{webhook.url}
							</span>
						</div>
					</Button>
				{/each}
			</div>
		</div>
	{/if}
</div>

<style>
	/* Custom scrollbar */
	:global(.overflow-y-auto::-webkit-scrollbar) {
		width: 4px;
	}
	
	:global(.overflow-y-auto::-webkit-scrollbar-track) {
		background: transparent;
	}
	
	:global(.overflow-y-auto::-webkit-scrollbar-thumb) {
		background: #4B5563;
		border-radius: 2px;
	}
	
	:global(.overflow-y-auto::-webkit-scrollbar-thumb:hover) {
		background: #6B7280;
	}
</style>