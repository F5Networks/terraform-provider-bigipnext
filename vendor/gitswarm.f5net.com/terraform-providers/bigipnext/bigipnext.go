/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
// Package bigipnext interacts with BIGIP-NEXT systems using the OPEN API.
package bigipnext

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	uriRoot           = "/api/v1/"
	uriLogin          = "/api/v1/login"
	contentTypeHeader = "application/json"
	uriPlatformType   = "/openconfig-platform:components/component=platform/state/description"
	uriVlan           = "/openconfig-vlan:vlans"
	uriInterface      = "/openconfig-interfaces:interfaces"
	uriAs3Post        = "/mgmt/shared/appsvcs/declare"
	uriAs3            = "/mgmt/shared/appsvcs"
)

var f5osLogger hclog.Logger

var defaultConfigOptions = &ConfigOptions{
	APICallTimeout: 60 * time.Second,
}

type ConfigOptions struct {
	APICallTimeout time.Duration
}

type BigipNextConfig struct {
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

// BigipNext is a container for our session state.
type BigipNext struct {
	Host          string
	Token         string // if set, will be used instead of User/Password
	Transport     *http.Transport
	UserAgent     string
	Teem          bool
	ConfigOptions *ConfigOptions
	PlatformType  string
}

type NextLoginResp struct {
	Token            string `json:"token"`
	TokenType        string `json:"tokenType"`
	ExpiresIn        int    `json:"expiresIn"`
	RefreshToken     string `json:"refreshToken"`
	RefreshExpiresIn int    `json:"refreshExpiresIn"`
}
type BigipNextError struct {
	IetfRestconfErrors struct {
		Error []struct {
			ErrorType    string `json:"error-type"`
			ErrorTag     string `json:"error-tag"`
			ErrorPath    string `json:"error-path"`
			ErrorMessage string `json:"error-message"`
		} `json:"error"`
	} `json:"ietf-restconf:errors"`
}

// Upload contains information about a file upload status
type Upload struct {
	RemainingByteCount int64          `json:"remainingByteCount"`
	UsedChunks         map[string]int `json:"usedChunks"`
	TotalByteCount     int64          `json:"totalByteCount"`
	LocalFilePath      string         `json:"localFilePath"`
	TemporaryFilePath  string         `json:"temporaryFilePath"`
	Generation         int            `json:"generation"`
	LastUpdateMicros   int            `json:"lastUpdateMicros"`
}

// RequestError contains information about any error we get from a request.
type RequestError struct {
	Code       int      `json:"code,omitempty"`
	Message    string   `json:"message,omitempty"`
	ErrorStack []string `json:"errorStack,omitempty"`
}

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
	Status   string `json:"status,omitempty"`
}

type Vlan struct {
	Name               string   `json:"name,omitempty"`
	Tag                int      `json:"tag,omitempty"`
	TaggedInterfaces   []string `json:"taggedInterfaces,omitempty"`
	UntaggedInterfaces []string `json:"untaggedInterfaces,omitempty"`
	Mtu                int      `json:"mtu,omitempty"`
	SelfIps            []struct {
		Address    string `json:"address,omitempty"`
		DeviceName string `json:"deviceName,omitempty"`
	} `json:"selfIps,omitempty"`
}

