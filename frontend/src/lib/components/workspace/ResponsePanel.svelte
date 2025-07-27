<script lang="ts">
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

	export let response: any = null;
	export let isExecuting = false;
	export let activeResponseTab = 'body';
	
	let responseViewMode = 'raw'; // 'raw', 'preview', 'visualize'

	function formatResponseBody(body: string) {
		try {
			return JSON.stringify(JSON.parse(body), null, 2);
		} catch {
			return body;
		}
	}

	function getStatusText(status: number): string {
		const statusTexts: Record<number, string> = {
			200: 'OK',
			201: 'Created',
			204: 'No Content',
			400: 'Bad Request',
			401: 'Unauthorized',
			403: 'Forbidden',
			404: 'Not Found',
			500: 'Internal Server Error'
		};
		return statusTexts[status] || 'Unknown';
	}

	function getResponseSize(body: any) {
		if (!body) return '0 B';
		const size = new Blob([typeof body === 'string' ? body : JSON.stringify(body)]).size;
		if (size < 1024) return `${size} B`;
		if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
		return `${(size / (1024 * 1024)).toFixed(1)} MB`;
	}

	function isJsonResponse(response: any) {
		const contentType = response?.headers?.['content-type'] || response?.headers?.['Content-Type'] || '';
		return contentType.includes('application/json') || contentType.includes('text/json');
	}

	function isHtmlResponse(response: any) {
		const contentType = response?.headers?.['content-type'] || response?.headers?.['Content-Type'] || '';
		return contentType.includes('text/html');
	}
</script>

