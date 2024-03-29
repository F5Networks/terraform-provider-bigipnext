---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bigipnext_cm_global_resiliency Resource - terraform-provider-bigipnext"
subcategory: ""
description: |-
  Resource used to manage(CRUD) Global Resiliency resources onto BIG-IP Next CM.
---

# bigipnext_cm_global_resiliency (Resource)

Resource used to manage(CRUD) Global Resiliency resources onto BIG-IP Next CM.

## Example Usage

```terraform
resource "bigipnext_cm_global_resiliency" "sample4" {
  name              = "sample4"
  dns_listener_name = "dln"
  dns_listener_port = 10
  protocols         = ["tcp"]
  instances = [
    {
      address              = "10.145.71.115"
      dns_listener_address = "2.2.2.3"
      group_sync_address   = "10.10.1.2/24"
      hostname             = "big-ip-next"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `dns_listener_name` (String) DNS Listener Name. The DNS listener name must start with lowercase letters (a-z) and consist only of lowercase letters (a-z) and digits (0-9).
- `dns_listener_port` (Number) DNS Listener Port. Port number must be greater than or equal to 1. Port number must not exceed 65535. Port cannot be updated once created
- `instances` (Attributes List) List of Instances (see [below for nested schema](#nestedatt--instances))
- `name` (String) Global Resiliency Group Name. The group name must start with lowercase letters (a-z) and consist only of lowercase letters (a-z) and digits (0-9).
- `protocols` (List of String) Protocols to be added to the Global Resiliency Group. Protocols cannot be updated once created.

### Read-Only

- `id` (String) Unique Identifier for the resource

<a id="nestedatt--instances"></a>
### Nested Schema for `instances`

Required:

- `address` (String) Address of the Bip-IP Next. A valid IP Address is required
- `dns_listener_address` (String) DNS Listener Address. A valid IP Address is required
- `group_sync_address` (String) GR Group Sunc IP. A valid IP Address with mask is required
- `hostname` (String) Hostname of the Instance to be added