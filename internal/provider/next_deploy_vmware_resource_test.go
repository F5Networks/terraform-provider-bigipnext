package provider

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNextDeployVmwareResourceTC1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextDeployVmwareResourceTC1Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_deploy_vmware.vmware", "ntp_servers.0", "0.us.pool.ntp.org"),
					resource.TestCheckResourceAttr("bigipnext_cm_deploy_vmware.vmware", "dns_servers.0", "8.8.8.8"),
					// resource.TestCheckTypeSetElemNestedAttrs(
					// 	"bigipnext_cm_deploy_vmware.vmware",
					// 	"vsphere_provider.*",
					// 	map[string]string{
					// 		"provider_name":      "myvsphere03",
					// 		"content_library":    "CM-IOD",
					// 		"cluster_name":       "vSAN Cluster",
					// 		"datacenter_name":    "mbip-7.0",
					// 		"datastore_name":     "vsanDatastore",
					// 		"resource_pool_name": "INFRAANO",
					// 		"vm_template_name":   "BIG-IP-Next-20.1.0-2.279.0+0.0.75-VM-template",
					// 	},
					// ),
				),
			},
		},
	})
}

func TestUnitNextDeployVmwareResourceTC1(t *testing.T) {
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
	mux.HandleFunc("/api/v1/spaces/default/providers?filter=name+eq+'myvsphere03'", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/certificates/create"}},"path":"/v1/certificates/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2"}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testAccNextDeployVmwareResourceTC1Config,
				ExpectError: regexp.MustCompile(`Failed to get provider ID:, got error`),
			},
		},
	})
}

func TestUnitNextDeployVmwareResourceTC2(t *testing.T) {
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
	mux.HandleFunc("/api/v1/spaces/default/providers", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `[{"instances":[{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"},{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"},{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"}],"provider_host":"mbip-70-vcenter.pdsea.f5net.com","provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","provider_name":"myvsphere03","provider_type":"VSPHERE","provider_username":"r.chinthalapalli@f5.com","updated":"2023-12-12T11:53:58.614102Z"}]`)
	})
	mux.HandleFunc("/api/v1/system/infra/info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"app_version": "0.179.6","ha_status": "Not Running","num_of_nodes": 1,"version": "BIG-IP-Next-CentralManager-20.3.0-0.14.42"}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		t.Logf("RequestURI: %s", r.RequestURI)
		if r.RequestURI == "/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api?path=api/vcenter/datacenter" {
			_, _ = fmt.Fprintf(w, `[{"name": "mbip-7.0","datacenter_id": "datacenter-3"}]`)
		}
		if r.RequestURI == "/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api?path=api/vcenter/cluster?datacenters=datacenter-3" {
			_, _ = fmt.Fprintf(w, `[{"name": "vSAN Cluster","cluster_id": "domain-c8"}]`)
		}
		if r.RequestURI == "/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api?path=api/vcenter/resource-pool?clusters=domain-c8" {
			_, _ = fmt.Fprintf(w, `[{"name": "INFRAANO","resource_pool_id": "resgroup-4047"},{"name": "MBIPMP_System","resource_pool_id": "resgroup-4049"}]`)
		}
	})
	mux.HandleFunc("/api/device/v1/instances", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438"}},"path":"/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438"}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/instances/initialization/tasks/deacca61-3162-4672-aac8-2d6bd2b69438", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/4b6fadb6-d44f-4a40-8e73-c2d0be7ca55f"}},"completed":"2024-03-07T06:57:00.298908Z","created":"2024-03-07T06:48:38.583246Z","failure_reason":"","id":"4b6fadb6-d44f-4a40-8e73-c2d0be7ca55f","name":"instance creation","payload":{"discovery":{"port":5443,"address":"10.146.194.171","device_user":"admin","device_password":"*****","management_user":"admin-cm","management_password":"*****"},"onboarding":{"mode":"STANDALONE","nodes":[{"hostname":"infraanovm01","password":"*****","username":"admin","managementAddress":"10.146.194.171"}],"siteInfo":{"dnsServers":["8.8.8.8"],"ntpServers":["0.us.pool.ntp.org"]},"platformType":"VE"},"instantiation":{"Request":{"F5osRequest":null,"VsphereRequest":{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","polling_delay":15,"provider_type":"vsphere","next_instances":[{"memory":16384,"cluster":"vSAN Cluster","num_cpus":8,"datastore":"vsanDatastore","datacenter":"mbip-7.0","sleep_time":"360s","resource_pool":"INFRAANO","mgmt_dns_server":"","vsphere_template":"BIG-IP-Next-20.1.0-2.279.0+0.0.75-VM-template","bigipnext_vm_name":"infraanovm01","mgmt_ipv4_address":"10.146.194.171/23","mgmt_ipv4_gateway":"10.146.195.254","mgmt_network_name":"VM-mgmt","bigipnext_vm_password":"*****","external_network_name":"LocalTestVLAN-115","internal_network_name":"LocalTestVLAN-114","vsphere_content_library":"CM-IOD","ha_data_plane_network_name":"LocalTestVLAN-116","ha_control_plane_network_name":""}]}},"BaseTask":{"id":"","payload":null,"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","provider_type":"vsphere"},"VsphereRequest":null}},"stage":"Discovery","state":"discoveryDone","status":"completed","task_type":"instance_creation","updated":"2024-03-07T06:57:00.298908Z"}`)
	})
	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_embedded":{"devices":[{"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+infraanovm01/37021437-2b5f-4e44-9f5e-6ea9838c5f7e"}},"address":"10.146.194.171","certificate_validated":"2024-03-07T06:56:51.915056Z","certificate_validity":false,"hostname":"infraanovm01","id":"37021437-2b5f-4e44-9f5e-6ea9838c5f7e","mode":"STANDALONE","platform_name":"VMware","platform_type":"VE","port":5443,"version":"20.1.0-2.279.0+0.0.75"}]},"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+infraanovm01"}},"count":1,"total":1}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/instances/37021437-2b5f-4e44-9f5e-6ea9838c5f7e", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/deletion-tasks/02752890-5660-450c-ace9-b8e0a86a15ad"}},"path":"/v1/deletion-tasks/02752890-5660-450c-ace9-b8e0a86a15ad"}`)
	})
	mux.HandleFunc("/api/device/v1/deletion-tasks/02752890-5660-450c-ace9-b8e0a86a15ad", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/deletion-tasks/642d5964-8cd9-4881-9086-1ed5ca682101"}},"address":"10.146.168.20","created":"2023-11-28T07:55:50.924918Z","device_id":"8d6c8c85-1738-4a34-b57b-d3644a2ecfcc","id":"642d5964-8cd9-4881-9086-1ed5ca682101","state":"factoryResetInstance","status":"completed"}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextDeployVmwareResourceTC1Config,
			},
			{
				Config: testAccNextDeployVmwareResourceTC2Config,
			},
		},
	})
}

