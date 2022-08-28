package scheduler

import (
	"os"

	jira "github.com/andygrunwald/go-jira"
)

type Gojira struct {
	Client *jira.Client
}

func getJira() *Gojira {
	base := os.Getenv("BASE_URL")
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	client, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		panic(err)
	}

	return &Gojira{
		Client: client,
	}
}
