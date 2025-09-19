
provider "logto" {
  hostname           = "yourHostname"         //In case of cloud hosted use `yourTenantId.logto.app`
  resource           = "yourResourceEndpoint" // Only for self-hosted else it use hostname automatically
  application_id     = "yourApplicationId"
  application_secret = "yourApplicationSecret"
}