{#if response || isExecuting}
	<div class="w-full bg-gray-900 flex flex-col h-full">
		<!-- Response Header -->
		<div class="border-b border-gray-700 bg-gray-800">
			<!-- Response Tabs -->
			<div class="flex items-center justify-between px-4 py-2">
				<nav class="flex space-x-6" aria-label="Tabs">
					<button
						class="py-2 px-1 border-b-2 font-normal text-[10px] {
							activeResponseTab === 'body' 
								? 'border-blue-500 text-blue-400' 
								: 'border-transparent text-gray-400 hover:text-gray-200 hover:border-gray-600'
						}"
						on:click={() => activeResponseTab = 'body'}
					>
						Body
					</button>
					<button
						class="py-2 px-1 border-b-2 font-normal text-[10px] {
							activeResponseTab === 'cookies' 
								? 'border-blue-500 text-blue-400' 
								: 'border-transparent text-gray-400 hover:text-gray-200 hover:border-gray-600'
						}"
						on:click={() => activeResponseTab = 'cookies'}
					>
						Cookies ({0})
					</button>
					<button
						class="py-2 px-1 border-b-2 font-normal text-[10px] {
							activeResponseTab === 'headers' 
								? 'border-blue-500 text-blue-400' 
								: 'border-transparent text-gray-400 hover:text-gray-200 hover:border-gray-600'
						}"
						on:click={() => activeResponseTab = 'headers'}
					>
						Headers ({Object.keys(response?.headers || {}).length})
					</button>
					<button
						class="py-2 px-1 border-b-2 font-normal text-[10px] {
							activeResponseTab === 'test-results' 
								? 'border-blue-500 text-blue-400' 
								: 'border-transparent text-gray-400 hover:text-gray-200 hover:border-gray-600'
						}"
						on:click={() => activeResponseTab = 'test-results'}
					>
						Test Results
					</button>
				</nav>
				
				{#if response}
					<div class="flex items-center space-x-4 text-[10px]">
						<span class="font-medium {
							response.status >= 200 && response.status < 300 ? 'text-green-400' : 
							response.status >= 400 ? 'text-red-400' : 
							'text-yellow-400'
						}">
							{response.status} {response.statusText || getStatusText(response.status)}
						</span>
						<span class="text-gray-300">{response.time} ms</span>
						<span class="text-gray-300">{getResponseSize(response.body)}</span>
					</div>
				{/if}
			</div>
		</div>

		{#if isExecuting}
			<!-- Loading State -->
			<div class="flex-1 flex items-center justify-center bg-gray-800">
				<div class="text-center">
					<LoadingSpinner size="lg" color="blue" />
					<p class="mt-3 text-sm text-gray-300">Sending request...</p>
				</div>
			</div>
		{:else if response}

			<!-- Response Content -->
			<div class="flex-1 overflow-hidden">
				{#if activeResponseTab === 'body'}
					<div class="h-full p-3 overflow-y-auto bg-gray-800">
						<div class="mb-2 flex items-center justify-between">
							<div class="flex items-center space-x-2">
								{#if isJsonResponse(response)}
									<button 
										on:click={() => responseViewMode = 'raw'}
										class="px-2 py-1 text-[10px] rounded font-mono border transition-all {
											responseViewMode === 'raw' 
												? 'bg-blue-900 text-blue-300 border-blue-700' 
												: 'bg-gray-700 text-gray-400 border-gray-600 hover:bg-blue-900/50 hover:text-blue-300'
										}"
									>
										JSON
									</button>
								{:else if isHtmlResponse(response)}
									<button 
										on:click={() => responseViewMode = 'raw'}
										class="px-2 py-1 text-[10px] rounded font-mono border transition-all {
											responseViewMode === 'raw' 
												? 'bg-orange-900 text-orange-300 border-orange-700' 
												: 'bg-gray-700 text-gray-400 border-gray-600 hover:bg-orange-900/50 hover:text-orange-300'
										}"
									>
										HTML
									</button>
								{:else}
									<button 
										on:click={() => responseViewMode = 'raw'}
										class="px-2 py-1 text-[10px] rounded font-mono border transition-all {
											responseViewMode === 'raw' 
												? 'bg-gray-600 text-gray-200 border-gray-500' 
												: 'bg-gray-700 text-gray-400 border-gray-600 hover:bg-gray-600 hover:text-gray-200'
										}"
									>
										Text
									</button>
								{/if}
								<button 
									on:click={() => responseViewMode = 'preview'}
									class="px-2 py-1 text-[10px] transition-colors {
										responseViewMode === 'preview' 
											? 'text-blue-400 font-medium' 
											: 'text-gray-400 hover:text-gray-200'
									}"
								>
									Preview
								</button>
								<button 
									on:click={() => responseViewMode = 'visualize'}
									class="px-2 py-1 text-[10px] transition-colors {
										responseViewMode === 'visualize' 
											? 'text-blue-400 font-medium' 
											: 'text-gray-400 hover:text-gray-200'
									}"
								>
									Visualize
								</button>
							</div>
							<div class="flex items-center space-x-2">
								<button class="p-1 text-gray-400 hover:text-gray-200" title="Pretty print">
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
									</svg>
								</button>
								<button class="p-1 text-gray-400 hover:text-gray-200" title="Wrap lines">
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
									</svg>
								</button>
								<button class="p-1 text-gray-400 hover:text-gray-200" title="Search">
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
									</svg>
								</button>
								<button class="p-1 text-gray-400 hover:text-gray-200" title="Copy">
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
									</svg>
								</button>
							</div>
						</div>
						<!-- Response Content Display -->
						{#if responseViewMode === 'raw'}
							<div class="bg-gray-850 rounded border border-gray-700 text-[10px] font-mono leading-relaxed text-gray-200 overflow-auto" style="height: calc(100% - 60px);">
								<pre class="p-3 m-0 whitespace-pre-wrap break-words max-w-full overflow-wrap-anywhere">{formatResponseBody(response.body)}</pre>
							</div>
						{:else if responseViewMode === 'preview' && isHtmlResponse(response)}
							<div class="bg-white rounded border border-gray-700 overflow-auto" style="height: calc(100% - 60px);">
								<iframe 
									srcdoc={response.body} 
									class="w-full h-full border-0"
									title="HTML Preview"
									sandbox="allow-scripts allow-same-origin"
								></iframe>
							</div>
						{:else if responseViewMode === 'preview'}
							<div class="bg-gray-850 rounded border border-gray-700 text-[10px] font-mono leading-relaxed text-gray-200 overflow-auto" style="height: calc(100% - 60px);">
								<pre class="p-3 m-0 whitespace-pre-wrap break-words max-w-full overflow-wrap-anywhere">{formatResponseBody(response.body)}</pre>
							</div>
						{:else if responseViewMode === 'visualize'}
							<div class="bg-gray-850 rounded border border-gray-700 flex items-center justify-center" style="height: calc(100% - 60px);">
								<div class="text-center p-8">
									<svg class="mx-auto h-12 w-12 text-gray-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
									</svg>
									<h3 class="text-[10px] font-medium text-gray-300 mb-2">Response Visualization</h3>
									<p class="text-[10px] text-gray-500">JSON and structured data visualization will be available here</p>
								</div>
							</div>
						{/if}
					</div>
				{:else if activeResponseTab === 'cookies'}
					<div class="h-full p-4 overflow-y-auto bg-gray-800">
						<div class="text-center py-12">
							<svg class="mx-auto h-8 w-8 text-gray-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4"></path>
							</svg>
							<p class="text-sm text-gray-400">No cookies in this response</p>
						</div>
					</div>
				{:else if activeResponseTab === 'headers'}
					<div class="h-full overflow-y-auto bg-gray-800">
						<div class="p-4">
							<div class="space-y-0">
								{#each Object.entries(response.headers || {}) as [key, value]}
									<div class="flex py-2 px-3 hover:bg-gray-750 border-b border-gray-700/30 last:border-b-0">
										<div class="flex-1 font-mono text-[10px]">
											<span class="text-orange-400 font-medium">{key}</span>
											<span class="text-gray-200 ml-2">{value}</span>
										</div>
									</div>
								{/each}
							</div>
						</div>
					</div>
				{:else if activeResponseTab === 'test-results'}
					<div class="h-full p-4 overflow-y-auto bg-gray-800">
						<div class="text-center py-12">
							<svg class="mx-auto h-8 w-8 text-gray-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
							</svg>
							<p class="text-sm text-gray-400">Test results will appear here after running tests</p>
						</div>
					</div>
				{/if}
			</div>
		{/if}
	</div>
{/if}