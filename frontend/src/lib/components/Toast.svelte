<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	
	export let type: 'success' | 'error' | 'info' = 'info';
	export let message: string = '';
	export let duration: number = 3000;
	
	const dispatch = createEventDispatcher();
	
	let visible = true;
	
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

{#if visible}
	<div 
		class="fixed top-4 right-4 z-50 transform transition-all duration-300 ease-in-out {visible ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'}"
		role="alert"
	>
		<div class="flex items-center p-4 rounded-lg shadow-lg border max-w-sm {
			type === 'success' ? 'bg-green-50 border-green-200 text-green-800' :
			type === 'error' ? 'bg-red-50 border-red-200 text-red-800' :
			'bg-blue-50 border-blue-200 text-blue-800'
		}">
			<!-- Icon -->
			<div class="flex-shrink-0 mr-3">
				{#if type === 'success'}
					<svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
					</svg>
				{:else if type === 'error'}
					<svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
					</svg>
				{:else}
					<svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
					</svg>
				{/if}
			</div>
			
			<!-- Message -->
			<div class="flex-1 text-sm font-medium">
				{message}
			</div>
			
			<!-- Close button -->
			<button
				on:click={hide}
				class="flex-shrink-0 ml-3 text-gray-400 hover:text-gray-600 transition-colors"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>
			</button>
		</div>
	</div>
{/if}