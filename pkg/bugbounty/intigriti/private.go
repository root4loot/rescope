package intigriti

type PrivateProgramList struct {
	MaxCount int `json:"maxCount"`
	Records  []struct {
		ID        string `json:"id"`
		Handle    string `json:"handle"`
		Name      string `json:"name"`
		Following bool   `json:"following"`
		MinBounty struct {
			Value    float32 `json:"value"`
			Currency string  `json:"currency"`
		} `json:"minBounty"`
		MaxBounty struct {
			Value    float32 `json:"value"`
			Currency string  `json:"currency"`
		} `json:"maxBounty"`
		ConfidentialityLevel struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"confidentialityLevel"`
		Status struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"status"`
		Type struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"type"`
		WebLinks struct {
			Detail string `json:"detail"`
		} `json:"webLinks"`
	} `json:"records"`
}

type PrivateProgramDetail struct {
	ID                   string `json:"id"`
	Handle               string `json:"handle"`
	Name                 string `json:"name"`
	Following            bool   `json:"following"`
	ConfidentialityLevel struct {
		ID    int    `json:"id"`
		Value string `json:"value"`
	} `json:"confidentialityLevel"`
	Status struct {
		ID    int    `json:"id"`
		Value string `json:"value"`
	} `json:"status"`
	Type struct {
		ID    int    `json:"id"`
		Value string `json:"value"`
	} `json:"type"`
	Domains struct {
		ID        string `json:"id"`
		CreatedAt int    `json:"createdAt"`
		Content   []struct {
			ID   string `json:"id"`
			Type struct {
				ID    int    `json:"id"`
				Value string `json:"value"`
			} `json:"type"`
			Endpoint string `json:"endpoint"`
			Tier     struct {
				ID    int    `json:"id"`
				Value string `json:"value"`
			} `json:"tier"`
			Description interface{} `json:"description"`
		} `json:"content"`
	} `json:"domains"`
	RulesOfEngagement struct {
		Attachments []interface{} `json:"attachments"`
		ID          string        `json:"id"`
		CreatedAt   int           `json:"createdAt"`
		Content     struct {
			Description         string `json:"description"`
			TestingRequirements struct {
				IntigritiMe      bool        `json:"intigritiMe"`
				AutomatedTooling int         `json:"automatedTooling"`
				UserAgent        interface{} `json:"userAgent"`
				RequestHeader    interface{} `json:"requestHeader"`
			} `json:"testingRequirements"`
			SafeHarbour bool `json:"safeHarbour"`
		} `json:"content"`
	} `json:"rulesOfEngagement"`
	WebLinks struct {
		Detail string `json:"detail"`
	} `json:"webLinks"`
}
