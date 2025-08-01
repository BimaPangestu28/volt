<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { projectStore, type Project, type ProjectTask } from '$lib/stores/project';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import { X } from 'lucide-svelte';
	
	const dispatch = createEventDispatcher();
	
	export let show = false;
	export let project: Project;
	export let task: ProjectTask | null = null;
	export let defaultStatus: ProjectTask['status'] = 'todo';
	
	let isEditing = false;
	let taskTitle = '';
	let taskDescription = '';
	let taskStatus: ProjectTask['status'] = 'todo';
	let taskPriority: ProjectTask['priority'] = 'medium';
	let taskAssignedTo = '';
	let taskDueDate = '';
	let taskTags: string[] = [];
	let newTag = '';
	let isCreating = false;
	let error = '';
	
	const statusOptions = [
		{ value: 'todo', label: 'To Do', color: 'text-gray-400' },
		{ value: 'in_progress', label: 'In Progress', color: 'text-blue-400' },
		{ value: 'review', label: 'Review', color: 'text-yellow-400' },
		{ value: 'completed', label: 'Completed', color: 'text-green-400' }
	];
	
	const priorityOptions = [
		{ value: 'low', label: 'Low', color: 'text-green-400' },
		{ value: 'medium', label: 'Medium', color: 'text-yellow-400' },
		{ value: 'high', label: 'High', color: 'text-orange-400' },
		{ value: 'urgent', label: 'Urgent', color: 'text-red-400' }
	];
	
	// Watch for changes to show/task props
	$: if (show) {
		initializeForm();
	}
	
	function initializeForm() {
		isEditing = !!task;
		
		if (task) {
			// Edit mode - populate with existing task data
			taskTitle = task.title;
			taskDescription = task.description || '';
			taskStatus = task.status;
			taskPriority = task.priority;
			taskAssignedTo = task.assignedTo || '';
			taskDueDate = task.dueDate ? new Date(task.dueDate).toISOString().split('T')[0] : '';
			taskTags = [...task.tags];
		} else {
			// Create mode - reset form
			taskTitle = '';
			taskDescription = '';
			taskStatus = defaultStatus;
			taskPriority = 'medium';
			taskAssignedTo = '';
			taskDueDate = '';
			taskTags = [];
		}
		
		newTag = '';
		error = '';
	}
	
	function addTag() {
		if (newTag.trim() && !taskTags.includes(newTag.trim())) {
			taskTags = [...taskTags, newTag.trim()];
			newTag = '';
		}
	}
	
	function removeTag(tag: string) {
		taskTags = taskTags.filter(t => t !== tag);
	}
	
	async function saveTask() {
		if (!taskTitle.trim()) {
			error = 'Task title is required';
			return;
		}
		
		isCreating = true;
		error = '';
		
		try {
			const now = new Date();
			
			if (isEditing && task) {
				// Update existing task
				const updatedTask: Partial<ProjectTask> = {
					title: taskTitle.trim(),
					description: taskDescription.trim() || undefined,
					status: taskStatus,
					priority: taskPriority,
					assignedTo: taskAssignedTo.trim() || undefined,
					dueDate: taskDueDate ? new Date(taskDueDate) : undefined,
					tags: taskTags,
					updatedAt: now
				};
				
				projectStore.updateTask(task.id, updatedTask);
				dispatch('updated', { task: { ...task, ...updatedTask } });
			} else {
				// Create new task
				const newTask: ProjectTask = {
					id: `task_${Date.now()}`,
					title: taskTitle.trim(),
					description: taskDescription.trim() || undefined,
					status: taskStatus,
					priority: taskPriority,
					assignedTo: taskAssignedTo.trim() || undefined,
					createdAt: now,
					updatedAt: now,
					dueDate: taskDueDate ? new Date(taskDueDate) : undefined,
					tags: taskTags,
					projectId: project.id
				};
				
				projectStore.addTask(newTask);
				dispatch('created', { task: newTask });
			}
			
			closeModal();
		} catch (err: any) {
			error = err.message || `Failed to ${isEditing ? 'update' : 'create'} task`;
		} finally {
			isCreating = false;
		}
	}
	
	function closeModal() {
		show = false;
		resetForm();
	}
	
	function resetForm() {
		taskTitle = '';
		taskDescription = '';
		taskStatus = 'todo';
		taskPriority = 'medium';
		taskAssignedTo = '';
		taskDueDate = '';
		taskTags = [];
		newTag = '';
		error = '';
	}
