package provider

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitCMNextAddJwtTokenResourceTC1(t *testing.T) {
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
	mux.HandleFunc("/api/v1/spaces/default/license/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"NewToken": {
			  "id": "94ff6fac-8920-492d-9545-9e837fe31a23",
			  "nickName": "paid_test_jwt",
			  "entitlement": "string",
			  "orderType": "string",
			  "orderSubType": "string",
			  "subscriptionExpiry": "string"
			},
			"DuplicatesTokenValue": {
			  "id": "94ff6fac-8920-492d-9545-9e837fe31a23",
			  "nickName": "paid_test_jwt",
			  "entitlement": "string",
			  "orderType": "string",
			  "orderSubType": "string",
			  "subscriptionExpiry": "string"
			},
			"DuplicatesTokenNickName": {
			  "id": "94ff6fac-8920-492d-9545-9e837fe31a23",
			  "nickName": "paid_test_jwt",
			  "entitlement": "string",
			  "orderType": "string",
			  "orderSubType": "string",
			  "subscriptionExpiry": "string"
			},
			"_links": {
			  "self": {
				"href": "string"
			  }
			},
			"duplicateToken": {
			  "_links": {
				"self": {
				  "href": "string"
				}
			  }
			},
			"duplicateShortname": {
			  "_links": {
				"self": {
				  "href": "string"
				}
			  }
			}
		  }`)
	})
	mux.HandleFunc("/api/v1/spaces/default/license/tokens/94ff6fac-8920-492d-9545-9e837fe31a23", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"id": "94ff6fac-8920-492d-9545-9e837fe31a23",
			"nickName": "paid_test_jwt",
			"entitlement": "string",
			"orderType": "string",
			"orderSubType": "string",
			"subscriptionExpiry": "string"
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
				Config: testUnitCMNextAddJwtTokenResourceTC1,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: testUnitCMNextAddJwtTokenResourceTC2,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestUnitCMNextAddJwtTokenResourceTC2(t *testing.T) {
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
	mux.HandleFunc("/api/v1/spaces/default/license/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `
		{
			"message": "description explaining the error cause in short",
			"help": "help to recover from the error",
			"code": "LICENSING-2245",
			"status": 400,
			"details": "inner details of the error"
		  }`)
	})
	mux.HandleFunc("/api/v1/spaces/default/license/tokens/94ff6fac-8920-492d-9545-9e837fe31a23", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"id": "94ff6fac-8920-492d-9545-9e837fe31a23",
			"nickName": "paid_test_jwt",
			"entitlement": "string",
			"orderType": "string",
			"orderSubType": "string",
			"subscriptionExpiry": "string"
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
				Config:      testUnitCMNextAddJwtTokenResourceTC1,
				ExpectError: regexp.MustCompile(`Failed to Create Jwt token`),
			},
		},
	})
}

func TestUnitCMNextAddJwtTokenResourceTC3(t *testing.T) {
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
	mux.HandleFunc("/api/v1/spaces/default/license/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"NewToken": {
			  "id": "94ff6fac-8920-492d-9545-9e837fe31a23",
			  "nickName": "paid_test_jwt",
			  "entitlement": "string",
			  "orderType": "string",
			  "orderSubType": "string",
			  "subscriptionExpiry": "string"
			},
			"DuplicatesTokenValue": {
			  "id": "94ff6fac-8920-492d-9545-9e837fe31a23",
			  "nickName": "paid_test_jwt",
			  "entitlement": "string",
			  "orderType": "string",
			  "orderSubType": "string",
			  "subscriptionExpiry": "string"
			},
			"DuplicatesTokenNickName": {
			  "id": "94ff6fac-8920-492d-9545-9e837fe31a23",
			  "nickName": "paid_test_jwt",
			  "entitlement": "string",
			  "orderType": "string",
			  "orderSubType": "string",
			  "subscriptionExpiry": "string"
			},
			"_links": {
			  "self": {
				"href": "string"
			  }
			},
			"duplicateToken": {
			  "_links": {
				"self": {
				  "href": "string"
				}
			  }
			},
			"duplicateShortname": {
			  "_links": {
				"self": {
				  "href": "string"
				}
			  }
			}
		  }`)
	})
	mux.HandleFunc("/api/v1/spaces/default/license/tokens/94ff6fac-8920-492d-9545-9e837fe31a23", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"id": "94ff6fac-8920-492d-9545-9e837fe31a23",
			"nickName": "paid_test_jwt",
			"entitlement": "string",
			"orderType": "paid",
			"orderSubType": "string",
			"subscriptionExpiry": "string"
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
				Config: testUnitCMNextAddJwtTokenResourceTC1,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: testUnitCMNextAddJwtTokenResourceTC1Update,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

//	{
//	    "code": 400,
//	    "error": {
//	        "code": "LICENSING-0004",
//	        "message": "LICENSING-0004: JWT is not verifiable: token contains an invalid number of segments",
//	        "details": "",
//	        "help": "please retry or contact admin"
//	    }
//	}
//
// TC1: Test case when JWT token is not verifiable
func TestAccCMNextAddJwtTokenResourceTC1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testUnitCMNextAddJwtTokenResourceTC1,
				ExpectError: regexp.MustCompile(`Failed to Create Jwt token`),
			},
		},
	})
}

