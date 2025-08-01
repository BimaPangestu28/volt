<script lang="ts">
	import type { Project } from '$lib/stores/project';
	import { goto } from '$app/navigation';
	import { Server, Globe, Smartphone, Monitor, Database, Zap, Users } from 'lucide-svelte';
	
	export let project: Project;
	
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
			case 'urgent': return 'text-red-600 bg-red-50';
			case 'high': return 'text-orange-600 bg-orange-50';
			case 'medium': return 'text-amber-600 bg-amber-50';
			case 'low': return 'text-emerald-600 bg-emerald-50';
			default: return 'text-slate-600 bg-slate-50';
		}
	}
	
	function getTypeIcon(type: Project['type']) {
		switch (type) {
			case 'api': return Server;
			case 'web': return Globe;
			case 'mobile': return Smartphone;
			case 'desktop': return Monitor;
			case 'data': return Database;
			default: return Server;
		}
	}
	
	function formatDate(date: Date) {
		return new Intl.DateTimeFormat('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		}).format(new Date(date));
	}
	
	function getDaysRemaining(dueDate: Date | undefined) {
		if (!dueDate) return null;
		const today = new Date();
		const due = new Date(dueDate);
		const diffTime = due.getTime() - today.getTime();
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
		return diffDays;
	}
	
	function selectProject() {
		goto(`/project/${project.id}`);
	}
	
	$: onlineMembers = project.members.filter(m => m.status === 'online');
	$: daysRemaining = getDaysRemaining(project.dueDate);
	$: isOverdue = daysRemaining !== null && daysRemaining < 0;
	$: isDueSoon = daysRemaining !== null && daysRemaining <= 7 && daysRemaining >= 0;
</script>

<div 
	class="group bg-gradient-to-br from-white to-slate-50 rounded-xl border border-slate-200 shadow-sm hover:shadow-lg hover:border-slate-300 transition-all duration-300 cursor-pointer outline-none focus:outline-none"
	on:click={selectProject}
	on:keydown={(e) => e.key === 'Enter' && selectProject()}
	role="button"
	tabindex="0"
>
	<!-- Header -->
	<div class="p-5">
		<div class="flex items-start justify-between mb-4">
			<div class="flex items-center space-x-3">
				<div class="p-2 bg-blue-600 rounded-lg">
					<svelte:component this={getTypeIcon(project.type)} class="w-4 h-4 text-white" />
				</div>
				<div>
					<h3 class="text-sm font-semibold text-slate-900 group-hover:text-blue-600 transition-colors">
						{project.name}
					</h3>
					<div class="flex items-center space-x-2 mt-1">
						<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium {getPriorityColor(project.priority)}">
							{project.priority}
						</span>
						<div class="flex items-center space-x-1">
							<div class="w-2 h-2 rounded-full {getStatusColor(project.status)}"></div>
							<span class="text-xs text-slate-500 capitalize">{project.status.replace('_', ' ')}</span>
						</div>
					</div>
				</div>
			</div>
			
			<div class="flex items-center space-x-2">
				{#if project.settings.autoSync}
					<div class="p-1 bg-emerald-50 rounded" title="Auto-sync enabled">
						<Zap class="w-3 h-3 text-emerald-600" />
					</div>
				{/if}
				{#if project.settings.isPublic}
					<div class="p-1 bg-blue-50 rounded" title="Public project">
						<Users class="w-3 h-3 text-blue-600" />
					</div>
				{/if}
			</div>
		</div>
		
		{#if project.description}
			<p class="text-sm text-slate-600 mb-4 line-clamp-2">{project.description}</p>
		{/if}
		
		<!-- Progress Bar -->
		<div class="mb-4">
			<div class="flex items-center justify-between mb-2">
				<span class="text-xs text-slate-500">Progress</span>
				<span class="text-xs font-medium text-slate-700">{project.progress}%</span>
			</div>
			<div class="w-full bg-slate-200 rounded-full h-1.5">
				<div 
					class="bg-blue-600 h-1.5 rounded-full transition-all duration-300"
					style="width: {project.progress}%"
				></div>
			</div>
		</div>
		
		<!-- Stats Grid -->
		<div class="grid grid-cols-3 gap-3 mb-4">
			<div class="text-center bg-slate-50 rounded-lg py-2">
				<div class="text-sm font-semibold text-slate-900">{project.stats.totalTasks}</div>
				<div class="text-xs text-slate-500">Tasks</div>
			</div>
			<div class="text-center bg-emerald-50 rounded-lg py-2">
				<div class="text-sm font-semibold text-emerald-700">{project.stats.completedTasks}</div>
				<div class="text-xs text-slate-500">Done</div>
			</div>
			<div class="text-center bg-red-50 rounded-lg py-2">
				{#if project.stats.overdueTasks > 0}
					<div class="text-sm font-semibold text-red-600">{project.stats.overdueTasks}</div>
				{:else}
					<div class="text-sm font-semibold text-slate-400">0</div>
				{/if}
				<div class="text-xs text-slate-500">Overdue</div>
			</div>
		</div>
	</div>
	
	<!-- Footer -->
	<div class="px-5 py-3 bg-slate-50 rounded-b-xl border-t border-slate-200">
		<div class="flex items-center justify-between">
			<!-- Team Members -->
			<div class="flex items-center space-x-2">
				<div class="flex -space-x-1.5">
					{#each project.members.slice(0, 3) as member}
						<div 
							class="w-5 h-5 rounded-full bg-blue-600 border-2 border-white flex items-center justify-center relative"
							title={member.name}
						>
							{#if member.avatar}
								<img src={member.avatar} alt={member.name} class="w-full h-full rounded-full object-cover" />
							{:else}
								<span class="text-xs font-medium text-white">{member.name.charAt(0)}</span>
							{/if}
							{#if member.status === 'online'}
								<div class="absolute -bottom-0.5 -right-0.5 w-1.5 h-1.5 bg-emerald-500 border border-white rounded-full"></div>
							{/if}
						</div>
					{/each}
					{#if project.members.length > 3}
						<div class="w-5 h-5 rounded-full bg-slate-300 border-2 border-white flex items-center justify-center">
							<span class="text-xs font-medium text-slate-600">+{project.members.length - 3}</span>
						</div>
					{/if}
				</div>
				{#if onlineMembers.length > 0}
					<span class="text-xs text-emerald-600">{onlineMembers.length} online</span>
				{/if}
			</div>
			
			<!-- Due Date -->
			{#if project.dueDate}
				<div class="text-right">
					{#if isOverdue}
						<div class="text-xs text-red-600 font-medium">
							Overdue by {Math.abs(daysRemaining)} days
						</div>
					{:else if isDueSoon}
						<div class="text-xs text-amber-600 font-medium">
							Due in {daysRemaining} days
						</div>
					{:else}
						<div class="text-xs text-slate-500">
							Due {formatDate(project.dueDate)}
						</div>
					{/if}
				</div>
			{:else}
				<div class="text-xs text-slate-400">No due date</div>
			{/if}
		</div>
		
		<!-- Tags -->
		{#if project.tags.length > 0}
			<div class="flex flex-wrap gap-1 mt-3">
				{#each project.tags.slice(0, 3) as tag}
					<span class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-slate-100 text-slate-600">
						{tag}
					</span>
				{/each}
				{#if project.tags.length > 3}
					<span class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-slate-100 text-slate-500">
						+{project.tags.length - 3}
					</span>
				{/if}
			</div>
		{/if}
	</div>
</div>

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>