</script>

{#if show}
	<div class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center p-4 z-50">
		<div class="bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-hidden">
			<!-- Header -->
			<div class="px-6 py-4 border-b border-gray-700">
				<div class="flex items-center justify-between">
					<h3 class="text-xl font-semibold text-white">
						{isEditing ? 'Edit Task' : 'Create New Task'}
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
					<!-- Task Title -->
					<div>
						<label for="taskTitle" class="block text-sm font-medium text-gray-300 mb-2">
							Task Title *
						</label>
						<input
							id="taskTitle"
							type="text"
							bind:value={taskTitle}
							placeholder="Enter task title..."
							class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>
					
					<!-- Task Description -->
					<div>
						<label for="taskDescription" class="block text-sm font-medium text-gray-300 mb-2">
							Description
						</label>
						<textarea
							id="taskDescription"
							bind:value={taskDescription}
							placeholder="Describe the task in detail..."
							rows="3"
							class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						></textarea>
					</div>
					
					<!-- Status and Priority -->
					<div class="grid grid-cols-2 gap-4">
						<div>
							<label for="taskStatus" class="block text-sm font-medium text-gray-300 mb-2">
								Status
							</label>
							<select
								id="taskStatus"
								bind:value={taskStatus}
								class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							>
								{#each statusOptions as option}
									<option value={option.value}>{option.label}</option>
								{/each}
							</select>
						</div>
						
						<div>
							<label for="taskPriority" class="block text-sm font-medium text-gray-300 mb-2">
								Priority
							</label>
							<select
								id="taskPriority"
								bind:value={taskPriority}
								class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							>
								{#each priorityOptions as option}
									<option value={option.value}>{option.label}</option>
								{/each}
							</select>
						</div>
					</div>
					
					<!-- Assignee and Due Date -->
					<div class="grid grid-cols-2 gap-4">
						<div>
							<label for="taskAssignedTo" class="block text-sm font-medium text-gray-300 mb-2">
								Assign To
							</label>
							<select
								id="taskAssignedTo"
								bind:value={taskAssignedTo}
								class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							>
								<option value="">Unassigned</option>
								{#each project.members as member}
									<option value={member.name}>{member.name} ({member.role})</option>
								{/each}
							</select>
						</div>
						
						<div>
							<label for="taskDueDate" class="block text-sm font-medium text-gray-300 mb-2">
								Due Date
							</label>
							<input
								id="taskDueDate"
								type="date"
								bind:value={taskDueDate}
								class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							/>
						</div>
					</div>
					
					<!-- Tags -->
					<div>
						<label class="block text-sm font-medium text-gray-300 mb-2">
							Tags
						</label>
						<div class="flex space-x-2 mb-3">
							<input
								type="text"
								bind:value={newTag}
								placeholder="Add a tag..."
								class="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
								on:keydown={(e) => e.key === 'Enter' && addTag()}
							/>
							<Button variant="primary" size="sm" on:click={addTag}>
								Add
							</Button>
						</div>
						{#if taskTags.length > 0}
							<div class="flex flex-wrap gap-2">
								{#each taskTags as tag}
									<div class="flex items-center space-x-2 bg-gray-700 px-3 py-1 rounded-full">
										<span class="text-sm text-white">{tag}</span>
										<IconButton variant="ghost" size="xs" on:click={() => removeTag(tag)}>
											<X class="btn-icon-xs" />
										</IconButton>
									</div>
								{/each}
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
				
				<Button variant="primary" size="sm" disabled={isCreating || !taskTitle.trim()} loading={isCreating} on:click={saveTask}>
					{#if isCreating}
						{isEditing ? 'Updating...' : 'Creating...'}
					{:else}
						{isEditing ? 'Update Task' : 'Create Task'}
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
</style>