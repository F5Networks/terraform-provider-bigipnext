package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// {\"name\":\"tfnext-cm-vpshere01\",\"type\":\"VSPHERE\",\"connection\":{\"host\":\"mbip-70-vcenter.pdsea.f5net.com\",\"authentication\":{\"type\":\"basic\",\"username\":\"r.chinthalapalli@f5.com\",\"password\":\"TrisRav@2024\"}}}

// {"_links":{"self":{"href":"/api/v1/spaces/default/providers/vsphere/9527cb70-182e-4185-a2f8-b8e6d898d379"}},"connection":{"authentication":{"type":"basic","username":"r.chinthalapalli@f5.com"},"host":"mbip-70-vcenter.pdsea.f5net.com"},"id":"9527cb70-182e-4185-a2f8-b8e6d898d379","name":"tfnext-cm-vpshere01","type":"VSPHERE"}

// {"_links":{"self":{"href":"/api/v1/spaces/default/providers/vsphere/9527cb70-182e-4185-a2f8-b8e6d898d379"}},"connection":{"authentication":{"type":"basic","username":"r.chinthalapalli@f5.com"},"host":"mbip-70-vcenter.pdsea.f5net.com"},"id":"9527cb70-182e-4185-a2f8-b8e6d898d379","name":"tfnext-cm-vpshere01","type":"VSPHERE"}

func TestUnitNextCMDeviceProviderResourceTC1(t *testing.T) {
	testAccPreUnitCheck(t)
	count := 0
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"token": "eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.AbY1hUw8wHO8Vt1qxRd5xQj_21EQ1iaH6q9Z2XgRwQl98M7aCpyjiF2J16S4HrZ-",
			"tokenType": "Bearer",
			"expiresIn": 3600,
			"refreshToken": "ODA0MmQzZTctZTk1Mi00OTk1LWJmMjUtZWZmMjc1NDE3YzliOt4bKlRr6g7RdTtnBKhm2vzkgJeWqfvow68gyxTipleCq4AjR4nxZDBYKQaWyCWGeA",
			"refreshExpiresIn": 1209600
		}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/providers/f5os", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/providers/f5os/f05d420e-781b-409e-a17a-3c321371b3ea"}},"connection":{"authentication":{"type":"basic","username":"admin"},"host":"10.14.10.14:443"},"id":"f05d420e-781b-409e-a17a-3c321371b3ea","name":"tfnext-cm-rseries01","type":"RSERIES"}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/providers/f5os/f05d420e-781b-409e-a17a-3c321371b3ea", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		count++
		t.Logf("\n#####################count: %d\n", count)
		if count >= 4 {
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/providers/f5os/f05d420e-781b-409e-a17a-3c321371b3ea"}},"connection":{"authentication":{"type":"basic","username":"admin"},"host":"10.14.10.14:8888"},"id":"f05d420e-781b-409e-a17a-3c321371b3ea","name":"tfnext-cm-rseries01","type":"RSERIES"}`)
		} else {
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/providers/f5os/f05d420e-781b-409e-a17a-3c321371b3ea"}},"connection":{"authentication":{"type":"basic","username":"admin"},"host":"10.14.10.14:443"},"id":"f05d420e-781b-409e-a17a-3c321371b3ea","name":"tfnext-cm-rseries01","type":"RSERIES"}`)
		}
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMDeviceProviderResourceTC1Config,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config:             testAccNextCMDeviceProviderResourceTC1UpdateConfig,
				Check:              resource.ComposeAggregateTestCheckFunc(),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

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
  name     = "tfnext-cm-rseries01"
  address  = "10.14.10.14:443"
  type     = "RSERIES"
  username = "admin"
  password = "xxxxxxxx"
}
`
const testAccNextCMDeviceProviderResourceTC1UpdateConfig = `
resource "bigipnext_cm_provider" "rseries" {
  name     = "tfnext-cm-rseries01"
  address  = "10.14.10.14:8888"
  type     = "RSERIES"
  username = "admin"
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
