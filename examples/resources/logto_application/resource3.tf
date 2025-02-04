
resource "logto_application" "app" {
  name                      = "test"
  description               = "test app description"
  type                      = "SPA"
  redirect_uris 						= ["http://test_modified.test.fr", "http://test_modified.test.com"]
  post_logout_redirect_uris = ["http://redirect_modified.test.fr", "http://redirect_modified.test.com"]
  cors_allowed_origins      = ["http://cors_allowed_origin_test.fr", "http://cors_allowed_origin_test.com"]
}