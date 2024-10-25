package provider

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitCMNextUpgradeVE(t *testing.T) {
	testAccPreUnitCheck(t)
	defer teardown()

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

	dummyInstanceId := "dummyInstanceId"
	dummyUpgradeTaskId := "dummyUpgradeId"

	mux.HandleFunc("/api/v1/spaces/default/instances", func(w http.ResponseWriter, r *http.Request) {
		instanceIdJson := fmt.Sprintf(
			`{"_embedded":{"devices":[{"id":"%s","address":"10.218.140.52"}]}}`,
			dummyInstanceId,
		)
		fmt.Fprint(w, instanceIdJson)
	})

	mux.HandleFunc("/api/device/v1/proxy/dummyInstanceId", func(w http.ResponseWriter, r *http.Request) {
		resp, _ := os.ReadFile("fixtures/filesOnNextInstance.json")
		fmt.Fprint(w, string(resp))
	})

	mux.HandleFunc("/api/v1/spaces/default/instances/dummyInstanceId/upgrade", func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf(`{"path": "/api/v1/spaces/default/instances/upgrade-tasks/%s"}`, dummyUpgradeTaskId)
		fmt.Fprint(w, resp)
	})

	mux.HandleFunc(fmt.Sprintf("/api/v1/spaces/default/instances/upgrade-tasks/%s", dummyUpgradeTaskId), func(w http.ResponseWriter, r *http.Request) {
		resp, _ := os.ReadFile("fixtures/nextInstanceUpgradeTask.json")
		fmt.Fprint(w, string(resp))
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:  terraformUnitTestCfg,
				Destroy: false,
			},
			{
				Config:  terraformUnitTestUpdateCfg,
				Destroy: false,
			},
		},
	})
}

func TestUnitCMNextUpgradeVEInstanceId404(t *testing.T) {
	testAccPreUnitCheck(t)
	defer teardown()

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
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "NEXT instance not found")
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformUnitTestCfg,
				ExpectError: regexp.MustCompile(`{"code":404,"error":NEXT instance not found`),
			},
		},
	})
}

func TestUnitCMNextUpgradeVEUpdateInstanceId404(t *testing.T) {
	testAccPreUnitCheck(t)
	defer teardown()

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

	dummyInstanceId := "dummyInstanceId"
	dummyUpgradeTaskId := "dummyUpgradeId"

	count := 0

	mux.HandleFunc("/api/v1/spaces/default/instances", func(w http.ResponseWriter, r *http.Request) {
		if count == 0 {
			instanceIdJson := fmt.Sprintf(
				`{"_embedded":{"devices":[{"id":"%s","address":"10.218.140.52"}]}}`,
				dummyInstanceId,
			)
			fmt.Fprint(w, instanceIdJson)
			count++
		} else if count > 0 {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "NEXT instance not found")
		}
	})

	mux.HandleFunc("/api/device/v1/proxy/dummyInstanceId", func(w http.ResponseWriter, r *http.Request) {
		resp, _ := os.ReadFile("fixtures/filesOnNextInstance.json")
		fmt.Fprint(w, string(resp))
	})

	mux.HandleFunc("/api/v1/spaces/default/instances/dummyInstanceId/upgrade", func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf(`{"path": "/api/v1/spaces/default/instances/upgrade-tasks/%s"}`, dummyUpgradeTaskId)
		fmt.Fprint(w, resp)
	})

	mux.HandleFunc(fmt.Sprintf("/api/v1/spaces/default/instances/upgrade-tasks/%s", dummyUpgradeTaskId), func(w http.ResponseWriter, r *http.Request) {
		resp, _ := os.ReadFile("fixtures/nextInstanceUpgradeTask.json")
		fmt.Fprint(w, string(resp))
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:  terraformUnitTestCfg,
				Destroy: false,
			},
			{
				Config:      terraformUnitTestUpdateCfg,
				ExpectError: regexp.MustCompile(`{"code":404,"error":NEXT instance not found`),
				Destroy:     false,
			},
		},
	})
}