type OnboardNextConfig struct {
	Token             string `json:"token"`
	DeviceCertificate struct {
		Key struct {
			Name        string `json:"name,omitempty"`
			FileName    string `json:"fileName,omitempty"`
			FileType    string `json:"fileType,omitempty"`
			FileContent string `json:"fileContent,omitempty"`
		} `json:"key,omitempty"`
		Certificate struct {
			Name        string `json:"name,omitempty"`
			FileName    string `json:"fileName,omitempty"`
			FileType    string `json:"fileType,omitempty"`
			FileContent string `json:"fileContent,omitempty"`
		} `json:"certificate,omitempty"`
	} `json:"deviceCertificate,omitempty"`
	Users struct {
		Users []User `json:"users,omitempty"`
	} `json:"users,omitempty"`
	Cluster struct {
		Name                string `json:"name,omitempty"`
		ClusterManagementIP string `json:"clusterManagementIP,omitempty"`
		VirtualRouterId     int    `json:"virtualRouterId,omitempty"`
		AutoFailback        bool   `json:"autoFailback,omitempty"`
		Nodes               []struct {
			Name                      string `json:"name,omitempty"`
			ManagementAddress         string `json:"managementAddress,omitempty"`
			ControlPlaneAddress       string `json:"controlPlaneAddress,omitempty"`
			Username                  string `json:"username,omitempty"`
			Password                  string `json:"password,omitempty"`
			Token                     string `json:"token,omitempty"`
			DataPlanePrimaryAddress   string `json:"dataPlanePrimaryAddress,omitempty"`
			DataPlaneSecondaryAddress string `json:"dataPlaneSecondaryAddress,omitempty"`
		} `json:"nodes,omitempty"`
		DataPlaneVlan    string `json:"dataPlaneVlan,omitempty"`
		ControlPlaneVlan string `json:"controlPlaneVlan,omitempty"`
	} `json:"cluster,omitempty"`
	PlatformType        string   `json:"platformType,omitempty"`
	DefaultGateway      string   `json:"defaultGateway,omitempty"`
	DnsServers          string   `json:"dnsServers,omitempty"`
	Hostname            string   `json:"hostname,omitempty"`
	ManagementIps       string   `json:"managementIps,omitempty"`
	NtpServers          []string `json:"ntpServers,omitempty"`
	Vlans               []Vlan   `json:"vlans,omitempty"`
	RemoteSyslogServers []struct {
		Name       string `json:"name,omitempty"`
		Host       string `json:"host,omitempty"`
		RemotePort int    `json:"remotePort,omitempty"`
	} `json:"remoteSyslogServers,omitempty"`
}

type SystemsObj struct {
	Embedded struct {
		Systems []struct {
			Links struct {
				Cores             string `json:"cores"`
				DataplaneDebug    string `json:"dataplaneDebug"`
				DeviceCertificate string `json:"deviceCertificate"`
				HostDNS           string `json:"hostDns"`
				HostNetwork       string `json:"hostNetwork"`
				Interfaces        string `json:"interfaces"`
				Logs              string `json:"logs"`
				Manifest          string `json:"manifest"`
				Self              string `json:"self"`
			} `json:"_links"`
			Hostname      string   `json:"hostname"`
			ID            string   `json:"id"`
			MachineID     string   `json:"machineID"`
			ManagementIps []string `json:"managementIps"`
			Name          string   `json:"name"`
			PlatformType  string   `json:"platformType"`
		} `json:"systems"`
	} `json:"_embedded"`
	Count int `json:"count"`
	Total int `json:"total"`
}

