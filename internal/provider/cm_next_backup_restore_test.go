package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// backup with device Ip
func TestUnitNextCMBackupCreateResourceTC1(t *testing.T) {
	testAccPreUnitCheck(t)
	// var getCount = 0
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

	// Get Device info by address
	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "%s", loadFixtureString("./fixtures/getDeviceIdByIp.json"))
	})

	// Post backup call
	mux.HandleFunc("/api/device/v1/inventory/494975ee-5eba-420a-90a6-f3aeb75bbf5b/backup", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
			  "self": {
				"href": "/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a"
			  }
			},
			"path": "/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a"
		  }`)
	})

	// Check Status
	mux.HandleFunc("/api/device/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
				"self": {
					"href": "/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a"
				}
			},
			"completed": "2024-07-23T08:07:45.267699Z",
			"created": "2024-07-23T08:07:24.079439Z",
			"failure_reason": "",
			"file_name": "test",
			"file_path": "storage/backups",
			"id": "c41f1566-61f8-4be2-942a-9b01c687f81a",
			"instance_id": "494975ee-5eba-420a-90a6-f3aeb75bbf5b",
			"instance_name": "big-ip-next",
			"run_id": "",
			"state": "backupDone",
			"status": "completed"
		}`)
	})

	//Get/Delete call
	mux.HandleFunc("/api/device/v1/backups/test.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, ``)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{
					"_links": {
					  "self": {
						"href": "/v1/backups/test.tar.gz"
					  }
					},
					"file_date": "2024-07-23T10:49:22Z",
					"file_name": "test.tar.gz",
					"file_size": 1894032,
					"instance_id": "",
					"instance_name": ""
				  }`)
		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMBackupByIpResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// backup with device hostname
func TestUnitNextCMBackupCreateResourceTC2(t *testing.T) {
	testAccPreUnitCheck(t)
	// var getCount = 0
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

	// Get Device info by hostname
	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "%s", loadFixtureString("./fixtures/getDeviceIdByHostname.json"))
	})

	// Post backup call
	mux.HandleFunc("/api/device/v1/inventory/494975ee-5eba-420a-90a6-f3aeb75bbf5b/backup", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
			  "self": {
				"href": "/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a"
			  }
			},
			"path": "/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a"
		  }`)
	})

	// Check Status
	mux.HandleFunc("/api/device/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
				"self": {
					"href": "/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a"
				}
			},
			"completed": "2024-07-23T08:07:45.267699Z",
			"created": "2024-07-23T08:07:24.079439Z",
			"failure_reason": "",
			"file_name": "test",
			"file_path": "storage/backups",
			"id": "c41f1566-61f8-4be2-942a-9b01c687f81a",
			"instance_id": "494975ee-5eba-420a-90a6-f3aeb75bbf5b",
			"instance_name": "big-ip-next",
			"run_id": "",
			"state": "backupDone",
			"status": "completed"
		}`)
	})

	//Get/Delete call
	mux.HandleFunc("/api/device/v1/backups/test.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, ``)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{
					"_links": {
					  "self": {
						"href": "/v1/backups/test.tar.gz"
					  }
					},
					"file_date": "2024-07-23T10:49:22Z",
					"file_name": "test.tar.gz",
					"file_size": 1894032,
					"instance_id": "",
					"instance_name": ""
				  }`)
		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMBackupByHostnameResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// restore with device Ip
