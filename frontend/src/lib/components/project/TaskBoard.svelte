<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { ProjectTask } from '$lib/stores/project';
	import { projectStore } from '$lib/stores/project';
	import Button from '../ui/Button.svelte';
	import IconButton from '../ui/IconButton.svelte';
	import { Plus, Edit, Trash } from 'lucide-svelte';
	
	const dispatch = createEventDispatcher();
	
	export let tasks: ProjectTask[] = [];
	export let projectId: string;
	export let compact = false;
	
	let draggedTask: ProjectTask | null = null;
	let draggedFrom: string | null = null;
	
	const columns = [
		{ id: 'todo', title: 'To Do', color: 'bg-gray-600' },
		{ id: 'in_progress', title: 'In Progress', color: 'bg-blue-600' },
		{ id: 'review', title: 'Review', color: 'bg-yellow-600' },
		{ id: 'completed', title: 'Completed', color: 'bg-green-600' }
	];
	
	function getTasksByStatus(status: ProjectTask['status']) {
		return tasks.filter(task => task.status === status);
	}
	
	function getPriorityColor(priority: ProjectTask['priority']) {
		switch (priority) {
			case 'urgent': return 'border-l-red-500';
			case 'high': return 'border-l-orange-500';
			case 'medium': return 'border-l-yellow-500';
			case 'low': return 'border-l-green-500';
			default: return 'border-l-gray-500';
		}
	}
	
	function formatDate(date: Date) {
		return new Intl.DateTimeFormat('en-US', {
			month: 'short',
			day: 'numeric'
		}).format(new Date(date));
	}
	
	function isOverdue(dueDate: Date | undefined, status: ProjectTask['status']) {
		if (!dueDate || status === 'completed') return false;
		return new Date(dueDate) < new Date();
	}
	
	// Drag and Drop handlers
	function handleDragStart(event: DragEvent, task: ProjectTask) {
		if (!event.dataTransfer) return;
		
		draggedTask = task;
		draggedFrom = task.status;
		event.dataTransfer.effectAllowed = 'move';
		event.dataTransfer.setData('text/html', '');
		
		// Add visual feedback
		if (event.target instanceof HTMLElement) {
			event.target.style.opacity = '0.5';
		}
	}
	
	function handleDragEnd(event: DragEvent) {
		// Reset visual feedback
		if (event.target instanceof HTMLElement) {
			event.target.style.opacity = '1';
		}
		
		draggedTask = null;
		draggedFrom = null;
	}
	
	function handleDragOver(event: DragEvent) {
		event.preventDefault();
		if (event.dataTransfer) {
			event.dataTransfer.dropEffect = 'move';
		}
	}
	
	function handleDrop(event: DragEvent, newStatus: ProjectTask['status']) {
		event.preventDefault();
		
		if (!draggedTask || draggedTask.status === newStatus) return;
		
		// Update task status
		projectStore.updateTask(draggedTask.id, { status: newStatus });
		
		dispatch('taskMoved', {
			task: draggedTask,
			from: draggedFrom,
			to: newStatus
		});
	}
	
	function createTask(status: ProjectTask['status']) {
		dispatch('createTask', { status });
	}
	
	function editTask(task: ProjectTask) {
		dispatch('editTask', task);
	}
	
	function deleteTask(task: ProjectTask) {
		if (confirm('Are you sure you want to delete this task?')) {
			projectStore.deleteTask(task.id);
			dispatch('taskDeleted', task);
		}
	}
</script>

