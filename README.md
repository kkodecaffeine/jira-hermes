# jira-hermes
`jira-hermes` is a slack bot :robot_face: for tracking jira issues
This bot is responsible for three parts:
- Find old jira issues: find old issues by some criteria using jira apis
- Send slack messages to a specific channel and direct messages to project owners
- Build a dashboard: create a dashboard by some criteria
## Libraries used
- `gin-gonic`: https://github.com/gin-gonic/gin
- `go-jira`: https://github.com/andygrunwald/go-jira
- `lambda`: https://github.com/aws/aws-lambda-go/lambda
- `slack-go`: https://github.com/slack-go/slack
- `markdown`: https://github.com/fbiville/markdown-table-formatter/pkg/markdown
## Build and deploy
Using the aws-sam-cli tool (aws/serverless-application-model: AWS Serverless Application Model (SAM) is an open-source framework for building serverless applications (github.com))
## Build
`sam build`
(Use template.yml in project root directory for setting information required for build)
## Deploy
`sam deploy --config-env <env-name>`