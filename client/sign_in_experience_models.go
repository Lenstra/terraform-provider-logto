package client

type SignInExperienceModel struct {
	TenantId                     string                `json:"tenantId,omitempty"`
	ID                           string                `json:"id,omitempty"`
	Color                        *Color                `json:"color,omitempty"`
	Branding                     *Branding             `json:"branding,omitempty"`
	LanguageInfo                 *LanguageInfo         `json:"languageInfo,omitempty"`
	TermsOfUseUrl                string                `json:"termsOfUseUrl,omitempty"`
	PrivacyPolicyUrl             string                `json:"privacyPolicyUrl,omitempty"`
	AgreeToTermsPolicy           string                `json:"agreeToTermsPolicy,omitempty"`
	SignIn                       *SignIn               `json:"signIn,omitempty"`
	SignUp                       *SignUp               `json:"signUp,omitempty"`
	SocialSignIn                 *SocialSignIn         `json:"socialSignIn,omitempty"`
	SocialSignInConnectorTargets []string              `json:"socialSignInConnectorTargets,omitempty"`
	SignInMode                   string                `json:"signInMode,omitempty"`
	CustomCss                    string                `json:"customCss,omitempty"`
	CustomContent                map[string]string     `json:"customContent,omitempty"`
	CustomUiAssets               *CustomUiAssets       `json:"customUiAssets,omitempty"`
	PasswordPolicy               *PasswordPolicy       `json:"passwordPolicy,omitempty"`
	Mfa                          *Mfa                  `json:"mfa,omitempty"`
	SingleSignOnEnabled          bool                  `json:"singleSignOnEnabled,omitempty"`
	SupportEmail                 string                `json:"supportEmail,omitempty"`
	SupportWebsiteUrl            string                `json:"supportWebsiteUrl,omitempty"`
	UnknownSessionRedirectUrl    string                `json:"unknownSessionRedirectUrl,omitempty"`
	CaptchaPolicy                *CaptchaPolicy        `json:"captchaPolicy,omitempty"`
	SentinelPolicy               *SentinelPolicy       `json:"sentinelPolicy,omitempty"`
	EmailBlocklistPolicy         *EmailBlocklistPolicy `json:"emailBlocklistPolicy,omitempty"`
}

type CaptchaPolicy struct {
	Enabled bool `json:"enabled"`
}

type SocialSignIn struct {
	AutomaticAccountLinking bool `json:"automaticAccountLinking,omitempty"`
}

type Color struct {
	PrimaryColor      string `json:"primaryColor"`
	IsDarkModeEnabled bool   `json:"isDarkModeEnabled"`
	DarkPrimaryColor  string `json:"darkPrimaryColor"`
}

type Branding struct {
	LogoUrl     string `json:"logoUrl,omitempty"`
	DarkLogoUrl string `json:"darkLogoUrl,omitempty"`
	Favicon     string `json:"favicon,omitempty"`
	DarkFavicon string `json:"darkFavicon,omitempty"`
}

type CustomUiAssets struct {
	ID        string  `json:"id"`
	CreatedAt float64 `json:"createdAt"`
}

type LanguageInfo struct {
	AutoDetect       bool   `json:"autoDetect"`
	FallbackLanguage string `json:"fallbackLanguage"`
}

type SignIn struct {
	Methods []Methods `json:"methods"`
}

type Methods struct {
	Identifier        string `json:"identifier"`
	Password          bool   `json:"password"`
	VerificationCode  bool   `json:"verificationCode"`
	IsPasswordPrimary bool   `json:"isPasswordPrimary"`
}

type SignUp struct {
	Identifiers          []string                `json:"identifiers"`
	Password             bool                    `json:"password"`
	Verify               bool                    `json:"verify"`
	SecondaryIdentifiers *[]SecondaryIdentifiers `json:"secondaryIdentifiers,omitempty"`
}

type SecondaryIdentifiers struct {
	Identifier string `json:"identifier"`
	Verify     bool   `json:"verify,omitempty"`
}

type PasswordPolicy struct {
	Length         *Length         `json:"length,omitempty"`
	CharacterTypes *CharacterTypes `json:"characterTypes,omitempty"`
	Rejects        *Rejects        `json:"rejects,omitempty"`
}

type Length struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

type CharacterTypes struct {
	Min int64 `json:"min"`
}

type Rejects struct {
	Pwned                 bool     `json:"pwned"`
	RepetitionAndSequence bool     `json:"repetitionAndSequence"`
	UserInfo              bool     `json:"userInfo"`
	Words                 []string `json:"words"`
}

type SentinelPolicy struct {
	MaxAttempts     float64 `json:"maxAttempts,omitempty"`
	LockoutDuration float64 `json:"lockoutDuration,omitempty"`
}

type Mfa struct {
	Factors                       []string `json:"factors"`
	Policy                        string   `json:"policy"`
	OrganizationRequiredMfaPolicy string   `json:"organizationRequiredMfaPolicy,omitempty"`
}

type EmailBlocklistPolicy struct {
	BlockDisposableAddresses bool     `json:"blockDisposableAddresses,omitempty"`
	BlockSubaddressing       bool     `json:"blockSubaddressing,omitempty"`
	CustomBlocklist          []string `json:"customBlocklist,omitempty"`
}
