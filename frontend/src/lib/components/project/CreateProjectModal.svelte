<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { projectStore, defaultTemplates } from '$lib/stores/project';
	import type { Project, ProjectTemplate } from '$lib/stores/project';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import { X } from 'lucide-svelte';
	
	const dispatch = createEventDispatcher();
	
	export let show = false;
	
	let currentStep = 1;
	let projectName = '';
	let projectDescription = '';
	let projectType: Project['type'] = 'api';
	let projectPriority: Project['priority'] = 'medium';
	let selectedTemplate: ProjectTemplate | null = null;
	let startDate = new Date().toISOString().split('T')[0];
	let dueDate = '';
	let teamMembers: string[] = [];
	let newMemberEmail = '';
	let projectTags: string[] = [];
	let newTag = '';
	let isCreating = false;
	let error = '';
	
	// Auto-sync and privacy settings
	let autoSync = true;
	let isPublic = false;
	let allowGuestAccess = false;
	let notifications = true;
	
	const projectTypes = [
		{ value: 'api', label: 'API Development', description: 'RESTful APIs, GraphQL, webhooks' },
		{ value: 'web', label: 'Web Application', description: 'Frontend, full-stack web apps' },
		{ value: 'mobile', label: 'Mobile App', description: 'iOS, Android, cross-platform' },
		{ value: 'desktop', label: 'Desktop App', description: 'Native desktop applications' },
		{ value: 'data', label: 'Data Project', description: 'Analytics, ML, data processing' },
		{ value: 'other', label: 'Other', description: 'Custom project type' }
	];
	
	const priorities = [
		{ value: 'low', label: 'Low', color: 'text-green-400' },
		{ value: 'medium', label: 'Medium', color: 'text-yellow-400' },
		{ value: 'high', label: 'High', color: 'text-orange-400' },
		{ value: 'urgent', label: 'Urgent', color: 'text-red-400' }
	];
	
	$: availableTemplates = defaultTemplates.filter(t => t.type === projectType);
	$: selectedTypeInfo = projectTypes.find(t => t.value === projectType);
	
	function nextStep() {
		if (currentStep < 4) {
			currentStep++;
		}
	}
	
	function prevStep() {
		if (currentStep > 1) {
			currentStep--;
		}
	}
	
	function selectTemplate(template: ProjectTemplate) {
		selectedTemplate = template;
		projectTags = [...template.tags];
		if (template.estimatedDuration && !dueDate) {
			const estimated = new Date();
			estimated.setDate(estimated.getDate() + template.estimatedDuration);
			dueDate = estimated.toISOString().split('T')[0];
		}
	}
	
	function addTeamMember() {
		if (newMemberEmail.trim() && !teamMembers.includes(newMemberEmail.trim())) {
			teamMembers = [...teamMembers, newMemberEmail.trim()];
			newMemberEmail = '';
		}
	}
	
	function removeTeamMember(email: string) {
		teamMembers = teamMembers.filter(m => m !== email);
	}
	
	function addTag() {
		if (newTag.trim() && !projectTags.includes(newTag.trim())) {
			projectTags = [...projectTags, newTag.trim()];
			newTag = '';
		}
	}
	
	function removeTag(tag: string) {
		projectTags = projectTags.filter(t => t !== tag);
	}
	
	async function createProject() {
		if (!projectName.trim()) {
			error = 'Project name is required';
			return;
		}
		
		isCreating = true;
		error = '';
		
		try {
			const now = new Date();
			const project: Project = {
				id: `project_${Date.now()}`,
				name: projectName.trim(),
				description: projectDescription.trim() || undefined,
				type: projectType,
				status: 'planning',
				priority: projectPriority,
				progress: 0,
				startDate: new Date(startDate),
				dueDate: dueDate ? new Date(dueDate) : undefined,
				createdAt: now,
				updatedAt: now,
				ownerId: 'current_user', // TODO: Get from auth store
				members: [
					{
						id: 'current_user',
						name: 'You',
						email: 'you@example.com',
						role: 'owner',
						status: 'online',
						lastSeen: now
					}
				],
				tags: projectTags,
				milestones: selectedTemplate?.defaultMilestones.map((m, index) => ({
					...m,
					id: `milestone_${Date.now()}_${index}`,
					projectId: `project_${Date.now()}`
				})) || [],
				tasks: selectedTemplate?.defaultTasks.map((t, index) => ({
					...t,
					id: `task_${Date.now()}_${index}`,
					projectId: `project_${Date.now()}`,
					createdAt: now,
					updatedAt: now
				})) || [],
				settings: {
					isPublic,
					allowGuestAccess,
					autoSync,
					notifications
				},
				stats: {
					totalTasks: selectedTemplate?.defaultTasks.length || 0,
					completedTasks: 0,
					overdueTasks: 0,
					activeMembers: 1,
					totalApiCalls: 0,
					lastActivity: now
				}
			};
			
			projectStore.addProject(project);
			
			// Reset form
			resetForm();
			dispatch('created', project);
			show = false;
			
		} catch (err: any) {
			error = err.message || 'Failed to create project';
		} finally {
			isCreating = false;
		}
	}
	
	function resetForm() {
		currentStep = 1;
		projectName = '';
		projectDescription = '';
		projectType = 'api';
		projectPriority = 'medium';
		selectedTemplate = null;
		startDate = new Date().toISOString().split('T')[0];
		dueDate = '';
		teamMembers = [];
		newMemberEmail = '';
		projectTags = [];
		newTag = '';
		autoSync = true;
		isPublic = false;
		allowGuestAccess = false;
		notifications = true;
		error = '';
	}
	
	function closeModal() {
		show = false;
		resetForm();
	}
