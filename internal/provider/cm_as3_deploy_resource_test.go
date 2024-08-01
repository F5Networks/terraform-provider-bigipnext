package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitNextCMAS3CreateResourceTC1(t *testing.T) {
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
	mux.HandleFunc("/api/v1/spaces/default/appsvcs/documents", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"Message":"Application service created successfully","_links":{"self":{"href":"/api/v1/spaces/default/appsvcs/documents/9a807f7f-f91c-4fb5-abee-a708dc44a7b8"}},"id":"9a807f7f-f91c-4fb5-abee-a708dc44a7b8"}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/appsvcs/documents/9a807f7f-f91c-4fb5-abee-a708dc44a7b8/deployments", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"Message":"Deployment task created successfully","_links":{"self":{"href":"/api/v1/spaces/default/appsvcs/documents/9a807f7f-f91c-4fb5-abee-a708dc44a7b8/deployments/dddd9907-2e5e-413b-9087-8665c872a001"}},"id":"dddd9907-2e5e-413b-9087-8665c872a001","task_id":"060aaca6-aedc-43b0-b803-573bca61b25f"}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/appsvcs/documents/9a807f7f-f91c-4fb5-abee-a708dc44a7b8/deployments/dddd9907-2e5e-413b-9087-8665c872a001", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/appsvcs/documents/9a807f7f-f91c-4fb5-abee-a708dc44a7b8/deployments/dddd9907-2e5e-413b-9087-8665c872a001"}},"id":"dddd9907-2e5e-413b-9087-8665c872a001","instance_id":"3c815d06-0fe4-407d-b5b4-a380f550565d","records":[{"id":"881bc077-8cc9-4f16-97aa-bf7703de88b1","task_id":"060aaca6-aedc-43b0-b803-573bca61b25f","start_time":"2024-07-17T05:34:46.22289Z","end_time":"0001-01-01T00:00:00Z","status":"completed"}],"records_count":{"total":1},"request":{"DemoTenant3":{"TestApp33":{"Pool3":{"class":"Pool","loadBalancingMode":"round-robin","members":[{"serverAddresses":["15.6.17.10"],"servicePort":80}]},"class":"Application","serviceMain":{"class":"Service_HTTP","pool":"Pool3","snat":"auto","virtualAddresses":["15.6.17.9"],"virtualPort":80},"template":"http"},"class":"Tenant"},"class":"ADC","schemaVersion":"3.50.0"},"response":{"_links":{"self":"/mgmt/shared/appsvcs/task/b220d337-f713-4b9f-b620-297e5c69fc88","taskStatus":"/files/09289779-c994-4c74-842f-d9688749b889"},"created":"2024-07-17T05:34:53.407Z","declaration":{"DemoTenant3":{"TestApp33":{"Pool3":{"class":"Pool","loadBalancingMode":"round-robin","members":[{"serverAddresses":["15.6.17.10"],"servicePort":80}]},"class":"Application","serviceMain":{"class":"Service_HTTP","pool":"Pool3","snat":"auto","virtualAddresses":["15.6.17.9"],"virtualPort":80},"template":"http"},"class":"Tenant"},"class":"ADC","schemaVersion":"3.50.0"},"id":"b220d337-f713-4b9f-b620-297e5c69fc88","results":[{"code":202,"host":"demovm01-93748901","message":"in progress","runTime":0,"tenant":"DemoTenant3"}],"selfLink":"/mgmt/shared/appsvcs/task/b220d337-f713-4b9f-b620-297e5c69fc88"},"tenant_name":"DemoTenant3","type":"AS3"}`)
	})

	// Delete call
	mux.HandleFunc("/api/v1/spaces/default/appsvcs/documents/9a807f7f-f91c-4fb5-abee-a708dc44a7b8", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/appsvcs/documents/9a807f7f-f91c-4fb5-abee-a708dc44a7b8"}},"deployments":[{"Message":"Delete Deployment task created successfully","id":"dddd9907-2e5e-413b-9087-8665c872a001","task_id":"8e2773bc-1538-4ceb-8105-732f2a8bd1dd"}],"id":"9a807f7f-f91c-4fb5-abee-a708dc44a7b8","message":"The application delete has been submitted successfully"}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMAS3ResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: testAccNextCMAS3ResourceUpdateConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// PUT:
// {"_links":{"self":{"href":"/api/v1/spaces/default/appsvcs/documents/9a807f7f-f91c-4fb5-abee-a708dc44a7b8"}},"deployments":[{"Message":"Update deployment task created","id":"dddd9907-2e5e-413b-9087-8665c872a001","task_id":"a7a6c23a-d96b-4eac-a62b-43decd510059"}],"id":"9a807f7f-f91c-4fb5-abee-a708dc44a7b8","message":"Application service updated successfully"}

// /api/v1/spaces/default/appsvcs/documents/9a807f7f-f91c-4fb5-abee-a708dc44a7b8

// /api/v1/spaces/default/appsvcs/documents/9a807f7f-f91c-4fb5-abee-a708dc44a7b8/deployments/dddd9907-2e5e-413b-9087-8665c872a001

//
// {"_links":{"self":{"href":"/api/v1/spaces/default/appsvcs/documents/9a807f7f-f91c-4fb5-abee-a708dc44a7b8/deployments/dddd9907-2e5e-413b-9087-8665c872a001"}},"id":"dddd9907-2e5e-413b-9087-8665c872a001","instance_id":"3c815d06-0fe4-407d-b5b4-a380f550565d","last_app_modified":"2024-07-17T05:43:01.054333Z","last_successful_deploy_time":"2024-07-17T05:34:48.828081Z","records":[{"id":"f7d711c3-ce3f-4db9-926b-6a95c637bcf0","task_id":"a7a6c23a-d96b-4eac-a62b-43decd510059","start_time":"2024-07-17T05:43:01.087674Z","end_time":"0001-01-01T00:00:00Z","status":"running"},{"id":"881bc077-8cc9-4f16-97aa-bf7703de88b1","task_id":"060aaca6-aedc-43b0-b803-573bca61b25f","start_time":"2024-07-17T05:34:46.22289Z","end_time":"0001-01-01T00:00:00Z","status":"completed"}],"records_count":{"total":2,"completed":1},"request":{"DemoTenant3":{"TestApp33":{"Pool3":{"class":"Pool","loadBalancingMode":"round-robin","members":[{"serverAddresses":["15.6.17.11"],"servicePort":80}]},"class":"Application","serviceMain":{"class":"Service_HTTP","pool":"Pool3","snat":"auto","virtualAddresses":["15.6.17.9"],"virtualPort":80},"template":"http"},"class":"Tenant"},"class":"ADC","schemaVersion":"3.50.0"},"response":{"_links":{"self":"/mgmt/shared/appsvcs/task/1affa319-1071-4e20-930d-fc809839997a","taskStatus":"/files/eab0ad9d-9be1-4a68-98bf-054f9a1f5e45"},"created":"2024-07-17T05:43:08.230Z","declaration":{"DemoTenant3":{"TestApp33":{"Pool3":{"class":"Pool","loadBalancingMode":"round-robin","members":[{"serverAddresses":["15.6.17.11"],"servicePort":80}]},"class":"Application","serviceMain":{"class":"Service_HTTP","pool":"Pool3","snat":"auto","virtualAddresses":["15.6.17.9"],"virtualPort":80},"template":"http"},"class":"Tenant"},"class":"ADC","schemaVersion":"3.50.0"},"id":"1affa319-1071-4e20-930d-fc809839997a","results":[{"code":202,"host":"demovm01-93748901","message":"in progress","runTime":0,"tenant":"DemoTenant3"}],"selfLink":"/mgmt/shared/appsvcs/task/1affa319-1071-4e20-930d-fc809839997a"},"tenant_name":"DemoTenant3","type":"AS3"}

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
resource "bigipnext_cm_as3_deploy" "test3" {
  target_address = "10.10.10.10"
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
resource "bigipnext_cm_as3_deploy" "test3" {
  target_address = "10.10.10.10"
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
                        "10.0.11.10"
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
