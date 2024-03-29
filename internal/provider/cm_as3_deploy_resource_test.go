package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNextCMAS3CreateTC1Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMAS3ResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_as3.test", "id", "10.144.72.114"),
				),
			},
			{
				Config: testAccNextCMAS3ResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_as3.test", "id", "10.144.72.114"),
				),
			},
			//// ImportState testing
			//{
			//	ResourceName:      "bigipnext_cm__as3.test",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
		},
	})
}

const testAccNextCMAS3ResourceConfig = `
resource "bigipnext_cm_as3" "test" {
  as3_json = <<EOT
{
    "class": "AS3",
    "action": "deploy",
    "persist": true,
    "declaration": {
        "class": "ADC",
        "schemaVersion": "3.45.0",
        "id": "example-declaration-01",
        "label": "Sample 1",
        "remark": "Simple HTTP application with round robin pool",
        "target": {
            "address": "10.144.72.114"
        },
        "next-cm-tenant01": {
            "class": "Tenant",
            "next-cm-app01": {
                "class": "Application",
                "template": "http",
                "serviceMain": {
                    "class": "Service_HTTP",
                    "virtualAddresses": [
                        "10.0.1.10"
                    ],
                    "pool": "next-cm-pool01"
                },
                "next-cm-pool01": {
                    "class": "Pool",
                    "monitors": [
                        "http"
                    ],
                    "members": [
                        {
                            "servicePort": 80,
                            "serverAddresses": [
                                "192.0.2.100",
                                "192.0.2.110"
                            ]
                        }
                    ]
                }
            }
        }
    }
}

EOT
}
`

const testAccNextCMAS3ResourceUpdateConfig = `
resource "bigipnext_cm_as3" "test" {
  as3_json = <<EOT
{
    "class": "AS3",
    "action": "deploy",
    "persist": true,
    "declaration": {
        "class": "ADC",
        "schemaVersion": "3.45.0",
        "id": "example-declaration-01",
        "label": "Sample 1",
        "remark": "Simple HTTP application with round robin pool",
        "target": {
            "address": "10.144.72.114"
        },
        "next-cm-tenant01": {
            "class": "Tenant",
            "next-cm-app01": {
                "class": "Application",
                "template": "http",
                "serviceMain": {
                    "class": "Service_HTTP",
                    "virtualAddresses": [
                        "10.0.1.10"
                    ],
                    "pool": "next-cm-pool01"
                },
                "next-cm-pool01": {
                    "class": "Pool",
                    "monitors": [
                        "http"
                    ],
                    "members": [
                        {
                            "servicePort": 80,
                            "serverAddresses": [
                                "192.0.2.100",
                                "192.0.2.110"
                            ]
                        }
                    ]
                }
            }
        }
    }
}

EOT
}
`
