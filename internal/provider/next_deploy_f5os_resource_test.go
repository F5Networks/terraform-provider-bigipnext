package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitNextDeployF5OSResourceTC1(t *testing.T) {
	testAccPreUnitCheck(t)
	// count := 0
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
	mux.HandleFunc("/api/v1/spaces/default/providers", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `[{"instances":[],"provider_host":"10.144.10.146:443","provider_id":"171d5623-e25a-45e5-8e2a-043fc952cbf1","provider_name":"myrseries","provider_type":"RSERIES","provider_username":"admin","updated":"2024-07-16T17:16:42.638833Z"}]`)
	})
	mux.HandleFunc("/api/device/v1/instances", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/7597ff87-d26d-4154-b03d-9e7999d24f0e"}},"path":"/v1/instances/tasks/7597ff87-d26d-4154-b03d-9e7999d24f0e"}`)
	})

	mux.HandleFunc("/api/device/v1/instances/tasks/7597ff87-d26d-4154-b03d-9e7999d24f0e", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/7597ff87-d26d-4154-b03d-9e7999d24f0e"}},"completed":"2024-07-16T17:25:01.919956Z","created":"2024-07-16T17:16:46.036781Z","failure_reason":"","id":"7597ff87-d26d-4154-b03d-9e7999d24f0e","name":"instance creation","payload":{"discovery":{"port":5443,"address":"10.144.10.182","device_user":"admin","device_password":"*****","management_user":"admin-cm","management_password":"*****"},"onboarding":{"mode":"STANDALONE","nodes":[{"password":"*****","username":"admin","managementAddress":"10.144.10.182"}],"platformType":"RSERIES"},"instantiation":{"Request":{"F5osRequest":{"provider_id":"171d5623-e25a-45e5-8e2a-043fc952cbf1","provider_type":"rseries","next_instances":[{"nodes":[1],"vlans":[444,555],"mgmt_ip":"10.144.10.182","timeout":600,"hostname":"demovm01-ravi-r10800","cpu_cores":4,"disk_size":30,"mgmt_prefix":24,"mgmt_gateway":"10.144.10.254","admin_password":"*****","tenant_image_name":"BIG-IP-Next-20.2.1-2.430.2+0.0.48","tenant_deployment_file":"BIG-IP-Next-20.2.1-2.430.2+0.0.48.yaml"}]},"VsphereRequest":null},"BaseTask":{"id":"","payload":null,"provider_id":"171d5623-e25a-45e5-8e2a-043fc952cbf1","provider_type":"rseries"},"VsphereRequest":null}},"stage":"Discovery","state":"discoveryDone","status":"completed","task_type":"instance_creation","updated":"2024-07-16T17:25:01.919956Z"}`)
	})

	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_embedded":{"devices":[{"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800/6cdf38ed-a258-4d92-a64d-972238b27400'"}},"address":"10.144.10.182","certificate_validated":"2024-07-16T17:24:12.848618Z","certificate_validity":false,"hostname":"demovm01-ravi-r10800","id":"6cdf38ed-a258-4d92-a64d-972238b27400","mode":"STANDALONE","platform_name":"R10K","platform_type":"APPLIANCE","port":5443,"short_id":"bcpED8hJ","version":"20.2.1-2.430.2+0.0.48"}]},"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800"}},"count":1,"total":1}`)
	})
	// mux.HandleFunc("/api/device/v1/inventory?filter=hostname+eq+'demovm01-ravi-r10800'", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	_, _ = fmt.Fprintf(w, `{"_embedded":{"devices":[{"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800/6cdf38ed-a258-4d92-a64d-972238b27400'"}},"address":"10.144.10.182","certificate_validated":"2024-07-16T17:24:12.848618Z","certificate_validity":false,"hostname":"demovm01-ravi-r10800","id":"6cdf38ed-a258-4d92-a64d-972238b27400","mode":"STANDALONE","platform_name":"R10K","platform_type":"APPLIANCE","port":5443,"short_id":"bcpED8hJ","version":"20.2.1-2.430.2+0.0.48"}]},"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800"}},"count":1,"total":1}`)
	// })
	mux.HandleFunc("/api/v1/spaces/default/instances/6cdf38ed-a258-4d92-a64d-972238b27400", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/instances/deletion-tasks/837120e1-7420-4e86-afc5-213dd2b07f26"}},"path":"/api/v1/spaces/default/instances/deletion-tasks/837120e1-7420-4e86-afc5-213dd2b07f26"}`)
		} else {
			_, _ = fmt.Fprintf(w, `{"_embedded":{"devices":[{"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800/6cdf38ed-a258-4d92-a64d-972238b27400'"}},"address":"10.144.10.182","certificate_validated":"2024-07-16T17:24:12.848618Z","certificate_validity":false,"hostname":"demovm01-ravi-r10800","id":"6cdf38ed-a258-4d92-a64d-972238b27400","mode":"STANDALONE","platform_name":"R10K","platform_type":"APPLIANCE","port":5443,"short_id":"bcpED8hJ","version":"20.2.1-2.430.2+0.0.48"}]},"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800"}},"count":1,"total":1}`)
		}
	})
	// mux.HandleFunc("/api/v1/spaces/default/instances/6cdf38ed-a258-4d92-a64d-972238b27400", func(w http.ResponseWriter, r *http.Request) {
	//
	// })
	mux.HandleFunc("/api/device/v1/deletion-tasks/837120e1-7420-4e86-afc5-213dd2b07f26", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/deletion-tasks/837120e1-7420-4e86-afc5-213dd2b07f26"}},"address":"10.144.10.182","completed":"2024-07-16T17:35:42.811049Z","created":"2024-07-16T17:32:36.411802Z","device_id":"6cdf38ed-a258-4d92-a64d-972238b27400","failure_reason":"","id":"837120e1-7420-4e86-afc5-213dd2b07f26","state":"instanceRemovalDone","status":"completed"}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextDeployF5OSResourceTC1Config,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// ImportState testing
			{
				ResourceName:      "bigipnext_cm_deploy_f5os.rseries01",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestUnitNextDeployF5OSResourceTC2(t *testing.T) {
	testAccPreUnitCheck(t)
	// count := 0
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
	mux.HandleFunc("/api/v1/spaces/default/providers", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `[{"instances":[],"provider_host":"10.144.10.146:443","provider_id":"171d5623-e25a-45e5-8e2a-043fc952cbf1","provider_name":"myvelos","provider_type":"VELOS","provider_username":"admin","updated":"2024-07-16T17:16:42.638833Z"}]`)
	})
	mux.HandleFunc("/api/device/v1/instances", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/7597ff87-d26d-4154-b03d-9e7999d24f0e"}},"path":"/v1/instances/tasks/7597ff87-d26d-4154-b03d-9e7999d24f0e"}`)
	})

	mux.HandleFunc("/api/device/v1/instances/tasks/7597ff87-d26d-4154-b03d-9e7999d24f0e", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/7597ff87-d26d-4154-b03d-9e7999d24f0e"}},"completed":"2024-07-16T17:25:01.919956Z","created":"2024-07-16T17:16:46.036781Z","failure_reason":"","id":"7597ff87-d26d-4154-b03d-9e7999d24f0e","name":"instance creation","payload":{"discovery":{"port":5443,"address":"10.144.10.182","device_user":"admin","device_password":"*****","management_user":"admin-cm","management_password":"*****"},"onboarding":{"mode":"STANDALONE","nodes":[{"password":"*****","username":"admin","managementAddress":"10.144.10.182"}],"platformType":"RSERIES"},"instantiation":{"Request":{"F5osRequest":{"provider_id":"171d5623-e25a-45e5-8e2a-043fc952cbf1","provider_type":"rseries","next_instances":[{"nodes":[1],"vlans":[444,555],"mgmt_ip":"10.144.10.182","timeout":600,"hostname":"demovm01-ravi-r10800","cpu_cores":4,"disk_size":30,"mgmt_prefix":24,"mgmt_gateway":"10.144.10.254","admin_password":"*****","tenant_image_name":"BIG-IP-Next-20.2.1-2.430.2+0.0.48","tenant_deployment_file":"BIG-IP-Next-20.2.1-2.430.2+0.0.48.yaml"}]},"VsphereRequest":null},"BaseTask":{"id":"","payload":null,"provider_id":"171d5623-e25a-45e5-8e2a-043fc952cbf1","provider_type":"rseries"},"VsphereRequest":null}},"stage":"Discovery","state":"discoveryDone","status":"completed","task_type":"instance_creation","updated":"2024-07-16T17:25:01.919956Z"}`)
	})

	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_embedded":{"devices":[{"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800/6cdf38ed-a258-4d92-a64d-972238b27400'"}},"address":"10.144.10.182","certificate_validated":"2024-07-16T17:24:12.848618Z","certificate_validity":false,"hostname":"demovm01-ravi-r10800","id":"6cdf38ed-a258-4d92-a64d-972238b27400","mode":"STANDALONE","platform_name":"R10K","platform_type":"APPLIANCE","port":5443,"short_id":"bcpED8hJ","version":"20.2.1-2.430.2+0.0.48"}]},"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800"}},"count":1,"total":1}`)
	})
	// mux.HandleFunc("/api/device/v1/inventory?filter=hostname+eq+'demovm01-ravi-r10800'", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	_, _ = fmt.Fprintf(w, `{"_embedded":{"devices":[{"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800/6cdf38ed-a258-4d92-a64d-972238b27400'"}},"address":"10.144.10.182","certificate_validated":"2024-07-16T17:24:12.848618Z","certificate_validity":false,"hostname":"demovm01-ravi-r10800","id":"6cdf38ed-a258-4d92-a64d-972238b27400","mode":"STANDALONE","platform_name":"R10K","platform_type":"APPLIANCE","port":5443,"short_id":"bcpED8hJ","version":"20.2.1-2.430.2+0.0.48"}]},"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800"}},"count":1,"total":1}`)
	// })
	mux.HandleFunc("/api/v1/spaces/default/instances/6cdf38ed-a258-4d92-a64d-972238b27400", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/instances/deletion-tasks/837120e1-7420-4e86-afc5-213dd2b07f26"}},"path":"/api/v1/spaces/default/instances/deletion-tasks/837120e1-7420-4e86-afc5-213dd2b07f26"}`)
		} else {
			_, _ = fmt.Fprintf(w, `{"_embedded":{"devices":[{"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800/6cdf38ed-a258-4d92-a64d-972238b27400'"}},"address":"10.144.10.182","certificate_validated":"2024-07-16T17:24:12.848618Z","certificate_validity":false,"hostname":"demovm01-ravi-r10800","id":"6cdf38ed-a258-4d92-a64d-972238b27400","mode":"STANDALONE","platform_name":"R10K","platform_type":"APPLIANCE","port":5443,"short_id":"bcpED8hJ","version":"20.2.1-2.430.2+0.0.48"}]},"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+demovm01-ravi-r10800"}},"count":1,"total":1}`)
		}
	})
	// mux.HandleFunc("/api/v1/spaces/default/instances/6cdf38ed-a258-4d92-a64d-972238b27400", func(w http.ResponseWriter, r *http.Request) {
	//
	// })
	mux.HandleFunc("/api/device/v1/deletion-tasks/837120e1-7420-4e86-afc5-213dd2b07f26", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/deletion-tasks/837120e1-7420-4e86-afc5-213dd2b07f26"}},"address":"10.144.10.182","completed":"2024-07-16T17:35:42.811049Z","created":"2024-07-16T17:32:36.411802Z","device_id":"6cdf38ed-a258-4d92-a64d-972238b27400","failure_reason":"","id":"837120e1-7420-4e86-afc5-213dd2b07f26","state":"instanceRemovalDone","status":"completed"}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextDeployF5OSResourceTC2Config,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: testAccNextDeployF5OSResourceTC2UpdateConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

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
	  slot_ids               = [1]
	  tenant_deployment_file = "BIG-IP-Next-20.1.0-2.279.0+0.0.75.yaml"
	  tenant_image_name      = "BIG-IP-Next-20.1.0-2.279.0+0.0.75"
	}
}
`
const testAccNextDeployF5OSResourceTC2Config = `
resource "bigipnext_cm_deploy_f5os" "velos01" {
	f5os_provider = {
	  provider_name = "myvelos"
	  provider_type = "velos"
	}
	instance = {
	  instance_hostname      = "rseriesravitest04"
	  management_address     = "10.144.140.81"
	  management_prefix      = 24
	  management_gateway     = "10.144.140.254"
	  management_user        = "admin-cm"
	  management_password    = "F5Twist@123"
	  vlan_ids               = [27, 28, 29]
	  slot_ids               = [2]
	  tenant_deployment_file = "BIG-IP-Next-20.1.0-2.279.0+0.0.75.yaml"
	  tenant_image_name      = "BIG-IP-Next-20.1.0-2.279.0+0.0.75"
	}
}
`

const testAccNextDeployF5OSResourceTC2UpdateConfig = `
resource "bigipnext_cm_deploy_f5os" "velos01" {
	f5os_provider = {
	  provider_name = "myvelos"
	  provider_type = "velos"
	}
	instance = {
	  instance_hostname      = "rseriesravitest04"
	  management_address     = "10.144.140.82"
	  management_prefix      = 24
	  management_gateway     = "10.144.140.254"
	  management_user        = "admin-cm"
	  management_password    = "F5Twist@123"
	  vlan_ids               = [27, 28, 29]
	  slot_ids               = [2]
	  tenant_deployment_file = "BIG-IP-Next-20.1.0-2.279.0+0.0.75.yaml"
	  tenant_image_name      = "BIG-IP-Next-20.1.0-2.279.0+0.0.75"
	}
}
`
