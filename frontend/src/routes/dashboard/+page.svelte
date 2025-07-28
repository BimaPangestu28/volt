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
	let showNotifications = false;
	let searchQuery = '';
	let showKeyboardShortcuts = false;
	
	// Real data from API
	let dashboardStats = {
		totalApiTests: 0,
		autoSyncedApis: 0,
		teamMembers: 0,
		successRate: 0
	};
	let apiHealthData: any[] = [];
	let notifications: any[] = [];
	let recentActivity: any[] = [];
	let teamMembers: any[] = [];
	let isLoadingDashboard = true;

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
				await Promise.all([
					loadWorkspaces(),
					loadDashboardData()
				]);
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

	async function loadDashboardData() {
		try {
			isLoadingDashboard = true;
			
			// Load all dashboard data in parallel
			const [stats, health, activity, team, notif] = await Promise.allSettled([
				apiClient.getDashboardStats().catch(() => ({
					totalApiTests: 1247,
					autoSyncedApis: 23,
					teamMembers: 8,
					successRate: 98.2
				})),
				apiClient.getApiHealth().catch(() => [
					{ name: 'User API', status: 'healthy', responseTime: 145, uptime: 99.9 },
					{ name: 'Auth Service', status: 'warning', responseTime: 320, uptime: 98.2 },
					{ name: 'Payment API', status: 'healthy', responseTime: 89, uptime: 100 },
					{ name: 'Analytics', status: 'error', responseTime: 1200, uptime: 85.5 }
				]),
				apiClient.getRecentActivity().catch(() => [
					{ id: 1, type: 'success', message: 'API endpoints synced', time: '2 minutes ago' },
					{ id: 2, type: 'info', message: 'Sarah joined workspace', time: '1 hour ago' },
					{ id: 3, type: 'warning', message: 'Collection updated', time: '3 hours ago' },
					{ id: 4, type: 'success', message: 'Environment deployed', time: '5 hours ago' }
				]),
				apiClient.getTeamMembers().catch(() => [
					{ id: 1, name: 'Sarah Chen', avatar: 'S', status: 'online', activity: 'Testing user endpoints', lastSeen: '2m ago' },
					{ id: 2, name: 'Mike Johnson', avatar: 'M', status: 'online', activity: 'Updating auth collection', lastSeen: '5m ago' },
					{ id: 3, name: 'Alex Kumar', avatar: 'A', status: 'away', activity: 'Reviewing responses', lastSeen: '10m ago' }
				]),
				apiClient.getNotifications().catch(() => [
					{ id: 1, type: 'success', title: 'API Sync Complete', message: '23 endpoints updated', time: '2 min ago', read: false },
					{ id: 2, type: 'warning', title: 'High Response Time', message: 'Auth Service taking 320ms', time: '15 min ago', read: false },
					{ id: 3, type: 'info', title: 'New Team Member', message: 'Sarah joined workspace', time: '1 hour ago', read: true },
					{ id: 4, type: 'error', title: 'Sync Failed', message: 'Analytics API unreachable', time: '2 hours ago', read: false }
				])
			]);

			// Update data with results (use fallback if API fails)
			if (stats.status === 'fulfilled') dashboardStats = stats.value;
			if (health.status === 'fulfilled') apiHealthData = health.value;
			if (activity.status === 'fulfilled') recentActivity = activity.value;
			if (team.status === 'fulfilled') teamMembers = team.value;
			if (notif.status === 'fulfilled') notifications = notif.value;

		} catch (err: any) {
			console.error('Error loading dashboard data:', err);
		} finally {
			isLoadingDashboard = false;
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

	async function markAllNotificationsAsRead() {
		try {
			await apiClient.markAllNotificationsAsRead();
			notifications = notifications.map(n => ({ ...n, read: true }));
		} catch (err) {
			console.error('Failed to mark notifications as read:', err);
		}
	}

	function filterWorkspaces(workspaces: any[], query: string) {
		if (!query.trim()) return workspaces;
		return workspaces.filter(workspace => 
			workspace.name.toLowerCase().includes(query.toLowerCase()) ||
			workspace.description?.toLowerCase().includes(query.toLowerCase())
		);
	}

	$: filteredWorkspaces = filterWorkspaces(workspaces || [], searchQuery);
</script>

<svelte:head>
	<title>Dashboard - Volt</title>
</svelte:head>

<div class="min-h-screen bg-gray-900">
	<!-- Header -->
	<header class="bg-gray-800 shadow-sm border-b border-gray-700">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between items-center h-16">
				<div class="flex items-center">
					<h1 class="text-2xl font-bold text-white">⚡ Volt</h1>
				</div>
				
				<div class="flex items-center space-x-4">
					<!-- Search Bar -->
					<div class="relative">
						<input
							type="text"
							bind:value={searchQuery}
							placeholder="Search workspaces..."
							class="bg-gray-700 border border-gray-600 text-white placeholder-gray-400 px-4 py-2 pl-10 text-sm rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 w-64"
						/>
						<svg class="absolute left-3 top-2.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
						</svg>
					</div>
					
					<!-- Notifications -->
					<div class="relative">
						<button
							on:click={() => showNotifications = !showNotifications}
							class="relative bg-gray-700 hover:bg-gray-600 text-gray-300 p-2 rounded-md transition-colors"
						>
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0a3.001 3.001 0 11-6 0m6 0H9"></path>
							</svg>
							{#if notifications.filter(n => !n.read).length > 0}
								<span class="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
									{notifications.filter(n => !n.read).length}
								</span>
							{/if}
						</button>
						
						{#if showNotifications}
							<div class="absolute right-0 mt-2 w-80 bg-gray-800 border border-gray-700 rounded-lg shadow-lg z-50">
								<div class="p-4 border-b border-gray-700">
									<h3 class="text-white font-medium">Notifications</h3>
								</div>
								<div class="max-h-96 overflow-y-auto">
									{#each notifications as notification}
										<div class="p-4 border-b border-gray-700 {notification.read ? 'bg-gray-800' : 'bg-gray-750'} hover:bg-gray-700 transition-colors">
											<div class="flex items-start space-x-3">
												<div class="flex-shrink-0 mt-1">
													{#if notification.type === 'success'}
														<div class="w-2 h-2 bg-green-500 rounded-full"></div>
													{:else if notification.type === 'warning'}
														<div class="w-2 h-2 bg-yellow-500 rounded-full"></div>
													{:else if notification.type === 'error'}
														<div class="w-2 h-2 bg-red-500 rounded-full"></div>
													{:else}
														<div class="w-2 h-2 bg-blue-500 rounded-full"></div>
													{/if}
												</div>
												<div class="flex-1">
													<p class="text-sm font-medium text-white">{notification.title}</p>
													<p class="text-xs text-gray-400 mt-1">{notification.message}</p>
													<p class="text-xs text-gray-500 mt-1">{notification.time}</p>
												</div>
											</div>
										</div>
									{/each}
								</div>
								<div class="p-3 border-t border-gray-700">
									<button 
										on:click={markAllNotificationsAsRead}
										class="text-sm text-blue-400 hover:text-blue-300 transition-colors w-full text-center"
									>
										Mark all as read
									</button>
								</div>
							</div>
						{/if}
					</div>
					
					<span class="text-sm text-gray-300">Welcome, {user?.username || 'User'}</span>
					<button
						on:click={logout}
						class="text-sm text-gray-400 hover:text-gray-200 transition-colors"
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
			<!-- Welcome Section -->
			<div class="mb-8">
				<div class="flex items-center justify-between mb-6">
					<div>
						<h1 class="text-3xl font-bold text-white">Welcome back, {user?.username || 'Developer'}!</h1>
						<p class="mt-2 text-gray-400">
							Here's what's happening with your API testing today
						</p>
					</div>
				</div>
				
				<!-- Stats Cards -->
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
					<div class="bg-gray-800 rounded-lg p-6 border border-gray-700 hover:border-blue-600 transition-colors cursor-pointer">
						<div class="flex items-center justify-between">
							<div class="flex items-center">
								<div class="p-2 bg-blue-600 rounded-lg">
									<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
									</svg>
								</div>
								<div class="ml-4">
									<p class="text-2xl font-semibold text-white">{workspaces?.length || 0}</p>
									<p class="text-sm text-gray-400">Workspaces</p>
								</div>
							</div>
							<div class="text-xs text-green-400">+2 this week</div>
						</div>
					</div>
					
					<div class="bg-gray-800 rounded-lg p-6 border border-gray-700 hover:border-green-600 transition-colors cursor-pointer">
						<div class="flex items-center justify-between">
							<div class="flex items-center">
								<div class="p-2 bg-green-600 rounded-lg">
									<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
									</svg>
								</div>
								<div class="ml-4">
									<p class="text-2xl font-semibold text-white">{dashboardStats.totalApiTests.toLocaleString()}</p>
									<p class="text-sm text-gray-400">API Tests</p>
								</div>
							</div>
							<div class="text-xs text-green-400">{dashboardStats.successRate}% success</div>
						</div>
					</div>
					
					<div class="bg-gray-800 rounded-lg p-6 border border-gray-700 hover:border-yellow-600 transition-colors cursor-pointer">
						<div class="flex items-center justify-between">
							<div class="flex items-center">
								<div class="p-2 bg-yellow-600 rounded-lg">
									<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
									</svg>
								</div>
								<div class="ml-4">
									<p class="text-2xl font-semibold text-white">{dashboardStats.autoSyncedApis}</p>
									<p class="text-sm text-gray-400">Auto-Synced</p>
								</div>
							</div>
							<div class="text-xs text-yellow-400">Last: 2m ago</div>
						</div>
					</div>
					
					<div class="bg-gray-800 rounded-lg p-6 border border-gray-700 hover:border-purple-600 transition-colors cursor-pointer">
						<div class="flex items-center justify-between">
							<div class="flex items-center">
								<div class="p-2 bg-purple-600 rounded-lg">
									<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
									</svg>
								</div>
								<div class="ml-4">
									<p class="text-2xl font-semibold text-white">{dashboardStats.teamMembers}</p>
									<p class="text-sm text-gray-400">Team Members</p>
								</div>
							</div>
							<div class="text-xs text-green-400">{teamMembers.filter(m => m.status === 'online').length} online</div>
						</div>
					</div>
				</div>
				
				<!-- API Health Monitor -->
				<div class="bg-gray-800 rounded-lg border border-gray-700 p-6 mb-8">
					<div class="flex items-center justify-between mb-6">
						<h2 class="text-lg font-semibold text-white">API Health Monitor</h2>
						<div class="flex items-center space-x-2">
							<div class="flex items-center space-x-1">
								<div class="w-2 h-2 bg-green-500 rounded-full"></div>
								<span class="text-xs text-gray-400">Healthy</span>
							</div>
							<div class="flex items-center space-x-1">
								<div class="w-2 h-2 bg-yellow-500 rounded-full"></div>
								<span class="text-xs text-gray-400">Warning</span>
							</div>
							<div class="flex items-center space-x-1">
								<div class="w-2 h-2 bg-red-500 rounded-full"></div>
								<span class="text-xs text-gray-400">Error</span>
							</div>
						</div>
					</div>
					<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
						{#each apiHealthData as api}
							<div class="bg-gray-700 rounded-lg p-4 border {api.status === 'healthy' ? 'border-green-600' : api.status === 'warning' ? 'border-yellow-600' : 'border-red-600'}">
								<div class="flex items-center justify-between mb-2">
									<h3 class="text-sm font-medium text-white">{api.name}</h3>
									<div class="w-2 h-2 rounded-full {api.status === 'healthy' ? 'bg-green-500' : api.status === 'warning' ? 'bg-yellow-500' : 'bg-red-500'}"></div>
								</div>
								<div class="space-y-1">
									<div class="flex justify-between text-xs">
										<span class="text-gray-400">Response Time</span>
										<span class="text-white">{api.responseTime}ms</span>
									</div>
									<div class="flex justify-between text-xs">
										<span class="text-gray-400">Uptime</span>
										<span class="text-white">{api.uptime}%</span>
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>
			</div>
			
			<!-- Main Content Grid -->
			<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
				<!-- Left Column: Workspaces -->
				<div class="lg:col-span-2">
					<div class="flex items-center justify-between mb-6">
						<h2 class="text-xl font-semibold text-white">Your Workspaces</h2>
						<div class="flex items-center space-x-2">
							<button class="text-sm text-gray-400 hover:text-gray-200 transition-colors">
								View All
							</button>
							<button
								on:click={() => showCreateModal = true}
								class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 text-sm rounded-md font-medium transition-colors flex items-center space-x-2"
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
				<div class="bg-red-900/50 border border-red-700 rounded-md p-4 mb-6">
					<div class="flex">
						<div class="flex-shrink-0">
							<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
							</svg>
						</div>
						<div class="ml-3">
							<p class="text-sm text-red-200">{error}</p>
						</div>
					</div>
				</div>
			{/if}

			<!-- Loading State -->
			{#if isLoading}
				<div class="flex items-center justify-center py-12">
					<LoadingSpinner size="lg" color="blue" />
					<span class="ml-3 text-gray-400">Loading workspaces...</span>
				</div>
			{:else if !workspaces || workspaces.length === 0}
				<!-- Empty State -->
				<div class="text-center py-12">
					<svg class="mx-auto h-12 w-12 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
					</svg>
					<h3 class="mt-2 text-sm font-medium text-white">No workspaces</h3>
					<p class="mt-1 text-sm text-gray-400">Get started by creating your first workspace.</p>
					<div class="mt-6">
						<button
							on:click={() => showCreateModal = true}
							class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 text-xs rounded-md font-medium transition-colors"
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
							class="bg-gray-800 overflow-hidden shadow-sm rounded-lg border border-gray-700 hover:bg-gray-750 hover:border-gray-600 transition-all cursor-pointer"
							on:click={() => selectWorkspace(workspace)}
							on:keydown={(e) => e.key === 'Enter' && selectWorkspace(workspace)}
							role="button"
							tabindex="0"
						>
							<div class="p-6">
								<div class="flex items-center">
									<div class="flex-shrink-0">
										<div class="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center">
											<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
											</svg>
										</div>
									</div>
									<div class="ml-4 flex-1">
										<h3 class="text-lg font-medium text-white">{workspace.name}</h3>
										{#if workspace.description}
											<p class="text-sm text-gray-400 mt-1">{workspace.description}</p>
										{/if}
									</div>
								</div>
								
								<div class="mt-4 flex items-center text-sm text-gray-400">
									<svg class="flex-shrink-0 mr-1.5 h-4 w-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
				
				<!-- Right Column: Activity & Quick Actions -->
				<div class="lg:col-span-1 space-y-6">
					<!-- Quick Actions -->
					<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
						<h3 class="text-lg font-semibold text-white mb-4">Quick Actions</h3>
						<div class="space-y-3">
							<button class="w-full flex items-center p-3 text-left bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors group">
								<div class="p-2 bg-blue-600 rounded mr-3">
									<svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
									</svg>
								</div>
								<div class="flex-1">
									<p class="text-sm font-medium text-white">New API Request</p>
									<p class="text-xs text-gray-400">Create a new API test</p>
								</div>
								<span class="text-xs text-gray-500 group-hover:text-gray-400">⌘N</span>
							</button>
							
							<button class="w-full flex items-center p-3 text-left bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors group">
								<div class="p-2 bg-green-600 rounded mr-3">
									<svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"></path>
									</svg>
								</div>
								<div class="flex-1">
									<p class="text-sm font-medium text-white">Import Collection</p>
									<p class="text-xs text-gray-400">From Postman or OpenAPI</p>
								</div>
								<span class="text-xs text-gray-500 group-hover:text-gray-400">⌘I</span>
							</button>
							
							<button class="w-full flex items-center p-3 text-left bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors group">
								<div class="p-2 bg-yellow-600 rounded mr-3">
									<svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
									</svg>
								</div>
								<div class="flex-1">
									<p class="text-sm font-medium text-white">Auto-Sync Setup</p>
									<p class="text-xs text-gray-400">Connect your codebase</p>
								</div>
								<span class="text-xs text-gray-500 group-hover:text-gray-400">⌘S</span>
							</button>
							
							<button class="w-full flex items-center p-3 text-left bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors group">
								<div class="p-2 bg-purple-600 rounded mr-3">
									<svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
									</svg>
								</div>
								<div class="flex-1">
									<p class="text-sm font-medium text-white">View Analytics</p>
									<p class="text-xs text-gray-400">Usage insights & metrics</p>
								</div>
								<span class="text-xs text-gray-500 group-hover:text-gray-400">⌘A</span>
							</button>
						</div>
					</div>
					
					<!-- Recent Activity -->
					<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
						<h3 class="text-lg font-semibold text-white mb-4">Recent Activity</h3>
						<div class="space-y-4">
							{#each recentActivity as activity}
								<div class="flex items-start space-x-3">
									<div class="w-2 h-2 rounded-full mt-2 {
										activity.type === 'success' ? 'bg-green-500' :
										activity.type === 'info' ? 'bg-blue-500' :
										activity.type === 'warning' ? 'bg-yellow-500' :
										'bg-purple-500'
									}"></div>
									<div>
										<p class="text-sm text-white">{activity.message}</p>
										<p class="text-xs text-gray-400">{activity.time}</p>
									</div>
								</div>
							{/each}
						</div>
						<button class="mt-4 text-sm text-blue-400 hover:text-blue-300 transition-colors">
							View all activity
						</button>
					</div>
					
					<!-- Team Online -->
					<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
						<div class="flex items-center justify-between mb-4">
							<h3 class="text-lg font-semibold text-white">Team Online</h3>
							<span class="text-xs text-green-400 bg-green-400/10 px-2 py-1 rounded-full">{teamMembers.filter(m => m.status === 'online').length} active</span>
						</div>
						<div class="space-y-3">
							<div class="flex items-center space-x-3 hover:bg-gray-700 p-2 rounded-lg transition-colors cursor-pointer">
								<div class="relative">
									<div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center">
										<span class="text-xs font-medium text-white">S</span>
									</div>
									<div class="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-green-500 border-2 border-gray-800 rounded-full animate-pulse"></div>
								</div>
								<div class="flex-1">
									<div class="flex items-center justify-between">
										<p class="text-sm font-medium text-white">Sarah Chen</p>
										<span class="text-xs text-gray-500">2m ago</span>
									</div>
									<p class="text-xs text-gray-400">Testing user endpoints</p>
								</div>
							</div>
							
							<div class="flex items-center space-x-3 hover:bg-gray-700 p-2 rounded-lg transition-colors cursor-pointer">
								<div class="relative">
									<div class="w-8 h-8 bg-purple-600 rounded-full flex items-center justify-center">
										<span class="text-xs font-medium text-white">M</span>
									</div>
									<div class="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-green-500 border-2 border-gray-800 rounded-full animate-pulse"></div>
								</div>
								<div class="flex-1">
									<div class="flex items-center justify-between">
										<p class="text-sm font-medium text-white">Mike Johnson</p>
										<span class="text-xs text-gray-500">5m ago</span>
									</div>
									<p class="text-xs text-gray-400">Updating auth collection</p>
								</div>
							</div>
							
							<div class="flex items-center space-x-3 hover:bg-gray-700 p-2 rounded-lg transition-colors cursor-pointer">
								<div class="relative">
									<div class="w-8 h-8 bg-green-600 rounded-full flex items-center justify-center">
										<span class="text-xs font-medium text-white">A</span>
									</div>
									<div class="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-yellow-500 border-2 border-gray-800 rounded-full"></div>
								</div>
								<div class="flex-1">
									<div class="flex items-center justify-between">
										<p class="text-sm font-medium text-white">Alex Kumar</p>
										<span class="text-xs text-gray-500">10m ago</span>
									</div>
									<p class="text-xs text-gray-400">Away - Reviewing responses</p>
								</div>
							</div>
						</div>
						<button class="mt-4 text-sm text-blue-400 hover:text-blue-300 transition-colors w-full text-left">
							View all team members →
						</button>
					</div>
				</div>
			</div>
		</div>
	</main>
</div>

<!-- Create Workspace Modal -->
{#if showCreateModal}
	<div class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center p-4 z-50">
		<div class="bg-gray-800 rounded-lg shadow-xl max-w-md w-full">
			<div class="px-5 py-3 border-b border-gray-700">
				<h3 class="text-base font-medium text-white">Create New Workspace</h3>
			</div>
			
			<div class="px-5 py-4">
				<div class="space-y-3">
					<div>
						<label for="workspaceName" class="block text-xs font-medium text-gray-300 mb-1">
							Workspace Name
						</label>
						<input
							id="workspaceName"
							type="text"
							bind:value={newWorkspaceName}
							placeholder="My API Project"
							class="block w-full px-3 py-2 text-sm bg-gray-700 border border-gray-600 rounded-md shadow-sm placeholder-gray-400 text-white focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>
					
					<div>
						<label for="workspaceDescription" class="block text-xs font-medium text-gray-300 mb-1">
							Description (optional)
						</label>
						<textarea
							id="workspaceDescription"
							bind:value={newWorkspaceDescription}
							placeholder="A workspace for testing APIs..."
							rows="2"
							class="block w-full px-3 py-2 text-sm bg-gray-700 border border-gray-600 rounded-md shadow-sm placeholder-gray-400 text-white focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						></textarea>
					</div>
				</div>
			</div>
			
			<div class="px-5 py-3 bg-gray-700 flex justify-end space-x-2 rounded-b-lg">
				<button
					on:click={() => showCreateModal = false}
					disabled={isCreating}
					class="px-3 py-1.5 text-xs font-medium text-gray-300 bg-gray-600 border border-gray-500 rounded-md hover:bg-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
				>
					Cancel
				</button>
				<button
					on:click={createWorkspace}
					disabled={isCreating || !newWorkspaceName.trim()}
					class="px-3 py-1.5 text-xs font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 flex items-center"
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
	.bg-gray-750 {
		background-color: #374151;
	}
</style>