// {"code":409,"error":{"DuplicatesTokenNickName":{"orderType":"","shortName":"","tokenId":"00000000-0000-0000-0000-000000000000"},"DuplicatesTokenValue":{"entitlement":"{\"compliance\":{\"digitalAssetComplianceStatus\":\"valid\",\"digitalAssetDaysRemainingInState\":0,\"digitalAssetExpiringSoon\":false,\"digitalAssetOutOfComplianceDate\":\"\",\"entitlementCheckStatus\":\"valid\",\"entitlementExpiryStatus\":\"valid\",\"telemetryStatus\":\"valid\",\"usageExceededStatus\":\"valid\"},\"documentType\":\"\",\"documentVersion\":\"\",\"digitalAsset\":{\"digitalAssetId\":\"\",\"digitalAssetName\":\"\",\"digitalAssetVersion\":\"\",\"telemetryId\":\"\"},\"entitlementMetadata\":{\"complianceEnforcements\":[\"entitlement\"],\"complianceStates\":{\"device-twin\":[\"in-grace-period\",\"in-enforcement-period\",\"non-functional\"],\"entitlement\":[\"in-grace-period\",\"in-enforcement-period\",\"non-functional\"],\"telemetry\":[\"in-grace-period\",\"in-enforcement-period\",\"non-functional\"],\"usage\":[\"in-grace-period\",\"in-enforcement-period\"]},\"enforcementBehavior\":\"visibility\",\"enforcementPeriodDays\":0,\"entitlementModel\":\"aggregated\",\"expiringSoonNotificationDays\":7,\"entitlementExpiryDate\":\"2024-12-07T00:00:00Z\",\"gracePeriodDays\":0,\"nonContactPeriodHours\":0,\"nonFunctionalPeriodDays\":14,\"orderSubType\":\"internal\",\"orderType\":\"paid\"},\"subscriptionMetadata\":{\"programName\":\"big_ip_next_internal\",\"programTypeDescription\":\"big_ip_next_internal\",\"subscriptionId\":\"A-S00019374\",\"subscriptionExpiryDate\":\"2024-12-07T00:00:00.000Z\",\"subscriptionNotifyDays\":\"\"},\"RepositoryCertificateMetadata\":{\"sslCertificate\":\"\",\"privateKey\":\"\"},\"entitledFeatures\":[{\"entitledFeatureId\":\"\",\"featureFlag\":\"UNKNOWN_set-bigip_ltm_module_8_bigip_vcpus\",\"featurePermitted\":0,\"featureRemain\":0,\"featureUnlimited\":true,\"featureUsed\":14,\"featureValueType\":\"integral\",\"uomCode\":\"\",\"uomTerm\":\"\",\"uomTermStart\":0},{\"entitledFeatureId\":\"\",\"featureFlag\":\"UNKNOWN_set-bigip_waf_module_8_bigip_vcpus\",\"featurePermitted\":0,\"featureRemain\":0,\"featureUnlimited\":true,\"featureUsed\":14,\"featureValueType\":\"integral\",\"uomCode\":\"\",\"uomTerm\":\"\",\"uomTermStart\":0},{\"entitledFeatureId\":\"\",\"featureFlag\":\"UNKNOWN_set-bigip_ltm_module_1_bigip_system\",\"featurePermitted\":0,\"featureRemain\":0,\"featureUnlimited\":true,\"featureUsed\":1,\"featureValueType\":\"integral\",\"uomCode\":\"\",\"uomTerm\":\"\",\"uomTermStart\":0},{\"entitledFeatureId\":\"\",\"featureFlag\":\"UNKNOWN_set-bigip_waf_module_1_bigip_system\",\"featurePermitted\":0,\"featureRemain\":0,\"featureUnlimited\":true,\"featureUsed\":1,\"featureValueType\":\"integral\",\"uomCode\":\"\",\"uomTerm\":\"\",\"uomTermStart\":0},{\"entitledFeatureId\":\"\",\"featureFlag\":\"UNKNOWN_set-bigip_waf_module_2_bigip_vcpus\",\"featurePermitted\":0,\"featureRemain\":0,\"featureUnlimited\":true,\"featureUsed\":2,\"featureValueType\":\"integral\",\"uomCode\":\"\",\"uomTerm\":\"\",\"uomTermStart\":0},{\"entitledFeatureId\":\"\",\"featureFlag\":\"UNKNOWN_set-bigip_ltm_module_2_bigip_vcpus\",\"featurePermitted\":0,\"featureRemain\":0,\"featureUnlimited\":true,\"featureUsed\":2,\"featureValueType\":\"integral\",\"uomCode\":\"\",\"uomTerm\":\"\",\"uomTermStart\":0}]}","orderSubType":"internal","orderType":"paid","shortName":"token2","tokenId":"f88665b2-2ad3-4110-a7be-e0ccdf8667ff"},"NewToken":{"orderType":"","shortName":"","tokenId":"00000000-0000-0000-0000-000000000000"}}
// TC2: Test case when JWT token is already present
func TestAccCMNextAddJwtTokenResourceTC2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testAccCMNextAddJwtTokenResourceTC2,
				ExpectError: regexp.MustCompile(`Failed to Create Jwt token`),
			},
		},
	})
}

