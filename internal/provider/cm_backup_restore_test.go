package provider

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCMDeviceBackUpTC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCMBackupConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "name", "Test"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "encryption_password", "F5site02@1234"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "backup", "true"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "type", "Light"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "scheduled", "false"),
				),
			},
			{
				Config:      testAccCMBackupUpdateConfig,
				ExpectError: regexp.MustCompile(`Only Scheduled Backup can be updated`),
			},
		},
	})
}

func TestAccCMDeviceRestoreTC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCMRestoreConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "name", "Backup-20241007-075720_L_20.3.1-0.5.0_8.tgz"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "encryption_password", "F5site02@123"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "backup", "false"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "type", "Restore"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "scheduled", "false"),
				),
			},
			{
				Config:      testAccCMRestoreUpdateConfig,
				ExpectError: regexp.MustCompile(`Restore can't be updated`),
			},
		},
	})
}

func TestAccCMDeviceScheduledBackUpTC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCMScheduledBackupConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "name", "Test"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "encryption_password", "F5site02@1234"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "backup", "true"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "type", "DaysOfTheWeek"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "scheduled", "true"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "frequency", "Weekly"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "schedule.start_at", "2025-08-25T18:30:00Z"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "schedule.end_at", "2025-09-25T18:30:00Z"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "days_of_the_week_to_run.0", "1"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "days_of_the_week_to_run.1", "2"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "days_of_the_week_to_run.2", "3"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "days_of_the_week_to_run.3", "4"),
				),
			},
			{
				Config: testAccCMScheduledBackupUpdate1Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "name", "Test"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "encryption_password", "F5site02@1234"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "backup", "true"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "type", "DayAndTimeOfTheMonth"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "scheduled", "true"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "frequency", "Monthly"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "schedule.start_at", "2025-08-25T18:30:00Z"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "schedule.end_at", "2025-09-25T18:30:00Z"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "day_of_the_month_to_run", "10"),
				),
			},
			{
				Config: testAccCMScheduledBackupUpdate2Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "name", "Test"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "encryption_password", "F5site02@1234"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "backup", "true"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "type", "BasicWithInterval"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "scheduled", "true"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "frequency", "Daily"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "schedule.start_at", "2025-08-25T18:30:00Z"),
				),
			},
		},
	})
}

