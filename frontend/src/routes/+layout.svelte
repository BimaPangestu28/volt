<script lang="ts">
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth';
	import ToastContainer from '$lib/components/ToastContainer.svelte';
	import '../app.css';

	onMount(() => {
		authStore.initFromStorage();
	});
</script>

<main class="min-h-screen bg-gray-50">
	<slot />
</main>

<ToastContainer />