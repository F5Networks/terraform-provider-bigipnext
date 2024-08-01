package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

func TestAccCMBootstrap(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCMBootstrapConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_bootstrap.cm_bootstrap", "run_setup", "false"),
					resource.TestCheckResourceAttr("bigipnext_cm_bootstrap.cm_bootstrap", "external_storage.storage_type", "NFS"),
					resource.TestCheckResourceAttr("bigipnext_cm_bootstrap.cm_bootstrap", "external_storage.storage_address", "10.218.134.22"),
				),
			},
		},
	})
}

func TestUnitCMBootstrap(t *testing.T) {
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

	mux.HandleFunc("/api/v1/system/infra/external-storage", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			payload, err := io.ReadAll(r.Body)
			if err != nil && err.Error() != "EOF" {
				t.Fatal(err)
			}

			payloadStruct := &bigipnextsdk.CMExternalStorage{}
			err = json.Unmarshal(payload, payloadStruct)
			if err != nil {
				t.Fatal(err)
			}

			if payloadStruct.StorageAddress != "10.218.134.22" {
				t.Errorf("Expected storage address to be 10.218.134.22, got %v", payloadStruct.StorageAddress)
			}
			if payloadStruct.StorageType != "NFS" {
				t.Errorf("Expected storage type to be NFS, got %v", payloadStruct.StorageType)
			}
			if payloadStruct.StorageSharePath != "/exports/backup" {
				t.Errorf("Expected storage path to be /exports/backup, got %v", payloadStruct.StorageSharePath)
			}

			fmt.Fprint(w, `
			{
    			"spec": {
					"storage_address": "10.218.134.22",
					"storage_share_dir": "",
					"storage_share_path": "/exports/backup",
					"storage_type": "NFS"
    			},
				"status": {
					"setup": "SUCCESSFUL"
				}
			}`)
		}
		if r.Method == http.MethodGet {
			fmt.Fprint(w, `
			{
    			"spec": {
					"storage_address": "10.218.134.22",
					"storage_share_dir": "9fb2d367-e8e6-4c09-abd2-faaee96ebe93",
					"storage_share_path": "/exports/backup",
					"storage_type": "NFS"
    			},
				"status": {
					"setup": "SUCCESSFUL"
				}
			}`)
		}
	})

	i := 0
	bootstrapStatuses := []string{
		`{
			"created": "2024-07-15T06:31:49.322459826Z",
			"status": "RUNNING",
			"step": "Installing Elasticsearch",
			"updated": "2024-07-15T06:40:36.083148739Z"
		}`,
		`{
			"created": "2024-07-15T06:31:49.322459826Z",
			"status": "RUNNING",
			"step": "Installing Central Manager applications",
			"updated": "2024-07-15T06:40:36.083148739Z"
		}`,
		`{
			"created": "2024-07-15T06:31:49.322459826Z",
			"status": "COMPLETED",
			"step": "Done",
			"updated": "2024-07-15T06:40:36.083148739Z"
		}`,
		`{
			"created": "2024-07-15T06:31:49.322459826Z",
			"status": "COMPLETED",
			"step": "Done",
			"updated": "2024-07-15T06:40:36.083148739Z"
		}`,
	}

	mux.HandleFunc("/api/v1/system/infra/bootstrap", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Fprint(w, `
			{
    			"created": "2024-07-15T06:31:49.322459826Z",
				"status": "RUNNING",
				"step": "Installing Kafka",
				"updated": "2024-07-15T06:40:36.083148739Z"
			}
			`)
		}
		if r.Method == http.MethodGet {
			fmt.Fprint(w, bootstrapStatuses[i])
			i += 1
		}
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testUnitCMBootstrapConfig,
			},
		},
	})
}

func TestUnitCMBootstrapSAMBA(t *testing.T) {
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

	mux.HandleFunc("/api/v1/system/infra/external-storage", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			payload, err := io.ReadAll(r.Body)
			if err != nil && err.Error() != "EOF" {
				t.Fatal(err)
			}

			payloadStruct := &bigipnextsdk.CMExternalStorage{}
			err = json.Unmarshal(payload, payloadStruct)

			if err != nil {
				t.Fatal(err)
			}

			if payloadStruct.StorageAddress != "10.22.33.44" {
				t.Errorf("Expected storage address to be 10.22.33.44, got %v", payloadStruct.StorageAddress)
			}
			if payloadStruct.StorageType != "SAMBA" {
				t.Errorf("Expected storage type to be SAMBA, got %v", payloadStruct.StorageType)
			}
			if payloadStruct.StorageSharePath != "/exports/backup" {
				t.Errorf("Expected storage path to be /exports/backup, got %v", payloadStruct.StorageSharePath)
			}
			if payloadStruct.StorageShareDir != "backup" {
				t.Errorf("Expected storage path to be backup, got %v", payloadStruct.StorageShareDir)
			}

			fmt.Fprint(w, `
			{
    			"spec": {
					"storage_address": "10.22.33.44",
					"storage_share_dir": "",
					"storage_share_path": "/exports/backup",
					"storage_type": "SAMBA",
					"storage_user": {
						"username": "admin",
						"password": "password"
					}
    			},
				"status": {
					"setup": "SUCCESSFUL"
				}
			}`)
		}
		if r.Method == http.MethodGet {
			fmt.Fprint(w, `
			{
    			"spec": {
					"storage_address": "10.22.33.44",
					"storage_share_dir": "backup",
					"storage_share_path": "/exports/backup",
					"storage_type": "SAMBA",
					"storage_user": {
						"username": "admin",
						"password": "password"
					}
    			},
				"status": {
					"setup": "SUCCESSFUL"
				}
			}`)
		}
	})

	mux.HandleFunc("/api/v1/system/infra/bootstrap", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
		{
			"created": "2024-07-15T06:31:49.322459826Z",
			"status": "SUCCESSFUL",
			"step": "Installing Kafka",
			"updated": "2024-07-15T06:40:36.083148739Z"
		}
		`)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testUnitCMBootstrapConfigSAMBA,
			},
		},
	})
}

const testAccCMBootstrapConfig = `
resource "bigipnext_cm_bootstrap" "cm_bootstrap" {
  run_setup = false
  external_storage = {
    storage_type    = "NFS"
    storage_address = "10.218.134.22"
    storage_path    = "/exports/backup"
  }
}
`

const testUnitCMBootstrapConfig = `
resource "bigipnext_cm_bootstrap" "cm_bootstrap" {
  run_setup = true
  external_storage = {
    storage_type    = "NFS"
    storage_address = "10.218.134.22"
    storage_path    = "/exports/backup"
  }
}
`

const testUnitCMBootstrapConfigSAMBA = `
resource "bigipnext_cm_bootstrap" "cm_bootstrap" {
  run_setup = false
  external_storage = {
    storage_type    = "SAMBA"
	storage_address = "10.22.33.44"
	storage_path    = "/exports/backup"
	cm_storage_dir  = "backup"
	username		= "admin"
	password		= "password"
  }
}
`
