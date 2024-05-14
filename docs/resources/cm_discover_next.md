---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bigipnext_cm_discover_next Resource - terraform-provider-bigipnext"
subcategory: ""
description: |-
  Resource used for add   (discover)   BIG-IP Next instance to BIG-IP Next Central Manager for management
---

# bigipnext_cm_discover_next (Resource)

Resource used for add	(discover)	 BIG-IP Next instance to BIG-IP Next Central Manager for management

## Example Usage

```terraform
resource "bigipnext_cm_discover_next" "test" {
  address             = "10.10.10.10"
  port                = 5443
  device_user         = "admin"
  device_password     = "admin123"
  management_user     = "admin-cm"
  management_password = "admin@123"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `address` (String) IP Address of the BIG-IP Next instance to be discovered
- `device_password` (String, Sensitive) The password that the BIG-IP Next Central Manager uses before Instance discovery for BIG-IP Next management
- `device_user` (String) The username that the BIG-IP Next Central Manager uses before Instance discovery for BIG-IP Next management
- `management_password` (String, Sensitive) The password that the BIG-IP Next Central Manager uses after Instance Discovery for BIG-IP Next management
- `management_user` (String) The username that the BIG-IP Next Central Manager uses after Instance Discovery for BIG-IP Next management
- `port` (Number) Port number of the BIG-IP Next instance to be discovered

### Read-Only

- `id` (String) Unique Identifier for the resource