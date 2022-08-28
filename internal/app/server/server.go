package server

import (
	scheduler "jira-hermes/internal/app/scheduler"
)

// NewServer Return new server instance
func NewServer() {
	scheduler.CreateSchedulerApp()
}
