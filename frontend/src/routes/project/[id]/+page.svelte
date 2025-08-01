<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { projectStore, type Project, type ProjectTask, type ProjectMilestone } from '$lib/stores/project';
	import { authStore } from '$lib/stores/auth';
	import { toastStore } from '$lib/stores/toast';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import ProjectTimeline from '$lib/components/project/ProjectTimeline.svelte';
	import TaskBoard from '$lib/components/project/TaskBoard.svelte';
	import CreateTaskModal from '$lib/components/project/CreateTaskModal.svelte';
	import CreateMilestoneModal from '$lib/components/project/CreateMilestoneModal.svelte';
	import ProjectChat from '$lib/components/project/ProjectChat.svelte';
	import ProjectFiles from '$lib/components/project/ProjectFiles.svelte';
	import ProjectAnalytics from '$lib/components/project/ProjectAnalytics.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import { ChevronLeft, UserPlus } from 'lucide-svelte';

	let project: Project | null = null;
	let isLoading = true;
	let error = '';
	let activeTab: 'overview' | 'tasks' | 'timeline' | 'chat' | 'files' | 'analytics' = 'overview';
	let showTaskModal = false;
	let showMilestoneModal = false;
	let editingTask: ProjectTask | null = null;
	let editingMilestone: ProjectMilestone | null = null;
	let newTaskStatus: ProjectTask['status'] = 'todo';

	// Get project ID from URL params
	$: projectId = $page.params.id;

	// Subscribe to project store
	projectStore.subscribe(state => {
		if (projectId) {
			project = state.projects.find(p => p.id === projectId) || null;
		}
	});

	onMount(async () => {
		// Check authentication
		const unsubscribe = authStore.subscribe(async (state) => {
			if (!state.isAuthenticated) {
				goto('/login');
				return;
			}
			
			if (projectId) {
				await loadProject();
			}
			unsubscribe();
		});
	});

	async function loadProject() {
		try {
			isLoading = true;
			error = '';
			
			// TODO: Replace with actual API call
			// const projectData = await apiClient.getProject(projectId);
			// projectStore.setCurrentProject(projectData);
			
			// For now, use local store data
			if (!project) {
				error = 'Project not found';
			}
		} catch (err: any) {
			console.error('Error loading project:', err);
			error = err.message || 'Failed to load project';
		} finally {
			isLoading = false;
		}
	}

	function getStatusColor(status: Project['status']) {
		switch (status) {
			case 'active': return 'bg-green-500';
			case 'planning': return 'bg-blue-500';
			case 'on_hold': return 'bg-yellow-500';
			case 'completed': return 'bg-purple-500';
			case 'archived': return 'bg-gray-500';
			default: return 'bg-gray-500';
		}
	}

	function getPriorityColor(priority: Project['priority']) {
		switch (priority) {
			case 'urgent': return 'text-red-400 bg-red-400/10';
			case 'high': return 'text-orange-400 bg-orange-400/10';
			case 'medium': return 'text-yellow-400 bg-yellow-400/10';
			case 'low': return 'text-green-400 bg-green-400/10';
			default: return 'text-gray-400 bg-gray-400/10';
		}
	}

	function formatDate(date: Date) {
		return new Intl.DateTimeFormat('en-US', {
			month: 'long',
			day: 'numeric',
			year: 'numeric'
		}).format(new Date(date));
	}

	function calculateDaysRemaining() {
		if (!project?.dueDate) return null;
		const today = new Date();
		const due = new Date(project.dueDate);
		const diffTime = due.getTime() - today.getTime();
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
		return diffDays;
	}

	function handleCreateTask(event: CustomEvent) {
		newTaskStatus = event.detail.status;
		showTaskModal = true;
	}

	function handleEditTask(event: CustomEvent) {
		editingTask = event.detail;
		showTaskModal = true;
	}

	function handleTaskCreated() {
		showTaskModal = false;
		editingTask = null;
		toastStore.success('Task created successfully!');
	}

	function handleTaskUpdated() {
		showTaskModal = false;
		editingTask = null;
		toastStore.success('Task updated successfully!');
	}

	function handleCreateMilestone() {
		showMilestoneModal = true;
	}

	function handleEditMilestone(event: CustomEvent) {
		editingMilestone = event.detail;
		showMilestoneModal = true;
	}

	function handleMilestoneCreated() {
		showMilestoneModal = false;
		editingMilestone = null;
		toastStore.success('Milestone created successfully!');
	}

	function handleMilestoneUpdated() {
		showMilestoneModal = false;
		editingMilestone = null;
		toastStore.success('Milestone updated successfully!');
	}

	function goBack() {
		goto('/dashboard');
	}

	$: daysRemaining = calculateDaysRemaining();
	$: isOverdue = daysRemaining !== null && daysRemaining < 0;
	$: isDueSoon = daysRemaining !== null && daysRemaining <= 7 && daysRemaining >= 0;
	$: onlineMembers = project?.members.filter(m => m.status === 'online') || [];
