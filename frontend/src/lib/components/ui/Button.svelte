<script lang="ts">
	export let variant: 'primary' | 'secondary' | 'ghost' | 'danger' = 'primary';
	export let size: 'xs' | 'sm' | 'md' | 'lg' = 'md';
	export let disabled = false;
	export let loading = false;
	export let type: 'button' | 'submit' | 'reset' = 'button';
	export let href: string | undefined = undefined;

	// Button classes based on variant and size
	$: classes = [
		// Base classes
		'inline-flex items-center justify-center font-medium rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2',
		
		// Size classes
		size === 'xs' ? 'px-2 py-0.5 text-xs' :
		size === 'sm' ? 'px-3 py-2 text-sm' :
		size === 'md' ? 'px-4 py-2.5 text-sm' :
		size === 'lg' ? 'px-6 py-3 text-base' : '',
		
		// Variant classes
		variant === 'primary' ? 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500 disabled:bg-blue-300' :
		variant === 'secondary' ? 'bg-slate-100 text-slate-700 hover:bg-slate-200 focus:ring-slate-500 disabled:bg-slate-50 disabled:text-slate-400' :
		variant === 'ghost' ? 'text-slate-600 hover:bg-slate-100 focus:ring-slate-500 disabled:text-slate-400' :
		variant === 'danger' ? 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500 disabled:bg-red-300' : '',
		
		// Disabled state
		disabled ? 'cursor-not-allowed opacity-60' : 'cursor-pointer',
		
		// Custom classes
		$$props.class || ''
	].filter(Boolean).join(' ');
</script>

{#if href}
	<a {href} class={classes} role="button" tabindex={disabled ? -1 : 0}>
		{#if loading}
			<svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
		{/if}
		<slot />
	</a>
{:else}
	<button {type} {disabled} class={classes} on:click>
		{#if loading}
			<svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
		{/if}
		<slot />
	</button>
{/if}