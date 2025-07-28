<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { Search, Plus, MoreHorizontal, Copy, Trash2, Edit } from 'lucide-svelte';
	
	const dispatch = createEventDispatcher();

	export let environments: any[] = [];
	export let selectedEnvironment: any = null;
	export let isLoading = false;

	let searchQuery = '';
	let showOptionsMenu: { [key: string]: boolean } = {};

	// Filter environments based on search
	$: filteredEnvironments = environments.filter(env => 
		env.name.toLowerCase().includes(searchQuery.toLowerCase())
	);

	function selectEnvironment(environment: any) {
		dispatch('selectEnvironment', environment);
	}

	function createNewEnvironment() {
		dispatch('createEnvironment');
	}

	function toggleOptionsMenu(envId: string, event: Event) {
		event.stopPropagation();
		showOptionsMenu = {
			...showOptionsMenu,
			[envId]: !showOptionsMenu[envId]
		};
	}

	function handleRename(environment: any, event: Event) {
		event.stopPropagation();
		showOptionsMenu = {};
		dispatch('renameEnvironment', environment);
	}

	function handleDuplicate(environment: any, event: Event) {
		event.stopPropagation();
		showOptionsMenu = {};
		dispatch('duplicateEnvironment', environment);
	}

	function handleDelete(environment: any, event: Event) {
		event.stopPropagation();
		showOptionsMenu = {};
		dispatch('deleteEnvironment', environment);
	}

	// Close options menu when clicking outside
	function handleClickOutside(event: Event) {
		const target = event.target as Element;
		if (!target.closest('.options-menu')) {
			showOptionsMenu = {};
		}
	}
</script>

<svelte:window on:click={handleClickOutside} />

<div class="w-80 h-full bg-gray-900 border-r border-gray-700 flex flex-col">
	<!-- Header -->
	<div class="p-4 border-b border-gray-700">
		<div class="flex items-center justify-between mb-3">
			<h2 class="text-sm font-medium text-white">Environments</h2>
			<button
				on:click={createNewEnvironment}
				class="p-1.5 text-gray-400 hover:text-white hover:bg-gray-800 rounded transition-colors"
				title="Create new environment"
			>
				<Plus class="w-4 h-4" />
			</button>
		</div>

		<!-- Search -->
		<div class="relative">
			<Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-500" />
			<input
				bind:value={searchQuery}
				type="text"
				placeholder="Search environments"
				class="w-full pl-9 pr-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-sm text-white placeholder-gray-500 focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
			/>
		</div>
	</div>

	<!-- Environments List -->
	<div class="flex-1 overflow-y-auto">
		{#if isLoading}
			<!-- Loading skeleton -->
			<div class="p-2 space-y-1">
				{#each Array(5) as _}
					<div class="h-10 bg-gray-800/50 rounded animate-pulse"></div>
				{/each}
			</div>
		{:else if filteredEnvironments.length > 0}
			<div class="p-2 space-y-1">
				{#each filteredEnvironments as environment (environment.id)}
					<div
						class="group relative flex items-center justify-between p-3 rounded-md cursor-pointer transition-colors {selectedEnvironment?.id === environment.id 
							? 'bg-blue-600 text-white' 
							: 'hover:bg-gray-800 text-gray-300'}"
						on:click={() => selectEnvironment(environment)}
						role="button"
						tabindex="0"
					>
						<div class="flex-1 min-w-0">
							<div class="text-sm font-medium truncate">
								{environment.name}
							</div>
							<div class="text-xs {selectedEnvironment?.id === environment.id 
								? 'text-blue-100' 
								: 'text-gray-500'} mt-0.5">
								{environment.variables?.length || 0} variables
							</div>
						</div>

						<!-- Options Menu -->
						<div class="options-menu relative opacity-0 group-hover:opacity-100 transition-opacity">
							<button
								on:click={(e) => toggleOptionsMenu(environment.id, e)}
								class="p-1.5 {selectedEnvironment?.id === environment.id 
									? 'text-blue-100 hover:text-white hover:bg-blue-700' 
									: 'text-gray-400 hover:text-white hover:bg-gray-700'} rounded transition-colors"
								title="Environment options"
							>
								<MoreHorizontal class="w-4 h-4" />
							</button>

							{#if showOptionsMenu[environment.id]}
								<div class="absolute right-0 top-full mt-1 w-40 bg-gray-800 border border-gray-700 rounded-md shadow-lg z-10">
									<div class="py-1">
										<button
											on:click={(e) => handleRename(environment, e)}
											class="flex items-center space-x-2 w-full px-3 py-2 text-sm text-gray-300 hover:bg-gray-700 hover:text-white"
										>
											<Edit class="w-3.5 h-3.5" />
											<span>Rename</span>
										</button>
										<button
											on:click={(e) => handleDuplicate(environment, e)}
											class="flex items-center space-x-2 w-full px-3 py-2 text-sm text-gray-300 hover:bg-gray-700 hover:text-white"
										>
											<Copy class="w-3.5 h-3.5" />
											<span>Duplicate</span>
										</button>
										<hr class="border-gray-700 my-1" />
										<button
											on:click={(e) => handleDelete(environment, e)}
											class="flex items-center space-x-2 w-full px-3 py-2 text-sm text-red-400 hover:bg-gray-700 hover:text-red-300"
										>
											<Trash2 class="w-3.5 h-3.5" />
											<span>Delete</span>
										</button>
									</div>
								</div>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		{:else}
			<!-- Empty state -->
			<div class="flex flex-col items-center justify-center h-64 px-4">
				<div class="text-gray-500 text-center">
					<div class="w-12 h-12 mx-auto mb-3 rounded-full bg-gray-800/50 flex items-center justify-center">
						<Search class="w-6 h-6" />
					</div>
					{#if searchQuery}
						<p class="text-sm mb-2">No environments found</p>
						<p class="text-xs">Try adjusting your search</p>
					{:else}
						<p class="text-sm mb-2">No environments yet</p>
						<p class="text-xs">Create your first environment to get started</p>
					{/if}
				</div>
			</div>
		{/if}
	</div>
</div>