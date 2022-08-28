package repo

import (
	"fmt"
	"jira-hermes/internal/pkg/issues"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/andygrunwald/go-jira"
)

type entityMapper struct{}

func (e entityMapper) toUserEntity(user *jira.User) *issues.User {
	return &issues.User{
		AccountID:    user.AccountID,
		Name:         user.Name,
		EmailAddress: user.EmailAddress,
		DisplayName:  user.DisplayName,
	}
}

func (e entityMapper) toProjectEntity(project jira.Project) *issues.Project {
	return &issues.Project{
		ID:   project.ID,
		Key:  project.Key,
		Name: project.Name,
	}
}

func (e entityMapper) toIssueEntity(issue jira.Issue, owneremail string) *issues.Issue {
	r, _ := regexp.Compile("([A-Z])+")

	elapsed := time.Since(time.Time(issue.Fields.Updated))
	timeDuration := elapsed / (24 * time.Hour)
	days, _ := strconv.Atoi(fmt.Sprintf("%d", timeDuration))

	return &issues.Issue{
		ProjectName: r.FindString(issue.Key),
		IssueKey:    issue.Key,
		IssueLink:   os.Getenv("BASE_URL") + "/browse/" + issue.Key,
		Summary:     fmt.Sprintf("%v %v", issue.Key, issue.Fields.Summary),
		Status:      issue.Fields.Status,
		IssueType:   issue.Fields.Type,
		Assignee:    issue.Fields.Assignee,
		Owner:       issue.Fields,
		OwnerEmail:  owneremail,
		Priority:    issue.Fields.Priority,
		Labels:      issue.Fields.Labels,
		ElapsedDays: days,
	}
}