func TestUnitCMBackupUnitTC1Resource(t *testing.T) {
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

	// post call
	mux.HandleFunc("/api/v1/system/backups", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_links": {
				"self": {
					"href": "/api/v1/system/backup-tasks/29761bc5-7f44-4b7e-9bbf-184eec783653"
				}
			},
			"path": "29761bc5-7f44-4b7e-9bbf-184eec783653"
		}`)
	})

	// backup-tasks status check
	mux.HandleFunc("/api/v1/system/backup-tasks/29761bc5-7f44-4b7e-9bbf-184eec783653", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_embedded": {
				"tasks": [
					{
						"_links": {
							"self": {
								"href": "/api/v1/system/backup-tasks/29761bc5-7f44-4b7e-9bbf-184eec783653/29761bc5-7f44-4b7e-9bbf-184eec783653"
							}
						},
						"backup_type": "light",
						"cm_version": "20.3.0-0.14.14",
						"completed": "2024-08-27T10:57:07.5639Z",
						"created": "2024-08-27T10:56:32.233816Z",
						"failure_reason": "",
						"file_name": "Backup-20240827-105632_L_20.3.0-0.14.14_4.tgz",
						"id": "29761bc5-7f44-4b7e-9bbf-184eec783653",
						"percentage": 100,
						"progress": "4 of 4",
						"schedule_id": null,
						"state": null,
						"status": "COMPLETED"
					}
				]
			},
			"_links": {
				"self": {
					"href": "/api/v1/system/backup-tasks/29761bc5-7f44-4b7e-9bbf-184eec783653"
				}
			}
		}`)
	})

	// read call
	mux.HandleFunc("/api/system/v1/files", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
				"_embedded": {
					"files": [
						{
							"_links": {
								"self": {
									"href": "/v1/files?filter=file_name+eq+Backup-20240827-105632_L_20.3.0-0.14.14_4.tgz/29761bc5-7f44-4b7e-9bbf-184eec783653"
								}
							},
							"description": "CM Backup",
							"file_name": "Backup-20240827-105632_L_20.3.0-0.14.14_4.tgz",
							"file_size": 5059307,
							"file_type": "backup",
							"id": "29761bc5-7f44-4b7e-9bbf-184eec783653",
							"updated": "2024-08-27T10:57:05.403587Z"
						}
					]
				}
			}
			`)
	})

	// delete call
	mux.HandleFunc("/api/system/v1/files/29761bc5-7f44-4b7e-9bbf-184eec783653", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, ``)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCMBackupConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "id", "29761bc5-7f44-4b7e-9bbf-184eec783653"),
				),
			},
		},
	})
}

func TestUnitCMBackupRestoreUnitTC1Resource(t *testing.T) {
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

	// read call
	mux.HandleFunc("/api/v1/system/backups", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"_embedded": {
				"backups": [
					{
						"_links": {
							"self": {
								"href": "/api/v1/system/backups?filter=file_name eq Backup-20240827-113901_L_20.3.0-0.14.14_4.tgz&=/"
							}
						},
						"backup_type": "light",
						"cm_version": "20.3.0-0.14.14",
						"file_date": "2024-08-27T11:39:01.892947Z",
						"file_id": "bb475827-0377-455e-a4b3-32b9f76b668b",
						"file_name": "Backup-20240827-113901_L_20.3.0-0.14.14_4.tgz",
						"file_size": 5061237,
						"status": "COMPLETED"
					}
				]
			},
			"_links": {
				"self": {
					"href": "/api/v1/system/backups?filter=file_name eq Backup-20240827-113901_L_20.3.0-0.14.14_4.tgz&="
				}
			}
		}`)
	})

	// post call
	mux.HandleFunc("/api/v1/system/restore", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/system/restore-tasks/5b4c0da6-277e-4539-81b8-24d247985033"}},"path":"5b4c0da6-277e-4539-81b8-24d247985033"}`)
	})

	mux.HandleFunc("/api/v1/system/restore-tasks", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_embedded":{"tasks":[{"_links":{"self":{"href":"/api/v1/system/restore-tasks/"}},"backup_type":null,"cm_version":null,"completed":"","created":"","failure_reason":"","file_name":null,"id":"","percentage":100,"progress":"10 of 10 Restore Done","schedule_id":null,"state":"Restore Completed","status":"COMPLETED"}]},"_links":{"self":{"href":"/api/v1/system/restore-tasks"}}}`)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCMRestoreConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(
				// resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "common_name", "sample_cert_common_name"),
				),
			},
		},
	})
}

func TestUnitCMBackupScheduleUnitTC1Resource(t *testing.T) {
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

	// post call
	mux.HandleFunc("/api/v1/system/backups/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/backups/schedule"}},"id":"dff97737-8686-44c6-9e0c-830d7b49b94a","message":"Task is scheduled.First Schedule will be on 2025-08-26 10:30:00 +0000 UTC","path":"/api/v1/system/backups/schedule/dff97737-8686-44c6-9e0c-830d7b49b94a"}`)
	})

	mux.HandleFunc("/api/system/v1/schedules/dff97737-8686-44c6-9e0c-830d7b49b94a", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
				"_links": {
					"self": {
						"href": "/v1/schedules/dff97737-8686-44c6-9e0c-830d7b49b94a"
					}
				},
				"days_of_the_week": {
					"days_of_the_week_to_run": [
						1,
						2,
						3,
						4
					],
					"hour_to_run_on": 10,
					"interval": 0,
					"minute_to_run_on": 30
				},
				"description": "CM Backup",
				"end_date": "2025-09-25T18:30:00Z",
				"id": "dff97737-8686-44c6-9e0c-830d7b49b94a",
				"job_type": "NORMAL",
				"name": "test",
				"next_run_time": "2025-08-26T10:30:00Z",
				"request_timeout_period": 30,
				"schedule_type": "DAYS_OF_THE_WEEK",
				"start_date": "2025-08-26T10:30:00Z",
				"status": "ENABLED",
				"tag": "cm-backup-tag",
				"topic": "cm-backup"
			}`)
	})

	mux.HandleFunc("/api/v1/system/backups/schedule/dff97737-8686-44c6-9e0c-830d7b49b94a", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, ``)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCMScheduledBackupConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "id", "dff97737-8686-44c6-9e0c-830d7b49b94a"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "type", "DaysOfTheWeek"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "days_of_the_week_to_run.0", "1"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "days_of_the_week_to_run.1", "2"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "days_of_the_week_to_run.2", "3"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "days_of_the_week_to_run.3", "4"),
				),
			},
		},
	})
}

