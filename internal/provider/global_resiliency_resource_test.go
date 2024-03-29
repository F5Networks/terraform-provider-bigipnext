package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNextCMGlobalResiliencyCreateTC1Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMGlobalResiliencyResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "name", "sample2"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "dns_listener_name", "dln"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "dns_listener_port", "10"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_global_resiliency.sample2", "protocols.*", "tcp"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "instances.0.address", "10.145.71.115"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "instances.0.dns_listener_address", "2.2.2.3"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "instances.0.group_sync_address", "10.10.1.2/24"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "instances.0.hostname", "big-ip-next"),
				),
			},
			{
				Config: testAccNextCMGlobalResiliencyResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "name", "sample2"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "dns_listener_name", "dln"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "dns_listener_port", "10"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_global_resiliency.sample2", "protocols.*", "tcp"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "instances.1.address", "10.144.140.81"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "instances.1.dns_listener_address", "1.1.1.3"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "instances.1.group_sync_address", "192.168.1.1/24"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample2", "instances.1.hostname", "demovm01-lab-local"),
				),
			},
		},
	})
}

func TestUnitCMGlobalResiliencyCreateUnitTC1Resource(t *testing.T) {
	var createCount = 0
	testAccPreUnitCheck(t)
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"access_token": "Zj7421sAIkQF4YeRfKdMvQv+rJ1BvoWAFAycmKkQDTmxpMj5vK4o79Imsed6KBVzOxoeHRGwb08wEu7rrdeRU9HXXmukeNyaFRyXfYdxWyd2GNRYY+uqlFHiNb+2kr1UC114AFv65rRZ7tplfpFcJL39ETxSs1vjhcsBT+BClEEUP48fuYQv3htSQvNbs82i/DHYU9FYn2diUuuoOVuPhHj81Q7/Rk5FFea1NA0ahYO7GHckAtFWM/UlPgA9yyDFQ6dzyX6OuE9eppR+6/VY1t55oPxFMFdL0wkq8aulGxFWnHwgJ6EYx5KkuQvIDQOguXAb7+C+ffH2fWWh7QPnCQjddVssrbwUpbXZDgMptSyOWul6MudTVAbHfyJNMxaj159HJUv/NhrGnfu1S7A9++aYnTPJsGgqSXnASkNOqOiZr/wiRyCxzx9VF1kqgzCN8Mc+U8y2EHDveix7nF3BiQtIneYrt2ycGlqZFXkfRnQCYiOOWcAvvz2eTKYoZOPPXU9TCI4WzWnOKCGQYYRvt2uy74IOeBSexMt03EU3GA==",
			"refresh_token": "Zj7421sAIkQF4YeRfKdMvQv+rJ1BvoWAFAycmKkQDTmxpMj5vK4o79Imsed6KBVzOxoeHRGwb08wEu7rrdeRU9HXXmukeNyaFRyXfYdxWyd2GNRYY+uqlBngG6u2kr1UC114AFv65rRZ7tplfpFcJL39ETxSs1vjhcsBT+BClEEUP48fuYQv3htSQvNbs82i/DHYU9FYn2diUuuoOVuPhHj81Q7/Rk5FFea1NA0ahYDmD3AoAtFVMK0zBABy3A7NRadzyXCPglteppR9+fRY18xtnsZHJW1LnwoE2PKlHC9WnHwgJ6EYx5KkuQvIDQOguXAb7+C+ffH2fWWh7QPnCQjddVssrbwUpbXZDgMptSyOWul6MudTVAbHfyJNMxaj159HJUv/NhrGnfu1S7A9++aYnTPJsGgqSVvAXU9Hk6mAq6cvVQKFw2tCK3Vt6SS0tpsHJ46/BiXwFyEQs2fuxx52tY4Bs3OoNlbkVTcfQFYMRA7QXvQA+QeHQ60BS+H8EFd6L6sU6UP1LEKDfYNJH3fUAVrsmPgBN5H8G67wTg==",
			"user_id": "6dd0d482-267e-4916-b524-ee8e5dd1c78"
		}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/gslb/gr-groups", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/gslb/gr-groups/"}},"path":"/v1/gslb/gr-groups/efcd6a7d-b3e3-4220-9e08-c9604ff9222e"}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/gslb/gr-groups/efcd6a7d-b3e3-4220-9e08-c9604ff9222e", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusAccepted)
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/gslb/gr-groups/efcd6a7d-b3e3-4220-9e08-c9604ff9222e"}},"id":"efcd6a7d-b3e3-4220-9e08-c9604ff9222e","message":"Deleting the Global Resiliency Group"}`)
		} else {

			if createCount == 0 || createCount == 1 {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/gslb/gr-groups/efcd6a7d-b3e3-4220-9e08-c9604ff9222e"}},"dns_listener_name":"dln","dns_listener_port":10,"fqdn_count":{"active":0,"disabled":0},"health":"GOOD","id":"efcd6a7d-b3e3-4220-9e08-c9604ff9222e","instance_count":1,"instances":[{"id":"c1fd14d5-97e2-4cbd-9509-bc3df1c0aafb","hostname":"big-ip-next","address":"10.145.71.115","dns_listener_address":"2.2.2.3","group_sync_address":"10.10.1.2/24","health":"GOOD"}],"name":"sample4","protocols":["tcp"],"status":"DEPLOYED"}
								`)
			} else {
				w.WriteHeader(http.StatusNotFound)
				_, _ = fmt.Fprintf(w, `{"status":404,"message":"Requested Global Resiliency Group not found for the given group id"}`)
			}
			createCount++
		}
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNextCMGlobalResiliencyResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestUnitCMGlobalResiliencyCreateUnitTC2Resource(t *testing.T) {
	var createCount = 0
	var updateCount = 0
	testAccPreUnitCheck(t)
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"access_token": "Zj7421sAIkQF4YeRfKdMvQv+rJ1BvoWAFAycmKkQDTmxpMj5vK4o79Imsed6KBVzOxoeHRGwb08wEu7rrdeRU9HXXmukeNyaFRyXfYdxWyd2GNRYY+uqlFHiNb+2kr1UC114AFv65rRZ7tplfpFcJL39ETxSs1vjhcsBT+BClEEUP48fuYQv3htSQvNbs82i/DHYU9FYn2diUuuoOVuPhHj81Q7/Rk5FFea1NA0ahYO7GHckAtFWM/UlPgA9yyDFQ6dzyX6OuE9eppR+6/VY1t55oPxFMFdL0wkq8aulGxFWnHwgJ6EYx5KkuQvIDQOguXAb7+C+ffH2fWWh7QPnCQjddVssrbwUpbXZDgMptSyOWul6MudTVAbHfyJNMxaj159HJUv/NhrGnfu1S7A9++aYnTPJsGgqSXnASkNOqOiZr/wiRyCxzx9VF1kqgzCN8Mc+U8y2EHDveix7nF3BiQtIneYrt2ycGlqZFXkfRnQCYiOOWcAvvz2eTKYoZOPPXU9TCI4WzWnOKCGQYYRvt2uy74IOeBSexMt03EU3GA==",
			"refresh_token": "Zj7421sAIkQF4YeRfKdMvQv+rJ1BvoWAFAycmKkQDTmxpMj5vK4o79Imsed6KBVzOxoeHRGwb08wEu7rrdeRU9HXXmukeNyaFRyXfYdxWyd2GNRYY+uqlBngG6u2kr1UC114AFv65rRZ7tplfpFcJL39ETxSs1vjhcsBT+BClEEUP48fuYQv3htSQvNbs82i/DHYU9FYn2diUuuoOVuPhHj81Q7/Rk5FFea1NA0ahYDmD3AoAtFVMK0zBABy3A7NRadzyXCPglteppR9+fRY18xtnsZHJW1LnwoE2PKlHC9WnHwgJ6EYx5KkuQvIDQOguXAb7+C+ffH2fWWh7QPnCQjddVssrbwUpbXZDgMptSyOWul6MudTVAbHfyJNMxaj159HJUv/NhrGnfu1S7A9++aYnTPJsGgqSVvAXU9Hk6mAq6cvVQKFw2tCK3Vt6SS0tpsHJ46/BiXwFyEQs2fuxx52tY4Bs3OoNlbkVTcfQFYMRA7QXvQA+QeHQ60BS+H8EFd6L6sU6UP1LEKDfYNJH3fUAVrsmPgBN5H8G67wTg==",
			"user_id": "6dd0d482-267e-4916-b524-ee8e5dd1c78"
		}`)
	})

	mux.HandleFunc("/api/v1/spaces/default/gslb/gr-groups/efcd6a7d-b3e3-4220-9e08-c9604ff9222e", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusAccepted)
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/gslb/gr-groups/efcd6a7d-b3e3-4220-9e08-c9604ff9222e"}},"id":"efcd6a7d-b3e3-4220-9e08-c9604ff9222e","message":"Deleting the Global Resiliency Group"}`)
		} else if r.Method == "GET" {

			if createCount == 0 || createCount == 1 || createCount == 2 {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/gslb/gr-groups/efcd6a7d-b3e3-4220-9e08-c9604ff9222e"}},"dns_listener_name":"dln","dns_listener_port":10,"fqdn_count":{"active":0,"disabled":0},"health":"GOOD","id":"efcd6a7d-b3e3-4220-9e08-c9604ff9222e","instance_count":1,"instances":[{"id":"c1fd14d5-97e2-4cbd-9509-bc3df1c0aafb","hostname":"big-ip-next","address":"10.145.71.115","dns_listener_address":"2.2.2.3","group_sync_address":"10.10.1.2/24","health":"GOOD"}],"name":"sample4","protocols":["tcp"],"status":"DEPLOYED"}
								`)
				createCount++
			} else if updateCount == 0 || updateCount == 1 {
				_, _ = fmt.Fprintf(w, `{"_links": {"self": {"href": "/api/v1/spaces/default/gslb/gr-groups/efcd6a7d-b3e3-4220-9e08-c9604ff9222e"}},"dns_listener_name": "dln","dns_listener_port": 10,"fqdn_count": {"active": 0,"disabled": 0},"health": "GOOD","id": "efcd6a7d-b3e3-4220-9e08-c9604ff9222e","instance_count": 1,"instances": [{"id": "c1fd14d5-97e2-4cbd-9509-bc3df1c0aafb","hostname": "big-ip-next","address": "10.145.71.115","dns_listener_address": "2.2.2.3","group_sync_address": "10.10.1.2/24","health": "GOOD"},{"id": "c1fd14d5-97e2-4cbd-9509-bc3df1c0aafb","hostname": "demovm01-lab-local","address": "10.144.140.81","dns_listener_address": "1.1.1.3","group_sync_address": "192.168.1.1/24","health": "GOOD"}],"name": "sample4","protocols": ["tcp"],"status": "DEPLOYED"}`)
				updateCount++
			} else {
				w.WriteHeader(http.StatusNotFound)
				_, _ = fmt.Fprintf(w, `{"status":404,"message":"Requested Global Resiliency Group not found for the given group id"}`)
			}

		} else if r.Method == "PUT" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/gslb/gr-groups/"}},"path":"/v1/gslb/gr-groups/efcd6a7d-b3e3-4220-9e08-c9604ff9222e"}`)
		}
	})

	mux.HandleFunc("/api/v1/spaces/default/gslb/gr-groups", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/gslb/gr-groups/"}},"path":"/v1/gslb/gr-groups/efcd6a7d-b3e3-4220-9e08-c9604ff9222e"}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMGlobalResiliencyResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "name", "sample4"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "dns_listener_name", "dln"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "dns_listener_port", "10"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_global_resiliency.sample4", "protocols.*", "tcp"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "instances.0.address", "10.145.71.115"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "instances.0.dns_listener_address", "2.2.2.3"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "instances.0.group_sync_address", "10.10.1.2/24"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "instances.0.hostname", "big-ip-next"),
				),
			},
			{
				Config: testAccNextCMGlobalResiliencyResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "name", "sample4"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "dns_listener_name", "dln"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "dns_listener_port", "10"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_global_resiliency.sample4", "protocols.*", "tcp"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "instances.1.address", "10.144.140.81"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "instances.1.dns_listener_address", "1.1.1.3"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "instances.1.group_sync_address", "192.168.1.1/24"),
					resource.TestCheckResourceAttr("bigipnext_cm_global_resiliency.sample4", "instances.1.hostname", "demovm01-lab-local"),
				),
			},
		},
	})
}

const testAccNextCMGlobalResiliencyResourceConfig = `
resource "bigipnext_cm_global_resiliency" "sample4" {
	name      = "sample4"
	dns_listener_name = "dln"
	dns_listener_port = 10
	protocols = ["tcp"]
	instances = [
		{
			address = "10.145.71.115"
			dns_listener_address = "2.2.2.3"
			group_sync_address = "10.10.1.2/24"
			hostname = "big-ip-next"
		}
	]
}`

const testAccNextCMGlobalResiliencyResourceUpdateConfig = `
resource "bigipnext_cm_global_resiliency" "sample4" {
	name      = "sample4"
	dns_listener_name = "dln"
	dns_listener_port = 10
	protocols = ["tcp"]
	instances = [
		{
			address = "10.145.71.115"
			dns_listener_address = "2.2.2.3"
			group_sync_address = "10.10.1.2/24"
			hostname = "big-ip-next"
		},
		{
			address = "10.144.140.81"
			dns_listener_address = "1.1.1.3"
			group_sync_address = "192.168.1.1/24"
			hostname = "demovm01-lab-local"
		}
	]
}`
