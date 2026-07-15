package response

import (
	"time"
)

type ReportResponse struct {
	ID          uint      `json:"id"`
	ReporterID  uint      `json:"reporter_id"`
	Resource    string    `json:"resource"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"`
	Severity    string    `json:"severity"`
	ActionTaken string    `json:"action_taken,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type ModStatsResponse struct {
	PendingReports int64 `json:"pending_reports"`
	ActiveBans     int   `json:"active_bans"`
	ActiveFilters  int   `json:"active_filters"`
}
