package provider

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNextCMWafPolicyImportCreateTC1Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMWafPolicyImportResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_waf_policy_import.sample", "name", "new_waf_policy"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_policy_import.sample", "description", "new_waf_policy desc"),
				),
			},
		},
	})
}

func TestAccNextCMWafPolicyImportCreateTC2Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMWafPolicyImportResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_waf_policy_import.sample", "name", "new_waf_policy"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_policy_import.sample", "description", "new_waf_policy desc"),
				),
			},
			{
				Config: testAccNextCMWafPolicyImportUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_waf_policy_import.sample", "name", "new_waf_policy"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_policy_import.sample", "description", "new_waf_policy desc updated"),
				),
			},
		},
	})
}

func TestUnitNextCMWafPolicyImportCreateResourceTC1(t *testing.T) {
	testAccPreUnitCheck(t)
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

	// Post call
	mux.HandleFunc("/api/waf/v1/tasks/policy-import", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
				"_links": {
					"self": {
						"href": "/api/waf/v1/tasks/policy-import/1a4453fe-b37a-4212-a813-a3d2f789dad1"
					}
				},
				"path": "/v1/tasks/policy-import/1a4453fe-b37a-4212-a813-a3d2f789dad1"
			}`)
	})

	// Check Status
	mux.HandleFunc("/api/waf/v1/tasks/policy-import/1a4453fe-b37a-4212-a813-a3d2f789dad1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
				"self": {
					"href": "/api/waf/v1/tasks/policy-import/1a4453fe-b37a-4212-a813-a3d2f789dad1"
				}
			},
			"completed": "2024-07-17T11:59:25.7393Z",
			"created": "2024-07-17T11:59:21.058179Z",
			"failure_reason": "",
			"id": "1a4453fe-b37a-4212-a813-a3d2f789dad1",
			"policy_id": "1a4453fe-b37a-4212-a813-a3d2f789dad1",
			"policy_name": "new_waf_policy",
			"state": "updatingTaskDataTable",
			"status": "completed",
			"warnings": []
		}`)
	})

	// Get/Delete call
	mux.HandleFunc("/api/v1/spaces/default/security/waf-policies/1a4453fe-b37a-4212-a813-a3d2f789dad1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, ``)
		} else {
			if getCount < 2 {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, "%s", loadFixtureString("./fixtures/getWaf.json"))
			} else {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, "%s", loadFixtureString("./fixtures/getWafUpdated.json"))
			}

			getCount++

		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMWafPolicyImportResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: testAccNextCMWafPolicyImportUpdateResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

var dir, _ = os.Getwd()

var testAccNextCMWafPolicyImportResourceConfig = `
resource "bigipnext_cm_waf_policy_import" "sample" {
	name        = "new_waf_policy"
	description = "new_waf_policy desc"
	file_path   = "` + dir + `/../../configs/testwaf.json"
	file_md5    = md5(file("` + dir + `/../../configs/testwaf.json"))
	}`

var testAccNextCMWafPolicyImportUpdateResourceConfig = `
resource "bigipnext_cm_waf_policy_import" "sample" {
	name        = "new_waf_policy"
	description = "new_waf_policy desc updated"
	file_path   = "` + dir + `/../../configs/testwaf.json"
	file_md5    =  md5(file("` + dir + `/../../configs/testwaf.json"))
  }
`
