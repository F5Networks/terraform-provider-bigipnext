---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bigipnext Provider"
subcategory: ""
description: |-
  Provider plugin to interact with BIG-IP Next Central Manager(CM) Using OpenAPI
---

# bigipnext Provider

Provider plugin to interact with BIG-IP Next Central Manager(CM) Using OpenAPI

## Example Usage

```terraform
provider "bigipnext" {
  username = "education"
  password = "test123"
  host     = "http://localhost:19090"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `host` (String) URI for BigipNext Device. May also be provided via `BIGIPNEXT_HOST` environment variable.
- `password` (String, Sensitive) Password for BigipNext Device. May also be provided via `BIGIPNEXT_PASSWORD` environment variable.
- `port` (Number) Port Number to be used to make API calls to HOST
- `username` (String) Username for BigipNext Device. May also be provided via `BIGIPNEXT_USERNAME` environment variable.
