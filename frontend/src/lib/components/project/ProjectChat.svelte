<script lang="ts">
	import { onMount, afterUpdate } from 'svelte';
	import type { Project } from '$lib/stores/project';
	import Button from '$lib/components/ui/Button.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import { Reply, Smile, X, Paperclip, Send } from 'lucide-svelte';
	
	export let project: Project;
	
	interface ChatMessage {
		id: string;
		userId: string;
		userName: string;
		userAvatar?: string;
		message: string;
		timestamp: Date;
		type: 'message' | 'system' | 'file' | 'mention';
		edited?: boolean;
		reactions?: { emoji: string; users: string[] }[];
		replyTo?: string;
		attachments?: { name: string; url: string; type: string }[];
	}
	
	let messages: ChatMessage[] = [];
	let newMessage = '';
	let chatContainer: HTMLDivElement;
	let isTyping = false;
	let typingUsers: string[] = [];
	let replyingTo: ChatMessage | null = null;
	let showEmojiPicker = false;
	let selectedMessageId: string | null = null;
	
	// Sample messages for demonstration
	const sampleMessages: ChatMessage[] = [
		{
			id: '1',
			userId: 'user_2',
			userName: 'Sarah Chen',
			userAvatar: 'S',
			message: 'Hey team! I just finished the authentication endpoints. Ready for review 🚀',
			timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000),
			type: 'message'
		},
		{
			id: '2',
			userId: 'system',
			userName: 'System',
			message: 'Mike Johnson completed task "Setup Database Schema"',
			timestamp: new Date(Date.now() - 1.5 * 60 * 60 * 1000),
			type: 'system'
		},
		{
			id: '3',
			userId: 'user_3',
			userName: 'Mike Johnson',
			userAvatar: 'M',
			message: 'Great work @Sarah! I\'ll review it after lunch. Also uploaded the database migration files.',
			timestamp: new Date(Date.now() - 1 * 60 * 60 * 1000),
			type: 'message',
			attachments: [
				{ name: 'migration_001.sql', url: '#', type: 'sql' },
				{ name: 'schema_diagram.png', url: '#', type: 'image' }
			]
		},
		{
			id: '4',
			userId: 'user_1',
			userName: 'You',
			message: 'Thanks for the updates! The API endpoints look solid. Should we schedule a quick demo for the stakeholders?',
			timestamp: new Date(Date.now() - 30 * 60 * 1000),
			type: 'message',
			reactions: [
				{ emoji: '👍', users: ['user_2', 'user_3'] },
				{ emoji: '🎯', users: ['user_2'] }
			]
		},
		{
			id: '5',
			userId: 'user_4',
			userName: 'Alex Kumar',
			userAvatar: 'A',
			message: 'I can join the demo! Working on the mobile integration now.',
			timestamp: new Date(Date.now() - 15 * 60 * 1000),
			type: 'message'
		}
	];
	
	onMount(() => {
		messages = sampleMessages;
		scrollToBottom();
		
		// Simulate typing indicator
		setTimeout(() => {
			typingUsers = ['Sarah Chen'];
			setTimeout(() => {
				typingUsers = [];
			}, 3000);
		}, 5000);
	});
	
	afterUpdate(() => {
		if (chatContainer) {
			scrollToBottom();
		}
	});
	
	function scrollToBottom() {
		if (chatContainer) {
			chatContainer.scrollTop = chatContainer.scrollHeight;
		}
	}
	
	function sendMessage() {
		if (!newMessage.trim()) return;
		
		const message: ChatMessage = {
			id: Date.now().toString(),
			userId: 'user_1',
			userName: 'You',
			message: newMessage.trim(),
			timestamp: new Date(),
			type: 'message',
			replyTo: replyingTo?.id
		};
		
		messages = [...messages, message];
		newMessage = '';
		replyingTo = null;
		scrollToBottom();
	}
	
	function formatTime(date: Date) {
		return new Intl.DateTimeFormat('en-US', {
			hour: 'numeric',
			minute: '2-digit',
			hour12: true
		}).format(date);
	}
	
	function formatDate(date: Date) {
		const today = new Date();
		const yesterday = new Date(today);
		yesterday.setDate(yesterday.getDate() - 1);
		
		if (date.toDateString() === today.toDateString()) {
			return 'Today';
		} else if (date.toDateString() === yesterday.toDateString()) {
			return 'Yesterday';
		} else {
			return new Intl.DateTimeFormat('en-US', {
				month: 'short',
				day: 'numeric'
			}).format(date);
		}
	}
	
	function shouldShowDateSeparator(currentMsg: ChatMessage, previousMsg: ChatMessage | null) {
		if (!previousMsg) return true;
		
		const currentDate = new Date(currentMsg.timestamp).toDateString();
		const previousDate = new Date(previousMsg.timestamp).toDateString();
		
		return currentDate !== previousDate;
	}
	
	function replyToMessage(message: ChatMessage) {
		replyingTo = message;
		// Focus on input
		const input = document.querySelector('#messageInput') as HTMLInputElement;
		if (input) input.focus();
	}
	
	function cancelReply() {
		replyingTo = null;
	}
	
	function addReaction(messageId: string, emoji: string) {
		messages = messages.map(msg => {
			if (msg.id === messageId) {
				const reactions = msg.reactions || [];
				const existingReaction = reactions.find(r => r.emoji === emoji);
				
				if (existingReaction) {
					if (existingReaction.users.includes('user_1')) {
						// Remove reaction
						existingReaction.users = existingReaction.users.filter(u => u !== 'user_1');
						if (existingReaction.users.length === 0) {
							return {
								...msg,
								reactions: reactions.filter(r => r.emoji !== emoji)
							};
						}
					} else {
						// Add reaction
						existingReaction.users.push('user_1');
					}
				} else {
					// New reaction
					reactions.push({ emoji, users: ['user_1'] });
				}
				
				return { ...msg, reactions };
			}
			return msg;
		});
	}
	
	function getFileIcon(type: string) {
		switch (type) {
			case 'image':
				return '🖼️';
			case 'sql':
				return '🗃️';
			case 'pdf':
				return '📄';
			case 'code':
				return '💻';
			default:
				return '📎';
		}
	}
	
	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' && !event.shiftKey) {
			event.preventDefault();
			sendMessage();
		}
	}
	
	// Group messages by date
	$: groupedMessages = messages.reduce((groups, message, index) => {
		const previousMessage = index > 0 ? messages[index - 1] : null;
		
		if (shouldShowDateSeparator(message, previousMessage)) {
			const dateKey = formatDate(message.timestamp);
			if (!groups[dateKey]) {
				groups[dateKey] = [];
			}
			groups[dateKey].push(message);
		} else {
			const lastGroup = Object.keys(groups).pop();
			if (lastGroup) {
				groups[lastGroup].push(message);
			}
		}
		
		return groups;
	}, {} as Record<string, ChatMessage[]>);
