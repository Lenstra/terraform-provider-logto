
resource "logto_api_resource" "api_resource" {
  name             = "api_resource_name"
  indicator        = "https://api-resource.test"
  access_token_ttl = 3600
}

resource "logto_api_resource_scope" "api_resource_scope" {
  name        = "scope_name"
  resource_id = logto_api_resource.api_resource.id
  description = "test_scope_description"

  depends_on = [logto_api_resource.api_resource]
}

resource "logto_role" "test_role" {
  name        = "role_name"
  description = "role_description"
  type        = "User"
  scope_ids = [
    logto_api_resource_scope.api_resource_scope.id
  ]

  depends_on = [
    logto_api_resource_scope.api_resource_scope
  ]
}
