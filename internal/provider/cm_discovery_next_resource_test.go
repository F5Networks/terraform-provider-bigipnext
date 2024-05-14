package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitCMDiscoveryNextResourceTC1(t *testing.T) {
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
	mux.HandleFunc("/api/v1/spaces/default/instances", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links": {"self": {"href": "/api/v1/spaces/default/instances/discovery-tasks/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2"}},"path": "/api/v1/spaces/default/instances/discovery-tasks/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2"}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/instances/discovery-tasks/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
				"_links": {
				  "self": {
					"href": "/api/v1/spaces/default/instances/discovery-tasks/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2"
				  }
				},
				"address": "10.145.10.1",
				"completed": "2021-04-02T23:11:19.08911Z",
				"created": "2021-04-02T23:11:18.051859Z",
				"device_group": "default",
				"device_user": "admin",
				"discovered_device_id": "c9796e86-21f7-4182-be1c-c737ed430242",
				"discovered_device_path": "/api/v1/spaces/default/instances/c9796e86-21f7-4182-be1c-c737ed430242",
				"failure_reason": "",
				"id": "43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2",
				"port": 5443,
				"state": "addingDeviceToGroup",
				"status": "completed"
			  }`)
	})

	mux.HandleFunc("/api/v1/spaces/default/instances/c9796e86-21f7-4182-be1c-c737ed430242", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{
			"address": "10.10.10.10",
			"hostname": "10.10.10.10",
			"id": "c9796e86-21f7-4182-be1c-c737ed430242",
			"port": 0,
			"version": "string",
			"certificate_validity": true,
			"certificate_validation_error": "string",
			"certificate_validated": "2019-08-24T14:15:22Z",
			"mode": "STANDALONE",
			"platform_type": "VE",
			"status": "HEALTHY"
		  }`)
		} else if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"_links": {
				  "self": {
					"href": "/api/v1/spaces/default/instances/discovery-tasks/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2"
				  }
				},
				"path": "/api/v1/spaces/default/instances/discovery-tasks/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2"
			  }`)
		}
	})
	mux.HandleFunc("/api/device/v1/deletion-tasks/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
			  "self": {
				"href": "/api/v1/spaces/default/instances/deletion-tasks/e8dd5c25-3357-4a05-825d-b183f6999b3e"
			  }
			},
			"address": "10.145.66.196",
			"completed": "2021-04-07T17:48:45.630172Z",
			"created": "2021-04-07T17:48:45.113981Z",
			"device_id": "c9796e86-21f7-4182-be1c-c737ed430242",
			"failure_reason": "",
			"id": "e8dd5c25-3357-4a05-825d-b183f6999b3e",
			"state": "deletingDevice",
			"status": "completed"
		  }`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testUnitCMDiscoveryNextResourceTC1,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: testUnitCMDiscoveryNextResourceTC2,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

const testUnitCMDiscoveryNextResourceTC1 = `
resource "bigipnext_cm_discover_next" "test" {
	address             = "10.10.10.10"
	port                = 5443
	device_user         = "admin"
	device_password     = "admin123"
	management_user     = "admin-cm"
	management_password = "admin@123"
  }
`

const testUnitCMDiscoveryNextResourceTC2 = `
resource "bigipnext_cm_discover_next" "test" {
	address             = "10.10.10.10"
	port                = 5443
	device_user         = "admin"
	device_password     = "admin@123"
	management_user     = "admin-cm"
	management_password = "admin@123"
  }
`
