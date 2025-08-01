package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectMember represents a team member in a project
type ProjectMember struct {
	ID       string    `json:"id" bson:"id"`
	Name     string    `json:"name" bson:"name"`
	Email    string    `json:"email" bson:"email"`
	Avatar   string    `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Role     string    `json:"role" bson:"role"` // owner, admin, member, viewer
	Status   string    `json:"status" bson:"status"` // online, away, offline
	LastSeen time.Time `json:"last_seen" bson:"last_seen"`
}

// ProjectTask represents a task within a project
type ProjectTask struct {
	ID          string     `json:"id" bson:"id"`
	Title       string     `json:"title" bson:"title" binding:"required"`
	Description string     `json:"description,omitempty" bson:"description,omitempty"`
	Status      string     `json:"status" bson:"status"` // todo, in_progress, review, completed
	Priority    string     `json:"priority" bson:"priority"` // low, medium, high, urgent
	AssignedTo  string     `json:"assigned_to,omitempty" bson:"assigned_to,omitempty"`
	CreatedAt   time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" bson:"updated_at"`
	DueDate     *time.Time `json:"due_date,omitempty" bson:"due_date,omitempty"`
	Tags        []string   `json:"tags" bson:"tags"`
	ProjectID   string     `json:"project_id" bson:"project_id"`
}

// ProjectMilestone represents a milestone within a project
type ProjectMilestone struct {
	ID          string    `json:"id" bson:"id"`
	Title       string    `json:"title" bson:"title" binding:"required"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	DueDate     time.Time `json:"due_date" bson:"due_date" binding:"required"`
	Status      string    `json:"status" bson:"status"` // upcoming, in_progress, completed, overdue
	Progress    int       `json:"progress" bson:"progress"` // 0-100
	ProjectID   string    `json:"project_id" bson:"project_id"`
	Tasks       []string  `json:"tasks" bson:"tasks"` // task IDs associated with this milestone
}

// ProjectSettings represents project configuration
type ProjectSettings struct {
	IsPublic         bool `json:"is_public" bson:"is_public"`
	AllowGuestAccess bool `json:"allow_guest_access" bson:"allow_guest_access"`
	AutoSync         bool `json:"auto_sync" bson:"auto_sync"`
	Notifications    bool `json:"notifications" bson:"notifications"`
}

// ProjectStats represents project statistics
type ProjectStats struct {
	TotalTasks     int64     `json:"total_tasks" bson:"total_tasks"`
	CompletedTasks int64     `json:"completed_tasks" bson:"completed_tasks"`
	OverdueTasks   int64     `json:"overdue_tasks" bson:"overdue_tasks"`
	ActiveMembers  int64     `json:"active_members" bson:"active_members"`
	TotalAPICalls  int64     `json:"total_api_calls" bson:"total_api_calls"`
	LastActivity   time.Time `json:"last_activity" bson:"last_activity"`
}

// Project represents a development project
type Project struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name        string              `json:"name" bson:"name" binding:"required"`
	Description string              `json:"description,omitempty" bson:"description,omitempty"`
	Type        string              `json:"type" bson:"type" binding:"required"` // api, web, mobile, desktop, data, other
	Status      string              `json:"status" bson:"status"` // planning, active, on_hold, completed, archived
	Priority    string              `json:"priority" bson:"priority"` // low, medium, high, urgent
	Progress    int                 `json:"progress" bson:"progress"` // 0-100
	StartDate   time.Time           `json:"start_date" bson:"start_date"`
	DueDate     *time.Time          `json:"due_date,omitempty" bson:"due_date,omitempty"`
	CreatedAt   time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" bson:"updated_at"`
	OwnerID     string              `json:"owner_id" bson:"owner_id"`
	Members     []ProjectMember     `json:"members" bson:"members"`
	Tags        []string            `json:"tags" bson:"tags"`
	Milestones  []ProjectMilestone  `json:"milestones" bson:"milestones"`
	Tasks       []ProjectTask       `json:"tasks" bson:"tasks"`
	Settings    ProjectSettings     `json:"settings" bson:"settings"`
	Stats       ProjectStats        `json:"stats" bson:"stats"`
}

// ProjectTemplate represents a project template
type ProjectTemplate struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name               string             `json:"name" bson:"name" binding:"required"`
	Description        string             `json:"description" bson:"description"`
	Type               string             `json:"type" bson:"type"` // api, web, mobile, desktop, data, other
	DefaultTasks       []ProjectTask      `json:"default_tasks" bson:"default_tasks"`
	DefaultMilestones  []ProjectMilestone `json:"default_milestones" bson:"default_milestones"`
	Tags               []string           `json:"tags" bson:"tags"`
	EstimatedDuration  int                `json:"estimated_duration" bson:"estimated_duration"` // in days
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedBy          string             `json:"created_by" bson:"created_by"`
	IsPublic           bool               `json:"is_public" bson:"is_public"`
}

