resource "logto_user" "user" {
  name          = "test"
  primary_email = "test@test.com"
  username      = "test username"

  profile = {
    family_name = "test family name"
    given_name  = "test given name"
    middle_name = "test middle name"
    nickname    = "test nickname"
  }
}
