package provider

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitCMNextHA(t *testing.T) {
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

	mux.HandleFunc("/api/device/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		query, _ := url.Parse(r.URL.String())
		filterQuery := query.Query().Get("filter")

		queryStr := strings.Split(filterQuery, " ")

		if strings.Contains(queryStr[2], "10.218.33.22") {
			fmt.Fprint(w, `
			{
				"_embedded": {
					"devices": [
						{
							"id": "1"
						}
					]
				},
				"count": 1
			}
			`)
		}

		if strings.Contains(queryStr[2], "10.218.33.23") {
			fmt.Fprint(w, `
			{
				"_embedded": {
					"devices": [{
						"id": "2"
					}]
				},
				"count": 1
			}
			`)
		}

		if strings.Contains(queryStr[2], "10.218.46.27") {
			fmt.Fprint(w, `
			{
				"_embedded": {
					"devices": [{
						"id": "2",
						"mode": "HA"
					}]
				},
				"count": 1
			}
			`)
		}
	})

	mux.HandleFunc("/api/device/v1/inventory/1/ha", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Fprint(w, `
			{
				"_links":{
					"self":{
						"href":"/v1/ha-creation-tasks/task-id-1"
					}
				},
				"path":"/v1/ha-creation-tasks/task-id-1"
			}
			`)
		}
	})

	mux.HandleFunc("/api/device/v1/ha-creation-tasks/task-id-1", func(w http.ResponseWriter, r *http.Request) {
		getTaskResp, _ := os.ReadFile("fixtures/cm_next_ha_task_status.json")
		fmt.Fprint(w, string(getTaskResp))
	})

	mux.HandleFunc("/api/v1/spaces/default/instances/06aea4ed-7425-4db3-a728-2574929885d9", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprint(w, `
		{
			"_links":{
				"self":{
					"href":"/v1/deletion-tasks/delete-ha-task"
				}
			},
			"path":"/v1/deletion-tasks/delete-ha-task"
		}
		`)
	})

	mux.HandleFunc("/api/device/v1/deletion-tasks/delete-ha-task", func(w http.ResponseWriter, r *http.Request) {
		deleteTaskStatus, _ := os.ReadFile("fixtures/cm_next_ha_delete_task_status.json")
		fmt.Fprint(w, string(deleteTaskStatus))
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testUnitCMNextHAConfig,
			},
		},
	})
}

const testUnitCMNextHAConfig = `
resource "bigipnext_cm_next_ha" "test" {
  ha_name                       = "testnextha"
  ha_ip                         = "10.218.46.27"
  active_node_ip                = "10.218.33.22"
  standby_node_ip               = "10.218.33.23"
  control_plane_vlan            = "ha-cp-vlan"
  control_plane_vlan_tag        = 101
  data_plane_vlan               = "ha-dp-vlan"
  data_plane_vlan_tag           = 102
  active_node_control_plane_ip  = "10.211.44.0/8"
  standby_node_control_plane_ip = "10.211.31.0/8"
  active_node_data_plane_ip     = "10.211.98.0/8"
  standby_node_data_plane_ip    = "10.211.76.0/8"
}
`