// ProjectActivity represents activity log entries
type ProjectActivity struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProjectID   string             `json:"project_id" bson:"project_id"`
	UserID      string             `json:"user_id" bson:"user_id"`
	UserName    string             `json:"user_name" bson:"user_name"`
	Action      string             `json:"action" bson:"action"` // created, updated, deleted, completed, etc.
	EntityType  string             `json:"entity_type" bson:"entity_type"` // project, task, milestone, member
	EntityID    string             `json:"entity_id,omitempty" bson:"entity_id,omitempty"`
	EntityName  string             `json:"entity_name,omitempty" bson:"entity_name,omitempty"`
	Description string             `json:"description" bson:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
}

// ProjectFile represents files associated with a project
type ProjectFile struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProjectID     string             `json:"project_id" bson:"project_id"`
	Name          string             `json:"name" bson:"name" binding:"required"`
	OriginalName  string             `json:"original_name" bson:"original_name"`
	Type          string             `json:"type" bson:"type"` // image, document, code, archive, other
	Size          int64              `json:"size" bson:"size"`
	MimeType      string             `json:"mime_type" bson:"mime_type"`
	URL           string             `json:"url" bson:"url"`
	Description   string             `json:"description,omitempty" bson:"description,omitempty"`
	Tags          []string           `json:"tags" bson:"tags"`
	UploadedBy    string             `json:"uploaded_by" bson:"uploaded_by"`
	UploadedAt    time.Time          `json:"uploaded_at" bson:"uploaded_at"`
	IsShared      bool               `json:"is_shared" bson:"is_shared"`
	DownloadCount int64              `json:"download_count" bson:"download_count"`
	Version       string             `json:"version" bson:"version"`
	Folder        string             `json:"folder,omitempty" bson:"folder,omitempty"`
}

// ProjectComment represents comments on projects, tasks, or milestones
type ProjectComment struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProjectID   string             `json:"project_id" bson:"project_id"`
	EntityType  string             `json:"entity_type" bson:"entity_type"` // project, task, milestone
	EntityID    string             `json:"entity_id" bson:"entity_id"`
	UserID      string             `json:"user_id" bson:"user_id"`
	UserName    string             `json:"user_name" bson:"user_name"`
	Content     string             `json:"content" bson:"content" binding:"required"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	IsEdited    bool               `json:"is_edited" bson:"is_edited"`
	Mentions    []string           `json:"mentions,omitempty" bson:"mentions,omitempty"` // user IDs mentioned
	Attachments []string           `json:"attachments,omitempty" bson:"attachments,omitempty"` // file IDs
}

// ProjectTimeEntry represents time tracking entries
type ProjectTimeEntry struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProjectID   string             `json:"project_id" bson:"project_id"`
	TaskID      string             `json:"task_id,omitempty" bson:"task_id,omitempty"`
	UserID      string             `json:"user_id" bson:"user_id"`
	UserName    string             `json:"user_name" bson:"user_name"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	StartTime   time.Time          `json:"start_time" bson:"start_time"`
	EndTime     *time.Time         `json:"end_time,omitempty" bson:"end_time,omitempty"`
	Duration    int64              `json:"duration" bson:"duration"` // in seconds
	IsRunning   bool               `json:"is_running" bson:"is_running"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// ProjectAnalytics represents analytics data for a project
type ProjectAnalytics struct {
	ProjectID         string                 `json:"project_id" bson:"project_id"`
	DateRange         DateRange              `json:"date_range" bson:"date_range"`
	TaskCompletion    TaskCompletionMetrics  `json:"task_completion" bson:"task_completion"`
	TeamActivity      []TeamActivityMetrics  `json:"team_activity" bson:"team_activity"`
	TimeTracking      []TimeTrackingMetrics  `json:"time_tracking" bson:"time_tracking"`
	MilestoneProgress []MilestoneMetrics     `json:"milestone_progress" bson:"milestone_progress"`
	VelocityMetrics   []VelocityMetrics      `json:"velocity_metrics" bson:"velocity_metrics"`
	Priority          PriorityDistribution   `json:"priority_distribution" bson:"priority_distribution"`
	GeneratedAt       time.Time              `json:"generated_at" bson:"generated_at"`
}

type DateRange struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}

type TaskCompletionMetrics struct {
	Completed  int `json:"completed" bson:"completed"`
	InProgress int `json:"in_progress" bson:"in_progress"`
	Todo       int `json:"todo" bson:"todo"`
	Review     int `json:"review" bson:"review"`
}

type TeamActivityMetrics struct {
	UserID         string  `json:"user_id" bson:"user_id"`
	UserName       string  `json:"user_name" bson:"user_name"`
	TasksCompleted int     `json:"tasks_completed" bson:"tasks_completed"`
	HoursLogged    float64 `json:"hours_logged" bson:"hours_logged"`
	LastActive     time.Time `json:"last_active" bson:"last_active"`
}

type TimeTrackingMetrics struct {
	Date  string  `json:"date" bson:"date"`
	Hours float64 `json:"hours" bson:"hours"`
}

type MilestoneMetrics struct {
	ID       string    `json:"id" bson:"id"`
	Name     string    `json:"name" bson:"name"`
	Progress int       `json:"progress" bson:"progress"`
	DueDate  time.Time `json:"due_date" bson:"due_date"`
	Status   string    `json:"status" bson:"status"`
}

type VelocityMetrics struct {
	Week           string  `json:"week" bson:"week"`
	TasksCompleted int     `json:"tasks_completed" bson:"tasks_completed"`
	AverageTime    float64 `json:"average_time" bson:"average_time"` // hours
}

type PriorityDistribution struct {
	Urgent int `json:"urgent" bson:"urgent"`
	High   int `json:"high" bson:"high"`
	Medium int `json:"medium" bson:"medium"`
	Low    int `json:"low" bson:"low"`
}