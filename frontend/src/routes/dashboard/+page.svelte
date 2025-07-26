<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { workspaceStore } from '$lib/stores/workspace';
	import { apiClient } from '$lib/api/client';
	import { toastStore } from '$lib/stores/toast';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

	let isAuthenticated = false;
	let user = null;
	let workspaces: any[] = []; // Ensure it's always an array
	let isLoading = true;
	let error = '';
	let showCreateModal = false;
	let newWorkspaceName = '';
	let newWorkspaceDescription = '';
	let isCreating = false;

	// Subscribe to stores
	authStore.subscribe(state => {
		isAuthenticated = state.isAuthenticated;
		user = state.user;
	});

	workspaceStore.subscribe(state => {
		workspaces = state.workspaces || []; // Ensure it's never null/undefined
		isLoading = state.isLoading;
		error = state.error || '';
	});

	onMount(async () => {
		// Wait a bit for auth store to initialize
		const unsubscribe = authStore.subscribe(async (state) => {
			if (state.isAuthenticated === false && !state.user) {
				goto('/login');
				return;
			}
			
			if (state.isAuthenticated && state.user) {
				await loadWorkspaces();
				unsubscribe();
			}
		});
	});

	async function loadWorkspaces() {
		try {
			console.log('Starting to load workspaces...');
			workspaceStore.setLoading(true);
			workspaceStore.setError(null);
			
			const data = await apiClient.getWorkspaces();
			console.log('Workspaces received:', data);
			workspaceStore.setWorkspaces(data);
		} catch (err: any) {
			console.error('Error loading workspaces:', err);
			workspaceStore.setError(err.response?.data?.error || 'Failed to load workspaces');
		} finally {
			workspaceStore.setLoading(false);
			console.log('Loading finished');
		}
	}

	async function createWorkspace() {
		if (!newWorkspaceName.trim()) {
			return;
		}

		isCreating = true;
		try {
			const workspace = await apiClient.createWorkspace(newWorkspaceName, newWorkspaceDescription);
			
			// Refetch workspaces from server to ensure data consistency
			await loadWorkspaces();
			
			// Show success toast
			toastStore.success(`Workspace "${newWorkspaceName}" created successfully!`);
			
			// Reset form
			newWorkspaceName = '';
			newWorkspaceDescription = '';
			showCreateModal = false;
		} catch (err: any) {
			const errorMessage = err.response?.data?.error || 'Failed to create workspace';
			error = errorMessage;
			toastStore.error(errorMessage);
		} finally {
			isCreating = false;
		}
	}

	function selectWorkspace(workspace: any) {
		workspaceStore.setCurrentWorkspace(workspace);
		goto(`/workspace/${workspace.id}`);
	}

	function logout() {
		apiClient.logout();
		goto('/');
	}
</script>

