<script lang="ts">
	export let size: 'sm' | 'md' | 'lg' = 'md';
	export let color: 'blue' | 'white' | 'gray' = 'blue';
	
	$: sizeClass = {
		sm: 'h-4 w-4',
		md: 'h-6 w-6', 
		lg: 'h-8 w-8'
	}[size];
	
	$: colorClass = {
		blue: 'text-blue-600',
		white: 'text-white',
		gray: 'text-gray-600'
	}[color];
</script>

<div class="inline-flex items-center justify-center">
	<svg 
		class="animate-spin {sizeClass} {colorClass}" 
		xmlns="http://www.w3.org/2000/svg" 
		fill="none" 
		viewBox="0 0 24 24"
	>
		<!-- Outer ring -->
		<circle 
			class="opacity-25" 
			cx="12" 
			cy="12" 
			r="10" 
			stroke="currentColor" 
			stroke-width="3"
		></circle>
		
		<!-- Electric bolt -->
		<path 
			class="opacity-75" 
			fill="currentColor" 
			d="M12 2C6.477 2 2 6.477 2 12h3c0-3.866 3.134-7 7-7V2z"
		></path>
		
		<!-- Inner highlight -->
		<circle 
			class="opacity-50" 
			cx="12" 
			cy="12" 
			r="6" 
			stroke="currentColor" 
			stroke-width="2"
		></circle>
	</svg>
</div>

<style>
	/* Custom electric animation */
	@keyframes electric-spin {
		0% { 
			transform: rotate(0deg); 
			filter: drop-shadow(0 0 2px currentColor);
		}
		25% { 
			filter: drop-shadow(0 0 4px currentColor);
		}
		50% { 
			transform: rotate(180deg); 
			filter: drop-shadow(0 0 6px currentColor);
		}
		75% { 
			filter: drop-shadow(0 0 4px currentColor);
		}
		100% { 
			transform: rotate(360deg); 
			filter: drop-shadow(0 0 2px currentColor);
		}
	}
	
	svg {
		animation: electric-spin 1s linear infinite;
	}
</style>