func TestUnitNextDeployVmwareResourceTC3(t *testing.T) {
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
	mux.HandleFunc("/api/v1/spaces/default/providers", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `[{"instances":[{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"},{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"},{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"}],"provider_host":"mbip-70-vcenter.pdsea.f5net.com","provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","provider_name":"myvsphere03","provider_type":"VSPHERE","provider_username":"r.chinthalapalli@f5.com","updated":"2023-12-12T11:53:58.614102Z"}]`)
	})
	mux.HandleFunc("/api/v1/system/infra/info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"app_version": "0.179.6","ha_status": "Not Running","num_of_nodes": 1,"version": "BIG-IP-Next-CentralManager-20.3.0-0.14.42"}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		t.Logf("RequestURI: %s", r.RequestURI)
		if r.RequestURI == "/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api?path=api/vcenter/datacenter" {
			_, _ = fmt.Fprintf(w, `[{"name": "mbip-7.0","datacenter_id": "datacenter-3"}]`)
		}
		if r.RequestURI == "/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api?path=api/vcenter/cluster?datacenters=datacenter-3" {
			_, _ = fmt.Fprintf(w, `[{"name": "vSAN Cluster","cluster_id": "domain-c8"}]`)
		}
		if r.RequestURI == "/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api?path=api/vcenter/resource-pool?clusters=domain-c8" {
			_, _ = fmt.Fprintf(w, `[{"name": "INFRAANO","resource_pool_id": "resgroup-4047"},{"name": "MBIPMP_System","resource_pool_id": "resgroup-4049"}]`)
		}
	})
	mux.HandleFunc("/api/device/v1/instances", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438"}},"path":"/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438"}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/instances/initialization/tasks/deacca61-3162-4672-aac8-2d6bd2b69438", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/4b6fadb6-d44f-4a40-8e73-c2d0be7ca55f"}},"completed":"2024-03-07T06:57:00.298908Z","created":"2024-03-07T06:48:38.583246Z","failure_reason":"","id":"4b6fadb6-d44f-4a40-8e73-c2d0be7ca55f","name":"instance creation","payload":{"discovery":{"port":5443,"address":"10.146.194.171","device_user":"admin","device_password":"*****","management_user":"admin-cm","management_password":"*****"},"onboarding":{"mode":"STANDALONE","nodes":[{"hostname":"infraanovm01","password":"*****","username":"admin","managementAddress":"10.146.194.171"}],"siteInfo":{"dnsServers":["8.8.8.8"],"ntpServers":["0.us.pool.ntp.org"]},"platformType":"VE"},"instantiation":{"Request":{"F5osRequest":null,"VsphereRequest":{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","polling_delay":15,"provider_type":"vsphere","next_instances":[{"memory":16384,"cluster":"vSAN Cluster","num_cpus":8,"datastore":"vsanDatastore","datacenter":"mbip-7.0","sleep_time":"360s","resource_pool":"INFRAANO","mgmt_dns_server":"","vsphere_template":"BIG-IP-Next-20.1.0-2.279.0+0.0.75-VM-template","bigipnext_vm_name":"infraanovm01","mgmt_ipv4_address":"10.146.194.171/23","mgmt_ipv4_gateway":"10.146.195.254","mgmt_network_name":"VM-mgmt","bigipnext_vm_password":"*****","external_network_name":"LocalTestVLAN-115","internal_network_name":"LocalTestVLAN-114","vsphere_content_library":"CM-IOD","ha_data_plane_network_name":"LocalTestVLAN-116","ha_control_plane_network_name":""}]}},"BaseTask":{"id":"","payload":null,"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","provider_type":"vsphere"},"VsphereRequest":null}},"stage":"Discovery","state":"discoveryDone","status":"completed","task_type":"instance_creation","updated":"2024-03-07T06:57:00.298908Z"}`)
	})
	// mux.HandleFunc("/api/device/v1/instances", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438"}},"path":"/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438"}`)
	// })
	// mux.HandleFunc("/api/device/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/4b6fadb6-d44f-4a40-8e73-c2d0be7ca55f"}},"completed":"2024-03-07T06:57:00.298908Z","created":"2024-03-07T06:48:38.583246Z","failure_reason":"","id":"4b6fadb6-d44f-4a40-8e73-c2d0be7ca55f","name":"instance creation","payload":{"discovery":{"port":5443,"address":"10.146.194.171","device_user":"admin","device_password":"*****","management_user":"admin-cm","management_password":"*****"},"onboarding":{"mode":"STANDALONE","nodes":[{"hostname":"infraanovm01","password":"*****","username":"admin","managementAddress":"10.146.194.171"}],"siteInfo":{"dnsServers":["8.8.8.8"],"ntpServers":["0.us.pool.ntp.org"]},"platformType":"VE"},"instantiation":{"Request":{"F5osRequest":null,"VsphereRequest":{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","polling_delay":15,"provider_type":"vsphere","next_instances":[{"memory":16384,"cluster":"vSAN Cluster","num_cpus":8,"datastore":"vsanDatastore","datacenter":"mbip-7.0","sleep_time":"360s","resource_pool":"INFRAANO","mgmt_dns_server":"","vsphere_template":"BIG-IP-Next-20.1.0-2.279.0+0.0.75-VM-template","bigipnext_vm_name":"infraanovm01","mgmt_ipv4_address":"10.146.194.171/23","mgmt_ipv4_gateway":"10.146.195.254","mgmt_network_name":"VM-mgmt","bigipnext_vm_password":"*****","external_network_name":"LocalTestVLAN-115","internal_network_name":"LocalTestVLAN-114","vsphere_content_library":"CM-IOD","ha_data_plane_network_name":"LocalTestVLAN-116","ha_control_plane_network_name":""}]}},"BaseTask":{"id":"","payload":null,"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","provider_type":"vsphere"},"VsphereRequest":null}},"stage":"Discovery","state":"discoveryDone","status":"completed","task_type":"instance_creation","updated":"2024-03-07T06:57:00.298908Z"}`)
	// 	// {"_links":{"self":{"href":"/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438"}},"completed":"2024-03-06T18:54:20.948233Z","created":"2024-03-06T18:53:50.092597Z","failure_reason":"error setting up new vSphere SOAP client: ServerFaultCode: Cannot complete login due to an incorrect user name or password.","id":"deacca61-3162-4672-aac8-2d6bd2b69438","name":"instance creation","payload":{"discovery":{"port":5443,"address":"10.146.194.171","device_user":"admin","device_password":"*****","management_user":"admin-cm","management_password":"*****"},"onboarding":{"mode":"STANDALONE","nodes":[{"hostname":"infraanovm01","password":"*****","username":"admin","managementAddress":"10.146.194.171"}],"siteInfo":{"dnsServers":["8.8.8.8"],"ntpServers":["0.us.pool.ntp.org"]},"platformType":"VE"},"instantiation":{"Request":{"F5osRequest":null,"VsphereRequest":{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","polling_delay":15,"provider_type":"vsphere","next_instances":[{"memory":16384,"cluster":"vSAN Cluster","num_cpus":8,"datastore":"vsanDatastore","datacenter":"mbip-7.0","sleep_time":"360s","resource_pool":"INFRAANO","mgmt_dns_server":"","vsphere_template":"BIG-IP-Next-20.1.0-2.279.0+0.0.75-VM-template","bigipnext_vm_name":"infraanovm01","mgmt_ipv4_address":"10.146.194.171/23","mgmt_ipv4_gateway":"10.146.195.254","mgmt_network_name":"VM-mgmt","bigipnext_vm_password":"*****","external_network_name":"LocalTestVLAN-115","internal_network_name":"LocalTestVLAN-114","vsphere_content_library":"CM-IOD","ha_data_plane_network_name":"LocalTestVLAN-116","ha_control_plane_network_name":""}]}},"BaseTask":{"id":"","payload":null,"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","provider_type":"vsphere"},"VsphereRequest":null}},"stage":"Instantiation","state":"instantiateInstances","status":"completed","task_type":"instance_creation","updated":"2024-03-06T18:54:20.948233Z"}
	// })
	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testAccNextDeployVmwareResourceTC1Config,
				ExpectError: regexp.MustCompile(`Failed to Read Device Info, got error`),
			},
		},
	})
}

