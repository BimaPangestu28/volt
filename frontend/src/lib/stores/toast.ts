import { writable } from 'svelte/store';

export interface Toast {
	id: string;
	type: 'success' | 'error' | 'info';
	message: string;
	duration?: number;
}

function createToastStore() {
	const { subscribe, update } = writable<Toast[]>([]);

	return {
		subscribe,
		add: (toast: Omit<Toast, 'id'>) => {
			const id = Math.random().toString(36).substr(2, 9);
			update(toasts => [...toasts, { ...toast, id }]);
		},
		remove: (id: string) => {
			update(toasts => toasts.filter(t => t.id !== id));
		},
		success: (message: string, duration = 3000) => {
			const id = Math.random().toString(36).substr(2, 9);
			update(toasts => [...toasts, { id, type: 'success', message, duration }]);
		},
		error: (message: string, duration = 5000) => {
			const id = Math.random().toString(36).substr(2, 9);
			update(toasts => [...toasts, { id, type: 'error', message, duration }]);
		},
		info: (message: string, duration = 3000) => {
			const id = Math.random().toString(36).substr(2, 9);
			update(toasts => [...toasts, { id, type: 'info', message, duration }]);
		}
	};
}

export const toastStore = createToastStore();