// type OnboardNextConfig struct {
// 	Token             string `json:"token"`
// 	DeviceCertificate struct {
// 		Key struct {
// 			Name        string `json:"name,omitempty"`
// 			FileName    string `json:"fileName,omitempty"`
// 			FileType    string `json:"fileType,omitempty"`
// 			FileContent string `json:"fileContent,omitempty"`
// 		} `json:"key,omitempty"`
// 		Certificate struct {
// 			Name        string `json:"name,omitempty"`
// 			FileName    string `json:"fileName,omitempty"`
// 			FileType    string `json:"fileType,omitempty"`
// 			FileContent string `json:"fileContent,omitempty"`
// 		} `json:"certificate,omitempty"`
// 	} `json:"deviceCertificate,omitempty"`
// 	Users struct {
// 		Users []User `json:"users,omitempty"`
// 	} `json:"users,omitempty"`
// 	Cluster struct {
// 		Name                string `json:"name,omitempty"`
// 		ClusterManagementIP string `json:"clusterManagementIP,omitempty"`
// 		VirtualRouterId     int    `json:"virtualRouterId,omitempty"`
// 		AutoFailback        bool   `json:"autoFailback,omitempty"`
// 		Nodes               []struct {
// 			Name                      string `json:"name,omitempty"`
// 			ManagementAddress         string `json:"managementAddress,omitempty"`
// 			ControlPlaneAddress       string `json:"controlPlaneAddress,omitempty"`
// 			Username                  string `json:"username,omitempty"`
// 			Password                  string `json:"password,omitempty"`
// 			Token                     string `json:"token,omitempty"`
// 			DataPlanePrimaryAddress   string `json:"dataPlanePrimaryAddress,omitempty"`
// 			DataPlaneSecondaryAddress string `json:"dataPlaneSecondaryAddress,omitempty"`
// 		} `json:"nodes,omitempty"`
// 		DataPlaneVlan    string `json:"dataPlaneVlan,omitempty"`
// 		ControlPlaneVlan string `json:"controlPlaneVlan,omitempty"`
// 	} `json:"cluster,omitempty"`
// 	Platform struct {
// 		PlatformType        string   `json:"platformType,omitempty"`
// 		DefaultGateway      string   `json:"defaultGateway,omitempty"`
// 		DnsServers          string   `json:"dnsServers,omitempty"`
// 		Hostname            string   `json:"hostname,omitempty"`
// 		ManagementIps       string   `json:"managementIps,omitempty"`
// 		NtpServers          []string `json:"ntpServers,omitempty"`
// 		Vlans               []Vlan   `json:"vlans,omitempty"`
// 		RemoteSyslogServers []struct {
// 			Name       string `json:"name,omitempty"`
// 			Host       string `json:"host,omitempty"`
// 			RemotePort int    `json:"remotePort,omitempty"`
// 		} `json:"remoteSyslogServers,omitempty"`
// 	} `json:"platform,omitempty"`
// }

// Error returns the error message.
func (r *BigipNextError) Error() error {
	if len(r.IetfRestconfErrors.Error) > 0 {
		return errors.New(r.IetfRestconfErrors.Error[0].ErrorMessage)
	}
	return nil
}

