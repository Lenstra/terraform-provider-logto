
resource "logto_role" "role" {
  name        = "role_name"
  description = "role_description"
  type        = "User"
}

resource "logto_user" "user" {
  name          = "user_name"
  primary_email = "user_primary_email@example.fr"

  profile = {
    family_name = "family name"
    given_name  = "given name"
    middle_name = "middle name"
    nickname    = "nickname"
  }

  role_ids = [
    logto_role.role.id
  ]
}