func TestUnitCMBackupScheduleUnitTC2Resource(t *testing.T) {
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

	// post call
	mux.HandleFunc("/api/v1/system/backups/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/backups/schedule"}},"id":"dff97737-8686-44c6-9e0c-830d7b49b94a","message":"Task is scheduled.First Schedule will be on 2025-08-26 10:30:00 +0000 UTC","path":"/api/v1/system/backups/schedule/dff97737-8686-44c6-9e0c-830d7b49b94a"}`)
	})

	mux.HandleFunc("/api/system/v1/schedules/dff97737-8686-44c6-9e0c-830d7b49b94a", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
				"_links": {
					"self": {
						"href": "/v1/schedules/dff97737-8686-44c6-9e0c-830d7b49b94a"
					}
				},
				"day_and_time_of_month": {
					"day_of_the_month_to_run": 10,
					"hour_to_run_on": 10,
					"minute_to_run_on": 30
				},
				"description": "CM Backup",
				"end_date": "2025-09-25T18:30:00Z",
				"id": "dff97737-8686-44c6-9e0c-830d7b49b94a",
				"job_type": "NORMAL",
				"name": "test",
				"next_run_time": "2025-08-10T10:30:00Z",
				"request_timeout_period": 30,
				"schedule_type": "DAY_AND_TIME_OF_THE_MONTH",
				"start_date": "2025-08-10T10:30:00Z",
				"status": "ENABLED",
				"tag": "cm-backup-tag",
				"topic": "cm-backup"
			}`)
	})

	mux.HandleFunc("/api/v1/system/backups/schedule/dff97737-8686-44c6-9e0c-830d7b49b94a", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, ``)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCMScheduledBackupUpdate1Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "id", "dff97737-8686-44c6-9e0c-830d7b49b94a"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "type", "DayAndTimeOfTheMonth"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "day_of_the_month_to_run", "10"),
				),
			},
		},
	})
}

func TestUnitCMBackupScheduleUnitTC3Resource(t *testing.T) {
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

	// post call
	mux.HandleFunc("/api/v1/system/backups/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/backups/schedule"}},"id":"dff97737-8686-44c6-9e0c-830d7b49b94a","message":"Task is scheduled.First Schedule will be on 2025-08-26 10:30:00 +0000 UTC","path":"/api/v1/system/backups/schedule/dff97737-8686-44c6-9e0c-830d7b49b94a"}`)
	})

	mux.HandleFunc("/api/system/v1/schedules/dff97737-8686-44c6-9e0c-830d7b49b94a", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
				"_links": {
					"self": {
						"href": "/v1/schedules/dff97737-8686-44c6-9e0c-830d7b49b94a"
					}
				},
				"basic_with_interval": {
					"interval_to_run": 24,
					"interval_unit": "HOUR"
				},
				"description": "CM Backup",
				"id": "dff97737-8686-44c6-9e0c-830d7b49b94a",
				"job_type": "NORMAL",
				"name": "test",
				"next_run_time": "2025-08-25T18:30:00Z",
				"request_timeout_period": 30,
				"schedule_type": "BASIC_WITH_INTERVAL",
				"start_date": "2025-08-25T18:30:00Z",
				"status": "ENABLED",
				"tag": "cm-backup-tag",
				"topic": "cm-backup"
			}`)
	})

	mux.HandleFunc("/api/v1/system/backups/schedule/dff97737-8686-44c6-9e0c-830d7b49b94a", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, ``)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCMScheduledBackupUpdate2Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "id", "dff97737-8686-44c6-9e0c-830d7b49b94a"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "type", "BasicWithInterval"),
				),
			},
		},
	})
}

