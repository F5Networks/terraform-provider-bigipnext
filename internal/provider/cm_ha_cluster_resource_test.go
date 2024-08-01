package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"gitswarm.f5net.com/terraform-providers/bigipnext"
)

func TestAccCMHAClusterTC(t *testing.T) {
	control_node := os.Getenv("BIGIPNEXT_HOST")
	id := fmt.Sprintf("central-manager-server-%s", control_node)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCMHAClusterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_ha_cluster.cm_ha_2_nodes", "server_nodes.0", "10.146.164.150"),
					resource.TestCheckResourceAttr("bigipnext_cm_ha_cluster.cm_ha_2_nodes", "id", id),
				),
				Destroy: false,
			},
			{
				Config: testAccCMHAClusterUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_ha_cluster.cm_ha_2_nodes", "server_nodes.0", "10.146.164.150"),
					resource.TestCheckResourceAttr("bigipnext_cm_ha_cluster.cm_ha_2_nodes", "server_nodes.1", "10.146.165.89"),
					resource.TestCheckResourceAttr("bigipnext_cm_ha_cluster.cm_ha_2_nodes", "id", id),
				),
				Destroy: false,
			},
		},
	})
}

func TestUnitCMHAClusterCreate(t *testing.T) {
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

	afterDeleteCall := false

	mux.HandleFunc("/api/v1/system/infra/nodes", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			payload := make([]byte, r.ContentLength)

			_, err := r.Body.Read(payload)

			if err != nil && err.Error() != "EOF" {
				t.Errorf("Error reading request body: %v", err)
			}

			defer r.Body.Close()
			var nodes []bigipnext.CMHANodes

			err = json.Unmarshal(payload, &nodes)
			if err != nil {
				t.Errorf("Error unmarshalling request body: %v", err)
			}

			if nodes[0].NodeAddress != "10.146.164.150" {
				t.Errorf("Expected node address to be 10.146.164.150, got %s", nodes[0].NodeAddress)
			}
			if nodes[0].Username != "admin" {
				t.Errorf("Expected username to be admin, got %s", nodes[0].Username)
			}
			if nodes[0].Password != "F5site02@123" {
				t.Errorf("Expected password to be F5site02@123, got %s", nodes[0].Password)
			}
			if nodes[0].Fingerprint != "c2b2472624cd7f3a053c4f6e0bd4b322" {
				t.Errorf("Expected fingerprint to be c2b2472624cd7f3a053c4f6e0bd4b322, got %s", nodes[0].Fingerprint)
			}

			if nodes[1].NodeAddress != "10.146.165.89" {
				t.Errorf("Expected node address to be 10.146.165.89, got %s", nodes[1].NodeAddress)
			}
			if nodes[1].Username != "admin" {
				t.Errorf("Expected username to be admin, got %s", nodes[1].Username)
			}
			if nodes[1].Password != "F5site02@123" {
				t.Errorf("Expected password to be F5site02@123, got %s", nodes[1].Password)
			}
			if nodes[1].Fingerprint != "c2b2472624cd7f3a053c4f6e0bd4b322" {
				t.Errorf("Expected fingerprint to be c2b2472624cd7f3a053c4f6e0bd4b322, got %s", nodes[1].Fingerprint)
			}

			postResp, _ := os.ReadFile("fixtures/cm_ha_nodes_post.json")
			fmt.Fprint(w, postResp)
		}

		if r.Method == http.MethodGet {
			getResp, _ := os.ReadFile("fixtures/cm_ha_nodes_get.json")
			fmt.Fprint(w, getResp)
		}

		if r.Method == http.MethodGet && afterDeleteCall {
			getAfterDeleteResp, _ := os.ReadFile("fixtures/cm_ha_nodes_get_after_delete.json")
			fmt.Fprint(w, getAfterDeleteResp)
		}
	})

	mux.HandleFunc("/api/v1/system/infra/nodes/central-manager-10-146-165-89", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			afterDeleteCall = true
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, "node central-manager-10-146-165-89 is in process of deregistration.")
		}
	})

	mux.HandleFunc("/api/v1/system/infra/nodes/central-manager-10-146-164-150", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			afterDeleteCall = true
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, "node central-manager-10-146-164-150 is in process of deregistration.")
		}
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testUnitCMHAClusterConfig,
			},
		},
	})
}

