package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNextCMDeviceProviderResourceTC1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMDeviceProviderResourceTC1Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_provider.rseries", "type", "RSERIES"),
					resource.TestCheckResourceAttr("bigipnext_cm_provider.rseries", "address", "10.144.140.80:443"),
					resource.TestCheckResourceAttr("bigipnext_cm_provider.rseries", "name", "testrseriesprovvm01"),
				),
			},
		},
	})
}

func TestAccNextCMDeviceProviderResourceTC2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMDeviceProviderResourceTC2Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_provider.vpshere", "type", "VSPHERE"),
					resource.TestCheckResourceAttr("bigipnext_cm_provider.vpshere", "address", "mbip-70-vcenter.pdsea.f5net.com"),
					resource.TestCheckResourceAttr("bigipnext_cm_provider.vpshere", "name", "testvpshereprovvm01"),
				),
			},
		},
	})
}

const testAccNextCMDeviceProviderResourceTC1Config = `
resource "bigipnext_cm_provider" "rseries" {
  name     = "testrseriesprovvm01"
  address  = "10.14.1.80:443"
  type     = "RSERIES"
  username = "xxxxxx"
  password = "xxxxxxxx"
}
`
const testAccNextCMDeviceProviderResourceTC2Config = `
resource "bigipnext_cm_provider" "vpshere" {
  name     = "testvpshereprovvm01"
  address  = "mbip-70-vcenter.pdsea.f5net.com"
  type     = "VSPHERE"
  username = "xxxxxxxxxx"
  password = "xxxxxx"
}
`
