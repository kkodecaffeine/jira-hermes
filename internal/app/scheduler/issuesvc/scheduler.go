package issuesvc

import (
	"fmt"
	"jira-hermes/internal/app/scheduler/common"
	"jira-hermes/internal/pkg/issues"
	"time"
)

type Scheduler struct {
	usecase issues.UseCase
}

// NewScheduler returns new scheduler instance
func NewScheduler(uc issues.UseCase) Scheduler {
	s := Scheduler{uc}

	s.PostOldIssues()
	s.PostDirectMessages()
	s.PostBigSellerIssues()

	return s
}

func (ctrl *Scheduler) PostOldIssues() {
	t := time.Now()
	formatted := fmt.Sprintf("*%d-%02d-%02d*", t.Year(), t.Month(), t.Day())

	channelID, timestamp := common.PostMessage(formatted, "")

	projects, err := ctrl.usecase.GetProjects()
	if err != nil {
		panic("🐛 GET projects " + err.Error())
	}

	var keys []string
	for _, project := range projects {
		keys = append(keys, project.Key)
	}
	var count int

	issues, err := ctrl.usecase.GetIssues(keys)
	if err != nil {
		panic("🐛 GET issues " + err.Error())
	}

	count += len(issues)

	divided := ctrl.usecase.SplitToChunks(issues)

	for index, issue := range divided {
		message := ctrl.usecase.BuildBlockKitSection(issue[index].ProjectName, issue)

		if len(issues) > 0 {
			common.ReplyMessage(channelID, timestamp, message)
		}
	}

	common.UpdateMessage(channelID, timestamp, fmt.Sprintf("%v `total: %v`", formatted, count))
}

func (ctrl *Scheduler) PostDirectMessages() {
	t := time.Now()
	formatted := fmt.Sprintf("*%d-%02d-%02d*", t.Year(), t.Month(), t.Day())

	projects, err := ctrl.usecase.GetProjects()
	if err != nil {
		panic("🐛 GET projects " + err.Error())
	}

	var keys []string
	for _, project := range projects {
		keys = append(keys, project.Key)
	}

	issues, err := ctrl.usecase.GetIssues(keys)
	if err != nil {
		panic("🐛 GET issues " + err.Error())
	}

	ownersissues := ctrl.usecase.GetOwnersIssues(issues)

	for email, issues := range ownersissues {
		openedID := common.OpenConversation(email)
		if openedID == "" {
			continue
		}

		channelID, timestamp := common.PostMessage(fmt.Sprintf("%v `total: %v`", formatted, len(issues)), openedID)
		divided := ctrl.usecase.SplitToChunks(issues)

		for index, issue := range divided {
			message := ctrl.usecase.BuildBlockKitSection(issue[index].ProjectName, issue)

			if len(issues) > 0 {
				common.ReplyMessage(channelID, timestamp, message)
			}
		}
	}
}

func (ctrl *Scheduler) PostBigSellerIssues() {
	projects, err := ctrl.usecase.GetProjects()
	if err != nil {
		panic("🐛 GET projects " + err.Error())
	}

	var keys []string
	for _, project := range projects {
		keys = append(keys, project.Key)
	}

	issues, err := ctrl.usecase.GetBigSellerIssues(keys)
	if err != nil {
		panic("🐛 GET issues " + err.Error())
	}

	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02d (total: %d)", t.Year(), t.Month(), t.Day(), len(issues))

	rawmatrix := ctrl.usecase.BuildMatrixByStatus(issues)
	message := ctrl.usecase.BuildJiraLinksWithBlockKitSection(rawmatrix)
	builtmatrix := ctrl.usecase.BuildDashboard(formatted, rawmatrix)

	channelID, timestamp := common.PostMessage(builtmatrix, "")
	common.ReplyMessage(channelID, timestamp, message)
}
