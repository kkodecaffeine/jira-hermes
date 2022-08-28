package issues

import "github.com/andygrunwald/go-jira"

// User is
type User struct {
	AccountID    string `json:"accountId,omitempty" structs:"accountId,omitempty"`
	Name         string `json:"name,omitempty" structs:"name,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty" structs:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty" structs:"displayName,omitempty"`
}

// Users are

type Users []*User

// Project is
type Project struct {
	ID   string `json:"id,omitempty" structs:"id,omitempty"`
	Key  string `json:"key,omitempty" structs:"key,omitempty"`
	Name string `json:"name,omitempty" structs:"name,omitempty"`
}

// Projects are
type Projects []*Project

// Issue is
type Issue struct {
	ProjectName string            `json:"projectname,omitempty" structs:"projectname,omitempty"`
	IssueKey    string            `json:"issuekey,omitempty" structs:"issuekey,omitempty"`
	IssueLink   string            `json:"issuelink,omitempty" structs:"issuelink,omitempty"`
	Summary     string            `json:"summary,omitempty" structs:"summary,omitempty"`
	Status      *jira.Status      `json:"status,omitempty" structs:"status,omitempty"`
	IssueType   jira.IssueType    `json:"issuetype,omitempty" structs:"issuetype,omitempty"`
	Assignee    *jira.User        `json:"assignee,omitempty" structs:"assignee,omitempty"`
	Owner       *jira.IssueFields `json:"onwer,omitempty" structs:"onwer,omitempty"`
	OwnerEmail  string            `json:"owneremail,omitempty" structs:"owneremail,omitempty"`
	Priority    *jira.Priority    `json:"priority,omitempty" structs:"priority,omitempty"`
	Labels      []string          `json:"labels,omitempty" structs:"labels,omitempty"`
	ElapsedDays int               `json:"elapseddate,omitempty" structs:"elapseddate,omitempty"`
}

// Issues are
type Issues []*Issue
