package common

type BugBountyProgram struct {
	InputURL    string `json:"input_url"`
	Platform    string `json:"platform"`
	Business    string `json:"business"`
	ProgramName string `json:"program"`
	PolicyURL   string `json:"policy_url"`
	FetchedAt   string `json:"fetched_at"`
}

type Result struct {
	ProgramDetails BugBountyProgram `json:"program"`
	InScope        []string         `json:"in_scope"`
	OutScope       []string         `json:"out_scope"`
	FetchedAt      string           `json:"fetched_at"`
}
