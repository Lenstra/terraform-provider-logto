package client

type LogtoDefaultStruct struct {
	TenantId    string `json:"tenantId"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OidcClientMetadata struct {
	RedirectUris                     []interface{} `json:"redirectUris"`
	PostLogoutRedirectUris           []string      `json:"postLogoutRedirectUris"`
	BackchannelLogoutUri             string        `json:"backchannelLogoutUri"`
	BackchannelLogoutSessionRequired bool          `json:"backchannelLogoutSessionRequired"`
	LogoUri                          string        `json:"logoUri"`
}

type CustomClientMetadata struct {
	CorsAllowedOrigins      []string `json:"corsAllowedOrigins"`
	IdTokenTtl              float64  `json:"idTokenTtl"`
	RefreshTokenTtl         float64  `json:"refreshTokenTtl"`
	RefreshTokenTtlInDays   float64  `json:"refreshTokenTtlInDays"`
	TenantId                string   `json:"tenantId"`
	AlwaysIssueRefreshToken bool     `json:"alwaysIssueRefreshToken"`
	RotateRefreshToken      bool     `json:"rotateRefreshToken"`
}

type PageRule struct {
	Path string `json:"path"`
}

type ProtectedAppMetadata struct {
	Origin          string     `json:"origin"`
	SessionDuration float64    `json:"sessionDuration"`
	PageRules       []PageRule `json:"pageRules"`
}

type ApplicationModel struct {
	LogtoDefaultStruct
	Type                 string                 `json:"type,omitempty"`
	OidcClientMetadata   OidcClientMetadata     `json:"oidcClientMetadata"`
	CustomClientMetadata CustomClientMetadata   `json:"customClientMetadata"`
	CustomData           map[string]interface{} `json:"customData"`
	ProtectedAppMetadata ProtectedAppMetadata   `json:"protectedAppMetadata"`
	IsAdmin              bool                   `json:"isAdmin"`
}
