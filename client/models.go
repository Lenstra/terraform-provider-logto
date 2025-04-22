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
	ID           string   `json:"id,omitempty"`
	PrimaryEmail string   `json:"primaryEmail,omitempty"`
	Username     string   `json:"username,omitempty"`
	Name         string   `json:"name,omitempty"`
	Profile      *Profile `json:"profile,omitempty"`
}

type Profile struct {
	FamilyName string `json:"familyName,omitempty"`
	GivenName  string `json:"givenName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
	Nickname   string `json:"nickname,omitempty"`
}

type Secret struct {
	TenantId      string `json:"tenantId"`
	ApplicationId string `json:"applicationId"`
	Name          string `json:"name"`
	Value         string `json:"value"`
}

type ApiResourceModel struct {
	TenantId       string        `json:"tenantId,omitempty"`
	ID             string        `json:"id,omitempty"`
	Name           string        `json:"name"`
	Indicator      string        `json:"indicator,omitempty"`
	AccessTokenTtl *float64      `json:"accessTokenTtl,omitempty"`
	IsDefault      *bool         `json:"isDefault,omitempty"`
	Scopes         *[]ScopeModel `json:"scopes,omitempty"`
}

type ScopeModel struct {
	TenantId    string           `json:"tenantId,omitempty"`
	ID          string           `json:"id,omitempty"`
	ResourceId  string           `json:"resourceId,omitempty"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	CreatedAt   *float64         `json:"createdAt,omitempty"`
	Resource    ApiResourceModel `json:"resource,omitempty"`
}

type RoleModel struct {
	TenantId    string   `json:"tenantId,omitempty"`
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type,omitempty"`
	IsDefault   bool     `json:"isDefault,omitempty"`
	ScopeIds    []string `json:"scopeIds,omitempty"`
}
