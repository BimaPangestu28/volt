<script lang="ts">
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	let isLoggedIn = false;

	authStore.subscribe(state => {
		isLoggedIn = state.isAuthenticated;
	});

	onMount(() => {
		if (isLoggedIn) {
			goto('/dashboard');
		}
	});
</script>

<svelte:head>
	<title>Volt - API Client for Teams</title>
</svelte:head>

<div class="min-h-screen bg-white">
	<!-- Header -->
	<nav class="border-b border-gray-200">
		<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between items-center h-16">
				<div class="flex items-center">
					<div class="flex-shrink-0 flex items-center">
						<span class="text-xl font-semibold text-gray-900">Volt</span>
					</div>
				</div>
				<div class="flex items-center space-x-4">
					<a href="/docs" class="text-sm text-gray-500 hover:text-gray-900">Docs</a>
					<a href="/login" class="text-sm text-gray-500 hover:text-gray-900">Login</a>
					<a href="/register" class="bg-gray-900 text-white px-3 py-1.5 rounded-md text-sm font-medium hover:bg-gray-800">
						Sign up
					</a>
				</div>
			</div>
		</div>
	</nav>

	<!-- Hero Section -->
	<main class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="py-20">
			<div class="text-center">
				<h1 class="text-5xl font-bold text-gray-900 tracking-tight mb-6">
					Collaborative API client with auto-import
				</h1>
				<p class="text-xl text-gray-600 max-w-3xl mx-auto mb-10">
					Real-time collaboration and automatic API discovery from your NestJS, Express, and FastAPI projects. 
					Stop manually importing endpoints.
				</p>
				
				<div class="bg-gray-50 rounded-lg p-6 max-w-2xl mx-auto mb-12">
					<div class="text-sm text-gray-500 mb-2">Install the CLI</div>
					<div class="bg-black text-green-400 font-mono text-sm p-3 rounded mb-3">
						<span class="text-gray-400">$</span> npm install -g @volt/cli
					</div>
					<div class="text-sm text-gray-500 mb-2">Auto-import your NestJS APIs</div>
					<div class="bg-black text-green-400 font-mono text-sm p-3 rounded">
						<span class="text-gray-400">$</span> volt init --framework nestjs
					</div>
				</div>
				
				<div class="flex justify-center gap-4">
					<a href="/register" class="bg-gray-900 text-white px-6 py-3 rounded-md font-medium hover:bg-gray-800">
						Start free trial
					</a>
					<a href="/demo" class="border border-gray-300 px-6 py-3 rounded-md hover:border-gray-400">
						Watch demo
					</a>
				</div>
			</div>
		</div>

		<!-- Features Grid -->
		<div class="py-16">
			<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
				<div class="border border-gray-200 rounded-lg p-6">
					<div class="text-2xl mb-3">🔄</div>
					<h3 class="text-lg font-semibold text-gray-900 mb-3">Automatic API Discovery</h3>
					<p class="text-gray-600 text-sm leading-relaxed mb-3">
						Install <code class="bg-gray-100 px-1 rounded">@volt/nestjs</code> in your project and watch all your endpoints automatically appear in Volt.
					</p>
					<div class="text-xs text-gray-500">
						Supports: NestJS • Express (coming soon) • FastAPI (coming soon)
					</div>
				</div>
				
				<div class="border border-gray-200 rounded-lg p-6">
					<div class="text-2xl mb-3">👥</div>
					<h3 class="text-lg font-semibold text-gray-900 mb-3">Live Collaboration</h3>
					<p class="text-gray-600 text-sm leading-relaxed mb-3">
						See your team's changes in real-time. Live cursors, shared request execution, and instant sync across all devices.
					</p>
					<div class="text-xs text-gray-500">
						Multi-user editing • Live presence • Real-time comments
					</div>
				</div>
				
				<div class="border border-gray-200 rounded-lg p-6">
					<div class="text-2xl mb-3">⚡</div>
					<h3 class="text-lg font-semibold text-gray-900 mb-3">Built for Speed</h3>
					<p class="text-gray-600 text-sm leading-relaxed mb-3">
						Modern interface with keyboard shortcuts, instant request execution, and zero-config environment sync.
					</p>
					<div class="text-xs text-gray-500">
						WebSocket-powered • Environment variables • Request history
					</div>
				</div>
			</div>
		</div>
		<!-- Getting Started Section -->
		<div class="border-t border-gray-200 py-16">
			<div class="text-center mb-12">
				<h2 class="text-2xl font-semibold text-gray-900 mb-4">How it works</h2>
				<p class="text-gray-600">Three steps to collaborative API testing</p>
			</div>
			
			<div class="grid md:grid-cols-3 gap-8">
				<div class="text-center">
					<div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
						<span class="text-blue-600 font-semibold">1</span>
					</div>
					<h3 class="font-semibold mb-2">Install & Connect</h3>
					<p class="text-gray-600 text-sm">
						Add <code class="bg-gray-100 px-1 rounded">@volt/nestjs</code> to your project and run <code class="bg-gray-100 px-1 rounded">volt sync</code>
					</p>
				</div>
				
				<div class="text-center">
					<div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
						<span class="text-blue-600 font-semibold">2</span>
					</div>
					<h3 class="font-semibold mb-2">Auto-Import</h3>
					<p class="text-gray-600 text-sm">
						All your API endpoints appear automatically with types, docs, and sample data
					</p>
				</div>
				
				<div class="text-center">
					<div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
						<span class="text-blue-600 font-semibold">3</span>
					</div>
					<h3 class="font-semibold mb-2">Collaborate</h3>
					<p class="text-gray-600 text-sm">
						Invite your team and test APIs together in real-time with live sync
					</p>
				</div>
			</div>
		</div>

	</main>
	
	<!-- Footer -->
	<footer class="border-t border-gray-200 py-8">
		<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center text-sm text-gray-500">
				<p>&copy; 2024 Volt. Made by developers for developers.</p>
			</div>
		</div>
	</footer>
</div>