<svelte:head>
	<title>Dashboard - Volt</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<!-- Header -->
	<header class="bg-white shadow-sm border-b border-gray-200">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between items-center h-16">
				<div class="flex items-center">
					<h1 class="text-2xl font-bold electric-gradient bg-clip-text text-transparent">⚡ Volt</h1>
				</div>
				
				<div class="flex items-center space-x-4">
					<span class="text-sm text-gray-600">Welcome, {user?.username || 'User'}</span>
					<button
						on:click={logout}
						class="text-sm text-gray-500 hover:text-gray-700 transition-colors"
					>
						Logout
					</button>
				</div>
			</div>
		</div>
	</header>

	<!-- Main Content -->
	<main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
		<div class="px-4 py-6 sm:px-0">
			<!-- Page Header -->
			<div class="border-b border-gray-200 pb-5 mb-6">
				<div class="flex items-center justify-between">
					<div>
						<h1 class="text-2xl font-bold text-gray-900">Your Workspaces</h1>
						<p class="mt-1 text-sm text-gray-600">
							Select a workspace to start testing APIs or create a new one
						</p>
					</div>
					<button
						on:click={() => showCreateModal = true}
						class="electric-gradient text-white px-4 py-2 text-xs rounded-md font-medium hover:opacity-90 transition-opacity flex items-center space-x-1.5"
					>
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
						</svg>
						<span>New Workspace</span>
					</button>
				</div>
			</div>

			<!-- Error Message -->
			{#if error}
				<div class="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
					<div class="flex">
						<div class="flex-shrink-0">
							<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
							</svg>
						</div>
						<div class="ml-3">
							<p class="text-sm text-red-800">{error}</p>
						</div>
					</div>
				</div>
			{/if}

			<!-- Loading State -->
			{#if isLoading}
				<div class="flex items-center justify-center py-12">
					<LoadingSpinner size="lg" color="blue" />
					<span class="ml-3 text-gray-600">Loading workspaces...</span>
				</div>
			{:else if !workspaces || workspaces.length === 0}
				<!-- Empty State -->
				<div class="text-center py-12">
					<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
					</svg>
					<h3 class="mt-2 text-sm font-medium text-gray-900">No workspaces</h3>
					<p class="mt-1 text-sm text-gray-500">Get started by creating your first workspace.</p>
					<div class="mt-6">
						<button
							on:click={() => showCreateModal = true}
							class="electric-gradient text-white px-4 py-2 text-xs rounded-md font-medium hover:opacity-90 transition-opacity"
						>
							Create Workspace
						</button>
					</div>
				</div>
			{:else}
				<!-- Workspaces Grid -->
				<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
					{#each workspaces || [] as workspace}
						<div 
							class="bg-white overflow-hidden shadow-sm rounded-lg border border-gray-200 hover:shadow-md transition-shadow cursor-pointer"
							on:click={() => selectWorkspace(workspace)}
							on:keydown={(e) => e.key === 'Enter' && selectWorkspace(workspace)}
							role="button"
							tabindex="0"
						>
							<div class="p-6">
								<div class="flex items-center">
									<div class="flex-shrink-0">
										<div class="w-10 h-10 electric-gradient rounded-lg flex items-center justify-center">
											<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
											</svg>
										</div>
									</div>
									<div class="ml-4 flex-1">
										<h3 class="text-lg font-medium text-gray-900">{workspace.name}</h3>
										{#if workspace.description}
											<p class="text-sm text-gray-500 mt-1">{workspace.description}</p>
										{/if}
									</div>
								</div>
								
								<div class="mt-4 flex items-center text-sm text-gray-500">
									<svg class="flex-shrink-0 mr-1.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"></path>
									</svg>
									{workspace.members?.length || 1} member{workspace.members?.length !== 1 ? 's' : ''}
								</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</main>
</div>

<!-- Create Workspace Modal -->
{#if showCreateModal}
	<div class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-lg shadow-xl max-w-md w-full">
			<div class="px-5 py-3 border-b border-gray-200">
				<h3 class="text-base font-medium text-gray-900">Create New Workspace</h3>
			</div>
			
			<div class="px-5 py-4">
				<div class="space-y-3">
					<div>
						<label for="workspaceName" class="block text-xs font-medium text-gray-700 mb-1">
							Workspace Name
						</label>
						<input
							id="workspaceName"
							type="text"
							bind:value={newWorkspaceName}
							placeholder="My API Project"
							class="block w-full px-3 py-2 text-sm border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>
					
					<div>
						<label for="workspaceDescription" class="block text-xs font-medium text-gray-700 mb-1">
							Description (optional)
						</label>
						<textarea
							id="workspaceDescription"
							bind:value={newWorkspaceDescription}
							placeholder="A workspace for testing APIs..."
							rows="2"
							class="block w-full px-3 py-2 text-sm border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						></textarea>
					</div>
				</div>
			</div>
			
			<div class="px-5 py-3 bg-gray-50 flex justify-end space-x-2 rounded-b-lg">
				<button
					on:click={() => showCreateModal = false}
					disabled={isCreating}
					class="px-3 py-1.5 text-xs font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
				>
					Cancel
				</button>
				<button
					on:click={createWorkspace}
					disabled={isCreating || !newWorkspaceName.trim()}
					class="px-3 py-1.5 text-xs font-medium text-white electric-gradient rounded-md hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 flex items-center"
				>
					{#if isCreating}
						<LoadingSpinner size="sm" color="white" />
						<span class="ml-1.5">Creating...</span>
					{:else}
						Create Workspace
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.electric-gradient {
		background: linear-gradient(135deg, var(--electric-blue), var(--deep-blue));
	}
	
	.bg-clip-text {
		background-clip: text;
		-webkit-background-clip: text;
	}
</style>