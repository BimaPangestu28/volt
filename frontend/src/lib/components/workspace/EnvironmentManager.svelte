<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { Search, Plus, Trash2, Eye, EyeOff, Save } from 'lucide-svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	
	const dispatch = createEventDispatcher();

	export let environment: any = null;
	export let isLoading = false;
	export let isSaving = false;

	let filterQuery = '';
	let variables: any[] = [];
	let hasUnsavedChanges = false;

	// Watch for environment changes
	$: if (environment) {
		variables = environment.variables ? [...environment.variables] : [];
		hasUnsavedChanges = false;
	}

	// Filter variables based on search
	$: filteredVariables = variables.filter(variable => 
		variable.key.toLowerCase().includes(filterQuery.toLowerCase()) ||
		variable.value.toLowerCase().includes(filterQuery.toLowerCase())
	);

	function addNewVariable() {
		variables = [
			...variables,
			{
				key: '',
				value: '',
				initial_value: '',
				type: 'default',
				enabled: true
			}
		];
		hasUnsavedChanges = true;
	}

	function removeVariable(index: number) {
		variables = variables.filter((_, i) => i !== index);
		hasUnsavedChanges = true;
	}

	function updateVariable(index: number, field: string, value: any) {
		variables[index] = { ...variables[index], [field]: value };
		variables = [...variables];
		hasUnsavedChanges = true;
	}

	function toggleVariableEnabled(index: number) {
		variables[index].enabled = !variables[index].enabled;
		variables = [...variables];
		hasUnsavedChanges = true;
	}

	function toggleVariableType(index: number) {
		const newType = variables[index].type === 'default' ? 'secret' : 'default';
		variables[index].type = newType;
		variables = [...variables];
		hasUnsavedChanges = true;
	}

	function saveEnvironment() {
		if (!environment || !hasUnsavedChanges) return;
		
		dispatch('saveEnvironment', {
			...environment,
			variables: variables.filter(v => v.key.trim() !== '')
		});
	}

	// Handle keyboard shortcuts
	function handleKeyDown(event: KeyboardEvent) {
		if ((event.ctrlKey || event.metaKey) && event.key === 's') {
			event.preventDefault();
			saveEnvironment();
		}
	}
</script>

<svelte:window on:keydown={handleKeyDown} />

