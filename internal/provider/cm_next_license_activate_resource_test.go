package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitCMNextLicenseActivateResourceTC1(t *testing.T) {
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
	mux.HandleFunc("/api/v1/spaces/default/instances/license/activate", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"422b0cec-03b9-4499-a26e-c88f57869637": {
			"_links": {
				"self": {
					"href": "/license-task/41e49d68-d146-4e16-b286-7b57731fe14d"
				}
			},
			"accepted": true,
			"deviceId": "422b0cec-03b9-4499-a26e-c88f57869637",
			"reason": "",
			"taskId": "41e49d68-d146-4e16-b286-7b57731fe14d"}}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/license/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
		"41e49d68-d146-4e16-b286-7b57731fe14d": {
			"_links": {
				"self": {
					"href": "/license-task/41e49d68-d146-4e16-b286-7b57731fe14d"
				}
			},
			"taskExecutionStatus": {
				"created": "2024-06-27T15:59:25.928845Z",
				"failureReason": "",
				"status": "completed",
				"subStatus": "TERMINATE_ACK_VERIFICATION_COMPLETE",
				"taskType": "activate"
			}}}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/instances/license/license-info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"3c815d06-0fe4-407d-b5b4-a380f550565d":{"_links":{"self":{"href":"/license/3c815d06-0fe4-407d-b5b4-a380f550565d/status"}},"deviceLicenseStatus":{"enabledFeatures":"[{\"entitledFeatureId\":\"\",\"featureFlag\":\"bigip_ltm_module\",\"featurePermitted\":1,\"featureRemain\":0,\"featureUnlimited\":false,\"featureUsed\":0,\"featureValueType\":\"boolean\",\"uomCode\":\"\",\"uomTerm\":\"\",\"uomTermStart\":0},{\"entitledFeatureId\":\"\",\"featureFlag\":\"bigip_waf_module\",\"featurePermitted\":1,\"featureRemain\":0,\"featureUnlimited\":false,\"featureUsed\":0,\"featureValueType\":\"boolean\",\"uomCode\":\"\",\"uomTerm\":\"\",\"uomTermStart\":0}]","expiryDate":"2024-12-07T00:00:00Z","licenseStatus":"Active","licenseSubStatus":"ACK_VERIFICATION_COMPLETE","licenseToken":{"_links":{"self":{"href":"/token/dd50dae9-a9a7-49f0-ac30-5d6d44efbedb"}},"tokenId":"dd50dae9-a9a7-49f0-ac30-5d6d44efbedb","tokenName":"ravitoken"},"subscriptionSubType":"internal","subscriptionType":"paid"}}}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/instances/license/deactivate", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
		"3c815d06-0fe4-407d-b5b4-a380f550565d": {
			"accepted": true,
			"reason": "failed to create task",
			"taskId": "d290f1ee-6c54-4b01-90e6-d701748f0982",
			"deviceId": "3c815d06-0fe4-407d-b5b4-a380f550565d",
			"_links": {
			"self": {
				"href": "/license-task/a01eeeaa-7cb1-4ce1-9d7e-6ea20e7693bb"
			}}}}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testUnitCMNextLicenseActivateResourceTC1,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: testUnitCMNextLicenseActivateResourceTC2,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// TC1: Test case when JWT token is not verifiable
func TestAccCMNextLicenseActivateResourceTC1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCMNextLicenseActivateResourceTC1,
				// ExpectError: regexp.MustCompile(`Failed to Create Jwt token`),
			},
		},
	})
}

const testUnitCMNextLicenseActivateResourceTC1 = `
resource "bigipnext_cm_activate_instance_license" "tokenadd" {
	instances = [{
		instance_address = "10.10.10.10"
		jwt_id="eyJhbGciOi"
	}]
}
`

const testUnitCMNextLicenseActivateResourceTC2 = `
resource "bigipnext_cm_activate_instance_license" "tokenadd" {
	instances = [{
		instance_address = "10.10.10.10"
		jwt_id="eyJhbGciOi"
	}]
}
`

// jwtID := "8a3dc22e-dd51-4a5c-bb3a-cb239b904326"