// func TestUnitNextDeployVmwareResourceTC4(t *testing.T) {
// 	testAccPreUnitCheck(t)
// 	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = fmt.Fprintf(w, `{
// 			"token": "eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.AbY1hUw8wHO8Vt1qxRd5xQj_21EQ1iaH6q9Z2XgRwQl98M7aCpyjiF2J16S4HrZ-",
// 			"tokenType": "Bearer",
// 			"expiresIn": 3600,
// 			"refreshToken": "ODA0MmQzZTctZTk1Mi00OTk1LWJmMjUtZWZmMjc1NDE3YzliOt4bKlRr6g7RdTtnBKhm2vzkgJeWqfvow68gyxTipleCq4AjR4nxZDBYKQaWyCWGeA",
// 			"refreshExpiresIn": 1209600
// 		}`)
// 	})
// 	mux.HandleFunc("/api/v1/spaces/default/providers", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = fmt.Fprintf(w, `[{"instances":[{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"},{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"},{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"}],"provider_host":"mbip-70-vcenter.pdsea.f5net.com","provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","provider_name":"myvsphere03","provider_type":"VSPHERE","provider_username":"r.chinthalapalli@f5.com","updated":"2023-12-12T11:53:58.614102Z"}]`)
// 	})
// 	mux.HandleFunc("/api/v1/system/infra/info", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = fmt.Fprintf(w, `{"app_version": "0.179.6","ha_status": "Not Running","num_of_nodes": 1,"version": "BIG-IP-Next-CentralManager-20.3.0-0.14.42"}`)
// 	})
// 	mux.HandleFunc("/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		t.Logf("RequestURI: %s", r.RequestURI)
// 		if r.RequestURI == "/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api?path=api/vcenter/datacenter" {
// 			_, _ = fmt.Fprintf(w, `[{"name": "mbip-7.0","datacenter_id": "datacenter-3"}]`)
// 		}
// 		if r.RequestURI == "/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api?path=api/vcenter/cluster?datacenters=datacenter-3" {
// 			_, _ = fmt.Fprintf(w, `[{"name": "vSAN Cluster","cluster_id": "domain-c8"}]`)
// 		}
// 		if r.RequestURI == "/api/v1/spaces/default/providers/vsphere/9945d9cd-9f13-438f-8b97-f0cb1745d32b/api?path=api/vcenter/resource-pool?clusters=domain-c8" {
// 			_, _ = fmt.Fprintf(w, `[{"name": "INFRAANO","resource_pool_id": "resgroup-4047"},{"name": "MBIPMP_System","resource_pool_id": "resgroup-4049"}]`)
// 		}
// 	})
// 	defer teardown()
// 	resource.Test(t, resource.TestCase{
// 		// PreCheck:                 func() { testAccPreUnitCheck(t) },
// 		IsUnitTest:               true,
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Read testing
// 			{
// 				Config:      testAccNextDeployVmwareResourceTC1Config,
// 				ExpectError: regexp.MustCompile(`Failed to Deploy Instance, got error:`),
// 				// Check: resource.ComposeAggregateTestCheckFunc(
// 				// 	resource.TestCheckResourceAttr("bigipnext_cm_deploy_vmware.vmware", "ntp_servers.0", "0.us.pool.ntp.org"),
// 				// 	resource.TestCheckResourceAttr("bigipnext_cm_deploy_vmware.vmware", "dns_servers.0", "8.8.8.8"),
// 				// ),
// 			},
// 		},
// 	})
// }

