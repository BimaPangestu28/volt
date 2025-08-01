<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import Button from '$lib/components/ui/Button.svelte';

	export let show = false;
	export let collectionName = '';
	export let collectionDescription = '';
	export let isCreating = false;

	const dispatch = createEventDispatcher();

	function handleClose() {
		if (!isCreating) {
			dispatch('close');
		}
	}

	function handleCreate() {
		if (collectionName.trim()) {
			dispatch('create', {
				name: collectionName,
				description: collectionDescription
			});
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			handleClose();
		} else if (event.key === 'Enter' && (event.metaKey || event.ctrlKey)) {
			handleCreate();
		}
	}
</script>

{#if show}
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div 
		class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center p-4 z-50"
		on:click={handleClose}
		on:keydown={handleKeydown}
	>
		<!-- svelte-ignore a11y-no-static-element-interactions -->
		<div 
			class="bg-gray-800 rounded-lg shadow-xl max-w-md w-full border border-gray-700"
			on:click|stopPropagation
		>
			<div class="px-5 py-3 border-b border-gray-700">
				<h3 class="text-base font-medium text-white">Create New Collection</h3>
			</div>
			
			<div class="px-5 py-4">
				<div class="space-y-3">
					<div>
						<label for="collectionName" class="block text-xs font-medium text-gray-300 mb-1">
							Collection Name
						</label>
						<input
							id="collectionName"
							type="text"
							bind:value={collectionName}
							placeholder="User API"
							disabled={isCreating}
							class="block w-full px-3 py-2 text-sm border border-gray-600 rounded-md shadow-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 disabled:opacity-50"
						/>
					</div>
					
					<div>
						<label for="collectionDescription" class="block text-xs font-medium text-gray-300 mb-1">
							Description (optional)
						</label>
						<textarea
							id="collectionDescription"
							bind:value={collectionDescription}
							placeholder="API endpoints for user management..."
							rows="2"
							disabled={isCreating}
							class="block w-full px-3 py-2 text-sm border border-gray-600 rounded-md shadow-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 disabled:opacity-50"
						></textarea>
					</div>
				</div>
			</div>
			
			<div class="px-5 py-3 bg-gray-900 flex justify-end space-x-2 rounded-b-lg">
				<Button
					on:click={handleClose}
					disabled={isCreating}
					variant="ghost"
					size="sm"
				>
					Cancel
				</Button>
				<Button
					on:click={handleCreate}
					disabled={isCreating || !collectionName.trim()}
					variant="primary"
					size="sm"
					loading={isCreating}
				>
					{isCreating ? 'Creating...' : 'Create Collection'}
				</Button>
			</div>
		</div>
	</div>
{/if}

<style>
	.electric-gradient {
		background: linear-gradient(135deg, var(--electric-blue), var(--deep-blue));
	}
</style>