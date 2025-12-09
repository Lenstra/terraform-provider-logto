
resource "logto_api_resource_scope" "api_resource_scope" {
  name        = "scope_name"
  resource_id = logto_api_resource.api_resource.id
  description = "scope_description"

  depends_on = [logto_api_resource.api_resource]
}
