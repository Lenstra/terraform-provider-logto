resource "logto_sign_in_experience" "test_experience" {
  color = {
    primary_color        = "#4287f5"
    dark_primary_color   = "#5442f5"
    is_dark_mode_enabled = false
  }

  branding = {
    logo_url      = "https://logo_url.test"
    dark_logo_url = "https://dark_logo_url.test"
    favicon       = "https://favicon.test"
    dark_favicon  = "https://dark_favicon.test"
  }

  language_info = {
    auto_detect       = true
    fallback_language = "en"
  }

  terms_of_use_url      = "https://example.com/terms"
  privacy_policy_url    = "https://example.com/privacy"
  agree_to_terms_policy = "Automatic"

  sign_in = {
    methods = [
      {
        identifier          = "email"
        password            = true
        verification_code   = false # Needs connectors to be true
        is_password_primary = true
      },
      {
        identifier          = "phone"
        password            = true
        verification_code   = false # Needs connectors to be true
        is_password_primary = true
      },
    ]
  }

  ### Needs to connectors to use it
  # sign_up = {
  #   identifiers = ["email"],
  #   password    = false,
  #   verify      = false,
  # }

  social_sign_in = {
    automatic_account_linking = true
  }

  social_sign_in_connector_targets = []

  sign_in_mode = "SignIn"
  custom_css   = ".youCutomCss{visibility:hidden;}"

  password_policy = {
    length = {
      min = 8
      max = 128
    }

    character_types = {
      min = 2
    }

    rejects = {
      pwned                   = true
      repetition_and_sequence = true
      user_info               = true
      words                   = ["password", "123456", "admin"]
    }
  }

  mfa = {
    factors = ["Totp"]
    policy  = "UserControlled"
  }

  # Needs conditions to be false
  single_sign_on_enabled = true

  support_email                = "support@email.foo"
  support_website_url          = "https://support_website_url.foo"
  unknown_session_redirect_url = "https://unknown_session_redirect_url.foo"

  captcha_policy = {
    enabled = true
  }

  sentinel_policy = {
    max_attempts     = 10
    lockout_duration = 3600
  }

  email_blocklist_policy = {
    block_disposable_addresses = true
    block_subaddressing        = true
  }
}