<div class="flex-1 flex flex-col bg-gray-900 overflow-hidden">
	{#if environment}
		<!-- Header -->
		<div class="border-b border-gray-700 p-6">
			<div class="flex items-center justify-between mb-4">
				<div>
					<h1 class="text-xl font-semibold text-white">{environment.name}</h1>
					<p class="text-sm text-gray-400 mt-1">
						Global variables for a workspace are a set of variables that are always available within the scope of that workspace. They can be viewed and edited by anyone in that workspace.
						<a href="#" class="text-blue-400 hover:text-blue-300">Learn more about globals</a>
					</p>
				</div>
				<div class="flex items-center space-x-3">
					{#if hasUnsavedChanges}
						<span class="text-sm text-yellow-400">Unsaved changes</span>
					{/if}
					<Button
						on:click={saveEnvironment}
						disabled={!hasUnsavedChanges || isSaving}
						variant="primary"
						size="sm"
						loading={isSaving}
					>
						<Save class="w-4 h-4 mr-2" />
						{isSaving ? 'Saving...' : 'Save'}
					</Button>
				</div>
			</div>

			<!-- Filter -->
			<div class="relative max-w-md">
				<Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-500" />
				<input
					bind:value={filterQuery}
					type="text"
					placeholder="Filter variables"
					class="w-full pl-9 pr-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-sm text-white placeholder-gray-500 focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
				/>
			</div>
		</div>

		<!-- Variables Table -->
		<div class="flex-1 overflow-hidden">
			<div class="h-full overflow-y-auto">
				<table class="w-full">
					<thead class="bg-gray-800 sticky top-0 z-10">
						<tr>
							<th class="w-12 px-4 py-3 text-left">
								<div class="w-4 h-4"></div>
							</th>
							<th class="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Variable</th>
							<th class="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Type</th>
							<th class="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Initial value</th>
							<th class="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Current value</th>
							<th class="w-12 px-4 py-3 text-left">
								<div class="w-4 h-4"></div>
							</th>
						</tr>
					</thead>
					<tbody class="bg-gray-900">
						{#each filteredVariables as variable, index (index)}
							<tr class="border-b border-gray-800 hover:bg-gray-800/50 transition-colors">
								<!-- Enabled checkbox -->
								<td class="px-4 py-3">
									<input
										type="checkbox"
										checked={variable.enabled}
										on:change={() => toggleVariableEnabled(index)}
										class="w-4 h-4 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500 focus:ring-2"
									/>
								</td>

								<!-- Variable name -->
								<td class="px-4 py-3">
									<input
										type="text"
										bind:value={variable.key}
										on:input={() => updateVariable(index, 'key', variable.key)}
										placeholder="Variable name"
										class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-sm text-white placeholder-gray-500 focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
									/>
								</td>

								<!-- Type -->
								<td class="px-4 py-3">
									<Button
										on:click={() => toggleVariableType(index)}
										variant={variable.type === 'secret' ? 'danger' : 'ghost'}
										size="xs"
										class="!rounded-full {variable.type === 'secret' ? '!bg-red-900/20 !border-red-800 !text-red-400 hover:!bg-red-900/30' : ''}"
									>
										{variable.type}
									</Button>
								</td>

								<!-- Initial value -->
								<td class="px-4 py-3">
									<div class="relative">
										{#if variable.type === 'secret'}
											<input
												type="password"
												bind:value={variable.initial_value}
												on:input={() => updateVariable(index, 'initial_value', variable.initial_value)}
												placeholder="Initial value"
												class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-sm text-white placeholder-gray-500 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 pr-8"
											/>
											<IconButton
												type="button"
												variant="ghost"
												size="xs"
												class="!absolute !right-2 !top-1/2 !transform !-translate-y-1/2"
											>
												<EyeOff class="w-4 h-4" />
											</IconButton>
										{:else}
											<input
												type="text"
												bind:value={variable.initial_value}
												on:input={() => updateVariable(index, 'initial_value', variable.initial_value)}
												placeholder="Initial value"
												class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-sm text-white placeholder-gray-500 focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
											/>
										{/if}
									</div>
								</td>

								<!-- Current value -->
								<td class="px-4 py-3">
									<div class="relative">
										{#if variable.type === 'secret'}
											<input
												type="password"
												bind:value={variable.value}
												on:input={() => updateVariable(index, 'value', variable.value)}
												placeholder="Current value"
												class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-sm text-white placeholder-gray-500 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 pr-8"
											/>
											<IconButton
												type="button"
												variant="ghost"
												size="xs"
												class="!absolute !right-2 !top-1/2 !transform !-translate-y-1/2"
											>
												<EyeOff class="w-4 h-4" />
											</IconButton>
										{:else}
											<input
												type="text"
												bind:value={variable.value}
												on:input={() => updateVariable(index, 'value', variable.value)}
												placeholder="Current value"
												class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-sm text-white placeholder-gray-500 focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
											/>
										{/if}
									</div>
								</td>

								<!-- Delete button -->
								<td class="px-4 py-3">
									<IconButton
										on:click={() => removeVariable(index)}
										variant="ghost"
										size="xs"
										title="Remove variable"
										class="!text-gray-400 hover:!text-red-400"
									>
										<Trash2 class="w-4 h-4" />
									</IconButton>
								</td>
							</tr>
						{/each}

						<!-- Add new variable row -->
						<tr class="border-b border-gray-800">
							<td colspan="6" class="px-4 py-3">
								<Button
									on:click={addNewVariable}
									variant="ghost"
									size="sm"
									class="!justify-start !text-gray-400 hover:!text-white"
								>
									<Plus class="w-4 h-4 mr-2" />
									Add new variable
								</Button>
							</td>
						</tr>
					</tbody>
				</table>

				{#if filteredVariables.length === 0 && filterQuery}
					<div class="flex flex-col items-center justify-center py-12">
						<Search class="w-12 h-12 text-gray-600 mb-3" />
						<p class="text-gray-400 text-sm">No variables match your filter</p>
					</div>
				{/if}
			</div>
		</div>
	{:else}
		<!-- Welcome State -->
		<div class="flex-1 flex items-center justify-center">
			<div class="text-center max-w-md mx-auto px-6">
				<div class="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-800/50 flex items-center justify-center">
					<svg class="w-8 h-8 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
					</svg>
				</div>
				<h3 class="text-lg font-medium text-white mb-2">Select an Environment</h3>
				<p class="text-gray-400">Choose an environment from the sidebar to view and edit its variables.</p>
			</div>
		</div>
	{/if}
</div>