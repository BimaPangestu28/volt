import { writable } from 'svelte/store';

export interface ProjectTask {
    id: string;
    title: string;
    description?: string;
    status: 'todo' | 'in_progress' | 'review' | 'completed';
    priority: 'low' | 'medium' | 'high' | 'urgent';
    assignedTo?: string;
    createdAt: Date;
    updatedAt: Date;
    dueDate?: Date;
    tags: string[];
    projectId: string;
}

export interface ProjectMilestone {
    id: string;
    title: string;
    description?: string;
    dueDate: Date;
    status: 'upcoming' | 'in_progress' | 'completed' | 'overdue';
    progress: number; // 0-100
    projectId: string;
    tasks: string[]; // task IDs
}

export interface ProjectMember {
    id: string;
    name: string;
    email: string;
    avatar?: string;
    role: 'owner' | 'admin' | 'member' | 'viewer';
    status: 'online' | 'away' | 'offline';
    lastSeen: Date;
}

export interface Project {
    id: string;
    name: string;
    description?: string;
    type: 'api' | 'web' | 'mobile' | 'desktop' | 'data' | 'other';
    status: 'planning' | 'active' | 'on_hold' | 'completed' | 'archived';
    priority: 'low' | 'medium' | 'high' | 'urgent';
    progress: number; // 0-100
    startDate: Date;
    dueDate?: Date;
    createdAt: Date;
    updatedAt: Date;
    ownerId: string;
    members: ProjectMember[];
    tags: string[];
    milestones: ProjectMilestone[];
    tasks: ProjectTask[];
    settings: {
        isPublic: boolean;
        allowGuestAccess: boolean;
        autoSync: boolean;
        notifications: boolean;
    };
    stats: {
        totalTasks: number;
        completedTasks: number;
        overdueTasks: number;
        activeMembers: number;
        totalApiCalls: number;
        lastActivity: Date;
    };
}

export interface ProjectTemplate {
    id: string;
    name: string;
    description: string;
    type: Project['type'];
    defaultTasks: Omit<ProjectTask, 'id' | 'projectId' | 'createdAt' | 'updatedAt'>[];
    defaultMilestones: Omit<ProjectMilestone, 'id' | 'projectId' | 'tasks'>[];
    tags: string[];
    estimatedDuration: number; // in days
}

interface ProjectStoreState {
    projects: Project[];
    currentProject: Project | null;
    tasks: ProjectTask[];
    milestones: ProjectMilestone[];
    templates: ProjectTemplate[];
    isLoading: boolean;
    error: string | null;
    filters: {
        status: Project['status'] | 'all';
        type: Project['type'] | 'all';
        priority: Project['priority'] | 'all';
        search: string;
    };
    view: 'grid' | 'list' | 'kanban';
}

const initialState: ProjectStoreState = {
    projects: [],
    currentProject: null,
    tasks: [],
    milestones: [],
    templates: [],
    isLoading: false,
    error: null,
    filters: {
        status: 'all',
        type: 'all',
        priority: 'all',
        search: ''
    },
    view: 'grid'
};

