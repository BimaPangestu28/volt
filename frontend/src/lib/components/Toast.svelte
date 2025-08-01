<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import IconButton from './ui/IconButton.svelte';
	
	export let type: 'success' | 'error' | 'info' = 'info';
	export let message: string = '';
	export let duration: number = 3000;
	
	const dispatch = createEventDispatcher();
	
	let visible = false;
	let mounted = false;
	
	// Entrance animation
	setTimeout(() => {
		visible = true;
	}, 10);
	
	// Auto hide after duration
	if (duration > 0) {
		setTimeout(() => {
			hide();
		}, duration);
	}
	
	function hide() {
		visible = false;
		setTimeout(() => {
			dispatch('close');
		}, 300); // Wait for fade out animation
	}
</script>

<div 
	class="transform transition-all duration-300 ease-out {visible ? 'translate-x-0 opacity-100 scale-100' : 'translate-x-full opacity-0 scale-95'}"
	role="alert"
>
	<div class="relative overflow-hidden rounded-xl shadow-2xl border backdrop-blur-sm max-w-sm {
		type === 'success' ? 'bg-gray-800/95 border-green-500/30' :
		type === 'error' ? 'bg-gray-800/95 border-red-500/30' :
		'bg-gray-800/95 border-blue-500/30'
	}">
		<!-- Content -->
		<div class="flex items-center p-4">
			<!-- Icon with background -->
			<div class="flex-shrink-0 mr-3">
				<div class="w-8 h-8 rounded-full flex items-center justify-center {
					type === 'success' ? 'bg-green-500/20' :
					type === 'error' ? 'bg-red-500/20' :
					'bg-blue-500/20'
				}">
					{#if type === 'success'}
						<svg class="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
						</svg>
					{:else if type === 'error'}
						<svg class="w-4 h-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
						</svg>
					{:else}
						<svg class="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
						</svg>
					{/if}
				</div>
			</div>
			
			<!-- Message -->
			<div class="flex-1 text-sm font-medium text-white">
				{message}
			</div>
			
			<!-- Close button -->
			<div class="flex-shrink-0 ml-3">
				<IconButton
					variant="ghost"
					size="sm"
					on:click={hide}
					class="text-gray-400 hover:text-white hover:bg-gray-700/50"
					title="Close"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</IconButton>
			</div>
		</div>
		
		<!-- Progress bar -->
		{#if duration > 0}
			<div class="absolute bottom-0 left-0 w-full h-1 bg-gray-700/50">
				<div 
					class="h-full transition-all ease-linear {
						type === 'success' ? 'bg-green-500' :
						type === 'error' ? 'bg-red-500' :
						'bg-blue-500'
					}"
					style="animation: shrink {duration}ms linear forwards;"
				></div>
			</div>
		{/if}
	</div>
</div>

<style>
	@keyframes shrink {
		from { width: 100%; }
		to { width: 0%; }
	}
</style>