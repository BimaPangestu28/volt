import { writable } from 'svelte/store';
import type { Writable } from 'svelte/store';

export interface Workspace {
	id: string;
	name: string;
	description: string;
	owner_id: string;
	members: string[];
	created_at: string;
	updated_at: string;
}

export interface Collection {
	id: string;
	name: string;
	description: string;
	workspace_id: string;
	created_by: string;
	requests: string[];
	created_at: string;
	updated_at: string;
}

export interface WorkspaceState {
	workspaces: Workspace[];
	currentWorkspace: Workspace | null;
	collections: Collection[];
	isLoading: boolean;
	error: string | null;
	collectionsCache: Record<string, {
		data: Collection[];
		timestamp: number;
		isRefreshing: boolean;
	}>;
}

const CACHE_DURATION = 5 * 60 * 1000; // 5 minutes

const initialState: WorkspaceState = {
	workspaces: [],
	currentWorkspace: null,
	collections: [],
	isLoading: false,
	error: null,
	collectionsCache: {}
};

function createWorkspaceStore() {
	const { subscribe, set, update }: Writable<WorkspaceState> = writable(initialState);

	return {
		subscribe,
		setWorkspaces: (workspaces: Workspace[]) => {
			update(state => ({ ...state, workspaces }));
		},
		setCurrentWorkspace: (workspace: Workspace | null) => {
			update(state => ({ ...state, currentWorkspace: workspace }));
		},
		setCollections: (collections: Collection[]) => {
			update(state => ({ ...state, collections }));
		},
		addWorkspace: (workspace: Workspace) => {
			update(state => ({
				...state,
				workspaces: [...state.workspaces, workspace]
			}));
		},
		addCollection: (collection: Collection) => {
			update(state => ({
				...state,
				collections: [...state.collections, collection]
			}));
		},
		setLoading: (loading: boolean) => {
			update(state => ({ ...state, isLoading: loading }));
		},
		setError: (error: string | null) => {
			update(state => ({ ...state, error }));
		},
		// Cache methods
		getCachedCollections: (workspaceId: string): Collection[] | null => {
			let cachedData: Collection[] | null = null;
			update(state => {
				const cache = state.collectionsCache[workspaceId];
				if (cache && Date.now() - cache.timestamp < CACHE_DURATION) {
					cachedData = cache.data;
				}
				return state;
			});
			return cachedData;
		},
		setCachedCollections: (workspaceId: string, collections: Collection[]) => {
			update(state => ({
				...state,
				collections,
				collectionsCache: {
					...state.collectionsCache,
					[workspaceId]: {
						data: collections,
						timestamp: Date.now(),
						isRefreshing: false
					}
				}
			}));
		},
		setRefreshing: (workspaceId: string, isRefreshing: boolean) => {
			update(state => {
				const cache = state.collectionsCache[workspaceId];
				if (cache) {
					return {
						...state,
						collectionsCache: {
							...state.collectionsCache,
							[workspaceId]: {
								...cache,
								isRefreshing
							}
						}
					};
				}
				return state;
			});
		},
		isCacheValid: (workspaceId: string): boolean => {
			let isValid = false;
			update(state => {
				const cache = state.collectionsCache[workspaceId];
				isValid = cache && Date.now() - cache.timestamp < CACHE_DURATION;
				return state;
			});
			return isValid;
		},
		reset: () => {
			set(initialState);
		}
	};
}

export const workspaceStore = createWorkspaceStore();