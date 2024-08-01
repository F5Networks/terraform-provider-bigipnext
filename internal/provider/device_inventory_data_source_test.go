package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitNextCMDeviceInventoryDatasourceTC1(t *testing.T) {
	testAccPreUnitCheck(t)
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
	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_embedded":{"devices":[{"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800/6cdf38ed-a258-4d92-a64d-972238b27400'"}},"address":"10.144.10.182","certificate_validated":"2024-07-16T17:24:12.848618Z","certificate_validity":false,"hostname":"demovm01-ravi-r10800","id":"6cdf38ed-a258-4d92-a64d-972238b27400","mode":"STANDALONE","platform_name":"R10K","platform_type":"APPLIANCE","port":5443,"short_id":"bcpED8hJ","version":"20.2.1-2.430.2+0.0.48"}]},"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800"}},"count":1,"total":1}`)
	})
	// mux.HandleFunc("/api/v1/spaces/default/providers/f5os/f05d420e-781b-409e-a17a-3c321371b3ea", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	count++
	// 	t.Logf("\n#####################count: %d\n", count)
	// 	if count >= 4 {
	// 		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/providers/f5os/f05d420e-781b-409e-a17a-3c321371b3ea"}},"connection":{"authentication":{"type":"basic","username":"admin"},"host":"10.14.10.14:8888"},"id":"f05d420e-781b-409e-a17a-3c321371b3ea","name":"tfnext-cm-rseries01","type":"RSERIES"}`)
	// 	} else {
	// 		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/providers/f5os/f05d420e-781b-409e-a17a-3c321371b3ea"}},"connection":{"authentication":{"type":"basic","username":"admin"},"host":"10.14.10.14:443"},"id":"f05d420e-781b-409e-a17a-3c321371b3ea","name":"tfnext-cm-rseries01","type":"RSERIES"}`)
	// 	}
	// })
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMDeviceInventoryDarasourceTC1Config,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

const testAccNextCMDeviceInventoryDarasourceTC1Config = `
data "bigipnext_cm_device_inventory" "test" {}
`
