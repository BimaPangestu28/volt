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
}

const initialState: WorkspaceState = {
	workspaces: [],
	currentWorkspace: null,
	collections: [],
	isLoading: false,
	error: null
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
		reset: () => {
			set(initialState);
		}
	};
}

export const workspaceStore = createWorkspaceStore();