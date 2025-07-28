<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { toastStore } from '$lib/stores/toast';
	import JsonEditor from '$lib/components/JsonEditor.svelte';

	export let show = false;
	export let isCreating = false;

	const dispatch = createEventDispatcher();

	let webhookName = '';
	let webhookDescription = '';
	let expiresInDays = 0;
	let selectedMethods = ['POST'];
	let responseStatus = 200;
	let responseHeaders = [{ key: '', value: '' }];
	let responseBody = '{\n  "message": "Webhook received successfully",\n  "timestamp": "{{timestamp}}"\n}';
	let responseDelay = 0;
	let contentType = 'application/json';
	let corsEnabled = true;
	let authEnabled = false;
	let authType = 'none';
	let authUsername = '';
	let authPassword = '';
	let authToken = '';

	const availableMethods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'];
	const authTypes = [
		{ value: 'none', label: 'No Authentication' },
		{ value: 'basic', label: 'Basic Auth' },
		{ value: 'bearer', label: 'Bearer Token' }
	];

	function close() {
		show = false;
		resetForm();
	}

	function resetForm() {
		webhookName = '';
		webhookDescription = '';
		expiresInDays = 0;
		selectedMethods = ['POST'];
		responseStatus = 200;
		responseHeaders = [{ key: '', value: '' }];
		responseBody = '{\n  "message": "Webhook received successfully",\n  "timestamp": "{{timestamp}}"\n}';
		responseDelay = 0;
		contentType = 'application/json';
		corsEnabled = true;
		authEnabled = false;
		authType = 'none';
		authUsername = '';
		authPassword = '';
		authToken = '';
	}

	function addResponseHeader() {
		responseHeaders = [...responseHeaders, { key: '', value: '' }];
	}

	function removeResponseHeader(index: number) {
		responseHeaders = responseHeaders.filter((_, i) => i !== index);
	}

	function toggleMethod(method: string) {
		if (selectedMethods.includes(method)) {
			selectedMethods = selectedMethods.filter(m => m !== method);
		} else {
			selectedMethods = [...selectedMethods, method];
		}
	}

	function createWebhook() {
		if (!webhookName.trim()) {
			toastStore.error('Webhook name is required');
			return;
		}

		if (selectedMethods.length === 0) {
			toastStore.error('At least one HTTP method must be selected');
			return;
		}

		// Build response headers object
		const headers: Record<string, string> = {};
		responseHeaders.forEach(({ key, value }) => {
			if (key.trim() && value.trim()) {
				headers[key] = value;
			}
		});

		// Build authentication config
		const authentication = {
			enabled: authEnabled,
			type: authType,
			...(authType === 'basic' && authEnabled ? { username: authUsername, password: authPassword } : {}),
			...(authType === 'bearer' && authEnabled ? { token: authToken } : {})
		};

		const webhookData = {
			name: webhookName.trim(),
			description: webhookDescription.trim(),
			expires_in_days: expiresInDays > 0 ? expiresInDays : undefined,
			config: {
				methods: selectedMethods,
				response_status: responseStatus,
				response_headers: headers,
				response_body: responseBody,
				response_delay: responseDelay,
				content_type: contentType,
				cors_enabled: corsEnabled,
				log_requests: true,
				authentication
			}
		};

		dispatch('create', webhookData);
	}

	function setDefaultResponse(type: string) {
		switch (type) {
			case 'success':
				responseBody = '{\n  "message": "Webhook received successfully",\n  "timestamp": "{{timestamp}}"\n}';
				responseStatus = 200;
				break;
			case 'echo':
				responseBody = '{\n  "echo": "{{request_body}}",\n  "received_at": "{{timestamp}}"\n}';
				responseStatus = 200;
				break;
			case 'error':
				responseBody = '{\n  "error": "Internal server error",\n  "code": 500\n}';
				responseStatus = 500;
				break;
			case 'empty':
				responseBody = '';
				responseStatus = 204;
				break;
		}
	}
</script>

