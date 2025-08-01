<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { apiClient } from '$lib/api/client';
	import { toastStore } from '$lib/stores/toast';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import { X, Copy, Link, Check, RefreshCw } from 'lucide-svelte';

	export let show = false;
	export let workspaceId = '';

	const dispatch = createEventDispatcher();

	let role = 'member';
	let expiresInHours = 168; // 7 days default
	let isGenerating = false;
	let error = '';
	let invitationLink = '';
	let showLinkResult = false;

	const roles = [
		{ value: 'member', label: 'Member', description: 'Can view and edit collections' },
		{ value: 'admin', label: 'Admin', description: 'Can manage workspace and members' },
		{ value: 'viewer', label: 'Viewer', description: 'Can only view collections' }
	];

	const expiryOptions = [
		{ value: 24, label: '24 hours' },
		{ value: 72, label: '3 days' },
		{ value: 168, label: '7 days' },
		{ value: 720, label: '30 days' },
		{ value: 0, label: 'Never expires' }
	];

	function close() {
		show = false;
		role = 'member';
		expiresInHours = 168;
		error = '';
		invitationLink = '';
		showLinkResult = false;
		dispatch('close');
	}

	async function generateInviteLink() {
		isGenerating = true;
		error = '';

		try {
			const result = await apiClient.generateInvitationLink(workspaceId, role, expiresInHours);
			invitationLink = result.inviteLink || result.link || `${window.location.origin}/invite/${result.token}`;
			showLinkResult = true;
			toastStore.success('Invitation link generated successfully!');
			dispatch('linkGenerated', { link: invitationLink, role, expiresInHours });
		} catch (err: any) {
			error = err.response?.data?.error || 'Failed to generate invitation link';
		} finally {
			isGenerating = false;
		}
	}

	async function copyToClipboard() {
		try {
			await navigator.clipboard.writeText(invitationLink);
			toastStore.success('Link copied to clipboard!');
		} catch (err) {
			toastStore.error('Failed to copy link');
		}
	}

	function generateNewLink() {
		showLinkResult = false;
		invitationLink = '';
		error = '';
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			close();
		} else if (event.key === 'Enter' && !isGenerating && !showLinkResult) {
			generateInviteLink();
		}
	}
</script>

<svelte:window on:keydown={handleKeyDown} />

{#if show}
	<!-- Modal Backdrop -->
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
		<div class="bg-gray-800 rounded-lg shadow-xl w-full max-w-md">
			<!-- Header -->
			<div class="flex items-center justify-between p-6 border-b border-gray-700">
				<h3 class="text-lg font-semibold text-white">
					{showLinkResult ? 'Invitation Link Generated' : 'Generate Invitation Link'}
				</h3>
				<IconButton
					on:click={close}
					variant="ghost"
					size="sm"
				>
					<X class="w-5 h-5" />
				</IconButton>
			</div>

			<!-- Body -->
			<div class="p-6 space-y-4">
				{#if !showLinkResult}
					<!-- Link Generation Form -->
					
					<!-- Role Selection -->
					<div>
						<label for="role" class="block text-sm font-medium text-gray-300 mb-2">
							Role
						</label>
						<select
							id="role"
							bind:value={role}
							class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							disabled={isGenerating}
						>
							{#each roles as roleOption}
								<option value={roleOption.value}>{roleOption.label}</option>
							{/each}
						</select>
						
						<!-- Role Description -->
						{#each roles as roleOption}
							{#if role === roleOption.value}
								<p class="mt-1 text-xs text-gray-400">{roleOption.description}</p>
							{/if}
						{/each}
					</div>

					<!-- Expiry Selection -->
					<div>
						<label for="expiry" class="block text-sm font-medium text-gray-300 mb-2">
							Link Expires
						</label>
						<select
							id="expiry"
							bind:value={expiresInHours}
							class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							disabled={isGenerating}
						>
							{#each expiryOptions as option}
								<option value={option.value}>{option.label}</option>
							{/each}
						</select>
					</div>

				{:else}
					<!-- Generated Link Result -->
					<div class="space-y-4">
						<div class="bg-green-900/20 border border-green-700 rounded-md p-4">
							<div class="flex items-center mb-2">
								<svg class="w-5 h-5 text-green-400 mr-2" fill="currentColor" viewBox="0 0 20 20">
									<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
								</svg>
								<h4 class="text-sm font-medium text-green-200">Invitation link generated!</h4>
							</div>
							<p class="text-xs text-green-300">Share this link with the person you want to invite to the workspace.</p>
						</div>

						<!-- Copy Link Box -->
						<div class="bg-gray-700 border border-gray-600 rounded-md p-3">
							<label class="block text-xs font-medium text-gray-300 mb-2">Invitation Link</label>
							<div class="flex items-center space-x-2">
								<input
									type="text"
									value={invitationLink}
									readonly
									class="flex-1 px-3 py-2 bg-gray-800 border border-gray-600 rounded-md text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
								/>
								<Button
									on:click={copyToClipboard}
									variant="primary"
									size="sm"
								>
									<Copy class="w-4 h-4 mr-1" />
									Copy
								</Button>
							</div>
						</div>

						<!-- Link Details -->
						<div class="text-xs text-gray-400 space-y-1">
							<p>Role: <span class="text-white font-medium">{roles.find(r => r.value === role)?.label}</span></p>
							<p>Expires: <span class="text-white font-medium">{expiresInHours === 0 ? 'Never' : expiryOptions.find(e => e.value === expiresInHours)?.label}</span></p>
						</div>
					</div>
				{/if}

				<!-- Error Message -->
				{#if error}
					<div class="bg-red-900/50 border border-red-700 rounded-md p-3">
						<div class="flex">
							<svg class="w-5 h-5 text-red-400 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
							</svg>
							<p class="text-sm text-red-200">{error}</p>
						</div>
					</div>
				{/if}
			</div>

			<!-- Footer -->
			<div class="flex items-center justify-end space-x-3 p-6 border-t border-gray-700">
				{#if !showLinkResult}
					<Button
						on:click={close}
						variant="ghost"
						size="sm"
						disabled={isGenerating}
					>
						Cancel
					</Button>
					<Button
						on:click={generateInviteLink}
						variant="primary"
						size="sm"
						loading={isGenerating}
						disabled={isGenerating}
					>
						{#if isGenerating}
							<RefreshCw class="w-4 h-4 mr-1" />
							Generating...
						{:else}
							<Link class="w-4 h-4 mr-1" />
							Generate Link
						{/if}
					</Button>
				{:else}
					<Button
						on:click={generateNewLink}
						variant="ghost"
						size="sm"
					>
						Generate New Link
					</Button>
					<Button
						on:click={close}
						variant="primary"
						size="sm"
						class="!bg-green-600 hover:!bg-green-700"
					>
						<Check class="w-4 h-4 mr-1" />
						Done
					</Button>
				{/if}
			</div>
		</div>
	</div>
{/if}