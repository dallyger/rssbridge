package forgejo

type IssueRepo struct {
	FullName string `json:"full_name"`
	Id       uint
	Name     string
	Owner    string
}

type Issue struct {
	Body       string
	CreatedAt  Time   `json:"created_at"`
	HtmlUrl    string `json:"html_url"`
	Id         uint
	Number     uint
	Repository IssueRepo
	State      string
	Title      string
	UpdatedAt  Time `json:"updated_at"`
	Url        string
	User       User
}
