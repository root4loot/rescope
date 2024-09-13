package config

type BurpConfig struct {
	Target struct {
		Scope struct {
			AdvancedMode bool          `json:"advancedMode"`
			Exclude      []BurpExclude `json:"exclude"`
			Include      []BurpInclude `json:"include"`
		} `json:"scope"`
	} `json:"target"`
}

type BurpExclude struct {
	Enabled  bool   `json:"enabled"`
	File     string `json:"file"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

type BurpInclude struct {
	Enabled  bool   `json:"enabled"`
	File     string `json:"file"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}