// func TestUnitNextDeployVmwareResourceTC4(t *testing.T) {
// 	testAccPreUnitCheck(t)
// 	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = fmt.Fprintf(w, `{
// 			"token": "eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.AbY1hUw8wHO8Vt1qxRd5xQj_21EQ1iaH6q9Z2XgRwQl98M7aCpyjiF2J16S4HrZ-",
// 			"tokenType": "Bearer",
// 			"expiresIn": 3600,
// 			"refreshToken": "ODA0MmQzZTctZTk1Mi00OTk1LWJmMjUtZWZmMjc1NDE3YzliOt4bKlRr6g7RdTtnBKhm2vzkgJeWqfvow68gyxTipleCq4AjR4nxZDBYKQaWyCWGeA",
// 			"refreshExpiresIn": 1209600
// 		}`)
// 	})
// 	mux.HandleFunc("/api/v1/spaces/default/providers", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = fmt.Fprintf(w, `[{"instances":[{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"},{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"},{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b"}],"provider_host":"mbip-70-vcenter.pdsea.f5net.com","provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","provider_name":"myvsphere03","provider_type":"VSPHERE","provider_username":"r.chinthalapalli@f5.com","updated":"2023-12-12T11:53:58.614102Z"}]`)
// 	})
// 	mux.HandleFunc("/api/device/v1/instances", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438"}},"path":"/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438"}`)
// 	})