func TestUnitNextCMBackupCreateResourceTC3(t *testing.T) {
	testAccPreUnitCheck(t)
	// var getCount = 0
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

	// Get Device info by address
	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "%s", loadFixtureString("./fixtures/getDeviceIdByIp.json"))
	})

	// Post restore call
	mux.HandleFunc("/api/device/v1/inventory/494975ee-5eba-420a-90a6-f3aeb75bbf5b/restore", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
			  "self": {
				"href": "/v1/restore-tasks/016bba84-2abc-4acd-8acc-0239dbc2ec2d"
			  }
			},
			"path": "/v1/restore-tasks/016bba84-2abc-4acd-8acc-0239dbc2ec2d"
		  }`)
	})

	// Check Status
	mux.HandleFunc("/api/device/v1/restore-tasks/016bba84-2abc-4acd-8acc-0239dbc2ec2d", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
			  "self": {
				"href": "/v1/restore-tasks/016bba84-2abc-4acd-8acc-0239dbc2ec2d"
			  }
			},
			"completed": "2024-07-23T13:32:01.982649Z",
			"created": "2024-07-23T13:30:11.312178Z",
			"failure_reason": "",
			"file_name": "test.tar.gz",
			"id": "016bba84-2abc-4acd-8acc-0239dbc2ec2d",
			"instance_id": "494975ee-5eba-420a-90a6-f3aeb75bbf5b",
			"instance_name": "big-ip-next",
			"state": "restoreDone",
			"status": "completed"
		  }`)
	})

	//Get/Delete call
	mux.HandleFunc("/api/device/v1/backups/test.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, ``)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{
					"_links": {
					  "self": {
						"href": "/v1/backups/test.tar.gz"
					  }
					},
					"file_date": "2024-07-23T10:49:22Z",
					"file_name": "test.tar.gz",
					"file_size": 1894032,
					"instance_id": "",
					"instance_name": ""
				  }`)
		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMRestoreByIpResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// restore with device hostname
func TestUnitNextCMBackupCreateResourceTC4(t *testing.T) {
	testAccPreUnitCheck(t)
	// var getCount = 0
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

	// // Get Device info by address
	// mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	_, _ = fmt.Fprintf(w, "%s", loadFixtureString("./fixtures/getDeviceIdByIp.json"))
	// })

	// Get Device info by hostname
	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "%s", loadFixtureString("./fixtures/getDeviceIdByHostname.json"))
	})

	// Post restore call
	mux.HandleFunc("/api/device/v1/inventory/494975ee-5eba-420a-90a6-f3aeb75bbf5b/restore", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
			  "self": {
				"href": "/v1/restore-tasks/016bba84-2abc-4acd-8acc-0239dbc2ec2d"
			  }
			},
			"path": "/v1/restore-tasks/016bba84-2abc-4acd-8acc-0239dbc2ec2d"
		  }`)
	})

	// Check Status
	mux.HandleFunc("/api/device/v1/restore-tasks/016bba84-2abc-4acd-8acc-0239dbc2ec2d", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
			  "self": {
				"href": "/v1/restore-tasks/016bba84-2abc-4acd-8acc-0239dbc2ec2d"
			  }
			},
			"completed": "2024-07-23T13:32:01.982649Z",
			"created": "2024-07-23T13:30:11.312178Z",
			"failure_reason": "",
			"file_name": "test.tar.gz",
			"id": "016bba84-2abc-4acd-8acc-0239dbc2ec2d",
			"instance_id": "494975ee-5eba-420a-90a6-f3aeb75bbf5b",
			"instance_name": "big-ip-next",
			"state": "restoreDone",
			"status": "completed"
		  }`)
	})

	//Get/Delete call
	mux.HandleFunc("/api/device/v1/backups/test.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, ``)
		} else {
			// if getCount < 2 {
			// 	w.WriteHeader(http.StatusAccepted)
			// 	_, _ = fmt.Fprintf(w, "")
			// } else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{
					"_links": {
					  "self": {
						"href": "/v1/backups/test.tar.gz"
					  }
					},
					"file_date": "2024-07-23T10:49:22Z",
					"file_name": "test.tar.gz",
					"file_size": 1894032,
					"instance_id": "",
					"instance_name": ""
				  }`)
			// }

			// getCount++

		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMRestoreByHostnameResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// {
			// 	Config: testAccNextCMWafPolicyImportUpdateResourceConfig,
			// 	Check:  resource.ComposeAggregateTestCheckFunc(),
			// },
		},
	})
}

