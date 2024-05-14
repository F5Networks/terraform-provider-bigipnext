package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNextCMWAFReportCreateTC1Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNextCMWAFReportResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "name", "sample"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "description", "WAF security report description"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "request_type", "illegal"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "time_frame_in_days", "7"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "top_level", "5"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "categories.0.name", "Source IPs"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "scope.entity", "policies"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "scope.all", "false"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_waf_report.sample", "scope.names.*", "test_policy"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "user_defined", "true"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "created_by", "admin"),
				),
			},
			{
				Config: testAccNextCMWAFReportResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "name", "sample"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "description", "WAF security report description Updated"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "request_type", "alerted"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "time_frame_in_days", "30"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "top_level", "10"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "categories.0.name", "Source IPs"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "categories.1.name", "Geolocations"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "scope.entity", "applications"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "scope.all", "false"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_waf_report.sample", "scope.names.*", "test_policy"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_waf_report.sample", "scope.names.*", "test"),
				),
			},
		},
	})
}

func TestUnitCMWAFReportCreateUnitTC1Resource(t *testing.T) {
	testAccPreUnitCheck(t)
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"access_token": "Zj7421sAIkQF4YeRfKdMvQv+rJ1BvoWAFAycmKkQDTmxpMj5vK4o79Imsed6KBVzOxoeHRGwb08wEu7rrdeRU9HXXmukeNyaFRyXfYdxWyd2GNRYY+uqlFHiNb+2kr1UC114AFv65rRZ7tplfpFcJL39ETxSs1vjhcsBT+BClEEUP48fuYQv3htSQvNbs82i/DHYU9FYn2diUuuoOVuPhHj81Q7/Rk5FFea1NA0ahYO7GHckAtFWM/UlPgA9yyDFQ6dzyX6OuE9eppR+6/VY1t55oPxFMFdL0wkq8aulGxFWnHwgJ6EYx5KkuQvIDQOguXAb7+C+ffH2fWWh7QPnCQjddVssrbwUpbXZDgMptSyOWul6MudTVAbHfyJNMxaj159HJUv/NhrGnfu1S7A9++aYnTPJsGgqSXnASkNOqOiZr/wiRyCxzx9VF1kqgzCN8Mc+U8y2EHDveix7nF3BiQtIneYrt2ycGlqZFXkfRnQCYiOOWcAvvz2eTKYoZOPPXU9TCI4WzWnOKCGQYYRvt2uy74IOeBSexMt03EU3GA==",
			"refresh_token": "Zj7421sAIkQF4YeRfKdMvQv+rJ1BvoWAFAycmKkQDTmxpMj5vK4o79Imsed6KBVzOxoeHRGwb08wEu7rrdeRU9HXXmukeNyaFRyXfYdxWyd2GNRYY+uqlBngG6u2kr1UC114AFv65rRZ7tplfpFcJL39ETxSs1vjhcsBT+BClEEUP48fuYQv3htSQvNbs82i/DHYU9FYn2diUuuoOVuPhHj81Q7/Rk5FFea1NA0ahYDmD3AoAtFVMK0zBABy3A7NRadzyXCPglteppR9+fRY18xtnsZHJW1LnwoE2PKlHC9WnHwgJ6EYx5KkuQvIDQOguXAb7+C+ffH2fWWh7QPnCQjddVssrbwUpbXZDgMptSyOWul6MudTVAbHfyJNMxaj159HJUv/NhrGnfu1S7A9++aYnTPJsGgqSVvAXU9Hk6mAq6cvVQKFw2tCK3Vt6SS0tpsHJ46/BiXwFyEQs2fuxx52tY4Bs3OoNlbkVTcfQFYMRA7QXvQA+QeHQ60BS+H8EFd6L6sU6UP1LEKDfYNJH3fUAVrsmPgBN5H8G67wTg==",
			"user_id": "6dd0d482-267e-4916-b524-ee8e5dd1c78"
		}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/security/waf/reports", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/security/waf/reports//112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac"}},"categories":[{"name":"Source IPs"}],"created_by":"admin","description":"WAF security report description","id":"112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac","name":"sample","request_type":"illegal","scope":{"entity":"policies","all":false,"names":["test_policy"]},"time_frame_in_days":7,"top_level":5}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/security/waf/reports/112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
			_, _ = fmt.Fprintf(w, ``)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/security/waf/reports/112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac"}},"categories":[{"name":"Source IPs"}],"created_by":"admin","description":"WAF security report description","id":"112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac","name":"sample","request_type":"illegal","scope":{"entity":"policies","all":false,"names":["test_policy"]},"time_frame_in_days":7,"top_level":5,"user_defined":true}`)
		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNextCMWAFReportResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestUnitCMWAFReportCreateUnitTC2Resource(t *testing.T) {
	var count = 0
	testAccPreUnitCheck(t)
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"access_token": "Zj7421sAIkQF4YeRfKdMvQv+rJ1BvoWAFAycmKkQDTmxpMj5vK4o79Imsed6KBVzOxoeHRGwb08wEu7rrdeRU9HXXmukeNyaFRyXfYdxWyd2GNRYY+uqlFHiNb+2kr1UC114AFv65rRZ7tplfpFcJL39ETxSs1vjhcsBT+BClEEUP48fuYQv3htSQvNbs82i/DHYU9FYn2diUuuoOVuPhHj81Q7/Rk5FFea1NA0ahYO7GHckAtFWM/UlPgA9yyDFQ6dzyX6OuE9eppR+6/VY1t55oPxFMFdL0wkq8aulGxFWnHwgJ6EYx5KkuQvIDQOguXAb7+C+ffH2fWWh7QPnCQjddVssrbwUpbXZDgMptSyOWul6MudTVAbHfyJNMxaj159HJUv/NhrGnfu1S7A9++aYnTPJsGgqSXnASkNOqOiZr/wiRyCxzx9VF1kqgzCN8Mc+U8y2EHDveix7nF3BiQtIneYrt2ycGlqZFXkfRnQCYiOOWcAvvz2eTKYoZOPPXU9TCI4WzWnOKCGQYYRvt2uy74IOeBSexMt03EU3GA==",
			"refresh_token": "Zj7421sAIkQF4YeRfKdMvQv+rJ1BvoWAFAycmKkQDTmxpMj5vK4o79Imsed6KBVzOxoeHRGwb08wEu7rrdeRU9HXXmukeNyaFRyXfYdxWyd2GNRYY+uqlBngG6u2kr1UC114AFv65rRZ7tplfpFcJL39ETxSs1vjhcsBT+BClEEUP48fuYQv3htSQvNbs82i/DHYU9FYn2diUuuoOVuPhHj81Q7/Rk5FFea1NA0ahYDmD3AoAtFVMK0zBABy3A7NRadzyXCPglteppR9+fRY18xtnsZHJW1LnwoE2PKlHC9WnHwgJ6EYx5KkuQvIDQOguXAb7+C+ffH2fWWh7QPnCQjddVssrbwUpbXZDgMptSyOWul6MudTVAbHfyJNMxaj159HJUv/NhrGnfu1S7A9++aYnTPJsGgqSVvAXU9Hk6mAq6cvVQKFw2tCK3Vt6SS0tpsHJ46/BiXwFyEQs2fuxx52tY4Bs3OoNlbkVTcfQFYMRA7QXvQA+QeHQ60BS+H8EFd6L6sU6UP1LEKDfYNJH3fUAVrsmPgBN5H8G67wTg==",
			"user_id": "6dd0d482-267e-4916-b524-ee8e5dd1c78"
		}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/security/waf/reports", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/security/waf/reports//112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac"}},"categories":[{"name":"Source IPs"}],"created_by":"admin","description":"WAF security report description","id":"112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac","name":"sample","request_type":"illegal","scope":{"entity":"policies","all":false,"names":["test_policy"]},"time_frame_in_days":7,"top_level":5}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/security/waf/reports/112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
			_, _ = fmt.Fprintf(w, ``)
		} else if r.Method == "PUT" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/security/waf/reports/112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac"}},"categories":[{"name":"Source IPs"},{"name":"Geolocations"}],"description":"WAF security report description Updated","id":"112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac","name":"sample","request_type":"alerted","scope":{"entity":"applications","all":false,"names":["test_policy","test"]},"time_frame_in_days":30,"top_level":10}`)
		} else {
			if count == 0 || count == 1 || count == 2 {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/security/waf/reports/112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac"}},"categories":[{"name":"Source IPs"}],"created_by":"admin","description":"WAF security report description","id":"112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac","name":"sample","request_type":"illegal","scope":{"entity":"policies","all":false,"names":["test_policy"]},"time_frame_in_days":7,"top_level":5,"user_defined":true}`)
			} else {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/security/waf/reports/112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac"}},"categories":[{"name":"Source IPs"},{"name":"Geolocations"}],"created_by":"admin","description":"WAF security report description Updated","id":"112632cc-d9c9-4dc2-b2b3-1ddcc11a14ac","name":"sample","request_type":"alerted","scope":{"entity":"applications","all":false,"names":["test_policy","test"]},"time_frame_in_days":30,"top_level":10,"user_defined":true}`)
			}
			count++
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
				Config: testAccNextCMWAFReportResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "name", "sample"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "description", "WAF security report description"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "request_type", "illegal"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "time_frame_in_days", "7"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "top_level", "5"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "categories.0.name", "Source IPs"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "scope.entity", "policies"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "scope.all", "false"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_waf_report.sample", "scope.names.*", "test_policy"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "user_defined", "true"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "created_by", "admin"),
				),
			},
			{
				Config: testAccNextCMWAFReportResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "name", "sample"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "description", "WAF security report description Updated"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "request_type", "alerted"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "time_frame_in_days", "30"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "top_level", "10"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "categories.0.name", "Source IPs"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "categories.1.name", "Geolocations"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "scope.entity", "applications"),
					resource.TestCheckResourceAttr("bigipnext_cm_waf_report.sample", "scope.all", "false"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_waf_report.sample", "scope.names.*", "test_policy"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_waf_report.sample", "scope.names.*", "test"),
				),
			},
		},
	})
}

const testAccNextCMWAFReportResourceConfig = `
resource "bigipnext_cm_waf_report" "sample" {
	name      = "sample"
	description = "WAF security report description"
	request_type = "illegal"
	time_frame_in_days = 7
	top_level = 5
	categories = [{"name" : "Source IPs"}]
	scope = {
        entity = "policies",
        all = false
		names = ["test_policy"]
    }
}`

const testAccNextCMWAFReportResourceUpdateConfig = `
resource "bigipnext_cm_waf_report" "sample" {
	name      = "sample"
	description = "WAF security report description Updated"
	request_type = "alerted"
	time_frame_in_days = 30
	top_level = 10
	categories = [{"name" : "Source IPs"} , {"name" : "Geolocations"}]
	scope = {
        entity = "applications",
        all = false
		names = ["test_policy", "test"]
    }
}`
