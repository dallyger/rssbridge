package forgejo

type NotificationRepository struct {
	Archived      bool
	ArchivedAt    Time   `json:"archived_at"`
	AvatarUrl     string `json:"avatar_url"`
	CreatedAt     Time   `json:"created_at"`
	DefaultBranch string `json:"default_branch"`
	Description   string
	Fork          bool
	FullName      string `json:"full_name"`
	HtmlUrl       string `json:"html_url"`
	Id            uint
	Name          string
	Owner         User
	Private       bool
	Url           string
	Website       string
}

type NotificationSubject struct {
	Title                string
	Type                 string
	Url                  string
	HtmlUrl              string `json:"html_url"`
	LatestCommentUrl     string `json:"latest_comment_url"`
	LatestCommentHtmlUrl string `json:"latest_comment_html_url"`
	State                string
}

type Notification struct {
	Id              uint
	Pinned          bool
	Reason          string
	Repository      NotificationRepository
	Subject         NotificationSubject
	SubscriptionUrl string `json:"subscription_url"`
	Unread          bool
	UpdatedAt       Time `json:"updated_at"`
	Url             string
}
