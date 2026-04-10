package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	ScheduledAt *time.Time        `json:"scheduled_at,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	dto := taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	if !task.ScheduledAt.IsZero() {
		dto.ScheduledAt = &task.ScheduledAt
	}
	return dto
}

// generateTasksRequest is the request body for POST /api/v1/tasks/generate.
type generateTasksRequest struct {
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Status      taskdomain.Status         `json:"status"`
	StartDate   string                    `json:"start_date"` // YYYY-MM-DD
	EndDate     string                    `json:"end_date"`   // YYYY-MM-DD
	Recurrence  recurrenceRuleDTO         `json:"recurrence"`
}

type recurrenceRuleDTO struct {
	// Type is one of: daily, monthly, specific_dates, even_odd
	Type       taskdomain.RecurrenceType `json:"type"`
	// Interval is used with type=daily (every N days, >= 1)
	Interval   int                       `json:"interval,omitempty"`
	// DayOfMonth is used with type=monthly (1-30)
	DayOfMonth int                       `json:"day_of_month,omitempty"`
	// Dates is used with type=specific_dates (list of YYYY-MM-DD strings)
	Dates      []string                  `json:"dates,omitempty"`
	// EvenDays is used with type=even_odd (true = even days, false = odd days)
	EvenDays   bool                      `json:"even_days,omitempty"`
}
