---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bigipnext_cm_ha_cluster Resource - terraform-provider-bigipnext"
subcategory: ""
description: |-
  Create a HA Cluster of BIG-IP Next Central Manager instances
---

# bigipnext_cm_ha_cluster (Resource)

Create a HA Cluster of BIG-IP Next Central Manager instances



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `nodes` (Attributes List) (see [below for nested schema](#nestedatt--nodes))

### Read-Only

- `agent_nodes` (List of String) List of nodes that are marked as agent nodes
- `id` (String) The ID of the resource
- `server_nodes` (List of String) List of nodes that are marked as control plane nodes

<a id="nestedatt--nodes"></a>
### Nested Schema for `nodes`

Required:

- `node_ip` (String) IP address of the node that will be added to the cluster
- `password` (String, Sensitive) The password of the node
- `username` (String) The username of the node

Optional:

- `fingerprint` (String) The fingerprint of the node in the SHA256 format