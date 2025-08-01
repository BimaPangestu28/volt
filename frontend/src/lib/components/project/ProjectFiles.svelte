<script lang="ts">
	import { onMount } from 'svelte';
	import type { Project } from '$lib/stores/project';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import { Plus, Upload, ArrowUpDown, Grid3X3, List, FolderPlus, Download, Share2, X } from 'lucide-svelte';
	
	export let project: Project;
	
	interface ProjectFile {
		id: string;
		name: string;
		type: 'image' | 'document' | 'code' | 'archive' | 'other';
		size: number;
		uploadedBy: string;
		uploadedAt: Date;
		url: string;
		description?: string;
		tags: string[];
		isShared: boolean;
		downloadCount: number;
		version: string;
		folder?: string;
	}
	
	interface FileFolder {
		id: string;
		name: string;
		parent?: string;
		createdAt: Date;
		fileCount: number;
	}
	
	let files: ProjectFile[] = [];
	let folders: FileFolder[] = [];
	let currentFolder: string | null = null;
	let viewMode: 'grid' | 'list' = 'grid';
	let searchQuery = '';
	let selectedFiles: string[] = [];
	let showUploadModal = false;
	let showCreateFolderModal = false;
	let newFolderName = '';
	let sortBy: 'name' | 'date' | 'size' | 'type' = 'date';
	let sortOrder: 'asc' | 'desc' = 'desc';
	let filterType: 'all' | 'image' | 'document' | 'code' | 'archive' = 'all';
	
	// Sample files for demonstration
	const sampleFiles: ProjectFile[] = [
		{
			id: '1',
			name: 'API_Documentation.pdf',
			type: 'document',
			size: 1024000,
			uploadedBy: 'Sarah Chen',
			uploadedAt: new Date(Date.now() - 2 * 60 * 60 * 1000),
			url: '#',
			description: 'Complete API documentation with examples',
			tags: ['documentation', 'api'],
			isShared: true,
			downloadCount: 15,
			version: 'v2.1',
			folder: 'docs'
		},
		{
			id: '2',
			name: 'database_schema.sql',
			type: 'code',
			size: 45600,
			uploadedBy: 'Mike Johnson',
			uploadedAt: new Date(Date.now() - 4 * 60 * 60 * 1000),
			url: '#',
			description: 'Database schema for user management',
			tags: ['database', 'schema'],
			isShared: true,
			downloadCount: 8,
			version: 'v1.3',
			folder: 'database'
		},
		{
			id: '3',
			name: 'wireframes_mobile.png',
			type: 'image',
			size: 2500000,
			uploadedBy: 'Alex Kumar',
			uploadedAt: new Date(Date.now() - 6 * 60 * 60 * 1000),
			url: '#',
			description: 'Mobile app wireframes for user interface',
			tags: ['design', 'mobile', 'wireframes'],
			isShared: false,
			downloadCount: 3,
			version: 'v1.0',
			folder: 'designs'
		},
		{
			id: '4',
			name: 'project_backup.zip',
			type: 'archive',
			size: 15600000,
			uploadedBy: 'You',
			uploadedAt: new Date(Date.now() - 24 * 60 * 60 * 1000),
			url: '#',
			description: 'Complete project backup including source code',
			tags: ['backup', 'source'],
			isShared: false,
			downloadCount: 1,
			version: 'v1.0'
		},
		{
			id: '5',
			name: 'meeting_notes.md',
			type: 'document',
			size: 8900,
			uploadedBy: 'Sarah Chen',
			uploadedAt: new Date(Date.now() - 48 * 60 * 60 * 1000),
			url: '#',
			description: 'Weekly project meeting notes and action items',
			tags: ['meeting', 'notes'],
			isShared: true,
			downloadCount: 12,
			version: 'v1.0',
			folder: 'meetings'
		}
	];
	
	const sampleFolders: FileFolder[] = [
		{
			id: 'docs',
			name: 'Documentation',
			createdAt: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000),
			fileCount: 8
		},
		{
			id: 'designs',
			name: 'Designs',
			createdAt: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000),
			fileCount: 12
		},
		{
			id: 'database',
			name: 'Database',
			createdAt: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000),
			fileCount: 4
		},
		{
			id: 'meetings',
			name: 'Meeting Notes',
			createdAt: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000),
			fileCount: 6
		}
	];
	
	onMount(() => {
		files = sampleFiles;
		folders = sampleFolders;
	});
	
	function getFileIcon(type: ProjectFile['type']) {
		switch (type) {
			case 'image':
				return '🖼️';
			case 'document':
				return '📄';
			case 'code':
				return '💻';
			case 'archive':
				return '📦';
			default:
				return '📎';
		}
	}
	
	function getFileTypeColor(type: ProjectFile['type']) {
		switch (type) {
			case 'image':
				return 'text-purple-400';
			case 'document':
				return 'text-blue-400';
			case 'code':
				return 'text-green-400';
			case 'archive':
				return 'text-orange-400';
			default:
				return 'text-gray-400';
		}
	}
	
	function formatFileSize(bytes: number): string {
		if (bytes === 0) return '0 Bytes';
		
		const k = 1024;
		const sizes = ['Bytes', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}
	
	function formatDate(date: Date) {
		return new Intl.DateTimeFormat('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: 'numeric',
			minute: '2-digit'
		}).format(new Date(date));
	}
	
	function sortFiles(files: ProjectFile[]) {
		return [...files].sort((a, b) => {
			let comparison = 0;
			
			switch (sortBy) {
				case 'name':
					comparison = a.name.localeCompare(b.name);
					break;
				case 'date':
					comparison = new Date(a.uploadedAt).getTime() - new Date(b.uploadedAt).getTime();
					break;
				case 'size':
					comparison = a.size - b.size;
					break;
				case 'type':
					comparison = a.type.localeCompare(b.type);
					break;
			}
			
			return sortOrder === 'asc' ? comparison : -comparison;
		});
	}
	
	function filterFiles(files: ProjectFile[]) {
		let filtered = files;
		
		// Filter by current folder
		if (currentFolder) {
			filtered = filtered.filter(file => file.folder === currentFolder);
		} else {
			// Show only files not in folders when in root
			filtered = filtered.filter(file => !file.folder);
		}
		
		// Filter by type
		if (filterType !== 'all') {
			filtered = filtered.filter(file => file.type === filterType);
		}
		
		// Filter by search query
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase();
			filtered = filtered.filter(file => 
				file.name.toLowerCase().includes(query) ||
				file.description?.toLowerCase().includes(query) ||
				file.tags.some(tag => tag.toLowerCase().includes(query))
			);
		}
		
		return sortFiles(filtered);
	}
	
	function toggleFileSelection(fileId: string) {
		if (selectedFiles.includes(fileId)) {
			selectedFiles = selectedFiles.filter(id => id !== fileId);
		} else {
			selectedFiles = [...selectedFiles, fileId];
		}
	}
	
	function selectAllFiles() {
		const visibleFiles = filterFiles(files);
		selectedFiles = visibleFiles.map(f => f.id);
	}
	
	function clearSelection() {
		selectedFiles = [];
	}
	
	function downloadFile(file: ProjectFile) {
		// Update download count
		files = files.map(f => 
			f.id === file.id ? { ...f, downloadCount: f.downloadCount + 1 } : f
		);
		
		// Simulate download
		const link = document.createElement('a');
		link.href = file.url;
		link.download = file.name;
		link.click();
	}
	
	function downloadSelectedFiles() {
		const selectedFileData = files.filter(f => selectedFiles.includes(f.id));
		selectedFileData.forEach(file => downloadFile(file));
		clearSelection();
	}
	
	function deleteSelectedFiles() {
		if (confirm(`Are you sure you want to delete ${selectedFiles.length} file(s)?`)) {
			files = files.filter(f => !selectedFiles.includes(f.id));
			clearSelection();
		}
	}
	
	function openFolder(folderId: string) {
		currentFolder = folderId;
		clearSelection();
	}
	
	function goToRoot() {
		currentFolder = null;
		clearSelection();
	}
	
	function createFolder() {
		if (!newFolderName.trim()) return;
		
		const newFolder: FileFolder = {
			id: `folder_${Date.now()}`,
			name: newFolderName.trim(),
			parent: currentFolder || undefined,
			createdAt: new Date(),
			fileCount: 0
		};
		
		folders = [...folders, newFolder];
		newFolderName = '';
		showCreateFolderModal = false;
	}
	
	function getCurrentFolderName() {
		if (!currentFolder) return null;
		return folders.find(f => f.id === currentFolder)?.name;
	}
	
	$: filteredFiles = filterFiles(files);
	$: currentFolderFiles = currentFolder ? files.filter(f => f.folder === currentFolder) : folders;
	$: totalSize = filteredFiles.reduce((sum, file) => sum + file.size, 0);
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div class="flex items-center space-x-4">
			<h3 class="text-lg font-semibold text-white">Project Files</h3>
			{#if currentFolder}
				<nav class="flex items-center space-x-2 text-sm">
					<Button on:click={goToRoot} variant="ghost" size="sm" class="text-blue-400 hover:text-blue-300">
						Root
					</Button>
					<span class="text-gray-500">/</span>
					<span class="text-gray-300">{getCurrentFolderName()}</span>
				</nav>
			{/if}
		</div>
		
		<div class="flex items-center space-x-2">
			<Button
				on:click={() => showCreateFolderModal = true}
				variant="secondary"
				size="sm"
				class="bg-gray-700 hover:bg-gray-600 text-white flex items-center space-x-2"
			>
				<FolderPlus class="btn-icon-sm" />
				<span>New Folder</span>
			</Button>
			
			<Button
				on:click={() => showUploadModal = true}
				variant="primary"
				size="sm"
				class="flex items-center space-x-2"
			>
				<Upload class="btn-icon-sm" />
				<span>Upload Files</span>
			</Button>
		</div>
	</div>
	
	<!-- Toolbar -->
	<div class="bg-gray-800 rounded-lg border border-gray-700 p-4">
		<div class="flex items-center justify-between">
			<div class="flex items-center space-x-4">
				<!-- Search -->
				<div class="relative">
					<input
						type="text"
						bind:value={searchQuery}
						placeholder="Search files..."
						class="bg-gray-700 border border-gray-600 text-white placeholder-gray-400 px-3 py-2 pl-9 text-sm rounded focus:outline-none focus:ring-2 focus:ring-blue-500 w-64"
					/>
					<svg class="absolute left-3 top-2.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
					</svg>
				</div>
				
				<!-- Filter -->
				<select
					bind:value={filterType}
					class="bg-gray-700 border border-gray-600 text-white text-sm rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					<option value="all">All Types</option>
					<option value="image">Images</option>
					<option value="document">Documents</option>
					<option value="code">Code</option>
					<option value="archive">Archives</option>
				</select>
				
				<!-- Sort -->
				<select
					bind:value={sortBy}
					class="bg-gray-700 border border-gray-600 text-white text-sm rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					<option value="date">Sort by Date</option>
					<option value="name">Sort by Name</option>
					<option value="size">Sort by Size</option>
					<option value="type">Sort by Type</option>
				</select>
				
				<IconButton
					on:click={() => sortOrder = sortOrder === 'asc' ? 'desc' : 'asc'}
					variant="ghost"
					size="sm"
					title="Toggle sort order"
					class="text-gray-400 hover:text-gray-200"
				>
					<ArrowUpDown class="btn-icon-sm {sortOrder === 'desc' ? 'transform rotate-180' : ''}" />
				</IconButton>
			</div>
			
			<div class="flex items-center space-x-4">
				{#if selectedFiles.length > 0}
					<div class="flex items-center space-x-2">
						<span class="text-sm text-gray-400">{selectedFiles.length} selected</span>
						<Button
							on:click={downloadSelectedFiles}
							variant="primary"
							size="sm"
						>
							Download
						</Button>
						<Button
							on:click={deleteSelectedFiles}
							variant="danger"
							size="sm"
						>
							Delete
						</Button>
						<Button
							on:click={clearSelection}
							variant="secondary"
							size="sm"
							class="bg-gray-600 hover:bg-gray-500 text-white"
						>
							Clear
						</Button>
					</div>
				{:else}
					<div class="text-sm text-gray-400">
						{filteredFiles.length} files • {formatFileSize(totalSize)}
					</div>
				{/if}
				
				<!-- View mode toggle -->
				<div class="flex items-center bg-gray-700 rounded p-1">
					<IconButton
						on:click={() => viewMode = 'grid'}
						variant={viewMode === 'grid' ? 'primary' : 'ghost'}
						size="sm"
						class="rounded {viewMode === 'grid' ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-gray-200'}"
					>
						<Grid3X3 class="btn-icon-sm" />
					</IconButton>
					<IconButton
						on:click={() => viewMode = 'list'}
						variant={viewMode === 'list' ? 'primary' : 'ghost'}
						size="sm"
						class="rounded {viewMode === 'list' ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-gray-200'}"
					>
						<List class="btn-icon-sm" />
					</IconButton>
				</div>
			</div>
		</div>
	</div>
	
	<!-- File/Folder List -->
	<div class="bg-gray-800 rounded-lg border border-gray-700">
		{#if currentFolder === null}
			<!-- Show folders in root -->
			{#if folders.length > 0}
				<div class="border-b border-gray-700 p-4">
					<h4 class="text-sm font-medium text-gray-300 mb-3">Folders</h4>
					<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
						{#each folders as folder}
							<Button
								on:click={() => openFolder(folder.id)}
								variant="ghost"
								size="md"
								class="flex items-center space-x-3 p-3 bg-gray-750 hover:bg-gray-700 rounded-lg text-left w-full justify-start"
							>
								<div class="text-2xl">📁</div>
								<div class="flex-1 min-w-0">
									<p class="text-sm font-medium text-white truncate">{folder.name}</p>
									<p class="text-xs text-gray-400">{folder.fileCount} files</p>
								</div>
							</Button>
						{/each}
					</div>
				</div>
			{/if}
		{/if}
		
		{#if filteredFiles.length === 0}
			<!-- Empty state -->
			<div class="text-center py-12">
				<svg class="w-12 h-12 text-gray-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"></path>
				</svg>
				<h4 class="text-lg font-medium text-white mb-2">No files found</h4>
				<p class="text-gray-400 mb-4">
					{searchQuery ? 'Try adjusting your search or filters' : 'Upload your first file to get started'}
				</p>
				{#if !searchQuery}
					<Button
						on:click={() => showUploadModal = true}
						variant="primary"
						size="md"
					>
						Upload Files
					</Button>
				{/if}
			</div>
		{:else if viewMode === 'grid'}
			<!-- Grid View -->
			<div class="p-4">
				{#if filteredFiles.length > 5}
					<div class="mb-4 flex items-center space-x-2">
						<Button
							on:click={selectAllFiles}
							variant="ghost"
							size="sm"
							class="text-blue-400 hover:text-blue-300"
						>
							Select All
						</Button>
						<span class="text-gray-500">•</span>
						<Button
							on:click={clearSelection}
							variant="ghost"
							size="sm"
							class="text-gray-400 hover:text-gray-300"
						>
							Clear Selection
						</Button>
					</div>
				{/if}
				
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
					{#each filteredFiles as file}
						<div class="bg-gray-750 rounded-lg p-4 hover:bg-gray-700 transition-colors group relative {selectedFiles.includes(file.id) ? 'ring-2 ring-blue-500' : ''}">
							<!-- Selection checkbox -->
							<div class="absolute top-2 left-2 opacity-0 group-hover:opacity-100 transition-opacity">
								<input
									type="checkbox"
									checked={selectedFiles.includes(file.id)}
									on:change={() => toggleFileSelection(file.id)}
									class="w-4 h-4 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500"
								/>
							</div>
							
							<!-- File icon and type -->
							<div class="text-center mb-3">
								<div class="text-4xl mb-2">{getFileIcon(file.type)}</div>
								<div class="text-xs {getFileTypeColor(file.type)} font-medium uppercase">{file.type}</div>
							</div>
							
							<!-- File info -->
							<div class="text-center">
								<h5 class="text-sm font-medium text-white truncate mb-1" title={file.name}>
									{file.name}
								</h5>
								<p class="text-xs text-gray-400 mb-2">{formatFileSize(file.size)}</p>
								
								{#if file.description}
									<p class="text-xs text-gray-500 line-clamp-2 mb-2">{file.description}</p>
								{/if}
								
								<div class="flex items-center justify-between text-xs text-gray-500">
									<span>{file.uploadedBy}</span>
									<span>{file.downloadCount} downloads</span>
								</div>
							</div>
							
							<!-- Actions -->
							<div class="mt-3 flex justify-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
								<Button
									on:click={() => downloadFile(file)}
									variant="primary"
									size="xs"
								>
									Download
								</Button>
								<Button 
									variant="secondary"
									size="xs"
									class="bg-gray-600 hover:bg-gray-500 text-white"
								>
									Share
								</Button>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{:else}
			<!-- List View -->
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead class="bg-gray-750 border-b border-gray-600">
						<tr>
							<th class="text-left p-3 text-xs font-medium text-gray-300 uppercase">
								<input
									type="checkbox"
									checked={selectedFiles.length === filteredFiles.length && filteredFiles.length > 0}
									on:change={(e) => e.target.checked ? selectAllFiles() : clearSelection()}
									class="w-4 h-4 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500"
								/>
							</th>
							<th class="text-left p-3 text-xs font-medium text-gray-300 uppercase">Name</th>
							<th class="text-left p-3 text-xs font-medium text-gray-300 uppercase">Type</th>
							<th class="text-left p-3 text-xs font-medium text-gray-300 uppercase">Size</th>
							<th class="text-left p-3 text-xs font-medium text-gray-300 uppercase">Uploaded By</th>
							<th class="text-left p-3 text-xs font-medium text-gray-300 uppercase">Date</th>
							<th class="text-left p-3 text-xs font-medium text-gray-300 uppercase">Downloads</th>
							<th class="text-left p-3 text-xs font-medium text-gray-300 uppercase">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-700">
						{#each filteredFiles as file}
							<tr class="hover:bg-gray-750 transition-colors {selectedFiles.includes(file.id) ? 'bg-blue-900/20' : ''}">
								<td class="p-3">
									<input
										type="checkbox"
										checked={selectedFiles.includes(file.id)}
										on:change={() => toggleFileSelection(file.id)}
										class="w-4 h-4 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500"
									/>
								</td>
								<td class="p-3">
									<div class="flex items-center space-x-3">
										<div class="text-lg">{getFileIcon(file.type)}</div>
										<div>
											<p class="text-sm font-medium text-white">{file.name}</p>
											{#if file.description}
												<p class="text-xs text-gray-400 truncate max-w-xs">{file.description}</p>
											{/if}
										</div>
									</div>
								</td>
								<td class="p-3">
									<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium {getFileTypeColor(file.type)} bg-gray-700 capitalize">
										{file.type}
									</span>
								</td>
								<td class="p-3 text-sm text-gray-300">{formatFileSize(file.size)}</td>
								<td class="p-3 text-sm text-gray-300">{file.uploadedBy}</td>
								<td class="p-3 text-sm text-gray-400">{formatDate(file.uploadedAt)}</td>
								<td class="p-3 text-sm text-gray-400">{file.downloadCount}</td>
								<td class="p-3">
									<div class="flex items-center space-x-2">
										<Button
											on:click={() => downloadFile(file)}
											variant="primary"
											size="xs"
										>
											Download
										</Button>
										<Button 
											variant="secondary"
											size="xs"
											class="bg-gray-600 hover:bg-gray-500 text-white"
										>
											Share
										</Button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>

<!-- Create Folder Modal -->
{#if showCreateFolderModal}
	<div class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center p-4 z-50">
		<div class="bg-gray-800 rounded-lg shadow-xl max-w-md w-full">
			<div class="px-5 py-3 border-b border-gray-700">
				<h3 class="text-base font-medium text-white">Create New Folder</h3>
			</div>
			
			<div class="px-5 py-4">
				<input
					type="text"
					bind:value={newFolderName}
					placeholder="Folder name"
					class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
				/>
			</div>
			
			<div class="px-5 py-3 bg-gray-700 flex justify-end space-x-2 rounded-b-lg">
				<Button
					on:click={() => showCreateFolderModal = false}
					variant="secondary"
					size="sm"
					class="bg-gray-600 border border-gray-500 text-gray-300 hover:bg-gray-500"
				>
					Cancel
				</Button>
				<Button
					on:click={createFolder}
					disabled={!newFolderName.trim()}
					variant="primary"
					size="sm"
				>
					Create Folder
				</Button>
			</div>
		</div>
	</div>
{/if}

<style>
	.bg-gray-750 {
		background-color: #374151;
	}
	
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
	
	.btn-icon-sm {
		width: 16px;
		height: 16px;
	}
</style>