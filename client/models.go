package client

type LogtoDefaultStruct struct {
	TenantId    string `json:"tenantId"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OidcClientMetadata struct {
	RedirectUris                     []string `json:"redirectUris"`
	PostLogoutRedirectUris           []string `json:"postLogoutRedirectUris"`
	BackchannelLogoutUri             string   `json:"backchannelLogoutUri,omitempty"`
	BackchannelLogoutSessionRequired bool     `json:"backchannelLogoutSessionRequired,omitempty"`
	LogoUri                          string   `json:"logoUri,omitempty"`
}

type CustomClientMetadata struct {
	CorsAllowedOrigins      []string `json:"corsAllowedOrigins"`
	IdTokenTtl              float64  `json:"idTokenTtl,omitempty"`
	RefreshTokenTtl         float64  `json:"refreshTokenTtl,omitempty"`
	RefreshTokenTtlInDays   float64  `json:"refreshTokenTtlInDays,omitempty"`
	TenantId                string   `json:"tenantId,omitempty"`
	AlwaysIssueRefreshToken bool     `json:"alwaysIssueRefreshToken,omitempty"`
	RotateRefreshToken      bool     `json:"rotateRefreshToken,omitempty"`
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
	Type                 string                 `json:"type"`
	OidcClientMetadata   OidcClientMetadata     `json:"oidcClientMetadata"`
	CustomClientMetadata CustomClientMetadata   `json:"customClientMetadata"`
	CustomData           map[string]interface{} `json:"customData"`
	ProtectedAppMetadata ProtectedAppMetadata   `json:"protectedAppMetadata"`
	IsAdmin              bool                   `json:"isAdmin"`
	Secrets              map[string]string
}

type Secret struct {
	TenantId      string `json:"tenantId"`
	ApplicationId string `json:"applicationId"`
	Name          string `json:"name"`
	Value         string `json:"value"`
}
