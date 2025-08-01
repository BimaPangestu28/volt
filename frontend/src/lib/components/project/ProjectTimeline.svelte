<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Project, ProjectMilestone } from '$lib/stores/project';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import { Plus, Edit, Trash2 } from 'lucide-svelte';
	
	const dispatch = createEventDispatcher();
	
	export let project: Project;
	export let compact = false;
	
	let selectedMilestone: ProjectMilestone | null = null;
	let showCreateModal = false;
	
	function formatDate(date: Date) {
		return new Intl.DateTimeFormat('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		}).format(new Date(date));
	}
	
	function formatDateShort(date: Date) {
		return new Intl.DateTimeFormat('en-US', {
			month: 'short',
			day: 'numeric'
		}).format(new Date(date));
	}
	
	function getMilestoneStatus(milestone: ProjectMilestone): 'completed' | 'current' | 'upcoming' | 'overdue' {
		const today = new Date();
		const dueDate = new Date(milestone.dueDate);
		
		if (milestone.status === 'completed') return 'completed';
		if (dueDate < today && milestone.status !== 'completed') return 'overdue';
		if (milestone.status === 'in_progress') return 'current';
		return 'upcoming';
	}
	
	function getStatusColor(status: string) {
		switch (status) {
			case 'completed': return 'bg-green-500';
			case 'current': return 'bg-blue-500';
			case 'overdue': return 'bg-red-500';
			case 'upcoming': return 'bg-gray-400';
			default: return 'bg-gray-400';
		}
	}
	
	function getStatusBorderColor(status: string) {
		switch (status) {
			case 'completed': return 'border-green-500';
			case 'current': return 'border-blue-500';
			case 'overdue': return 'border-red-500';
			case 'upcoming': return 'border-gray-400';
			default: return 'border-gray-400';
		}
	}
	
	function getDaysBetween(start: Date, end: Date): number {
		const diffTime = new Date(end).getTime() - new Date(start).getTime();
		return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
	}
	
	function createMilestone() {
		showCreateModal = true;
	}
	
	function editMilestone(milestone: ProjectMilestone) {
		selectedMilestone = milestone;
		dispatch('editMilestone', milestone);
	}
	
	function deleteMilestone(milestone: ProjectMilestone) {
		if (confirm('Are you sure you want to delete this milestone?')) {
			dispatch('deleteMilestone', milestone);
		}
	}
	
	// Sort milestones by due date
	$: sortedMilestones = [...project.milestones].sort((a, b) => 
		new Date(a.dueDate).getTime() - new Date(b.dueDate).getTime()
	);
	
	// Calculate timeline positions
	$: timelineData = (() => {
		if (sortedMilestones.length === 0) return [];
		
		const startDate = new Date(project.startDate);
		const endDate = project.dueDate ? new Date(project.dueDate) : 
			new Date(Math.max(...sortedMilestones.map(m => new Date(m.dueDate).getTime())));
		
		const totalDays = getDaysBetween(startDate, endDate);
		
		return sortedMilestones.map(milestone => {
			const daysSinceStart = getDaysBetween(startDate, milestone.dueDate);
			const position = totalDays > 0 ? (daysSinceStart / totalDays) * 100 : 0;
			
			return {
				...milestone,
				position: Math.max(0, Math.min(100, position)),
				status: getMilestoneStatus(milestone)
			};
		});
	})();
</script>

<div class="space-y-6">
	<!-- Header -->
	{#if !compact}
		<div class="flex items-center justify-between">
			<h3 class="text-lg font-semibold text-white">Project Timeline</h3>
			<Button
				on:click={createMilestone}
				variant="primary"
				size="sm"
				class="flex items-center space-x-2"
			>
				<Plus class="btn-icon-sm" />
				<span>Add Milestone</span>
			</Button>
		</div>
	{/if}
	
	{#if sortedMilestones.length === 0}
		<!-- Empty State -->
		<div class="bg-gray-800 rounded-lg border border-gray-700 p-8 text-center">
			<svg class="w-12 h-12 text-gray-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
			</svg>
			<h4 class="text-lg font-medium text-white mb-2">No Milestones Yet</h4>
			<p class="text-gray-400 mb-4">Create milestones to track important project deliverables and deadlines.</p>
			<Button
				on:click={createMilestone}
				variant="primary"
				size="md"
			>
				Create First Milestone
			</Button>
		</div>
	{:else}
		<!-- Timeline Visualization -->
		<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
			<!-- Timeline Header -->
			<div class="flex items-center justify-between mb-6">
				<div class="flex items-center space-x-4">
					<div class="text-sm text-gray-400">
						<span class="font-medium text-white">Start:</span> {formatDate(project.startDate)}
					</div>
					{#if project.dueDate}
						<div class="text-sm text-gray-400">
							<span class="font-medium text-white">Due:</span> {formatDate(project.dueDate)}
						</div>
					{/if}
				</div>
				<div class="text-sm text-gray-400">
					{sortedMilestones.length} milestone{sortedMilestones.length !== 1 ? 's' : ''}
				</div>
			</div>
			
			<!-- Timeline Line -->
			<div class="relative">
				<!-- Background line -->
				<div class="absolute inset-x-0 top-4 h-0.5 bg-gray-600"></div>
				
				<!-- Progress line -->
				<div 
					class="absolute inset-x-0 top-4 h-0.5 bg-blue-500 transition-all duration-500"
					style="width: {project.progress}%"
				></div>
				
				<!-- Milestones -->
				<div class="relative flex">
					{#each timelineData as milestone, index}
						<div 
							class="absolute flex flex-col items-center transform -translate-x-1/2 cursor-pointer group"
							style="left: {milestone.position}%"
							on:click={() => editMilestone(milestone)}
							on:keydown={(e) => e.key === 'Enter' && editMilestone(milestone)}
							role="button"
							tabindex="0"
						>
							<!-- Milestone dot -->
							<div class="relative">
								<div class="w-8 h-8 rounded-full border-4 {getStatusBorderColor(milestone.status)} bg-gray-800 flex items-center justify-center group-hover:scale-110 transition-transform">
									{#if milestone.status === 'completed'}
										<svg class="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
										</svg>
									{:else if milestone.status === 'current'}
										<div class="w-3 h-3 bg-blue-500 rounded-full animate-pulse"></div>
									{:else if milestone.status === 'overdue'}
										<svg class="w-4 h-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
										</svg>
									{:else}
										<div class="w-2 h-2 bg-gray-400 rounded-full"></div>
									{/if}
								</div>
								
								<!-- Status indicator -->
								<div class="absolute -top-1 -right-1 w-3 h-3 {getStatusColor(milestone.status)} rounded-full border-2 border-gray-800"></div>
							</div>
							
							<!-- Milestone info tooltip -->
							<div class="absolute top-12 left-1/2 transform -translate-x-1/2 bg-gray-900 border border-gray-600 rounded-lg p-3 min-w-48 opacity-0 group-hover:opacity-100 transition-opacity z-10 pointer-events-none">
								<div class="text-sm font-medium text-white mb-1">{milestone.title}</div>
								{#if milestone.description}
									<div class="text-xs text-gray-400 mb-2">{milestone.description}</div>
								{/if}
								<div class="flex items-center justify-between text-xs">
									<span class="text-gray-500">Due: {formatDateShort(milestone.dueDate)}</span>
									<span class="text-{milestone.status === 'completed' ? 'green' : milestone.status === 'overdue' ? 'red' : milestone.status === 'current' ? 'blue' : 'gray'}-400 font-medium capitalize">
										{milestone.status}
									</span>
								</div>
								{#if milestone.progress > 0}
									<div class="mt-2">
										<div class="flex items-center justify-between mb-1">
											<span class="text-xs text-gray-400">Progress</span>
											<span class="text-xs text-white">{milestone.progress}%</span>
										</div>
										<div class="w-full bg-gray-700 rounded-full h-1">
											<div class="bg-blue-500 h-1 rounded-full transition-all" style="width: {milestone.progress}%"></div>
										</div>
									</div>
								{/if}
							</div>
						</div>
					{/each}
				</div>
				
				<!-- Timeline labels -->
				<div class="mt-16 space-y-3">
					{#each timelineData as milestone, index}
						<div class="flex items-center justify-between p-3 bg-gray-750 rounded-lg hover:bg-gray-700 transition-colors group">
							<div class="flex items-center space-x-3">
								<div class="w-3 h-3 {getStatusColor(milestone.status)} rounded-full"></div>
								<div>
									<h4 class="text-sm font-medium text-white">{milestone.title}</h4>
									{#if milestone.description}
										<p class="text-xs text-gray-400 mt-1">{milestone.description}</p>
									{/if}
								</div>
							</div>
							
							<div class="flex items-center space-x-4">
								<div class="text-xs text-gray-400">
									{formatDateShort(milestone.dueDate)}
								</div>
								{#if milestone.progress > 0}
									<div class="text-xs text-blue-400">
										{milestone.progress}%
									</div>
								{/if}
								<div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
									<IconButton
										on:click={(e) => { e.stopPropagation(); editMilestone(milestone); }}
										variant="ghost"
										size="xs"
										title="Edit milestone"
										class="text-gray-400 hover:text-gray-200"
									>
										<Edit class="btn-icon-xs" />
									</IconButton>
									<IconButton
										on:click={(e) => { e.stopPropagation(); deleteMilestone(milestone); }}
										variant="ghost"
										size="xs"
										title="Delete milestone"
										class="text-gray-400 hover:text-red-400"
									>
										<Trash2 class="btn-icon-xs" />
									</IconButton>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
		
		<!-- Milestone Statistics -->
		{#if !compact}
			<div class="grid grid-cols-4 gap-4">
				<div class="bg-gray-800 rounded-lg border border-gray-700 p-4 text-center">
					<div class="text-2xl font-semibold text-white">
						{sortedMilestones.filter(m => getMilestoneStatus(m) === 'completed').length}
					</div>
					<div class="text-sm text-green-400">Completed</div>
				</div>
				<div class="bg-gray-800 rounded-lg border border-gray-700 p-4 text-center">
					<div class="text-2xl font-semibold text-white">
						{sortedMilestones.filter(m => getMilestoneStatus(m) === 'current').length}
					</div>
					<div class="text-sm text-blue-400">In Progress</div>
				</div>
				<div class="bg-gray-800 rounded-lg border border-gray-700 p-4 text-center">
					<div class="text-2xl font-semibold text-white">
						{sortedMilestones.filter(m => getMilestoneStatus(m) === 'upcoming').length}
					</div>
					<div class="text-sm text-gray-400">Upcoming</div>
				</div>
				<div class="bg-gray-800 rounded-lg border border-gray-700 p-4 text-center">
					<div class="text-2xl font-semibold text-white">
						{sortedMilestones.filter(m => getMilestoneStatus(m) === 'overdue').length}
					</div>
					<div class="text-sm text-red-400">Overdue</div>
				</div>
			</div>
		{/if}
	{/if}
</div>

<style>
	.bg-gray-750 {
		background-color: #374151;
	}
	
	.btn-icon-xs {
		width: 12px;
		height: 12px;
	}
	
	.btn-icon-sm {
		width: 16px;
		height: 16px;
	}
</style>