// 	mux.HandleFunc("/api/device/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/instances/tasks/4b6fadb6-d44f-4a40-8e73-c2d0be7ca55f"}},"completed":"2024-03-07T06:57:00.298908Z","created":"2024-03-07T06:48:38.583246Z","failure_reason":"","id":"4b6fadb6-d44f-4a40-8e73-c2d0be7ca55f","name":"instance creation","payload":{"discovery":{"port":5443,"address":"10.146.194.171","device_user":"admin","device_password":"*****","management_user":"admin-cm","management_password":"*****"},"onboarding":{"mode":"STANDALONE","nodes":[{"hostname":"infraanovm01","password":"*****","username":"admin","managementAddress":"10.146.194.171"}],"siteInfo":{"dnsServers":["8.8.8.8"],"ntpServers":["0.us.pool.ntp.org"]},"platformType":"VE"},"instantiation":{"Request":{"F5osRequest":null,"VsphereRequest":{"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","polling_delay":15,"provider_type":"vsphere","next_instances":[{"memory":16384,"cluster":"vSAN Cluster","num_cpus":8,"datastore":"vsanDatastore","datacenter":"mbip-7.0","sleep_time":"360s","resource_pool":"INFRAANO","mgmt_dns_server":"","vsphere_template":"BIG-IP-Next-20.1.0-2.279.0+0.0.75-VM-template","bigipnext_vm_name":"infraanovm01","mgmt_ipv4_address":"10.146.194.171/23","mgmt_ipv4_gateway":"10.146.195.254","mgmt_network_name":"VM-mgmt","bigipnext_vm_password":"*****","external_network_name":"LocalTestVLAN-115","internal_network_name":"LocalTestVLAN-114","vsphere_content_library":"CM-IOD","ha_data_plane_network_name":"LocalTestVLAN-116","ha_control_plane_network_name":""}]}},"BaseTask":{"id":"","payload":null,"provider_id":"9945d9cd-9f13-438f-8b97-f0cb1745d32b","provider_type":"vsphere"},"VsphereRequest":null}},"stage":"Discovery","state":"discoveryDone","status":"completed","task_type":"instance_creation","updated":"2024-03-07T06:57:00.298908Z"}`)
// 	})
// 	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = fmt.Fprintf(w, `{"_embedded":{"devices":[{"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+infraanovm01/37021437-2b5f-4e44-9f5e-6ea9838c5f7e"}},"address":"10.146.194.171","certificate_validated":"2024-03-07T06:56:51.915056Z","certificate_validity":false,"hostname":"infraanovm01","id":"37021437-2b5f-4e44-9f5e-6ea9838c5f7e","mode":"STANDALONE","platform_name":"VMware","platform_type":"VE","port":5443,"version":"20.1.0-2.279.0+0.0.75"}]},"_links":{"self":{"href":"/v1/inventory?filter=hostname+eq+infraanovm01"}},"count":1,"total":1}`)
// 	})
// 	mux.HandleFunc("/api/device/v1/inventory/37021437-2b5f-4e44-9f5e-6ea9838c5f7e", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusBadRequest)
// 		_, _ = fmt.Fprintf(w, ``)
// 	})
// 	defer teardown()
// 	resource.Test(t, resource.TestCase{
// 		// PreCheck:                 func() { testAccPreUnitCheck(t) },
// 		IsUnitTest:               true,
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Read testing
// 			{
// 				Config:      testAccNextDeployVmwareResourceTC1Config,
// 				ExpectError: regexp.MustCompile(`Unable to Delete Instance, got error:`),
// 			},
// 		},
// 	})
// }