function createProjectStore() {
    const { subscribe, set, update } = writable(initialState);

    return {
        subscribe,
        
        // Project actions
        setProjects: (projects: Project[]) => 
            update(state => ({ ...state, projects })),
            
        addProject: (project: Project) =>
            update(state => ({ 
                ...state, 
                projects: [...state.projects, project] 
            })),
            
        updateProject: (projectId: string, updates: Partial<Project>) =>
            update(state => ({
                ...state,
                projects: state.projects.map(p => 
                    p.id === projectId ? { ...p, ...updates, updatedAt: new Date() } : p
                ),
                currentProject: state.currentProject?.id === projectId 
                    ? { ...state.currentProject, ...updates, updatedAt: new Date() }
                    : state.currentProject
            })),
            
        deleteProject: (projectId: string) =>
            update(state => ({
                ...state,
                projects: state.projects.filter(p => p.id !== projectId),
                currentProject: state.currentProject?.id === projectId ? null : state.currentProject
            })),
            
        setCurrentProject: (project: Project | null) =>
            update(state => ({ ...state, currentProject: project })),
            
        // Task actions
        setTasks: (tasks: ProjectTask[]) =>
            update(state => ({ ...state, tasks })),
            
        addTask: (task: ProjectTask) =>
            update(state => ({ 
                ...state, 
                tasks: [...state.tasks, task],
                projects: state.projects.map(p => 
                    p.id === task.projectId 
                        ? { 
                            ...p, 
                            tasks: [...p.tasks, task],
                            stats: { 
                                ...p.stats, 
                                totalTasks: p.stats.totalTasks + 1,
                                lastActivity: new Date()
                            }
                        }
                        : p
                )
            })),
            
        updateTask: (taskId: string, updates: Partial<ProjectTask>) =>
            update(state => {
                const updatedTask = state.tasks.find(t => t.id === taskId);
                if (!updatedTask) return state;
                
                const newTask = { ...updatedTask, ...updates, updatedAt: new Date() };
                const statusChanged = updatedTask.status !== newTask.status;
                
                return {
                    ...state,
                    tasks: state.tasks.map(t => t.id === taskId ? newTask : t),
                    projects: state.projects.map(p => {
                        if (p.id === newTask.projectId) {
                            const updatedTasks = p.tasks.map(t => t.id === taskId ? newTask : t);
                            const completedTasks = updatedTasks.filter(t => t.status === 'completed').length;
                            const overdueTasks = updatedTasks.filter(t => 
                                t.dueDate && new Date(t.dueDate) < new Date() && t.status !== 'completed'
                            ).length;
                            
                            return {
                                ...p,
                                tasks: updatedTasks,
                                stats: {
                                    ...p.stats,
                                    completedTasks: statusChanged ? completedTasks : p.stats.completedTasks,
                                    overdueTasks: statusChanged ? overdueTasks : p.stats.overdueTasks,
                                    lastActivity: new Date()
                                }
                            };
                        }
                        return p;
                    })
                };
            }),
            
        deleteTask: (taskId: string) =>
            update(state => {
                const task = state.tasks.find(t => t.id === taskId);
                if (!task) return state;
                
                return {
                    ...state,
                    tasks: state.tasks.filter(t => t.id !== taskId),
                    projects: state.projects.map(p => 
                        p.id === task.projectId
                            ? {
                                ...p,
                                tasks: p.tasks.filter(t => t.id !== taskId),
                                stats: {
                                    ...p.stats,
                                    totalTasks: p.stats.totalTasks - 1,
                                    completedTasks: task.status === 'completed' 
                                        ? p.stats.completedTasks - 1 
                                        : p.stats.completedTasks,
                                    lastActivity: new Date()
                                }
                            }
                            : p
                    )
                };
            }),
            
        // Milestone actions
        setMilestones: (milestones: ProjectMilestone[]) =>
            update(state => ({ ...state, milestones })),
            
        addMilestone: (milestone: ProjectMilestone) =>
            update(state => ({ 
                ...state, 
                milestones: [...state.milestones, milestone],
                projects: state.projects.map(p => 
                    p.id === milestone.projectId
                        ? { ...p, milestones: [...p.milestones, milestone] }
                        : p
                )
            })),
            
        updateMilestone: (milestoneId: string, updates: Partial<ProjectMilestone>) =>
            update(state => ({
                ...state,
                milestones: state.milestones.map(m => 
                    m.id === milestoneId ? { ...m, ...updates } : m
                ),
                projects: state.projects.map(p => ({
                    ...p,
                    milestones: p.milestones.map(m => 
                        m.id === milestoneId ? { ...m, ...updates } : m
                    )
                }))
            })),
            
        // Filter actions
        setFilters: (filters: Partial<ProjectStoreState['filters']>) =>
            update(state => ({ 
                ...state, 
                filters: { ...state.filters, ...filters } 
            })),
            
        setView: (view: ProjectStoreState['view']) =>
            update(state => ({ ...state, view })),
            
        // Loading and error states
        setLoading: (isLoading: boolean) =>
            update(state => ({ ...state, isLoading })),
            
        setError: (error: string | null) =>
            update(state => ({ ...state, error })),
            
        // Templates
        setTemplates: (templates: ProjectTemplate[]) =>
            update(state => ({ ...state, templates })),
            
        // Utility functions
        getProjectById: (projectId: string): Project | undefined => {
            let project: Project | undefined;
            update(state => {
                project = state.projects.find(p => p.id === projectId);
                return state;
            });
            return project;
        },
        
        getTasksByProject: (projectId: string): ProjectTask[] => {
            let tasks: ProjectTask[] = [];
            update(state => {
                tasks = state.tasks.filter(t => t.projectId === projectId);
                return state;
            });
            return tasks;
        },
        
        getMilestonesByProject: (projectId: string): ProjectMilestone[] => {
            let milestones: ProjectMilestone[] = [];
            update(state => {
                milestones = state.milestones.filter(m => m.projectId === projectId);
                return state;
            });
            return milestones;
        }
    };
}

export const projectStore = createProjectStore();