func init() {
	val, ok := os.LookupEnv("TF_LOG")
	if !ok {
		val, ok = os.LookupEnv("TF_LOG_PROVIDER_BIGIPNEXT")
		if !ok {
			val = "INFO"
		}
	}
	f5osLogger = hclog.New(&hclog.LoggerOptions{
		Name:  "[BIGIPNEXT]",
		Level: hclog.LevelFromString(val),
	})
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

// NewSession sets up connection to the BIG-IP Next system.

func NewSession(bigipNextobj *BigipNextConfig) (*BigipNext, error) {
	f5osLogger.Info("[NewSession] Session creation Starts...")
	var urlString string
	bigipNextSession := &BigipNext{}
	if !strings.HasPrefix(bigipNextobj.Host, "http") {
		urlString = fmt.Sprintf("https://%s", bigipNextobj.Host)
	} else {
		urlString = bigipNextobj.Host
	}
	u, _ := url.Parse(urlString)
	_, port, _ := net.SplitHostPort(u.Host)

	if bigipNextobj.Port != 0 && port == "" {
		urlString = fmt.Sprintf("%s:%d", urlString, bigipNextobj.Port)
	}
	if bigipNextobj.ConfigOptions == nil {
		bigipNextobj.ConfigOptions = defaultConfigOptions
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	bigipNextSession.Host = urlString
	bigipNextSession.Transport = tr
	bigipNextSession.ConfigOptions = bigipNextobj.ConfigOptions
	client := &http.Client{
		Transport: tr,
	}
	method := "GET"
	urlString = fmt.Sprintf("%s%s", urlString, uriLogin)

	f5osLogger.Info("[NewSession]", "URL", hclog.Fmt("%+v", urlString))
	req, err := http.NewRequest(method, urlString, nil)
	req.Header.Set("Content-Type", contentTypeHeader)
	req.SetBasicAuth(bigipNextobj.User, bigipNextobj.Password)
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
	var resp NextLoginResp
	json.Unmarshal(bodyResp, &resp)
	bigipNextSession.Token = resp.Token
	f5osLogger.Info("[NewSession] Session creation Success")
	return bigipNextSession, nil
}

func (p *BigipNext) doRequest(op, path string, body []byte) ([]byte, error) {
	f5osLogger.Debug("[doRequest]", "Request path", hclog.Fmt("%+v", path))
	if len(body) > 0 {
		f5osLogger.Debug("[doRequest]", "Request body", hclog.Fmt("%+v", string(body)))
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
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	f5osLogger.Debug("[doRequest]", "Resp CODE", hclog.Fmt("%+v", resp.StatusCode))
	if resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 202 {
		return io.ReadAll(resp.Body)
	}
	if resp.StatusCode == 404 {
		return io.ReadAll(resp.Body)
	}
	if resp.StatusCode >= 400 {
		byteData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s", byteData)
	}
	return nil, nil
}

func (p *BigipNext) GetRequest(path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s%s", p.Host, uriRoot, path)
	f5osLogger.Info("[GetRequest]", "Request path", hclog.Fmt("%+v", url))
	return p.doRequest("GET", url, nil)
}

func (p *BigipNext) DeleteRequest(path string) error {
	url := fmt.Sprintf("%s%s%s", p.Host, uriRoot, path)
	f5osLogger.Debug("[DeleteRequest]", "Request path", hclog.Fmt("%+v", url))
	if resp, err := p.doRequest("DELETE", url, nil); err != nil {
		return err
	} else if len(resp) > 0 {
		f5osLogger.Trace("[DeleteRequest]", "Response", hclog.Fmt("%+v", string(resp)))
	}
	return nil
}

func (p *BigipNext) DeleteRequestNew(path string, data []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s%s", p.Host, uriRoot, path)
	f5osLogger.Info("[DeleteRequestNew]", "Request path", hclog.Fmt("%+v", url))
	resp, err := p.doRequest("DELETE", url, data)
	if err != nil {
		return []byte(""), err
	}
	f5osLogger.Info("[DeleteRequestNew]", "Resp Data", hclog.Fmt("%+v", string(resp)))
	return resp, nil
}

func (p *BigipNext) PatchRequest(path string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s%s", p.Host, uriRoot, path)
	f5osLogger.Debug("[PatchRequest]", "Request path", hclog.Fmt("%+v", url))
	return p.doRequest("PATCH", url, body)
}

func (p *BigipNext) PostRequest(path string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s", p.Host, path)
	f5osLogger.Debug("[PostRequest]", "Request path", hclog.Fmt("%+v", url))
	return p.doRequest("POST", url, body)
}

func (p *BigipNext) PutRequest(path string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s%s", p.Host, uriRoot, path)
	f5osLogger.Info("[PutRequest]", "Request path", hclog.Fmt("%+v", url))
	return p.doRequest("PUT", url, body)
}

func (p *BigipNext) PostAs3(as3Json string, tenantFilter string) (string, string, error) {
	var as3Endpoint string
	asyncTrue := "?async=true"
	if tenantFilter != "" {
		as3Endpoint = fmt.Sprintf("%s/%s%s", uriAs3Post, tenantFilter, asyncTrue)
	} else {
		as3Endpoint = fmt.Sprintf("%s%s", uriAs3Post, asyncTrue)
	}
	// tenantPath := tenantFilter + "?async=true"
	// as3Endpoint := uriAs3Post + tenantPath
	successfulTenants := make([]string, 0)

	resp, err := p.PostRequest(as3Endpoint, []byte(as3Json))
	if err != nil {
		return "", "", err
	}
	jsonResp := make(map[string]any)
	json.Unmarshal(resp, &jsonResp)
	f5osLogger.Info("Response", fmt.Sprintf("--------------------%s-------------------", jsonResp))
	taskId := jsonResp["id"].(string)
	taskStatus, err := p.getAs3TaskStatus(taskId)
	if err != nil {
		return "", "", err
	}
	respCode := taskStatus["results"].([]any)[0].(map[string]any)["code"].(float64)

	for respCode != 200 {
		time.Sleep(3 * time.Second)
		fastTask, err := p.getAs3TaskStatus(taskId)
		if err != nil {
			return "", taskId, err
		}
		respCode = fastTask["results"].([]interface{})[0].(map[string]interface{})["code"].(float64)
		if respCode == 200 {
			log.Printf("[DEBUG]Sucessfully Created Application with ID  = %v", taskId)
			break
		}
		if respCode >= 400 {
			j, _ := json.MarshalIndent(fastTask["results"].([]interface{}), "", "\t")
			return "", taskId, fmt.Errorf("tenant Creation failed. Response: %+v", string(j))
		}
		if respCode != 0 && respCode != 503 {
			tenant_list, tenant_count, _ := p.GetTenantList(as3Json)
			if tenantCompare(tenant_list, tenantFilter) == 1 {
				taskMsg := fastTask["results"].([]any)[0].(map[string]any)["message"].(string)
				if len(fastTask["results"].([]any)) == 1 && taskMsg == "declaration is invalid" {
					return "", taskId, fmt.Errorf("Error :%+v", fastTask["results"].([]interface{})[0].(map[string]interface{})["errors"])
				}
				if len(fastTask["results"].([]any)) == 1 && taskMsg == "no change" {
					return "", taskId, fmt.Errorf("Error:%+v", taskMsg)
				}
				i := tenant_count - 1
				success_count := 0
				for i >= 0 {
					if fastTask["results"].([]interface{})[i].(map[string]interface{})["code"].(float64) == 200 {
						successfulTenants = append(successfulTenants, fastTask["results"].([]interface{})[i].(map[string]interface{})["tenant"].(string))
						success_count++
					}
					if fastTask["results"].([]interface{})[i].(map[string]interface{})["code"].(float64) >= 400 {
						log.Printf("[ERROR] : HTTP %v :: %s for tenant %v", fastTask["results"].([]interface{})[i].(map[string]interface{})["code"].(float64), taskMsg, fastTask["results"].([]interface{})[i].(map[string]interface{})["tenant"])
					}
					i = i - 1
				}
				if success_count == tenant_count {
					log.Printf("[DEBUG]Sucessfully Created Application with ID  = %v", taskId)
					break // break here
				} else if success_count == 0 {
					j, _ := json.MarshalIndent(fastTask["results"].([]interface{}), "", "\t")
					return "", taskId, fmt.Errorf("tenant Creation failed. Response: %+v", string(j))
				} else {
					finallist := strings.Join(successfulTenants[:], ",")
					j, _ := json.MarshalIndent(fastTask["results"].([]interface{}), "", "\t")
					return finallist, taskId, fmt.Errorf("as3 config post error response %+v", string(j))
				}
			}
		}
	}
	return strings.Join(successfulTenants[:], ","), taskId, nil
}

func (p *BigipNext) DeleteAs3All() error {
	err := p.DeleteRequest(uriAs3Post)
	if err != nil {
		return err
	}
	return nil
}

func (p *BigipNext) DeleteAs3(tenantName string) (string, error) {
	// tenant := tenantName + "?async=true"
	// as3Endpoint := uriAs3Post + "/" + tenantName
	if tenantName == "" {
		return "", fmt.Errorf("name of the tenant to be deleted is not provided")
	}
	failedTenants := make([]string, 0)
	err := p.DeleteRequest(uriAs3Post + "/" + tenantName)

	if err != nil {
		return "", err
	}

	tenantList := strings.Split(tenantName, ",")

	resp, err := p.GetRequest(uriAs3Post)
	if err != nil {
		return "", err
	}
	existingTenants, _, _ := p.GetTenantList(string(resp))
	nonDeletedTenants := strings.Split(existingTenants, ",")

	for _, v1 := range tenantList {
		for _, v2 := range nonDeletedTenants {
			if v1 == v2 {
				failedTenants = append(failedTenants, v1)
			}
		}
	}

	return strings.Join(failedTenants, ","), nil
}

func (p *BigipNext) GetAs3(tenantList string) (string, error) {
	as3Json := make(map[string]interface{})
	as3Json["class"] = "AS3"
	as3Json["action"] = "deploy"
	as3Json["persist"] = true
	adcJson := make(map[string]interface{})

	as3GetUrl := "/mgmt/shared/appsvcs/declare/" + tenantList
	resp, err := p.GetRequest(as3GetUrl)

	if err != nil {
		return string(resp), err
	}
	delete(adcJson, "updateMode")
	delete(adcJson, "controls")
	as3Json["declaration"] = adcJson
	out, _ := json.Marshal(as3Json)
	as3String := string(out)

	return as3String, nil
}

func tenantCompare(t1 string, t2 string) int {
	tenantList1 := strings.Split(t1, ",")
	tenantList2 := strings.Split(t2, ",")
	if len(tenantList1) == len(tenantList2) {
		return 1
	}
	return 0
}

func (p *BigipNext) getAs3TaskStatus(id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/mgmt/shared/appsvcs/task/%s", id)
	var taskList map[string]interface{}
	resp, err := p.GetRequest(path)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(resp, &taskList)
	return taskList, nil
}

func (p *BigipNext) TenantDifference(slice1 []string, slice2 []string) string {
	var diff []string
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, s1)
		}
	}
	diff_tenant_list := strings.Join(diff[:], ",")
	return diff_tenant_list
}

