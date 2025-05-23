package client

type SignInExperienceModel struct {
	TenantId           string       `json:"tenantId,omitempty"`
	Id                 string       `json:"id,omitempty"`
	Color              Color        `json:"color,omitempty"`
	Branding           Branding     `json:"branding,omitempty"`
	LanguageInfo       LanguageInfo `json:"languageInfo,omitempty"`
	TermsOfUseUrl      string       `json:"termsOfUseUrl,omitempty"`
	PrivacyPolicyUrl   string       `json:"privacyPolicyUrl,omitempty"`
	AgreeToTermsPolicy string       `json:"agreeToTermsPolicy,omitempty"`
	SignIn             SignIn       `json:"signIn,omitempty"`
	SignUp             SignUp       `json:"signUp,omitempty"`

	SocialSignIn struct {
		AutomaticAccountLinking bool `json:"automaticAccountLinking,omitempty"`
	} `json:"socialSignIn,omitempty"`

	SocialSignInConnectorTargets []string          `json:"socialSignInConnectorTargets,omitempty"`
	SignInMode                   string            `json:"signInMode,omitempty"`
	CustomCss                    string            `json:"customCss,omitempty"`
	CustomContent                map[string]string `json:"customContent,omitempty"`
	PasswordPolicy               PasswordPolicy    `json:"passwordPolicy,omitempty"`
	Mfa                          Mfa               `json:"mfa,omitempty"`
	SingleSignOnEnabled          bool              `json:"singleSignOnEnabled,omitempty"`
	SupportEmail                 string            `json:"supportEmail,omitempty"`
	SupportWebsiteUrl            string            `json:"supportWebsiteUrl,omitempty"`
	UnknownSessionRedirectUrl    string            `json:"unknownSessionRedirectUrl,omitempty"`

	CaptchaPolicy struct {
		Enabled bool `json:"enabled,omitempty"`
	} `json:"captchaPolicy,omitempty"`

	SentinelPolicy SentinelPolicy `json:"sentinelPolicy,omitempty"`

	EmailBlocklistPolicy EmailBlocklistPolicy `json:"emailBlocklistPolicy,omitempty"`
}

type Color struct {
	PrimaryColor      string `json:"primaryColor,omitempty"`
	IsDarkModeEnabled string `json:"isDarkModeEnabled,omitempty"`
	DarkPrimaryColor  string `json:"darkPrimaryColor,omitempty"`
}

type Branding struct {
	LogoUrl     string `json:"logoUrl,omitempty"`
	DarkLogoUrl string `json:"darkLogoUrl,omitempty"`
	Favicon     string `json:"favicon,omitempty"`
	DarkFavicon string `json:"darkFavicon,omitempty"`
}

type LanguageInfo struct {
	AutoDetect       bool   `json:"autoDetect,omitempty"`
	FallbackLanguage string `json:"fallbackLanguage,omitempty"`
}

type SignIn struct {
	Methods struct {
		Identifier        string `json:"identifier,omitempty"`
		Password          bool   `json:"password,omitempty"`
		VerificationCode  bool   `json:"verificationCode,omitempty"`
		IsPasswordPrimary bool   `json:"isPasswordPrimary,omitempty"`
	} `json:"methods,omitempty"`
}

type SignUp struct {
	Identifiers          []string `json:"identifiers,omitempty"`
	Password             bool     `json:"password,omitempty"`
	Verify               bool     `json:"verify,omitempty"`
	SecondaryIdentifiers []any    `json:"secondaryIdentifiers,omitempty"`
}

type PasswordPolicy struct {
	Length struct {
		Min int64 `json:"min,omitempty"`
		Max int64 `json:"max,omitempty"`
	} `json:"length,omitempty"`

	CharacterTypes struct {
		Min int64 `json:"min,omitempty"`
	} `json:"characterTypes,omitempty"`

	Rejects Rejects `json:"rejects,omitempty"`
}

type Rejects struct {
	Pwned                 bool     `json:"pwned,omitempty"`
	RepetitionAndSequence bool     `json:"repetitionAndSequence,omitempty"`
	UserInfo              bool     `json:"userInfo,omitempty"`
	Words                 []string `json:"words,omitempty"`
}

type SentinelPolicy struct {
	MaxAttempts     float64 `json:"maxAttempts,omitempty"`
	LockoutDuration float64 `json:"lockoutDuration,omitempty"`
}

type Mfa struct {
	Factors                       []string `json:"factors,omitempty"`
	Policy                        string   `json:"policy,omitempty"`
	OrganizationRequiredMfaPolicy string   `json:"organizationRequiredMfaPolicy,omitempty"`
}

type EmailBlocklistPolicy struct {
	BlockDisposableAddresses bool `json:"blockDisposableAddresses,omitempty"`
	BlockSubaddressing       bool `json:"blockSubaddressing,omitempty"`
	CustomBlocklist          bool `json:"customBlocklist,omitempty"`
}
