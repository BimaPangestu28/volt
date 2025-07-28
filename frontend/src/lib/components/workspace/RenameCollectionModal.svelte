<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	
	export let show = false;
	export let collectionName = '';
	export let isRenaming = false;
	
	const dispatch = createEventDispatcher();
	
	let newName = '';
	let nameInput: HTMLInputElement;
	let initialized = false;
	
	$: if (show && collectionName && !initialized) {
		newName = collectionName;
		initialized = true;
		// Focus input after modal opens
		setTimeout(() => {
			if (nameInput) {
				nameInput.focus();
				nameInput.select();
			}
		}, 100);
	}
	
	$: if (!show) {
		initialized = false;
	}
	
	function handleSubmit() {
		const trimmedName = newName.trim();
		if (!trimmedName) return;
		
		if (trimmedName === collectionName) {
			// No changes, just close
			handleClose();
			return;
		}
		
		dispatch('rename', {
			newName: trimmedName
		});
	}
	
	function handleClose() {
		show = false;
		newName = '';
		initialized = false;
		dispatch('close');
	}
	
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			handleClose();
		} else if (event.key === 'Enter') {
			event.preventDefault();
			handleSubmit();
		}
	}
	
	function handleBackdropClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			handleClose();
		}
	}
</script>

<!-- Modal backdrop -->
{#if show}
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div 
		class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
		on:click={handleBackdropClick}
		on:keydown={handleKeydown}
	>
		<!-- Modal content -->
		<div class="bg-gray-800 rounded-lg border border-gray-700 shadow-xl w-full max-w-md mx-4">
			<!-- Header -->
			<div class="px-6 py-4 border-b border-gray-700">
				<h3 class="text-lg font-semibold text-white">Rename Collection</h3>
			</div>
			
			<!-- Body -->
			<div class="px-6 py-4">
				<div class="space-y-4">
					<div>
						<label for="collection-name" class="block text-sm font-medium text-gray-300 mb-2">
							Collection name
						</label>
						<input
							id="collection-name"
							bind:this={nameInput}
							bind:value={newName}
							type="text"
							placeholder="Enter collection name"
							disabled={isRenaming}
							class="w-full px-3 py-2 border border-gray-600 rounded-md text-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
							on:keydown={handleKeydown}
						/>
					</div>
				</div>
			</div>
			
			<!-- Footer -->
			<div class="px-6 py-4 border-t border-gray-700 flex justify-end space-x-3">
				<button
					type="button"
					on:click={handleClose}
					disabled={isRenaming}
					class="px-4 py-2 text-sm font-medium text-gray-300 border border-gray-600 rounded-md hover:bg-gray-700 hover:text-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
				>
					Cancel
				</button>
				<button
					type="button"
					on:click={handleSubmit}
					disabled={isRenaming || !newName.trim() || newName.trim() === collectionName}
					class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
				>
					{#if isRenaming}
						<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Renaming...
					{:else}
						Rename
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}