func TestUnitCMNextUpgradeVEImageIdError(t *testing.T) {
	testAccPreUnitCheck(t)
	defer teardown()

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

	dummyInstanceId := "dummyInstanceId"

	mux.HandleFunc("/api/v1/spaces/default/instances", func(w http.ResponseWriter, r *http.Request) {
		instanceIdJson := fmt.Sprintf(
			`{"_embedded":{"devices":[{"id":"%s","address":"10.218.140.52"}]}}`,
			dummyInstanceId,
		)
		fmt.Fprint(w, instanceIdJson)
	})

	mux.HandleFunc("/api/device/v1/proxy/dummyInstanceId", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"_embedded":{"files":[{}]}}`)
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformUnitTestCfg,
				ExpectError: regexp.MustCompile(`image or signature file not found`),
			},
		},
	})
}

func TestUnitCMNextUpgradeVEUpdateImageId404(t *testing.T) {
	testAccPreUnitCheck(t)
	defer teardown()

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

	dummyInstanceId := "dummyInstanceId"
	dummyUpgradeTaskId := "dummyUpgradeId"

	mux.HandleFunc("/api/v1/spaces/default/instances", func(w http.ResponseWriter, r *http.Request) {
		instanceIdJson := fmt.Sprintf(
			`{"_embedded":{"devices":[{"id":"%s","address":"10.218.140.52"}]}}`,
			dummyInstanceId,
		)
		fmt.Fprint(w, instanceIdJson)
	})

	count := 0
	mux.HandleFunc("/api/device/v1/proxy/dummyInstanceId", func(w http.ResponseWriter, r *http.Request) {
		if count == 0 {
			resp, _ := os.ReadFile("fixtures/filesOnNextInstance.json")
			fmt.Fprint(w, string(resp))
			count++
		} else if count > 0 {
			fmt.Fprint(w, `{"_embedded":{"files":[{}]}}`)
		}
	})

	mux.HandleFunc("/api/v1/spaces/default/instances/dummyInstanceId/upgrade", func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf(`{"path": "/api/v1/spaces/default/instances/upgrade-tasks/%s"}`, dummyUpgradeTaskId)
		fmt.Fprint(w, resp)
	})

	mux.HandleFunc(fmt.Sprintf("/api/v1/spaces/default/instances/upgrade-tasks/%s", dummyUpgradeTaskId), func(w http.ResponseWriter, r *http.Request) {
		resp, _ := os.ReadFile("fixtures/nextInstanceUpgradeTask.json")
		fmt.Fprint(w, string(resp))
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:  terraformUnitTestCfg,
				Destroy: false,
			},
			{
				Config:      terraformUnitTestUpdateCfg,
				ExpectError: regexp.MustCompile(`image or signature file not found`),
				Destroy:     false,
			},
		},
	})
}

func TestUnitCMNextUpgradeVEUpgradeError(t *testing.T) {
	testAccPreUnitCheck(t)
	defer teardown()

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

	dummyInstanceId := "dummyInstanceId"

	mux.HandleFunc("/api/v1/spaces/default/instances", func(w http.ResponseWriter, r *http.Request) {
		instanceIdJson := fmt.Sprintf(
			`{"_embedded":{"devices":[{"id":"%s","address":"10.218.140.52"}]}}`,
			dummyInstanceId,
		)
		fmt.Fprint(w, instanceIdJson)
	})

	mux.HandleFunc("/api/device/v1/proxy/dummyInstanceId", func(w http.ResponseWriter, r *http.Request) {
		resp, _ := os.ReadFile("fixtures/filesOnNextInstance.json")
		fmt.Fprint(w, string(resp))
	})

	mux.HandleFunc("/api/v1/spaces/default/instances/dummyInstanceId/upgrade", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "server error")
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformUnitTestCfg,
				ExpectError: regexp.MustCompile(`{"code":503,"error":server error`),
			},
		},
	})
}

func TestUnitCMNextUpgradeApplianceUpgradeError(t *testing.T) {
	testAccPreUnitCheck(t)
	defer teardown()

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

	dummyInstanceId := "dummyInstanceId"

	mux.HandleFunc("/api/v1/spaces/default/instances", func(w http.ResponseWriter, r *http.Request) {
		instanceIdJson := fmt.Sprintf(
			`{"_embedded":{"devices":[{"id":"%s","address":"10.144.140.85"}]}}`,
			dummyInstanceId,
		)
		fmt.Fprint(w, instanceIdJson)
	})

	mux.HandleFunc("/api/v1/spaces/default/instances/dummyInstanceId/upgrade", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "server error")
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformUnitTestCfgAppliance,
				ExpectError: regexp.MustCompile(`{"code":503,"error":server error`),
			},
		},
	})
}

func TestUnitCMNextUpgradeApplianceUpgradeUpdateError(t *testing.T) {
	testAccPreUnitCheck(t)
	defer teardown()

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

	dummyInstanceId := "dummyInstanceId"
	dummyUpgradeTaskId := "dummyUpgradeId"

	mux.HandleFunc("/api/v1/spaces/default/instances", func(w http.ResponseWriter, r *http.Request) {
		instanceIdJson := fmt.Sprintf(
			`{"_embedded":{"devices":[{"id":"%s","address":"10.144.140.85"}]}}`,
			dummyInstanceId,
		)
		fmt.Fprint(w, instanceIdJson)
	})

	count := 0
	mux.HandleFunc("/api/v1/spaces/default/instances/dummyInstanceId/upgrade", func(w http.ResponseWriter, r *http.Request) {
		if count == 0 {
			resp := fmt.Sprintf(`{"path": "/api/v1/spaces/default/instances/upgrade-tasks/%s"}`, dummyUpgradeTaskId)
			fmt.Fprint(w, resp)
			count++
		} else if count > 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "server error")
		}
	})

	mux.HandleFunc(fmt.Sprintf("/api/v1/spaces/default/instances/upgrade-tasks/%s", dummyUpgradeTaskId), func(w http.ResponseWriter, r *http.Request) {
		resp, _ := os.ReadFile("fixtures/nextInstanceUpgradeTask.json")
		fmt.Fprint(w, string(resp))
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:  terraformUnitTestCfgAppliance,
				Destroy: false,
			},
			{
				Config:      terraformUnitTestCfgApplianceUpdate,
				ExpectError: regexp.MustCompile(`{"code":503,"error":server error`),
				Destroy:     false,
			},
		},
	})
}

func TestUnitCMNextUpgradeAppliance(t *testing.T) {
	testAccPreUnitCheck(t)
	defer teardown()

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

	dummyInstanceId := "dummyInstanceId"
	dummyUpgradeTaskId := "dummyUpgradeId"

	mux.HandleFunc("/api/v1/spaces/default/instances", func(w http.ResponseWriter, r *http.Request) {
		instanceIdJson := fmt.Sprintf(
			`{"_embedded":{"devices":[{"id":"%s","address":"10.144.140.85"}]}}`,
			dummyInstanceId,
		)
		fmt.Fprint(w, instanceIdJson)
	})

	mux.HandleFunc("/api/v1/spaces/default/instances/dummyInstanceId/upgrade", func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf(`{"path": "/api/v1/spaces/default/instances/upgrade-tasks/%s"}`, dummyUpgradeTaskId)
		fmt.Fprint(w, resp)
	})

	mux.HandleFunc(fmt.Sprintf("/api/v1/spaces/default/instances/upgrade-tasks/%s", dummyUpgradeTaskId), func(w http.ResponseWriter, r *http.Request) {
		resp, _ := os.ReadFile("fixtures/nextInstanceUpgradeTask.json")
		fmt.Fprint(w, string(resp))
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:  terraformUnitTestCfgAppliance,
				Destroy: false,
			},
			{
				Config:  terraformUnitTestCfgApplianceUpdate,
				Destroy: false,
			},
		},
	})
}

const terraformUnitTestCfg = `
resource "bigipnext_cm_next_upgrade" "upgrade_next_ve" {
  upgrade_type       = "ve"
  next_instance_ip   = "10.218.140.52"
  image_name         = "BIG-IP-Next-20.2.1-2.430.2+0.0.48.tgz"
  signature_filename = "BIG-IP-Next-20.2.1-2.430.2+0.0.48.tgz.512.sig"
  timeout            = 900
}
`

const terraformUnitTestUpdateCfg = `
resource "bigipnext_cm_next_upgrade" "upgrade_next_ve" {
  upgrade_type       = "ve"
  next_instance_ip   = "10.218.140.52"
  image_name         = "BIG-IP-Next-20.3.0-2.713.1.tgz"
  signature_filename = "BIG-IP-Next-20.3.0-2.713.1.tgz.512.sig"
  timeout            = 900
}
`

const terraformUnitTestCfgAppliance = `
resource "bigipnext_cm_next_upgrade" "upgrade_next" {
  upgrade_type = "appliance"
  tenant_name = "testtenantx"
  next_instance_ip = "10.144.140.85"
  image_name = "BIG-IP-Next-20.3.0-2.538.0"
  partition_address = "10.144.140.80"
  partition_port = 8888
  partition_username = "admin"
  partition_password = "ess-pwe-f5site02"
}
`

const terraformUnitTestCfgApplianceUpdate = `
resource "bigipnext_cm_next_upgrade" "upgrade_next" {
  upgrade_type       = "appliance"
  tenant_name        = "testtenantx"
  next_instance_ip   = "10.144.140.85"
  image_name         = "BIG-IP-Next-20.3.1-2.538.0"
  partition_address  = "10.144.140.80"
  partition_port     = 8888
  partition_username = "admin"
  partition_password = "ess-pwe-f5site02"
}
`
