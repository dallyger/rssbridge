package github

type Discussion struct {
	Body string
	Comments uint
	Created_at string
	Html_url string
	Id uint
	Locked bool
	Number uint
	State string
	State_reason string
	Title string
	Updated_at Time
	User User
}
