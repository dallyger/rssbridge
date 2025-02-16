package forgejo

type Comment struct {
	Body           string
	CreatedAt      Time   `json:"created_at"`
	HtmlUrl        string `json:"html_url"`
	Id             uint
	IssueUrl       string `json:"issue_url"`
	PullRequestUrl string `json:"pull_request_url"`
	UpdatedAt      Time   `json:"updated_at"`
	User           User
}
