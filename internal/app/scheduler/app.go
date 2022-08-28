package scheduler

import (
	issuesvc "jira-hermes/internal/app/scheduler/issuesvc"
	"jira-hermes/internal/pkg/issues"
	issuetRepo "jira-hermes/internal/pkg/issues/repo"

	"github.com/gin-gonic/gin"
)

type App interface {
	Init()
	RegisterRoute(driver *gin.Engine)
	Clean() error
}

type apiApp struct {
	gojira *Gojira
}

func (app *apiApp) Init() {
	app.gojira = getJira()
}

func (app *apiApp) RegisterScheduler() {
	ic := issues.NewUseCase(issuetRepo.New(app.gojira.Client))
	issuesvc.NewScheduler(ic)
}

func (app *apiApp) Clean() error {
	return nil
}

// CreateSchedulerApp returns new core.App implementation
func CreateSchedulerApp() {
	hermesApp := &apiApp{}
	hermesApp.Init()
	hermesApp.RegisterScheduler()
	hermesApp.Clean()
}
