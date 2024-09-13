package intigriti

type PublicProgramDetail struct {
	ProgramID            string `json:"programId"`
	Status               int    `json:"status"`
	ConfidentialityLevel int    `json:"confidentialityLevel"`
	CompanyHandle        string `json:"companyHandle"`
	CompanyName          string `json:"companyName"`
	Handle               string `json:"handle"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Domains              []struct {
		Content []struct {
			ID           string      `json:"id"`
			Type         int         `json:"type"`
			Endpoint     string      `json:"endpoint"`
			BountyTierID int         `json:"bountyTierId"`
			Description  interface{} `json:"description"`
		} `json:"content"`
		CreatedAt int `json:"createdAt"`
	} `json:"domains"`
	InScopes []struct {
		Content struct {
			Content     string        `json:"content"`
			Attachments []interface{} `json:"attachments"`
		} `json:"content"`
		CreatedAt int `json:"createdAt"`
	} `json:"inScopes"`
	OutOfScopes []struct {
		Content struct {
			Content     string        `json:"content"`
			Attachments []interface{} `json:"attachments"`
		} `json:"content"`
		CreatedAt int `json:"createdAt"`
	} `json:"outOfScopes"`
	Faqs []struct {
		Content struct {
			Content     string        `json:"content"`
			Attachments []interface{} `json:"attachments"`
		} `json:"content"`
		CreatedAt int `json:"createdAt"`
	} `json:"faqs"`
	SeverityAssessments []struct {
		Content struct {
			Content     string        `json:"content"`
			Attachments []interface{} `json:"attachments"`
		} `json:"content"`
		CreatedAt int `json:"createdAt"`
	} `json:"severityAssessments"`
	RulesOfEngagements []struct {
		Content struct {
			Content struct {
				Description         string `json:"description"`
				TestingRequirements struct {
					IntigritiMe      bool        `json:"intigritiMe"`
					AutomatedTooling int         `json:"automatedTooling"`
					UserAgent        interface{} `json:"userAgent"`
					RequestHeader    interface{} `json:"requestHeader"`
				} `json:"testingRequirements"`
				SafeHarbour bool `json:"safeHarbour"`
				CreatedAt   int  `json:"createdAt"`
			} `json:"content"`
			Attachments []interface{} `json:"attachments"`
		} `json:"content"`
		CreatedAt int `json:"createdAt"`
	} `json:"rulesOfEngagements"`
	BountyTables []struct {
		Content struct {
			Currency   string `json:"currency"`
			BountyRows []struct {
				BountyRanges []struct {
					MinScore  float64 `json:"minScore"`
					MaxScore  float64 `json:"maxScore"`
					MinBounty struct {
						Value    float64 `json:"value"`
						Currency string  `json:"currency"`
					} `json:"minBounty"`
					MaxBounty struct {
						Value    float64 `json:"value"`
						Currency string  `json:"currency"`
					} `json:"maxBounty"`
				} `json:"bountyRanges"`
				BountyTierID int `json:"bountyTierId"`
			} `json:"bountyRows"`
			RewardPolicy interface{} `json:"rewardPolicy"`
		} `json:"content"`
		CreatedAt int `json:"createdAt"`
	} `json:"bountyTables"`
	LastContributors []struct {
		Role            string `json:"role"`
		IdentityChecked bool   `json:"identityChecked"`
		UserID          string `json:"userId"`
		AvatarID        string `json:"avatarId"`
		UserName        string `json:"userName"`
	} `json:"lastContributors"`
	LastActivity []struct {
		CompanyName string `json:"companyName,omitempty"`
		LogoID      string `json:"logoId,omitempty"`
		Timestamp   int    `json:"timestamp"`
		Type        int    `json:"type"`
		Researcher  struct {
			Role            string      `json:"role"`
			IdentityChecked bool        `json:"identityChecked"`
			UserID          string      `json:"userId"`
			AvatarID        interface{} `json:"avatarId"`
			UserName        string      `json:"userName"`
		} `json:"researcher,omitempty"`
	} `json:"lastActivity"`
	AveragePayout           interface{} `json:"averagePayout"`
	SubmissionCount         int         `json:"submissionCount"`
	AcceptedSubmissionCount int         `json:"acceptedSubmissionCount"`
	TotalPayout             interface{} `json:"totalPayout"`
	IdentityCheckedRequired bool        `json:"identityCheckedRequired"`
	AwardRep                bool        `json:"awardRep"`
	SkipTriage              bool        `json:"skipTriage"`
	LogoID                  string      `json:"logoId"`
	HasUpdates              bool        `json:"hasUpdates"`
	AllowCollaboration      bool        `json:"allowCollaboration"`
}
