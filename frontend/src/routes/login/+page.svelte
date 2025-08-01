<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { apiClient } from '$lib/api/client';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import { Eye, EyeOff } from 'lucide-svelte';

	let email = '';
	let password = '';
	let isLoading = false;
	let error = '';
	let showPassword = false;

	let isAuthenticated = false;

	authStore.subscribe(state => {
		isAuthenticated = state.isAuthenticated;
	});

	onMount(() => {
		if (isAuthenticated) {
			goto('/dashboard');
		}
	});

	async function handleLogin() {
		if (!email || !password) {
			error = 'Please fill in all fields';
			return;
		}

		isLoading = true;
		error = '';

		try {
			const response = await apiClient.login(email, password);
			authStore.setUser(response.user);
			goto('/dashboard');
		} catch (err: any) {
			error = err.response?.data?.error || 'Login failed. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	function handleKeyPress(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			handleLogin();
		}
	}
</script>

<svelte:head>
	<title>Login - Volt</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<!-- Logo -->
		<div class="flex justify-center">
			<h1 class="text-4xl font-bold electric-gradient bg-clip-text text-transparent">⚡ Volt</h1>
		</div>
		<h2 class="mt-6 text-center text-3xl font-bold text-gray-900">
			Sign in to your account
		</h2>
		<p class="mt-2 text-center text-sm text-gray-600">
			Or 
			<a href="/register" class="font-medium text-blue-600 hover:text-blue-500 transition-colors">
				create a new account
			</a>
		</p>
	</div>

	<div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
		<div class="bg-white py-8 px-4 shadow-lg sm:rounded-lg sm:px-10 border border-gray-200">
			<form class="space-y-6" on:submit|preventDefault={handleLogin}>
				<!-- Email Field -->
				<div>
					<label for="email" class="block text-sm font-medium text-gray-700">
						Email address
					</label>
					<div class="mt-1">
						<input
							id="email"
							name="email"
							type="email"
							autocomplete="email"
							required
							bind:value={email}
							on:keypress={handleKeyPress}
							class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
							placeholder="Enter your email"
						/>
					</div>
				</div>

				<!-- Password Field -->
				<div>
					<label for="password" class="block text-sm font-medium text-gray-700">
						Password
					</label>
					<div class="mt-1 relative">
						{#if showPassword}
							<input
								id="password"
								name="password"
								type="text"
								autocomplete="current-password"
								required
								bind:value={password}
								on:keypress={handleKeyPress}
								class="appearance-none block w-full px-3 py-2 pr-10 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Enter your password"
							/>
						{:else}
							<input
								id="password"
								name="password"
								type="password"
								autocomplete="current-password"
								required
								bind:value={password}
								on:keypress={handleKeyPress}
								class="appearance-none block w-full px-3 py-2 pr-10 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Enter your password"
							/>
						{/if}
						<IconButton
							variant="ghost"
							size="sm"
							on:click={() => showPassword = !showPassword}
							class="absolute inset-y-0 right-0 pr-3 text-gray-400 hover:text-gray-600"
						>
							{#if showPassword}
								<EyeOff class="w-5 h-5" />
							{:else}
								<Eye class="w-5 h-5" />
							{/if}
						</IconButton>
					</div>
				</div>

				<!-- Error Message -->
				{#if error}
					<div class="bg-red-50 border border-red-200 rounded-md p-3">
						<p class="text-sm text-red-800">{error}</p>
					</div>
				{/if}

				<!-- Submit Button -->
				<div>
					<Button
						type="submit"
						variant="primary"
						size="md"
						disabled={isLoading}
						loading={isLoading}
						class="w-full electric-gradient"
					>
						{isLoading ? 'Signing in...' : 'Sign in'}
					</Button>
				</div>
			</form>

			<!-- Divider -->
			<div class="mt-6">
				<div class="relative">
					<div class="absolute inset-0 flex items-center">
						<div class="w-full border-t border-gray-300" />
					</div>
					<div class="relative flex justify-center text-sm">
						<span class="px-2 bg-white text-gray-500">New to Volt?</span>
					</div>
				</div>

				<div class="mt-6">
					<a
						href="/register"
						class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 transition-colors"
					>
						Create an account
					</a>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.electric-gradient {
		background: linear-gradient(135deg, var(--electric-blue), var(--deep-blue));
	}
	
	.bg-clip-text {
		background-clip: text;
		-webkit-background-clip: text;
	}
</style>