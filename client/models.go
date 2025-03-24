package client

type OidcClientMetadata struct {
	RedirectUris                     []string `json:"redirectUris"`
	PostLogoutRedirectUris           []string `json:"postLogoutRedirectUris"`
	BackchannelLogoutUri             string   `json:"backchannelLogoutUri,omitempty"`
	BackchannelLogoutSessionRequired bool     `json:"backchannelLogoutSessionRequired,omitempty"`
	LogoUri                          string   `json:"logoUri,omitempty"`
}

type CustomClientMetadata struct {
	CorsAllowedOrigins      []string `json:"corsAllowedOrigins,omitempty"`
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
	TenantId             string                 `json:"tenantId,omitempty"`
	ID                   string                 `json:"id,omitempty"`
	Name                 string                 `json:"name"`
	Description          string                 `json:"description,omitempty"`
	Type                 string                 `json:"type"`
	OidcClientMetadata   *OidcClientMetadata    `json:"oidcClientMetadata,omitempty"`
	CustomClientMetadata *CustomClientMetadata  `json:"customClientMetadata,omitempty"`
	CustomData           map[string]interface{} `json:"customData,omitempty"`
	ProtectedAppMetadata *ProtectedAppMetadata  `json:"protectedAppMetadata,omitempty"`
	IsAdmin              bool                   `json:"isAdmin"`
	IsThirdParty         bool                   `json:"isThirdParty"`
}

type UserModel struct {
	ID       string   `json:"id,omitempty"`
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Profile  *Profile `json:"profile"`
}

type Profile struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
	MiddleName string `json:"middleName"`
	Nickname   string `json:"nickname"`
}

type Secret struct {
	TenantId      string `json:"tenantId"`
	ApplicationId string `json:"applicationId"`
	Name          string `json:"name"`
	Value         string `json:"value"`
}
