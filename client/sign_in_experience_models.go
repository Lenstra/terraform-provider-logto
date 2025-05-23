package client

type SignInExperienceModel struct {
	TenantId           string       `json:"tenantId"`
	Id                 string       `json:"id"`
	Color              Color        `json:"color"`
	Branding           Branding     `json:"branding"`
	LanguageInfo       LanguageInfo `json:"languageInfo"`
	TermsOfUseUrl      string       `json:"termsOfUseUrl"`
	PrivacyPolicyUrl   string       `json:"privacyPolicyUrl"`
	AgreeToTermsPolicy string       `json:"agreeToTermsPolicy"`
	SignIn             SignIn       `json:"signIn"`
	SignUp             SignUp       `json:"signUp"`

	SocialSignIn struct {
		AutomaticAccountLinking bool `json:"automaticAccountLinking"`
	} `json:"socialSignIn"`

	SocialSignInConnectorTargets []string          `json:"socialSignInConnectorTargets"`
	SignInMode                   string            `json:"signInMode"`
	CustomCss                    string            `json:"customCss"`
	CustomContent                map[string]string `json:"customContent"`
	PasswordPolicy               PasswordPolicy    `json:"passwordPolicy"`
	Mfa                          Mfa               `json:"mfa"`
	SingleSignOnEnabled          bool              `json:"singleSignOnEnabled"`
	SupportEmail                 string            `json:"supportEmail"`
	SupportWebsiteUrl            string            `json:"supportWebsiteUrl"`
	UnknownSessionRedirectUrl    string            `json:"unknownSessionRedirectUrl"`

	CaptchaPolicy struct {
		Enabled bool `json:"enabled"`
	} `json:"captchaPolicy"`

	SentinelPolicy SentinelPolicy `json:"sentinelPolicy"`

	EmailBlocklistPolicy EmailBlocklistPolicy `json:"emailBlocklistPolicy"`
}

type Color struct {
	PrimaryColor      string `json:"primaryColor"`
	IsDarkModeEnabled string `json:"isDarkModeEnabled"`
	DarkPrimaryColor  string `json:"darkPrimaryColor"`
}

type Branding struct {
	LogoUrl     string `json:"logoUrl"`
	DarkLogoUrl string `json:"darkLogoUrl"`
	Favicon     string `json:"favicon"`
	DarkFavicon string `json:"darkFavicon"`
}

type LanguageInfo struct {
	AutoDetect       bool   `json:"autoDetect"`
	FallbackLanguage string `json:"fallbackLanguage"`
}

type SignIn struct {
	Methods struct {
		Identifier        string `json:"identifier"`
		Password          bool   `json:"password"`
		VerificationCode  bool   `json:"verificationCode"`
		IsPasswordPrimary bool   `json:"isPasswordPrimary"`
	} `json:"methods"`
}

type SignUp struct {
	Identifiers          []string `json:"identifiers"`
	Password             bool     `json:"password"`
	Verify               bool     `json:"verify"`
	SecondaryIdentifiers []any    `json:"secondaryIdentifiers"`
}

type PasswordPolicy struct {
	Length struct {
		Min int64 `json:"min"`
		Max int64 `json:"max"`
	} `json:"length"`

	CharacterTypes struct {
		Min int64 `json:"min"`
	} `json:"characterTypes"`

	Rejects Rejects `json:"rejects"`
}

type Rejects struct {
	Pwned                 bool     `json:"pwned"`
	RepetitionAndSequence bool     `json:"repetitionAndSequence"`
	UserInfo              bool     `json:"userInfo"`
	Words                 []string `json:"words"`
}

type SentinelPolicy struct {
	MaxAttempts     float64 `json:"maxAttempts"`
	LockoutDuration float64 `json:"lockoutDuration"`
}

type Mfa struct {
	Factors                       []string `json:"factors"`
	Policy                        string   `json:"policy"`
	OrganizationRequiredMfaPolicy string   `json:"organizationRequiredMfaPolicy"`
}

type EmailBlocklistPolicy struct {
	BlockDisposableAddresses bool `json:"blockDisposableAddresses"`
	BlockSubaddressing       bool `json:"blockSubaddressing"`
	CustomBlocklist          bool `json:"customBlocklist"`
}
