<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Button from '$lib/components/ui/Button.svelte';
	
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
				<Button
					type="button"
					on:click={handleClose}
					disabled={isRenaming}
					variant="ghost"
					size="sm"
				>
					Cancel
				</Button>
				<Button
					type="button"
					on:click={handleSubmit}
					disabled={isRenaming || !newName.trim() || newName.trim() === collectionName}
					variant="primary"
					size="sm"
					loading={isRenaming}
				>
					{isRenaming ? 'Renaming...' : 'Rename'}
				</Button>
			</div>
		</div>
	</div>
{/if}