</script>

<svelte:head>
	<title>{project?.name || 'Project'} - Volt</title>
</svelte:head>

<div class="min-h-screen bg-gray-900">
	{#if isLoading}
		<div class="flex items-center justify-center min-h-screen">
			<LoadingSpinner size="lg" color="blue" />
			<span class="ml-3 text-gray-400">Loading project...</span>
		</div>
	{:else if error}
		<div class="flex items-center justify-center min-h-screen">
			<div class="text-center">
				<svg class="w-16 h-16 text-gray-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
				</svg>
				<h2 class="text-xl font-semibold text-white mb-2">Project Not Found</h2>
				<p class="text-gray-400 mb-4">{error}</p>
				<Button
					variant="primary"
					size="sm"
					on:click={goBack}
				>
					Back to Dashboard
				</Button>
			</div>
		</div>
	{:else if project}
		<!-- Header -->
		<header class="bg-gray-800 shadow-sm border-b border-gray-700">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="flex items-center justify-between h-16">
					<div class="flex items-center space-x-4">
						<IconButton
							variant="ghost"
							size="md"
							on:click={goBack}
							class="text-gray-400 hover:text-gray-200"
						>
							<ChevronLeft class="w-6 h-6" />
						</IconButton>
						<h1 class="text-xl font-bold text-white">{project.name}</h1>
						<div class="flex items-center space-x-2">
							<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium {getPriorityColor(project.priority)}">
								{project.priority}
							</span>
							<div class="flex items-center space-x-1">
								<div class="w-2 h-2 rounded-full {getStatusColor(project.status)}"></div>
								<span class="text-xs text-gray-400 capitalize">{project.status.replace('_', ' ')}</span>
							</div>
						</div>
					</div>
					
					<div class="flex items-center space-x-4">
						<!-- Team avatars -->
						<div class="flex -space-x-2">
							{#each project.members.slice(0, 4) as member}
								<div 
									class="w-8 h-8 rounded-full bg-blue-600 border-2 border-gray-800 flex items-center justify-center relative"
									title={member.name}
								>
									{#if member.avatar}
										<img src={member.avatar} alt={member.name} class="w-full h-full rounded-full object-cover" />
									{:else}
										<span class="text-xs font-medium text-white">{member.name.charAt(0)}</span>
									{/if}
									{#if member.status === 'online'}
										<div class="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-green-500 border border-gray-800 rounded-full"></div>
									{/if}
								</div>
							{/each}
							{#if project.members.length > 4}
								<div class="w-8 h-8 rounded-full bg-gray-600 border-2 border-gray-800 flex items-center justify-center">
									<span class="text-xs font-medium text-gray-300">+{project.members.length - 4}</span>
								</div>
							{/if}
						</div>
						
						<Button
							variant="primary"
							size="sm"
						>
							<UserPlus class="w-4 h-4 mr-2" />
							Invite Members
						</Button>
					</div>
				</div>
			</div>
		</header>

		<!-- Project Stats Bar -->
		<div class="bg-gray-800 border-b border-gray-700">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
				<div class="grid grid-cols-4 gap-8">
					<div class="text-center">
						<div class="text-2xl font-semibold text-white">{project.progress}%</div>
						<div class="text-sm text-gray-400">Progress</div>
					</div>
					<div class="text-center">
						<div class="text-2xl font-semibold text-white">{project.stats.totalTasks}</div>
						<div class="text-sm text-gray-400">Total Tasks</div>
					</div>
					<div class="text-center">
						<div class="text-2xl font-semibold text-green-400">{project.stats.completedTasks}</div>
						<div class="text-sm text-gray-400">Completed</div>
					</div>
					<div class="text-center">
						{#if daysRemaining !== null}
							{#if isOverdue}
								<div class="text-2xl font-semibold text-red-400">{Math.abs(daysRemaining)}</div>
								<div class="text-sm text-red-400">Days Overdue</div>
							{:else}
								<div class="text-2xl font-semibold text-white">{daysRemaining}</div>
								<div class="text-sm text-gray-400">Days Left</div>
							{/if}
						{:else}
							<div class="text-2xl font-semibold text-gray-400">∞</div>
							<div class="text-sm text-gray-400">No Due Date</div>
						{/if}
					</div>
				</div>
			</div>
		</div>

		<!-- Tab Navigation -->
		<div class="border-b border-gray-700">
			<nav class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="flex space-x-8">
					<Button
						variant="ghost"
						size="sm"
						on:click={() => activeTab = 'overview'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab === 'overview' ? 'border-blue-500 text-blue-400' : 'border-transparent text-gray-400 hover:text-gray-300'}"
					>
						Overview
					</Button>
					<Button
						variant="ghost"
						size="sm"
						on:click={() => activeTab = 'tasks'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab === 'tasks' ? 'border-blue-500 text-blue-400' : 'border-transparent text-gray-400 hover:text-gray-300'}"
					>
						Tasks
						<span class="ml-2 bg-gray-700 text-gray-300 px-2 py-1 rounded-full text-xs">
							{project.stats.totalTasks}
						</span>
					</Button>
					<Button
						variant="ghost"
						size="sm"
						on:click={() => activeTab = 'timeline'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab === 'timeline' ? 'border-blue-500 text-blue-400' : 'border-transparent text-gray-400 hover:text-gray-300'}"
					>
						Timeline
						<span class="ml-2 bg-gray-700 text-gray-300 px-2 py-1 rounded-full text-xs">
							{project.milestones.length}
						</span>
					</Button>
					<Button
						variant="ghost"
						size="sm"
						on:click={() => activeTab = 'chat'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab === 'chat' ? 'border-blue-500 text-blue-400' : 'border-transparent text-gray-400 hover:text-gray-300'}"
					>
						Chat
						{#if onlineMembers.length > 0}
							<span class="ml-2 bg-green-600 text-white px-2 py-1 rounded-full text-xs">
								{onlineMembers.length}
							</span>
						{/if}
					</Button>
					<Button
						variant="ghost"
						size="sm"
						on:click={() => activeTab = 'files'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab === 'files' ? 'border-blue-500 text-blue-400' : 'border-transparent text-gray-400 hover:text-gray-300'}"
					>
						Files
					</Button>
					<Button
						variant="ghost"
						size="sm"
						on:click={() => activeTab = 'analytics'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab === 'analytics' ? 'border-blue-500 text-blue-400' : 'border-transparent text-gray-400 hover:text-gray-300'}"
					>
						Analytics
					</Button>
				</div>
			</nav>
		</div>

		<!-- Tab Content -->
		<main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
			<div class="px-4 py-6 sm:px-0">
				{#if activeTab === 'overview'}
					<!-- Overview Tab -->
					<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
						<!-- Left Column -->
						<div class="lg:col-span-2 space-y-6">
							<!-- Project Description -->
							<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
								<h3 class="text-lg font-semibold text-white mb-4">About This Project</h3>
								{#if project.description}
									<p class="text-gray-300 leading-relaxed">{project.description}</p>
								{:else}
									<p class="text-gray-500 italic">No description provided.</p>
								{/if}
								
								<!-- Project Details -->
								<div class="mt-6 grid grid-cols-2 gap-4">
									<div>
										<h4 class="text-sm font-medium text-gray-400 mb-2">Start Date</h4>
										<p class="text-white">{formatDate(project.startDate)}</p>
									</div>
									{#if project.dueDate}
										<div>
											<h4 class="text-sm font-medium text-gray-400 mb-2">Due Date</h4>
											<p class="text-white {isOverdue ? 'text-red-400' : isDueSoon ? 'text-yellow-400' : ''}">{formatDate(project.dueDate)}</p>
										</div>
									{/if}
								</div>
								
								<!-- Tags -->
								{#if project.tags.length > 0}
									<div class="mt-6">
										<h4 class="text-sm font-medium text-gray-400 mb-2">Tags</h4>
										<div class="flex flex-wrap gap-2">
											{#each project.tags as tag}
												<span class="inline-flex items-center px-3 py-1 rounded-full text-sm bg-gray-700 text-gray-300">
													{tag}
												</span>
											{/each}
										</div>
									</div>
								{/if}
							</div>
							
							<!-- Recent Tasks -->
							<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
								<div class="flex items-center justify-between mb-4">
									<h3 class="text-lg font-semibold text-white">Recent Tasks</h3>
									<Button
										variant="ghost"
										size="xs"
										on:click={() => activeTab = 'tasks'}
										class="text-sm text-blue-400 hover:text-blue-300"
									>
										View All →
									</Button>
								</div>
								<TaskBoard 
									tasks={project.tasks.slice(0, 8)} 
									projectId={project.id} 
									compact={true}
									on:createTask={handleCreateTask}
									on:editTask={handleEditTask}
								/>
							</div>
						</div>
						
						<!-- Right Column -->
						<div class="space-y-6">
							<!-- Team Members -->
							<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
								<h3 class="text-lg font-semibold text-white mb-4">Team Members</h3>
								<div class="space-y-3">
									{#each project.members as member}
										<div class="flex items-center space-x-3">
											<div class="relative">
												<div class="w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center">
													{#if member.avatar}
														<img src={member.avatar} alt={member.name} class="w-full h-full rounded-full object-cover" />
													{:else}
														<span class="text-sm font-medium text-white">{member.name.charAt(0)}</span>
													{/if}
												</div>
												<div class="absolute -bottom-0.5 -right-0.5 w-3 h-3 border-2 border-gray-800 rounded-full {
													member.status === 'online' ? 'bg-green-500' :
													member.status === 'away' ? 'bg-yellow-500' :
													'bg-gray-500'
												}"></div>
											</div>
											<div class="flex-1 min-w-0">
												<p class="text-sm font-medium text-white truncate">{member.name}</p>
												<p class="text-xs text-gray-400 capitalize">{member.role}</p>
											</div>
											<div class="text-xs text-gray-500">
												{member.status === 'online' ? 'online' : `${new Date(member.lastSeen).toLocaleTimeString()}`}
											</div>
										</div>
									{/each}
								</div>
								<Button
									variant="secondary"
									size="sm"
									class="mt-4 w-full border-gray-600 text-gray-300 hover:bg-gray-700"
								>
									<UserPlus class="w-4 h-4 mr-2" />
									Invite More Members
								</Button>
							</div>
							
							<!-- Progress Overview -->
							<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
								<h3 class="text-lg font-semibold text-white mb-4">Progress Overview</h3>
								<div class="space-y-4">
									<div>
										<div class="flex items-center justify-between mb-2">
											<span class="text-sm text-gray-400">Overall Progress</span>
											<span class="text-sm font-medium text-white">{project.progress}%</span>
										</div>
										<div class="w-full bg-gray-700 rounded-full h-2">
											<div class="bg-blue-600 h-2 rounded-full transition-all" style="width: {project.progress}%"></div>
										</div>
									</div>
									
									<div class="grid grid-cols-2 gap-4 pt-4">
										<div class="text-center">
											<div class="text-lg font-semibold text-green-400">{project.stats.completedTasks}</div>
											<div class="text-xs text-gray-400">Done</div>
										</div>
										<div class="text-center">
											<div class="text-lg font-semibold text-red-400">{project.stats.overdueTasks}</div>
											<div class="text-xs text-gray-400">Overdue</div>
										</div>
									</div>
								</div>
							</div>
							
							<!-- Milestones -->
							<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
								<div class="flex items-center justify-between mb-4">
									<h3 class="text-lg font-semibold text-white">Upcoming Milestones</h3>
									<Button
										variant="ghost"
										size="xs"
										on:click={() => activeTab = 'timeline'}
										class="text-sm text-blue-400 hover:text-blue-300"
									>
										View All →
									</Button>
								</div>
								{#if project.milestones.length > 0}
									<div class="space-y-3">
										{#each project.milestones.slice(0, 3) as milestone}
											<div class="flex items-center space-x-3 p-3 bg-gray-750 rounded-lg">
												<div class="w-2 h-2 bg-blue-500 rounded-full"></div>
												<div class="flex-1">
													<p class="text-sm font-medium text-white">{milestone.title}</p>
													<p class="text-xs text-gray-400">{formatDate(milestone.dueDate)}</p>
												</div>
												<div class="text-xs text-blue-400">{milestone.progress}%</div>
											</div>
										{/each}
									</div>
								{:else}
									<p class="text-gray-500 text-sm">No milestones yet. Create your first milestone to track progress.</p>
								{/if}
							</div>
						</div>
					</div>
				{:else if activeTab === 'tasks'}
					<!-- Tasks Tab -->
					<TaskBoard 
						tasks={project.tasks} 
						projectId={project.id}
						on:createTask={handleCreateTask}
						on:editTask={handleEditTask}
					/>
				{:else if activeTab === 'timeline'}
					<!-- Timeline Tab -->
					<ProjectTimeline 
						{project}
						on:createMilestone={handleCreateMilestone}
						on:editMilestone={handleEditMilestone}
					/>
				{:else if activeTab === 'chat'}
					<!-- Chat Tab -->
					<ProjectChat {project} />
				{:else if activeTab === 'files'}
					<!-- Files Tab -->
					<ProjectFiles {project} />
				{:else if activeTab === 'analytics'}
					<!-- Analytics Tab -->
					<ProjectAnalytics {project} />
				{/if}
			</div>
		</main>
	{/if}
</div>

<!-- Modals -->
<CreateTaskModal 
	bind:show={showTaskModal} 
	{project}
	task={editingTask}
	defaultStatus={newTaskStatus}
	on:created={handleTaskCreated}
	on:updated={handleTaskUpdated}
/>

<CreateMilestoneModal 
	bind:show={showMilestoneModal} 
	{project}
	milestone={editingMilestone}
	on:created={handleMilestoneCreated}
	on:updated={handleMilestoneUpdated}
/>

<style>
	.bg-gray-750 {
		background-color: #374151;
	}
</style>