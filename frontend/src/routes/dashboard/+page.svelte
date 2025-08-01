<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { workspaceStore } from '$lib/stores/workspace';
	import { projectStore, defaultTemplates } from '$lib/stores/project';
	import type { Project } from '$lib/stores/project';
	import { apiClient } from '$lib/api/client';
	import { toastStore } from '$lib/stores/toast';
	import ProjectCard from '$lib/components/project/ProjectCard.svelte';
	import CreateProjectModal from '$lib/components/project/CreateProjectModal.svelte';
	import TaskBoard from '$lib/components/project/TaskBoard.svelte';
	import { 
		Briefcase, Code, Database, Globe, Laptop, Smartphone, Users, Zap, Building, Rocket, Target, Settings,
		Search, Bell, LogOut, Archive, ClipboardList, CheckCircle, Activity, Plus, ChevronDown,
		FolderOpen, Calendar, AlertCircle, UserPlus, MoreHorizontal
	} from 'lucide-svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import ActionCard from '$lib/components/ui/ActionCard.svelte';

	let isAuthenticated = false;
	let user = null;
	let workspaces: any[] = [];
	let projects: Project[] = [];
	let isLoading = true;
	let error = '';
	let showCreateProjectModal = false;
	let showCreateWorkspaceModal = false;
	let showEditWorkspaceModal = false;
	let selectedWorkspace = null;
	let newWorkspaceName = '';
	let newWorkspaceDescription = '';
	let workspaceIcon = 'default';
	let workspaceImage = null;
	let uploadedImageUrl = '';
	let isCreating = false;
	let isUpdating = false;
	let showNotifications = false;
	let searchQuery = '';
	let expandedWorkspaces: Set<string> = new Set();
	let projectFilters = {
		status: 'all',
		type: 'all',
		priority: 'all'
	};
	
	// Available workspace icons
	const workspaceIcons = [
		{ id: 'briefcase', component: Briefcase, name: 'Briefcase', color: 'from-blue-500 to-blue-600' },
		{ id: 'code', component: Code, name: 'Code', color: 'from-green-500 to-green-600' },
		{ id: 'database', component: Database, name: 'Database', color: 'from-purple-500 to-purple-600' },
		{ id: 'globe', component: Globe, name: 'Globe', color: 'from-indigo-500 to-indigo-600' },
		{ id: 'laptop', component: Laptop, name: 'Laptop', color: 'from-gray-500 to-gray-600' },
		{ id: 'smartphone', component: Smartphone, name: 'Mobile', color: 'from-pink-500 to-pink-600' },
		{ id: 'users', component: Users, name: 'Team', color: 'from-orange-500 to-orange-600' },
		{ id: 'zap', component: Zap, name: 'Lightning', color: 'from-yellow-500 to-yellow-600' },
		{ id: 'building', component: Building, name: 'Building', color: 'from-slate-500 to-slate-600' },
		{ id: 'rocket', component: Rocket, name: 'Rocket', color: 'from-red-500 to-red-600' },
		{ id: 'target', component: Target, name: 'Target', color: 'from-emerald-500 to-emerald-600' },
		{ id: 'settings', component: Settings, name: 'Settings', color: 'from-cyan-500 to-cyan-600' }
	];
	
	// Project-focused stats
	let dashboardStats = {
		totalProjects: 0,
		activeProjects: 0,
		completedTasks: 0,
		overdueTasks: 0,
		teamMembers: 0,
		totalWorkspaces: 0
	};
	let recentTasks: any[] = [];
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
		workspaces = state.workspaces || [];
		isLoading = state.isLoading;
		error = state.error || '';
	});

	projectStore.subscribe(state => {
		projects = state.projects || [];
		updateProjectStats();
	});

	onMount(async () => {
		// Wait for auth store to initialize
		const unsubscribe = authStore.subscribe(async (state) => {
			console.log('Auth state changed:', state);
			if (state.isAuthenticated === false && !state.user) {
				console.log('Not authenticated, redirecting to login');
				goto('/login');
				return;
			}
			
			if (state.isAuthenticated && state.user) {
				console.log('User authenticated, loading dashboard data');
				await Promise.all([
					loadWorkspaces(),
					loadProjectData(),
					loadDashboardData()
				]);
				
				
				unsubscribe();
			}
		});
	});

	async function loadWorkspaces() {
		try {
			workspaceStore.setLoading(true);
			workspaceStore.setError(null);
			
			const data = await apiClient.getWorkspaces();
			workspaceStore.setWorkspaces(data);
			
			// Expand all workspaces by default to show projects (like skeleton)
			expandedWorkspaces = new Set(data.map((ws: any) => ws.id));
		} catch (err: any) {
			console.error('Error loading workspaces:', err);
			workspaceStore.setError(err.response?.data?.error || 'Failed to load workspaces');
		} finally {
			workspaceStore.setLoading(false);
		}
	}

	async function loadProjectData() {
		try {
			projectStore.setLoading(true);
			const data = await apiClient.getProjects();
			projectStore.setProjects(data);
		} catch (err: any) {
			console.error('Error loading projects:', err);
			projectStore.setError(err.response?.data?.error || 'Failed to load projects');
		} finally {
			projectStore.setLoading(false);
		}
	}

	async function loadDashboardData() {
		try {
			console.log('Loading dashboard data...');
			isLoadingDashboard = true;
			
			// Load dashboard activity and notifications
			const [activity, team, notif] = await Promise.allSettled([
				apiClient.getRecentActivity(),
				apiClient.getTeamMembers(),
				apiClient.getNotifications()
			]);

			console.log('Dashboard API results:', { activity: activity.status, team: team.status, notif: notif.status });

			if (activity.status === 'fulfilled') {
				recentActivity = activity.value?.activities || [];
				console.log('Recent activity loaded:', recentActivity.length, activity.value);
			} else {
				console.error('Activity API failed:', activity.reason);
				// Provide mock data as fallback
				recentActivity = [
					{
						id: '1',
						type: 'project_created',
						description: 'Created new project "E-commerce Platform"',
						timestamp: Date.now() - 2 * 60 * 60 * 1000,
						user: 'You'
					},
					{
						id: '2', 
						type: 'task_completed',
						description: 'Completed task "Setup Authentication"',
						timestamp: Date.now() - 4 * 60 * 60 * 1000,
						user: 'You'
					},
					{
						id: '3',
						type: 'workspace_updated',
						description: 'Updated workspace "Development Team"',
						timestamp: Date.now() - 24 * 60 * 60 * 1000,
						user: 'You'
					}
				];
			}

			if (team.status === 'fulfilled') {
				const teamData = team.value?.team_members;
				teamMembers = Array.isArray(teamData) ? teamData : [];
				console.log('Team members loaded:', teamMembers.length, team.value);
			} else {
				console.error('Team API failed:', team.reason);
				// Provide mock data as fallback
				teamMembers = [
					{
						id: '1',
						name: 'You',
						email: 'you@example.com',
						avatar: '',
						status: 'online',
						role: 'owner'
					},
					{
						id: '2',
						name: 'Alice Smith',
						email: 'alice@example.com',
						avatar: '',
						status: 'online',
						role: 'member'
					},
					{
						id: '3',
						name: 'Bob Johnson',
						email: 'bob@example.com',
						avatar: '',
						status: 'offline',
						role: 'member'
					}
				];
			}

			if (notif.status === 'fulfilled') {
				notifications = notif.value || [];
				console.log('Notifications loaded:', notifications.length);
			} else {
				console.error('Notifications API failed:', notif.reason);
				notifications = [];
			}

		} catch (err: any) {
			console.error('Error loading dashboard data:', err);
			// Set fallback mock data
			recentActivity = [
				{
					id: '1',
					type: 'project_created',
					description: 'Created new project "E-commerce Platform"',
					timestamp: Date.now() - 2 * 60 * 60 * 1000,
					user: 'You'
				}
			];
			teamMembers = [
				{
					id: '1',
					name: 'You',
					email: 'you@example.com',
					avatar: '',
					status: 'online',
					role: 'owner'
				}
			];
			notifications = [];
		} finally {
			console.log('Dashboard loading complete, setting isLoadingDashboard to false');
			isLoadingDashboard = false;
		}
	}


	function updateProjectStats() {
		dashboardStats = {
			totalProjects: projects.length,
			activeProjects: projects.filter(p => p.status === 'active' || p.status === 'planning').length,
			completedTasks: projects.reduce((sum, p) => sum + p.stats.completedTasks, 0),
			overdueTasks: projects.reduce((sum, p) => sum + p.stats.overdueTasks, 0),
			teamMembers: new Set(projects.flatMap(p => p.members.map(m => m.id))).size,
			totalWorkspaces: workspaces.length
		};
	}

	async function createWorkspace() {
		if (!newWorkspaceName.trim()) {
			return;
		}

		isCreating = true;
		try {
			const workspaceData = {
				name: newWorkspaceName,
				description: newWorkspaceDescription,
				icon: workspaceIcon,
				imageUrl: uploadedImageUrl
			};
			
			const workspace = await apiClient.createWorkspace(workspaceData.name, workspaceData.description);
			
			await loadWorkspaces();
			
			toastStore.success(`Workspace "${newWorkspaceName}" created successfully!`);
			
			resetWorkspaceForm();
		} catch (err: any) {
			const errorMessage = err.response?.data?.error || 'Failed to create workspace';
			error = errorMessage;
			toastStore.error(errorMessage);
		} finally {
			isCreating = false;
		}
	}

	function openEditWorkspace(workspace: any) {
		selectedWorkspace = workspace;
		newWorkspaceName = workspace.name;
		newWorkspaceDescription = workspace.description || '';
		workspaceIcon = workspace.icon || 'briefcase';
		uploadedImageUrl = workspace.imageUrl || '';
		showEditWorkspaceModal = true;
	}

	async function updateWorkspace() {
		if (!newWorkspaceName.trim() || !selectedWorkspace) {
			return;
		}

		isUpdating = true;
		try {
			const workspaceData = {
				name: newWorkspaceName,
				description: newWorkspaceDescription,
				icon: workspaceIcon,
				imageUrl: uploadedImageUrl
			};
			
			// TODO: Add API endpoint for updating workspace
			// await apiClient.updateWorkspace(selectedWorkspace.id, workspaceData);
			
			// For now, update locally
			const updatedWorkspaces = workspaces.map(ws => 
				ws.id === selectedWorkspace.id 
					? { ...ws, ...workspaceData }
					: ws
			);
			workspaceStore.setWorkspaces(updatedWorkspaces);
			
			toastStore.success(`Workspace "${newWorkspaceName}" updated successfully!`);
			
			resetWorkspaceForm();
			showEditWorkspaceModal = false;
		} catch (err: any) {
			const errorMessage = err.response?.data?.error || 'Failed to update workspace';
			toastStore.error(errorMessage);
		} finally {
			isUpdating = false;
		}
	}

	function resetWorkspaceForm() {
		newWorkspaceName = '';
		newWorkspaceDescription = '';
		workspaceIcon = 'briefcase';
		uploadedImageUrl = '';
		workspaceImage = null;
		selectedWorkspace = null;
		showCreateWorkspaceModal = false;
		showEditWorkspaceModal = false;
	}

	function handleImageUpload(event: Event) {
		const file = (event.target as HTMLInputElement).files?.[0];
		if (file) {
			if (file.size > 2 * 1024 * 1024) { // 2MB limit
				toastStore.error('Image size must be less than 2MB');
				return;
			}
			
			if (!file.type.startsWith('image/')) {
				toastStore.error('Please select a valid image file');
				return;
			}
			
			workspaceImage = file;
			
			// Create preview URL
			const reader = new FileReader();
			reader.onload = (e) => {
				uploadedImageUrl = e.target?.result as string;
			};
			reader.readAsDataURL(file);
		}
	}

	function removeImage() {
		workspaceImage = null;
		uploadedImageUrl = '';
		// Reset file input
		const fileInput = document.getElementById('workspace-image-upload') as HTMLInputElement;
		if (fileInput) fileInput.value = '';
	}

	function getWorkspaceIcon(workspace: any) {
		if (workspace.imageUrl) return null;
		const iconData = workspaceIcons.find(icon => icon.id === (workspace.icon || 'briefcase'));
		return iconData || workspaceIcons[0];
	}

	function selectWorkspace(workspace: any) {
		workspaceStore.setCurrentWorkspace(workspace);
		goto(`/workspace/${workspace.id}`);
	}

	function selectProject(project: Project) {
		projectStore.setCurrentProject(project);
		goto(`/project/${project.id}`);
	}

	function toggleWorkspace(workspaceId: string) {
		if (expandedWorkspaces.has(workspaceId)) {
			expandedWorkspaces.delete(workspaceId);
		} else {
			expandedWorkspaces.add(workspaceId);
		}
		expandedWorkspaces = expandedWorkspaces;
	}

	function getWorkspaceProjects(workspaceId: string) {
		// For now, since we don't have workspace-project association yet,
		// we'll show all projects for demonstration with proper filtering
		return filterProjects(projects);
	}

	function logout() {
		apiClient.logout();
		goto('/');
	}

	async function markAllNotificationsAsRead() {
		try {
			await apiClient.markAllNotificationsAsRead();
			notifications = Array.isArray(notifications) ? notifications.map(n => ({ ...n, read: true })) : [];
		} catch (err) {
			console.error('Failed to mark notifications as read:', err);
		}
	}

	function filterProjects(projects: Project[]) {
		return projects.filter(project => {
			const matchesSearch = !searchQuery.trim() || 
				project.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				project.description?.toLowerCase().includes(searchQuery.toLowerCase()) ||
				project.tags.some(tag => tag.toLowerCase().includes(searchQuery.toLowerCase()));
			
			const matchesStatus = projectFilters.status === 'all' || project.status === projectFilters.status;
			const matchesType = projectFilters.type === 'all' || project.type === projectFilters.type;
			const matchesPriority = projectFilters.priority === 'all' || project.priority === projectFilters.priority;
			
			return matchesSearch && matchesStatus && matchesType && matchesPriority;
		});
	}

	$: filteredWorkspaces = workspaces.filter(workspace => 
		!searchQuery.trim() || 
		workspace.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
		workspace.description?.toLowerCase().includes(searchQuery.toLowerCase())
	);
