---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bigipnext_cm_add_jwt_token Resource - terraform-provider-bigipnext"
subcategory: ""
description: |-
  Resource used for add/copy JWT Token on Central Manager
---

# bigipnext_cm_add_jwt_token (Resource)

Resource used for add/copy JWT Token on Central Manager

## Example Usage

```terraform
resource "bigipnext_cm_add_jwt_token" "tokenadd" {
  token_name = "paid_test_jwt"
  jwt_token  = "eyJhbG"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `jwt_token` (String, Sensitive) JWT token to be added on Central Manager
- `token_name` (String) Nickname to be used to add the JWT token on Central Manager

### Read-Only

- `id` (String) Unique Identifier for the resource
- `order_type` (String) JWT token to be added on Central Manager
- `subscription_expiry` (String) JWT token to be added on Central Manager