func TestUnitCMBackupScheduleUnitTC4Resource(t *testing.T) {
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

	mux.HandleFunc("/api/v1/system/backups/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/backups/schedule"}},"id":"dff97737-8686-44c6-9e0c-830d7b49b94a","message":"Task is scheduled.First Schedule will be on 2025-08-26 10:30:00 +0000 UTC","path":"/api/v1/system/backups/schedule/dff97737-8686-44c6-9e0c-830d7b49b94a"}`)

	})

	var getCount = 0
	mux.HandleFunc("/api/system/v1/schedules/dff97737-8686-44c6-9e0c-830d7b49b94a", func(w http.ResponseWriter, r *http.Request) {
		if getCount <= 1 {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{
				"_links": {
					"self": {
						"href": "/v1/schedules/dff97737-8686-44c6-9e0c-830d7b49b94a"
					}
				},
				"day_and_time_of_month": {
					"day_of_the_month_to_run": 10,
					"hour_to_run_on": 10,
					"minute_to_run_on": 30
				},
				"description": "CM Backup",
				"end_date": "2025-09-25T18:30:00Z",
				"id": "dff97737-8686-44c6-9e0c-830d7b49b94a",
				"job_type": "NORMAL",
				"name": "test",
				"next_run_time": "2025-08-10T10:30:00Z",
				"request_timeout_period": 30,
				"schedule_type": "DAY_AND_TIME_OF_THE_MONTH",
				"start_date": "2025-08-10T10:30:00Z",
				"status": "ENABLED",
				"tag": "cm-backup-tag",
				"topic": "cm-backup"
			}`)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{
				"_links": {
					"self": {
						"href": "/v1/schedules/dff97737-8686-44c6-9e0c-830d7b49b94a"
					}
				},
				"basic_with_interval": {
					"interval_to_run": 24,
					"interval_unit": "HOUR"
				},
				"description": "CM Backup",
				"id": "dff97737-8686-44c6-9e0c-830d7b49b94a",
				"job_type": "NORMAL",
				"name": "test",
				"next_run_time": "2025-08-25T18:30:00Z",
				"request_timeout_period": 30,
				"schedule_type": "BASIC_WITH_INTERVAL",
				"start_date": "2025-08-25T18:30:00Z",
				"status": "ENABLED",
				"tag": "cm-backup-tag",
				"topic": "cm-backup"
			}`)
		}
		getCount++
	})

	mux.HandleFunc("/api/v1/system/backups/schedule/dff97737-8686-44c6-9e0c-830d7b49b94a", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, ``)
		} else if r.Method == "PUT" {
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/v1/backups/schedule"}},"id":"dff97737-8686-44c6-9e0c-830d7b49b94a","message":"Task is scheduled.First Schedule will be on 2025-08-26 10:30:00 +0000 UTC","path":"/api/v1/system/backups/schedule/dff97737-8686-44c6-9e0c-830d7b49b94a"}`)
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
				Config: testAccCMScheduledBackupUpdate1Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "id", "dff97737-8686-44c6-9e0c-830d7b49b94a"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "type", "DayAndTimeOfTheMonth"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "day_of_the_month_to_run", "10"),
				),
			},
			{
				Config: testAccCMScheduledBackupUpdate2Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "id", "dff97737-8686-44c6-9e0c-830d7b49b94a"),
					resource.TestCheckResourceAttr("bigipnext_cm_backup_restore.test", "type", "BasicWithInterval"),
				),
			},
		},
	})
}

const testAccCMBackupConfig = `
resource "bigipnext_cm_backup_restore" "test" {
	name = "Test"
	encryption_password = "F5site02@1234"
  backup = true
}
`
const testAccCMBackupUpdateConfig = `
resource "bigipnext_cm_backup_restore" "test" {
	name = "Test"
	encryption_password = "F5site02@1233"
  backup = true
}
`

const testAccCMRestoreConfig = `
resource "bigipnext_cm_backup_restore" "test" {
	name = "Backup-20241007-075720_L_20.3.1-0.5.0_8.tgz"
	encryption_password = "F5site02@123"
  backup = false
}
`
const testAccCMRestoreUpdateConfig = `
resource "bigipnext_cm_backup_restore" "test" {
	name = "Backup-20241007-075720_L_20.3.1-0.5.0_8.tgz"
	encryption_password = "F5site02@1234"
  backup = false
}
`

const testAccCMScheduledBackupConfig = `
resource "bigipnext_cm_backup_restore" "test" {
	name = "Test"
	encryption_password = "F5site02@1234"
  	backup = true
	schedule = {
		start_at = "2025-08-25T18:30:00Z"
    	end_at = "2025-09-25T18:30:00Z"
	}
	frequency = "Weekly"
	days_of_the_week_to_run = [1,2,3,4]
	// day_of_the_month_to_run = 10
}
`
const testAccCMScheduledBackupUpdate1Config = `
resource "bigipnext_cm_backup_restore" "test" {
	name = "Test"
	encryption_password = "F5site02@1234"
  	backup = true
	schedule = {
		start_at = "2025-08-25T18:30:00Z"
    	end_at = "2025-09-25T18:30:00Z"
	}
	frequency = "Monthly"
	day_of_the_month_to_run = 10
}
`
const testAccCMScheduledBackupUpdate2Config = `
resource "bigipnext_cm_backup_restore" "test" {
	name = "Test"
	encryption_password = "F5site02@1234"
  	backup = true
	schedule = {
		start_at = "2025-08-25T18:30:00Z"
	}
	frequency = "Daily"
}
`
