package client

type LogtoDefaultStruct struct {
	TenantId    string `json:"tenantId"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ApplicationModel struct {
	LogtoDefaultStruct
	Type string `json:"type"`
}
