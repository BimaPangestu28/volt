<script lang="ts">
	import type { ComponentType } from 'svelte';
	import Button from './Button.svelte';

	export let title: string;
	export let description: string;
	export let icon: ComponentType;
	export let iconColor: string = 'bg-blue-600';
	export let disabled: boolean = false;
</script>

<Button 
	variant="ghost"
	class="w-full flex items-center p-2 text-left bg-slate-50 hover:bg-slate-100 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition-colors group"
	{disabled}
	on:click
>
	<div class="p-1.5 {iconColor} rounded-md mr-2 flex-shrink-0">
		<svelte:component this={icon} class="w-3 h-3 text-white" />
	</div>
	<div class="flex-1 min-w-0">
		<p class="text-xs font-medium text-slate-900 truncate">{title}</p>
		<p class="text-xs text-slate-500 truncate">{description}</p>
	</div>
</Button>