func TestUnitCMHAClusterUpdate(t *testing.T) {
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

	afterDeleteCall := false
	updateCall := false
	mux.HandleFunc("/api/v1/system/infra/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postResp, _ := os.ReadFile("fixtures/cm_ha_nodes_post.json")
			fmt.Fprint(w, postResp)
			updateCall = true
		}

		if r.Method == http.MethodPost && updateCall {
			postResp, _ := os.ReadFile("fixtures/cm_ha_nodes_post_update.json")
			fmt.Fprint(w, postResp)
		}

		if r.Method == http.MethodGet {
			getResp, _ := os.ReadFile("fixtures/cm_ha_nodes_get.json")
			fmt.Fprint(w, getResp)
		}

		if r.Method == http.MethodGet && afterDeleteCall {
			getAfterDeleteResp, _ := os.ReadFile("fixtures/cm_ha_nodes_get_after_delete.json")
			fmt.Fprint(w, getAfterDeleteResp)
		}
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:  testUnitCMHAClusterConfig,
				Destroy: false,
			},
			{
				Config:  testUnitCMHAClusterUpdateConfig,
				Destroy: false,
			},
		},
	})
}

func TestUnitCMHAClusterGetFingerprint(t *testing.T) {
	fp1, err := getFingerPrint("8.8.8.8")
	if len(fp1) != 64 {
		t.Errorf("Expected fingerprint length to be 64, got %d", len(fp1))
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, err = getFingerPrint("a.b.c.d")
	if err == nil {
		t.Errorf("Expected error for a bogus IP address, got nil")
	}
}

func TestUnitCMHAClusterGetServerAndAgentNodes(t *testing.T) {
	nodes := []bigipnext.CMHANodesStatus{
		{
			Spec: bigipnext.CMHANodeSpec{
				NodeAddress: "1.2.3.4",
				NodeType:    "server",
			},
		},
		{
			Spec: bigipnext.CMHANodeSpec{
				NodeAddress: "5.6.7.8",
				NodeType:    "agent",
			},
		},
	}

	serverNodes, agentNodes := getServerAndAgentNodes(nodes)
	if len(serverNodes) != 1 {
		t.Errorf("Expected server nodes length to be 1, got %d", len(serverNodes))
	}
	if len(agentNodes) != 1 {
		t.Errorf("Expected agent nodes length to be 1, got %d", len(agentNodes))
	}

	if serverNodes[0] != "1.2.3.4" {
		t.Errorf("Expected server node address to be 1.2.3.4, got %s", nodes[0].Spec.NodeAddress)
	}
	if agentNodes[0] != "5.6.7.8" {
		t.Errorf("Expected agent node address to be 5.6.7.8, got %s", nodes[1].Spec.NodeAddress)
	}
}

const testAccCMHAClusterConfig = `
resource "bigipnext_cm_ha_cluster" "cm_ha_2_nodes" {
  nodes = [
    {
        node_ip  = "10.146.164.150"
        username = "admin",
        password = "F5site02@123"
    }
  ]
}
`

const testAccCMHAClusterUpdateConfig = `
resource "bigipnext_cm_ha_cluster" "cm_ha_2_nodes" {
  nodes = [
    {
        node_ip  = "10.146.164.150"
        username = "admin",
        password = "F5site02@123"
    },
    {
        node_ip  = "10.146.165.89"
        username = "admin",
        password = "F5site02@123"
    }
  ]
}
`

const testUnitCMHAClusterConfig = `
resource "bigipnext_cm_ha_cluster" "cm_ha_2_nodes" {
  nodes = [
    {
        node_ip     = "10.146.164.150"
        username    = "admin",
        password    = "F5site02@123"
		fingerprint = "c2b2472624cd7f3a053c4f6e0bd4b322"
    },
    {
        node_ip     = "10.146.165.89"
        username    = "admin",
        password    = "F5site02@123"
		fingerprint = "c2b2472624cd7f3a053c4f6e0bd4b322"
    }
  ]
}
`

const testUnitCMHAClusterUpdateConfig = `
resource "bigipnext_cm_ha_cluster" "cm_ha_2_nodes" {
  nodes = [
    {
        node_ip     = "10.146.164.150"
        username    = "admin",
        password    = "F5site02@123"
		fingerprint = "c2b2472624cd7f3a053c4f6e0bd4b322"
    },
    {
        node_ip     = "10.146.165.89"
        username    = "admin",
        password    = "F5site02@123"
		fingerprint = "c2b2472624cd7f3a053c4f6e0bd4b322"
    },
	{
		node_ip     = "12.34.56.77"
		username    = "admin",
		password    = "F5site02@123"
		fingerprint = "c2b2472624cd7f3a053c4f6e0bd4b322"
	}
  ]
}
`
