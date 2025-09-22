## 0.0.13

IMPROVEMENTS:

- The `logto_user` now has a `role_ids` attribute. (#22)

NOTES:
- Update of the documentation examples (#21)

## 0.0.12

FEATURES:

- **New Resource:** `role` (#11)
- **New Resource:** `api_resource` (#11)
- **New Resource:** `api_resource_scope` (#11)

ENHANCEMENTS:

- `terraform-generator`: Add support for multiple word resource names (#17)

NOTES:
- Update of the golangci-lint (#16)

## 0.0.11

BUG FIXES:

- Fix handling of argument `cors_allowed_origins` in resource `logto_application`.

## 0.0.10

BUG FIXES:

- Fix handling of `profile` argument in `logto_user` resource.

## 0.0.9

IMPROVEMENTS:

- Log more information on errors when calling the API.

## 0.0.8

BUG FIXES:

- Revert changes from v0.0.7.

## 0.0.7

BUG FIXES:

- Fixed user creation error by handling empty profile fields correctly in state update.

## 0.0.6

BUG FIXES:

- Detect when a user has been deleted from Logto.

## 0.0.5

IMPROVEMENTS:

- The `logto_user` now has a `primary_email` attribute.

## 0.0.4

FEATURES:

- **New `logto_user` resource.**

## 0.0.3

IMPROVEMENTS:

- Add support for `is_admin` and `is_third_party` to `logto_application` resource.

## 0.0.2

BUG FIXES:

- Fix authentication for local deployment of Logto.

## 0.0.1

FEATURES:

- **New `logto_application` resource.**
