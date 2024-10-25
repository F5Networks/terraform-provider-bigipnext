/*
Copyright 2023 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
// Package bigipnext interacts with BIGIP-NEXT/CM systems using the OPEN API.
package bigipnext

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
)

const (
	uriAS3Root = "/api/v1/spaces/default/appsvcs"
	uriFast    = "/mgmt/shared/fast"
)

// /api/v1/spaces/default/appsvcs/documents
func (p *BigipNextCM) PostAS3DraftDocument(config string) (string, error) {
	as3DraftURL := fmt.Sprintf("%s%s%s", p.Host, uriAS3Root, "/documents")
	f5osLogger.Info("[PostAS3DraftDocument]", "URI Path", as3DraftURL)
	f5osLogger.Info("[PostAS3DraftDocument]", "Config", hclog.Fmt("%+v", config))
	respData, err := p.doCMRequest("POST", as3DraftURL, []byte(config))
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostAS3DraftDocument]", "Data::", hclog.Fmt("%+v", string(respData)))
	//{"Message":"Application service created successfully","_links":{"self":{"href":"/api/v1/spaces/default/appsvcs/documents/3a220683-6527-4443-8da7-279680c21ac5"}},"id":"3a220683-6527-4443-8da7-279680c21ac5"}
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostAS3DraftDocument]", "Document Drart", hclog.Fmt("%+v", respString["id"].(string)))
	return respString["id"].(string), nil
}

// /api/v1/spaces/default/appsvcs/documents/3a220683-6527-4443-8da7-279680c21ac5
func (p *BigipNextCM) GetAS3DraftDocument(docID string) ([]byte, error) {
	as3DraftURL := fmt.Sprintf("%s%s%s/%s", p.Host, uriAS3Root, "/documents", docID)
	f5osLogger.Info("[GetAS3DraftDocument]", "URI Path", as3DraftURL)
	respData, err := p.doCMRequest("GET", as3DraftURL, nil)
	if err != nil {
		return []byte(""), err
	}
	f5osLogger.Info("[GetAS3DraftDocument]", "Data::", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

func (p *BigipNextCM) PutAS3DraftDocument(docID, config string) error {
	as3DraftURL := fmt.Sprintf("%s%s%s/%s", p.Host, uriAS3Root, "/documents", docID)
	f5osLogger.Info("[PutAS3DraftDocument]", "URI Path", as3DraftURL)
	f5osLogger.Info("[PutAS3DraftDocument]", "Config", hclog.Fmt("%+v", config))
	respData, err := p.doCMRequest("PUT", as3DraftURL, []byte(config))
	if err != nil {
		return err
	}
	f5osLogger.Info("[PutAS3DraftDocument]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// /api/v1/spaces/default/appsvcs/documents/3a220683-6527-4443-8da7-279680c21ac5
func (p *BigipNextCM) DeleteAS3DraftDocument(docID string) error {
	as3DraftURL := fmt.Sprintf("%s%s%s/%s", p.Host, uriAS3Root, "/documents", docID)
	f5osLogger.Info("[DeleteAS3DraftDocument]", "URI Path", as3DraftURL)
	respData, err := p.doCMRequest("DELETE", as3DraftURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteAS3DraftDocument]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// https://clouddocs.f5.com/api/v1/spaces/default/appsvcs/documents/{document-id}/deployments
func (p *BigipNextCM) CMAS3DeployNext(draftID, target string, timeOut int) (string, error) {
	as3DeployUrl := fmt.Sprintf("%s%s%s/%s/%s", p.Host, uriAS3Root, "/documents", draftID, "deployments")
	f5osLogger.Info("[CMAS3DeployNext]", "URI Path", as3DeployUrl)
	as3Json := make(map[string]interface{})
	as3Json["target"] = target
	as3data, err := json.Marshal(as3Json)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[CMAS3DeployNext]", "Data::", hclog.Fmt("%+v", string(as3data)))
	respData, err := p.doCMRequest("POST", as3DeployUrl, as3data)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[CMAS3DeployNext]", "Data::", hclog.Fmt("%+v", string(respData)))
	//{ "Message": "Deployment task created successfully", "_links": { "self": { "href": "/declare/1a5a6049-8220-483a-8cbc-275a4b190d35/deployments/2ceb048a-0ee6-4a2d-8952-cd15583bb5e8" } }, "id": "2ceb048a-0ee6-4a2d-8952-cd15583bb5e8" }
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[CMAS3DeployNext]", "Deployment Task", hclog.Fmt("%+v", respString["id"].(string)))
	_, err = p.getAS3DeploymentTaskStatus(draftID, respString["id"].(string), timeOut)
	if err != nil {
		return "", err
	}
	return respString["id"].(string), nil
}

// https://clouddocs.f5.com/api/v1/spaces/default/appsvcs/documents/{document-id}/deployments/{deployment-id}
func (p *BigipNextCM) GetAS3DeploymentTaskStatus(docID, deployID string) (interface{}, error) {
	as3DeployUrl := fmt.Sprintf("%s%s%s/%s/%s/%s", p.Host, uriAS3Root, "/documents", docID, "deployments", deployID)
	f5osLogger.Info("[GetAS3DeploymentTaskStatus]", "URI Path", as3DeployUrl)
	return p.getAS3DeploymentTaskStatus(docID, deployID, 60)
}

// https://clouddocs.f5.com/api/v1/spaces/default/appsvcs/documents/{document-id}/deployments/{deployment-id}
func (p *BigipNextCM) getAS3DeploymentTaskStatus(docID, deployID string, timeOut int) (interface{}, error) {
	as3DeployUrl := fmt.Sprintf("%s%s%s/%s/%s/%s", p.Host, uriAS3Root, "/documents", docID, "deployments", deployID)
	f5osLogger.Info("[getAS3DeploymentTaskStatus]", "URI Path", as3DeployUrl)
	responseData := make(map[string]interface{})
	timeout := time.Duration(timeOut) * time.Second
	endtime := time.Now().Add(timeout)
outerfor:
	for time.Now().Before(endtime) {
		respData, err := p.doCMRequest("GET", as3DeployUrl, nil)
		if err != nil {
			return nil, err
		}
		f5osLogger.Info("[getAS3DeploymentTaskStatus]", "respData:", hclog.Fmt("%+v", string(respData)))
		err = json.Unmarshal(respData, &responseData)
		if err != nil {
			return nil, err
		}
		f5osLogger.Info("[getAS3DeploymentTaskStatus]", "Status:", hclog.Fmt("%+v", responseData["records"].([]interface{})[0].(map[string]interface{})["status"].(string)))
		for _, v := range responseData["records"].([]interface{}) {
			if v.(map[string]interface{})["status"].(string) == "failed" {
				return nil, fmt.Errorf("%v", v.(map[string]interface{})["failure_reason"].(string))
			}
			if v.(map[string]interface{})["status"].(string) != "completed" {
				time.Sleep(5 * time.Second)
				continue
			} else {
				break outerfor
			}
		}
	}
	for _, v := range responseData["response"].(map[string]interface{})["results"].([]interface{}) {
		if v.(map[string]interface{})["message"].(string) == "failed" {
			tenantName := v.(map[string]interface{})["tenant"].(string)
			return nil, fmt.Errorf("%v deployment failed", tenantName)
		}
	}
	f5osLogger.Info("[getAS3DeploymentTaskStatus]", "Response Result:", hclog.Fmt("%+v", responseData["response"]))
	// .(map[string]interface{})["results"].([]interface{})[0].(map[string]interface{})["status"].(string)))
	// if responseData["records"].([]interface{})[0].(map[string]interface{})["status"].(string) != "completed" {
	// 	return nil, fmt.Errorf("AS3 service deployment failed with :%+v", responseData["records"].([]interface{})[0].(map[string]interface{})["status"].(string))
	// }
	byteData, err := json.Marshal(responseData["request"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	// appData := strings.Join([]string{strings.TrimSpace(string(byteData))}, "")
	return string(byteData), nil
}

// /api/v1/spaces/default/appsvcs/documents/83ff823d-477c-4666-a4c7-6b0563bb7be6/deployments/f1f55f4b-5bad-4f67-8ac2-83551502a7c8
func (p *BigipNextCM) DeleteAS3DeploymentTask(docID string) error {
	// as3DeployUrl := fmt.Sprintf("%s%s%s/%s/%s/%s", p.Host, uriAS3Root, "/documents", docID, "deployments", deployID)
	as3DeployUrl := fmt.Sprintf("%s%s%s/%s", p.Host, uriAS3Root, "/documents", docID)
	f5osLogger.Info("[DeleteAS3DeploymentTask]", "URI Path", as3DeployUrl)
	respData, err := p.doCMRequest("DELETE", as3DeployUrl, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteAS3DeploymentTask]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// create Get request to get Fast openapi sepcification
func (p *BigipNextCM) GetFastSpecificationOpenAPI() error {
	fastURL := fmt.Sprintf("%s%s%s", p.Host, uriFast, uriOpenAPI)
	f5osLogger.Info("[GetFastSpecificationOpenAPI]", "URI Path", fastURL)
	respData, err := p.doCMRequest("GET", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[GetFastSpecificationOpenAPI]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// create GET request to get Fast templates
func (p *BigipNextCM) GetFastTemplates() error {
	fastURL := fmt.Sprintf("%s%s%s", p.Host, uriFast, "/templates")
	f5osLogger.Info("[GetFastTemplates]", "URI Path", fastURL)
	respData, err := p.doCMRequest("GET", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[GetFastTemplates]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// create GET request to get Fast version
func (p *BigipNextCM) GetFastVersion() error {
	fastURL := fmt.Sprintf("%s%s%s", p.Host, uriFast, "/v1/version")
	f5osLogger.Info("[GetFastVersion]", "URI Path", fastURL)
	respData, err := p.doCMRequest("GET", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[GetFastVersion]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// create GET request to get Fast applications
func (p *BigipNextCM) GetFastApplications() error {
	fastURL := fmt.Sprintf("%s%s%s", p.Host, uriFast, "/applications")
	f5osLogger.Info("[GetFastApplications]", "URI Path", fastURL)
	respData, err := p.doCMRequest("GET", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[GetFastApplications]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// https://clouddocs.f5.com/api/v1/spaces/default/application-templates

func (p *BigipNextCM) GetApplicationTemplates() error {
	fastURL := fmt.Sprintf("%s%s%s", p.Host, uriDefault, "/application-templates")
	f5osLogger.Info("[GetApplicationTemplates]", "URI Path", fastURL)
	respData, err := p.doCMRequest("GET", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[GetApplicationTemplates]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// {"name":"teststandard","template_name":"http","set_name":"Examples","parameters":{"application_name":"teststandard","application_description":""},"allowOverwrite":true}
// create struct for above template

type FastRequestDraft struct {
	Name       string `json:"name,omitempty"`
	Parameters struct {
		ApplicationDescription string     `json:"application_description,omitempty"`
		ApplicationName        string     `json:"application_name,omitempty"`
		Pools                  []FastPool `json:"pools,omitempty"`
		Virtuals               []Virtual  `json:"virtuals,omitempty"`
	} `json:"parameters,omitempty"`
	SetName        string `json:"set_name,omitempty"`
	TemplateName   string `json:"template_name,omitempty"`
	TenantName     string `json:"tenant_name,omitempty"`
	AllowOverwrite bool   `json:"allowOverwrite,omitempty"`
}

type ApplicationTemplate struct {
	Name           string    `json:"name,omitempty"`
	TemplateName   string    `json:"template_name,omitempty"`
	SetName        string    `json:"set_name,omitempty"`
	Parameters     Parameter `json:"parameters,omitempty"`
	AllowOverwrite bool      `json:"allowOverwrite,omitempty"`
}

type Parameter struct {
	ApplicationName        string     `json:"application_name,omitempty"`
	ApplicationDescription string     `json:"application_description,omitempty"`
	Pools                  []FastPool `json:"pools,omitempty"`
	Virtuals               []Virtual  `json:"virtuals,omitempty"`
}

// create struct for above fast request
type FastPool struct {
	LoadBalancingMode string   `json:"loadBalancingMode,omitempty"`
	MonitorType       []string `json:"monitorType,omitempty"`
	PoolName          string   `json:"poolName,omitempty"`
	ServicePort       int      `json:"servicePort,omitempty"`
}

type Virtual struct {
	FastL4TOS                     int      `json:"FastL4_TOS,omitempty"`
	FastL4IdleTimeout             int      `json:"FastL4_idleTimeout,omitempty"`
	FastL4LooseClose              *bool    `json:"FastL4_looseClose,omitempty"`
	FastL4LooseInitialization     *bool    `json:"FastL4_looseInitialization,omitempty"`
	FastL4PvaAcceleration         string   `json:"FastL4_pvaAcceleration,omitempty"`
	FastL4PvaDynamicClientPackets int      `json:"FastL4_pvaDynamicClientPackets,omitempty"`
	FastL4PvaDynamicServerPackets int      `json:"FastL4_pvaDynamicServerPackets,omitempty"`
	FastL4ResetOnTimeout          *bool    `json:"FastL4_resetOnTimeout,omitempty"`
	FastL4TcpCloseTimeout         int      `json:"FastL4_tcpCloseTimeout,omitempty"`
	FastL4TcpHandshakeTimeout     int      `json:"FastL4_tcpHandshakeTimeout,omitempty"`
	InspectionServicesEnum        []string `json:"InspectionServicesEnum,omitempty"`
	TCPIdleTimeout                int      `json:"TCP_idle_timeout,omitempty"`
	UDPIdleTimeout                int      `json:"UDP_idle_timeout,omitempty"`
	WAFPolicyName                 string   `json:"WAFPolicyName,omitempty"`
	AutoLastHop                   string   `json:"auto_last_hop,omitempty"`
	EnableAccess                  bool     `json:"enable_Access,omitempty"`
	EnableFastL4                  bool     `json:"enable_FastL4,omitempty"`
	EnableFastL4DSR               bool     `json:"enable_FastL4_DSR,omitempty"`
	EnableHTTP2Profile            bool     `json:"enable_HTTP2_Profile,omitempty"`
	EnableHTTPProfile             bool     `json:"enable_HTTP_Profile,omitempty"`
	EnableInspectionServices      bool     `json:"enable_InspectionServices,omitempty"`
	EnableSsloPolicy              bool     `json:"enable_SsloPolicy,omitempty"`
	EnableTCPProfile              bool     `json:"enable_TCP_Profile,omitempty"`
	EnableTLSClient               bool     `json:"enable_TLS_Client,omitempty"`
	EnableTLSServer               bool     `json:"enable_TLS_Server,omitempty"`
	EnableUDPProfile              bool     `json:"enable_UDP_Profile,omitempty"`
	EnableWAF                     bool     `json:"enable_WAF,omitempty"`
	EnableIRules                  bool     `json:"enable_iRules,omitempty"`
	EnableMirroring               bool     `json:"enable_mirroring,omitempty"`
	EnableSnat                    bool     `json:"enable_snat,omitempty"`
	IRulesEnum                    []string `json:"iRulesEnum,omitempty"`
	IsVirtualTypeStandard         bool     `json:"is_virtual_type_standard,omitempty"`
	MultiCertificatesEnum         []string `json:"multiCertificatesEnum,omitempty"`
	PerRequestAccessPolicyEnum    string   `json:"perRequestAccessPolicyEnum,omitempty"`
	Pool                          string   `json:"pool,omitempty"`
	ServerAddressTranslation      bool     `json:"server_address_translation,omitempty"`
	SnatAddresses                 []string `json:"snat_addresses,omitempty"`
	SnatAutomap                   bool     `json:"snat_automap,omitempty"`
	TLSServerCertificates         []string `json:"tls_server_certificates,omitempty"`
	VirtualName                   string   `json:"virtualName,omitempty"`
	VirtualPort                   int      `json:"virtualPort,omitempty"`
	VirtualType                   string   `json:"virtualType,omitempty"`
}

// type VirtualServer struct {
// 	FastL4IdleTimeout         int      `json:"FastL4_idleTimeout,omitempty"`
// 	FastL4LooseClose          bool     `json:"FastL4_looseClose,omitempty"`
// 	FastL4LooseInitialization bool     `json:"FastL4_looseInitialization,omitempty"`
// 	FastL4ResetOnTimeout      bool     `json:"FastL4_resetOnTimeout,omitempty"`
// 	FastL4TcpCloseTimeout     int      `json:"FastL4_tcpCloseTimeout,omitempty"`
// 	FastL4TcpHandshakeTimeout int      `json:"FastL4_tcpHandshakeTimeout,omitempty"`
// 	TCPIdleTimeout            int      `json:"TCP_idle_timeout,omitempty"`
// 	UDPIdleTimeout            int      `json:"UDP_idle_timeout,omitempty"`
// 	Ciphers                   string   `json:"ciphers,omitempty"`
// 	CiphersServer             string   `json:"ciphers_server,omitempty"`
// 	EnableAccess              bool     `json:"enable_Access,omitempty"`
// 	EnableFastL4              bool     `json:"enable_FastL4,omitempty"`
// 	EnableHTTP2Profile        bool     `json:"enable_HTTP2_Profile,omitempty"`
// 	EnableTCPProfile          bool     `json:"enable_TCP_Profile,omitempty"`
// 	EnableTLSClient           bool     `json:"enable_TLS_Client,omitempty"`
// 	EnableTLSServer           bool     `json:"enable_TLS_Server,omitempty"`
// 	EnableUDPProfile          bool     `json:"enable_UDP_Profile,omitempty"`
// 	EnableWAF                 bool     `json:"enable_WAF,omitempty"`
// 	WAFPolicyName             string   `json:"WAF_Policy_Name,omitempty"`
// 	EnableIRules              bool     `json:"enable_iRules,omitempty"`
// 	EnableMirroring           bool     `json:"enable_mirroring,omitempty"`
// 	EnableSnat                bool     `json:"enable_snat,omitempty"`
// 	IRulesEnum                []string `json:"iRulesEnum,omitempty"`
// 	Pool                      string   `json:"pool,omitempty"`
// 	SnatAddresses             []string `json:"snat_addresses,omitempty"`
// 	SnatAutomap               bool     `json:"snat_automap,omitempty"`
// 	TlsC12                    bool     `json:"tls_c_1_2,omitempty"`
// 	TlsC13                    bool     `json:"tls_c_1_3,omitempty"`
// 	TlsS12                    bool     `json:"tls_s_1_2,omitempty"`
// 	TlsS13                    bool     `json:"tls_s_1_3,omitempty"`
// 	VirtualName               string   `json:"virtualName,omitempty"`
// 	VirtualPort               int      `json:"virtualPort,omitempty"`
// }

// create FAST app draft request using above json payload
func (p *BigipNextCM) PostFastApplicationDraft(config *FastRequestDraft) (string, error) {
	// fastURL := fmt.Sprintf("%s%s%s", p.Host, uriFast, "/appsvcs")
	fastURL := fmt.Sprintf("%s%s%s", p.Host, uriDefault, "/appsvcs/blueprints")
	f5osLogger.Info("[PostFastApplicationDraft]", "URI Path", fastURL)
	f5osLogger.Info("[PostFastApplicationDraft]", "Config", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	respData, err := p.doCMRequest("POST", fastURL, body)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostFastApplicationDraft]", "Data::", hclog.Fmt("%+v", string(respData)))
	// {"_links":{"self":{"href":"/api/v1/spaces/default/appsvcs/blueprints/d59b1bf8-9e4d-47ea-bafa-6986479fee0e"}},"id":"d59b1bf8-9e4d-47ea-bafa-6986479fee0e","message":"application created successfully","status":200}

	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostFastApplicationDraft]", "Draft ID", hclog.Fmt("%+v", respString["id"].(string)))
	return respString["id"].(string), nil
}

// /api/v1/spaces/default/
func (p *BigipNextCM) PostApplicationTemplate(template *ApplicationTemplate) ([]byte, error) {
	fastURL := fmt.Sprintf("%s%s%s", p.Host, uriDefault, "/appsvcs/blueprints")
	f5osLogger.Info("[PostApplicationTemplate]", "URI Path", fastURL)
	jsonData, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[PostApplicationTemplate]", "Data::", hclog.Fmt("%+v", string(jsonData)))
	respData, err := p.doCMRequest("POST", fastURL, jsonData)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[PostApplicationTemplate]", "Data::", hclog.Fmt("%+v", string(respData)))
	// {"_links":{"self":{"href":"/api/v1/spaces/default/appsvcs/blueprints/fbb7a8c9-4c58-4f06-bf32-ab262a6ae2d5"}},"id":"fbb7a8c9-4c58-4f06-bf32-ab262a6ae2d5","message":"application created successfully","status":200}

	return respData, nil
}

func (p *BigipNextCM) PatchApplicationTemplate(appID, config interface{}) ([]byte, error) {
	fastURL := fmt.Sprintf("%s%s%s/%s", p.Host, uriDefault, "/appsvcs/blueprints", appID)
	f5osLogger.Info("[GetApplicationTemplate]", "URI Path", fastURL)
	jsonData, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.doCMRequest("PATCH", fastURL, jsonData)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetApplicationTemplate]", "Data::", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

// /api/v1/spaces/default/appsvcs/blueprints
func (p *BigipNextCM) GetApplicationBlueprints(blueprintID string) ([]byte, error) {
	fastURL := fmt.Sprintf("%s%s%s/%s", p.Host, uriDefault, "/appsvcs/blueprints", blueprintID)
	f5osLogger.Info("[GetApplicationBlueprints]", "URI Path", fastURL)
	respData, err := p.doCMRequest("GET", fastURL, nil)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetApplicationBlueprints]", "Data::", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

// /api/v1/spaces/default/appsvcs/blueprints/7fe953de-8390-4a5f-ae28-385d99f284c9
func (p *BigipNextCM) DeleteApplicationBlueprint(blueprintID string) error {
	fastURL := fmt.Sprintf("%s%s%s/%s", p.Host, uriDefault, "/appsvcs/blueprints", blueprintID)
	f5osLogger.Info("[DeleteApplicationBlueprint]", "URI Path", fastURL)
	respData, err := p.doCMRequest("DELETE", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteApplicationBlueprint]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// Helper function to create a pointer to a boolean
func BoolPtr(b bool) *bool {
	return &b
}