</script>

<svelte:head>
	<title>Project Dashboard - Volt</title>
</svelte:head>

<div class="min-h-screen bg-slate-50">
	<!-- Header -->
	<header class="bg-white border-b border-slate-200">
		<div class="max-w-7xl mx-auto px-4 lg:px-6">
			<div class="flex justify-between items-center h-12">
				<div class="flex items-center space-x-2">
					<div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
						<Zap class="w-5 h-5 text-white" />
					</div>
					<h1 class="text-lg font-bold text-slate-900">Volt</h1>
				</div>
				
				<div class="flex items-center space-x-2">
					<!-- Search Bar -->
					<div class="relative">
						<input
							type="text"
							bind:value={searchQuery}
							placeholder="Search..."
							class="bg-slate-50 border border-slate-200 text-slate-900 placeholder-slate-400 px-2 py-1.5 pl-6 text-sm rounded focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 w-48"
						/>
						<Search class="absolute left-2 top-2 w-3.5 h-3.5 text-slate-400" />
					</div>
					
					<!-- Notifications -->
					<div class="relative">
						<IconButton
							variant="ghost"
							size="sm"
							on:click={() => showNotifications = !showNotifications}
							class="relative"
						>
							<Bell class="w-4 h-4" />
							{#if Array.isArray(notifications) && notifications.filter(n => !n.read).length > 0}
								<span class="absolute -top-0.5 -right-0.5 bg-red-500 text-white rounded-full h-4 w-4 flex items-center justify-center text-[10px]">
									{Array.isArray(notifications) ? notifications.filter(n => !n.read).length : 0}
								</span>
							{/if}
						</IconButton>
						
						{#if showNotifications}
							<div class="absolute right-0 mt-2 w-72 bg-white border border-slate-200 rounded-lg shadow-lg z-50">
								<div class="p-3 border-b border-slate-100">
									<h3 class="text-label-lg text-slate-900">Notifications</h3>
								</div>
								<div class="max-h-64 overflow-y-auto">
									{#each Array.isArray(notifications) ? notifications.slice(0, 4) : [] as notification}
										<div class="p-3 border-b border-slate-100 {notification.read ? 'bg-white' : 'bg-blue-50'} hover:bg-slate-50 transition-colors">
											<div class="flex items-start space-x-2">
												<div class="flex-shrink-0 mt-1">
													{#if notification.type === 'success'}
														<div class="w-1.5 h-1.5 bg-green-500 rounded-full"></div>
													{:else if notification.type === 'warning'}
														<div class="w-1.5 h-1.5 bg-amber-500 rounded-full"></div>
													{:else if notification.type === 'error'}
														<div class="w-1.5 h-1.5 bg-red-500 rounded-full"></div>
													{:else}
														<div class="w-1.5 h-1.5 bg-blue-500 rounded-full"></div>
													{/if}
												</div>
												<div class="flex-1">
													<p class="text-xs font-medium text-slate-900">{notification.title}</p>
													<p class="text-xs text-slate-600 mt-0.5">{notification.message}</p>
													<p class="text-xs text-slate-400 mt-0.5">{notification.time}</p>
												</div>
											</div>
										</div>
									{/each}
								</div>
								<div class="p-2 border-t border-slate-100">
									<Button
										variant="ghost"
										size="xs"
										on:click={markAllNotificationsAsRead}
										class="text-blue-600 hover:text-blue-700 w-full"
									>
										Mark all as read
									</Button>
								</div>
							</div>
						{/if}
					</div>
					
					<div class="flex items-center space-x-1.5">
						<span class="text-slate-600 text-sm">{user?.username || 'User'}</span>
						<IconButton
							variant="ghost"
							size="sm"
							on:click={logout}
							class="text-slate-400 hover:text-slate-600"
						>
							<LogOut class="w-4 h-4" />
						</IconButton>
					</div>
				</div>
			</div>
		</div>
	</header>

	<!-- Main Content -->
	<main class="max-w-7xl mx-auto py-4 px-4 lg:px-6">
		<!-- Welcome Section -->
		<div class="mb-6">
			<div class="mb-4">
				<h1 class="text-xl font-bold text-slate-900 mb-1">Welcome back, {user?.username || 'Developer'}!</h1>
				<p class="text-sm text-slate-600">
					Manage your workspaces and projects in one place.
				</p>
			</div>
				
			<!-- Stats Cards -->
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
				<div class="bg-white rounded-lg p-4 shadow-sm border border-slate-200 hover:shadow-md transition-all duration-200">
					<div class="flex items-center justify-between mb-2">
						<div class="p-2 bg-blue-50 rounded-lg">
							<Archive class="w-4 h-4 text-blue-600" />
						</div>
						<span class="text-xs font-medium text-blue-600 bg-blue-50 px-2 py-0.5 rounded-full">+2</span>
					</div>
					<div>
						<p class="text-2xl font-bold text-slate-900 mb-0.5">{workspaces?.length || 0}</p>
						<p class="text-xs text-slate-500">Workspaces</p>
					</div>
				</div>
				
				<div class="bg-white rounded-lg p-4 shadow-sm border border-slate-200 hover:shadow-md transition-all duration-200">
					<div class="flex items-center justify-between mb-2">
						<div class="p-2.5 bg-green-50 rounded-lg">
							<Activity class="icon-md text-green-600" />
						</div>
						<span class="text-xs font-medium text-green-600 bg-green-50 px-2 py-0.5 rounded-full">{dashboardStats.activeProjects}</span>
					</div>
					<div>
						<p class="text-2xl font-bold text-slate-900 mb-0.5">{dashboardStats.totalProjects}</p>
						<p class="text-xs text-slate-500">Projects</p>
					</div>
				</div>
				
				<div class="bg-white rounded-lg p-4 shadow-sm border border-slate-200 hover:shadow-md transition-all duration-200">
					<div class="flex items-center justify-between mb-2">
						<div class="p-2.5 bg-emerald-50 rounded-lg">
							<CheckCircle class="icon-md text-emerald-600" />
						</div>
						<span class="text-xs font-medium text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-full">Done</span>
					</div>
					<div>
						<p class="text-2xl font-bold text-slate-900 mb-0.5">{dashboardStats.completedTasks}</p>
						<p class="text-xs text-slate-500">Tasks</p>
					</div>
				</div>
				
				<div class="bg-white rounded-lg p-4 shadow-sm border border-slate-200 hover:shadow-md transition-all duration-200">
					<div class="flex items-center justify-between mb-2">
						<div class="p-2.5 bg-purple-50 rounded-lg">
							<Users class="icon-md text-purple-600" />
						</div>
						<span class="text-xs font-medium text-purple-600 bg-purple-50 px-2 py-0.5 rounded-full">Online</span>
					</div>
					<div>
						<p class="text-2xl font-bold text-slate-900 mb-0.5">{dashboardStats.teamMembers}</p>
						<p class="text-xs text-slate-500">Members</p>
					</div>
				</div>
			</div>
		</div>
			
		<!-- Main Content Grid -->
		<div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
			<!-- Left Column: Workspaces with nested Projects -->
			<div class="lg:col-span-3">
				<div class="flex items-center justify-between mb-4">
					<h2 class="text-lg font-semibold text-slate-900">Workspaces & Projects</h2>
					<div class="flex items-center space-x-2">
						<!-- Filters -->
						<select
							bind:value={projectFilters.status}
							class="bg-white border border-slate-200 text-slate-900 text-xs rounded-md px-2 py-1.5 focus:outline-none focus:ring-1 focus:ring-blue-500"
						>
							<option value="all">All Status</option>
							<option value="planning">Planning</option>
							<option value="active">Active</option>
							<option value="on_hold">On Hold</option>
							<option value="completed">Completed</option>
							<option value="archived">Archived</option>
						</select>
						
						<select
							bind:value={projectFilters.type}
							class="bg-white border border-slate-200 text-slate-900 text-xs rounded-md px-2 py-1.5 focus:outline-none focus:ring-1 focus:ring-blue-500"
						>
							<option value="all">All Types</option>
							<option value="api">API</option>
							<option value="web">Web App</option>
							<option value="mobile">Mobile</option>
							<option value="desktop">Desktop</option>
							<option value="data">Data</option>
							<option value="other">Other</option>
						</select>
						
						<Button
							variant="primary"
							size="xs"
							on:click={() => showCreateWorkspaceModal = true}
							class="bg-indigo-600 hover:bg-indigo-700"
						>
							<Briefcase class="btn-icon-xs mr-1" />
							Workspace
						</Button>
						
						<Button
							variant="primary"
							size="xs"
							on:click={() => showCreateProjectModal = true}
						>
							<Plus class="btn-icon-xs mr-1" />
							Project
						</Button>
					</div>
				</div>

				{#if error}
					<div class="bg-red-50 border border-red-200 rounded-lg p-3 mb-4">
						<div class="flex">
							<div class="flex-shrink-0">
								<svg class="h-4 w-4 text-red-600" viewBox="0 0 20 20" fill="currentColor">
									<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
								</svg>
							</div>
							<div class="ml-2">
								<p class="text-xs text-red-800">{error}</p>
							</div>
						</div>
					</div>
				{/if}

				{#if isLoading}
					<!-- Skeleton Loading -->
					<div class="space-y-6">
						{#each Array(2) as _, i}
							<div class="bg-gradient-to-br from-white to-slate-50 rounded-xl border border-slate-200 shadow-sm">
								<!-- Workspace Header Skeleton -->
								<div class="p-5">
									<div class="flex items-start justify-between">
										<div class="flex items-start space-x-4 flex-1">
											<div class="w-12 h-12 bg-slate-200 rounded-xl animate-pulse"></div>
											<div class="flex-1">
												<div class="flex items-center space-x-3 mb-2">
													<div class="h-5 bg-slate-200 rounded animate-pulse" style="width: {140 + i * 20}px"></div>
													<div class="h-5 bg-slate-200 rounded-full animate-pulse w-16"></div>
												</div>
												<div class="h-4 bg-slate-200 rounded animate-pulse mb-3" style="width: {220 + i * 30}px"></div>
												
												<!-- Stats Skeleton -->
												<div class="flex items-center space-x-6">
													<div class="flex items-center space-x-2">
														<div class="flex -space-x-1">
															{#each Array(3) as _}
																<div class="w-6 h-6 bg-slate-200 rounded-full animate-pulse"></div>
															{/each}
														</div>
														<div class="h-3 bg-slate-200 rounded animate-pulse w-16"></div>
													</div>
													<div class="flex items-center space-x-1">
														<div class="w-2 h-2 bg-slate-200 rounded-full animate-pulse"></div>
														<div class="h-3 bg-slate-200 rounded animate-pulse w-12"></div>
													</div>
												</div>
											</div>
										</div>
										<div class="flex items-start space-x-2">
											<div class="h-7 bg-slate-200 rounded-lg animate-pulse w-16"></div>
											<div class="w-8 h-8 bg-slate-200 rounded-lg animate-pulse"></div>
										</div>
									</div>
								</div>
								
								<!-- Projects Skeleton -->
								<div class="border-t border-slate-100">
									<div class="p-5">
										<div class="flex items-center justify-between mb-4">
											<div class="h-4 bg-slate-200 rounded animate-pulse w-24"></div>
											<div class="h-6 bg-slate-200 rounded-md animate-pulse w-20"></div>
										</div>
										<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
											{#each Array(2) as _, j}
												<div class="bg-slate-50 rounded-lg p-4 border border-slate-100">
													<div class="flex items-start justify-between mb-3">
														<div class="flex-1">
															<div class="h-4 bg-slate-200 rounded animate-pulse mb-2" style="width: {100 + j * 25}px"></div>
															<div class="h-3 bg-slate-200 rounded animate-pulse" style="width: {150 + j * 40}px"></div>
														</div>
														<div class="w-6 h-6 bg-slate-200 rounded animate-pulse"></div>
													</div>
													<div class="flex items-center justify-between mb-3">
														<div class="h-2 bg-slate-200 rounded animate-pulse w-20"></div>
														<div class="h-4 bg-slate-200 rounded animate-pulse w-12"></div>
													</div>
													<div class="w-full bg-slate-200 rounded-full h-1.5 animate-pulse mb-3"></div>
													<div class="flex items-center justify-between">
														<div class="flex space-x-1">
															{#each Array(3) as _}
																<div class="w-5 h-5 bg-slate-200 rounded-full animate-pulse"></div>
															{/each}
														</div>
														<div class="h-3 bg-slate-200 rounded animate-pulse w-16"></div>
													</div>
												</div>
											{/each}
										</div>
									</div>
								</div>
							</div>
						{/each}
					</div>
				{:else if !workspaces || workspaces.length === 0}
					<div class="text-center py-12">
						<div class="bg-slate-100 rounded-full w-12 h-12 flex items-center justify-center mx-auto mb-3">
							<svg class="w-6 h-6 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
							</svg>
						</div>
						<h3 class="text-base font-semibold text-slate-900 mb-1">No workspaces found</h3>
						<p class="text-sm text-slate-600 mb-4">Get started by creating your first workspace to organize your projects.</p>
						<div class="flex justify-center space-x-3">
							<Button
								variant="primary"
								size="sm"
								on:click={() => showCreateWorkspaceModal = true}
								class="bg-indigo-600 hover:bg-indigo-700"
							>
								Create Your First Workspace
							</Button>
							<Button
								variant="primary"
								size="sm"
								on:click={() => showCreateProjectModal = true}
							>
								Create a Project
							</Button>
						</div>
					</div>
				{:else}
					<!-- Workspace Cards with nested Projects -->
					<div class="space-y-6">
						{#each filteredWorkspaces as workspace}
							<div class="group bg-gradient-to-br from-white to-slate-50 rounded-xl border border-slate-200 shadow-sm hover:shadow-lg hover:border-slate-300 transition-all duration-300">
								<!-- Workspace Header -->
								<div class="p-5">
									<div class="flex items-start justify-between">
										<div class="flex items-start space-x-4 flex-1">
											<div class="flex-shrink-0">
												{#if workspace.imageUrl}
													<img 
														src={workspace.imageUrl} 
														alt={workspace.name}
														class="w-12 h-12 rounded-xl object-cover shadow-lg"
													/>
												{:else}
													{@const iconData = getWorkspaceIcon(workspace)}
													<div class="w-12 h-12 bg-gradient-to-br {iconData.color} rounded-xl flex items-center justify-center shadow-lg">
														<svelte:component this={iconData.component} class="w-6 h-6 text-white" />
													</div>
												{/if}
											</div>
											<div class="flex-1 min-w-0">
												<div class="flex items-center space-x-3 mb-2">
													<h3 class="text-lg font-bold text-slate-900 truncate">{workspace.name}</h3>
													<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
														{getWorkspaceProjects(workspace.id).length} projects
													</span>
												</div>
												{#if workspace.description}
													<p class="text-sm text-slate-600 mb-3 leading-relaxed">{workspace.description}</p>
												{/if}
												
												<!-- Workspace Stats -->
												<div class="flex items-center space-x-6">
													<div class="flex items-center space-x-2">
														<div class="flex -space-x-1">
															{#each (workspace.members || []).slice(0, 3) as member, i}
																<div class="w-6 h-6 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 border-2 border-white flex items-center justify-center">
																	<span class="text-xs font-medium text-white">{member.name?.[0] || 'U'}</span>
																</div>
															{/each}
															{#if (workspace.members?.length || 0) > 3}
																<div class="w-6 h-6 rounded-full bg-slate-300 border-2 border-white flex items-center justify-center">
																	<span class="text-xs font-medium text-slate-600">+{(workspace.members?.length || 0) - 3}</span>
																</div>
															{/if}
														</div>
														<span class="text-xs text-slate-500 font-medium">
															{workspace.members?.length || 1} member{workspace.members?.length !== 1 ? 's' : ''}
														</span>
													</div>
													
													<div class="flex items-center space-x-1">
														<div class="w-2 h-2 rounded-full bg-green-500"></div>
														<span class="text-xs text-slate-500">Active</span>
													</div>
												</div>
											</div>
										</div>
										
										<div class="flex items-start space-x-2 ml-4">
											<Button
												variant="secondary"
												size="xs"
												on:click={(e) => { e.stopPropagation(); openEditWorkspace(workspace); }}
												class="text-slate-600"
											>
												<Settings class="btn-icon-xs mr-1" />
												Edit
											</Button>
											<Button
												variant="ghost"
												size="xs"
												on:click={(e) => { e.stopPropagation(); selectWorkspace(workspace); }}
												class="text-indigo-700 hover:bg-indigo-50"
											>
												<FolderOpen class="btn-icon-xs mr-1" />
												Open
											</Button>
											<IconButton
												variant="ghost"
												size="sm"
												on:click={(e) => { e.stopPropagation(); toggleWorkspace(workspace.id); }}
												class="text-slate-400 hover:text-slate-600"
											>
												<ChevronDown class="icon-sm transform transition-transform duration-200 {expandedWorkspaces.has(workspace.id) ? 'rotate-180' : ''}" />
											</IconButton>
										</div>
									</div>
								</div>
								
								<!-- Projects within Workspace -->
								{#if expandedWorkspaces.has(workspace.id)}
									{@const workspaceProjects = getWorkspaceProjects(workspace.id)}
									<div class="border-t border-slate-100">
										{#if workspaceProjects.length === 0}
											<div class="text-center py-8">
												<div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-full w-12 h-12 flex items-center justify-center mx-auto mb-3">
													<svg class="w-6 h-6 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
													</svg>
												</div>
												<h4 class="text-sm font-medium text-slate-900 mb-1">No projects yet</h4>
												<p class="text-xs text-slate-500 mb-4">Get started by creating your first project in this workspace</p>
												<Button
													variant="primary"
													size="sm"
													on:click={() => showCreateProjectModal = true}
												>
													<Plus class="btn-icon-sm mr-1" />
													Create Project
												</Button>
											</div>
										{:else}
											<div class="p-5">
												<div class="flex items-center justify-between mb-4">
													<h4 class="text-sm font-semibold text-slate-700">Projects ({workspaceProjects.length})</h4>
													<Button
														variant="ghost"
														size="xs"
														on:click={() => showCreateProjectModal = true}
														class="text-slate-600 hover:text-slate-900"
													>
														<Plus class="btn-icon-xs mr-1" />
														Add Project
													</Button>
												</div>
												<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
													{#each workspaceProjects as project}
														<ProjectCard {project} />
													{/each}
												</div>
											</div>
										{/if}
									</div>
								{/if}
							</div>
						{/each}
					</div>
				{/if}
			</div>
				
			<!-- Right Column: Activity & Quick Actions -->
			<div class="lg:col-span-1 space-y-4">
				<!-- Quick Actions -->
				<div class="bg-white rounded-lg border border-slate-200 p-4 shadow-sm">
					<h3 class="text-sm font-semibold text-slate-900 mb-3">Quick Actions</h3>
					<div class="space-y-2">
						<ActionCard
							title="New Project"
							description="Start development project"
							icon={Briefcase}
							iconColor="bg-blue-600"
							on:click={() => showCreateProjectModal = true}
						/>
						
						<ActionCard
							title="Create Task"
							description="Add task to project"
							icon={CheckCircle}
							iconColor="bg-green-600"
							disabled={true}
						/>
						
						<ActionCard
							title="New Workspace"
							description="Organize projects"
							icon={Building}
							iconColor="bg-indigo-600"
							on:click={() => showCreateWorkspaceModal = true}
						/>
					</div>
				</div>
					
				<!-- Recent Activity -->
				<div class="bg-white rounded-lg border border-slate-200 p-4 shadow-sm">
					<h3 class="text-sm font-semibold text-slate-900 mb-3">Recent Activity</h3>
					{#if isLoadingDashboard}
						<!-- Activity Skeleton -->
						<div class="space-y-2">
							{#each Array(4) as _, i}
								<div class="flex items-start space-x-2">
									<div class="w-1.5 h-1.5 rounded-full mt-1.5 bg-slate-200 animate-pulse"></div>
									<div class="flex-1">
										<div class="h-3 bg-slate-200 rounded animate-pulse mb-1" style="width: {120 + i * 20}px"></div>
										<div class="h-2 bg-slate-200 rounded animate-pulse w-16"></div>
									</div>
								</div>
							{/each}
						</div>
						<div class="mt-3 h-3 bg-slate-200 rounded animate-pulse w-16"></div>
					{:else}
						{#if Array.isArray(recentActivity) && recentActivity.length > 0}
							<div class="space-y-2">
								{#each recentActivity.slice(0, 4) as activity}
									<div class="flex items-start space-x-2">
										<div class="w-1.5 h-1.5 rounded-full mt-1.5 {
											activity.type === 'project_created' ? 'bg-blue-500' :
											activity.type === 'task_completed' ? 'bg-green-500' :
											activity.type === 'workspace_updated' ? 'bg-purple-500' :
											activity.type === 'request_created' ? 'bg-amber-500' :
											'bg-slate-400'
										}"></div>
										<div>
											<p class="text-xs text-slate-900">{activity.description || activity.message || 'No description'}</p>
											<p class="text-xs text-slate-500">{
												activity.time || 
												(activity.timestamp ? new Date(activity.timestamp * 1000).toLocaleString() : 'No time')
											}</p>
										</div>
									</div>
								{/each}
							</div>
							<Button
								variant="ghost"
								size="sm"
								class="text-blue-600 hover:text-blue-700 mt-3"
							>
								View all →
							</Button>
						{:else}
							<div class="text-center py-6">
								<div class="w-12 h-12 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-3">
									<ClipboardList class="w-6 h-6 text-slate-400" />
								</div>
								<p class="text-sm text-slate-500">No recent activity</p>
								<p class="text-xs text-slate-400 mt-1">Activity will appear here as you work on projects</p>
							</div>
						{/if}
					{/if}
				</div>
					
				<!-- Team Online -->
				<div class="bg-white rounded-lg border border-slate-200 p-4 shadow-sm">
					<div class="flex items-center justify-between mb-3">
						<h3 class="text-sm font-semibold text-slate-900">Team Online</h3>
						{#if isLoadingDashboard}
							<div class="h-4 bg-slate-200 rounded-full animate-pulse w-6"></div>
						{:else}
							<span class="text-xs font-medium text-green-600 bg-green-50 px-2 py-0.5 rounded-full">{Array.isArray(teamMembers) ? teamMembers.filter(m => m.status === 'online').length : 0}</span>
						{/if}
					</div>
					{#if isLoadingDashboard}
						<!-- Team Skeleton -->
						<div class="space-y-2">
							{#each Array(3) as _, i}
								<div class="flex items-center space-x-2 p-2">
									<div class="relative">
										<div class="w-6 h-6 bg-slate-200 rounded-full animate-pulse"></div>
										<div class="absolute -bottom-0.5 -right-0.5 w-2 h-2 bg-slate-200 rounded-full animate-pulse"></div>
									</div>
									<div class="flex-1">
										<div class="flex items-center justify-between">
											<div class="h-3 bg-slate-200 rounded animate-pulse" style="width: {60 + i * 15}px"></div>
											<div class="h-2 bg-slate-200 rounded animate-pulse w-8"></div>
										</div>
										<div class="h-2 bg-slate-200 rounded animate-pulse mt-1" style="width: {80 + i * 25}px"></div>
									</div>
								</div>
							{/each}
						</div>
						<div class="mt-3 h-3 bg-slate-200 rounded animate-pulse w-24"></div>
					{:else}
						{#if Array.isArray(teamMembers) && teamMembers.length > 0}
							<div class="space-y-2">
								{#each teamMembers.slice(0, 3) as member}
									<div class="flex items-center space-x-2 hover:bg-slate-50 p-2 rounded-lg transition-colors cursor-pointer">
										<div class="relative">
											<div class="w-6 h-6 bg-blue-600 rounded-full flex items-center justify-center">
												{#if member.avatar}
													<img src={member.avatar} alt={member.name} class="w-full h-full rounded-full object-cover" />
												{:else}
													<span class="text-xs font-medium text-white">{member.name ? member.name.charAt(0).toUpperCase() : 'U'}</span>
												{/if}
											</div>
											<div class="absolute -bottom-0.5 -right-0.5 w-2 h-2 border border-white rounded-full {
												member.status === 'online' ? 'bg-green-500' :
												member.status === 'away' ? 'bg-amber-500' :
												'bg-slate-400'
											}"></div>
										</div>
										<div class="flex-1 min-w-0">
											<div class="flex items-center justify-between">
												<p class="text-xs font-medium text-slate-900 truncate">{member.name || 'Unknown User'}</p>
												<span class="text-xs text-slate-400">{member.lastSeen || 'Just now'}</span>
											</div>
											<p class="text-xs text-slate-500 truncate">{member.role || 'Member'}</p>
										</div>
									</div>
								{/each}
							</div>
							<Button
								variant="ghost"
								size="sm"
								class="text-blue-600 hover:text-blue-700 mt-3 w-full text-left"
							>
								View all members →
							</Button>
						{:else}
							<div class="text-center py-6">
								<div class="w-12 h-12 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-3">
									<UserPlus class="w-6 h-6 text-slate-400" />
								</div>
								<p class="text-sm text-slate-500">No team members</p>
								<p class="text-xs text-slate-400 mt-1">Invite team members to collaborate</p>
							</div>
						{/if}
					{/if}
				</div>
			</div>
		</div>
	</main>
</div>

<!-- Create Project Modal -->
<CreateProjectModal bind:show={showCreateProjectModal} />

<!-- Create Workspace Modal -->
{#if showCreateWorkspaceModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-2xl shadow-xl max-w-md w-full">
			<div class="px-6 py-4 border-b border-slate-200">
				<h3 class="text-lg font-semibold text-slate-900">Create New Workspace</h3>
			</div>
			
			<div class="px-6 py-6">
				<div class="space-y-4">
					<div>
						<label for="workspaceName" class="block text-sm font-medium text-slate-700 mb-2">
							Workspace Name
						</label>
						<input
							id="workspaceName"
							type="text"
							bind:value={newWorkspaceName}
							placeholder="My API Project"
							class="block w-full px-4 py-3 text-sm bg-white border border-slate-200 rounded-lg shadow-sm placeholder-slate-400 text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>
					
					<div>
						<label for="workspaceDescription" class="block text-sm font-medium text-slate-700 mb-2">
							Description (optional)
						</label>
						<textarea
							id="workspaceDescription"
							bind:value={newWorkspaceDescription}
							placeholder="A workspace for testing APIs..."
							rows="3"
							class="block w-full px-4 py-3 text-sm bg-white border border-slate-200 rounded-lg shadow-sm placeholder-slate-400 text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						></textarea>
					</div>
				</div>
			</div>
			
			<div class="px-6 py-4 bg-slate-50 flex justify-end space-x-3 rounded-b-2xl">
				<Button
					variant="secondary"
					size="xs"
					on:click={() => showCreateWorkspaceModal = false}
					disabled={isCreating}
				>
					Cancel
				</Button>
				<Button
					variant="primary"
					size="xs"
					on:click={createWorkspace}
					disabled={isCreating || !newWorkspaceName.trim()}
					loading={isCreating}
				>
					{isCreating ? 'Creating...' : 'Create Workspace'}
				</Button>
			</div>
		</div>
	</div>
{/if}

<!-- Edit Workspace Modal -->
{#if showEditWorkspaceModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-2xl shadow-2xl max-w-lg w-full max-h-[90vh] overflow-y-auto">
			<div class="px-6 py-4 border-b border-slate-100">
				<h3 class="text-md font-semibold text-slate-900">Edit Workspace</h3>
				<p class="text-xs text-slate-600 mt-1">Customize your workspace appearance and details</p>
			</div>
			
			<div class="px-6 py-4 space-y-4">
				<!-- Icon/Image Selection -->
				<div>
					<label class="block text-xs font-medium text-slate-900 mb-2">Workspace Icon</label>
					
					<!-- Current Icon Preview -->
					<div class="flex items-center space-x-3 mb-3">
						<div class="flex-shrink-0">
							{#if uploadedImageUrl}
								<img 
									src={uploadedImageUrl} 
									alt="Workspace preview"
									class="w-12 h-12 rounded-lg object-cover shadow-md"
								/>
							{:else}
								{@const iconData = workspaceIcons.find(icon => icon.id === workspaceIcon) || workspaceIcons[0]}
								<div class="w-12 h-12 bg-gradient-to-br {iconData.color} rounded-lg flex items-center justify-center shadow-md">
									<svelte:component this={iconData.component} class="w-6 h-6 text-white" />
								</div>
							{/if}
						</div>
						{#if uploadedImageUrl}
							<Button
								variant="ghost"
								size="xs"
								on:click={removeImage}
								class="text-red-600 hover:text-red-700"
							>
								Remove Image
							</Button>
						{/if}
					</div>
					
					<!-- Image Upload -->
					<div class="mb-3">
						<label for="workspace-image-upload" class="block text-xs font-medium text-slate-700 mb-1">
							Upload Custom Image (max 2MB)
						</label>
						<input
							id="workspace-image-upload"
							type="file"
							accept="image/*"
							on:change={handleImageUpload}
							class="block w-full text-xs text-slate-500 file:mr-3 file:py-1.5 file:px-3 file:rounded-md file:border-0 file:text-xs file:font-medium file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 transition-colors"
						/>
					</div>
					
					<!-- Icon Picker -->
					{#if !uploadedImageUrl}
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Or choose an icon</label>
							<div class="grid grid-cols-8 gap-1.5">
								{#each workspaceIcons as icon}
									<Button
										variant="ghost"
										size="sm"
										on:click={() => workspaceIcon = icon.id}
										class="p-2 rounded-md border-2 transition-all {workspaceIcon === icon.id ? 'border-blue-500 bg-blue-50' : 'border-slate-200 hover:border-slate-300'}"
									>
										<div class="w-6 h-6 bg-gradient-to-br {icon.color} rounded-md flex items-center justify-center">
											<svelte:component this={icon.component} class="w-3 h-3 text-white" />
										</div>
									</Button>
								{/each}
							</div>
						</div>
					{/if}
				</div>
				
				<!-- Workspace Name -->
				<div>
					<label for="edit-workspace-name" class="block text-xs font-medium text-slate-900 mb-1">
						Workspace Name *
					</label>
					<input
						id="edit-workspace-name"
						type="text"
						bind:value={newWorkspaceName}
						placeholder="Enter workspace name"
						class="w-full px-2 py-1.5 border border-slate-200 rounded-md text-sm text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
						required
					/>
				</div>
				
				<!-- Workspace Description -->
				<div>
					<label for="edit-workspace-description" class="block text-xs font-medium text-slate-900 mb-1">
						Description
					</label>
					<textarea
						id="edit-workspace-description"
						bind:value={newWorkspaceDescription}
						placeholder="Brief description of your workspace"
						rows="2"
						class="w-full px-2 py-1.5 border border-slate-200 rounded-md text-sm text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 resize-none"
					></textarea>
				</div>
			</div>
			
			<div class="px-6 py-4 bg-slate-50 flex justify-end space-x-3 rounded-b-2xl">
				<Button
					variant="secondary"
					size="xs"
					on:click={resetWorkspaceForm}
					disabled={isUpdating}
				>
					Cancel
				</Button>
				<Button
					variant="primary"
					size="xs"
					on:click={updateWorkspace}
					disabled={isUpdating || !newWorkspaceName.trim()}
					loading={isUpdating}
				>
					{isUpdating ? 'Updating...' : 'Update Workspace'}
				</Button>
			</div>
		</div>
	</div>
{/if}

