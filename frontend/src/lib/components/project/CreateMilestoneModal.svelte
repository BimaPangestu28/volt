<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { projectStore, type Project, type ProjectMilestone } from '$lib/stores/project';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import { X } from 'lucide-svelte';
	
	const dispatch = createEventDispatcher();
	
	export let show = false;
	export let project: Project;
	export let milestone: ProjectMilestone | null = null;
	
	let isEditing = false;
	let milestoneTitle = '';
	let milestoneDescription = '';
	let milestoneDueDate = '';
	let milestoneProgress = 0;
	let milestoneStatus: ProjectMilestone['status'] = 'upcoming';
	let selectedTasks: string[] = [];
	let isCreating = false;
	let error = '';
	
	const statusOptions = [
		{ value: 'upcoming', label: 'Upcoming', color: 'text-gray-400' },
		{ value: 'in_progress', label: 'In Progress', color: 'text-blue-400' },
		{ value: 'completed', label: 'Completed', color: 'text-green-400' },
		{ value: 'overdue', label: 'Overdue', color: 'text-red-400' }
	];
	
	// Watch for changes to show/milestone props
	$: if (show) {
		initializeForm();
	}
	
	function initializeForm() {
		isEditing = !!milestone;
		
		if (milestone) {
			// Edit mode - populate with existing milestone data
			milestoneTitle = milestone.title;
			milestoneDescription = milestone.description || '';
			milestoneDueDate = new Date(milestone.dueDate).toISOString().split('T')[0];
			milestoneProgress = milestone.progress;
			milestoneStatus = milestone.status;
			selectedTasks = [...milestone.tasks];
		} else {
			// Create mode - reset form
			milestoneTitle = '';
			milestoneDescription = '';
			milestoneDueDate = '';
			milestoneProgress = 0;
			milestoneStatus = 'upcoming';
			selectedTasks = [];
		}
		
		error = '';
	}
	
	function toggleTask(taskId: string) {
		if (selectedTasks.includes(taskId)) {
			selectedTasks = selectedTasks.filter(id => id !== taskId);
		} else {
			selectedTasks = [...selectedTasks, taskId];
		}
	}
	
	function calculateProgressFromTasks() {
		if (selectedTasks.length === 0) return 0;
		
		const tasksData = project.tasks.filter(task => selectedTasks.includes(task.id));
		const completedTasks = tasksData.filter(task => task.status === 'completed').length;
		
		return Math.round((completedTasks / tasksData.length) * 100);
	}
	
	function updateProgressFromTasks() {
		milestoneProgress = calculateProgressFromTasks();
	}
	
	async function saveMilestone() {
		if (!milestoneTitle.trim()) {
			error = 'Milestone title is required';
			return;
		}
		
		if (!milestoneDueDate) {
			error = 'Due date is required';
			return;
		}
		
		isCreating = true;
		error = '';
		
		try {
			if (isEditing && milestone) {
				// Update existing milestone
				const updatedMilestone: Partial<ProjectMilestone> = {
					title: milestoneTitle.trim(),
					description: milestoneDescription.trim() || undefined,
					dueDate: new Date(milestoneDueDate),
					progress: milestoneProgress,
					status: milestoneStatus,
					tasks: selectedTasks
				};
				
				projectStore.updateMilestone(milestone.id, updatedMilestone);
				dispatch('updated', { milestone: { ...milestone, ...updatedMilestone } });
			} else {
				// Create new milestone
				const newMilestone: ProjectMilestone = {
					id: `milestone_${Date.now()}`,
					title: milestoneTitle.trim(),
					description: milestoneDescription.trim() || undefined,
					dueDate: new Date(milestoneDueDate),
					status: milestoneStatus,
					progress: milestoneProgress,
					projectId: project.id,
					tasks: selectedTasks
				};
				
				projectStore.addMilestone(newMilestone);
				dispatch('created', { milestone: newMilestone });
			}
			
			closeModal();
		} catch (err: any) {
			error = err.message || `Failed to ${isEditing ? 'update' : 'create'} milestone`;
		} finally {
			isCreating = false;
		}
	}
	
	function closeModal() {
		show = false;
		resetForm();
	}
	
	function resetForm() {
		milestoneTitle = '';
		milestoneDescription = '';
		milestoneDueDate = '';
		milestoneProgress = 0;
		milestoneStatus = 'upcoming';
		selectedTasks = [];
		error = '';
	}
	
	$: autoProgress = calculateProgressFromTasks();
