package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCMHAClusterTC(t *testing.T) {
	control_node := os.Getenv("BIGIPNEXT_HOST")
	id := fmt.Sprintf("central-manager-server-%s", control_node)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCMHAClusterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_ha_cluster.cm_ha_2_nodes", "server_nodes.0", "10.146.164.150"),
					resource.TestCheckResourceAttr("bigipnext_cm_ha_cluster.cm_ha_2_nodes", "id", id),
				),
				Destroy: false,
			},
			{
				Config: testAccCMHAClusterUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_ha_cluster.cm_ha_2_nodes", "server_nodes.0", "10.146.164.150"),
					resource.TestCheckResourceAttr("bigipnext_cm_ha_cluster.cm_ha_2_nodes", "server_nodes.1", "10.146.165.89"),
					resource.TestCheckResourceAttr("bigipnext_cm_ha_cluster.cm_ha_2_nodes", "id", id),
				),
				Destroy: false,
			},
		},
	})
}

const testAccCMHAClusterConfig = `
resource "bigipnext_cm_ha_cluster" "cm_ha_2_nodes" {
  nodes = [
    {
        node_ip = "10.146.164.150"
        username = "admin",
        password = "F5site02@123"
    }
  ]
}
`

const testAccCMHAClusterUpdateConfig = `
resource "bigipnext_cm_ha_cluster" "cm_ha_2_nodes" {
  nodes = [
    {
        node_ip = "10.146.164.150"
        username = "admin",
        password = "F5site02@123"
    },
    {
        node_ip = "10.146.165.89"
        username = "admin",
        password = "F5site02@123"
    }
  ]
}
`
