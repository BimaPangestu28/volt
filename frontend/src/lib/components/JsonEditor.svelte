<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import loader from '@monaco-editor/loader';
	import type * as monaco from 'monaco-editor';

	export let value = '';
	export let readOnly = false;
	export let height = '300px';
	export let placeholder = '';

	let container: HTMLDivElement;
	let editor: monaco.editor.IStandaloneCodeEditor;
	let Monaco: typeof monaco;

	// Debounce value updates to avoid excessive reactivity
	let debounceTimer: number;
	function updateValue(newValue: string) {
		clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => {
			value = newValue;
		}, 300) as any;
	}

	onMount(async () => {
		try {
			// Configure Monaco loader
			loader.config({
				paths: {
					vs: 'https://cdn.jsdelivr.net/npm/monaco-editor@0.45.0/min/vs'
				}
			});

			// Load Monaco Editor
			Monaco = await loader.init();

			// Configure JSON language options
			Monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
				validate: true,
				allowComments: false,
				schemas: [],
				enableSchemaRequest: false
			});

			// Create editor
			editor = Monaco.editor.create(container, {
				value: value || placeholder,
				language: 'json',
				theme: 'vs-dark',
				readOnly,
				automaticLayout: true,
				minimap: { enabled: false },
				scrollBeyondLastLine: false,
				fontSize: 12,
				fontFamily: 'JetBrains Mono, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
				lineNumbers: 'on',
				renderWhitespace: 'selection',
				bracketPairColorization: { enabled: true },
				formatOnPaste: true,
				formatOnType: true,
				tabSize: 2,
				insertSpaces: true,
				wordWrap: 'on',
				folding: true,
				foldingStrategy: 'indentation',
				showFoldingControls: 'always',
				contextmenu: true,
				selectOnLineNumbers: true,
				roundedSelection: false,
				cursorStyle: 'line',
				cursorBlinking: 'blink',
				smoothScrolling: true,
				mouseWheelZoom: true,
				suggest: {
					showKeywords: true,
					showSnippets: true
				}
			});

			// Listen for content changes
			editor.onDidChangeModelContent(() => {
				const currentValue = editor.getValue();
				updateValue(currentValue);
			});

			// Format JSON on Ctrl+Shift+F
			editor.addCommand(Monaco.KeyMod.CtrlCmd | Monaco.KeyMod.Shift | Monaco.KeyCode.KeyF, () => {
				editor.getAction('editor.action.formatDocument')?.run();
			});

		} catch (error) {
			console.error('Failed to initialize Monaco Editor:', error);
		}
	});

	onDestroy(() => {
		if (editor) {
			editor.dispose();
		}
		clearTimeout(debounceTimer);
	});

	// Update editor value when prop changes
	$: if (editor && value !== editor.getValue()) {
		const position = editor.getPosition();
		editor.setValue(value || '');
		if (position) {
			editor.setPosition(position);
		}
	}

	// Utility functions to expose editor capabilities
	export function formatJSON() {
		if (editor) {
			editor.getAction('editor.action.formatDocument')?.run();
		}
	}

	export function validateJSON(): boolean {
		if (!editor) return false;
		
		try {
			const content = editor.getValue().trim();
			if (!content) return true;
			JSON.parse(content);
			return true;
		} catch {
			return false;
		}
	}

	export function getErrors(): string[] {
		if (!editor || !Monaco) return [];
		
		const model = editor.getModel();
		if (!model) return [];
		
		const markers = Monaco.editor.getModelMarkers({ resource: model.uri });
		return markers.map(marker => marker.message);
	}

	export function insertSnippet(snippet: string) {
		if (editor) {
			const position = editor.getPosition();
			if (position) {
				editor.executeEdits('', [{
					range: new Monaco.Range(position.lineNumber, position.column, position.lineNumber, position.column),
					text: snippet
				}]);
			}
		}
	}

	// Common JSON snippets
	export function insertObjectSnippet() {
		insertSnippet('{\n  "": ""\n}');
	}

	export function insertArraySnippet() {
		insertSnippet('[\n  \n]');
	}
</script>

<div class="json-editor-container" style="height: {height};">
	<div bind:this={container} class="w-full h-full" />
</div>

<style>
	.json-editor-container {
		border: 1px solid #374151;
		border-radius: 0.375rem;
		overflow: hidden;
		background: #1f2937;
	}

	:global(.monaco-editor) {
		background: #1f2937 !important;
	}

	:global(.monaco-editor .margin) {
		background: #1f2937 !important;
	}

	:global(.monaco-editor-background) {
		background: #1f2937 !important;
	}

	:global(.monaco-editor .monaco-editor-background) {
		background: #1f2937 !important;
	}
</style>