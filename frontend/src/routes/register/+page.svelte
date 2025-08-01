<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { apiClient } from '$lib/api/client';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import { Eye, EyeOff } from 'lucide-svelte';

	let username = '';
	let email = '';
	let password = '';
	let confirmPassword = '';
	let isLoading = false;
	let error = '';
	let showPassword = false;
	let showConfirmPassword = false;

	let isAuthenticated = false;

	authStore.subscribe(state => {
		isAuthenticated = state.isAuthenticated;
	});

	onMount(() => {
		if (isAuthenticated) {
			goto('/dashboard');
		}
	});

	async function handleRegister() {
		// Validation
		if (!username || !email || !password || !confirmPassword) {
			error = 'Please fill in all fields';
			return;
		}

		if (username.length < 3) {
			error = 'Username must be at least 3 characters long';
			return;
		}

		if (password.length < 8) {
			error = 'Password must be at least 8 characters long';
			return;
		}

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		isLoading = true;
		error = '';

		try {
			const response = await apiClient.register(username, email, password);
			authStore.setUser(response.user);
			goto('/dashboard');
		} catch (err: any) {
			error = err.response?.data?.error || 'Registration failed. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	function handleKeyPress(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			handleRegister();
		}
	}
</script>

<svelte:head>
	<title>Sign Up - Volt</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<!-- Logo -->
		<div class="flex justify-center">
			<h1 class="text-4xl font-bold electric-gradient bg-clip-text text-transparent">⚡ Volt</h1>
		</div>
		<h2 class="mt-6 text-center text-3xl font-bold text-gray-900">
			Create your account
		</h2>
		<p class="mt-2 text-center text-sm text-gray-600">
			Or 
			<a href="/login" class="font-medium text-blue-600 hover:text-blue-500 transition-colors">
				sign in to your existing account
			</a>
		</p>
	</div>

	<div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
		<div class="bg-white py-8 px-4 shadow-lg sm:rounded-lg sm:px-10 border border-gray-200">
			<form class="space-y-6" on:submit|preventDefault={handleRegister}>
				<!-- Username Field -->
				<div>
					<label for="username" class="block text-sm font-medium text-gray-700">
						Username
					</label>
					<div class="mt-1">
						<input
							id="username"
							name="username"
							type="text"
							autocomplete="username"
							required
							bind:value={username}
							on:keypress={handleKeyPress}
							class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
							placeholder="Choose a username"
						/>
					</div>
					<p class="mt-1 text-xs text-gray-500">At least 3 characters</p>
				</div>

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
								autocomplete="new-password"
								required
								bind:value={password}
								on:keypress={handleKeyPress}
								class="appearance-none block w-full px-3 py-2 pr-10 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Create a password"
							/>
						{:else}
							<input
								id="password"
								name="password"
								type="password"
								autocomplete="new-password"
								required
								bind:value={password}
								on:keypress={handleKeyPress}
								class="appearance-none block w-full px-3 py-2 pr-10 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Create a password"
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
					<p class="mt-1 text-xs text-gray-500">At least 8 characters</p>
				</div>

				<!-- Confirm Password Field -->
				<div>
					<label for="confirmPassword" class="block text-sm font-medium text-gray-700">
						Confirm Password
					</label>
					<div class="mt-1 relative">
						{#if showConfirmPassword}
							<input
								id="confirmPassword"
								name="confirmPassword"
								type="text"
								autocomplete="new-password"
								required
								bind:value={confirmPassword}
								on:keypress={handleKeyPress}
								class="appearance-none block w-full px-3 py-2 pr-10 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Confirm your password"
							/>
						{:else}
							<input
								id="confirmPassword"
								name="confirmPassword"
								type="password"
								autocomplete="new-password"
								required
								bind:value={confirmPassword}
								on:keypress={handleKeyPress}
								class="appearance-none block w-full px-3 py-2 pr-10 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Confirm your password"
							/>
						{/if}
						<IconButton
							variant="ghost"
							size="sm"
							on:click={() => showConfirmPassword = !showConfirmPassword}
							class="absolute inset-y-0 right-0 pr-3 text-gray-400 hover:text-gray-600"
						>
							{#if showConfirmPassword}
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
						{isLoading ? 'Creating account...' : 'Create account'}
					</Button>
				</div>
			</form>

			<!-- Terms -->
			<div class="mt-6">
				<p class="text-xs text-center text-gray-500">
					By creating an account, you agree to our 
					<a href="#" class="text-blue-600 hover:text-blue-500">Terms of Service</a> 
					and 
					<a href="#" class="text-blue-600 hover:text-blue-500">Privacy Policy</a>
				</p>
			</div>

			<!-- Divider -->
			<div class="mt-6">
				<div class="relative">
					<div class="absolute inset-0 flex items-center">
						<div class="w-full border-t border-gray-300" />
					</div>
					<div class="relative flex justify-center text-sm">
						<span class="px-2 bg-white text-gray-500">Already have an account?</span>
					</div>
				</div>

				<div class="mt-6">
					<a
						href="/login"
						class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 transition-colors"
					>
						Sign in instead
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