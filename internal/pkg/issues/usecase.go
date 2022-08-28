package issues

import (
	"fmt"
	"jira-hermes/internal/app/scheduler/common"
	"regexp"
	"sort"
	"strings"

	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
	"github.com/slack-go/slack"
)

// UseCase interface definition
type UseCase interface {
	GetProjects() (Projects, error)
	GetIssues(keys []string) (Issues, error)
	GetBigSellerIssues(keys []string) (Issues, error)
	GetOwnersIssues(issues Issues) map[string]Issues
	SplitToChunks(issues Issues) []Issues
	BuildBlockKitSection(projectName string, issues Issues) slack.Message
	BuildDashboard(title string, matrix [][]string) string
	BuildJiraLinksWithBlockKitSection(matrix [][]string) slack.Message
	BuildMatrixByStatus(issues Issues) [][]string
	BuildMatrixByPeriod(issues Issues) [][]string
}

type usecase struct {
	repo Repository
}

func (u *usecase) GetProjects() (Projects, error) {
	return u.repo.FindProjects()
}

func (u *usecase) GetIssues(keys []string) (Issues, error) {
	return u.repo.FindIssues(keys)
}

func (u *usecase) GetBigSellerIssues(keys []string) (Issues, error) {
	return u.repo.FindBigSellerIssues(keys)
}

func (u *usecase) GetOwnersIssues(issues Issues) map[string]Issues {
	ownersissues := make(map[string]Issues)
	for _, element := range issues {
		if _, ok := ownersissues[element.OwnerEmail]; ok {
			ownersissues[element.OwnerEmail] = append(ownersissues[element.OwnerEmail], element)
		} else {
			newmap := make(map[string]Issues)
			newmap[element.OwnerEmail] = append(newmap[element.OwnerEmail], element)
			ownersissues[element.OwnerEmail] = newmap[element.OwnerEmail]
		}
	}

	return ownersissues
}

func (u *usecase) SplitToChunks(issues Issues) []Issues {
	var divided []Issues
	chunksize := 50

	for i := 0; i < len(issues); i += chunksize {
		end := i + chunksize

		if end > len(issues) {
			end = len(issues)
		}

		divided = append(divided, issues[i:end])
	}

	return divided
}

func (u *usecase) BuildBlockKitSection(projectName string, issues Issues) slack.Message {
	var message slack.Message

	for _, element := range issues {

		fixedSummary := strings.Replace(element.Summary, ">", "", 1)
		fixedSummary = strings.Replace(fixedSummary, "<", "", 1)

		aaa := "*<@issueLink|@summary>*```Labels: @labels\nDays elapsed: @elapsedDays days \nStatus: @status\nAssignee: @assignee\nOwner: @owner\nPriority: @priority```"
		aaa = strings.Replace(aaa, "@issueLink", element.IssueLink, 1)
		aaa = strings.Replace(aaa, "@summary", fixedSummary, 1)
		aaa = strings.Replace(aaa, "@status", element.Status.Name, 1)
		aaa = strings.Replace(aaa, "@issueType", element.IssueType.Name, 1)

		if element.Assignee != nil {
			aaa = strings.Replace(aaa, "@assignee", element.Assignee.DisplayName, 1)
		} else {
			aaa = strings.Replace(aaa, "@assignee", "Unassigned", 1)
		}
		unknowns := element.Owner.Unknowns["customfield_10037"]
		owner, exist := common.Find(unknowns, "displayName")

		if exist {
			aaa = strings.Replace(aaa, "@owner", fmt.Sprintf("%v", owner), 1)
		}

		aaa = strings.Replace(aaa, "@labels", strings.Join(element.Labels, ","), 1)
		aaa = strings.Replace(aaa, "@priority", element.Priority.Name, 1)
		aaa = strings.Replace(aaa, "@elapsedDays", fmt.Sprintf("%v", element.ElapsedDays), 1)

		aTextBlock := slack.NewTextBlockObject("mrkdwn", aaa, false, false)
		aSectionBlock := slack.NewSectionBlock(aTextBlock, nil, nil)
		temp := slack.AddBlockMessage(message, aSectionBlock)
		message = temp
	}

	return message
}

func (u *usecase) BuildDashboard(title string, matrix [][]string) string {
	for i := 0; i < len(matrix); i++ {
		nums := matrix[i]
		for c, v := range nums {
			if c == 0 {
				continue
			}
			if v != "" {
				count := strings.Count(v, ",") + 1
				matrix[i][c] = fmt.Sprintf("%v", count)
			}
		}
	}

	table, err := markdown.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build("", "BUG", "STORY", "TASK", "SUB-TASK").
		Format(matrix)

	if err != nil {
		panic("🐛 Build dashboard " + err.Error())
	}

	return "```" + title + "\n\n" + table + "```"
}

