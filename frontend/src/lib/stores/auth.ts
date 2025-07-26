import { writable } from 'svelte/store';
import type { Writable } from 'svelte/store';

export interface User {
	id: string;
	username: string;
	email: string;
	created_at: string;
	updated_at: string;
}

export interface AuthState {
	user: User | null;
	token: string | null;
	isAuthenticated: boolean;
	isLoading: boolean;
}

const initialState: AuthState = {
	user: null,
	token: null,
	isAuthenticated: false,
	isLoading: false
};

function createAuthStore() {
	const { subscribe, set, update }: Writable<AuthState> = writable(initialState);

	return {
		subscribe,
		setUser: (user: User) => {
			update(state => ({
				...state,
				user,
				token: null, // Token is now in httpOnly cookie
				isAuthenticated: true,
				isLoading: false
			}));
			// Store only user data in localStorage (token is in httpOnly cookie)
			if (typeof window !== 'undefined') {
				localStorage.setItem('volt_user', JSON.stringify(user));
			}
		},
		logout: () => {
			set(initialState);
			// Remove user data from localStorage (cookie will be cleared by server)
			if (typeof window !== 'undefined') {
				localStorage.removeItem('volt_user');
			}
		},
		setLoading: (loading: boolean) => {
			update(state => ({ ...state, isLoading: loading }));
		},
		initFromStorage: () => {
			if (typeof window !== 'undefined') {
				const userData = localStorage.getItem('volt_user');
				if (userData) {
					try {
						const user = JSON.parse(userData);
						update(state => ({
							...state,
							user,
							isAuthenticated: true
						}));
					} catch (error) {
						// Invalid stored data, clear it
						localStorage.removeItem('volt_user');
					}
				}
			}
		}
	};
}

export const authStore = createAuthStore();