/*
Copyright 2023 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
// Package bigipnext interacts with BIGIP-NEXT/CM systems using the OPEN API.
package bigipnext

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
)

const (
	uriCMRoot            = "/api"
	uriCMLogin           = "/api/login"
	uriCMProxyFileUpload = "/device/v1/proxy-file-upload"
	uriCMFileUpload      = "/system/v1/files"
	uriFast              = "/mgmt/shared/fast"
	uriOpenAPI           = "/openapi"
	uriInventory         = "/device/v1/inventory"
	uriBackups           = "/device/v1/backups"
	uriBackupTasks       = "/device/v1/backup-tasks"
	uriRestoreTasks      = "/device/v1/restore-tasks"
	uriCertificate       = "/api/v1/spaces/default/certificates"
	uriCMUpgradeTask     = "/upgrade-manager/v1/upgrade-tasks"
	uriAS3Root           = "/api/v1/spaces/default/appsvcs"
	uriDiscoverInstance  = "/v1/spaces/default/instances"
	// uriCertificateUpdate = "/api/certificate/v1/certificates"
	uriGlobalResiliency    = "/api/v1/spaces/default/gslb/gr-groups"
	uriGetGlobalResiliency = "/v1/spaces/default/gslb/gr-groups"
	uriGetAlert            = "/alert/v1/alerts/?limit=1&sort=-start_time&filter=source%20eq%20%27DNS%27%20and%20status%20eq%20%27ACTIVE%27&select=summary"
	uriWafReport           = "/v1/spaces/default/security/waf/reports"
	// uriGetWafReport        = "/v1/spaces/default/security/waf/reports"
	uriWafPolicy    = "/api/v1/spaces/default/security/waf-policies"
	uriGetWafPolicy = "/v1/spaces/default/security/waf-policies"
)

// BIG IP Next CM Config Request structure
type BigipNextCMReqConfig struct {
	Host      string
	User      string
	Password  string
	Port      int
	Transport *http.Transport
	// UserAgent is an optional field that specifies the caller of this request.
	UserAgent     string
	Teem          bool
	ConfigOptions *ConfigOptions
}

type BigipNextCMLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BigipNextCMLoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserId       string `json:"user_id"`
}

// BigipCMNext is a container for our session state.
type BigipNextCM struct {
	Host          string
	Token         string // The authentication token/token has a short lifetime and should be included in subsequent API requests
	RefreshToken  string // The refresh token has a longer lifetime and can be used to get a new authentication token
	Transport     *http.Transport
	UserAgent     string
	Teem          bool
	ConfigOptions *ConfigOptions
	PlatformType  string
}

type CMError struct {
	Error   error
	Message string
	Code    int
}

type DeviceInventoryList struct {
	Embedded struct {
		Devices []struct {
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
			Address                    string    `json:"address"`
			CertificateValidated       time.Time `json:"certificate_validated"`
			CertificateValidationError string    `json:"certificate_validation_error"`
			CertificateValidity        bool      `json:"certificate_validity"`
			Hostname                   string    `json:"hostname"`
			Id                         string    `json:"id"`
			Mode                       string    `json:"mode"`
			PlatformType               string    `json:"platform_type"`
			Port                       int       `json:"port"`
			Version                    string    `json:"version"`
		} `json:"devices"`
	} `json:"_embedded"`
	Count int `json:"count"`
	Total int `json:"total"`
}

// CmNewSession sets up connection to the BIG-IP Next CM system.

func CmNewSession(bigipNextCmObj *BigipNextCMReqConfig) (*BigipNextCM, error) {
	f5osLogger.Info("[NewSession] Session creation Starts...")
	var urlString string
	bigipNextCmSession := &BigipNextCM{}
	if !strings.HasPrefix(bigipNextCmObj.Host, "http") {
		urlString = fmt.Sprintf("https://%s", bigipNextCmObj.Host)
	} else {
		urlString = bigipNextCmObj.Host
	}
	u, _ := url.Parse(urlString)
	_, port, _ := net.SplitHostPort(u.Host)

	if bigipNextCmObj.Port != 0 && port == "" {
		urlString = fmt.Sprintf("%s:%d", urlString, bigipNextCmObj.Port)
	}
	if bigipNextCmObj.ConfigOptions == nil {
		bigipNextCmObj.ConfigOptions = defaultConfigOptions
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	bigipNextCmSession.Host = urlString
	bigipNextCmSession.Transport = tr
	bigipNextCmSession.ConfigOptions = bigipNextCmObj.ConfigOptions
	client := &http.Client{
		Transport: tr,
	}
	method := "POST"
	urlString = fmt.Sprintf("%s%s", urlString, uriCMLogin)
	reqBody := &BigipNextCMLoginReq{}
	reqBody.Username = bigipNextCmObj.User
	reqBody.Password = bigipNextCmObj.Password
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[NewSession]", "URL", hclog.Fmt("%+v", urlString))

	req, err := http.NewRequest(method, urlString, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentTypeHeader)
	// req.SetBasicAuth(bigipNextCmObj.User, bigipNextCmObj.Password)
	// TODO Retry Logic
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bodyResp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error message: %+v (Body:%+v)", res.Status, string(bodyResp))
	}
	var resp BigipNextCMLoginResp
	err = json.Unmarshal(bodyResp, &resp)
	if err != nil {
		return nil, err
	}
	bigipNextCmSession.Token = resp.AccessToken
	bigipNextCmSession.RefreshToken = resp.RefreshToken
	f5osLogger.Info("[NewSession] Session creation Success")
	return bigipNextCmSession, nil
}

func (p *BigipNextCM) GetDeviceInventory() (*DeviceInventoryList, error) {
	f5osLogger.Debug("[GetDeviceInventory]", "URI Path", "/device/v1/inventory")
	listBigipNext := &DeviceInventoryList{}
	respData, err := p.GetCMRequest("/device/v1/inventory")
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[GetDeviceInventory]", "List of BIG-IP Next:", hclog.Fmt("%+v", string(respData)))
	json.Unmarshal(respData, listBigipNext)
	return listBigipNext, nil
}

func (p *BigipNextCM) PostCMRequest(path string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, path)
	f5osLogger.Info("[PostCMRequest]", "Request path", hclog.Fmt("%+v", url))
	f5osLogger.Info("[PostCMRequest]", "Request body", hclog.Fmt("%+v", string(body)))
	return p.doCMRequest("POST", url, body)
}

func (p *BigipNextCM) PutCMRequest(path string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, path)
	f5osLogger.Info("[PutCMRequest]", "Request path", hclog.Fmt("%+v", url))
	f5osLogger.Info("[PutCMRequest]", "Request body", hclog.Fmt("%+v", string(body)))
	return p.doCMRequest("PUT", url, body)
}

func (p *BigipNextCM) GetCMRequest(path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, path)
	f5osLogger.Info("[GetCMRequest]", "Request path", hclog.Fmt("%+v", url))
	return p.doCMRequest("GET", url, nil)
}

func (p *BigipNextCM) DeleteCMRequest(path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, path)
	f5osLogger.Info("[DeleteCMRequest]", "Request path", hclog.Fmt("%+v", url))
	return p.doCMRequest("DELETE", url, nil)
}

func (p *BigipNextCM) doCMRequest(op, path string, body []byte) ([]byte, error) {
	f5osLogger.Info("[doCMRequest]", "Request path", hclog.Fmt("%+v", path))
	if len(body) > 0 {
		f5osLogger.Debug("[doCMRequest]", "Request body", hclog.Fmt("%+v", string(body)))
	}
	req, err := http.NewRequest(op, path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Token))
	req.Header.Set("Content-Type", contentTypeHeader)
	client := &http.Client{
		Transport: p.Transport,
		Timeout:   p.ConfigOptions.APICallTimeout,
	}
	// f5osLogger.Debug("[doCMRequest]", "Request path", hclog.Fmt("%+v", req.))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	f5osLogger.Debug("[doCMRequest]", "Resp CODE", hclog.Fmt("%+v", resp.StatusCode))
	if resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 202 || resp.StatusCode == 204 {
		return io.ReadAll(resp.Body)
	}
	if resp.StatusCode == 401 {
		//{"code":401,"error":{"status":401,"message":"GATEWAY-00023: The access token expired."}
		byteData, _ := io.ReadAll(resp.Body)
		// check if message is "GATEWAY-00023: The access token expired."
		if strings.Contains(string(byteData), "GATEWAY-00023: The access token expired.") {
			//
			f5osLogger.Info("[doCMRequest]", "Refresh Token", hclog.Fmt("%+v", p.CMTokenRefresh))
			err = p.CMTokenRefreshNew()
			if err != nil {
				return nil, err
			}
			// retry request
			return p.doCMRequest(op, path, body)
		}
	}
	if resp.StatusCode >= 400 {
		byteData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(`{"code":%d,"error":%s`, resp.StatusCode, byteData)
	}
	return nil, nil
}

func (p *BigipNextCM) GetAs3SpecificationOpenAPI() error {
	as3URL := fmt.Sprintf("%s%s/%s", p.Host, uriAs3, "openapi")
	f5osLogger.Info("[GetAs3SpecificationOpenAPI]", "URI Path", as3URL)
	respData, err := p.doCMRequest("GET", as3URL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[GetAs3SpecificationOpenAPI]", "Data::", hclog.Fmt("%+v", string(respData)))
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

// // create struct for above fast request
// type FastRequest struct {
// 	Name       string `json:"name,omitempty"`
// 	Parameters struct {
// 		ApplicationName string   `json:"application_name,omitempty"`
// 		ServerAddresses []string `json:"servers,omitempty"`
// 		ServicePort     int      `json:"servicePort,omitempty"`
// 		VirtualAddress  string   `json:"virtualAddress,omitempty"`
// 		VirtualPort     int      `json:"virtualPort,omitempty"`
// 	} `json:"parameters,omitempty"`
// 	Target struct {
// 		Address  string `json:"address,omitempty"`
// 		Hostname string `json:"hostname,omitempty"`
// 	} `json:"target,omitempty"`
// }

type Server struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}

type FastRequest struct {
	AllowOverwrite bool   `json:"allowOverwrite,omitempty"`
	Name           string `json:"name,omitempty"`
	Parameters     struct {
		EnableAccess                   bool          `json:"enable_Access,omitempty"`
		EnableHTTP2Profile             bool          `json:"enable_HTTP2_Profile,omitempty"`
		EnableTLSServer                bool          `json:"enable_TLS_Server,omitempty"`
		EnableWAF                      bool          `json:"enable_WAF,omitempty"`
		EnableIRules                   bool          `json:"enable_iRules,omitempty"`
		EnableSnat                     bool          `json:"enable_snat,omitempty"`
		LoadBalancingMode              string        `json:"loadBalancingMode,omitempty"`
		MonitorType                    string        `json:"monitorType,omitempty"`
		ServicePort                    int           `json:"servicePort,omitempty"`
		VirtualPort                    int           `json:"virtualPort,omitempty"`
		EnableFastL4                   bool          `json:"enable_FastL4,omitempty"`
		EnableTCPProfile               bool          `json:"enable_TCP_Profile,omitempty"`
		EnableUDPProfile               bool          `json:"enable_UDP_Profile,omitempty"`
		AccessAdditionalConfigurations string        `json:"accessAdditionalConfigurations,omitempty"`
		SnatAutomap                    bool          `json:"snat_automap,omitempty"`
		SnatAddresses                  []interface{} `json:"snat_addresses,omitempty"`
		IRulesContent                  []interface{} `json:"iRulesContent,omitempty"`
		IRulesList                     []interface{} `json:"iRulesList,omitempty"`
		UDPIdleTimeout                 int           `json:"UDP_idle_timeout,omitempty"`
		TCPIdleTimeout                 int           `json:"TCP_idle_timeout,omitempty"`
		FastL4IdleTimeout              int           `json:"FastL4_idleTimeout,omitempty"`
		FastL4LooseClose               bool          `json:"FastL4_looseClose,omitempty"`
		FastL4LooseInitialization      bool          `json:"FastL4_looseInitialization,omitempty"`
		FastL4ResetOnTimeout           bool          `json:"FastL4_resetOnTimeout,omitempty"`
		FastL4TcpCloseTimeout          int           `json:"FastL4_tcpCloseTimeout,omitempty"`
		FastL4TcpHandshakeTimeout      int           `json:"FastL4_tcpHandshakeTimeout,omitempty"`
		ApplicationName                string        `json:"application_name,omitempty"`
		VirtualAddress                 string        `json:"virtualAddress,omitempty"`
		Servers                        []Server      `json:"servers,omitempty"`
	} `json:"parameters,omitempty"`
	Target struct {
		Address  string `json:"address,omitempty"`
		Hostname string `json:"hostname,omitempty"`
	} `json:"target,omitempty"`
}

// Create POST request to render Fast application
func (p *BigipNextCM) PostFastRenderApplication(config *FastRequest) error {
	fastURL := fmt.Sprintf("%s%s%s", p.Host, uriFast, "/render")
	f5osLogger.Info("[PostFastRender]", "URI Path", fastURL)
	f5osLogger.Info("[PostFastRender]", "Config", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return err
	}
	respData, err := p.doCMRequest("POST", fastURL, body)
	if err != nil {
		return err
	}
	f5osLogger.Info("[PostFastRender]", "Data::", hclog.Fmt("%+v", string(respData)))
	// {"_links":{"self":{"href":"/applications/tasks/e42bdc83-1da4-4dfd-902b-0c27dd8a8f53"}},"path":"/applications/tasks/e42bdc83-1da4-4dfd-902b-0c27dd8a8f53"}
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return err
	}
	f5osLogger.Info("[PostFastRender]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	for i := 0; i < 10; i++ {
		err = p.GetFastApplicationTaskStatus(respString["path"].(string))
		if err != nil {
			return err
		}
	}
	// err = p.GetFastApplicationTaskStatus(respString["path"].(string))
	// if err != nil {
	// 	return err
	// }
	return nil
}

// Create POST request to create Fast application
func (p *BigipNextCM) PostFastApplication(config *FastRequest) error {
	fastURL := fmt.Sprintf("%s%s%s", p.Host, uriFast, "/appsvcs")
	// fastURL := fmt.Sprintf("%s%s%s", p.Host, uriFast, "/applications")
	f5osLogger.Info("[PostFastApplication]", "URI Path", fastURL)
	f5osLogger.Info("[PostFastApplication]", "Config", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return err
	}
	respData, err := p.doCMRequest("POST", fastURL, body)
	if err != nil {
		return err
	}
	f5osLogger.Info("[PostFastApplication]", "Data::", hclog.Fmt("%+v", string(respData)))
	// {"_links":{"self":{"href":"/applications/tasks/e42bdc83-1da4-4dfd-902b-0c27dd8a8f53"}},"path":"/applications/tasks/e42bdc83-1da4-4dfd-902b-0c27dd8a8f53"}
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return err
	}
	f5osLogger.Info("[PostFastApplication]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	for i := 0; i < 10; i++ {
		err = p.GetFastApplicationTaskStatus(respString["path"].(string))
		if err != nil {
			return err
		}
	}
	// err = p.GetFastApplicationTaskStatus(respString["path"].(string))
	// if err != nil {
	// 	return err
	// }
	return nil
}

// https://<BIG-IP-Next-Central-Manager-IP-Address>/mgmt/shared/fast/applications/tasks/16114eee-8622-4227-a3f7-a2758214677f
// create GET request to get Fast application task status
func (p *BigipNextCM) GetFastApplicationTaskStatus(taskid string) error {
	fastURL := fmt.Sprintf("%s%s%s", p.Host, uriFast, taskid)
	f5osLogger.Info("[GetFastApplicationTaskStatus]", "URI Path", fastURL)
	respData, err := p.doCMRequest("GET", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[GetFastApplicationTaskStatus]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// https://<BIG-IP-Next-Central-Manager-IP-Address>/mgmt/shared/fast/applications/{tenantName}/{appName}
// create DELETE request to delete Fast application
func (p *BigipNextCM) DeleteFastApplication(tenantName, appName string) error {
	fastURL := fmt.Sprintf("%s%s%s/%s/%s", p.Host, uriFast, "/applications", tenantName, appName)
	f5osLogger.Info("[DeleteFastApplication]", "URI Path", fastURL)
	respData, err := p.doCMRequest("DELETE", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteFastApplication]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// https://<BIG-IP-Next-Central-Manager-IP-Address>/mgmt/shared/fast/applications/{tenantName}/{appName}
// create delete request to delete Fast application with filter query parameters
func (p *BigipNextCM) DeleteFastApplicationWithFilter(tenantName, appName, filter string) error {
	// filter = url.QueryEscape(filter)
	fastURL := fmt.Sprintf("%s%s%s/%s/%s?filter=%s", p.Host, uriFast, "/applications", tenantName, appName, filter)
	f5osLogger.Info("[DeleteFastApplicationWithFilter]", "URI Path", fastURL)
	respData, err := p.doCMRequest("DELETE", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteFastApplicationWithFilter]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// create struct for above fast request
type FastPool struct {
	LoadBalancingMode string   `json:"loadBalancingMode,omitempty"`
	MonitorType       []string `json:"monitorType,omitempty"`
	PoolName          string   `json:"poolName,omitempty"`
	ServicePort       int      `json:"servicePort,omitempty"`
}
type VirtualServer struct {
	FastL4IdleTimeout         int      `json:"FastL4_idleTimeout,omitempty"`
	FastL4LooseClose          bool     `json:"FastL4_looseClose,omitempty"`
	FastL4LooseInitialization bool     `json:"FastL4_looseInitialization,omitempty"`
	FastL4ResetOnTimeout      bool     `json:"FastL4_resetOnTimeout,omitempty"`
	FastL4TcpCloseTimeout     int      `json:"FastL4_tcpCloseTimeout,omitempty"`
	FastL4TcpHandshakeTimeout int      `json:"FastL4_tcpHandshakeTimeout,omitempty"`
	TCPIdleTimeout            int      `json:"TCP_idle_timeout,omitempty"`
	UDPIdleTimeout            int      `json:"UDP_idle_timeout,omitempty"`
	Ciphers                   string   `json:"ciphers,omitempty"`
	CiphersServer             string   `json:"ciphers_server,omitempty"`
	EnableAccess              bool     `json:"enable_Access,omitempty"`
	EnableFastL4              bool     `json:"enable_FastL4,omitempty"`
	EnableHTTP2Profile        bool     `json:"enable_HTTP2_Profile,omitempty"`
	EnableTCPProfile          bool     `json:"enable_TCP_Profile,omitempty"`
	EnableTLSClient           bool     `json:"enable_TLS_Client,omitempty"`
	EnableTLSServer           bool     `json:"enable_TLS_Server,omitempty"`
	EnableUDPProfile          bool     `json:"enable_UDP_Profile,omitempty"`
	EnableWAF                 bool     `json:"enable_WAF,omitempty"`
	EnableIRules              bool     `json:"enable_iRules,omitempty"`
	EnableMirroring           bool     `json:"enable_mirroring,omitempty"`
	EnableSnat                bool     `json:"enable_snat,omitempty"`
	IRulesEnum                []string `json:"iRulesEnum,omitempty"`
	Pool                      string   `json:"pool,omitempty"`
	SnatAddresses             []string `json:"snat_addresses,omitempty"`
	SnatAutomap               bool     `json:"snat_automap,omitempty"`
	TlsC12                    bool     `json:"tls_c_1_2,omitempty"`
	TlsC13                    bool     `json:"tls_c_1_3,omitempty"`
	TlsS12                    bool     `json:"tls_s_1_2,omitempty"`
	TlsS13                    bool     `json:"tls_s_1_3,omitempty"`
	VirtualName               string   `json:"virtualName,omitempty"`
	VirtualPort               int      `json:"virtualPort,omitempty"`
}

type FastRequestDraft struct {
	Name       string `json:"name,omitempty"`
	Parameters struct {
		ApplicationDescription string          `json:"application_description,omitempty"`
		ApplicationName        string          `json:"application_name,omitempty"`
		Pools                  []FastPool      `json:"pools,omitempty"`
		Virtuals               []VirtualServer `json:"virtuals,omitempty"`
	} `json:"parameters,omitempty"`
	SetName        string `json:"set_name,omitempty"`
	TemplateName   string `json:"template_name,omitempty"`
	TenantName     string `json:"tenant_name,omitempty"`
	AllowOverwrite bool   `json:"allowOverwrite,omitempty"`
}

// create FAST app draft request using above json payload
func (p *BigipNextCM) PostFastApplicationDraft(config *FastRequestDraft) (string, error) {
	fastURL := fmt.Sprintf("%s%s%s", p.Host, uriFast, "/appsvcs")
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
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostFastApplicationDraft]", "Task Path", hclog.Fmt("%+v", respString["id"].(string)))
	return respString["id"].(string), nil
}

type CertificateRequestDraft struct {
	Issuer                 string   `json:"issuer,omitempty"`
	Name                   string   `json:"name,omitempty"`
	CommonName             string   `json:"common_name,omitempty"`
	Division               []string `json:"division,omitempty"`
	Organization           []string `json:"organization,omitempty"`
	Locality               []string `json:"locality,omitempty"`
	State                  []string `json:"state,omitempty"`
	Country                []string `json:"country,omitempty"`
	Email                  []string `json:"email,omitempty"`
	SubjectAlternativeName string   `json:"subject_alternative_name,omitempty"`
	DurationInDays         int      `json:"duration_in_days,omitempty"`
	KeyType                string   `json:"key_type,omitempty"`
	KeySecurityType        string   `json:"key_security_type,omitempty"`
	KeySize                int      `json:"key_size,omitempty"`
	KeyCurveName           string   `json:"key_curve_name,omitempty"`
	KeyPassphrase          string   `json:"key_passphrase,omitempty"`
	AdministratorEmail     string   `json:"administrator_email,omitempty"`
	// ChallengePassword      string   `json:"challenge_password,omitempty"`
}

// create Certificate draft request using above json payload
func (p *BigipNextCM) PostCertificateCreate(config interface{}, op string) (string, error) {
	createCertificateURL := fmt.Sprintf("%s%s%s", p.Host, uriCertificate, "/create")
	if op == "UPDATE" {
		createCertificateURL = fmt.Sprintf("%s%s%s", p.Host, uriCertificate, "/renew")
	}
	if op == "IMPORT" {
		createCertificateURL = fmt.Sprintf("%s%s%s", p.Host, uriCertificate, "/import")
	}
	if op == "UPDATEIMPORT" {
		createCertificateURL = fmt.Sprintf("%s%s%s", p.Host, uriCertificate, "/import")
	}
	f5osLogger.Info("[PostCertificateCreate]", "URI Path", createCertificateURL)
	f5osLogger.Debug("[PostCertificateCreate]", "Config", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	respData, err := p.doCMRequest("POST", createCertificateURL, body)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostCertificateCreate]", "Data::", hclog.Fmt("%+v", string(respData)))
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostCertificateCreate]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	pathList := strings.Split(respString["path"].(string), "/")
	id := pathList[len(pathList)-1]

	getCertificateURL := fmt.Sprintf("%s%s/%s", p.Host, uriCertificate, id)
	time.Sleep(5 * time.Second)
	certData, err := p.doCMRequest("GET", getCertificateURL, nil)
	if err != nil {
		return "", err
	}
	certString := make(map[string]interface{})
	err = json.Unmarshal(certData, &certString)
	if err != nil {
		return "", err
	}
	status := certString["status"].(string)
	f5osLogger.Info("[PostCertificateCreate]", "Certificate Status::", hclog.Fmt("%+v", string(status)))
	if status == "failed" {
		return "", fmt.Errorf("certificate failure reason is :%+v ", certString["failure_reason"].(string))
	}
	return id, nil
}

func (p *BigipNextCM) GetNextCMCertificate(id string) (interface{}, error) {
	getCertificateURL := fmt.Sprintf("%s%s/%s", p.Host, uriCertificate, id)
	f5osLogger.Info("[GetNextCMCertificate]", "URI Path", getCertificateURL)
	respData, err := p.doCMRequest("GET", getCertificateURL, nil)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetNextCMCertificate]", "Before Resp ", hclog.Fmt("%+v", string(respData)))
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetNextCMCertificate]", "Resp Message", hclog.Fmt("%+v", respString))
	return respString, nil
}

func (p *BigipNextCM) GetNextCMImportCertificateKeyData(id string) (interface{}, error) {
	getCertificateURL := fmt.Sprintf("%s%s/%s/%s", p.Host, uriCertificate, id, "crt")
	f5osLogger.Info("[GetNextCMImportCertificateKeyData]", "URI Path for cert", getCertificateURL)
	certData, err := p.doCMRequest("GET", getCertificateURL, nil)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetNextCMImportCertificateKeyData]", "Before Resp ", hclog.Fmt("%+v", string(certData)))
	// certString := make(map[string]interface{})
	// err = json.Unmarshal(certData, &certString)
	// if err != nil {
	// 	return nil, err
	// }

	getKeyURL := fmt.Sprintf("%s%s/%s/%s", p.Host, uriCertificate, id, "key")
	f5osLogger.Info("[GetNextCMImportCertificateKeyData]", "URI Path for key", getKeyURL)
	keyData, err := p.doCMRequest("GET", getKeyURL, nil)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetNextCMImportCertificateKeyData]", "Before Resp ", hclog.Fmt("%+v", string(keyData)))
	// keyString := make(map[string]interface{})
	// err = json.Unmarshal(keyData, &keyString)
	// if err != nil {
	// 	return nil, err
	// }
	keycertData := make(map[string]interface{})
	keycertData["key_data"] = string(keyData)
	keycertData["cert_data"] = string(certData)

	f5osLogger.Info("[GetNextCMCertificate]", "Resp Message", hclog.Fmt("%+v", keycertData))
	return keycertData, nil
}

func (p *BigipNextCM) DeleteNextCMCertificate(id string) error {
	deleteCertificateURL := fmt.Sprintf("%s%s/%s", p.Host, uriCertificate, id)
	f5osLogger.Info("[DeleteNextCMCertificate]", "URI Path", deleteCertificateURL)
	respData, err := p.doCMRequest("DELETE", deleteCertificateURL, nil)
	if err != nil {
		return err
	}
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteNextCMCertificate]", "Task Message", hclog.Fmt("%+v", respString["message"].(string)))

	if respString["message"] != "deleted" {
		return err
	}
	return nil
}

type ImportCertificateRequestDraft struct {
	Name           string `json:"name,omitempty"`
	KeyPassphrase  string `json:"key_passphrase,omitempty"`
	CertPassphrase string `json:"cert_passphrase,omitempty"`
	KeyText        string `json:"key_text,omitempty"`
	CertText       string `json:"cert_text,omitempty"`
	ImportType     string `json:"import_type,omitempty"`
	Id             string `json:"id,omitempty"`
}

type GlobalResiliencyRequestDraft struct {
	Name            string     `json:"name,omitempty"`
	DNSListenerName string     `json:"dns_listener_name,omitempty"`
	Protocols       []string   `json:"protocols,omitempty"`
	DNSListenerPort int        `json:"dns_listener_port,omitempty"`
	Id              string     `json:"id,omitempty"`
	Instances       []Instance `json:"instances,omitempty"`
}

type Instance struct {
	Hostname           string `json:"hostname,omitempty"`
	Address            string `json:"address,omitempty"`
	DNSListenerAddress string `json:"dns_listener_address,omitempty"`
	GroupSyncAddress   string `json:"group_sync_address,omitempty"`
}

type CMWAFReportRequestDraft struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	TimeFrameInDays int    `json:"time_frame_in_days,omitempty"`
	TopLevel        int    `json:"top_level,omitempty"`
	RequestType     string `json:"request_type,omitempty"`
	UserDefined     string `json:"user_defined,omitempty"`
	CreatedBy       string `json:"created_by,omitempty"`
	Scope           struct {
		Entity string   `json:"entity,omitempty"`
		All    bool     `json:"all"`
		Names  []string `json:"names,omitempty"`
	} `json:"scope,omitempty"`
	Categories []Category `json:"categories,omitempty"`
	Id         string     `json:"id,omitempty"`
}

type Category struct {
	Name string `json:"name,omitempty"`
}
type FastDeployVirtual struct {
	VirtualName    string `json:"virtualName,omitempty"`
	VirtualAddress string `json:"virtualAddress,omitempty"`
}
type FastDeployPoolMember struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}
type FastDeployPool struct {
	PoolName    string                 `json:"poolName,omitempty"`
	PoolMembers []FastDeployPoolMember `json:"poolMembers,omitempty"`
}

// create struct for above deploy request
type FastDeployRequest struct {
	Deployments []struct {
		Parameters struct {
			Pools    []FastDeployPool    `json:"pools,omitempty"`
			Virtuals []FastDeployVirtual `json:"virtuals,omitempty"`
		} `json:"parameters,omitempty"`
		Target struct {
			Address string `json:"address,omitempty"`
		} `json:"target,omitempty"`
		AllowOverwrite bool `json:"allow_overwrite,omitempty"`
	} `json:"deployments,omitempty"`
}

// create payload for Fast application draft deploy request using draftID
// mgmt/shared/fast/appsvcs/ac6b8145-c1bf-4140-b2cb-4358cc742931/deployments
// create POST request to deploy Fast application draft using FastDeployRequest and draftID
func (p *BigipNextCM) PostFastApplicationDraftDeployments(draftID string, config *FastDeployRequest) error {
	fastURL := fmt.Sprintf("%s%s%s/%s/%s", p.Host, uriFast, "/appsvcs", draftID, "deployments")
	f5osLogger.Info("[PostFastApplicationDraftDeployments]", "URI Path", fastURL)
	f5osLogger.Info("[PostFastApplicationDraftDeployments]", "Config", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return err
	}
	respData, err := p.doCMRequest("POST", fastURL, body)
	if err != nil {
		return err
	}
	f5osLogger.Info("[PostFastApplicationDraftDeployments]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// Create a Global Resiliency Groups
// /api/v1/spaces/default/gslb/gr-groups
func (p *BigipNextCM) PostGlobalResiliencyGroup(op string, config *GlobalResiliencyRequestDraft) (string, error) {
	grURL := fmt.Sprintf("%s%s", p.Host, uriGlobalResiliency)
	if op == "PUT" {
		grURL = fmt.Sprintf("%s%s/%s", p.Host, uriGlobalResiliency, config.Id)
	}
	f5osLogger.Info("[PostGlobalResiliencyGroup]", "URI Path", grURL)
	f5osLogger.Info("[PostGlobalResiliencyGroup]", "Config", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostGlobalResiliencyGroup]", "Body", hclog.Fmt("%+v", string(body)))
	respData, err := p.doCMRequest(op, grURL, body)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostGlobalResiliencyGroup]", "Data::", hclog.Fmt("%+v", string(respData)))

	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostGlobalResiliencyGroup]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	pathList := strings.Split(respString["path"].(string), "/")
	return p.GetGlobalResiliencyTaskStatus(pathList[len(pathList)-1])

}

// GET request to check the status of the Global Resiliency Group created
// /v1/spaces/default/gslb/gr-groups/{id}
func (p *BigipNextCM) GetGlobalResiliencyTaskStatus(taskId string) (string, error) {
	getTaskUrl := fmt.Sprintf("%s/%s", uriGetGlobalResiliency, taskId)
	f5osLogger.Info("[GetGlobalResiliencyTaskStatus]", "getTaskUrl", getTaskUrl)
	// poll the status is completed or failed until timeout
	var respInfo map[string]interface{}
	timeout := 60 * time.Second
	endtime := time.Now().Add(timeout)
	for time.Now().Before(endtime) {
		respData, err := p.GetCMRequest(getTaskUrl)
		if err != nil {
			return "", err
		}
		f5osLogger.Info("[GetGlobalResiliencyTaskStatus]", "Data::", hclog.Fmt("%+v", string(respData)))

		err = json.Unmarshal(respData, &respInfo)
		if err != nil {
			return "", err
		}
		//DEPLOYING
		if _, ok := respInfo["status"]; ok && respInfo["status"].(string) == "DEPLOYED" {
			return respInfo["id"].(string), nil
		}

		if _, ok := respInfo["status"]; ok && respInfo["status"].(string) == "FAILED" {
			return "", p.GetAlertMessage()
		}
		time.Sleep(10 * time.Second)
	}
	return "", fmt.Errorf("task status is still in :%+v within timeout period of:%+v", respInfo["status"].(string), timeout)
}

// GET request to get the details of the Global Resiliency Group
// /v1/spaces/default/gslb/gr-groups/{id}
func (p *BigipNextCM) GetGlobalResiliencyGroupDetails(id string) (interface{}, error) {
	getGlobalResiliencyGroupDetailsUrl := fmt.Sprintf("%s/%s", uriGetGlobalResiliency, id)
	f5osLogger.Info("[GetGlobalResiliencyGroupDetails]", "getGlobalResiliencyGroupDetails Url", getGlobalResiliencyGroupDetailsUrl)

	respData, err := p.GetCMRequest(getGlobalResiliencyGroupDetailsUrl)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[GetGlobalResiliencyGroupDetails]", "Data::", hclog.Fmt("%+v", string(respData)))

	var respInfo map[string]interface{}
	err = json.Unmarshal(respData, &respInfo)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[GetGlobalResiliencyGroupDetails]", "Resp Message", hclog.Fmt("%+v", respInfo))

	return respInfo, nil
}

// DELETE request to delete the Global Resiliency Group
// /v1/spaces/default/gslb/gr-groups/{id}
func (p *BigipNextCM) DeleteGlobalResiliencyGroup(id string) error {

	deleteGlobalResiliencyGroupUrl := fmt.Sprintf("%s%s/%s", p.Host, uriGlobalResiliency, id)
	f5osLogger.Info("[DeleteGlobalResiliencyGroup]", "URI Path", deleteGlobalResiliencyGroupUrl)

	respData, err := p.doCMRequest("DELETE", deleteGlobalResiliencyGroupUrl, nil)
	if err != nil {
		return err
	}
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteGlobalResiliencyGroup]", "Delete Response", hclog.Fmt("%+v", respString))

	timeout := 60 * time.Second
	endtime := time.Now().Add(timeout)
	for time.Now().Before(endtime) {
		respData, err := p.doCMRequest("GET", deleteGlobalResiliencyGroupUrl, nil)

		if err != nil {
			f5osLogger.Info("[DeleteGlobalResiliencyGroup]", "err status / code ", err.Error())

			if strings.Contains(err.Error(), "Requested Global Resiliency Group not found for the given group id") {
				f5osLogger.Info("[DeleteGlobalResiliencyGroup] Resiliency Group already deleted")
				return nil
			}
		}

		f5osLogger.Info("[DeleteGlobalResiliencyGroup]", "Data::", hclog.Fmt("%+v", string(respData)))
		time.Sleep(10 * time.Second)
	}

	return fmt.Errorf("task status is still in Running State within timeout period of:%+v", timeout)
}

// GET request to check summary of the last GR Group Deployed in case POST status is Failed
// /api/alert/v1/alerts/?limit=1&sort=-start_time&filter=source%20eq%20%27WAF%27%20and%20status%20eq%20%27ACTIVE%27&select=summary
func (p *BigipNextCM) GetAlertMessage() error {
	// getAlertUrl := fmt.Sprintf("%s", uriGetAlert)
	f5osLogger.Info("[GetAlertMessage]", "getAlertUrl", uriGetAlert)

	respData, err := p.GetCMRequest(uriGetAlert)
	if err != nil {
		return err
	}
	var respString map[string]interface{}
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return err
	}
	summary := respString["_embedded"].(map[string]interface{})["alerts"].([]interface{})[0].(map[string]interface{})["summary"].(string)
	f5osLogger.Info("[GetAlertMessage]", "Summary::", hclog.Fmt("%+v", string(summary)))
	return fmt.Errorf("task failed, summary : %+v ", summary)
}

// Create a WAF Security Report
// /api/v1/spaces/default/security/waf/reports
func (p *BigipNextCM) PostWAFReport(op string, config *CMWAFReportRequestDraft) (string, string, bool, error) {
	wafURL := fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, uriWafReport)
	if op == "PUT" {
		wafURL = fmt.Sprintf("%s%s%s/%s", p.Host, uriCMRoot, uriWafReport, config.Id)
	}
	f5osLogger.Info("[PostWAFReport]", "URI Path", wafURL)
	f5osLogger.Info("[PostWAFReport]", "Config", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return "", "", false, err
	}
	f5osLogger.Info("[PostWAFReport]", "Body", hclog.Fmt("%+v", string(body)))
	respData, err := p.doCMRequest(op, wafURL, body)
	if err != nil {
		return "", "", false, err
	}
	f5osLogger.Info("[PostWAFReport]", "Data::", hclog.Fmt("%+v", string(respData)))
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return "", "", false, err
	}
	f5osLogger.Info("[PostWAFReport]", "Task ID", hclog.Fmt("%+v", respString["id"].(string)))
	wafData, err := p.GetWAFReportDetails(respString["id"].(string))
	if err != nil {
		return "", "", false, err
	}
	return respString["id"].(string), wafData.(map[string]interface{})["created_by"].(string), wafData.(map[string]interface{})["user_defined"].(bool), nil

}

// GET request to get the details of the WAF Security Report
//
//	/api/v1/spaces/default/security/waf/reports{id}
func (p *BigipNextCM) GetWAFReportDetails(id string) (interface{}, error) {
	getWAFReportDetailsUrl := fmt.Sprintf("%s/%s", uriWafReport, id)
	f5osLogger.Info("[GetWAFReportDetails]", "GetWAFReportDetails Url", getWAFReportDetailsUrl)

	respData, err := p.GetCMRequest(getWAFReportDetailsUrl)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[GetWAFReportDetails]", "Data::", hclog.Fmt("%+v", string(respData)))

	var respInfo map[string]interface{}
	err = json.Unmarshal(respData, &respInfo)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[GetWAFReportDetails]", "Resp Message", hclog.Fmt("%+v", respInfo))

	return respInfo, nil
}

// DELETE request to delete the WAF Security Report
//
//	/api/v1/spaces/default/security/waf/reports{id}
func (p *BigipNextCM) DeleteWAFReport(id string) error {

	deleteWAFReportUrl := fmt.Sprintf("%s%s%s/%s", p.Host, uriCMRoot, uriWafReport, id)
	f5osLogger.Info("[DeleteWAFReport]", "URI Path", deleteWAFReportUrl)

	_, err := p.doCMRequest("DELETE", deleteWAFReportUrl, nil)
	if err != nil {
		return err
	}

	f5osLogger.Info("[DeleteWAFReport]", "WAF Report Deleted Successfully")
	return nil
}

type CMWAFPolicyRequestDraft struct {
	Name                string   `json:"name,omitempty"`
	Description         string   `json:"description,omitempty"`
	Tags                []string `json:"tags,omitempty"`
	EnforecementMode    string   `json:"enforcement_mode,omitempty"`
	ApplicationLanguage string   `json:"application_language,omitempty"`
	TemplateName        string   `json:"template_name,omitempty"`
	Declaration         struct {
		Policy struct {
			Name        string `json:"name,omitempty"`
			Description string `json:"description,omitempty"`
			Template    struct {
				Name string `json:"name,omitempty"`
			} `json:"template,omitempty"`
			BotDefense struct {
				Settings struct {
					IsEnabled bool `json:"isEnabled"`
				} `json:"settings"`
			} `json:"bot-defense"`
			IpIntelligence struct {
				Enabled bool `json:"enabled"`
			} `json:"ip-intelligence"`
			DosProtection struct {
				Enabled bool `json:"enabled"`
			} `json:"dos-protection"`
			BlockingSettings struct {
				Violations []Violation `json:"violations"`
			} `json:"blocking-settings"`
		} `json:"policy"`
	} `json:"declaration"`
	Id string `json:"id,omitempty"`
}

type Violation struct {
	Alarm       bool   `json:"alarm"`
	Block       bool   `json:"block"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
}

// Create a WAF Security Policy
// /api/v1/spaces/default/security/waf-policies
func (p *BigipNextCM) PostWAFPolicy(op string, config *CMWAFPolicyRequestDraft) (string, error) {
	wafURL := fmt.Sprintf("%s%s", p.Host, uriWafPolicy)
	if op == "PUT" {
		wafURL = fmt.Sprintf("%s%s/%s", p.Host, uriWafPolicy, config.Id)
	}
	f5osLogger.Info("[PostWAFPolicy]", "URI Path", wafURL)
	f5osLogger.Info("[PostWAFPolicy]", "Config", hclog.Fmt("%+v", config))

	body, err := json.Marshal(config)

	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostWAFPolicy]", "Body", hclog.Fmt("%+v", string(body)))
	respData, err := p.doCMRequest(op, wafURL, body)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostWAFPolicy]", "Data::", hclog.Fmt("%+v", string(respData)))

	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return "", err
	}
	f5osLogger.Info("[PostWAFPolicy]", "Task ID", hclog.Fmt("%+v", respString["id"].(string)))

	return respString["id"].(string), nil
}

// GET request to get the details of the WAF Policy
//
//	/api/v1/spaces/default/security/waf-policies/{id}
func (p *BigipNextCM) GetWAFPolicyDetails(id string) (interface{}, error) {
	getWAFPolicyDetailsUrl := fmt.Sprintf("%s/%s", uriGetWafPolicy, id)
	f5osLogger.Info("[GetWAFPolicyDetails]", "GetWAFPolicyDetails Url", getWAFPolicyDetailsUrl)
	respData, err := p.GetCMRequest(getWAFPolicyDetailsUrl)
	if err != nil {
		return "", err
	}
	// f5osLogger.Info("[GetWAFPolicyDetails]", "Data::", hclog.Fmt("%+v", string(respData)))
	var respInfo map[string]interface{}
	err = json.Unmarshal(respData, &respInfo)
	if err != nil {
		return "", err
	}
	return respInfo, nil
}

// DELETE request to delete the WAF Security Policy
//
//	/api/v1/spaces/default/security/waf-policies{id}
func (p *BigipNextCM) DeleteWAFPolicy(id string) error {
	deleteWAFPolicyUrl := fmt.Sprintf("%s%s/%s", p.Host, uriWafPolicy, id)
	f5osLogger.Info("[DeleteWAFPolicy]", "URI Path", deleteWAFPolicyUrl)
	resp, err := p.doCMRequest("DELETE", deleteWAFPolicyUrl, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteWAFPolicy]", "resp", string(resp))
	f5osLogger.Info("[DeleteWAFPolicy]", "WAF Policy Deleted Successfully")
	return nil
}

// https://clouddocs.f5.com/api/waf/v1/tasks/policy-import
// create a function to make policy-import

// func (p *BigipNextCM) PostWAFPolicyImport(config interface{}) (string, error) {
// 	wafPolicyImportURL := fmt.Sprintf("%s%s", p.Host, "/api/waf/v1/tasks/policy-import")
// 	f5osLogger.Info("[PostWAFPolicyImport]", "URI Path", wafPolicyImportURL)
// 	f5osLogger.Info("[PostWAFPolicyImport]", "Config", hclog.Fmt("%+v", config))
// 	body, err := json.Marshal(config)
// 	if err != nil {
// 		return "", err
// 	}
// 	respData, err := p.doCMRequest("POST", wafPolicyImportURL, body)
// 	if err != nil {
// 		return "", err
// 	}
// 	f5osLogger.Info("[PostWAFPolicyImport]", "Data::", hclog.Fmt("%+v", string(respData)))

// 	return "", nil
// }

// mgmt/shared/fast/appsvcs/ac6b8145-c1bf-4140-b2cb-4358cc742931/deployments
// create GET request to get Fast application draft deployments
func (p *BigipNextCM) GetFastApplicationDraftDeployments(draftID string) error {
	fastURL := fmt.Sprintf("%s%s%s/%s", p.Host, uriFast, "/appsvcs", draftID)
	f5osLogger.Info("[GetFastApplicationDraftDeployments]", "URI Path", fastURL)
	respData, err := p.doCMRequest("GET", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[GetFastApplicationDraftDeployments]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// mgmt/shared/fast/appsvcs/ac6b8145-c1bf-4140-b2cb-4358cc742931/deployments
// create DELETE request to delete Fast application draft deployments
func (p *BigipNextCM) DeleteFastApplicationDraftDeployments(draftID string) error {
	// fastURL := fmt.Sprintf("%s%s%s/%s/%s", p.Host, uriFast, "/appsvcs", draftID, "deployments")
	fastURL := fmt.Sprintf("%s%s%s/%s", p.Host, uriFast, "/appsvcs", draftID)
	f5osLogger.Info("[DeleteFastApplicationDraftDeployments]", "URI Path", fastURL)
	respData, err := p.doCMRequest("DELETE", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteFastApplicationDraftDeployments]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// mgmt/shared/fast/schemas/tasks/ID
// create GET request to get Fast Schema task status
func (p *BigipNextCM) GetFastSchemaTaskStatus(taskid string) error {
	fastURL := fmt.Sprintf("%s%s%s%s", p.Host, uriFast, "/schemas/", taskid)
	f5osLogger.Info("[GetFastSchemaTaskStatus]", "URI Path", fastURL)
	respData, err := p.doCMRequest("GET", fastURL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[GetFastSchemaTaskStatus]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// https://{{bigip_next_cm_mgmt_ip}}/mgmt/shared/fast/applications?filter=instanceID eq '{{bigip_next_cm_device_id}}'
func (p *BigipNextCM) GetAs3TargetInfo(target string) error {
	as3URL := fmt.Sprintf("%s%s/%s?target=%s", p.Host, uriAs3, "info", target)
	f5osLogger.Info("[GetAs3TargetInfo]", "URI Path", as3URL)
	respData, err := p.doCMRequest("GET", as3URL, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[GetAs3TargetInfo]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

func (p *BigipNextCM) CMTokenRefresh() (*BigipNextCM, error) {
	tokenRefreshUrl := fmt.Sprintf("%s", "/token-refresh")
	f5osLogger.Info("[CMTokenRefresh]", "tokenRefreshUrl", tokenRefreshUrl)
	tokenPayload := make(map[string]interface{})
	tokenPayload["refresh_token"] = p.RefreshToken
	tokenData, err := json.Marshal(tokenPayload)
	if err != nil {
		return nil, err
	}
	respData, err := p.PostCMRequest(tokenRefreshUrl, tokenData)
	if err != nil {
		return nil, err
	}
	var resp BigipNextCMLoginResp
	err = json.Unmarshal(respData, &resp)
	if err != nil {
		return nil, err
	}
	p.Token = resp.AccessToken
	if resp.RefreshToken != "" {
		p.RefreshToken = resp.RefreshToken
	}
	return p, nil
}

func (p *BigipNextCM) CMTokenRefreshNew() error {
	tokenRefreshUrl := "/token-refresh"
	f5osLogger.Info("[CMTokenRefresh]", "tokenRefreshUrl", tokenRefreshUrl)
	tokenPayload := make(map[string]interface{})
	tokenPayload["refresh_token"] = p.RefreshToken
	tokenData, err := json.Marshal(tokenPayload)
	if err != nil {
		return err
	}
	respData, err := p.PostCMRequest(tokenRefreshUrl, tokenData)
	if err != nil {
		return err
	}
	var resp BigipNextCMLoginResp
	err = json.Unmarshal(respData, &resp)
	if err != nil {
		return err
	}
	p.Token = resp.AccessToken
	if resp.RefreshToken != "" {
		p.RefreshToken = resp.RefreshToken
	}
	return nil
}

func (p *BigipNextCM) GetProxyFiles(proxyID string) ([]byte, error) {
	proxyFileUrl := fmt.Sprintf("%s/%s?path=/%s", "/device/v1/proxy/", proxyID, "files")
	f5osLogger.Info("[GetProxyFiles]", "proxyFileUrl", proxyFileUrl)
	respData, err := p.GetCMRequest(proxyFileUrl)
	if err != nil {
		return []byte(""), err
	}
	f5osLogger.Info("[GetAs3TargetInfo]", "Data::", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

func (p *BigipNextCM) ProxyFileUpload(proxyID, filePath string) ([]byte, error) {
	proxyFileUploadUrl := fmt.Sprintf("%s/%s", uriCMProxyFileUpload, proxyID)
	f5osLogger.Info("[ProxyFileUpload]", "proxyFileUploadUrl", proxyFileUploadUrl)
	fileObj, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	fileInfo, err := fileObj.Stat()
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[ProxyFileUpload]", "fileInfo", fileInfo)
	return nil, nil
}

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
	f5osLogger.Info("[getAS3DeploymentTaskStatus]", "Response Result:", hclog.Fmt("%+v", responseData["response"]))
	// .(map[string]interface{})["results"].([]interface{})[0].(map[string]interface{})["status"].(string)))
	// if responseData["records"].([]interface{})[0].(map[string]interface{})["status"].(string) != "completed" {
	// 	return nil, fmt.Errorf("AS3 service deployment failed with :%+v", responseData["records"].([]interface{})[0].(map[string]interface{})["status"].(string))
	// }
	byteData, err := json.Marshal(responseData["app_data"].(map[string]interface{}))
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

type DiscoverInstanceRequest struct {
	Address            string `json:"address,omitempty"`
	Port               int    `json:"port,omitempty"`
	DeviceUser         string `json:"device_user,omitempty"`
	DevicePassword     string `json:"device_password,omitempty"`
	ManagementUser     string `json:"management_user,omitempty"`
	ManagementPassword string `json:"management_password,omitempty"`
}

// create POST request to Add instance to CM
func (p *BigipNextCM) DiscoverInstance(config *DiscoverInstanceRequest) ([]byte, error) {
	if config.DevicePassword == "admin" {
		err := config.resetDevicePassword()
		if err != nil {
			return nil, err
		}
		f5osLogger.Info("[DiscoverInstance]", "admin password reset successfully")
		time.Sleep(2 * time.Second)
		config.DevicePassword = config.ManagementPassword
	}
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.PostCMRequest(uriDiscoverInstance, body)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[DiscoverInstance]", "Data::", hclog.Fmt("%+v", string(respData)))
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[DiscoverInstance]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	pathList := strings.Split(respString["path"].(string), "/")

	err = p.acceptUntrustedCertificate(pathList[len(pathList)-1])
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[getDiscoverInstanceTaskStatus]", "Data::", hclog.Fmt("%+v", string(respData)))
	respData, err = p.getDiscoverInstanceTaskStatus(pathList[len(pathList)-1])
	if err != nil {
		return nil, err
	}
	return respData, nil
}

// Accept untrusted certificate for the device
func (p *BigipNextCM) acceptUntrustedCertificate(taskid string) error {
	acceptUntrust := true
	unTrust := make(map[string]interface{})
	if acceptUntrust {
		unTrust["is_user_accepted_untrusted_cert"] = true
	}
	body, err := json.Marshal(unTrust)
	if err != nil {
		return err
	}
	getTaskUrl := fmt.Sprintf("%s/api%s%s/%s", p.Host, uriDiscoverInstance, "/discovery-tasks", taskid)
	respData, err := p.doCMRequest("PATCH", getTaskUrl, body)
	if err != nil {
		return err
	}
	f5osLogger.Info("[AcceptUntrustedCertificate]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// Check the status of the discovery task
func (p *BigipNextCM) getDiscoverInstanceTaskStatus(taskid string) ([]byte, error) {
	getTaskUrl := fmt.Sprintf("%s%s/%s", uriDiscoverInstance, "/discovery-tasks", taskid)
	f5osLogger.Info("[getDiscoverInstanceTaskStatus]", "getTaskUrl", getTaskUrl)
	var respInfo map[string]interface{}
	timeout := 360 * time.Second
	endtime := time.Now().Add(timeout)
	for time.Now().Before(endtime) {
		respData, err := p.GetCMRequest(getTaskUrl)
		if err != nil {
			return nil, err
		}
		f5osLogger.Info("[getDiscoverInstanceTaskStatus]", "Data::", hclog.Fmt("%+v", string(respData)))
		// {"_links":{"self":{"href":"/api/v1/spaces/default/instances/discovery-tasks/2e718d16-66af-4a11-960a-cd2dfcf48229"}},"address":"10.145.71.115","created":"2024-04-05T17:43:27.382035Z","device_group":"default","device_user":"admin","fingerprint":"771caf5eaf0718911c4da754fd7bc998797066992c6ebb6129f5dcf58528aba4","id":"2e718d16-66af-4a11-960a-cd2dfcf48229","port":5443,"state":"discoveryWaitForUserInput","status":"running"}
		err = json.Unmarshal(respData, &respInfo)
		if err != nil {
			return nil, err
		}
		if respInfo["status"].(string) == "completed" {
			return []byte(respInfo["discovered_device_id"].(string)), nil
		}
		if respInfo["status"].(string) == "failed" {
			return respData, fmt.Errorf("discovery-tasks failed with :%+v", respInfo["failure_reason"].(string))
		}
		time.Sleep(10 * time.Second)
	}
	return []byte(""), fmt.Errorf("task status is still in :%+v within timeout period of:%+v", respInfo["status"].(string), timeout)
}

// reset the device password
func (d *DiscoverInstanceRequest) resetDevicePassword() error {
	urlString := fmt.Sprintf("https://%s:%d%s", d.Address, d.Port, "/api/v1/me")
	f5osLogger.Info("[resetDevicePassword]", "getTaskUrl", urlString)
	resetPassword := make(map[string]interface{})
	resetPassword["currentPassword"] = d.DevicePassword
	resetPassword["newPassword"] = d.ManagementPassword
	body, err := json.Marshal(resetPassword)
	if err != nil {
		return err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: tr,
	}
	method := "PUT"
	f5osLogger.Info("[resetDevicePassword]", "URL", hclog.Fmt("%+v", urlString))
	req, err := http.NewRequest(method, urlString, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentTypeHeader)
	req.SetBasicAuth(d.DeviceUser, string(d.DevicePassword))
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	bodyResp, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 204 {
		return fmt.Errorf("error message: %+v (Body:%+v)", res.Status, string(bodyResp))
	}
	return nil
}

// func standardizeSpaces(s string) string {
// 	return strings.Join(strings.Fields(s), " ")
// }

func (p *BigipNextCM) PostNextCMAs3(as3Json string) error {
	as3URL := fmt.Sprintf("%s%s/%s", p.Host, uriAs3, "declare")
	f5osLogger.Info("[PostNextCMAs3]", "URI Path", as3URL)
	respData, err := p.doCMRequest("POST", as3URL, []byte(as3Json))
	if err != nil {
		return err
	}
	f5osLogger.Info("[PostNextCMAs3]", "AS3 Task Status:", hclog.Fmt("%+v", string(respData)))
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return err
	}
	for _, v := range respString["results"].([]interface{}) {
		if int(v.(map[string]interface{})["code"].(float64)) != 200 {
			return fmt.Errorf("posting AS3 failed with :%+v", respString["results"])
		}
	}
	return nil
}

func (p *BigipNextCM) GetBigIPNextCMAs3(target, tenantList string) ([]byte, error) {
	as3URL := fmt.Sprintf("%s%s/%s?target=%s", p.Host, uriAs3, "declare", target)
	if tenantList != "" {
		as3URL = fmt.Sprintf("%s%s/%s/%s?target=%s", p.Host, uriAs3, "declare", tenantList, target)
	}
	f5osLogger.Info("[GetBigIPNextCMAs3]", "URI Path", as3URL)
	respData, err := p.doCMRequest("GET", as3URL, nil)
	if err != nil {
		return []byte(""), err
	}
	as3Json := make(map[string]interface{})
	as3Json["class"] = "AS3"
	as3Json["action"] = "deploy"
	as3Json["persist"] = true
	adcJson := make(map[string]interface{})
	err = json.Unmarshal(respData, &adcJson)
	if err != nil {
		return []byte(""), err
	}
	adcJson["target"] = map[string]interface{}{"address": target}
	as3Json["declaration"] = adcJson
	as3data, err := json.Marshal(as3Json)
	if err != nil {
		return []byte(""), err
	}
	f5osLogger.Info("[GetBigIPNextCMAs3]", "Data::", hclog.Fmt("%+v", string(as3data)))
	return as3data, nil
}

func (p *BigipNextCM) DeleteNextCMAs3(target, tenantName string) error {
	as3URL := fmt.Sprintf("%s%s/%s/%s?target=%s", p.Host, uriAs3, "declare", tenantName, target)
	f5osLogger.Info("[DeleteNextCMAs3]", "URI Path", as3URL)
	respData, err := p.doCMRequest("DELETE", as3URL, nil)
	// {"_links":{"self":{"href":"/delete-tenant-tasks/58c864e5-fe44-4270-9e16-386d06b19a40"}},"path":"/delete-tenant-tasks/58c864e5-fe44-4270-9e16-386d06b19a40"}
	if err != nil {
		return err
	}
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteNextCMAs3]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	err = p.deleteTenantTaskStatus(respString["path"].(string))
	if err != nil {
		return err
	}
	return nil
}

func (p *BigipNextCM) deleteTenantTaskStatus(taskidPath string) error {
	as3URL := fmt.Sprintf("%s%s%s", p.Host, uriAs3, taskidPath)
	f5osLogger.Info("[deleteTenantTaskStatus]", "URI Path", as3URL)
	timeout := 60 * time.Second
	endtime := time.Now().Add(timeout)
	respString := make(map[string]interface{})
	for time.Now().Before(endtime) {
		respData, err := p.doCMRequest("GET", as3URL, nil)
		// {"_links":{"self":{"href":"/delete-tenant-tasks/58c864e5-fe44-4270-9e16-386d06b19a40"}},"completed":"2023-07-24T14:40:00.668098Z","created":"2023-07-24T14:39:43.382163Z","failure_reason":"","id":"58c864e5-fe44-4270-9e16-386d06b19a40","instance_id":"ca588f7a-5ecf-4080-b168-733b55636cfc","name":"delete AS3 tenant next-cm-tenant02","state":"delDone","status":"completed","task_type":"as3_tenant_deletion","tenant_name":"next-cm-tenant02"}
		if err != nil {
			return err
		}
		f5osLogger.Info("[deleteTenantTaskStatus]", "Task Status:\t", hclog.Fmt("%+v", string(respData)))
		err = json.Unmarshal(respData, &respString)
		if err != nil {
			return err
		}
		f5osLogger.Info("[deleteTenantTaskStatus]", "Task Status", hclog.Fmt("%+v", respString["status"].(string)))
		if respString["status"].(string) == "completed" && respString["state"].(string) == "delDone" {
			return nil
		}
		if respString["status"].(string) == "failed" {
			return fmt.Errorf("%s", respString)
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("task status is still in :%+v within timeout period of:%+v", respString["status"].(string), timeout)
}

func (p *BigipNextCM) GetTargetTenantList(body interface{}) (string, string) {
	tenantList := make([]string, 0)
	applicationList := make([]string, 0)
	as3json := body.(string)
	resp := []byte(as3json)
	var targetDevice string
	jsonRef := make(map[string]interface{})
	json.Unmarshal(resp, &jsonRef)
	for key, value := range jsonRef {
		if rec, ok := value.(map[string]interface{}); ok && key == "declaration" {
			f5osLogger.Info("[GetTargetTenantList]", "key value", hclog.Fmt("%+v", rec["target"]))
			if rec["target"] != nil && rec["target"].(map[string]interface{})["address"].(string) != "" {
				targetDevice = rec["target"].(map[string]interface{})["address"].(string)
			}
			for k, v := range rec {

				if rec2, ok := v.(map[string]interface{}); ok {
					found := 0
					for k1, v1 := range rec2 {
						if k1 == "class" && v1 == "Tenant" {
							found = 1
						}
						if rec3, ok := v1.(map[string]interface{}); ok {
							found1 := 0
							for k2, v2 := range rec3 {
								if k2 == "class" && v2 == "Application" {
									found1 = 1
								}
							}
							if found1 == 1 {
								applicationList = append(applicationList, k1)
							}

						}
					}
					if found == 1 {
						tenantList = append(tenantList, k)
					}
				}
			}
		}
	}
	finalTenantlist := strings.Join(tenantList[:], ",")
	return finalTenantlist, targetDevice
}

// https://10.144.73.240/api/system/v1/files
// multipart/form-data

// ------WebKitFormBoundary1mZMGjaktKFrIxOX
// Content-Disposition: form-data; name="content"; filename="BIG-IP-Next-CentralManager-20.1.0-0.8.114-Update.tgz"
// Content-Type: application/octet-stream

// ------WebKitFormBoundary1mZMGjaktKFrIxOX
// Content-Disposition: form-data; name="file_name"

// BIG-IP-Next-CentralManager-20.1.0-0.8.114-Update.tgz
// ------WebKitFormBoundary1mZMGjaktKFrIxOX
// Content-Disposition: form-data; name="description"

// CM upgrade
// ------WebKitFormBoundary1mZMGjaktKFrIxOX--

// reqBody := new(bytes.Buffer)
// mp := multipart.NewWriter(reqBody)
// for k, v := range bodyInput {
//   str, ok := v.(string)
//   if !ok {
//     return fmt.Errorf("converting %v to string", v)
//   }
//   mp.WriteField(k, str)
// }
// mp.Close()

// req, err := http.NewRequest(http.MethodPost, "https://my-website.com/endpoint/path", reqBody)
// if err != nil {
// // handle err
// }
// req.Header["Content-Type"] = []string{mp.FormDataContentType()}

type PolicyimportReqObj struct {
	FilePath    string `json:"file_path,omitempty"`
	PolicyName  string `json:"policy_name,omitempty"`
	Description string `json:"description,omitempty"`
	Override    string `json:"override,omitempty"`
}

const filechunk = 8192

// func (p *BigipNextCM) GetMd5ofFile(filePath string) (interface{}, error) {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer file.Close()

// 	// Get file info
// 	info, err := file.Stat()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get the filesize
// 	filesize := info.Size()

// 	// Calculate the number of blocks
// 	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

// 	// Start hash
// 	hash := md5.New()

// 	// Check each block
// 	for i := uint64(0); i < blocks; i++ {
// 		// Calculate block size
// 		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))

// 		// Make a buffer
// 		buf := make([]byte, blocksize)

// 		// Make a buffer
// 		file.Read(buf)

// 		// Write to the buffer
// 		io.WriteString(hash, string(buf))
// 	}
// 	return hash.Sum(nil), nil
// }

func (p *BigipNextCM) PolicyImport(config *PolicyimportReqObj) (interface{}, error) {
	// func (p *BigipNextCM) PolicyImport(filePath, policyName, description, override string) ([]byte, error) {
	body := &bytes.Buffer{}
	file, err := os.Open(config.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	writer := multipart.NewWriter(body)
	policyName := config.PolicyName
	override := config.Override
	description := config.Description
	if err = WriteField(writer, "policy_name", policyName); err != nil {
		return nil, err
	}
	if err = WriteField(writer, "description", description); err != nil {
		return nil, err
	}
	if err = WriteField(writer, "override", override); err != nil {
		return nil, err
	}
	if err = WriteFile(writer, "content", file); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil { // finishes the multipart message and writes the trailing boundary
		return nil, err
	}
	url := "/waf/v1/tasks/policy-import"
	url = fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, url)
	f5osLogger.Info("[PolicyImport]", "URL ", hclog.Fmt("%+v", url))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Token))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Transport: p.Transport,
		Timeout:   30 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 202 {
		//{"_links":{"self":{"href":"/api/waf/v1/tasks/policy-import/661f6053-7a8c-4f1d-9259-2f1c001490f4"}},"path":"/v1/tasks/policy-import/661f6053-7a8c-4f1d-9259-2f1c001490f4"}
		dataResp, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil
		}
		respData := make(map[string]interface{})
		err = json.Unmarshal(dataResp, &respData)
		if err != nil {
			return nil, nil
		}
		// get ID from path key
		pathList := strings.Split(respData["path"].(string), "/")
		return p.PolicyImportStatus(pathList[len(pathList)-1], 100)
	}
	if resp.StatusCode >= 400 {
		byteData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(`{"code":%d,"error":%s`, resp.StatusCode, byteData)
	}
	return nil, nil
}

func (p *BigipNextCM) PolicyImportStatus(taskID string, timeOut int) (interface{}, error) {
	importUrl := fmt.Sprintf("%s%s", "/waf/v1/tasks/policy-import/", taskID)
	f5osLogger.Info("[PolicyImportStatus]", "URI Path", importUrl)
	taskData := make(map[string]interface{})
	// var timeout time.Duration
	timeout := time.Duration(timeOut) * time.Second
	endtime := time.Now().Add(timeout)
	for time.Now().Before(endtime) {
		respData, err := p.GetCMRequest(importUrl)
		if err != nil {
			return nil, err
		}
		f5osLogger.Info("[PolicyImportStatus]", "Task Status:\t", hclog.Fmt("%+v", string(respData)))
		err = json.Unmarshal(respData, &taskData)
		if err != nil {
			return nil, err
		}
		if taskData["status"] == "completed" {
			return taskData, nil
		}
		if taskData["status"] == "failed" {
			return nil, fmt.Errorf("%s", taskData["failure_reason"])
		}
		inVal := timeOut / 10
		time.Sleep(time.Duration(inVal) * time.Second)
	}
	return nil, fmt.Errorf("task Status is still in :%+v within timeout period of:%+v", taskData["status"], timeout)
}

func (p *BigipNextCM) FileImportBackup(url string, values map[string]io.Reader) ([]byte, error) {
	// Prepare a form that you will submit to that URL.
	// var b bytes.Buffer
	b := &bytes.Buffer{}
	var err error
	w := multipart.NewWriter(b)
	// writer := multipart.NewWriter(body)

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	defer w.Close()

	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return []byte(""), nil
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return []byte(""), nil
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return []byte(""), nil
		}

	}

	url = fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, url)
	f5osLogger.Info("[FileImport]", "URL ", hclog.Fmt("%+v", url))
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return []byte(""), nil
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Token))
	req.Header.Add("Content-Type", w.FormDataContentType())

	client := &http.Client{
		Transport: p.Transport,
		Timeout:   30 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 202 {
		return io.ReadAll(resp.Body)
	}
	if resp.StatusCode >= 400 {
		byteData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(`{"code":%d,"error":%s`, resp.StatusCode, byteData)
	}
	return []byte("FileImport Success"), nil
}

// create multi-part form data for file upload
func (p *BigipNextCM) createMultipartFormData(filePath string) (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	f5osLogger.Info("[createMultipartFormData]", "file_name ", hclog.Fmt("%+v", filepath.Base(filePath)))
	// Adding MIMEHeader
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "content", filepath.Base(file.Name())))
	h.Set("Content-Type", "application/octet-stream")
	part, err := writer.CreatePart(h)

	// part, err := writer.CreateFormFile("content", filepath.Base(filePath))

	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, "", err
	}
	fw, err := writer.CreateFormField("description")
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(fw, strings.NewReader("CM upgrade"))
	if err != nil {
		return nil, "", err
	}
	fw1, err := writer.CreateFormField("file_name")
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(fw1, strings.NewReader(filepath.Base(filePath)))
	if err != nil {
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return body, writer.FormDataContentType(), nil
}

// curl -vk --location -H "Authorization: Bearer $TOKEN" https://10.145.67.139:443/api/system/v1/files \
//     -H 'Content-Type: multipart/form-data' \
//     -F "file_name=$LOCAL_LAST_NEXT_VE_UPGRADE" \
//     -F "content=@/Users/r.chinthalapalli/Downloads/$LOCAL_LAST_NEXT_VE_UPGRADE" \
//     -F "description=CM upgrade"

// convert above curl command to GO Lang code

// upload file using multi-part form data
func (p *BigipNextCM) uploadFile(filePath string) ([]byte, error) {
	body, contentType, err := p.createMultipartFormData(filePath)
	if err != nil {
		return nil, err
	}
	// url := fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, uriCMFileUpload)
	url := fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, "/v1/spaces/default/files")
	f5osLogger.Info("[uploadFile]", "url", hclog.Fmt("%+v", url))
	f5osLogger.Info("[uploadFile]", "content-type", hclog.Fmt("%+v", contentType))

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Token))
	req.Header.Add("Content-Type", contentType)
	client := &http.Client{
		Transport: p.Transport,
		Timeout:   30 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 202 {
		return io.ReadAll(resp.Body)
	}
	if resp.StatusCode >= 400 {
		byteData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(`{"code":%d,"error":%s`, resp.StatusCode, byteData)
	}

	// nextChan := make(chan *BigipNextCM)
	// go func() {
	// 	for {
	// 		time.Sleep(3 * time.Minute)
	// 		err := p.CMTokenRefreshNew()
	// 		if err != nil {
	// 			f5osLogger.Error("Error refreshing CM token", "error", err)
	// 		}
	// 		nextChan <- p
	// 	}
	// }()
	// go func() {
	// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Token))
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	defer resp.Body.Close()
	// 	if resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 202 {
	// 		return io.ReadAll(resp.Body)
	// 	}

	// 	if resp.StatusCode == 401 {
	// 		continue
	// 	}
	// 	if resp.StatusCode >= 400 {
	// 		byteData, _ := io.ReadAll(resp.Body)
	// 		return nil, fmt.Errorf(`{"code":%d,"error":%s`, resp.StatusCode, byteData)
	// 	}
	// }()

	// for {
	// 	p = <-nextChan

	// }

	return []byte(""), nil
}

// make request to upload muti-part form data file
func (p *BigipNextCM) UploadFile(filePath string) ([]byte, error) {
	// return p.uploadFile(filePath)
	// return p.UploadFileWithTokenRefresh(filePath)
	// return p.UploadFileWithTokenRefreshbackup(filePath)
	return p.uploadFileWithRefresh(filePath)
}

// make request to upload muti-part form data file with token refresh
func (p *BigipNextCM) UploadFileWithTokenRefreshbacup(filePath string) ([]byte, error) {
	var err error
	var respData []byte
	for i := 0; i < 10; i++ {
		respData, err = p.uploadFile(filePath)
		if err != nil {
			if strings.Contains(err.Error(), "401") {
				p, err = p.CMTokenRefresh()
				if err != nil {
					return nil, err
				}
				continue
			}
			return nil, err
		}
		break
	}
	return respData, nil
}

// make request to upload muti-part form data file with token refresh
func (p *BigipNextCM) UploadFileWithTokenRefresh(filePath string) ([]byte, error) {
	go func() {
		for {
			time.Sleep(2 * time.Minute)
			err := p.CMTokenRefreshNew()
			if err != nil {
				f5osLogger.Error("Error refreshing CM token", "error", err)
			} else {
				f5osLogger.Info("Refreshed CM token successfully")
			}
		}
	}()

	respData, err := p.uploadFile(filePath)
	if err != nil {
		return []byte(""), err
	}
	return respData, nil

	// var err error
	// var respData []byte
	// for i := 0; i < 10; i++ {
	// 	respData, err = p.uploadFile(filePath)
	// 	if err != nil {
	// 		if strings.Contains(err.Error(), "401") {
	// 			p, err = p.CMTokenRefresh()
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			continue
	// 		}
	// 		return nil, err
	// 	}
	// 	break
	// }
	// return respData, nil
}

// create func to run CMRefresh Token and uploadFile parallel
func (p *BigipNextCM) uploadFileWithRefresh(filePath string) ([]byte, error) {
	nextChan := make(chan *BigipNextCM)
	var resp []byte
	respChan := make(chan []byte)
	errChan := make(chan error)
	var err error
	go func() {
		nextChan <- p
		resp, err = p.uploadFile(filePath)
		if err != nil {
			errChan <- err
		}
		respChan <- resp
	}()
	go func() {
		p = <-nextChan
		for {
			time.Sleep(3 * time.Minute)
			err := p.CMTokenRefreshNew()
			if err != nil {
				f5osLogger.Error("Error refreshing CM token", "error", err)
			}
		}
	}()
	err = <-errChan
	if err != nil {
		return nil, err
	}
	return <-respChan, nil
}

// write above logic using channel

// create http request with large file using multi part form-data and token refresh

// create a function to add two numbers

// https://10.144.73.240/api/system/v1/files
// filename query parameter
func (p *BigipNextCM) GetFile(fileName string) ([]byte, error) {
	fileUrl := fmt.Sprintf("%s?filter=file_name+eq+'%s'", uriCMFileUpload, fileName)
	f5osLogger.Info("[GetFile]", "fileUrl", fileUrl)
	respData, err := p.GetCMRequest(fileUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetFile]", "Data::", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

// https://10.145.79.7/api/upgrade-manager/v1/upgrade-tasks/
// { "file_id": "8281d20b-34f1-4c14-9812-9af468b472bd", "file_format":"tar"}
func (p *BigipNextCM) PostUpgradeTask(fileId string) (interface{}, error) {
	// upgradeTaskUrl := fmt.Sprintf("%s", uriCMUpgradeTask)
	f5osLogger.Info("[PostUpgradeTask]", "upgradeTaskUrl", uriCMUpgradeTask)
	upgradeTaskPayload := make(map[string]interface{})
	upgradeTaskPayload["file_id"] = fileId
	upgradeTaskPayload["file_format"] = "tar"
	upgradeTaskData, err := json.Marshal(upgradeTaskPayload)
	if err != nil {
		return nil, err
	}
	respData, err := p.PutCMRequest(uriCMUpgradeTask, upgradeTaskData)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[PostUpgradeTask]", "Data::", hclog.Fmt("%+v", string(respData)))
	// {"_links":{"self":{"href":"/v1/upgrade-tasks/169653ef-1a15-4d05-b090-c07cb8ac5a43"}},"path":"/v1/upgrade-tasks/169653ef-1a15-4d05-b090-c07cb8ac5a43"}
	var respInfo map[string]interface{}
	err = json.Unmarshal(respData, &respInfo)
	if err != nil {
		return nil, err
	}
	// check if path key is present in response
	if _, ok := respInfo["path"]; !ok {
		return nil, fmt.Errorf("upgrade task failed with :%+v", respInfo)
	}
	// get ID from path key
	pathList := strings.Split(respInfo["path"].(string), "/")
	return p.GetUpgradeTaskStatus(pathList[len(pathList)-1], 1200)
	// return []byte(pathList[len(pathList)-1]), nil
}

// https://10.144.73.240/api/upgrade-manager/v1/upgrade-tasks
func (p *BigipNextCM) GetUpgradeTaskStatus(taskId string, timeOut int) (interface{}, error) {
	upgradeTaskUrl := fmt.Sprintf("%s/%s", uriCMUpgradeTask, taskId)
	f5osLogger.Info("[GetUpgradeTaskStatus]", "upgradeTaskUrl", upgradeTaskUrl)
	respData, err := p.GetCMRequest(upgradeTaskUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetUpgradeTaskStatus]", "Data::", hclog.Fmt("%+v", string(respData)))
	var respInfo map[string]interface{}
	err = json.Unmarshal(respData, &respInfo)
	if err != nil {
		return nil, err
	}
	// poll the status is completed or failed until timeout
	timeout := time.Duration(timeOut) * time.Second
	endtime := time.Now().Add(timeout)
	for time.Now().Before(endtime) {
		if _, ok := respInfo["status"]; ok && respInfo["status"].(string) == "completed" {
			return respInfo, nil
		}
		if _, ok := respInfo["status"]; ok && respInfo["status"].(string) == "failed" {
			// {
			// 	"_links": {
			// 		"self": {
			// 			"href": "/v1/upgrade-tasks/f3fbba78-8d87-46f6-a18c-ab3b5486bf42"
			// 		}
			// 	},
			// 	"completed": "2024-02-21T17:24:31.022646Z",
			// 	"created": "2024-02-21T17:24:29.989491Z",
			// 	"failure_reason": "unable to unarchive tgz file opening tar archive for reading: wrapping file reader: gzip: invalid header",
			// 	"file_id": "793bd34e-9a39-4299-a1c0-8c0d5e1ade6a",
			// 	"id": "f3fbba78-8d87-46f6-a18c-ab3b5486bf42",
			// 	"state": "unpackUpgradeFiles",
			// 	"status": "failed"
			// }
			return nil, fmt.Errorf("upgrade task failed with :%+v", respInfo)
		}
		time.Sleep(time.Duration(timeOut/10) * time.Second)
		respData, err = p.GetCMRequest(upgradeTaskUrl)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(respData, &respInfo)
		if err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("upgrade task status is still in :%+v within timeout period of:%+v", respInfo["status"].(string), timeout)
}

// curl -ks -H "Authorization: Bearer $TOKEN" -F file_name=cm-install-bundle.tgz -F content=@cm-install-bundle.tgz';'type=application/octet-stream 'https://10.145.79.7/api/system/v1/files

//
// //check if status key is present in response and status is not completed
// {
// 	"_links": {
// 		"self": {
// 			"href": "/v1/upgrade-tasks/f3fbba78-8d87-46f6-a18c-ab3b5486bf42"
// 		}
// 	},
// 	"created": "2024-02-21T17:24:29.989491Z",
// 	"file_id": "793bd34e-9a39-4299-a1c0-8c0d5e1ade6a",
// 	"id": "f3fbba78-8d87-46f6-a18c-ab3b5486bf42",
// 	"state": "unpackUpgradeFiles",
// 	"status": "running"
// }
// {
// 	"_links": {
// 		"self": {
// 			"href": "/v1/upgrade-tasks/bfb8100b-7798-4262-8873-61e856652843"
// 		}
// 	},
// 	"completed": "2023-12-28T15:00:07.241268Z",
// 	"created": "2023-12-28T14:44:57.955703Z",
// 	"failure_reason": "",
// 	"file_id": "bd069deb-00af-429a-bfc0-24e11cbd152f",
// 	"id": "bfb8100b-7798-4262-8873-61e856652843",
// 	"state": "done",
// 	"status": "completed"
// }

func (p *BigipNextCM) DeleteFile(fileName string) ([]byte, error) {
	fileUrl := fmt.Sprintf("%s?filter=file_name+eq+'%s'", uriCMFileUpload, fileName)
	f5osLogger.Info("[DeleteFile]", "fileUrl", fileUrl)
	respData, err := p.GetCMRequest(fileUrl)
	if err != nil {
		return nil, err
	}
	var respInfo map[string]interface{}
	json.Unmarshal(respData, &respInfo)
	if _, ok := respInfo["_embedded"]; !ok {
		return nil, fmt.Errorf("the requested file:%s, was not found", fileName)
	}
	//get the ID of the file
	//{"_embedded":{"files":[{"_links":{"self":{"href":"/v1/files?filter=file_name+eq+%27BIG-IP-Next-CentralManager-20.1.0-0.8.115-Update.tgz%27/ae0a842a-7ed5-44b6-98ae-eed553695818"}},"description":"CM upgrade","file_name":"BIG-IP-Next-CentralManager-20.1.0-0.8.115-Update.tgz","file_size":345,"hash":"a240520a9fc7532de3786b82e0d4068357836e8ba7ca23098f8f4b28cc8f2573","id":"ae0a842a-7ed5-44b6-98ae-eed553695818","updated":"2024-02-21T16:23:34.978231Z"}]},"_links":{"self":{"href":"/v1/files?filter=file_name+eq+%27BIG-IP-Next-CentralManager-20.1.0-0.8.115-Update.tgz%27"}}}
	fileId := respInfo["_embedded"].(map[string]interface{})["files"].([]interface{})[0].(map[string]interface{})["id"].(string)
	//delete the file
	fileUrl = fmt.Sprintf("%s/%s", uriCMFileUpload, fileId)
	f5osLogger.Info("[DeleteFile]", "fileUrl", fileUrl)
	respData, err = p.DeleteCMRequest(fileUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[DeleteFile]", "Data::", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

// cm next backup restore resource structs
type BackupRestoreTenantRequest struct {
	FileName string `json:"file_name,omitempty"`
	Password string `json:"encryption_password,omitempty"`
}

type TenantBackupFile struct {
	FileDate        string `json:"file_date,omitempty"`
	FileName        string `json:"file_name,omitempty"`
	FileSize        int    `json:"file_size,omitempty"`
	ExternalStorage bool   `json:"in_external_storage,omitempty"`
	InstanceId      string `json:"instance_id,omitempty"`
	InstanceName    string `json:"instance_name,omitempty"`
}

type TenantBackupRestoreTaskStatus struct {
	Links struct {
		Self struct {
			Link string `json:"href,omitempty"`
		} `json:"self,omitempty"`
	} `json:"_links,omitempty"`
	CompletionDate  string `json:"completed,omitempty"`
	CreationDate    string `json:"created,omitempty"`
	FailureReason   string `json:"failure_reason,omitempty"`
	FilePath        string `json:"file_path,omitempty"`
	FileName        string `json:"file_name,omitempty"`
	Id              string `json:"id,omitempty"`
	ExternalStorage bool   `json:"in_external_storage,omitempty"`
	InstanceId      string `json:"instance_id,omitempty"`
	InstanceName    string `json:"instance_name,omitempty"`
	RunId           string `json:"run_id,omitempty"`
	State           string `json:"state,omitempty"`
	Status          string `json:"status,omitempty"`
}

// /device/v1/inventory?filter=address+eq+'10.10.10.10'

func (p *BigipNextCM) GetDeviceIdByIp(deviceIp string) (deviceId *string, err error) {
	deviceUrl := fmt.Sprintf("%s?filter=address+eq+'%s'", uriInventory, deviceIp)
	f5osLogger.Debug("[GetDeviceInventory]", "URI Path", deviceUrl)
	bigipNextDevice := &DeviceInventoryList{}
	respData, err := p.GetCMRequest(deviceUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[GetDeviceIdByIp]", "Requested BIG-IP Next:", hclog.Fmt("%+v", string(respData)))
	json.Unmarshal(respData, bigipNextDevice)
	if bigipNextDevice.Count == 1 {
		deviceId := bigipNextDevice.Embedded.Devices[0].Id
		return &deviceId, nil
	}
	return nil, fmt.Errorf("the requested device:%s, was not found", deviceIp)
}

func (p *BigipNextCM) GetDeviceInfoByIp(deviceIp string) (deviceInfo interface{}, err error) {
	deviceUrl := fmt.Sprintf("%s?filter=address+eq+'%s'", uriInventory, deviceIp)
	f5osLogger.Debug("[GetDeviceInfoByIp]", "URI Path", deviceUrl)
	respData, err := p.GetCMRequest(deviceUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[GetDeviceInfoByIp]", "Requested BIG-IP Next:", hclog.Fmt("%+v", string(respData)))
	deviceList := make(map[string]interface{})
	json.Unmarshal(respData, &deviceList)
	if len(deviceList["_embedded"].(map[string]interface{})["devices"].([]interface{})) == 1 {
		deviceInfo := deviceList["_embedded"].(map[string]interface{})["devices"].([]interface{})[0]
		return deviceInfo, nil
	}
	return nil, fmt.Errorf("the requested device:%s, was not found", deviceIp)
}

func (p *BigipNextCM) GetDeviceIdByHostname(deviceHostname string) (deviceId *string, err error) {
	deviceUrl := fmt.Sprintf("%s?filter=hostname+eq+'%s'", uriInventory, deviceHostname)
	f5osLogger.Info("[GetDeviceIdByHostname]", "URI Path", deviceUrl)
	bigipNextDevice := &DeviceInventoryList{}
	respData, err := p.GetCMRequest(deviceUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetDeviceIdByHostname]", "Resp BIG-IP Next:", hclog.Fmt("%+v", string(respData)))
	err = json.Unmarshal(respData, &bigipNextDevice)
	if err != nil {
		return nil, err
	}

	if bigipNextDevice.Count == 1 {
		deviceId := bigipNextDevice.Embedded.Devices[0].Id
		return &deviceId, nil
	}
	return nil, fmt.Errorf("the requested device:%s, was not found", deviceHostname)
}

func (p *BigipNextCM) GetDeviceInfoByID(deviceId string) (interface{}, error) {
	// deviceUrl := fmt.Sprintf("%s/%s", uriInventory, deviceId)
	deviceUrl := fmt.Sprintf("%s/%s", uriDiscoverInstance, deviceId)
	url := fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, deviceUrl)
	f5osLogger.Info("[GetDeviceInfoByID]", "Request path", hclog.Fmt("%+v", url))
	dataResource, err := p.doCMRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetDeviceInfoByID]", "Data::", hclog.Fmt("%+v", string(dataResource)))
	var deviceInfo interface{}
	err = json.Unmarshal(dataResource, &deviceInfo)
	if err != nil {
		return nil, err
	}
	return deviceInfo, nil
}

// delete device from CM
func (p *BigipNextCM) DeleteDevice(deviceId string) error {
	// deviceUrl := fmt.Sprintf("%s/%s", uriInventory, deviceId)
	deviceUrl := fmt.Sprintf("%s/%s", uriDiscoverInstance, deviceId)
	url := fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, deviceUrl)
	f5osLogger.Info("[DeleteDevice]", "Request path", hclog.Fmt("%+v", url))
	//{"save_backup":false}
	var data = []byte(`{"save_backup":false}`)
	respData, err := p.doCMRequest("DELETE", url, data)
	// respData, err := p.DeleteCMRequest("DELETE", deviceUrl, data)
	if err != nil {
		return err
	}
	// {"_links":{"self":{"href":"/v1/deletion-tasks/02752890-5660-450c-ace9-b8e0a86a15ad"}},"path":"/v1/deletion-tasks/02752890-5660-450c-ace9-b8e0a86a15ad"}
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteDevice]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	//get task id from path
	pathList := strings.Split(respString["path"].(string), "/")
	taskId := pathList[len(pathList)-1]
	f5osLogger.Info("[DeleteDevice]", "Task Id", hclog.Fmt("%+v", taskId))
	err = p.deleteTaskStatus(taskId)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteDevice]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// https://10.10.10.10/api/device/v1/deletion-tasks/02752890-5660-450c-ace9-b8e0a86a15ad
// verify device deletion task status
func (p *BigipNextCM) deleteTaskStatus(taskID string) error {
	deviceUrl := fmt.Sprintf("%s/%s", "/device/v1/deletion-tasks", taskID)
	url := fmt.Sprintf("%s%s%s", p.Host, uriCMRoot, deviceUrl)
	f5osLogger.Info("[deleteTaskStatus]", "Request path", hclog.Fmt("%+v", url))
	timeout := 360 * time.Second
	endtime := time.Now().Add(timeout)
	respString := make(map[string]interface{})
	for time.Now().Before(endtime) {
		respData, err := p.doCMRequest("GET", url, nil)
		if err != nil {
			return err
		}
		f5osLogger.Debug("[deleteTaskStatus]", "Data::", hclog.Fmt("%+v", string(respData)))
		// {"_links":{"self":{"href":"/v1/deletion-tasks/642d5964-8cd9-4881-9086-1ed5ca682101"}},"address":"10.146.168.20","created":"2023-11-28T07:55:50.924918Z","device_id":"8d6c8c85-1738-4a34-b57b-d3644a2ecfcc","id":"642d5964-8cd9-4881-9086-1ed5ca682101","state":"factoryResetInstance","status":"running"}
		err = json.Unmarshal(respData, &respString)
		if err != nil {
			return err
		}
		f5osLogger.Info("[deleteTaskStatus]", "Task Status", hclog.Fmt("%+v", respString["status"].(string)))
		if respString["status"].(string) == "completed" {
			return nil
		}
		if respString["status"].(string) == "failed" {
			return fmt.Errorf("%s", respString)
		}
		time.Sleep(10 * time.Second)
	}
	return fmt.Errorf("%s", respString)
}

func (p *BigipNextCM) backupTenantTaskStatus(taskidPath string, timeOut int) (*TenantBackupRestoreTaskStatus, error) {
	taskData := &TenantBackupRestoreTaskStatus{}
	backupUrl := fmt.Sprintf("%s%s", "/device", taskidPath)
	f5osLogger.Info("[backupTenantTaskStatus]", "URI Path", backupUrl)
	timeout := time.Duration(timeOut) * time.Second
	endtime := time.Now().Add(timeout)
	for time.Now().Before(endtime) {
		respData, err := p.GetCMRequest(backupUrl)
		if err != nil {
			return nil, err
		}
		f5osLogger.Info("[backupTenantTaskStatus]", "Task Status:\t", hclog.Fmt("%+v", string(respData)))
		err = json.Unmarshal(respData, &taskData)
		if err != nil {
			return nil, err
		}
		if taskData.Status == "completed" && taskData.State == "backupDone" {
			return taskData, nil
		}
		if taskData.Status == "failed" {
			return nil, fmt.Errorf("%s", taskData.FailureReason)
		}
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("task Status is still in :%+v within timeout period of:%+v", taskData.Status, timeout)
}

func (p *BigipNextCM) restoreTenantTaskStatus(taskidPath string, timeOut int) (*TenantBackupRestoreTaskStatus, error) {
	taskData := &TenantBackupRestoreTaskStatus{}
	restoreUrl := fmt.Sprintf("%s%s", "/device", taskidPath)
	f5osLogger.Info("[restoreTenantTaskStatus]", "URI Path", restoreUrl)
	// var timeout time.Duration
	timeout := time.Duration(timeOut) * time.Second
	endtime := time.Now().Add(timeout)
	for time.Now().Before(endtime) {
		respData, err := p.GetCMRequest(restoreUrl)
		if err != nil {
			return nil, err
		}
		f5osLogger.Info("[backupTenantTaskStatus]", "Task Status:\t", hclog.Fmt("%+v", string(respData)))
		err = json.Unmarshal(respData, &taskData)
		if err != nil {
			return nil, err
		}
		if taskData.Status == "completed" && taskData.State == "restoreDone" {
			return taskData, nil
		}
		if taskData.Status == "failed" {
			return nil, fmt.Errorf("%s", taskData.FailureReason)
		}
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("task Status is still in :%+v within timeout period of:%+v", taskData.Status, timeout)
}

func (p *BigipNextCM) BackupTenant(tenantId *string, config *BackupRestoreTenantRequest, timeOut int) (*TenantBackupRestoreTaskStatus, error) {
	backupUrl := fmt.Sprintf("%s/%s/backup", uriInventory, *tenantId)
	f5osLogger.Info("[BackupTenant]", "URI Path", backupUrl)
	f5osLogger.Info("[BackupTenant]", "Backup", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.PostCMRequest(backupUrl, body)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[BackupTenant]", "Data::", hclog.Fmt("%+v", string(respData)))
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[BackupTenant]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	taskData, err := p.backupTenantTaskStatus(respString["path"].(string), timeOut)
	if err != nil {
		return nil, err
	}
	return taskData, nil
}

func (p *BigipNextCM) RestoreTenant(tenantId *string, config *BackupRestoreTenantRequest, timeOut int) (*TenantBackupRestoreTaskStatus, error) {
	restoreUrl := fmt.Sprintf("%s/%s/restore", uriInventory, *tenantId)
	f5osLogger.Info("[RestoreTenant]", "URI Path", restoreUrl)
	f5osLogger.Info("[RestoreTenant]", "Restore", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.PostCMRequest(restoreUrl, body)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[RestoreTenant]", "Data::", hclog.Fmt("%+v", string(respData)))
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[RestoreTenant]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	taskData, err := p.restoreTenantTaskStatus(respString["path"].(string), timeOut)
	if err != nil {
		return nil, err
	}
	return taskData, nil
}

func (p *BigipNextCM) GetBackupFile(fileName *string) (fileData *TenantBackupFile, err error) {
	fileUrl := fmt.Sprintf("%s/%s", uriBackups, *fileName)
	f5osLogger.Debug("[GetBackupFile]", "URI Path", fileUrl)
	respData, err := p.GetCMRequest(fileUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[GetBackupFile]", "Requested Backup File:", hclog.Fmt("%+v", string(respData)))
	err = json.Unmarshal(respData, &fileData)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func (p *BigipNextCM) DeleteBackupFile(fileName *string) error {
	fileUrl := fmt.Sprintf("%s%s%s/%s", p.Host, uriCMRoot, uriBackups, *fileName)
	f5osLogger.Debug("[DeleteBackupFile]", "URI Path", fileUrl)
	respData, err := p.doCMRequest("DELETE", fileUrl, nil)
	if err != nil {
		return err
	}
	f5osLogger.Info("[DeleteBackupFile]", "Data::", hclog.Fmt("%+v", string(respData)))
	return nil
}

// https://10.192.75.131/api/device/v1/providers
// https://10.145.75.237

func (p *BigipNextCM) GetDeviceProviders() ([]byte, error) {
	uriProviders := "/device/v1/providers"
	// providerUrl := fmt.Sprintf("%s", uriProviders)
	f5osLogger.Debug("[GetDeviceProviders]", "URI Path", uriProviders)
	respData, err := p.GetCMRequest(uriProviders)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[GetDeviceProviders]", "Requested Providers:", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

type DeviceProviderReq struct {
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	Connection struct {
		Host           string `json:"host,omitempty"`
		Authentication struct {
			Type     string `json:"type,omitempty"`
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		} `json:"authentication,omitempty"`
		CertFingerprint string `json:"cert_fingerprint,omitempty"`
	} `json:"connection,omitempty"`
}
type DeviceProviderResponse struct {
	Connection struct {
		Authentication struct {
			Type     string `json:"type,omitempty"`
			Username string `json:"username,omitempty"`
		} `json:"authentication,omitempty"`
		Host string `json:"host,omitempty"`
	} `json:"connection,omitempty"`
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// https://10.145.75.237/api/device/v1/providers/vsphere
// Create POST call to create device provider
func (p *BigipNextCM) PostDeviceProvider(config interface{}) (*DeviceProviderResponse, error) {
	uriProviders := "/device/v1/providers/vsphere"
	if config.(*DeviceProviderReq).Type == "VELOS" || config.(*DeviceProviderReq).Type == "RSERIES" {
		uriProviders = "/device/v1/providers/f5os"
		//uriProviders = "/v1/spaces/default/providers/f5os"
	}
	providerUrl := fmt.Sprintf("%s", uriProviders)
	f5osLogger.Debug("[PostDeviceProvider]", "URI Path", providerUrl)
	f5osLogger.Debug("[PostDeviceProvider]", "Config", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.PostCMRequest(providerUrl, body)
	if err != nil {
		if strings.Contains(err.Error(), "cert_fingerprint") {
			config.(*DeviceProviderReq).Connection.CertFingerprint = strings.ReplaceAll(strings.Split(string(strings.Split(err.Error(), ",")[1]), ":")[2], "\"", "")
			return p.PostDeviceProvider(config)
		}
		return nil, err
	}
	f5osLogger.Debug("[PostDeviceProvider]", "Resp::", hclog.Fmt("%+v", string(respData)))
	var providerResp DeviceProviderResponse
	err = json.Unmarshal(respData, &providerResp)
	if err != nil {
		return nil, err
	}
	return &providerResp, nil
}

func (p *BigipNextCM) UpdateDeviceProvider(providerId string, config interface{}) (*DeviceProviderResponse, error) {
	uriProviders := "/device/v1/providers/vsphere"
	if config.(*DeviceProviderReq).Type == "VELOS" || config.(*DeviceProviderReq).Type == "RSERIES" {
		uriProviders = "/device/v1/providers/f5os"
		//uriProviders = "/v1/spaces/default/providers/f5os"
	}
	providerUrl := fmt.Sprintf("%s/%s", uriProviders, providerId)
	f5osLogger.Debug("[UpdateDeviceProvider]", "URI Path", providerUrl)
	f5osLogger.Debug("[UpdateDeviceProvider]", "Config", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.PutCMRequest(providerUrl, body)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[UpdateDeviceProvider]", "Resp::", hclog.Fmt("%+v", string(respData)))
	var providerResp DeviceProviderResponse
	err = json.Unmarshal(respData, &providerResp)
	if err != nil {
		return nil, err
	}
	return &providerResp, nil
}

func (p *BigipNextCM) GetDeviceProvider(providerId, providerType string) (*DeviceProviderResponse, error) {
	uriProviders := "/device/v1/providers/vsphere"
	if stringToUppercase(providerType) == "VELOS" || stringToUppercase(providerType) == "RSERIES" {
		uriProviders = "/device/v1/providers/f5os"
		//uriProviders = "/v1/spaces/default/providers/f5os"
	}
	providerUrl := fmt.Sprintf("%s/%s", uriProviders, providerId)
	f5osLogger.Debug("[GetDeviceProvider]", "URI Path", providerUrl)
	respData, err := p.GetCMRequest(providerUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetDeviceProvider]", "Data::", hclog.Fmt("%+v", string(respData)))
	var providerResp DeviceProviderResponse
	err = json.Unmarshal(respData, &providerResp)
	if err != nil {
		return nil, err
	}
	return &providerResp, nil
}

// https://10.145.75.237/api/device/v1/providers/vsphere/85bc71c3-0bfc-4b28-bb86-13f7e1c1d7af
// Create function to delete device provider using provider id
func (p *BigipNextCM) DeleteDeviceProvider(providerId, providerType string) ([]byte, error) {
	uriProviders := "/device/v1/providers/vsphere"
	if stringToUppercase(providerType) == "VELOS" || stringToUppercase(providerType) == "RSERIES" {
		uriProviders = "/device/v1/providers/f5os"
		//uriProviders = "/v1/spaces/default/providers/f5os"
	}
	providerUrl := fmt.Sprintf("%s/%s", uriProviders, providerId)
	f5osLogger.Debug("[DeleteDeviceProvider]", "URI Path", providerUrl)
	respData, err := p.DeleteCMRequest(providerUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[DeleteDeviceProvider]", "Data::", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

func stringToUppercase(str string) string {
	return strings.ToUpper(str)
}

// {
//     "template_name": "default-standalone-ve",
//     "parameters": {
//         "instantiation_provider": [
//             {
//                 "id": "aec1af9c-cf91-4d2d-9056-c4ea204bb307",
//                 "name": "myvsphere",
//                 "type": "vsphere"
//             }
//         ],
//         "vSphere_properties": [
//             {
//                 "num_cpus": 8,
//                 "memory": 16384,
//                 "datacenter_name": "mbip-7.0",
//                 "cluster_name": "vSAN Cluster",
//                 "datastore_name": "vsanDatastore",
//                 "resource_pool_name": "Earthlings",
//                 "vsphere_content_library": "CM-IOD",
//                 "vm_template_name": "BIG-IP-Next-20.0.1-2.139.10-0.0.136-VM-template"
//             }
//         ],
//         "vsphere_network_adapter_settings": [
//             {
//                 "internal_network_name": "LocalTestVLAN-286",
//                 "ha_data_plane_network_name": "",
//                 "ha_control_plane_network_name": "",
//                 "mgmt_network_name": "VM-mgmt",
//                 "external_network_name": "LocalTestVLAN-196"
//             }
//         ],
//         "dns_servers": [
//             "128.95.112.1"
//         ],
//         "ntp_servers": [],
//         "management_address": "10.146.194.142",
//         "management_network_width": 23,
//         "default_gateway": "10.146.195.254",
//         "l1Networks": [
//             {
//                 "vlans": [
//                     {
//                         "selfIps": [
//                             {
//                                 "address": "10.10.10.10/24",
//                                 "deviceName": "device1"
//                             }
//                         ],
//                         "name": "externalvlan",
//                         "tag": 100
//                     }
//                 ],
//                 "l1Link": {
//                     "linkType": "Interface",
//                     "name": "1.1"
//                 },
//                 "name": "LocalTestVLAN-196"
//             },
//             {
//                 "vlans": [
//                     {
//                         "selfIps": [
//                             {
//                                 "address": "10.10.20.10/24",
//                                 "deviceName": "device2"
//                             }
//                         ],
//                         "name": "internalvlan",
//                         "tag": 110
//                     }
//                 ],
//                 "l1Link": {
//                     "linkType": "Interface",
//                     "name": "1.2"
//                 },
//                 "name": "LocalTestVLAN-286"
//             }
//         ],
//         "management_credentials_username": "admin-cm",
//         "management_credentials_password": "F5Twist@123",
//         "instance_one_time_password": ":#YgMK&wEKhv",
//         "hostname": "testecosyshydvm05"
//     }
// }

// function to get CMReqDeviceInstance struct

// func (p *BigipNextCM) GetCMReqDeviceInstance() (*CMReqDeviceInstance, error) {
// 	var cmReqDeviceInstance CMReqDeviceInstance
// 	cmReqDeviceInstance.TemplateName = "default-standalone-ve"
// 	cmReqDeviceInstance.Parameters.InstantiationProvider = append(cmReqDeviceInstance.Parameters.InstantiationProvider, CMReqInstantiationProvider{
// 		Id:   p.ProviderId,
// 		Name: p.ProviderName,
// 		Type: p.ProviderType,
// 	})
// 	cmReqDeviceInstance.Parameters.VSphereProperties = append(cmReqDeviceInstance.Parameters.VSphereProperties, CMReqVsphereProperties{
// 		NumCpus:               8,
// 		Memory:                16384,
// 		DatacenterName:        p.DatacenterName,
// 		ClusterName:           p.ClusterName,
// 		DatastoreName:         p.DatastoreName,
// 		ResourcePoolName:      p.ResourcePoolName,
// 		VsphereContentLibrary: p.VsphereContentLibrary,
// 		VmTemplateName:        p.VmTemplateName,
// 	})
// 	cmReqDeviceInstance.Parameters.VsphereNetworkAdapterSettings = append(cmReqDeviceInstance.Parameters.VsphereNetworkAdapterSettings, CMReqVsphereNetworkAdapterSettings{
// 		InternalNetworkName:       p.InternalNetworkName,
// 		HaDataPlaneNetworkName:    p.HaDataPlaneNetworkName,
// 		HaControlPlaneNetworkName: p.HaControlPlaneNetworkName,
// 		MgmtNetworkName:           p.MgmtNetworkName,
// 		ExternalNetworkName:       p.ExternalNetworkName,
// 	})
// 	cmReqDeviceInstance.Parameters.DnsServers = append(cmReqDeviceInstance.Parameters.DnsServers, p.DnsServers...)
// 	cmReqDeviceInstance.Parameters.NtpServers = append(cmReqDeviceInstance.Parameters.NtpServers, p.NtpServers...)
// 	cmReqDeviceInstance.Parameters.ManagementAddress = p.ManagementAddress
// 	cmReqDeviceInstance.Parameters.ManagementNetworkWidth = p.ManagementNetworkWidth
// 	cmReqDeviceInstance.Parameters.DefaultGateway = p.DefaultGateway
// 	cmReqDeviceInstance.Parameters.ManagementCredentialsUsername = p.ManagementCredentialsUsername
// 	cmReqDeviceInstance.Parameters.ManagementCredentialsPassword = p.ManagementCredentialsPassword
// 	cmReqDeviceInstance.Parameters.InstanceOneTimePassword = p.InstanceOneTimePassword
// 	cmReqDeviceInstance.Parameters.Hostname = p.Hostname
// 	for _, l1Network := range p.L1Networks {
// 		var cmReqL1Networks CMReqL1Networks
// 		cmReqL1Networks.Name = l1Network.Name
// 		for _, vlan := range l1Network.Vlans {
// 			var cmReqVlans CMReqVlans
// 			cmReqVlans.Name = vlan.Name
// 			cmReqVlans.Tag = vlan.Tag
// 			for _, selfIp := range vlan.SelfIps {
// 				var cmReqSelfIps CMReqSelfIps
// 				cmReqSelfIps.Address = selfIp.Address
// 				cmReqSelfIps.DeviceName = selfIp.DeviceName
// 				cmReqVlans.SelfIps = append(cmReqVlans.SelfIps, cmReqSelfIps)
// 			}
// 			cmReqL1Networks.Vlans = append(cmReqL1Networks.Vlans, cmReqVlans)
// 		}
// 		var cmReqL1Link CMReqL1Link
// 		cmReqL1Link.LinkType = l1Network.L1Link.LinkType
// 		cmReqL1Link.Name = l1Network.L1Link.Name
// 		cmReqL1Networks.L1Link = cmReqL1Link
// 		cmReqDeviceInstance.Parameters.L1Networks = append(cmReqDeviceInstance.Parameters.L1Networks, cmReqL1Networks)
// 	}
// 	return &cmReqDeviceInstance, nil
// }

type CMReqRseriesProperties struct {
	TenantImageName      string `json:"tenant_image_name"`
	TenantDeploymentFile string `json:"tenant_deployment_file"`
	VlanIds              []int  `json:"vlan_ids"`
	DiskSize             int    `json:"disk_size"`
	CpuCores             int    `json:"cpu_cores"`
	// ManagementAddress      string   `json:"management_address"`
	// ManagementNetworkWidth int      `json:"management_network_width"`
	// L1Networks             []string `json:"l1Networks"`
	// ManagementCredentials  struct {
	// 	Username string `json:"username"`
	// 	Password string `json:"password"`
	// } `json:"management_credentials"`
	// InstanceOneTimePassword string `json:"instance_one_time_password"`
	// Hostname                string `json:"hostname"`
}

type CMReqVelosProperties struct {
	TenantImageName      string `json:"tenant_image_name"`
	TenantDeploymentFile string `json:"tenant_deployment_file"`
	VlanIds              []int  `json:"vlan_ids"`
	SlotIds              []int  `json:"slot_ids"`
	DiskSize             int    `json:"disk_size"`
	CpuCores             int    `json:"cpu_cores"`
}

type CMReqDeviceInstance struct {
	TemplateName string `json:"template_name,omitempty"`
	Parameters   struct {
		InstantiationProvider         []CMReqInstantiationProvider         `json:"instantiation_provider,omitempty"`
		VSphereProperties             []CMReqVsphereProperties             `json:"vSphere_properties,omitempty"`
		VsphereNetworkAdapterSettings []CMReqVsphereNetworkAdapterSettings `json:"vsphere_network_adapter_settings,omitempty"`
		RseriesProperties             []CMReqRseriesProperties             `json:"rseries_properties,omitempty"`
		VelosProperties               []CMReqVelosProperties               `json:"velos_properties,omitempty"`
		DnsServers                    []string                             `json:"dns_servers,omitempty"`
		NtpServers                    []string                             `json:"ntp_servers,omitempty"`
		ManagementAddress             string                               `json:"management_address,omitempty"`
		ManagementNetworkWidth        int                                  `json:"management_network_width,omitempty"`
		DefaultGateway                string                               `json:"default_gateway,omitempty"`
		L1Networks                    []CMReqL1Networks                    `json:"l1Networks,omitempty"`
		ManagementCredentialsUsername string                               `json:"management_credentials_username,omitempty"`
		ManagementCredentialsPassword string                               `json:"management_credentials_password,omitempty"`
		InstanceOneTimePassword       string                               `json:"instance_one_time_password,omitempty"`
		Hostname                      string                               `json:"hostname,omitempty"`
	} `json:"parameters,omitempty"`
}

type CMReqInstantiationProvider struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}
type CMReqVsphereProperties struct {
	NumCpus               int    `json:"num_cpus,omitempty"`
	Memory                int    `json:"memory,omitempty"`
	DatacenterName        string `json:"datacenter_name,omitempty"`
	ClusterName           string `json:"cluster_name,omitempty"`
	DatastoreName         string `json:"datastore_name,omitempty"`
	ResourcePoolName      string `json:"resource_pool_name,omitempty"`
	VsphereContentLibrary string `json:"vsphere_content_library,omitempty"`
	VmTemplateName        string `json:"vm_template_name,omitempty"`
}
type CMReqVsphereNetworkAdapterSettings struct {
	InternalNetworkName       string `json:"internal_network_name,omitempty"`
	HaDataPlaneNetworkName    string `json:"ha_data_plane_network_name,omitempty"`
	HaControlPlaneNetworkName string `json:"ha_control_plane_network_name,omitempty"`
	MgmtNetworkName           string `json:"mgmt_network_name,omitempty"`
	ExternalNetworkName       string `json:"external_network_name,omitempty"`
}

type CMReqSelfIps struct {
	Address    string `json:"address,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
}

type CMReqVlans struct {
	SelfIps    []CMReqSelfIps `json:"selfIps,omitempty"`
	Name       string         `json:"name,omitempty"`
	Tag        int            `json:"tag,omitempty"`
	DefaultVrf bool           `json:"defaultVrf,omitempty"`
}

type CMReqL1Networks struct {
	Vlans  []CMReqVlans `json:"vlans,omitempty"`
	L1Link struct {
		LinkType string `json:"linkType,omitempty"`
		Name     string `json:"name,omitempty"`
	} `json:"l1Link,omitempty"`
	Name string `json:"name,omitempty"`
}

type CMReqDeviceInstanceBackup struct {
	TemplateName string `json:"template_name,omitempty"`
	Parameters   struct {
		InstantiationProvider         []CMReqInstantiationProvider         `json:"instantiation_provider,omitempty"`
		VSphereProperties             []CMReqVsphereProperties             `json:"vSphere_properties,omitempty"`
		VsphereNetworkAdapterSettings []CMReqVsphereNetworkAdapterSettings `json:"vsphere_network_adapter_settings,omitempty"`
		DnsServers                    []string                             `json:"dns_servers,omitempty"`
		NtpServers                    []string                             `json:"ntp_servers,omitempty"`
		ManagementAddress             string                               `json:"management_address,omitempty"`
		ManagementNetworkWidth        int                                  `json:"management_network_width,omitempty"`
		DefaultGateway                string                               `json:"default_gateway,omitempty"`
		L1Networks                    []CMReqL1Networks                    `json:"l1Networks,omitempty"`
		ManagementCredentialsUsername string                               `json:"management_credentials_username,omitempty"`
		ManagementCredentialsPassword string                               `json:"management_credentials_password,omitempty"`
		InstanceOneTimePassword       string                               `json:"instance_one_time_password,omitempty"`
		Hostname                      string                               `json:"hostname,omitempty"`
	} `json:"parameters,omitempty"`
}

func (p *BigipNextCM) PostDeviceInstance(config *CMReqDeviceInstance, timeout int) ([]byte, error) {
	uriInstances := "/device/v1/instances"
	instanceUrl := fmt.Sprintf("%s", uriInstances)
	f5osLogger.Debug("[PostDeviceInstance]", "URI Path", instanceUrl)
	f5osLogger.Debug("[PostDeviceInstance]", "Config", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.PostCMRequest(instanceUrl, body)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[PostDeviceInstance]", "Data::", hclog.Fmt("%+v", string(respData)))
	// {"_links":{"self":{"href":"/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438"}},"path":"/v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438"}
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[PostDeviceInstance]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	// split path string to get task id
	taskId := strings.Split(respString["path"].(string), "/")
	f5osLogger.Info("[PostDeviceInstance]", "Task Id", hclog.Fmt("%+v", taskId[len(taskId)-1]))
	// get task status
	taskData, err := p.GetDeviceInstanceTaskStatus(taskId[len(taskId)-1], timeout)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[PostDeviceInstance]", "Task Status", hclog.Fmt("%+v", taskData))
	return respData, nil
}

// /v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438
// get device instance task status
func (p *BigipNextCM) GetDeviceInstanceTaskStatus(taskID string, timeOut int) (map[string]interface{}, error) {
	// taskData := &DeviceInstanceTaskStatus{}
	taskData := make(map[string]interface{})
	instanceUrl := fmt.Sprintf("%s%s", "/device/v1/instances/tasks/", taskID)
	f5osLogger.Debug("[GetDeviceInstanceTaskStatus]", "URI Path", instanceUrl)
	// var timeout time.Duration
	timeout := time.Duration(timeOut) * time.Second
	endtime := time.Now().Add(timeout)
	for time.Now().Before(endtime) {
		respData, err := p.GetCMRequest(instanceUrl)
		if err != nil {
			return nil, err
		}
		f5osLogger.Info("[GetDeviceInstanceTaskStatus]", "Task Status:\t", hclog.Fmt("%+v", string(respData)))
		err = json.Unmarshal(respData, &taskData)
		if err != nil {
			return nil, err
		}
		if taskData["status"] == "completed" {
			return taskData, nil
		}
		if taskData["status"] == "failed" {
			return nil, fmt.Errorf("%s", taskData["failure_reason"])
		}
		inVal := timeOut / 10
		time.Sleep(time.Duration(inVal) * time.Second)
	}
	return nil, fmt.Errorf("task Status is still in :%+v within timeout period of:%+v", taskData["status"], timeout)
}

// // convert a string to byte array
// func stringToByteArray(str string) []byte {
// 	return []byte(str)
// }

func (p *BigipNextCM) GetDeviceProviderIDByHostname(hostname string) (interface{}, error) {
	uriProviders := "/device/v1/providers"
	providerUrl := fmt.Sprintf("%s?filter=name+eq+'%s'", uriProviders, hostname)
	f5osLogger.Info("[GetDeviceProviderIDByHostname]", "URI Path", providerUrl)
	respData, err := p.GetCMRequest(providerUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetDeviceProviderIDByHostname]", "Requested Providers:", hclog.Fmt("%+v", string(respData)))
	var providerResp []interface{}
	err = json.Unmarshal(respData, &providerResp)
	if err != nil {
		return nil, err
	}
	if len(providerResp) == 1 && providerResp[0].(map[string]interface{})["provider_name"].(string) == hostname {
		return providerResp[0].(map[string]interface{})["provider_id"].(string), nil
	}
	return nil, fmt.Errorf("failed to get ID for provider: %+v", hostname)
}

// https://10.145.75.237/api/llm/license/a2064013-659d-4de0-8c22-773d21414885/status
// get device license status
func (p *BigipNextCM) GetDeviceLicenseStatus(deviceId *string) ([]byte, error) {
	uriLicense := "/llm/license"
	licenseUrl := fmt.Sprintf("%s/%s/status", uriLicense, *deviceId)
	f5osLogger.Debug("[GetDeviceLicenseStatus]", "URI Path", licenseUrl)
	respData, err := p.GetCMRequest(licenseUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[GetDeviceLicenseStatus]", "Requested License Status:", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

func encodeUrl(urlName string) string {
	// Encode the urlName
	urlNameEncoded := url.QueryEscape(urlName)
	return urlNameEncoded
}

// ####################
// POST
// https://10.145.75.237/api/llm/tasks/token/verify
// create post call to token verify

// func (p *BigipNextCM) PostTokenVerify(config *TokenVerify) ([]byte, error) {
// 	uriToken := "/llm/tasks/token/verify"
// 	tokenUrl := fmt.Sprintf("%s", uriToken)
// 	f5osLogger.Debug("[PostTokenVerify]", "URI Path", tokenUrl)
// 	f5osLogger.Debug("[PostTokenVerify]", "Config", hclog.Fmt("%+v", config))
// 	body, err := json.Marshal(config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	respData, err := p.PostCMRequest(tokenUrl, body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	f5osLogger.Debug("[PostTokenVerify]", "Data::", hclog.Fmt("%+v", string(respData)))
// 	return respData, nil
// }

// create TokenVerify struct with below payload
//{"jwt":"eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6InYxIiwiamt1IjoiaHR0cHM6Ly9wcm9kdWN0LmFwaXMuZjUuY29tL2VlL3Y.eyJzdWIiOiJGTkktNzYyNmVmM2QtNTc4ZC00MGU1LWI5YjYtZDI3ZmRkZTRlODBkIiwiaWF0IjoxNjQ3NDQ1NjMxLCJpc3MiOiJGNSBJbmMuIiwiYXVkIjoidXJuOmY1OnRlZW0iLCJqdGkiOiI1N2VmM2ZmMC1hNTQwLTExZWMtYjE3Ny1lYmRiMzNmYmNjZmYiLCJmNV9vcmRlcl90eXBlIjoiZXZhbCIsImY1X29yZGVyX3N1YnR5cGUiOiJ0cmlhbCJ9.cBfrqxn09rTGiSKpIu6PpDZoCKOY2BRtm6Q9xfAf0iv6IdY3YZn3iqSR1Qrl5Wgwx1uEDsNasFELdvynAQ1vDTG0QNFgSR5HKC9rFS_QBXK8G2XZuJr_XLQxKeOztzbYTn1V2aoVBeZXawcQG9YVu_MXdkDG2LL7LhWgXWVyckuF99cW1ndwsbucx2nXW7-fcU_TsDnTryt8nwQi0hiw-0DlYXEVYfHxndg_JNRlNtKL8aAgf5rUACrhQTVag7in_UuGV7jhKIk5SjVR2-lUnKA2w3Oo6WCeJv9DyIULWfkJwasBlyF9hiYiMUiTyaW7MK-Kx0w9IamlYy0KBzepFvUsUfYIsRJUnqjFHn_S1Rcg7cGiJyl4XUtVP0LKB80xfxYN2ThAiW7usNSAchepSbUXHatxyWWZavxTu1B48tQmiwBb6_OFSxw_GP1SOlE5v539uObPsJ7cTA-OiWby3VgaU4SgHuLg_ITlwcSc3FSZnQY4qYcu9k8nbhbx-2UmN6C0lkaW5ha2xe08kJWXdPzQQ3Y1bJb9T5IXWeGdGb_ppqHsV7LzuEVZJUACOVFvqXDnJPXggLnI8G_w1aCWLWpxZeRpY7iqWjVGzD5cD_eAwLYoSEUtSyG83dbRSZeDXFQSkZ6ZuLz94iySLZPR98Mi0rLfwpRlkIXBXZ_ZMDs"}

// OUTPUT: {"isValid":true}
// ####################

// ####################
// POST
// https://10.145.75.237/api/llm/token

//{"jwt":"eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6InYxIiwiamt1IjoiaHR0cHM6LaXMuZjUuY29tL2VlL3YxL2tleXMvandrcyJ9.eyJzdWIiOiJGTkktNzYyNmVmM2QtNTc4ZC00MGU1LWI5YjYtZDI3ZmRkZTRlODBkIiwiaWF0IjoxNjQ3NDQ1NjMxLCJpc3MiOiJGNSBJbmMuIiwiYXVkIjoidXJuOmY1OnRlZW0iLCJqdGkiOiI1N2VmM2ZmMC1hNTQwLTExZWMtYjE3Ny1lYmRiMzNmYmNjZmYiLCJmNV9vcmRlcl90eXBlIjoiZXZhbCIsImY1X29yZGVyX3N1YnR5cGUiOiJ0cmlhbCJ9.cBfrqxn09rTGiSKpIu6PpDZoCKOY2BRtm6Q9xfAf0iv6IdY3YZn3iqSR1Qrl5Wgwx1uEDsNasFELdvynAQ1vDTG0QNFgSR5HKC9rFS_QBXK8G2XZuJr_XLQxKeOztzbYTn1V2aoVBeZXawcQG9YVu_MXdkDG2LL7LhWgXWVyckuF99cW1ndwsbucx2nXW7-fcU_TsDnTryt8nwQi0hiw-0DlYXEVYfHxndg_JNRlNtKL8aAgf5rUACrhQTVag7in_UuGV7jhKIk5SjVR2-lUnKA2w3Oo6WCeJv9DyIULWfkJwasBlyF9hiYiMUiTyaW7MK-Kx0w9IamlYy0KBzepFvUsUfYIsRJUnqjFHn_S1Rcg7cGiJyl4XUtVP0LKB80xfxYN2ThAiW7usNSAchepSbUXHatxyWWZavxTu1B48tQmiwBb6_OFSxw_GP1SOlE5v539uObPsJ7cTA-OiWby3VgaU4SgHuLg_ITlwcSc3FSZnQY4qYcu9k8nbhbx-2UmN6C0lkaW5ha2xe08kJWXdPzQQ3Y1bJb9T5IXWeGdGb_ppqHsV7LzuEVZJUACOVFvqXDnJPXggLnI8G_w1aCWLWpxZeRpY7iqWjVGzD5cD_eAwLYoSEUtSyG83dbRSZeDXFQSkZ6ZuLz94iySLZPR98Mi0rLfwpRlkIXBXZ_ZMDs","nickName":"testtoken_ravi"}

// {
//     "DuplicatesTokenNickName": {
//         "orderType": "",
//         "shortName": "",
//         "tokenId": "00000000-0000-0000-0000-000000000000"
//     },
//     "DuplicatesTokenValue": {
//         "orderType": "",
//         "shortName": "",
//         "tokenId": "00000000-0000-0000-0000-000000000000"
//     },
//     "NewToken": {
//         "entitlement": "{\"compliance\":{\"digitalAssetComplianceStatus\":\"\",\"digitalAssetDaysRemainingInState\":0,\"digitalAssetExpiringSoon\":false,\"digitalAssetOutOfComplianceDate\":\"\",\"entitlementCheckStatus\":\"\",\"entitlementExpiryStatus\":\"\",\"telemetryStatus\":\"\",\"usageExceededStatus\":\"\"},\"documentType\":\"BIG-IP Next License\",\"documentVersion\":\"1\",\"digitalAsset\":{\"digitalAssetId\":\"\",\"digitalAssetName\":\"\",\"digitalAssetVersion\":\"\",\"telemetryId\":\"\"},\"entitlementMetadata\":{\"complianceEnforcements\":null,\"complianceStates\":null,\"enforcementBehavior\":\"\",\"enforcementPeriodDays\":0,\"entitlementModel\":\"\",\"expiringSoonNotificationDays\":0,\"entitlementExpiryDate\":\"0001-01-01T00:00:00Z\",\"gracePeriodDays\":0,\"nonContactPeriodHours\":0,\"nonFunctionalPeriodDays\":0,\"orderSubType\":\"\",\"orderType\":\"\"},\"subscriptionMetadata\":{\"programName\":\"big_ip_next_trial\",\"programTypeDescription\":\"big_ip_next_trial\",\"subscriptionId\":\"NGX-Subscription-1-TRL-076761\",\"subscriptionExpiryDate\":\"\",\"subscriptionNotifyDays\":\"\"},\"RepositoryCertificateMetadata\":{\"sslCertificate\":\"\",\"privateKey\":\"\"},\"entitledFeatures\":[{\"entitledFeatureId\":\"c18ab1e4-7801-4840-8241-e50d161c9e0d\",\"featureFlag\":\"bigip_active_assets\",\"featurePermitted\":1250,\"featureRemain\":0,\"featureUnlimited\":false,\"featureUsed\":0,\"featureValueType\":\"integral\",\"uomCode\":\"\",\"uomTerm\":\"\",\"uomTermStart\":0}]}",
//         "orderSubType": "trial",
//         "orderType": "eval",
//         "shortName": "testtoken_ravi",
//         "tokenId": "caee7bc2-b5d1-41e3-8ec8-6f2d0a669810"
//     }
// }

// ####################

// ####################

// POST
// https://10.145.75.237/api/llm/tasks/license/activate
// [{"digitalAssetId":"a2064013-659d-4de0-8c22-773d21414885","jwtId":"caee7bc2-b5d1-41e3-8ec8-6f2d0a669810"}]
// {
//     "a2064013-659d-4de0-8c22-773d21414885": {
//         "_links": {
//             "self": {
//                 "href": "/license-task/c713f2eb-41c6-4c2f-91aa-f01697b27a5b"
//             }
//         },
//         "accepted": true,
//         "deviceId": "a2064013-659d-4de0-8c22-773d21414885",
//         "reason": "",
//         "taskId": "c713f2eb-41c6-4c2f-91aa-f01697b27a5b"
//     }
// }

// GET
// https://10.145.75.237/api/llm/license-task/c713f2eb-41c6-4c2f-91aa-f01697b27a5b
// {
//     "_links": {
//         "self": {
//             "href": "/license-task/c713f2eb-41c6-4c2f-91aa-f01697b27a5b"
//         }
//     },
//     "created": "2023-10-17T17:12:40.372114858Z",
//     "failureReason": "",
//     "status": "running",
//     "subStatus": "TASK_NOT_STARTED",
//     "taskType": "activation"
// }

// {
//     "_links": {
//         "self": {
//             "href": "/license-task/c713f2eb-41c6-4c2f-91aa-f01697b27a5b"
//         }
//     },
//     "created": "2023-10-17T17:12:40.372114858Z",
//     "failureReason": "",
//     "status": "completed",
//     "subStatus": "ACK_VERIFICATION_COMPLETE",
//     "taskType": "activation"
// }

// {
//     "_links": {
//         "self": {
//             "href": "/license/a2064013-659d-4de0-8c22-773d21414885/status"
//         }
//     },
//     "deviceActions": [],
//     "deviceId": "a2064013-659d-4de0-8c22-773d21414885",
//     "expiryDate": "2023-11-16T17:12:48Z",
//     "licenseStatus": "Active",
//     "licenseSubStatus": "ACK_VERIFICATION_COMPLETE",
//     "licenseToken": {
//         "_links": {
//             "self": {
//                 "href": "/token/caee7bc2-b5d1-41e3-8ec8-6f2d0a669810"
//             }
//         },
//         "tokenId": "caee7bc2-b5d1-41e3-8ec8-6f2d0a669810",
//         "tokenName": "testtoken_ravi"
//     },
//     "subscriptionSubType": "trial",
//     "subscriptionType": "eval"
// }

// POST
// https://10.145.75.237/api/llm/tasks/license/deactivate
// {"digitalAssetIds":["a2064013-659d-4de0-8c22-773d21414885"]}
// {
//     "a2064013-659d-4de0-8c22-773d21414885": {
//         "_links": {
//             "self": {
//                 "href": "/license-task/de9bb0c6-20f7-40a3-b4a1-5734abb2b0bd"
//             }
//         },
//         "accepted": true,
//         "deviceId": "a2064013-659d-4de0-8c22-773d21414885",
//         "reason": "",
//         "taskId": "de9bb0c6-20f7-40a3-b4a1-5734abb2b0bd"
//     }
// }

// GET
// https://10.145.75.237/api/llm/license-task/de9bb0c6-20f7-40a3-b4a1-5734abb2b0bd
// // {
//     "_links": {
//         "self": {
//             "href": "/license-task/de9bb0c6-20f7-40a3-b4a1-5734abb2b0bd"
//         }
//     },
//     "created": "2023-10-17T17:22:11.229226131Z",
//     "failureReason": "",
//     "status": "running",
//     "subStatus": "INITIALISE_TERMINATION_INPROGRESS",
//     "taskType": "deactivate"
// }

// // Completion
// {
//     "_links": {
//         "self": {
//             "href": "/license-task/de9bb0c6-20f7-40a3-b4a1-5734abb2b0bd"
//         }
//     },
//     "created": "2023-10-17T17:22:11.229226131Z",
//     "failureReason": "",
//     "status": "completed",
//     "subStatus": "TERMINATE_ACK_VERIFICATION_COMPLETE",
//     "taskType": "deactivate"
// }

// GET
// https://10.145.75.237/api/llm/license/a2064013-659d-4de0-8c22-773d21414885/status
// {
//     "_links": {
//         "self": {
//             "href": "/license/a2064013-659d-4de0-8c22-773d21414885/status"
//         }
//     },
//     "deviceActions": [],
//     "deviceId": "a2064013-659d-4de0-8c22-773d21414885",
//     "expiryDate": "2023-11-16T17:12:48Z",
//     "licenseStatus": "Deactivated",
//     "licenseSubStatus": "TERMINATE_ACK_VERIFICATION_COMPLETE",
//     "licenseToken": {
//         "_links": {
//             "self": {
//                 "href": "/token/caee7bc2-b5d1-41e3-8ec8-6f2d0a669810"
//             }
//         },
//         "tokenId": "caee7bc2-b5d1-41e3-8ec8-6f2d0a669810",
//         "tokenName": "testtoken_ravi"
//     },
//     "subscriptionSubType": "trial",
//     "subscriptionType": "eval"
// }

// INSTANCE DEPLOYMENT

// {
//     "template_name": "default-standalone-ve",
//     "parameters": {
//         "instantiation_provider": [
//             {
//                 "id": "aec1af9c-cf91-4d2d-9056-c4ea204bb307",
//                 "name": "myvsphere",
//                 "type": "vsphere"
//             }
//         ],
//         "vSphere_properties": [
//             {
//                 "num_cpus": 8,
//                 "memory": 16384,
//                 "datacenter_name": "mbip-7.0",
//                 "cluster_name": "vSAN Cluster",
//                 "datastore_name": "vsanDatastore",
//                 "resource_pool_name": "Earthlings",
//                 "vsphere_content_library": "CM-IOD",
//                 "vm_template_name": "BIG-IP-Next-20.0.1-2.139.10-0.0.136-VM-template"
//             }
//         ],
//         "vsphere_network_adapter_settings": [
//             {
//                 "internal_network_name": "LocalTestVLAN-286",
//                 "ha_data_plane_network_name": "",
//                 "ha_control_plane_network_name": "",
//                 "mgmt_network_name": "VM-mgmt",
//                 "external_network_name": "LocalTestVLAN-196"
//             }
//         ],
//         "dns_servers": [
//             "128.95.112.1"
//         ],
//         "ntp_servers": [],
//         "management_address": "10.146.194.142",
//         "management_network_width": 23,
//         "default_gateway": "10.146.195.254",
//         "l1Networks": [
//             {
//                 "vlans": [
//                     {
//                         "selfIps": [
//                             {
//                                 "address": "10.10.10.10/24",
//                                 "deviceName": "device1"
//                             }
//                         ],
//                         "name": "externalvlan",
//                         "tag": 100
//                     }
//                 ],
//                 "l1Link": {
//                     "linkType": "Interface",
//                     "name": "1.1"
//                 },
//                 "name": "LocalTestVLAN-196"
//             },
//             {
//                 "vlans": [
//                     {
//                         "selfIps": [
//                             {
//                                 "address": "10.10.20.10/24",
//                                 "deviceName": "device2"
//                             }
//                         ],
//                         "name": "internalvlan",
//                         "tag": 110
//                     }
//                 ],
//                 "l1Link": {
//                     "linkType": "Interface",
//                     "name": "1.2"
//                 },
//                 "name": "LocalTestVLAN-286"
//             }
//         ],
//         "management_credentials_username": "admin-cm",
//         "management_credentials_password": "F5Twist@123",
//         "instance_one_time_password": ":#YgMK&wEKhv",
//         "hostname": "testecosyshydvm05"
//     }
// }

// create DeviceInstance struct

// https://10.192.75.131/api/device/v1/instances/tasks/7f584bfd-7838-4efc-8ac3-2ce900df25d4

// {
//     "auto_failback": false,
//     "cluster_management_ip": "10.146.168.20",
//     "cluster_name": "raviecosyshydha",
//     "control_plane_vlan": {
//         "name": "ha-cp-vlan",
//         "tag": 101
//     },
//     "data_plane_vlan": {
//         "name": "ha-dp-vlan",
//         "tag": 100,
//         "networkInterface": "1.3"
//     },
//     "nodes": [
//         {
//             "name": "active-node",
//             "control_plane_address": "10.146.168.21/16",
//             "data_plane_primary_address": "10.3.0.10/16",
//             "data_plane_secondary_address": ""
//         },
//         {
//             "name": "standby-node",
//             "control_plane_address": "10.146.168.22/16",
//             "data_plane_primary_address": "10.3.0.11/16",
//             "data_plane_secondary_address": ""
//         }
//     ],
//     "standby_instance_id": "e674140f-9765-4d8b-9b2b-ee2e35059695",
//     "traffic_vlan": []
// }

// create DeviceHA from above json

type CMReqDeviceHA struct {
	AutoFailback        bool   `json:"auto_failback,omitempty"`
	ClusterManagementIP string `json:"cluster_management_ip,omitempty"`
	ClusterName         string `json:"cluster_name,omitempty"`
	ControlPlaneVlan    struct {
		Name string `json:"name,omitempty"`
		Tag  int    `json:"tag,omitempty"`
	} `json:"control_plane_vlan,omitempty"`
	DataPlaneVlan struct {
		Name             string `json:"name,omitempty"`
		Tag              int    `json:"tag,omitempty"`
		NetworkInterface string `json:"networkInterface,omitempty"`
	} `json:"data_plane_vlan,omitempty"`
	Nodes             []CMReqHANode `json:"nodes,omitempty"`
	StandbyInstanceID string        `json:"standby_instance_id,omitempty"`
	TrafficVlan       []interface{} `json:"traffic_vlan,omitempty"`
}

type CMReqHANode struct {
	Name                      string `json:"name,omitempty"`
	ControlPlaneAddress       string `json:"control_plane_address,omitempty"`
	DataPlanePrimaryAddress   string `json:"data_plane_primary_address,omitempty"`
	DataPlaneSecondaryAddress string `json:"data_plane_secondary_address,omitempty"`
}

// https://10.192.75.131/api/device/v1/inventory/7f584bfd-7838-4efc-8ac3-2ce900df25d4/ha
// create POST call to create HA
func (p *BigipNextCM) PostDeviceHA(activeID string, config *CMReqDeviceHA, timeOut int) (interface{}, error) {
	uriHA := "/device/v1/inventory"
	haUrl := fmt.Sprintf("%s/%s/ha", uriHA, activeID)
	f5osLogger.Info("[PostDeviceHA]", "URI Path", haUrl)
	f5osLogger.Debug("[PostDeviceHA]", "Config", hclog.Fmt("%+v", config))
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.PostCMRequest(haUrl, body)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[PostDeviceHA]", "Data::", hclog.Fmt("%+v", string(respData)))
	// {"_links":{"self":{"href":"/v1/ha-creation-tasks/267acbc5-3242-4812-ba88-cd865f8ed41e"}},"path":"/v1/ha-creation-tasks/267acbc5-3242-4812-ba88-cd865f8ed41e"}
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[PostDeviceHA]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	// split path string to get task id
	taskId := strings.Split(respString["path"].(string), "/")
	f5osLogger.Info("[PostDeviceHA]", "Task Id", hclog.Fmt("%+v", taskId[len(taskId)-1]))
	// get task status
	taskData, err := p.GetDeviceHATaskStatus(taskId[len(taskId)-1], timeOut)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[PostDeviceHA]", "Task Status", hclog.Fmt("%+v", taskData))
	return taskData, nil
}

type HATaskResp struct {
	ActiveInstanceId    string `json:"active_instance_id,omitempty"`
	AutoFailback        bool   `json:"auto_failback,omitempty"`
	ClusterManagementIp string `json:"cluster_management_ip,omitempty"`
	ClusterName         string `json:"cluster_name,omitempty"`
	ControlPlaneVlan    struct {
		Tag  int    `json:"tag,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"control_plane_vlan,omitempty"`
	DataPlaneVlan struct {
		Tag              int    `json:"tag,omitempty"`
		Name             string `json:"name,omitempty"`
		NetworkInterface string `json:"NetworkInterface,omitempty"`
	} `json:"data_plane_vlan,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
	Id            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Nodes         []struct {
		Name                    string `json:"name,omitempty"`
		ControlPlaneAddress     string `json:"control_plane_address,omitempty"`
		DataPlanePrimaryAddress string `json:"data_plane_primary_address,omitempty"`
	} `json:"nodes,omitempty"`
	StandbyInstanceId string      `json:"standby_instance_id,omitempty"`
	State             string      `json:"state,omitempty"`
	Status            string      `json:"status,omitempty"`
	TaskType          string      `json:"task_type,omitempty"`
	TrafficVlan       interface{} `json:"traffic_vlan,omitempty"`
}

// /v1/ha-creation-tasks/267acbc5-3242-4812-ba88-cd865f8ed41e
// get device HA task status
func (p *BigipNextCM) GetDeviceHATaskStatus(taskID string, timeOut int) (map[string]interface{}, error) {
	// taskData := &DeviceInstanceTaskStatus{}
	taskData := make(map[string]interface{})
	instanceUrl := fmt.Sprintf("%s%s", "/device/v1/ha-creation-tasks/", taskID)
	f5osLogger.Debug("[GetDeviceHATaskStatus]", "URI Path", instanceUrl)
	timeout := time.Duration(timeOut) * time.Second
	endtime := time.Now().Add(timeout)
	for time.Now().Before(endtime) {
		respData, err := p.GetCMRequest(instanceUrl)
		if err != nil {
			return nil, err
		}
		f5osLogger.Info("[GetDeviceHATaskStatus]", "Task Status:\t", hclog.Fmt("%+v", string(respData)))
		// {"_links":{"self":{"href":"/v1/ha-creation-tasks/06aea4ed-7425-4db3-a728-2574929885d9"}},"active_instance_id":"8d6c8c85-1738-4a34-b57b-d3644a2ecfcc","auto_failback":false,"cluster_management_ip":"10.146.168.20","cluster_name":"raviecosyshydha","control_plane_vlan":{"tag":101,"name":"ha-cp-vlan"},"created":"2023-11-28T05:56:57.618962Z","data_plane_vlan":{"tag":102,"name":"ha-dp-vlan","NetworkInterface":"1.3"},"id":"06aea4ed-7425-4db3-a728-2574929885d9","name":"create HA from 8d6c8c85-1738-4a34-b57b-d3644a2ecfcc","nodes":[{"name":"active-node","control_plane_address":"10.146.168.21/16","data_plane_primary_address":"10.3.0.10/16"},{"name":"standby-node","control_plane_address":"10.146.168.22/16","data_plane_primary_address":"10.3.0.10/16"}],"standby_instance_id":"d0e9cda1-4460-4132-87fd-0f3aa18f3872","state":"haGetNodesLoginInfo","status":"running","task_type":"instance_ha_creation","traffic_vlan":null,"updated":"2023-11-28T05:56:57.712342Z"}
		err = json.Unmarshal(respData, &taskData)
		if err != nil {
			return nil, err
		}
		if taskData["status"] == "completed" {
			return taskData, nil
		}
		if taskData["status"] == "failed" {
			return nil, fmt.Errorf("%s", taskData["failure_reason"])
		}
		inVal := timeOut / 10
		time.Sleep(time.Duration(inVal) * time.Second)
	}
	return nil, fmt.Errorf("task status is still in :%+v within timeout period of:%+v", taskData["status"], timeout)
}

type CMHANodes struct {
	NodeAddress string `json:"node_address,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

type CMHANodesStatus struct {
	Metadata CMHANodeMetadata   `json:"metadata,omitempty"`
	Spec     CMHANodeSpec       `json:"spec,omitempty"`
	Status   CMHANodeReadyState `json:"status,omitempty"`
}

type CMHANodeMetadata struct {
	Id        string `json:"id,omitempty"`
	MachineId string `json:"machine_id,omitempty"`
	Name      string `json:"name,omitempty"`
}

type CMHANodeSpec struct {
	NodeAddress string `json:"node_address,omitempty"`
	NodeType    string `json:"node_type,omitempty"`
}

type CMHANodeReadyState struct {
	Ready        bool   `json:"ready,omitempty"`
	Registration string `json:"registration,omitempty"`
}

func (p *BigipNextCM) CreateCMHACluster(nodes []CMHANodes) (string, error) {
	uriNodes := "/v1/system/infra/nodes"
	payload, err := json.Marshal(nodes)
	if err != nil {
		return "", err
	}

	resp, err := p.PostCMRequest(uriNodes, payload)

	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func (p *BigipNextCM) GetCMHANodes() ([]CMHANodesStatus, error) {
	uriNodes := "/v1/system/infra/nodes"

	var resp []CMHANodesStatus
	ret, err := p.GetCMRequest(uriNodes)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(ret, &resp)
	return resp, nil
}

func (p *BigipNextCM) CheckCMHANodesStatus(nodes []string) ([]CMHANodesStatus, error) {
	uriNodes := "/v1/system/infra/nodes"

	var resp []CMHANodesStatus
	start := time.Now()
	waitTime := time.Duration(5) * time.Minute

	for time.Since(start) < waitTime {
		ret, err := p.GetCMRequest(uriNodes)

		if err != nil {
			f5osLogger.Error("[CheckCMHANodesStatus]", "Error", err)
			time.Sleep(3 * time.Second)
			continue
		}

		json.Unmarshal(ret, &resp)

		r := p.CheckCMHANodeReadyState(resp, nodes)

		if len(r) == 0 {
			return resp, nil
		} else if len(r) > 0 {
			time.Sleep(3 * time.Second)
			continue
		}
	}

	return nil, fmt.Errorf("timeout waiting for nodes to be ready")
}

func (p *BigipNextCM) CheckCMHANodeReadyState(resp []CMHANodesStatus, nodes []string) []string {
	var waitingNodes []string

	for _, node := range resp {
		for _, n := range nodes {
			if node.Metadata.Name != "central-manager-server" && node.Spec.NodeAddress == n {
				if node.Status.Ready {
					f5osLogger.Info("[CheckCMHANodesStatus]", "Info", fmt.Sprintf("Node %v is ready", node.Spec.NodeAddress))
				} else {
					waitingNodes = append(waitingNodes, node.Spec.NodeAddress)
				}
			}
		}
	}

	if len(waitingNodes) > 0 {
		f5osLogger.Info("[CheckCMHANodesStatus]", "Info", fmt.Sprintf("Waiting for nodes to be ready: %v", waitingNodes))
		return waitingNodes
	}
	f5osLogger.Info("[CheckCMHANodesStatus]", "Info", fmt.Sprintf("All nodes are ready: %v", nodes))
	return waitingNodes
}

func (p *BigipNextCM) DeleteCMHANodes(deleteNodes []string) {
	uriNodes := "/api/v1/system/infra/nodes"

	for i, node := range deleteNodes {
		nodeName := "central-manager-" + strings.ReplaceAll(node, ".", "-")
		nodeName = strings.ReplaceAll(nodeName, "\"", "")
		uri := p.Host + uriNodes + "/" + nodeName
		res, err := p.doCMRequest("DELETE", uri, nil)
		if err != nil {
			f5osLogger.Error("[DeleteCMHANodes]", "Error", fmt.Sprintf("%v, retrying in 10 seconds", err))
			time.Sleep(10 * time.Second)
			res, err = p.doCMRequest("DELETE", uri, nil)
			if err != nil {
				f5osLogger.Error("[DeleteCMHANodes]", "Error", fmt.Sprintf("%v, unable to delete node %v after retry", err, nodeName))
				continue
			}
		}
		f5osLogger.Info("[DeleteCMHANodes]", "Info", string(res))
		if i != len(deleteNodes)-1 {
			time.Sleep(10 * time.Second)
		}
	}
}

// https://10.144.73.240/api/device/v1/instances

// Req:
// {
//     "template_name": "default-standalone-rseries",
//     "parameters": {
//         "instantiation_provider": [
//             {
//                 "id": "1dbe5b14-ea5d-48f9-a3ad-96ca52eba694",
//                 "name": "myrseries",
//                 "type": "rseries"
//             }
//         ],
//         "rseries_properties": [
//             {
//                 "tenant_image_name": "BIG-IP-Next-20.1.0-2.264.6",
//                 "tenant_deployment_file": "BIG-IP-Next-20.1.0-2.264.6.yaml",
//                 "vlan_ids": [
//                     27,
//                     28,
//                     29
//                 ],
//                 "disk_size": 30,
//                 "cpu_cores": 4
//             }
//         ],
//         "default_gateway": "10.144.140.254",
//         "management_address": "10.144.140.81",
//         "management_network_width": 24,
//         "l1Networks": [],
//         "management_credentials_username": "admin-cm",
//         "management_credentials_password": "F5Twist@123",
//         "instance_one_time_password": "8&(IFi/]kAdX",
//         "hostname": "testecosyshydvm01"
//     }
// }

// Response:

// {
//     "_links": {
//         "self": {
//             "href": "/v1/instances/tasks/cf081285-9df9-44e3-9856-7e412610b6d0"
//         }
//     },
//     "path": "/v1/instances/tasks/cf081285-9df9-44e3-9856-7e412610b6d0"
// }

// https://10.144.73.240/api/device/v1/instances/tasks/cf081285-9df9-44e3-9856-7e412610b6d0

// {
//     "_links": {
//         "self": {
//             "href": "/v1/instances/tasks/cf081285-9df9-44e3-9856-7e412610b6d0"
//         }
//     },
//     "created": "2023-12-04T18:44:26.765093Z",
//     "id": "cf081285-9df9-44e3-9856-7e412610b6d0",
//     "name": "instance creation",
//     "payload": {
//         "discovery": {
//             "port": 5443,
//             "address": "10.144.140.81",
//             "device_user": "admin",
//             "device_password": "*****",
//             "management_user": "admin-cm",
//             "management_password": "*****"
//         },
//         "onboarding": {
//             "mode": "STANDALONE",
//             "nodes": [
//                 {
//                     "password": "*****",
//                     "username": "admin",
//                     "managementAddress": "10.144.140.81"
//                 }
//             ],
//             "platformType": "RSERIES"
//         },
//         "instantiation": {
//             "Request": {
//                 "F5osRequest": {
//                     "provider_id": "1dbe5b14-ea5d-48f9-a3ad-96ca52eba694",
//                     "provider_type": "rseries",
//                     "next_instances": [
//                         {
//                             "nodes": [
//                                 1
//                             ],
//                             "vlans": [
//                                 27,
//                                 28,
//                                 29
//                             ],
//                             "mgmt_ip": "10.144.140.81",
//                             "timeout": 360,
//                             "hostname": "testecosyshydvm01",
//                             "cpu_cores": 4,
//                             "disk_size": 30,
//                             "mgmt_prefix": 24,
//                             "mgmt_gateway": "10.144.140.254",
//                             "admin_password": "*****",
//                             "tenant_image_name": "BIG-IP-Next-20.1.0-2.264.6",
//                             "tenant_deployment_file": "BIG-IP-Next-20.1.0-2.264.6.yaml"
//                         }
//                     ]
//                 },
//                 "VsphereRequest": null
//             },
//             "BaseTask": {
//                 "id": "",
//                 "payload": null,
//                 "provider_id": "1dbe5b14-ea5d-48f9-a3ad-96ca52eba694",
//                 "provider_type": "rseries"
//             },
//             "VsphereRequest": null
//         }
//     },
//     "stage": "Instantiation",
//     "state": "instantiateInstances",
//     "status": "running",
//     "task_type": "instance_creation",
//     "updated": "2023-12-04T18:44:26.843488Z"
// }
