package github

type NotificationRepository struct {
	Description string
	Fork        bool
	Full_name   string
	Id          float32
	Name        string
	Private     bool
	HtmlUrl     string
}

type NotificationSubject struct {
	Title string
	Type  string
	Url   string
}

type Notification struct {
	Id               string
	Last_read_at     Time
	Reason           string
	Subject          NotificationSubject
	Unread           bool
	Updated_at       Time
	Repository       NotificationRepository
	Url              string
	Subscription_url string
}

