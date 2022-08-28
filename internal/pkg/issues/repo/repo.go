package repo

import (
	"fmt"
	"jira-hermes/internal/app/scheduler/common"
	"jira-hermes/internal/pkg/issues"
	"strings"

	"github.com/andygrunwald/go-jira"
)

type issuesRepo struct {
	client *jira.Client
	mapper entityMapper
}

var _ issues.Repository = &issuesRepo{}

func (r *issuesRepo) FindUser(ID string) (*issues.User, error) {
	client := r.client
	result, _, err := client.User.GetByAccountID(ID)
	if err != nil {
		panic("GET user" + err.Error())
	} else if result == nil {
		fmt.Println("Expected user. User is nil")
	}

	user := r.mapper.toUserEntity(result)

	return user, err
}

func (r *issuesRepo) FindProjects() (issues.Projects, error) {
	client := r.client
	request, _ := client.NewRawRequest("GET", "rest/api/3/project", nil)
	reponse := new([]jira.Project)

	_, err := client.Do(request, reponse)
	if err != nil {
		panic("GET Projects" + err.Error())
	}

	projects := make(issues.Projects, 0)
	for _, element := range *reponse {
		if element.Key == "TTP" || element.Key == "TSE" {
			projects = append(projects, r.mapper.toProjectEntity(element))
		}
	}

	return projects, nil
}

func (r *issuesRepo) FindIssues(keys []string) (issues.Issues, error) {
	client := r.client

	last := 0
	var jiraissues []jira.Issue

	for {
		options := &jira.SearchOptions{
			MaxResults: 1000, // Max results can go up to 1000
			StartAt:    last,
		}

		justString := strings.Join(keys, "','")

		chunk, response, err := client.Issue.Search(fmt.Sprintf(`project IN ('%s') AND status NOT IN (Developed, Developing, "Done DEV Testing", "On DEV", "On QA/STAGING", "Ready to Release", Released, "Selected for Development", Done, Closed, Archived) AND type = Bug AND (labels IN (big_seller_request) OR (updated <= -8d AND status CHANGED))`, justString), options)
		if err != nil {
			return nil, err
		}

		total := response.Total
		if jiraissues == nil {
			jiraissues = make([]jira.Issue, 0, total)
		}
		jiraissues = append(jiraissues, chunk...)
		last = response.StartAt + len(chunk)
		if last >= total {
			result := make(issues.Issues, 0)
			for _, element := range jiraissues {

				unknowns := element.Fields.Unknowns["customfield_10037"]
				owner, exist := common.Find(unknowns, "accountId")

				if exist {
					owners, _ := r.FindUser(fmt.Sprintf("%v", owner))
					email := owners.EmailAddress

					result = append(result, r.mapper.toIssueEntity(element, email))
				}
			}

			return result, nil
		}
	}
}

func (r *issuesRepo) FindBigSellerIssues(keys []string) (issues.Issues, error) {
	client := r.client

	last := 0
	var jiraissues []jira.Issue

	for {
		options := &jira.SearchOptions{
			MaxResults: 1000, // Max results can go up to 1000
			StartAt:    last,
		}

		justString := strings.Join(keys, "','")

		chunk, response, err := client.Issue.Search(fmt.Sprintf(`project IN ('%s') AND status NOT IN (Developed, Developing, "Done DEV Testing", "On DEV", "On QA/STAGING", "Ready to Release", Released, "Selected for Development", Done, Closed, Archived) AND labels IN (big_seller_request) AND updated <= -8d AND status CHANGED ORDER BY status ASC`, justString), options)
		if err != nil {
			return nil, err
		}

		total := response.Total
		if jiraissues == nil {
			jiraissues = make([]jira.Issue, 0, total)
		}
		jiraissues = append(jiraissues, chunk...)
		last = response.StartAt + len(chunk)
		if last >= total {
			result := make(issues.Issues, 0)
			for _, element := range jiraissues {

				unknowns := element.Fields.Unknowns["customfield_10037"]
				owner, exist := common.Find(unknowns, "accountId")

				if exist {
					owners, _ := r.FindUser(fmt.Sprintf("%v", owner))
					email := owners.EmailAddress

					result = append(result, r.mapper.toIssueEntity(element, email))
				}
			}

			return result, nil
		}
	}
}

func New(client *jira.Client) issues.Repository {
	return &issuesRepo{client, entityMapper{}}
}