</script>

{#if show}
	<div class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center p-4 z-50">
		<div class="bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-hidden">
			<!-- Header -->
			<div class="px-6 py-4 border-b border-gray-700">
				<div class="flex items-center justify-between">
					<h3 class="text-xl font-semibold text-white">Create New Project</h3>
					<IconButton variant="ghost" size="sm" on:click={closeModal}>
						<X class="w-3 h-3" />
					</IconButton>
				</div>
				
				<!-- Progress Steps -->
				<div class="flex items-center mt-4 space-x-4">
					{#each Array(4) as _, i}
						<div class="flex items-center">
							<div class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium
								{currentStep > i + 1 ? 'bg-blue-600 text-white' : 
								 currentStep === i + 1 ? 'bg-blue-600 text-white' : 
								 'bg-gray-700 text-gray-400'}">
								{i + 1}
							</div>
							{#if i < 3}
								<div class="w-12 h-0.5 mx-2 {currentStep > i + 1 ? 'bg-blue-600' : 'bg-gray-700'}"></div>
							{/if}
						</div>
					{/each}
				</div>
			</div>
			
			<div class="px-6 py-4 max-h-[60vh] overflow-y-auto">
				{#if error}
					<div class="bg-red-900/50 border border-red-700 rounded-md p-3 mb-4">
						<p class="text-sm text-red-200">{error}</p>
					</div>
				{/if}
				
				<!-- Step 1: Basic Info -->
				{#if currentStep === 1}
					<div class="space-y-4">
						<h4 class="text-lg font-medium text-white mb-4">Project Details</h4>
						
						<div>
							<label for="projectName" class="block text-sm font-medium text-gray-300 mb-2">
								Project Name *
							</label>
							<input
								id="projectName"
								type="text"
								bind:value={projectName}
								placeholder="My Awesome Project"
								class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							/>
						</div>
						
						<div>
							<label for="projectDescription" class="block text-sm font-medium text-gray-300 mb-2">
								Description
							</label>
							<textarea
								id="projectDescription"
								bind:value={projectDescription}
								placeholder="Brief description of your project..."
								rows="3"
								class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							></textarea>
						</div>
						
						<div class="grid grid-cols-2 gap-4">
							<div>
								<label for="projectType" class="block text-sm font-medium text-gray-300 mb-2">
									Project Type
								</label>
								<select
									id="projectType"
									bind:value={projectType}
									class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
								>
									{#each projectTypes as type}
										<option value={type.value}>{type.label}</option>
									{/each}
								</select>
								{#if selectedTypeInfo}
									<p class="text-xs text-gray-400 mt-1">{selectedTypeInfo.description}</p>
								{/if}
							</div>
							
							<div>
								<label for="projectPriority" class="block text-sm font-medium text-gray-300 mb-2">
									Priority
								</label>
								<select
									id="projectPriority"
									bind:value={projectPriority}
									class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
								>
									{#each priorities as priority}
										<option value={priority.value}>{priority.label}</option>
									{/each}
								</select>
							</div>
						</div>
					</div>
				{/if}
				
				<!-- Step 2: Template Selection -->
				{#if currentStep === 2}
					<div class="space-y-4">
						<h4 class="text-lg font-medium text-white mb-4">Choose Template</h4>
						
						{#if availableTemplates.length > 0}
							<div class="grid gap-3">
								<Button
									variant="ghost"
									size="lg"
									class="p-4 border-2 rounded-lg text-left transition-colors w-full justify-start
										{selectedTemplate === null ? 'border-blue-600 bg-blue-600/10' : 'border-gray-600 hover:border-gray-500'}"
									on:click={() => selectedTemplate = null}
								>
									<div class="flex items-center justify-between">
										<div>
											<h5 class="font-medium text-white">Start from scratch</h5>
											<p class="text-sm text-gray-400">Create an empty project</p>
										</div>
										{#if selectedTemplate === null}
											<div class="w-4 h-4 bg-blue-600 rounded-full"></div>
										{/if}
									</div>
								</Button>
								
								{#each availableTemplates as template}
									<Button
										variant="ghost"
										size="lg"
										class="p-4 border-2 rounded-lg text-left transition-colors w-full justify-start
											{selectedTemplate?.id === template.id ? 'border-blue-600 bg-blue-600/10' : 'border-gray-600 hover:border-gray-500'}"
										on:click={() => selectTemplate(template)}
									>
										<div class="flex items-center justify-between">
											<div class="flex-1">
												<h5 class="font-medium text-white">{template.name}</h5>
												<p class="text-sm text-gray-400 mb-2">{template.description}</p>
												<div class="flex items-center space-x-4 text-xs text-gray-500">
													<span>{template.defaultTasks.length} tasks</span>
													<span>{template.defaultMilestones.length} milestones</span>
													<span>~{template.estimatedDuration} days</span>
												</div>
											</div>
											{#if selectedTemplate?.id === template.id}
												<div class="w-4 h-4 bg-blue-600 rounded-full"></div>
											{/if}
										</div>
									</Button>
								{/each}
							</div>
						{:else}
							<div class="text-center py-8">
								<p class="text-gray-400">No templates available for {selectedTypeInfo?.label} projects</p>
								<p class="text-sm text-gray-500 mt-2">You can start from scratch or choose a different project type</p>
							</div>
						{/if}
					</div>
				{/if}
				
				<!-- Step 3: Timeline & Team -->
				{#if currentStep === 3}
					<div class="space-y-6">
						<h4 class="text-lg font-medium text-white mb-4">Timeline & Team</h4>
						
						<!-- Timeline -->
						<div class="grid grid-cols-2 gap-4">
							<div>
								<label for="startDate" class="block text-sm font-medium text-gray-300 mb-2">
									Start Date
								</label>
								<input
									id="startDate"
									type="date"
									bind:value={startDate}
									class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
								/>
							</div>
							
							<div>
								<label for="dueDate" class="block text-sm font-medium text-gray-300 mb-2">
									Due Date (optional)
								</label>
								<input
									id="dueDate"
									type="date"
									bind:value={dueDate}
									class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
								/>
							</div>
						</div>
						
						<!-- Team Members -->
						<div>
							<label class="block text-sm font-medium text-gray-300 mb-2">
								Team Members
							</label>
							<div class="flex space-x-2 mb-3">
								<input
									type="email"
									bind:value={newMemberEmail}
									placeholder="colleague@example.com"
									class="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
									on:keydown={(e) => e.key === 'Enter' && addTeamMember()}
								/>
								<Button variant="primary" size="sm" on:click={addTeamMember}>
									Add
								</Button>
							</div>
							{#if teamMembers.length > 0}
								<div class="flex flex-wrap gap-2">
									{#each teamMembers as member}
										<div class="flex items-center space-x-2 bg-gray-700 px-3 py-1 rounded-full">
											<span class="text-sm text-white">{member}</span>
											<IconButton variant="ghost" size="xs" on:click={() => removeTeamMember(member)}>
												<X class="btn-icon-xs" />
											</IconButton>
										</div>
									{/each}
								</div>
							{/if}
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
							{#if projectTags.length > 0}
								<div class="flex flex-wrap gap-2">
									{#each projectTags as tag}
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
				{/if}
				
				<!-- Step 4: Settings -->
				{#if currentStep === 4}
					<div class="space-y-6">
						<h4 class="text-lg font-medium text-white mb-4">Project Settings</h4>
						
						<div class="space-y-4">
							<div class="flex items-center justify-between">
								<div>
									<h5 class="text-sm font-medium text-white">Auto-sync APIs</h5>
									<p class="text-xs text-gray-400">Automatically sync API endpoints from your codebase</p>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input type="checkbox" bind:checked={autoSync} class="sr-only peer" />
									<div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
								</label>
							</div>
							
							<div class="flex items-center justify-between">
								<div>
									<h5 class="text-sm font-medium text-white">Public project</h5>
									<p class="text-xs text-gray-400">Make this project visible to everyone</p>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input type="checkbox" bind:checked={isPublic} class="sr-only peer" />
									<div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
								</label>
							</div>
							
							<div class="flex items-center justify-between">
								<div>
									<h5 class="text-sm font-medium text-white">Guest access</h5>
									<p class="text-xs text-gray-400">Allow non-members to view the project</p>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input type="checkbox" bind:checked={allowGuestAccess} disabled={!isPublic} class="sr-only peer" />
									<div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600 {!isPublic ? 'opacity-50' : ''}"></div>
								</label>
							</div>
							
							<div class="flex items-center justify-between">
								<div>
									<h5 class="text-sm font-medium text-white">Notifications</h5>
									<p class="text-xs text-gray-400">Receive updates about project activity</p>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input type="checkbox" bind:checked={notifications} class="sr-only peer" />
									<div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
								</label>
							</div>
						</div>
					</div>
				{/if}
			</div>
			
			<!-- Footer -->
			<div class="px-6 py-4 bg-gray-750 border-t border-gray-700 flex justify-between">
				<div>
					{#if currentStep > 1}
						<Button variant="ghost" size="sm" on:click={prevStep}>
							← Previous
						</Button>
					{/if}
				</div>
				
				<div class="flex space-x-3">
					<Button variant="ghost" size="sm" disabled={isCreating} on:click={closeModal}>
						Cancel
					</Button>
					
					{#if currentStep < 4}
						<Button variant="primary" size="sm" disabled={currentStep === 1 && !projectName.trim()} on:click={nextStep}>
							Next →
						</Button>
					{:else}
						<Button variant="primary" size="sm" disabled={isCreating || !projectName.trim()} loading={isCreating} on:click={createProject}>
							{#if isCreating}
								Creating...
							{:else}
								Create Project
							{/if}
						</Button>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.bg-gray-750 {
		background-color: #374151;
	}
</style>