
resource "logto_role" "role" {
  name        = "role_name"
  description = "role_description"
  type        = "User"
}

resource "logto_user" "test_user" {
  name          = "test_user_modified"
  primary_email = "test_user@test.fr"

  profile = {
    family_name = "test_family_name_modified"
    given_name  = "test_given_name_modified"
    middle_name = "test_middle_name_modified"
    nickname    = "test_nickname_modified"
  }

  role_ids = [
    logto_role.role.id
  ]
}