</script>

{#if show}
	<div class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center p-4 z-50">
		<div class="bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-hidden">
			<!-- Header -->
			<div class="px-6 py-4 border-b border-gray-700">
				<div class="flex items-center justify-between">
					<h3 class="text-xl font-semibold text-white">
						{isEditing ? 'Edit Milestone' : 'Create New Milestone'}
					</h3>
					<IconButton variant="ghost" size="sm" on:click={closeModal}>
						<X class="w-3 h-3" />
					</IconButton>
				</div>
			</div>
			
			<div class="px-6 py-4 max-h-[60vh] overflow-y-auto">
				{#if error}
					<div class="bg-red-900/50 border border-red-700 rounded-md p-3 mb-4">
						<p class="text-sm text-red-200">{error}</p>
					</div>
				{/if}
				
				<div class="space-y-4">
					<!-- Milestone Title -->
					<div>
						<label for="milestoneTitle" class="block text-sm font-medium text-gray-300 mb-2">
							Milestone Title *
						</label>
						<input
							id="milestoneTitle"
							type="text"
							bind:value={milestoneTitle}
							placeholder="Enter milestone title..."
							class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>
					
					<!-- Milestone Description -->
					<div>
						<label for="milestoneDescription" class="block text-sm font-medium text-gray-300 mb-2">
							Description
						</label>
						<textarea
							id="milestoneDescription"
							bind:value={milestoneDescription}
							placeholder="Describe the milestone and its deliverables..."
							rows="3"
							class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						></textarea>
					</div>
					
					<!-- Due Date and Status -->
					<div class="grid grid-cols-2 gap-4">
						<div>
							<label for="milestoneDueDate" class="block text-sm font-medium text-gray-300 mb-2">
								Due Date *
							</label>
							<input
								id="milestoneDueDate"
								type="date"
								bind:value={milestoneDueDate}
								class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							/>
						</div>
						
						<div>
							<label for="milestoneStatus" class="block text-sm font-medium text-gray-300 mb-2">
								Status
							</label>
							<select
								id="milestoneStatus"
								bind:value={milestoneStatus}
								class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							>
								{#each statusOptions as option}
									<option value={option.value}>{option.label}</option>
								{/each}
							</select>
						</div>
					</div>
					
					<!-- Progress -->
					<div>
						<div class="flex items-center justify-between mb-2">
							<label for="milestoneProgress" class="block text-sm font-medium text-gray-300">
								Progress
							</label>
							{#if selectedTasks.length > 0}
								<Button variant="ghost" size="xs" on:click={updateProgressFromTasks}>
									Auto-calculate from tasks ({autoProgress}%)
								</Button>
							{/if}
						</div>
						<div class="flex items-center space-x-4">
							<input
								id="milestoneProgress"
								type="range"
								min="0"
								max="100"
								bind:value={milestoneProgress}
								class="flex-1 h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer"
							/>
							<div class="flex items-center space-x-2 min-w-16">
								<input
									type="number"
									min="0"
									max="100"
									bind:value={milestoneProgress}
									class="w-16 px-2 py-1 bg-gray-700 border border-gray-600 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
								/>
								<span class="text-sm text-gray-400">%</span>
							</div>
						</div>
						<div class="mt-2 w-full bg-gray-700 rounded-full h-2">
							<div class="bg-blue-600 h-2 rounded-full transition-all" style="width: {milestoneProgress}%"></div>
						</div>
					</div>
					
					<!-- Associated Tasks -->
					<div>
						<label class="block text-sm font-medium text-gray-300 mb-2">
							Associated Tasks
						</label>
						<p class="text-xs text-gray-400 mb-3">
							Select tasks that contribute to this milestone. Progress will be calculated based on completed tasks.
						</p>
						
						{#if project.tasks.length === 0}
							<div class="text-center py-4 bg-gray-750 rounded-lg">
								<p class="text-sm text-gray-400">No tasks available yet. Create some tasks first.</p>
							</div>
						{:else}
							<div class="max-h-48 overflow-y-auto bg-gray-750 rounded-lg p-3">
								<div class="space-y-2">
									{#each project.tasks as task}
										<label class="flex items-center space-x-3 p-2 rounded hover:bg-gray-700 cursor-pointer">
											<input
												type="checkbox"
												checked={selectedTasks.includes(task.id)}
												on:change={() => toggleTask(task.id)}
												class="w-4 h-4 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500 focus:ring-2"
											/>
											<div class="flex-1 min-w-0">
												<div class="flex items-center space-x-2">
													<span class="text-sm font-medium text-white truncate">{task.title}</span>
													<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium {
														task.priority === 'urgent' ? 'text-red-400 bg-red-400/10' :
														task.priority === 'high' ? 'text-orange-400 bg-orange-400/10' :
														task.priority === 'medium' ? 'text-yellow-400 bg-yellow-400/10' :
														'text-green-400 bg-green-400/10'
													}">
														{task.priority}
													</span>
												</div>
												{#if task.description}
													<p class="text-xs text-gray-400 mt-1 truncate">{task.description}</p>
												{/if}
											</div>
											<div class="flex items-center space-x-2">
												<div class="w-2 h-2 rounded-full {
													task.status === 'completed' ? 'bg-green-500' :
													task.status === 'in_progress' ? 'bg-blue-500' :
													task.status === 'review' ? 'bg-yellow-500' :
													'bg-gray-500'
												}"></div>
												<span class="text-xs text-gray-400 capitalize">{task.status.replace('_', ' ')}</span>
											</div>
										</label>
									{/each}
								</div>
							</div>
						{/if}
						
						{#if selectedTasks.length > 0}
							<div class="mt-3 text-sm text-gray-400">
								{selectedTasks.length} task{selectedTasks.length !== 1 ? 's' : ''} selected
								{#if autoProgress !== milestoneProgress}
									• Auto-calculated progress: {autoProgress}%
								{/if}
							</div>
						{/if}
					</div>
				</div>
			</div>
			
			<!-- Footer -->
			<div class="px-6 py-4 bg-gray-750 border-t border-gray-700 flex justify-end space-x-3">
				<Button variant="ghost" size="sm" disabled={isCreating} on:click={closeModal}>
					Cancel
				</Button>
				
				<Button variant="primary" size="sm" disabled={isCreating || !milestoneTitle.trim() || !milestoneDueDate} loading={isCreating} on:click={saveMilestone}>
					{#if isCreating}
						{isEditing ? 'Updating...' : 'Creating...'}
					{:else}
						{isEditing ? 'Update Milestone' : 'Create Milestone'}
					{/if}
				</Button>
			</div>
		</div>
	</div>
{/if}

<style>
	.bg-gray-750 {
		background-color: #374151;
	}
	
	/* Custom range slider styles */
	input[type="range"]::-webkit-slider-thumb {
		appearance: none;
		height: 20px;
		width: 20px;
		border-radius: 50%;
		background: #2563eb;
		cursor: pointer;
		border: 2px solid #1f2937;
	}
	
	input[type="range"]::-moz-range-thumb {
		height: 20px;
		width: 20px;
		border-radius: 50%;
		background: #2563eb;
		cursor: pointer;
		border: 2px solid #1f2937;
	}
</style>