{#if show}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" on:click|self={close}>
		<div class="bg-gray-800 rounded-lg w-full max-w-4xl max-h-[90vh] overflow-y-auto m-4">
			<!-- Header -->
			<div class="sticky top-0 bg-gray-800 border-b border-gray-700 p-6 rounded-t-lg">
				<div class="flex items-center justify-between">
					<h2 class="text-lg font-semibold text-white flex items-center">
						<svg class="w-5 h-5 mr-2 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path>
						</svg>
						Create New Webhook
					</h2>
					<button on:click={close} class="text-gray-400 hover:text-white transition-colors">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
						</svg>
					</button>
				</div>
			</div>

			<!-- Content -->
			<div class="p-6 space-y-6">
				<!-- Basic Information -->
				<div class="bg-gray-700/30 rounded-lg p-4 space-y-4">
					<h3 class="text-base font-medium text-white">Basic Information</h3>
					
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div>
							<label class="block text-sm font-medium text-gray-300 mb-2">
								Webhook Name <span class="text-red-400">*</span>
							</label>
							<input
								type="text"
								bind:value={webhookName}
								placeholder="My API Webhook"
								class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
							/>
						</div>

						<div>
							<label class="block text-sm font-medium text-gray-300 mb-2">Expires in Days (Optional)</label>
							<input
								type="number"
								bind:value={expiresInDays}
								min="0"
								max="365"
								placeholder="0 = Never expires"
								class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
							/>
						</div>
					</div>

					<div>
						<label class="block text-sm font-medium text-gray-300 mb-2">Description (Optional)</label>
						<textarea
							bind:value={webhookDescription}
							placeholder="Webhook for receiving payment notifications..."
							rows="2"
							class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 resize-none"
						></textarea>
					</div>
				</div>

				<!-- HTTP Methods -->
				<div class="bg-gray-700/30 rounded-lg p-4 space-y-4">
					<h3 class="text-base font-medium text-white">Allowed HTTP Methods</h3>
					<div class="grid grid-cols-3 md:grid-cols-5 gap-3">
						{#each availableMethods as method}
							<label class="flex items-center space-x-2 cursor-pointer">
								<input
									type="checkbox"
									checked={selectedMethods.includes(method)}
									on:change={() => toggleMethod(method)}
									class="w-4 h-4 text-purple-600 bg-gray-700 border-gray-600 rounded focus:ring-purple-500"
								/>
								<span class="text-sm text-gray-300 font-medium">{method}</span>
							</label>
						{/each}
					</div>
				</div>

				<!-- Response Configuration -->
				<div class="bg-gray-700/30 rounded-lg p-4 space-y-4">
					<div class="flex items-center justify-between">
						<h3 class="text-base font-medium text-white">Response Configuration</h3>
						<div class="flex space-x-2">
							<button
								on:click={() => setDefaultResponse('success')}
								class="px-2 py-1 text-xs bg-green-600 hover:bg-green-700 text-white rounded transition-colors"
							>
								Success
							</button>
							<button
								on:click={() => setDefaultResponse('echo')}
								class="px-2 py-1 text-xs bg-blue-600 hover:bg-blue-700 text-white rounded transition-colors"
							>
								Echo
							</button>
							<button
								on:click={() => setDefaultResponse('error')}
								class="px-2 py-1 text-xs bg-red-600 hover:bg-red-700 text-white rounded transition-colors"
							>
								Error
							</button>
							<button
								on:click={() => setDefaultResponse('empty')}
								class="px-2 py-1 text-xs bg-gray-600 hover:bg-gray-700 text-white rounded transition-colors"
							>
								Empty
							</button>
						</div>
					</div>

					<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
						<div>
							<label class="block text-sm font-medium text-gray-300 mb-2">Status Code</label>
							<input
								type="number"
								bind:value={responseStatus}
								min="100"
								max="599"
								class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
							/>
						</div>

						<div>
							<label class="block text-sm font-medium text-gray-300 mb-2">Content Type</label>
							<select
								bind:value={contentType}
								class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
							>
								<option value="application/json">application/json</option>
								<option value="text/plain">text/plain</option>
								<option value="text/html">text/html</option>
								<option value="application/xml">application/xml</option>
								<option value="application/x-www-form-urlencoded">application/x-www-form-urlencoded</option>
							</select>
						</div>

						<div>
							<label class="block text-sm font-medium text-gray-300 mb-2">Response Delay (ms)</label>
							<input
								type="number"
								bind:value={responseDelay}
								min="0"
								max="30000"
								placeholder="0"
								class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
							/>
						</div>
					</div>

					<!-- Response Headers -->
					<div>
						<div class="flex items-center justify-between mb-2">
							<label class="block text-sm font-medium text-gray-300">Response Headers</label>
							<button
								on:click={addResponseHeader}
								class="text-xs text-purple-400 hover:text-purple-300"
							>
								Add Header
							</button>
						</div>
						<div class="space-y-2">
							{#each responseHeaders as header, index}
								<div class="flex space-x-2">
									<input
										type="text"
										bind:value={header.key}
										placeholder="Header name"
										class="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
									/>
									<input
										type="text"
										bind:value={header.value}
										placeholder="Header value"
										class="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
									/>
									<button
										on:click={() => removeResponseHeader(index)}
										class="px-2 py-2 text-red-400 hover:text-red-300"
									>
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
										</svg>
									</button>
								</div>
							{/each}
						</div>
					</div>

					<!-- Response Body -->
					<div>
						<label class="block text-sm font-medium text-gray-300 mb-2">Response Body</label>
						<JsonEditor
							bind:value={responseBody}
							height="150px"
							placeholder="Enter response body content..."
						/>
						<p class="text-xs text-gray-500 mt-1">
							Use <code>{'{{timestamp}}'}</code> for current timestamp, <code>{'{{request_body}}'}</code> for echoing request body
						</p>
					</div>
				</div>

				<!-- Authentication -->
				<div class="bg-gray-700/30 rounded-lg p-4 space-y-4">
					<div class="flex items-center justify-between">
						<h3 class="text-base font-medium text-white">Authentication</h3>
						<label class="relative inline-flex items-center cursor-pointer">
							<input type="checkbox" bind:checked={authEnabled} class="sr-only peer">
							<div class="w-11 h-6 bg-gray-600 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-600"></div>
						</label>
					</div>

					{#if authEnabled}
						<div class="space-y-4">
							<div>
								<label class="block text-sm font-medium text-gray-300 mb-2">Authentication Type</label>
								<select
									bind:value={authType}
									class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
								>
									{#each authTypes as type}
										<option value={type.value}>{type.label}</option>
									{/each}
								</select>
							</div>

							{#if authType === 'basic'}
								<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
									<div>
										<label class="block text-sm font-medium text-gray-300 mb-2">Username</label>
										<input
											type="text"
											bind:value={authUsername}
											placeholder="username"
											class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
										/>
									</div>
									<div>
										<label class="block text-sm font-medium text-gray-300 mb-2">Password</label>
										<input
											type="password"
											bind:value={authPassword}
											placeholder="password"
											class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
										/>
									</div>
								</div>
							{:else if authType === 'bearer'}
								<div>
									<label class="block text-sm font-medium text-gray-300 mb-2">Bearer Token</label>
									<input
										type="text"
										bind:value={authToken}
										placeholder="your-bearer-token"
										class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-sm text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
									/>
								</div>
							{/if}
						</div>
					{:else}
						<p class="text-gray-400 text-sm">No authentication required. Anyone with the webhook URL can send requests.</p>
					{/if}
				</div>

				<!-- Advanced Options -->
				<div class="bg-gray-700/30 rounded-lg p-4 space-y-4">
					<h3 class="text-base font-medium text-white">Advanced Options</h3>
					
					<div class="flex items-center justify-between">
						<div>
							<div class="text-white font-medium">Enable CORS</div>
							<div class="text-gray-400 text-sm">Allow cross-origin requests from web browsers</div>
						</div>
						<label class="relative inline-flex items-center cursor-pointer">
							<input type="checkbox" bind:checked={corsEnabled} class="sr-only peer">
							<div class="w-11 h-6 bg-gray-600 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-600"></div>
						</label>
					</div>
				</div>
			</div>

			<!-- Footer -->
			<div class="sticky bottom-0 bg-gray-800 border-t border-gray-700 p-6 rounded-b-lg">
				<div class="flex items-center justify-end space-x-3">
					<button
						on:click={close}
						class="px-4 py-2 text-sm text-gray-300 hover:text-white transition-colors"
					>
						Cancel
					</button>
					<button
						on:click={createWebhook}
						disabled={isCreating || !webhookName.trim() || selectedMethods.length === 0}
						class="px-6 py-2 bg-purple-600 hover:bg-purple-700 disabled:opacity-50 disabled:cursor-not-allowed text-white text-sm font-medium rounded transition-colors flex items-center"
					>
						{#if isCreating}
							<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							Creating...
						{:else}
							Create Webhook
						{/if}
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}