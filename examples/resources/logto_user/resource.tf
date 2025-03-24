
resource "logto_user" "user" {
  name     = "test"
  username = "test username"

  profile = {
    family_name = "test family name"
    given_name = "test given name"
    middle_name = "test middle name"
    nickname = "test nickname"
  }
}
