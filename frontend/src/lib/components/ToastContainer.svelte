<script lang="ts">
	import { toastStore } from '$lib/stores/toast';
	import Toast from './Toast.svelte';
	
	$: toasts = $toastStore;
</script>

<div class="fixed top-6 right-6 z-50 space-y-3 pointer-events-none">
	{#each toasts as toast, index (toast.id)}
		<div 
			class="pointer-events-auto"
			style="transform: translateY({index * 4}px); z-index: {50 - index};"
		>
			<Toast
				type={toast.type}
				message={toast.message}
				duration={toast.duration}
				on:close={() => toastStore.remove(toast.id)}
			/>
		</div>
	{/each}
</div>