// TC3: Test case when JWT token is not verifiable
func TestAccCMNextAddJwtTokenResourceTC3(t *testing.T) {
	os.Setenv("BIGIPNEXT_HOST", "10.145.64.3")
	os.Setenv("BIGIPNEXT_USERNAME", "admin")
	os.Setenv("BIGIPNEXT_PASSWORD", "F5site02@123")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCMNextAddJwtTokenResourceTC2,
				// ExpectError: regexp.MustCompile(`Failed to Create Jwt token`),
			},
		},
	})
}

const testUnitCMNextAddJwtTokenResourceTC1 = `
resource "bigipnext_cm_add_jwt_token" "tokenadd" {
	token_name = "paid_test_jwt"
	jwt_token  = "eyJhbG"
  }
`

const testUnitCMNextAddJwtTokenResourceTC2 = `
resource "bigipnext_cm_add_jwt_token" "tokenadd" {
	token_name = "paid_test_jwt"
	jwt_token  = "eyJhbG"
  }
`

const testUnitCMNextAddJwtTokenResourceTC1Update = `
resource "bigipnext_cm_add_jwt_token" "tokenadd" {
	token_name = "paid_test_jwt"
	jwt_token  = "eyJhbGC"
  }
`

const testAccCMNextAddJwtTokenResourceTC2 = `
resource "bigipnext_cm_add_jwt_token" "tokenadd2" {
	token_name = "paid_test_jwt"
	jwt_token  = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6InYxIiwiamt1IjoiaHR0cHM6Ly9wcm9kdWN0LmFwaXMuZjUuY29tL2VlL3YxL2tleXMvandrcyJ9.eyJzdWIiOiJGTkktZGM5MTFhYjctNDlhYi00ZWUzLTkwYjItNTNmNmVmOGY5ZTc0IiwiaWF0IjoxNzA1NjIyMDY5LCJpc3MiOiJGNSBJbmMuIiwiYXVkIjoidXJuOmY1OnRlZW0iLCJqdGkiOiJlYjIyZDhjMS1iNjVjLTExZWUtODQ1MC1mM2M0OTkzMmJmOTIiLCJmNV9vcmRlcl90eXBlIjoicGFpZCIsImY1X29yZGVyX3N1YnR5cGUiOiJpbnRlcm5hbCJ9.c6V5gAVGd-Pj62krRj8740bLq0_YyuRUvtKmat-oEJdiQn10GFbHNqBH8l3x0stcdE0UldrGszVQI3CmukDceRYi1XiTQpB69EubbOpx8Pe4qc6ht7kErmkDvsLlpy6ALYhdl8j2m5_npy3HvmoDnE2jjzWQkiQeFZjdxT4Gqc05LmRsO4_RnMOkvZFECfbRU6dEhtmP1es7L2FxJaJyJ8JEL2mz9kC8XtwaoW_jS_lxq_l5brDCnXJuFmLF882xyCReCT62FvIb4P4vzN1OQzYkRFVOJeodhHy2OdckgJC6yFlFBL0LmyA2lXUpy8mFtqxQuelmQZbD-wsxrNDzJ72CkIdg1fD6MLHQpmKQDIYEMaSFkz68nfQLsoIEOKVq6UPr0Yc-4YTsmqogaF_YN4lbUn5czmhHZgBtieitwhr4uKGJskj090kWPOQGAJae0GSUelPTk03v4vTP-efKVBb1Rj4INL1R6-la41HNJi5YXwUj1yU_gMaqzTD7G3et-fbVeMG0HuJfSENVAwzdcwmivBH-C6Iaq9g9B0BP0gC1HH_L6NSceOdSKWBYRnQZg67S3C9cGyDL8Sf29zcyNKDmjdyRnISxIy0sZZ3X3svt5QVlPS9xODl1ZztefeyfnDdgT1Rlw6bpj-UMhQSqMztclYvJXfvZOt4vQ97HYU8"
  }
`
