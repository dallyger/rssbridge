package github

type Release struct {
	Author       User
	Body         string
	Created_at   Time
	Published_at Time
	Draft        bool
	Html_url     string
	Id           uint
	Name         string
	Tag_name     string
	Prerelease   bool
}