<div class="space-y-6">
	<!-- Header -->
	{#if !compact}
		<div class="flex items-center justify-between">
			<h3 class="text-lg font-semibold text-white">Task Board</h3>
			<div class="flex items-center space-x-2">
				<span class="text-sm text-gray-400">{tasks.length} total tasks</span>
				<Button
					on:click={() => createTask('todo')}
					variant="primary"
					size="sm"
					class="px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded transition-colors"
				>
					Add Task
				</Button>
			</div>
		</div>
	{/if}
	
	<!-- Kanban Board -->
	<div class="grid grid-cols-4 gap-4 min-h-[500px]">
		{#each columns as column}
			{@const columnTasks = getTasksByStatus(column.id)}
			<div 
				class="bg-gray-800 rounded-lg border border-gray-700"
				on:dragover={handleDragOver}
				on:drop={(e) => handleDrop(e, column.id)}
			>
				<!-- Column Header -->
				<div class="p-4 border-b border-gray-700">
					<div class="flex items-center justify-between">
						<div class="flex items-center space-x-2">
							<div class="w-3 h-3 rounded-full {column.color}"></div>
							<h4 class="font-medium text-white">{column.title}</h4>
						</div>
						<div class="flex items-center space-x-2">
							<span class="bg-gray-700 text-gray-300 text-xs px-2 py-1 rounded-full">
								{columnTasks.length}
							</span>
							<IconButton
								on:click={() => createTask(column.id)}
								size="sm"
								class="text-gray-400 hover:text-gray-200 transition-colors"
								title="Add task"
							>
								<Plus class="w-4 h-4" />
							</IconButton>
						</div>
					</div>
				</div>
				
				<!-- Tasks -->
				<div class="p-3 space-y-3 min-h-[400px]">
					{#each columnTasks as task}
						<div
							class="bg-gray-700 rounded-lg p-3 border-l-2 {getPriorityColor(task.priority)} cursor-move hover:bg-gray-650 transition-colors group"
							draggable="true"
							on:dragstart={(e) => handleDragStart(e, task)}
							on:dragend={handleDragEnd}
							role="button"
							tabindex="0"
						>
							<!-- Task Header -->
							<div class="flex items-start justify-between mb-2">
								<h5 class="text-sm font-medium text-white line-clamp-2 flex-1">
									{task.title}
								</h5>
								<div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
									<IconButton
										on:click={(e) => { e.stopPropagation(); editTask(task); }}
										size="sm"
										class="p-1 text-gray-400 hover:text-gray-200 transition-colors"
										title="Edit task"
									>
										<Edit class="w-3 h-3" />
									</IconButton>
									<IconButton
										on:click={(e) => { e.stopPropagation(); deleteTask(task); }}
										size="sm"
										class="p-1 text-gray-400 hover:text-red-400 transition-colors"
										title="Delete task"
									>
										<Trash class="w-3 h-3" />
									</IconButton>
								</div>
							</div>
							
							<!-- Task Description -->
							{#if task.description}
								<p class="text-xs text-gray-400 mb-3 line-clamp-2">{task.description}</p>
							{/if}
							
							<!-- Task Meta -->
							<div class="flex items-center justify-between text-xs">
								<div class="flex items-center space-x-2">
									<!-- Priority -->
									<span class="inline-flex items-center px-2 py-1 rounded-full font-medium
										{task.priority === 'urgent' ? 'text-red-400 bg-red-400/10' :
										 task.priority === 'high' ? 'text-orange-400 bg-orange-400/10' :
										 task.priority === 'medium' ? 'text-yellow-400 bg-yellow-400/10' :
										 'text-green-400 bg-green-400/10'}">
										{task.priority}
									</span>
									
									<!-- Tags -->
									{#if task.tags.length > 0}
										<div class="flex space-x-1">
											{#each task.tags.slice(0, 2) as tag}
												<span class="text-gray-500 bg-gray-600 px-1 py-0.5 rounded text-xs">
													{tag}
												</span>
											{/each}
											{#if task.tags.length > 2}
												<span class="text-gray-500">+{task.tags.length - 2}</span>
											{/if}
										</div>
									{/if}
								</div>
								
								<!-- Due Date -->
								{#if task.dueDate}
									<div class="flex items-center space-x-1">
										{#if isOverdue(task.dueDate, task.status)}
											<svg class="w-3 h-3 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
											</svg>
											<span class="text-red-400 font-medium">
												{formatDate(task.dueDate)}
											</span>
										{:else}
											<svg class="w-3 h-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
											</svg>
											<span class="text-gray-400">
												{formatDate(task.dueDate)}
											</span>
										{/if}
									</div>
								{/if}
							</div>
							
							<!-- Assignee -->
							{#if task.assignedTo}
								<div class="flex items-center mt-2">
									<div class="w-5 h-5 bg-blue-600 rounded-full flex items-center justify-center">
										<span class="text-xs font-medium text-white">
											{task.assignedTo.charAt(0).toUpperCase()}
										</span>
									</div>
									<span class="text-xs text-gray-400 ml-2">Assigned to {task.assignedTo}</span>
								</div>
							{/if}
						</div>
					{/each}
					
					<!-- Empty State -->
					{#if columnTasks.length === 0}
						<div class="text-center py-8">
							<svg class="w-8 h-8 text-gray-600 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
							</svg>
							<p class="text-xs text-gray-500">No tasks</p>
							<Button
								on:click={() => createTask(column.id)}
								variant="ghost"
								size="xs"
								class="text-blue-400 hover:text-blue-300 mt-1"
							>
								Add first task
							</Button>
						</div>
					{/if}
				</div>
			</div>
		{/each}
	</div>
	
	<!-- Summary Stats -->
	{#if !compact}
		<div class="bg-gray-800 rounded-lg border border-gray-700 p-4">
			<div class="grid grid-cols-4 gap-4 text-center">
				<div>
					<div class="text-2xl font-semibold text-white">{getTasksByStatus('todo').length}</div>
					<div class="text-sm text-gray-400">To Do</div>
				</div>
				<div>
					<div class="text-2xl font-semibold text-blue-400">{getTasksByStatus('in_progress').length}</div>
					<div class="text-sm text-gray-400">In Progress</div>
				</div>
				<div>
					<div class="text-2xl font-semibold text-yellow-400">{getTasksByStatus('review').length}</div>
					<div class="text-sm text-gray-400">In Review</div>
				</div>
				<div>
					<div class="text-2xl font-semibold text-green-400">{getTasksByStatus('completed').length}</div>
					<div class="text-sm text-gray-400">Completed</div>
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.bg-gray-650 {
		background-color: #4B5563;
	}
	
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
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