</script>

<div class="flex flex-col h-[600px] bg-gray-800 rounded-lg border border-gray-700">
	<!-- Chat Header -->
	<div class="p-4 border-b border-gray-700">
		<div class="flex items-center justify-between">
			<h3 class="text-lg font-semibold text-white">Project Chat</h3>
			<div class="flex items-center space-x-2">
				<!-- Online members -->
				<div class="flex -space-x-2">
					{#each project.members.filter(m => m.status === 'online').slice(0, 3) as member}
						<div 
							class="w-6 h-6 rounded-full bg-blue-600 border-2 border-gray-800 flex items-center justify-center relative"
							title="{member.name} (online)"
						>
							{#if member.avatar}
								<img src={member.avatar} alt={member.name} class="w-full h-full rounded-full object-cover" />
							{:else}
								<span class="text-xs font-medium text-white">{member.name.charAt(0)}</span>
							{/if}
							<div class="absolute -bottom-0.5 -right-0.5 w-2 h-2 bg-green-500 border border-gray-800 rounded-full"></div>
						</div>
					{/each}
				</div>
				<span class="text-sm text-green-400">{project.members.filter(m => m.status === 'online').length} online</span>
			</div>
		</div>
	</div>
	
	<!-- Messages -->
	<div bind:this={chatContainer} class="flex-1 overflow-y-auto p-4 space-y-4">
		{#each Object.entries(groupedMessages) as [date, dayMessages]}
			<!-- Date separator -->
			<div class="flex items-center justify-center">
				<div class="bg-gray-700 px-3 py-1 rounded-full">
					<span class="text-xs text-gray-400">{date}</span>
				</div>
			</div>
			
			<!-- Messages for this date -->
			{#each dayMessages as message}
				<div class="flex items-start space-x-3 group {message.userId === 'user_1' ? 'flex-row-reverse space-x-reverse' : ''}">
					<!-- Avatar -->
					{#if message.type !== 'system'}
						<div class="flex-shrink-0">
							<div class="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center">
								{#if message.userAvatar}
									<span class="text-sm font-medium text-white">{message.userAvatar}</span>
								{:else}
									<span class="text-sm font-medium text-white">{message.userName.charAt(0)}</span>
								{/if}
							</div>
						</div>
					{/if}
					
					<!-- Message content -->
					<div class="flex-1 min-w-0 {message.userId === 'user_1' ? 'text-right' : ''} {message.type === 'system' ? 'text-center' : ''}">
						{#if message.type === 'system'}
							<!-- System message -->
							<div class="inline-flex items-center space-x-2 bg-gray-700 px-3 py-2 rounded-full">
								<svg class="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
								</svg>
								<span class="text-sm text-gray-300">{message.message}</span>
							</div>
						{:else}
							<!-- Regular message -->
							<div class="max-w-md {message.userId === 'user_1' ? 'ml-auto' : ''}">
								<!-- Reply indicator -->
								{#if message.replyTo}
									{@const replyMsg = messages.find(m => m.id === message.replyTo)}
									{#if replyMsg}
										<div class="mb-2 pl-3 border-l-2 border-blue-500 bg-gray-750 rounded p-2">
											<div class="text-xs text-blue-400 font-medium">{replyMsg.userName}</div>
											<div class="text-xs text-gray-400 truncate">{replyMsg.message}</div>
										</div>
									{/if}
								{/if}
								
								<!-- Message bubble -->
								<div class="relative {message.userId === 'user_1' ? 'bg-blue-600' : 'bg-gray-700'} rounded-lg p-3">
									<!-- User name and time -->
									<div class="flex items-center justify-between mb-1">
										<span class="text-xs font-medium {message.userId === 'user_1' ? 'text-blue-100' : 'text-gray-300'}">{message.userName}</span>
										<span class="text-xs {message.userId === 'user_1' ? 'text-blue-200' : 'text-gray-500'}">{formatTime(message.timestamp)}</span>
									</div>
									
									<!-- Message text -->
									<div class="text-sm {message.userId === 'user_1' ? 'text-white' : 'text-gray-100'}">{message.message}</div>
									
									<!-- Attachments -->
									{#if message.attachments}
										<div class="mt-3 space-y-2">
											{#each message.attachments as attachment}
												<div class="flex items-center space-x-2 bg-black bg-opacity-20 rounded p-2">
													<span class="text-lg">{getFileIcon(attachment.type)}</span>
													<span class="text-xs {message.userId === 'user_1' ? 'text-blue-100' : 'text-gray-300'} truncate">{attachment.name}</span>
													<Button 
														variant="ghost"
														size="xs"
														class="text-xs {message.userId === 'user_1' ? 'text-blue-200' : 'text-gray-400'} hover:underline p-0"
													>
														Download
													</Button>
												</div>
											{/each}
										</div>
									{/if}
									
									<!-- Message actions (show on hover) -->
									<div class="absolute top-2 {message.userId === 'user_1' ? 'left-2' : 'right-2'} opacity-0 group-hover:opacity-100 transition-opacity">
										<div class="flex items-center space-x-1">
											<IconButton
												on:click={() => replyToMessage(message)}
												variant="ghost"
												size="xs"
												title="Reply"
												class="hover:bg-black hover:bg-opacity-20 text-gray-300"
											>
												<Reply class="btn-icon-xs" />
											</IconButton>
											<IconButton
												on:click={() => selectedMessageId = selectedMessageId === message.id ? null : message.id}
												variant="ghost"
												size="xs"
												title="React"
												class="hover:bg-black hover:bg-opacity-20 text-gray-300"
											>
												<Smile class="btn-icon-xs" />
											</IconButton>
										</div>
									</div>
									
									<!-- Quick reactions -->
									{#if selectedMessageId === message.id}
										<div class="absolute top-full mt-1 {message.userId === 'user_1' ? 'right-0' : 'left-0'} bg-gray-900 border border-gray-600 rounded-lg p-2 flex space-x-1 z-10">
											{#each ['👍', '❤️', '😊', '🎉', '🚀', '👎'] as emoji}
												<Button
													on:click={() => addReaction(message.id, emoji)}
													variant="ghost"
													size="xs"
													class="p-1 hover:bg-gray-700 rounded"
												>
													{emoji}
												</Button>
											{/each}
										</div>
									{/if}
								</div>
								
								<!-- Reactions -->
								{#if message.reactions && message.reactions.length > 0}
									<div class="flex flex-wrap gap-1 mt-2 {message.userId === 'user_1' ? 'justify-end' : ''}">
										{#each message.reactions as reaction}
											<Button
												on:click={() => addReaction(message.id, reaction.emoji)}
												variant="ghost"
												size="xs"
												class="inline-flex items-center space-x-1 bg-gray-700 hover:bg-gray-600 px-2 py-1 rounded-full text-xs transition-colors {reaction.users.includes('user_1') ? 'ring-1 ring-blue-500' : ''}"
											>
												<span>{reaction.emoji}</span>
												<span class="text-gray-300">{reaction.users.length}</span>
											</Button>
										{/each}
									</div>
								{/if}
							</div>
						{/if}
					</div>
				</div>
			{/each}
		{/each}
		
		<!-- Typing indicator -->
		{#if typingUsers.length > 0}
			<div class="flex items-center space-x-3">
				<div class="w-8 h-8 rounded-full bg-gray-600 flex items-center justify-center">
					<div class="flex space-x-1">
						<div class="w-1 h-1 bg-gray-400 rounded-full animate-bounce"></div>
						<div class="w-1 h-1 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.1s"></div>
						<div class="w-1 h-1 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
					</div>
				</div>
				<div class="text-sm text-gray-400">
					{typingUsers.join(', ')} {typingUsers.length === 1 ? 'is' : 'are'} typing...
				</div>
			</div>
		{/if}
	</div>
	
	<!-- Message Input -->
	<div class="border-t border-gray-700 p-4">
		{#if replyingTo}
			<div class="mb-3 flex items-center justify-between bg-gray-750 rounded p-2">
				<div class="flex items-center space-x-2">
					<svg class="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"></path>
					</svg>
					<div>
						<div class="text-xs text-blue-400 font-medium">Replying to {replyingTo.userName}</div>
						<div class="text-xs text-gray-400 truncate max-w-md">{replyingTo.message}</div>
					</div>
				</div>
				<IconButton 
					on:click={cancelReply} 
					variant="ghost"
					size="sm"
					class="text-gray-400 hover:text-gray-200"
				>
					<X class="btn-icon-sm" />
				</IconButton>
			</div>
		{/if}
		
		<div class="flex items-end space-x-3">
			<div class="flex-1">
				<textarea
					id="messageInput"
					bind:value={newMessage}
					on:keydown={handleKeyDown}
					placeholder="Type a message... (Press Enter to send, Shift+Enter for new line)"
					rows="1"
					class="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none"
				></textarea>
			</div>
			
			<div class="flex items-center space-x-2">
				<IconButton 
					variant="ghost"
					size="sm"
					title="Attach file"
					class="text-gray-400 hover:text-gray-200"
				>
					<Paperclip class="btn-icon-sm" />
				</IconButton>
				
				<Button
					on:click={sendMessage}
					disabled={!newMessage.trim()}
					variant="primary"
					size="sm"
					class="flex items-center space-x-2"
				>
					<span>Send</span>
					<Send class="btn-icon-sm" />
				</Button>
			</div>
		</div>
	</div>
</div>

<style>
	.bg-gray-750 {
		background-color: #374151;
	}
	
	.btn-icon-xs {
		width: 12px;
		height: 12px;
	}
	
	.btn-icon-sm {
		width: 16px;
		height: 16px;
	}
</style>