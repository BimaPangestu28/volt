<script lang="ts">
	import { onMount } from 'svelte';
	import type { Project } from '$lib/stores/project';
	import Button from '$lib/components/ui/Button.svelte';
	
	export let project: Project;
	
	interface AnalyticsData {
		taskCompletion: {
			completed: number;
			inProgress: number;
			todo: number;
			review: number;
		};
		teamActivity: {
			name: string;
			tasksCompleted: number;
			hoursLogged: number;
			lastActive: Date;
		}[];
		timeTracking: {
			date: string;
			hours: number;
		}[];
		milestoneProgress: {
			name: string;
			progress: number;
			dueDate: Date;
			status: 'on_track' | 'at_risk' | 'overdue';
		}[];
		velocityMetrics: {
			week: string;
			tasksCompleted: number;
			averageTime: number; // hours
		}[];
		priorityDistribution: {
			urgent: number;
			high: number;
			medium: number;
			low: number;
		};
	}
	
	let analytics: AnalyticsData | null = null;
	let isLoading = true;
	let selectedTimeRange: '7d' | '30d' | '90d' | '1y' = '30d';
	let selectedMetric: 'overview' | 'team' | 'velocity' | 'milestones' = 'overview';
	
	// Sample analytics data
	const sampleAnalytics: AnalyticsData = {
		taskCompletion: {
			completed: 28,
			inProgress: 12,
			todo: 18,
			review: 6
		},
		teamActivity: [
			{
				name: 'Sarah Chen',
				tasksCompleted: 15,
				hoursLogged: 84.5,
				lastActive: new Date(Date.now() - 2 * 60 * 60 * 1000)
			},
			{
				name: 'Mike Johnson',
				tasksCompleted: 12,
				hoursLogged: 72.3,
				lastActive: new Date(Date.now() - 4 * 60 * 60 * 1000)
			},
			{
				name: 'Alex Kumar',
				tasksCompleted: 8,
				hoursLogged: 45.8,
				lastActive: new Date(Date.now() - 8 * 60 * 60 * 1000)
			},
			{
				name: 'You',
				tasksCompleted: 18,
				hoursLogged: 92.4,
				lastActive: new Date(Date.now() - 30 * 60 * 1000)
			}
		],
		timeTracking: [
			{ date: '2024-01-15', hours: 8.5 },
			{ date: '2024-01-16', hours: 7.2 },
			{ date: '2024-01-17', hours: 9.1 },
			{ date: '2024-01-18', hours: 6.8 },
			{ date: '2024-01-19', hours: 8.0 },
			{ date: '2024-01-22', hours: 9.5 },
			{ date: '2024-01-23', hours: 7.8 },
			{ date: '2024-01-24', hours: 8.2 },
			{ date: '2024-01-25', hours: 7.5 },
			{ date: '2024-01-26', hours: 8.8 }
		],
		milestoneProgress: [
			{
				name: 'MVP Release',
				progress: 78,
				dueDate: new Date(Date.now() + 14 * 24 * 60 * 60 * 1000),
				status: 'on_track'
			},
			{
				name: 'User Testing',
				progress: 45,
				dueDate: new Date(Date.now() + 21 * 24 * 60 * 60 * 1000),
				status: 'at_risk'
			},
			{
				name: 'Beta Launch',
				progress: 15,
				dueDate: new Date(Date.now() + 35 * 24 * 60 * 60 * 1000),
				status: 'on_track'
			}
		],
		velocityMetrics: [
			{ week: 'Week 1', tasksCompleted: 8, averageTime: 6.2 },
			{ week: 'Week 2', tasksCompleted: 12, averageTime: 5.8 },
			{ week: 'Week 3', tasksCompleted: 15, averageTime: 4.9 },
			{ week: 'Week 4', tasksCompleted: 11, averageTime: 5.5 }
		],
		priorityDistribution: {
			urgent: 8,
			high: 15,
			medium: 28,
			low: 13
		}
	};
	
	onMount(() => {
		// Simulate loading
		setTimeout(() => {
			analytics = sampleAnalytics;
			isLoading = false;
		}, 1000);
	});
	
	function formatDate(date: Date) {
		return new Intl.DateTimeFormat('en-US', {
			month: 'short',
			day: 'numeric'
		}).format(new Date(date));
	}
	
	function getStatusColor(status: 'on_track' | 'at_risk' | 'overdue') {
		switch (status) {
			case 'on_track':
				return 'text-green-400 bg-green-400/10';
			case 'at_risk':
				return 'text-yellow-400 bg-yellow-400/10';
			case 'overdue':
				return 'text-red-400 bg-red-400/10';
		}
	}
	
	function calculateTaskCompletionRate() {
		if (!analytics) return 0;
		const total = Object.values(analytics.taskCompletion).reduce((sum, count) => sum + count, 0);
		return total > 0 ? Math.round((analytics.taskCompletion.completed / total) * 100) : 0;
	}
	
	function calculateAverageVelocity() {
		if (!analytics) return 0;
		const total = analytics.velocityMetrics.reduce((sum, week) => sum + week.tasksCompleted, 0);
		return Math.round(total / analytics.velocityMetrics.length);
	}
	
	function getTotalHoursLogged() {
		if (!analytics) return 0;
		return analytics.teamActivity.reduce((sum, member) => sum + member.hoursLogged, 0);
	}
	
	$: completionRate = calculateTaskCompletionRate();
	$: averageVelocity = calculateAverageVelocity();
	$: totalHours = getTotalHoursLogged();
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<h3 class="text-lg font-semibold text-white">Project Analytics</h3>
		<div class="flex items-center space-x-2">
			<!-- Time Range Selector -->
			<select
				bind:value={selectedTimeRange}
				class="bg-gray-700 border border-gray-600 text-white text-sm rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
			>
				<option value="7d">Last 7 days</option>
				<option value="30d">Last 30 days</option>
				<option value="90d">Last 90 days</option>
				<option value="1y">Last year</option>
			</select>
		</div>
	</div>
	
	{#if isLoading}
		<!-- Loading State -->
		<div class="flex items-center justify-center py-12">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
			<span class="ml-3 text-gray-400">Loading analytics...</span>
		</div>
	{:else if analytics}
		<!-- Metric Navigation -->
		<div class="bg-gray-800 rounded-lg border border-gray-700 p-1">
			<nav class="flex space-x-1">
				<Button
					variant={selectedMetric === 'overview' ? 'primary' : 'ghost'}
					size="sm"
					on:click={() => selectedMetric = 'overview'}
				>
					Overview
				</Button>
				<Button
					variant={selectedMetric === 'team' ? 'primary' : 'ghost'}
					size="sm"
					on:click={() => selectedMetric = 'team'}
				>
					Team Performance
				</Button>
				<Button
					variant={selectedMetric === 'velocity' ? 'primary' : 'ghost'}
					size="sm"
					on:click={() => selectedMetric = 'velocity'}
				>
					Velocity
				</Button>
				<Button
					variant={selectedMetric === 'milestones' ? 'primary' : 'ghost'}
					size="sm"
					on:click={() => selectedMetric = 'milestones'}
				>
					Milestones
				</Button>
			</nav>
		</div>
		
		{#if selectedMetric === 'overview'}
			<!-- Overview Metrics -->
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
				<!-- Completion Rate -->
				<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
					<div class="flex items-center justify-between">
						<div>
							<p class="text-sm text-gray-400">Completion Rate</p>
							<p class="text-3xl font-semibold text-green-400">{completionRate}%</p>
						</div>
						<div class="p-3 bg-green-400/10 rounded-lg">
							<svg class="w-6 h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
							</svg>
						</div>
					</div>
					<div class="mt-4 w-full bg-gray-700 rounded-full h-2">
						<div class="bg-green-400 h-2 rounded-full transition-all" style="width: {completionRate}%"></div>
					</div>
				</div>
				
				<!-- Average Velocity -->
				<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
					<div class="flex items-center justify-between">
						<div>
							<p class="text-sm text-gray-400">Avg Velocity</p>
							<p class="text-3xl font-semibold text-blue-400">{averageVelocity}</p>
							<p class="text-xs text-gray-500">tasks/week</p>
						</div>
						<div class="p-3 bg-blue-400/10 rounded-lg">
							<svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
							</svg>
						</div>
					</div>
				</div>
				
				<!-- Total Hours -->
				<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
					<div class="flex items-center justify-between">
						<div>
							<p class="text-sm text-gray-400">Total Hours</p>
							<p class="text-3xl font-semibold text-purple-400">{totalHours.toFixed(1)}</p>
							<p class="text-xs text-gray-500">logged</p>
						</div>
						<div class="p-3 bg-purple-400/10 rounded-lg">
							<svg class="w-6 h-6 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
							</svg>
						</div>
					</div>
				</div>
				
				<!-- Active Members -->
				<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
					<div class="flex items-center justify-between">
						<div>
							<p class="text-sm text-gray-400">Active Members</p>
							<p class="text-3xl font-semibold text-yellow-400">{project.members.filter(m => m.status === 'online').length}</p>
							<p class="text-xs text-gray-500">of {project.members.length}</p>
						</div>
						<div class="p-3 bg-yellow-400/10 rounded-lg">
							<svg class="w-6 h-6 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
							</svg>
						</div>
					</div>
				</div>
			</div>
			
			<!-- Task Completion Breakdown -->
			<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
				<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
					<h4 class="text-lg font-semibold text-white mb-4">Task Status Distribution</h4>
					<div class="space-y-4">
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-3">
								<div class="w-3 h-3 bg-green-500 rounded-full"></div>
								<span class="text-sm text-gray-300">Completed</span>
							</div>
							<div class="flex items-center space-x-3">
								<span class="text-sm text-white font-medium">{analytics.taskCompletion.completed}</span>
								<div class="w-20 bg-gray-700 rounded-full h-2">
									<div class="bg-green-500 h-2 rounded-full" style="width: {(analytics.taskCompletion.completed / 64) * 100}%"></div>
								</div>
							</div>
						</div>
						
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-3">
								<div class="w-3 h-3 bg-blue-500 rounded-full"></div>
								<span class="text-sm text-gray-300">In Progress</span>
							</div>
							<div class="flex items-center space-x-3">
								<span class="text-sm text-white font-medium">{analytics.taskCompletion.inProgress}</span>
								<div class="w-20 bg-gray-700 rounded-full h-2">
									<div class="bg-blue-500 h-2 rounded-full" style="width: {(analytics.taskCompletion.inProgress / 64) * 100}%"></div>
								</div>
							</div>
						</div>
						
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-3">
								<div class="w-3 h-3 bg-yellow-500 rounded-full"></div>
								<span class="text-sm text-gray-300">Review</span>
							</div>
							<div class="flex items-center space-x-3">
								<span class="text-sm text-white font-medium">{analytics.taskCompletion.review}</span>
								<div class="w-20 bg-gray-700 rounded-full h-2">
									<div class="bg-yellow-500 h-2 rounded-full" style="width: {(analytics.taskCompletion.review / 64) * 100}%"></div>
								</div>
							</div>
						</div>
						
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-3">
								<div class="w-3 h-3 bg-gray-500 rounded-full"></div>
								<span class="text-sm text-gray-300">To Do</span>
							</div>
							<div class="flex items-center space-x-3">
								<span class="text-sm text-white font-medium">{analytics.taskCompletion.todo}</span>
								<div class="w-20 bg-gray-700 rounded-full h-2">
									<div class="bg-gray-500 h-2 rounded-full" style="width: {(analytics.taskCompletion.todo / 64) * 100}%"></div>
								</div>
							</div>
						</div>
					</div>
				</div>
				
				<!-- Priority Distribution -->
				<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
					<h4 class="text-lg font-semibold text-white mb-4">Priority Distribution</h4>
					<div class="space-y-4">
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-3">
								<div class="w-3 h-3 bg-red-500 rounded-full"></div>
								<span class="text-sm text-gray-300">Urgent</span>
							</div>
							<div class="flex items-center space-x-3">
								<span class="text-sm text-white font-medium">{analytics.priorityDistribution.urgent}</span>
								<div class="w-20 bg-gray-700 rounded-full h-2">
									<div class="bg-red-500 h-2 rounded-full" style="width: {(analytics.priorityDistribution.urgent / 64) * 100}%"></div>
								</div>
							</div>
						</div>
						
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-3">
								<div class="w-3 h-3 bg-orange-500 rounded-full"></div>
								<span class="text-sm text-gray-300">High</span>
							</div>
							<div class="flex items-center space-x-3">
								<span class="text-sm text-white font-medium">{analytics.priorityDistribution.high}</span>
								<div class="w-20 bg-gray-700 rounded-full h-2">
									<div class="bg-orange-500 h-2 rounded-full" style="width: {(analytics.priorityDistribution.high / 64) * 100}%"></div>
								</div>
							</div>
						</div>
						
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-3">
								<div class="w-3 h-3 bg-yellow-500 rounded-full"></div>
								<span class="text-sm text-gray-300">Medium</span>
							</div>
							<div class="flex items-center space-x-3">
								<span class="text-sm text-white font-medium">{analytics.priorityDistribution.medium}</span>
								<div class="w-20 bg-gray-700 rounded-full h-2">
									<div class="bg-yellow-500 h-2 rounded-full" style="width: {(analytics.priorityDistribution.medium / 64) * 100}%"></div>
								</div>
							</div>
						</div>
						
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-3">
								<div class="w-3 h-3 bg-green-500 rounded-full"></div>
								<span class="text-sm text-gray-300">Low</span>
							</div>
							<div class="flex items-center space-x-3">
								<span class="text-sm text-white font-medium">{analytics.priorityDistribution.low}</span>
								<div class="w-20 bg-gray-700 rounded-full h-2">
									<div class="bg-green-500 h-2 rounded-full" style="width: {(analytics.priorityDistribution.low / 64) * 100}%"></div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
			
		{:else if selectedMetric === 'team'}
			<!-- Team Performance -->
			<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
				<h4 class="text-lg font-semibold text-white mb-6">Team Performance</h4>
				<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
					{#each analytics.teamActivity as member}
						<div class="bg-gray-750 rounded-lg p-4">
							<div class="flex items-center space-x-3 mb-4">
								<div class="w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center">
									<span class="text-sm font-medium text-white">{member.name.charAt(0)}</span>
								</div>
								<div>
									<h5 class="text-sm font-medium text-white">{member.name}</h5>
									<p class="text-xs text-gray-400">Last active: {formatDate(member.lastActive)}</p>
								</div>
							</div>
							
							<div class="space-y-3">
								<div class="flex items-center justify-between">
									<span class="text-sm text-gray-400">Tasks Completed</span>
									<div class="flex items-center space-x-2">
										<span class="text-sm font-medium text-white">{member.tasksCompleted}</span>
										<div class="w-16 bg-gray-700 rounded-full h-2">
											<div class="bg-green-500 h-2 rounded-full" style="width: {(member.tasksCompleted / 20) * 100}%"></div>
										</div>
									</div>
								</div>
								
								<div class="flex items-center justify-between">
									<span class="text-sm text-gray-400">Hours Logged</span>
									<div class="flex items-center space-x-2">
										<span class="text-sm font-medium text-white">{member.hoursLogged.toFixed(1)}h</span>
										<div class="w-16 bg-gray-700 rounded-full h-2">
											<div class="bg-blue-500 h-2 rounded-full" style="width: {(member.hoursLogged / 100) * 100}%"></div>
										</div>
									</div>
								</div>
								
								<div class="flex items-center justify-between">
									<span class="text-sm text-gray-400">Avg per Task</span>
									<span class="text-sm font-medium text-white">{(member.hoursLogged / member.tasksCompleted).toFixed(1)}h</span>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
			
		{:else if selectedMetric === 'velocity'}
			<!-- Velocity Charts -->
			<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
				<h4 class="text-lg font-semibold text-white mb-6">Sprint Velocity</h4>
				<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
					<!-- Tasks Completed Chart -->
					<div>
						<h5 class="text-sm font-medium text-gray-300 mb-4">Tasks Completed per Week</h5>
						<div class="space-y-3">
							{#each analytics.velocityMetrics as week}
								<div class="flex items-center space-x-3">
									<div class="w-16 text-sm text-gray-400">{week.week}</div>
									<div class="flex-1 bg-gray-700 rounded-full h-6 relative">
										<div 
											class="bg-blue-500 h-6 rounded-full flex items-center justify-end pr-3 transition-all"
											style="width: {(week.tasksCompleted / 20) * 100}%"
										>
											<span class="text-xs font-medium text-white">{week.tasksCompleted}</span>
										</div>
									</div>
								</div>
							{/each}
						</div>
					</div>
					
					<!-- Average Time Chart -->
					<div>
						<h5 class="text-sm font-medium text-gray-300 mb-4">Average Time per Task</h5>
						<div class="space-y-3">
							{#each analytics.velocityMetrics as week}
								<div class="flex items-center space-x-3">
									<div class="w-16 text-sm text-gray-400">{week.week}</div>
									<div class="flex-1 bg-gray-700 rounded-full h-6 relative">
										<div 
											class="bg-purple-500 h-6 rounded-full flex items-center justify-end pr-3 transition-all"
											style="width: {(week.averageTime / 10) * 100}%"
										>
											<span class="text-xs font-medium text-white">{week.averageTime}h</span>
										</div>
									</div>
								</div>
							{/each}
						</div>
					</div>
				</div>
			</div>
			
		{:else if selectedMetric === 'milestones'}
			<!-- Milestone Progress -->
			<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
				<h4 class="text-lg font-semibold text-white mb-6">Milestone Progress</h4>
				<div class="space-y-6">
					{#each analytics.milestoneProgress as milestone}
						<div class="bg-gray-750 rounded-lg p-4">
							<div class="flex items-center justify-between mb-3">
								<div>
									<h5 class="text-sm font-medium text-white">{milestone.name}</h5>
									<p class="text-xs text-gray-400">Due: {formatDate(milestone.dueDate)}</p>
								</div>
								<div class="flex items-center space-x-2">
									<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium {getStatusColor(milestone.status)} capitalize">
										{milestone.status.replace('_', ' ')}
									</span>
									<span class="text-sm font-medium text-white">{milestone.progress}%</span>
								</div>
							</div>
							
							<div class="w-full bg-gray-700 rounded-full h-3">
								<div 
									class="h-3 rounded-full transition-all {
										milestone.status === 'on_track' ? 'bg-green-500' :
										milestone.status === 'at_risk' ? 'bg-yellow-500' :
										'bg-red-500'
									}"
									style="width: {milestone.progress}%"
								></div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
		
		<!-- Time Tracking Chart -->
		<div class="bg-gray-800 rounded-lg border border-gray-700 p-6">
			<h4 class="text-lg font-semibold text-white mb-6">Daily Time Tracking</h4>
			<div class="flex items-end space-x-2 h-32">
				{#each analytics.timeTracking as day}
					<div class="flex-1 flex flex-col items-center">
						<div 
							class="w-full bg-blue-500 rounded-t transition-all hover:bg-blue-400 cursor-pointer"
							style="height: {(day.hours / 10) * 100}%"
							title="{day.date}: {day.hours} hours"
						></div>
						<span class="text-xs text-gray-400 mt-2 transform rotate-45 origin-bottom-left">
							{new Date(day.date).getDate()}
						</span>
					</div>
				{/each}
			</div>
			<div class="flex items-center justify-between mt-4 text-xs text-gray-400">
				<span>Daily Hours</span>
				<span>Average: {(analytics.timeTracking.reduce((sum, day) => sum + day.hours, 0) / analytics.timeTracking.length).toFixed(1)}h</span>
			</div>
		</div>
	{/if}
</div>

<style>
	.bg-gray-750 {
		background-color: #374151;
	}
</style>