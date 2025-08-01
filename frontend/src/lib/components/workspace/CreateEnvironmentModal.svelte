<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { X } from 'lucide-svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	
	const dispatch = createEventDispatcher();
	
	export let show = false;
	export let environmentName = '';
	export let isCreating = false;
	
	let nameInput: HTMLInputElement;

	$: if (show && nameInput) {
		setTimeout(() => {
			nameInput.focus();
		}, 100);
	}

	function handleSubmit() {
		if (!environmentName.trim()) return;
		
		dispatch('create', {
			name: environmentName.trim()
		});
	}

	function handleClose() {
		dispatch('close');
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			handleClose();
		} else if (event.key === 'Enter') {
			handleSubmit();
		}
	}
</script>

{#if show}
	<div class="fixed inset-0 z-50 flex items-center justify-center">
		<!-- Backdrop -->
		<div class="absolute inset-0 bg-black/50" on:click={handleClose}></div>
		
		<!-- Modal -->
		<div class="relative bg-gray-800 border border-gray-700 rounded-lg shadow-xl w-full max-w-md mx-4" on:keydown={handleKeyDown}>
			<!-- Header -->
			<div class="flex items-center justify-between p-6 border-b border-gray-700">
				<h3 class="text-lg font-medium text-white">Create Environment</h3>
				<IconButton 
					on:click={handleClose}
					variant="ghost"
					size="sm"
					disabled={isCreating}
				>
					<X class="w-4 h-4" />
				</IconButton>
			</div>
			
			<!-- Content -->
			<div class="p-6">
				<div class="space-y-4">
					<div>
						<label for="environment-name" class="block text-sm font-medium text-gray-300 mb-2">
							Environment Name
						</label>
						<input
							id="environment-name"
							bind:this={nameInput}
							bind:value={environmentName}
							type="text"
							placeholder="e.g. Production, Development, Staging"
							class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
							disabled={isCreating}
						/>
					</div>
					
					<div class="bg-gray-700/50 rounded-md p-3">
						<p class="text-sm text-gray-400">
							Environments allow you to store sets of variables that you can reuse across your requests. 
							You can switch between environments to quickly change the values used in your API calls.
						</p>
					</div>
				</div>
			</div>
			
			<!-- Footer -->
			<div class="flex items-center justify-end space-x-3 p-6 border-t border-gray-700">
				<Button 
					on:click={handleClose}
					variant="ghost"
					size="sm"
					disabled={isCreating}
				>
					Cancel
				</Button>
				<Button 
					on:click={handleSubmit}
					disabled={!environmentName.trim() || isCreating}
					variant="primary"
					size="sm"
					loading={isCreating}
				>
					{isCreating ? 'Creating...' : 'Create Environment'}
				</Button>
			</div>
		</div>
	</div>
{/if}