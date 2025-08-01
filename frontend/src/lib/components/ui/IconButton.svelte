<script lang="ts">
	export let variant: 'primary' | 'secondary' | 'ghost' | 'danger' = 'ghost';
	export let size: 'xs' | 'sm' | 'md' | 'lg' = 'md';
	export let disabled = false;
	export let type: 'button' | 'submit' | 'reset' = 'button';

	// Icon button classes
	$: classes = [
		// Base classes
		'inline-flex items-center justify-center font-medium rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2',
		
		// Size classes (square)
		size === 'xs' ? 'w-5 h-5 text-xs' :
		size === 'sm' ? 'w-8 h-8 text-sm' :
		size === 'md' ? 'w-10 h-10 text-sm' :
		size === 'lg' ? 'w-12 h-12 text-base' : '',
		
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

<button {type} {disabled} class={classes} on:click>
	<slot />
</button>