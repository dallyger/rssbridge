package github

type Issue struct {
	Body       string
	Created_at Time
	Diff_url   string
	Draft      bool
	Html_url   string
	Id         uint
	Number     uint
	Patch_url  string
	Sha        string
	State      string
	Title      string
	Updated_at Time
	User       User
}