// update backup with device Ip
func TestUnitNextCMBackupUpdateResourceTC1(t *testing.T) {
	testAccPreUnitCheck(t)
	// var getCount = 0
	var postCount = 0
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

	// Get Device info by address
	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "%s", loadFixtureString("./fixtures/getDeviceIdByIp.json"))
	})

	// Post backup call
	mux.HandleFunc("/api/device/v1/inventory/494975ee-5eba-420a-90a6-f3aeb75bbf5b/backup", func(w http.ResponseWriter, r *http.Request) {
		t.Log("Post call")
		if postCount == 0 {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{
				"_links": {
				  "self": {
					"href": "/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a"
				  }
				},
				"path": "/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a"
			  }`)
		} else {
			_, _ = fmt.Fprintf(w, `{
				"_links": {
				  "self": {
					"href": "/v1/backup-tasks/0d1d9fb0-74d5-4924-ab11-85db571d26af"
				  }
				},
				"path": "/v1/backup-tasks/0d1d9fb0-74d5-4924-ab11-85db571d26af"
			  }`)
		}
		postCount++
	})

	// Check Status
	mux.HandleFunc("/api/device/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
				"self": {
					"href": "/v1/backup-tasks/c41f1566-61f8-4be2-942a-9b01c687f81a"
				}
			},
			"completed": "2024-07-23T08:07:45.267699Z",
			"created": "2024-07-23T08:07:24.079439Z",
			"failure_reason": "",
			"file_name": "test",
			"file_path": "storage/backups",
			"id": "c41f1566-61f8-4be2-942a-9b01c687f81a",
			"instance_id": "494975ee-5eba-420a-90a6-f3aeb75bbf5b",
			"instance_name": "big-ip-next",
			"run_id": "",
			"state": "backupDone",
			"status": "completed"
		}`)
	})

	mux.HandleFunc("/api/device/v1/backup-tasks/0d1d9fb0-74d5-4924-ab11-85db571d26af", func(w http.ResponseWriter, r *http.Request) {
		t.Log("Checking the status of task 6af")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
				"self": {
					"href": "/v1/backup-tasks/0d1d9fb0-74d5-4924-ab11-85db571d26af"
				}
			},
			"completed": "2024-07-23T08:07:45.267699Z",
			"created": "2024-07-23T08:07:24.079439Z",
			"failure_reason": "",
			"file_name": "test",
			"file_path": "storage/backups",
			"id": "0d1d9fb0-74d5-4924-ab11-85db571d26af",
			"instance_id": "494975ee-5eba-420a-90a6-f3aeb75bbf5b",
			"instance_name": "big-ip-next",
			"run_id": "",
			"state": "backupDone",
			"status": "completed"
		}`)
	})

	//Get/Delete call
	mux.HandleFunc("/api/device/v1/backups/test.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			t.Log("Inside the Delete Condition")
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, ``)
		} else {
			t.Log("Inside the Get Condition")
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{
					"_links": {
					  "self": {
						"href": "/v1/backups/test.tar.gz"
					  }
					},
					"file_date": "2024-07-23T10:49:22Z",
					"file_name": "test.tar.gz",
					"file_size": 1894032,
					"instance_id": "",
					"instance_name": ""
				  }`)
		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMBackupByIpResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// {
			// 	Config: testAccNextCMBackupByIpUpdateResourceConfig,
			// 	Check:  resource.ComposeAggregateTestCheckFunc(),
			// },
		},
	})
}

const testAccNextCMBackupByIpResourceConfig = `
resource "bigipnext_cm_backup_restore" "sample" {
	backup_password = "F5site02@123"
	operation       = "backup"
	file_name       = "test.tar.gz"
	device_ip = "10.218.132.39"
  }
`

// const testAccNextCMBackupByIpUpdateResourceConfig = `
// resource "bigipnext_cm_backup_restore" "sample" {
// 	backup_password = "F5site02@1234"
// 	operation       = "backup"
// 	file_name       = "test.tar.gz"
// 	device_ip = "10.218.132.39"
//   }
// `

const testAccNextCMBackupByHostnameResourceConfig = `
resource "bigipnext_cm_backup_restore" "sample" {
	backup_password = "F5site02@123"
	operation       = "backup"
	file_name       = "test.tar.gz"
	device_hostname = "big-ip-next"
  }
`

const testAccNextCMRestoreByIpResourceConfig = `
resource "bigipnext_cm_backup_restore" "sample" {
	backup_password = "F5site02@123"
	operation       = "restore"
	file_name       = "test.tar.gz"
	device_ip = "10.218.132.39"
  }
`

const testAccNextCMRestoreByHostnameResourceConfig = `
resource "bigipnext_cm_backup_restore" "sample" {
	backup_password = "F5site02@123"
	operation       = "restore"
	file_name       = "test.tar.gz"
	device_hostname = "big-ip-next"
  }
`
