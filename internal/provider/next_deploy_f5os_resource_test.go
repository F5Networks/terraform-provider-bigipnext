package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNextDeployF5OSResourceTC1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextDeployF5OSResourceTC1Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs(
						"bigipnext_cm_deploy_f5os.rseries01",
						"instance",
						map[string]string{
							"instance_hostname":      "rseriesravitest04",
							"management_address":     "10.144.140.81",
							"management_prefix":      "24",
							"management_gateway":     "10.144.140.254",
							"management_user":        "admin-cm",
							"tenant_image_name":      "BIG-IP-Next-20.1.0-2.279.0+0.0.75",
							"tenant_deployment_file": "BIG-IP-Next-20.1.0-2.279.0+0.0.75.yaml",
						},
					),
				),
			},
		},
	})
}

const testAccNextDeployF5OSResourceTC1Config = `
resource "bigipnext_cm_deploy_f5os" "rseries01" {
	f5os_provider = {
	  provider_name = "myrseries"
	  provider_type = "rseries"
	}
	instance = {
	  instance_hostname      = "rseriesravitest04"
	  management_address     = "10.144.140.81"
	  management_prefix      = 24
	  management_gateway     = "10.144.140.254"
	  management_user        = "admin-cm"
	  management_password    = "F5Twist@123"
	  vlan_ids               = [27, 28, 29]
	  tenant_deployment_file = "BIG-IP-Next-20.1.0-2.279.0+0.0.75.yaml"
	  tenant_image_name      = "BIG-IP-Next-20.1.0-2.279.0+0.0.75"
	}
}
`
