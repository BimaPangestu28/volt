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

<div class="min-h-screen flex flex-col">
	<!-- Header -->
	<header class="electric-gradient text-white">
		<nav class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex items-center justify-between h-16">
				<div class="flex items-center">
					<h1 class="text-2xl font-bold">⚡ Volt</h1>
				</div>
				<div class="flex space-x-4">
					<a href="/login" class="bg-white/10 hover:bg-white/20 px-4 py-2 rounded-lg transition-colors">
						Login
					</a>
					<a href="/register" class="bg-yellow-400 hover:bg-yellow-300 text-gray-900 px-4 py-2 rounded-lg font-medium transition-colors">
						Get Started
					</a>
				</div>
			</div>
		</nav>
	</header>

	<!-- Hero Section -->
	<main class="flex-1 flex items-center justify-center px-4">
		<div class="max-w-4xl mx-auto text-center">
			<h2 class="text-5xl font-bold text-gray-900 mb-6">
				API Testing Made <span class="text-blue-600">Electric</span>
			</h2>
			<p class="text-xl text-gray-600 mb-8 leading-relaxed">
				Collaborative API client with real-time features, auto-import from your codebase, 
				and team synchronization. Built for modern development teams.
			</p>
			
			<div class="flex flex-col sm:flex-row gap-4 justify-center mb-12">
				<a href="/register" class="electric-gradient text-white px-8 py-4 rounded-lg font-medium text-lg hover:opacity-90 transition-opacity">
					Start Free Trial
				</a>
				<a href="/demo" class="border-2 border-blue-600 text-blue-600 px-8 py-4 rounded-lg font-medium text-lg hover:bg-blue-50 transition-colors">
					View Demo
				</a>
			</div>

			<!-- Features Grid -->
			<div class="grid md:grid-cols-3 gap-8 mt-16">
				<div class="card">
					<div class="w-12 h-12 electric-gradient rounded-lg flex items-center justify-center mb-4 mx-auto">
						<span class="text-white text-xl">🚀</span>
					</div>
					<h3 class="text-xl font-semibold mb-2">Auto-Import APIs</h3>
					<p class="text-gray-600">
						Automatically discover and sync API endpoints from NestJS, Express, and other frameworks.
					</p>
				</div>
				
				<div class="card">
					<div class="w-12 h-12 electric-gradient rounded-lg flex items-center justify-center mb-4 mx-auto">
						<span class="text-white text-xl">👥</span>
					</div>
					<h3 class="text-xl font-semibold mb-2">Real-time Collaboration</h3>
					<p class="text-gray-600">
						Work together with your team in real-time. See live cursors, share executions, and collaborate seamlessly.
					</p>
				</div>
				
				<div class="card">
					<div class="w-12 h-12 electric-gradient rounded-lg flex items-center justify-center mb-4 mx-auto">
						<span class="text-white text-xl">⚡</span>
					</div>
					<h3 class="text-xl font-semibold mb-2">Lightning Fast</h3>
					<p class="text-gray-600">
						Built for speed with modern web technologies. Test APIs instantly and get immediate feedback.
					</p>
				</div>
			</div>
		</div>
	</main>

	<!-- Footer -->
	<footer class="bg-gray-900 text-white py-8">
		<div class="max-w-7xl mx-auto px-4 text-center">
			<p>&copy; 2024 Volt API Client. Built with ⚡ for developers.</p>
		</div>
	</footer>
</div>

<style>
	.electric-gradient {
		background: linear-gradient(135deg, var(--electric-blue), var(--deep-blue));
	}
</style>