// Default project templates
export const defaultTemplates: ProjectTemplate[] = [
    {
        id: 'api-project',
        name: 'API Development',
        description: 'RESTful API development with testing and documentation',
        type: 'api',
        defaultTasks: [
            {
                title: 'Design API Architecture',
                description: 'Define endpoints, data models, and authentication',
                status: 'todo',
                priority: 'high',
                tags: ['planning', 'architecture']
            },
            {
                title: 'Setup Development Environment',
                description: 'Configure database, server, and testing tools',
                status: 'todo',
                priority: 'high',
                tags: ['setup', 'devops']
            },
            {
                title: 'Implement Core Endpoints',
                description: 'Build CRUD operations for main resources',
                status: 'todo',
                priority: 'medium',
                tags: ['development', 'backend']
            },
            {
                title: 'Add Authentication & Authorization',
                description: 'Implement JWT-based authentication system',
                status: 'todo',
                priority: 'high',
                tags: ['security', 'auth']
            },
            {
                title: 'Write API Documentation',
                description: 'Create comprehensive API documentation',
                status: 'todo',
                priority: 'medium',
                tags: ['documentation']
            },
            {
                title: 'Setup API Testing',
                description: 'Create test collections and automated tests',
                status: 'todo',
                priority: 'medium',
                tags: ['testing', 'qa']
            }
        ],
        defaultMilestones: [
            {
                title: 'MVP API Complete',
                description: 'Core functionality implemented and tested',
                dueDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000), // 30 days
                status: 'upcoming',
                progress: 0
            },
            {
                title: 'Production Ready',
                description: 'Security, documentation, and deployment complete',
                dueDate: new Date(Date.now() + 60 * 24 * 60 * 60 * 1000), // 60 days
                status: 'upcoming',
                progress: 0
            }
        ],
        tags: ['api', 'backend', 'rest'],
        estimatedDuration: 45
    },
    {
        id: 'web-app',
        name: 'Web Application',
        description: 'Full-stack web application development',
        type: 'web',
        defaultTasks: [
            {
                title: 'Project Setup & Architecture',
                description: 'Setup frontend, backend, and database',
                status: 'todo',
                priority: 'high',
                tags: ['setup', 'architecture']
            },
            {
                title: 'Design System & UI Components',
                description: 'Create reusable UI components and design system',
                status: 'todo',
                priority: 'medium',
                tags: ['design', 'frontend']
            },
            {
                title: 'User Authentication',
                description: 'Implement user registration, login, and session management',
                status: 'todo',
                priority: 'high',
                tags: ['auth', 'security']
            },
            {
                title: 'Core Features Development',
                description: 'Build main application features',
                status: 'todo',
                priority: 'high',
                tags: ['development', 'features']
            },
            {
                title: 'Testing & Quality Assurance',
                description: 'Unit tests, integration tests, and manual testing',
                status: 'todo',
                priority: 'medium',
                tags: ['testing', 'qa']
            },
            {
                title: 'Deployment & DevOps',
                description: 'Setup CI/CD pipeline and production deployment',
                status: 'todo',
                priority: 'medium',
                tags: ['devops', 'deployment']
            }
        ],
        defaultMilestones: [
            {
                title: 'Alpha Version',
                description: 'Core features implemented for internal testing',
                dueDate: new Date(Date.now() + 45 * 24 * 60 * 60 * 1000), // 45 days
                status: 'upcoming',
                progress: 0
            },
            {
                title: 'Beta Release',
                description: 'Feature-complete version for user testing',
                dueDate: new Date(Date.now() + 75 * 24 * 60 * 60 * 1000), // 75 days
                status: 'upcoming',
                progress: 0
            },
            {
                title: 'Production Launch',
                description: 'Stable version ready for public release',
                dueDate: new Date(Date.now() + 90 * 24 * 60 * 60 * 1000), // 90 days
                status: 'upcoming',
                progress: 0
            }
        ],
        tags: ['web', 'fullstack', 'frontend', 'backend'],
        estimatedDuration: 90
    },
    {
        id: 'mobile-app',
        name: 'Mobile Application',
        description: 'Cross-platform mobile application development',
        type: 'mobile',
        defaultTasks: [
            {
                title: 'Mobile App Architecture',
                description: 'Define app structure, navigation, and state management',
                status: 'todo',
                priority: 'high',
                tags: ['architecture', 'planning']
            },
            {
                title: 'UI/UX Design Implementation',
                description: 'Create mobile-responsive design and components',
                status: 'todo',
                priority: 'high',
                tags: ['design', 'ui', 'mobile']
            },
            {
                title: 'API Integration',
                description: 'Connect mobile app to backend services',
                status: 'todo',
                priority: 'medium',
                tags: ['api', 'integration']
            },
            {
                title: 'Offline Functionality',
                description: 'Implement local storage and offline capabilities',
                status: 'todo',
                priority: 'medium',
                tags: ['offline', 'storage']
            },
            {
                title: 'Testing on Devices',
                description: 'Test on various mobile devices and screen sizes',
                status: 'todo',
                priority: 'medium',
                tags: ['testing', 'devices']
            },
            {
                title: 'App Store Deployment',
                description: 'Prepare and submit to app stores',
                status: 'todo',
                priority: 'low',
                tags: ['deployment', 'stores']
            }
        ],
        defaultMilestones: [
            {
                title: 'Prototype Ready',
                description: 'Basic functionality working on mobile devices',
                dueDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000), // 30 days
                status: 'upcoming',
                progress: 0
            },
            {
                title: 'Beta Testing',
                description: 'Feature-complete app ready for beta testing',
                dueDate: new Date(Date.now() + 60 * 24 * 60 * 60 * 1000), // 60 days
                status: 'upcoming',
                progress: 0
            },
            {
                title: 'App Store Release',
                description: 'App published and available for download',
                dueDate: new Date(Date.now() + 90 * 24 * 60 * 60 * 1000), // 90 days
                status: 'upcoming',
                progress: 0
            }
        ],
        tags: ['mobile', 'cross-platform', 'app'],
        estimatedDuration: 75
    }
];