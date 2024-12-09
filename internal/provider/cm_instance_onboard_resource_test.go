package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCMOnboardResourceTC1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testUnitCMOnboardResourceTC1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "management_address", "10.218.135.67"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "dns_servers.*", "2.2.2.4"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "ntp_servers.*", "4.pool.com")),
			},
			{
				Config: testUnitCMOnboardResourceTC1_2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "management_address", "10.218.135.67"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "dns_servers.*", "2.2.2.4"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "ntp_servers.*", "4.pool.com"),
					resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "l1_networks.0.name", "l1network4"),
				),
			},
			{
				Config: testUnitCMOnboardResourceTC1_3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "management_address", "10.218.135.67"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "dns_servers.*", "2.2.2.4"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "ntp_servers.*", "4.pool.com"),
					resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "l1_networks.0.name", "l1network4"),
					resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "l1_networks.0.vlans.0.name", "vlan104"),
				),
			},
			{
				Config: testUnitCMOnboardResourceTC2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "management_address", "10.218.135.67"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "dns_servers.*", "2.2.2.5"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "ntp_servers.*", "5.pool.com"),
					resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "l1_networks.0.name", "l1network4"),
					resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "l1_networks.0.vlans.0.name", "vlan104"),
				),
			},
		},
	})
}

func TestUnitCMOnboardResourceTC1(t *testing.T) {
	testAccPreUnitCheck(t)
	var count = 0
	var getCount = 0
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
		_, _ = fmt.Fprintf(w, `{
					"_embedded": {
						"devices": [
							{
								"address": "10.218.135.67",
								"certificate_validated": "2024-07-01T08:55:37.817999Z",
								"certificate_validation_error": "tls: failed to verify certificate: x509: certificate signed by unknown authority (possibly because of \"x509: invalid signature: parent certificate cannot sign this kind of certificate\" while trying to verify candidate authority certificate \"localhost\")",
								"certificate_validity": false,
								"hostname": "big-ip-next",
								"id": "6a19b54a-723a-4dea-85fb-46c75c6acfd2",
								"mode": "STANDALONE",
								"platform_name": "KVM",
								"platform_type": "VE",
								"port": 5443,
								"short_id": "S-Rg6ftH",
								"version": "20.3.0-2.489.0"
							}
						]
					},
					"count": 1,
					"total": 1
				}`)
	})
	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
				"_embedded": {
					"devices": [
						{
							"_links": {
								"self": {
									"href": "/v1/inventory"
								}
							},
							"address": "10.218.135.67",
							"analytics_node_address": "10.218.135.237",
							"certificate_validated": "2024-11-12T11:53:34.309759Z",
							"certificate_validation_error": "tls: failed to verify certificate: x509: certificate signed by unknown authority (possibly because of \"x509: invalid signature: parent certificate cannot sign this kind of certificate\" while trying to verify candidate authority certificate \"localhost\")",
							"certificate_validity": false,
							"hostname": "big-ip-next",
							"id": "6a19b54a-723a-4dea-85fb-46c75c6acfd2",
							"mode": "STANDALONE",
							"platform_name": "KVM",
							"platform_type": "VE",
							"port": 5443,
							"short_id": "Y2B_FmIG",
							"version": "20.4.0-2.851.0"
						}
					]
				},
				"_links": {
					"self": {
						"href": "/v1/inventory"
					}
				},
				"count": 1,
				"total": 1
			}`)
	})

	mux.HandleFunc("/api/v1/spaces/default/instances/initialization/6a19b54a-723a-4dea-85fb-46c75c6acfd2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			if count == 0 {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{
				"_links": {
					"self": {
						"href": "/api/v1/spaces/default/instances/initialization/tasks/2bf5d785-d0cd-4b02-8810-fbcc070aef15"
					}
				},
				"path": "/api/v1/spaces/default/instances/initialization/tasks/2bf5d785-d0cd-4b02-8810-fbcc070aef15"
			}`)
			} else if count == 1 {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{
					"_links": {
						"self": {
							"href": "/api/v1/spaces/default/instances/initialization/tasks/1400a2ea-e5ea-4cfe-a220-89230e318620"
						}
					},
					"path": "/api/v1/spaces/default/instances/initialization/tasks/1400a2ea-e5ea-4cfe-a220-89230e318620"
				}`)
			}
			count++
		} else if r.Method == http.MethodGet {
			if getCount < 2 {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{
					"_links": {
						"self": {
							"href": "/api/v1/spaces/default/instances/initialization/6a19b54a-723a-4dea-85fb-46c75c6acfd2"
						}
					},
					"id": "1e5ae167-f2bb-4c97-a931-ae06f739515a",
					"instance_id": "6a19b54a-723a-4dea-85fb-46c75c6acfd2",
					"onboarding_manifest_id": "c4d28a61-059c-46a7-b213-8b625aa64ce0",
					"parameters": {
						"hostname": "big-ip-next",
						"l1Networks": [
							{
								"name": "l1network4",
								"vlans": [
									{
										"tag": 104,
										"name": "vlan104",
										"selfIps": [
											{
												"address": "20.20.20.24/24",
												"deviceName": "device4"
											}
										]
									}
								],
								"l1Link": {
									"name": "1.1",
									"linkType": "Interface"
								}
							}
						],
						"dns_servers": [
							"2.2.2.5"
						],
						"ntp_servers": [
							"5.pool.com"
						],
						"default_gateway": "",
						"management_address": "10.218.135.67",
						"vSphere_properties": [
							{
								"memory": 16384,
								"num_cpus": 8,
								"cluster_name": "none",
								"datastore_name": "none",
								"datacenter_name": "none",
								"vm_template_name": "none",
								"resource_pool_name": "none",
								"vsphere_content_library": "none"
							}
						],
						"instantiation_provider": [
							{
								"id": "00000000-0000-0000-0000-000000000000",
								"type": "vsphere"
							}
						],
						"management_network_width": 24,
						"management_credentials_password": "*****",
						"management_credentials_username": "admin-cm",
						"vsphere_network_adapter_settings": [
							{
								"mgmt_network_name": "mgmt-placeholder",
								"external_network_name": "external-placeholder",
								"internal_network_name": "internal-placeholder",
								"ha_data_plane_network_name": "dataplane-placeholder",
								"ha_control_plane_network_name": "controlplane-placeholder"
							}
						]
					},
					"template_name": "default-standalone-ve"
				}`)
			} else if getCount >= 2 {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{
					"_links": {
						"self": {
							"href": "/api/v1/spaces/default/instances/initialization/6a19b54a-723a-4dea-85fb-46c75c6acfd2"
						}
					},
					"id": "1e5ae167-f2bb-4c97-a931-ae06f739515a",
					"instance_id": "6a19b54a-723a-4dea-85fb-46c75c6acfd2",
					"onboarding_manifest_id": "f042507a-8f49-4fe7-9a68-773621ec2b1f",
					"parameters": {
						"hostname": "big-ip-next",
						"l1Networks": [
							{
								"name": "l1network4",
								"vlans": [
									{
										"tag": 105,
										"name": "vlan105",
										"selfIps": [
											{
												"address": "20.20.20.25/24",
												"deviceName": "device4"
											}
										]
									}
								],
								"l1Link": {
									"name": "1.1",
									"linkType": "Interface"
								}
							}
						],
						"dns_servers": [
							"2.2.2.6"
						],
						"ntp_servers": [
							"6.pool.com"
						],
						"default_gateway": "",
						"management_address": "10.218.135.67",
						"vSphere_properties": [
							{
								"memory": 16384,
								"num_cpus": 8,
								"cluster_name": "none",
								"datastore_name": "none",
								"datacenter_name": "none",
								"vm_template_name": "none",
								"resource_pool_name": "none",
								"vsphere_content_library": "none"
							}
						],
						"instantiation_provider": [
							{
								"id": "00000000-0000-0000-0000-000000000000",
								"type": "vsphere"
							}
						],
						"management_network_width": 24,
						"management_credentials_password": "*****",
						"management_credentials_username": "admin-cm",
						"vsphere_network_adapter_settings": [
							{
								"mgmt_network_name": "mgmt-placeholder",
								"external_network_name": "external-placeholder",
								"internal_network_name": "internal-placeholder",
								"ha_data_plane_network_name": "dataplane-placeholder",
								"ha_control_plane_network_name": "controlplane-placeholder"
							}
						]
					},
					"template_name": "default-standalone-ve"
				}`)
			}
			getCount++
		}
	})

	mux.HandleFunc("/api/v1/spaces/default/instances/initialization/tasks/2bf5d785-d0cd-4b02-8810-fbcc070aef15", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{
				"_links": {
					"self": {
						"href": "/v1/instances/tasks/2bf5d785-d0cd-4b02-8810-fbcc070aef15"
					}
				},
				"completed": "2024-07-02T06:00:35.776456Z",
				"created": "2024-07-02T06:00:35.572952Z",
				"id": "2bf5d785-d0cd-4b02-8810-fbcc070aef15",
				"instance_id": "6a19b54a-723a-4dea-85fb-46c75c6acfd2",
				"name": "instance edit",
				"payload": {
					"is_edit": true,
					"instance_id": "6a19b54a-723a-4dea-85fb-46c75c6acfd2",
					"template_name": "default-standalone-ve"
				},
				"state": "ProcessTemplate",
				"status": "completed",
				"task_type": "instance_edit",
				"updated": "2024-07-02T06:00:35.776456Z"
			}`)
		}
	})

	mux.HandleFunc("/api/v1/spaces/default/instances/initialization/tasks/1400a2ea-e5ea-4cfe-a220-89230e318620", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{
				"_links": {
					"self": {
						"href": "/v1/instances/tasks/1400a2ea-e5ea-4cfe-a220-89230e318620"
					}
				},
				"created": "2024-07-02T09:43:57.563675Z",
				"id": "1400a2ea-e5ea-4cfe-a220-89230e318620",
				"instance_id": "6a19b54a-723a-4dea-85fb-46c75c6acfd2",
				"name": "instance edit",
				"payload": {
					"discovery": {
						"port": 5443,
						"address": "10.218.135.67",
						"device_user": "admin",
						"device_password": "*****",
						"management_user": "admin-cm",
						"management_password": "*****"
					},
					"onboarding": {
						"mode": "STANDALONE",
						"nodes": [
							{
								"hostname": "big-ip-next",
								"password": "*****",
								"username": "admin",
								"managementAddress": "10.218.135.67"
							}
						],
						"siteInfo": {
							"dnsServers": [
								"2.2.2.6"
							],
							"ntpServers": [
								"6.pool.com"
							]
						},
						"l1Networks": [
							{
								"name": "l1network4",
								"vlans": [
									{
										"tag": 105,
										"name": "vlan105",
										"selfIps": [
											{
												"address": "20.20.20.25/24",
												"deviceName": "device4"
											}
										],
										"defaultVrf": true
									}
								],
								"l1Link": {
									"name": "1.1",
									"linkType": "Interface"
								}
							}
						],
						"platformType": "VE"
					},
					"instantiation": {
						"Request": {
							"F5osRequest": null,
							"VsphereRequest": {
								"provider_id": "00000000-0000-0000-0000-000000000000",
								"polling_delay": 15,
								"provider_type": "vsphere",
								"next_instances": [
									{
										"memory": 16384,
										"cluster": "none",
										"num_cpus": 8,
										"datastore": "none",
										"datacenter": "none",
										"sleep_time": "360s",
										"resource_pool": "none",
										"mgmt_dns_server": "",
										"vsphere_template": "none",
										"bigipnext_vm_name": "big-ip-next",
										"mgmt_ipv4_address": "10.218.135.67/24",
										"mgmt_ipv4_gateway": "",
										"mgmt_network_name": "mgmt-placeholder",
										"bigipnext_vm_password": "*****",
										"external_network_name": "external-placeholder",
										"internal_network_name": "internal-placeholder",
										"vsphere_content_library": "none",
										"ha_data_plane_network_name": "dataplane-placeholder",
										"ha_control_plane_network_name": "controlplane-placeholder"
									}
								]
							}
						},
						"BaseTask": {
							"id": "",
							"payload": null,
							"provider_id": "00000000-0000-0000-0000-000000000000",
							"provider_type": "vsphere"
						},
						"VsphereRequest": null
					}
				},
				"stage": "Onboarding",
				"state": "onboardHandleStandaloneVE",
				"status": "completed",
				"task_type": "instance_edit",
				"updated": "2024-07-02T09:43:57.708384Z"
			}`)
		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testUnitCMOnboardResourceTC2,
				Check: resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "management_address", "10.218.135.67"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "dns_servers.*", "2.2.2.5"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "ntp_servers.*", "5.pool.com"),
					resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "l1_networks.0.vlans.0.name", "vlan104"),
				),
			},
			{
				Config: testUnitCMOnboardResourceTC2_1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "dns_servers.*", "2.2.2.6"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_instance_onboard.test", "ntp_servers.*", "6.pool.com"),
					resource.TestCheckResourceAttr("bigipnext_cm_instance_onboard.test", "l1_networks.0.vlans.0.name", "vlan105"),
				),
			},
		},
	})
}

const testUnitCMOnboardResourceTC1 = `
resource "bigipnext_cm_instance_onboard" "test" {
	dns_servers        = ["2.2.2.4"]
	ntp_servers        = ["4.pool.com"]
	management_address = "10.218.135.67"
	timeout = 300
  }
`

const testUnitCMOnboardResourceTC1_2 = `
resource "bigipnext_cm_instance_onboard" "test" {
	dns_servers        = ["2.2.2.4"]
	ntp_servers        = ["4.pool.com"]
	management_address = "10.218.135.67"
	l1_networks = [{
		name = "l1network4"
		l1_link = {
		  name = "1.1"
		  link_type : "Interface"
		}
	  }]
	timeout = 300
  }
`

const testUnitCMOnboardResourceTC1_3 = `
resource "bigipnext_cm_instance_onboard" "test" {
	dns_servers        = ["2.2.2.4"]
	ntp_servers        = ["4.pool.com"]
	management_address = "10.218.135.67"
	l1_networks = [{
		name = "l1network4"
		vlans = [
		{
		  tag  = 104
		  name = "vlan104"
		}
	  ]
		l1_link = {
		  name = "1.1"
		  link_type : "Interface"
		}
	  }]
	timeout = 300
  }
`

const testUnitCMOnboardResourceTC2 = `
resource "bigipnext_cm_instance_onboard" "test" {
	dns_servers        = ["2.2.2.5"]
	ntp_servers        = ["5.pool.com"]
	management_address = "10.218.135.67"
	l1_networks = [{
	  name = "l1network4"
	  vlans = [
		{
		  tag  = 104
		  name = "vlan104"
		  self_ips = [
			{
			  address    = "20.20.20.24/24"
			  device_name = "device4"
			}
		  ]
		}
	  ]
	  l1_link = {
		name = "1.1"
		link_type : "Interface"
	  }
	}]
	timeout = 300
  }`

const testUnitCMOnboardResourceTC2_1 = `
resource "bigipnext_cm_instance_onboard" "test" {
	dns_servers        = ["2.2.2.6"]
	ntp_servers        = ["6.pool.com"]
	management_address = "10.218.135.67"
	l1_networks = [{
	name = "l1network4"
	vlans = [
		{
		tag  = 105
		name = "vlan105"
		self_ips = [
			{
			address    = "20.20.20.25/24"
			device_name = "device4"
			}
		]
		}
	]
	l1_link = {
		name = "1.1"
		link_type : "Interface"
	}
	}]
	timeout = 300
}
`