func (p *BigipNext) AddUser(userConfig *User) ([]byte, error) {
	byteBody, err := json.Marshal(userConfig)
	if err != nil {
		return byteBody, err
	}
	respData, err := p.PutRequest("users", byteBody)
	if err != nil {
		return respData, err
	}
	f5osLogger.Info("[AddUser]", "userConfig:", hclog.Fmt("%+v", string(respData)))
	return respData, nil

}

func (p *BigipNext) OnboardNext(userConfig *OnboardNextConfig) ([]byte, error) {
	// byteBody, err := json.Marshal(userConfig)
	byteBody, err := jsonMarshal(userConfig)
	if err != nil {
		return byteBody, err
	}
	f5osLogger.Info("[OnboardNext]", "userConfig:", hclog.Fmt("%+v", string(byteBody)))
	respData, err := p.PutRequest("onboard", byteBody)
	if err != nil {
		return respData, err
	}
	f5osLogger.Info("[OnboardNext]", "userConfig:", hclog.Fmt("%+v", string(respData)))
	return respData, nil

}

func (p *BigipNext) GetSystems() (*SystemsObj, error) {
	f5osLogger.Info("[GetSystems]", "get system info")
	bigipNextSystems := &SystemsObj{}
	respData, err := p.GetRequest("systems")
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetSystems]", "userConfig:", hclog.Fmt("%+v", string(respData)))
	json.Unmarshal(respData, bigipNextSystems)
	return bigipNextSystems, nil
}

func (p *BigipNext) GetTenantList(body interface{}) (string, int, string) {
	tenantList := make([]string, 0)
	applicationList := make([]string, 0)
	as3json := body.(string)
	resp := []byte(as3json)
	jsonRef := make(map[string]interface{})
	json.Unmarshal(resp, &jsonRef)
	for key, value := range jsonRef {
		if rec, ok := value.(map[string]interface{}); ok && key == "declaration" {
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
	finalApplicationList := strings.Join(applicationList[:], ",")
	return finalTenantlist, len(tenantList), finalApplicationList
}

// contains checks if a int is present in
// a slice
func contains(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// jsonMarshal specifies an encoder with 'SetEscapeHTML' set to 'false' so that <, >, and & are not escaped. https://golang.org/pkg/encoding/json/#Marshal
// https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