func (u *usecase) BuildMatrixByStatus(issues Issues) [][]string {
	matrix := [][]string{}

	for _, element := range issues {
		r, c := common.SearchMatrix(matrix, element.Status.Name)

		match := regexp.MustCompile(`[(].*[)]`)
		jiraurl := match.ReplaceAllString(common.JIRAURL, "("+element.IssueKey+")")

		if r == -1 && c == -1 {
			var temp []string
			if element.IssueType.Name == "Bug" {
				temp = []string{element.Status.Name, jiraurl, "", "", ""}
			} else if element.IssueType.Name == "Story" {
				temp = []string{element.Status.Name, "", jiraurl, "", ""}
			} else if element.IssueType.Name == "Task" {
				temp = []string{element.Status.Name, "", "", jiraurl, ""}
			} else if element.IssueType.Name == "Sub-task" {
				temp = []string{element.Status.Name, "", "", "", jiraurl}
			} else {
				fmt.Println(element.IssueType.Name)
			}

			if temp == nil {
				continue
			}

			matrix = append(matrix, temp)
		} else {
			col := 1
			if element.IssueType.Name == "Bug" {
				col = 1
			} else if element.IssueType.Name == "Story" {
				col = 2
			} else if element.IssueType.Name == "Task" {
				col = 3
			} else if element.IssueType.Name == "Sub-task" {
				col = 4
			}

			if matrix[r][col] == "" {
				matrix[r][col] = jiraurl
			} else {
				var tobereplaced string

				found := match.FindAllString(matrix[r][col], 1)[0]
				found = strings.ReplaceAll(found, "(", "") // TODO
				found = strings.ReplaceAll(found, ")", "") // TODO
				tobereplaced = found + "," + element.IssueKey

				matrix[r][col] = match.ReplaceAllString(matrix[r][col], "("+tobereplaced+")")
			}
		}
	}

	return matrix
}

func (u *usecase) BuildMatrixByPeriod(issues Issues) [][]string {
	matrix := [][]string{}
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].ElapsedDays < issues[j].ElapsedDays
	})

	for _, element := range issues {
		var period string
		if element.ElapsedDays >= 7 && element.ElapsedDays < 14 {
			period = "1W-2W"
		} else if element.ElapsedDays >= 14 && element.ElapsedDays < 21 {
			period = "2W-3W"
		} else if element.ElapsedDays >= 21 && element.ElapsedDays < 28 {
			period = "3W-4W"
		} else {
			period = "4W"
		}
		r, c := common.SearchMatrix(matrix, period)

		match := regexp.MustCompile(`[(].*[)]`)
		jiraurl := match.ReplaceAllString(common.JIRAURL, "("+element.IssueKey+")")

		if r == -1 && c == -1 {
			var temp []string
			if element.IssueType.Name == "Bug" {
				temp = []string{period, jiraurl, "", "", ""}
			} else if element.IssueType.Name == "Story" {
				temp = []string{period, "", jiraurl, "", ""}
			} else if element.IssueType.Name == "Task" {
				temp = []string{period, "", "", jiraurl, ""}
			} else if element.IssueType.Name == "Sub-task" {
				temp = []string{period, "", "", "", jiraurl}
			}

			if temp == nil {
				continue
			}

			matrix = append(matrix, temp)
		} else {
			col := 1
			if element.IssueType.Name == "Bug" {
				col = 1
			} else if element.IssueType.Name == "Story" {
				col = 2
			} else if element.IssueType.Name == "Task" {
				col = 3
			} else if element.IssueType.Name == "Sub-task" {
				col = 4
			}

			if matrix[r][col] == "" {
				matrix[r][col] = jiraurl
			} else {
				var tobereplaced string

				found := match.FindAllString(matrix[r][col], 1)[0]
				found = strings.ReplaceAll(found, "(", "") // TODO
				found = strings.ReplaceAll(found, ")", "") // TODO
				tobereplaced = found + "," + element.IssueKey

				matrix[r][col] = match.ReplaceAllString(matrix[r][col], "("+tobereplaced+")")
			}
		}
	}

	return matrix
}

func (u *usecase) BuildJiraLinksWithBlockKitSection(matrix [][]string) slack.Message {
	var message slack.Message

	for i := 0; i < len(matrix); i++ {
		nums := matrix[i]
		for col, v := range nums {
			if col == 0 {
				continue
			}
			if v != "" {
				var typename string
				if col == 1 {
					typename = "Bug"
				} else if col == 2 {
					typename = "Story"
				} else if col == 3 {
					typename = "Task"
				} else if col == 4 {
					typename = "Sub-task"
				}

				count := strings.Count(v, ",") + 1
				aTextBlock := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*<%v|%v, %v, %v>*", v, matrix[i][0], typename, count), false, false)
				aSectionBlock := slack.NewSectionBlock(aTextBlock, nil, nil)
				temp := slack.AddBlockMessage(message, aSectionBlock)
				message = temp
			}
		}
	}

	return message
}

// NewUseCase returns new UseCase implementation
func NewUseCase(issueRepo Repository) UseCase {
	return &usecase{repo: issueRepo}
}

var _ UseCase = &usecase{}