const testAccNextDeployVmwareResourceTC1Config = `
resource "bigipnext_cm_deploy_vmware" "vmware" {
  vsphere_provider = {
    provider_name      = "myvsphere03"
    content_library    = "CM-IOD"
    cluster_name       = "vSAN Cluster"
    datacenter_name    = "mbip-7.0"
    datastore_name     = "vsanDatastore"
    resource_pool_name = "INFRAANO"
    vm_template_name   = "BIG-IP-Next-20.1.0-2.279.0+0.0.75-VM-template"
  }
  instance = {
    instance_hostname          = "infraanovm01"
    mgmt_address               = "10.146.194.171"
    mgmt_prefix                = 23
    mgmt_gateway               = "10.146.195.254"
    mgmt_network_name          = "VM-mgmt"
    mgmt_user                  = "admin-cm"
    mgmt_password              = "F5Twist@123"
    external_network_name      = "LocalTestVLAN-115"
    internal_network_name      = "LocalTestVLAN-114"
    ha_data_plane_network_name = "LocalTestVLAN-116"
  }
  l1_networks = [{
    name          = "demonetwork1"
    vlans = [{
      vlan_tag = 115
      vlan_name = "vlan-115"
      self_ips=["10.101.10.10/24","10.101.10.11/24"]}]
  }]
  ntp_servers = ["0.us.pool.ntp.org"]
  dns_servers = ["8.8.8.8"]
  timeout     = 1200
}
`

const testAccNextDeployVmwareResourceTC2Config = `
resource "bigipnext_cm_deploy_vmware" "vmware" {
  vsphere_provider = {
    provider_name      = "myvsphere03"
    content_library    = "CM-IOD"
    cluster_name       = "vSAN Cluster"
    datacenter_name    = "mbip-7.0"
    datastore_name     = "vsanDatastore"
    resource_pool_name = "INFRAANO"
    vm_template_name   = "BIG-IP-Next-20.1.0-2.279.0+0.0.75-VM-template"
  }
  instance = {
    instance_hostname          = "infraanovm01"
    mgmt_address               = "10.100.100.171"
    mgmt_prefix                = 23
    mgmt_gateway               = "10.100.100.254"
    mgmt_network_name          = "VM-mgmt"
    mgmt_user                  = "admin-cm"
    mgmt_password              = "F5Twist@123"
    external_network_name      = "LocalTestVLAN-115"
    internal_network_name      = "LocalTestVLAN-114"
    ha_data_plane_network_name = "LocalTestVLAN-116"
  }
  ntp_servers = ["0.us.pool.ntp.org"]
  dns_servers = ["8.8.8.8"]
  timeout     = 1200
}
`