const testAccCMNextLicenseActivateResourceTC1 = `
resource "bigipnext_cm_add_jwt_token" "tokenadd2" {
  token_name = "ravitoken"
  jwt_token  = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6InYxIiwiamt1IjoiaHR0cHM6Ly9wcm9kdWN0LmFwaXMuZjUuY29tL2VlL3YxL2tleXMvandrcyJ9.eyJzdWIiOiJGTkktZGM5MTFhYjctNDlhYi00ZWUzLTkwYjItNTNmNmVmOGY5ZTc0IiwiaWF0IjoxNzA1NjIyMDY5LCJpc3MiOiJGNSBJbmMuIiwiYXVkIjoidXJuOmY1OnRlZW0iLCJqdGkiOiJlYjIyZDhjMS1iNjVjLTExZWUtODQ1MC1mM2M0OTkzMmJmOTIiLCJmNV9vcmRlcl90eXBlIjoicGFpZCIsImY1X29yZGVyX3N1YnR5cGUiOiJpbnRlcm5hbCJ9.c6V5gAVGd-Pj62krRj8740bLq0_YyuRUvtKmat-oEJdiQn10GFbHNqBH8l3x0stcdE0UldrGszVQI3CmukDceRYi1XiTQpB69EubbOpx8Pe4qc6ht7kErmkDvsLlpy6ALYhdl8j2m5_npy3HvmoDnE2jjzWQkiQeFZjdxT4Gqc05LmRsO4_RnMOkvZFECfbRU6dEhtmP1es7L2FxJaJyJ8JEL2mz9kC8XtwaoW_jS_lxq_l5brDCnXJuFmLF882xyCReCT62FvIb4P4vzN1OQzYkRFVOJeodhHy2OdckgJC6yFlFBL0LmyA2lXUpy8mFtqxQuelmQZbD-wsxrNDzJ72CkIdg1fD6MLHQpmKQDIYEMaSFkz68nfQLsoIEOKVq6UPr0Yc-4YTsmqogaF_YN4lbUn5czmhHZgBtieitwhr4uKGJskj090kWPOQGAJae0GSUelPTk03v4vTP-efKVBb1Rj4INL1R6-la41HNJi5YXwUj1yU_gMaqzTD7G3et-fbVeMG0HuJfSENVAwzdcwmivBH-C6Iaq9g9B0BP0gC1HH_L6NSceOdSKWBYRnQZg67S3C9cGyDL8Sf29zcyNKDmjdyRnISxIy0sZZ3X3svt5QVlPS9xODl1ZztefeyfnDdgT1Rlw6bpj-UMhQSqMztclYvJXfvZOt4vQ97HYU8"
}

resource "bigipnext_cm_activate_instance_license" "tokenadd" {
  instances = [{
    instance_address = "10.144.10.181"
    jwt_id           = bigipnext_cm_add_jwt_token.tokenadd2.id
  }]
}
`

// const testAccCMNextLicenseActivateResourceTC2 = `
// resource "bigipnext_cm_add_jwt_token" "tokenadd2" {
// 	token_name = "paid_test_jwt"
// 	jwt_token  = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6InYxIiwiamt1IjoiaHR0cHM6Ly9wcm9kdWN0LmFwaXMuZjUuY29tL2VlL3YxL2tleXMvandrcyJ9.eyJzdWIiOiJGTkktZGM5MTFhYjctNDlhYi00ZWUzLTkwYjItNTNmNmVmOGY5ZTc0IiwiaWF0IjoxNzA1NjIyMDY5LCJpc3MiOiJGNSBJbmMuIiwiYXVkIjoidXJuOmY1OnRlZW0iLCJqdGkiOiJlYjIyZDhjMS1iNjVjLTExZWUtODQ1MC1mM2M0OTkzMmJmOTIiLCJmNV9vcmRlcl90eXBlIjoicGFpZCIsImY1X29yZGVyX3N1YnR5cGUiOiJpbnRlcm5hbCJ9.c6V5gAVGd-Pj62krRj8740bLq0_YyuRUvtKmat-oEJdiQn10GFbHNqBH8l3x0stcdE0UldrGszVQI3CmukDceRYi1XiTQpB69EubbOpx8Pe4qc6ht7kErmkDvsLlpy6ALYhdl8j2m5_npy3HvmoDnE2jjzWQkiQeFZjdxT4Gqc05LmRsO4_RnMOkvZFECfbRU6dEhtmP1es7L2FxJaJyJ8JEL2mz9kC8XtwaoW_jS_lxq_l5brDCnXJuFmLF882xyCReCT62FvIb4P4vzN1OQzYkRFVOJeodhHy2OdckgJC6yFlFBL0LmyA2lXUpy8mFtqxQuelmQZbD-wsxrNDzJ72CkIdg1fD6MLHQpmKQDIYEMaSFkz68nfQLsoIEOKVq6UPr0Yc-4YTsmqogaF_YN4lbUn5czmhHZgBtieitwhr4uKGJskj090kWPOQGAJae0GSUelPTk03v4vTP-efKVBb1Rj4INL1R6-la41HNJi5YXwUj1yU_gMaqzTD7G3et-fbVeMG0HuJfSENVAwzdcwmivBH-C6Iaq9g9B0BP0gC1HH_L6NSceOdSKWBYRnQZg67S3C9cGyDL8Sf29zcyNKDmjdyRnISxIy0sZZ3X3svt5QVlPS9xODl1ZztefeyfnDdgT1Rlw6bpj-UMhQSqMztclYvJXfvZOt4vQ97HYU8"
//   }
// `
