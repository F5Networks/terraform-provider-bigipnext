/*
Copyright 2024 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
// Package bigipnext interacts with BIGIP-NEXT/CM systems using the OPEN API.
package bigipnext

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
)

const (
	uriInventory        = "/device/v1/inventory"
	uriProviders        = "/v1/spaces/default/providers"
	uriDiscoverInstance = "/v1/spaces/default/instances"
	uriLicense          = "/v1/spaces/default/instances/license"
	// uriLicenseActivate  = "/v1/spaces/default/instances/license/activate"
)

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
	ResourcePoolId        string `json:"resource_pool_id,omitempty"`
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

func (p *BigipNextCM) PostDeviceProvider(config interface{}) (*DeviceProviderResponse, error) {
	providerUrl := fmt.Sprintf("%s/vsphere", uriProviders)
	if config.(*DeviceProviderReq).Type == "VELOS" || config.(*DeviceProviderReq).Type == "RSERIES" {
		providerUrl = fmt.Sprintf("%s/f5os", uriProviders)
	}
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
	providerUrl := fmt.Sprintf("%s/vsphere", uriProviders)
	if config.(*DeviceProviderReq).Type == "VELOS" || config.(*DeviceProviderReq).Type == "RSERIES" {
		providerUrl = fmt.Sprintf("%s/f5os", uriProviders)
	}
	providerUrl = fmt.Sprintf("%s/%s", providerUrl, providerId)
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
	providerUrl := fmt.Sprintf("%s/vsphere", uriProviders)
	if stringToUppercase(providerType) == "VELOS" || stringToUppercase(providerType) == "RSERIES" {
		providerUrl = fmt.Sprintf("%s/f5os", uriProviders)
	}
	providerUrl = fmt.Sprintf("%s/%s", providerUrl, providerId)
	f5osLogger.Debug("[GetDeviceProvider]", "URI Path", providerUrl)
	respData, err := p.GetCMRequest(providerUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetDeviceProvider]", "\n--------Resp--------\n", hclog.Fmt("%+v", string(respData)))
	var providerResp DeviceProviderResponse
	err = json.Unmarshal(respData, &providerResp)
	if err != nil {
		return nil, err
	}
	return &providerResp, nil
}

func (p *BigipNextCM) DeleteDeviceProvider(providerId, providerType string) ([]byte, error) {
	providerUrl := fmt.Sprintf("%s/vsphere", uriProviders)
	if stringToUppercase(providerType) == "VELOS" || stringToUppercase(providerType) == "RSERIES" {
		providerUrl = fmt.Sprintf("%s/f5os", uriProviders)
	}
	providerUrl = fmt.Sprintf("%s/%s", providerUrl, providerId)
	f5osLogger.Debug("[DeleteDeviceProvider]", "URI Path", providerUrl)
	respData, err := p.DeleteCMRequest(providerUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[DeleteDeviceProvider]", "Data::", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

func (p *BigipNextCM) GetDeviceProviderIDByHostname(hostname string) (interface{}, error) {
	providerUrl := fmt.Sprintf("%s?filter=name+eq+'%s'", uriProviders, hostname)
	f5osLogger.Info("[GetDeviceProviderIDByHostname]", "URI Path", providerUrl)
	respData, err := p.GetCMRequest(providerUrl)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[GetDeviceProviderIDByHostname]", "provider query response:", hclog.Fmt("%+v", string(respData)))
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

func (p *BigipNextCM) GetResourcePoolID(providerID, dataCenterName, clusterName, resourcePoolName string) ([]byte, error) {
	providerUrl := fmt.Sprintf("%s/vsphere/%s/api?path=api/vcenter/datacenter", uriProviders, providerID)
	respData, err := p.GetCMRequest(providerUrl)
	if err != nil {
		return nil, err
	}
	var mapDatacenter []interface{}
	err = json.Unmarshal(respData, &mapDatacenter)
	if err != nil {
		return nil, err
	}
	for _, datacenter := range mapDatacenter {
		if datacenter.(map[string]interface{})["name"].(string) == dataCenterName {
			datacenterID := datacenter.(map[string]interface{})["datacenter_id"].(string)
			providerUrl := fmt.Sprintf("%s/vsphere/%s/api?path=api/vcenter/cluster?datacenters=%s", uriProviders, providerID, datacenterID)
			respData, err := p.GetCMRequest(providerUrl)
			if err != nil {
				return nil, err
			}
			var mapCluster []interface{}
			err = json.Unmarshal(respData, &mapCluster)
			if err != nil {
				return nil, err
			}
			for _, cluster := range mapCluster {
				if cluster.(map[string]interface{})["name"].(string) == clusterName {
					clusterID := cluster.(map[string]interface{})["cluster_id"].(string)
					providerUrl := fmt.Sprintf("%s/vsphere/%s/api?path=api/vcenter/resource-pool?clusters=%s", uriProviders, providerID, clusterID)
					respData, err := p.GetCMRequest(providerUrl)
					if err != nil {
						return nil, err
					}
					var mapResourcePool []interface{}
					err = json.Unmarshal(respData, &mapResourcePool)
					if err != nil {
						return nil, err
					}
					for _, resourcePool := range mapResourcePool {
						if resourcePool.(map[string]interface{})["name"].(string) == resourcePoolName {
							resourcePoolID := resourcePool.(map[string]interface{})["resource_pool_id"].(string)
							return []byte(resourcePoolID), nil
						}
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to get ID for resource pool: %+v", resourcePoolName)

	// [{\"name\":\"vSAN Cluster\",\"cluster_id\":\"domain-c8\"}]
	// providerUrl := fmt.Sprintf("%s/vsphere/%s/api?path=api/vcenter/resource-pool?clusters=domain-c8", uriProviders, providerID)
	// [{"name":"CM-Applications","resource_pool_id":"resgroup-10283"},{"name":"Infrastructure","resource_pool_id":"resgroup-2015"},{"name":"SPK-TEST","resource_pool_id":"resgroup-4008"},{"name":"Mbip-waf","resource_pool_id":"resgroup-4010"},{"name":"mbipauto_spk_team","resource_pool_id":"resgroup-4011"},{"name":"CNF-TEST","resource_pool_id":"resgroup-4014"},{"name":"Piecemakers","resource_pool_id":"resgroup-4016"},{"name":"BIG-IP-Next-License","resource_pool_id":"resgroup-4018"},{"name":"BIG-IP-Next-Platform","resource_pool_id":"resgroup-4025"},{"name":"Eternals-Mavericks","resource_pool_id":"resgroup-4029"},{"name":"LTM-CM-PDM","resource_pool_id":"resgroup-4030"},{"name":"Multiverse","resource_pool_id":"resgroup-4032"},{"name":"Prodigy","resource_pool_id":"resgroup-4033"},{"name":"SSLO","resource_pool_id":"resgroup-4035"},{"name":"Sunrisers","resource_pool_id":"resgroup-4036"},{"name":"Central-Manager","resource_pool_id":"resgroup-4039"},{"name":"CM Solution","resource_pool_id":"resgroup-4040"},{"name":"CM-Access","resource_pool_id":"resgroup-4041"},{"name":"CM-MBIP-Solution","resource_pool_id":"resgroup-4043"},{"name":"EQ","resource_pool_id":"resgroup-4045"},{"name":"INFRAANO","resource_pool_id":"resgroup-4047"},{"name":"MBIPMP_System","resource_pool_id":"resgroup-4049"},{"name":"Simplifiers","resource_pool_id":"resgroup-4053"},{"name":"CM Pipeline Pool","resource_pool_id":"resgroup-4055"},{"name":"Sales","resource_pool_id":"resgroup-4061"},{"name":"Simplifiers","resource_pool_id":"resgroup-4063"},{"name":"test","resource_pool_id":"resgroup-4064"},{"name":"SPK-DEV","resource_pool_id":"resgroup-4065"},{"name":"Journeys","resource_pool_id":"resgroup-4067"},{"name":"samurais","resource_pool_id":"resgroup-4069"},{"name":"Sparatans","resource_pool_id":"resgroup-4070"},{"name":"Synergy","resource_pool_id":"resgroup-4071"},{"name":"TestPoolDEV-DND","resource_pool_id":"resgroup-4073"},{"name":"XC-BIGIP-SaaS","resource_pool_id":"resgroup-4075"},{"name":"LooneyTunes","resource_pool_id":"resgroup-4077"},{"name":"Team-Cosmos","resource_pool_id":"resgroup-4078"},{"name":"Team-Jazz","resource_pool_id":"resgroup-4080"},{"name":"ESXI-1","resource_pool_id":"resgroup-5116"},{"name":"automation-test","resource_pool_id":"resgroup-5304"},{"name":"Automation-Pipeline-Traffic-Machines","resource_pool_id":"resgroup-6542"},{"name":"Koolminds","resource_pool_id":"resgroup-8411"},{"name":"Resources","resource_pool_id":"resgroup-9"}]
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

func (p *BigipNextCM) PostDeviceInstance(config *CMReqDeviceInstance, timeout int) ([]byte, error) {
	instanceUrl := "/device/v1/instances"
	// instanceUrl := fmt.Sprintf("%s", uriInstances)
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

// https://<BIG-IP-Next-Central-Manager-IP-Address>/api/device/api/v1/spaces/default/instances/initialization/{instance_id}

func (p *BigipNextCM) UpdateNextInstanceConfig(instanceID string, config *CMReqDeviceInstance, timeout int) ([]byte, error) {
	instanceUrl := fmt.Sprintf("%s/%s%s/%s", p.Host, uriDefault, "/instances/initialization", instanceID)
	f5osLogger.Debug("[UpdateNextInstanceConfig]", "URI Path", instanceUrl)
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[UpdateNextInstanceConfig]", "Config", hclog.Fmt("%+v", string(body)))
	respData, err := p.doCMRequest("PATCH", instanceUrl, body)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[UpdateNextInstanceConfig]", "Data::", hclog.Fmt("%+v", string(respData)))
	// return respData, nil
	// {"_links":{"self":{"href":"/api/v1/spaces/default/instances/initialization/tasks/9d1a572a-6a28-4af2-a555-c61f77dbbb96"}},"path":"/api/v1/spaces/default/instances/initialization/tasks/9d1a572a-6a28-4af2-a555-c61f77dbbb96"}
	respString := make(map[string]interface{})
	err = json.Unmarshal(respData, &respString)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[UpdateNextInstanceConfig]", "Task Path", hclog.Fmt("%+v", respString["path"].(string)))
	// split path string to get task id
	taskId := strings.Split(respString["path"].(string), "/")
	f5osLogger.Info("[UpdateNextInstanceConfig]", "Task Id", hclog.Fmt("%+v", taskId[len(taskId)-1]))
	// get task status
	taskData, err := p.GetDeviceInstanceTaskStatus(taskId[len(taskId)-1], timeout)
	if err != nil {
		return nil, err
	}
	f5osLogger.Info("[UpdateNextInstanceConfig]", "Task Status", hclog.Fmt("%+v", taskData))
	return respData, nil
}

// /v1/instances/tasks/deacca61-3162-4672-aac8-2d6bd2b69438
// get device instance task status
func (p *BigipNextCM) GetDeviceInstanceTaskStatus(taskID string, timeOut int) (map[string]interface{}, error) {
	taskData := make(map[string]interface{})
	// "/api/v1/spaces/default/instances/initialization/tasks"
	instanceUrl := fmt.Sprintf("%s%s%s%s", p.Host, uriDefault, "/instances/initialization/tasks/", taskID)
	f5osLogger.Debug("[GetDeviceInstanceTaskStatus]", "URI Path", instanceUrl)
	// var timeout time.Duration
	timeout := time.Duration(timeOut) * time.Second
	endtime := time.Now().Add(timeout)
	for time.Now().Before(endtime) {
		respData, err := p.doCMRequest("GET", instanceUrl, nil)
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

type LicenseReq struct {
	DigitalAssetId string `json:"digitalAssetId,omitempty"`
	JwtId          string `json:"jwtId,omitempty"`
}

// https://clouddocs.f5.com/api/v1/spaces/default/instances/license/activate
// Activate License Post Req
func (p *BigipNextCM) PostActivateLicense(config interface{}) (interface{}, error) {
	uriLicenseActivate := fmt.Sprintf("%s%s", uriLicense, "/activate")
	f5osLogger.Debug("[PostActivateLicense]", "URI Path", uriLicenseActivate)
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.PostCMRequest(uriLicenseActivate, body)
	if err != nil {
		return nil, err
	}

	// {
	// 	"422b0cec-03b9-4499-a26e-c88f57869637": {
	// 		"_links": {
	// 			"self": {
	// 				"href": "/license-task/41e49d68-d146-4e16-b286-7b57731fe14d"
	// 			}
	// 		},
	// 		"accepted": true,
	// 		"deviceId": "422b0cec-03b9-4499-a26e-c88f57869637",
	// 		"reason": "",
	// 		"taskId": "41e49d68-d146-4e16-b286-7b57731fe14d"
	// 	},
	// 	"bf89ae4b-c8f1-4c93-b47e-f2051513ad2f": {
	// 		"_links": {
	// 			"self": {
	// 				"href": "/license-task/10786da5-a6a0-45fd-83d3-9db89a8f0a33"
	// 			}
	// 		},
	// 		"accepted": true,
	// 		"deviceId": "bf89ae4b-c8f1-4c93-b47e-f2051513ad2f",
	// 		"reason": "",
	// 		"taskId": "10786da5-a6a0-45fd-83d3-9db89a8f0a33"
	// 	}
	// }
	// get taskid

	respMap := make(map[string]interface{})
	err = json.Unmarshal(respData, &respMap)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[PostActivateLicense]", "Task Path", hclog.Fmt("%+v", respMap))
	// get task id from respMap for each device
	var taskIds []string
	for _, v := range respMap {
		f5osLogger.Info("[PostActivateLicense]", "Task Id", hclog.Fmt("%+v", v.(map[string]interface{})["taskId"].(string)))
		taskIds = append(taskIds, v.(map[string]interface{})["taskId"].(string))
	}
	f5osLogger.Debug("[PostActivateLicense]", "taskIds:", hclog.Fmt("%+v", taskIds))
	lictskReq := &LicenseTaskReq{}
	lictskReq.LicenseTaskIds = taskIds
	return p.PostLicenseTaskStatus(lictskReq)
	// return taskIds, nil
}

// {
// 	"licenseTaskIds": [
// 	  "d290f1ee-6c54-4b01-90e6-d701748f0851",
// 	  "d290f1ee-6c54-4b01-90e6-d701748f0852",
// 	  "d290f1ee-6c54-4b01-90e6-d701748f0853"
// 	]
//   }

type LicenseTaskReq struct {
	LicenseTaskIds []string `json:"licenseTaskIds,omitempty"`
}

// https://clouddocs.f5.com/api/v1/spaces/default/license/tasks
// Create POST call to get license task status
func (p *BigipNextCM) PostLicenseTaskStatus(config interface{}) (interface{}, error) {
	uriLicenseTasks := "/v1/spaces/default/license/tasks"
	f5osLogger.Debug("[PostLicenseTaskStatus]", "URI Path", uriLicenseTasks)
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.PostCMRequest(uriLicenseTasks, body)
	if err != nil {
		return nil, err
	}

	// {
	// 	"3e45e2bd-4e01-4926-8794-67bf8ceb4f61": {
	// 		"_links": {
	// 			"self": {
	// 				"href": "/license-task/3e45e2bd-4e01-4926-8794-67bf8ceb4f61"
	// 			}
	// 		},
	// 		"taskExecutionStatus": {
	// 			"created": "2024-06-27T15:59:25.928845Z",
	// 			"failureReason": "",
	// 			"status": "completed",
	// 			"subStatus": "TERMINATE_ACK_VERIFICATION_COMPLETE",
	// 			"taskType": "deactivate"
	// 		}
	// 	},
	// 	"dfeb50ae-c664-4b76-ac29-540f5e5178ab": {
	// 		"_links": {
	// 			"self": {
	// 				"href": "/license-task/dfeb50ae-c664-4b76-ac29-540f5e5178ab"
	// 			}
	// 		},
	// 		"taskExecutionStatus": {
	// 			"created": "2024-06-27T15:59:25.914384Z",
	// 			"failureReason": "",
	// 			"status": "completed",
	// 			"subStatus": "TERMINATE_ACK_VERIFICATION_COMPLETE",
	// 			"taskType": "deactivate"
	// 		}
	// 	}
	// }

	respMap := make(map[string]interface{})
	err = json.Unmarshal(respData, &respMap)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[PostLicenseTaskStatus]", "Task Path", hclog.Fmt("%+v", respMap))
	// verify taskExecutionStatus
	count := 0
	for k, v := range respMap {
		f5osLogger.Info("[PostLicenseTaskStatus]", "Task Id", hclog.Fmt("%+v", k))
		if v.(map[string]interface{})["taskExecutionStatus"].(map[string]interface{})["status"].(string) == "completed" {
			count++
		} else if v.(map[string]interface{})["taskExecutionStatus"].(map[string]interface{})["status"].(string) == "failed" {
			return nil, fmt.Errorf("%s", v.(map[string]interface{})["taskExecutionStatus"].(map[string]interface{})["failureReason"].(string))
		} else {
			time.Sleep(30 * time.Second)
			return p.PostLicenseTaskStatus(config)
		}
	}
	if count == len(respMap) {
		return respMap, nil
	}
	return respMap, nil
}

// https://clouddocs.f5.com/api/v1/spaces/default/instances/license/license-info
// create POST call to get license info
func (p *BigipNextCM) PostLicenseInfo(config interface{}) (interface{}, error) {
	uriLicenseInfo := fmt.Sprintf("%s%s", uriLicense, "/license-info")
	// uriLicenseInfo := "/v1/spaces/default/instances/license/license-info"
	f5osLogger.Debug("[PostLicenseInfo]", "URI Path", uriLicenseInfo)
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.PostCMRequest(uriLicenseInfo, body)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[PostLicenseInfo]", "Data::", hclog.Fmt("%+v", string(respData)))
	//conver to map
	respMap := make(map[string]interface{})
	err = json.Unmarshal(respData, &respMap)
	if err != nil {
		return nil, err
	}
	return respMap, nil
}

type LicenseDeactivaeReq struct {
	DigitalAssetIds []string `json:"digitalAssetIds,omitempty"`
}

// https://clouddocs.f5.com/api/v1/spaces/default/instances/license/deactivate
func (p *BigipNextCM) PostDeactivateLicense(config interface{}) (interface{}, error) {
	uriLicenseDeactivate := fmt.Sprintf("%s%s", uriLicense, "/deactivate")
	// uriLicenseDeactivate := "/v1/spaces/default/instances/license/deactivate"
	f5osLogger.Debug("[PostDeactivateLicense]", "URI Path", uriLicenseDeactivate)
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	respData, err := p.PostCMRequest(uriLicenseDeactivate, body)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(respData, &respMap)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[PostDeactivateLicense]", "Task Path", hclog.Fmt("%+v", respMap))
	// get task id from respMap for each device
	var taskIds []string
	for _, v := range respMap {
		f5osLogger.Info("[PostDeactivateLicense]", "Task Id", hclog.Fmt("%+v", v.(map[string]interface{})["taskId"].(string)))
		taskIds = append(taskIds, v.(map[string]interface{})["taskId"].(string))
	}
	lictskReq := &LicenseTaskReq{}
	lictskReq.LicenseTaskIds = taskIds
	return p.PostLicenseTaskStatus(lictskReq)
	// f5osLogger.Debug("[PostDeactivateLicense]", "taskIds:", hclog.Fmt("%+v", taskIds))
	// return taskIds, nil
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

// https://clouddocs.f5.com/api/v1/spaces/default/instances/onboarding-tasks
// create Get call to get onboarding task status
func (p *BigipNextCM) GetOnboardingTasks() ([]byte, error) {
	uriOnboardingTasks := fmt.Sprintf("%s/%s", uriDiscoverInstance, "onboarding-tasks")
	f5osLogger.Debug("[GetOnboardingTasks]", "URI Path", uriOnboardingTasks)
	respData, err := p.GetCMRequest(uriOnboardingTasks)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[GetOnboardingTasks]", "Requested Onboarding Tasks:", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

// https://<BIG-IP-Next-Central-Manager-IP-Address>/api/device/api/v1/spaces/default/instances/initialization/{instance_id}
// create Get call to get onboarding task status
func (p *BigipNextCM) GetInitializationTaskStatus(instanceId string) ([]byte, error) {
	uriInitializationTask := fmt.Sprintf("%s/%s/%s/%s", p.Host, uriDefault, "instances/initialization", instanceId)
	f5osLogger.Debug("[GetInitializationTaskStatus]", "URI Path", uriInitializationTask)
	respData, err := p.doCMRequest("GET", uriInitializationTask, nil)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[GetInitializationTaskStatus]", "Requested Initialization Task:", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}

// https://clouddocs.f5.com/api/v1/spaces/default/instances/{instanceID}/interfaces
func (p *BigipNextCM) GetDeviceInterfaces(instanceId string) ([]byte, error) {
	uriInterfaces := fmt.Sprintf("%s/%s/%s", uriDiscoverInstance, instanceId, "interfaces")
	f5osLogger.Debug("[GetDeviceInterfaces]", "URI Path", uriInterfaces)
	respData, err := p.GetCMRequest(uriInterfaces)
	if err != nil {
		return nil, err
	}
	f5osLogger.Debug("[GetDeviceInterfaces]", "Requested Interfaces:", hclog.Fmt("%+v", string(respData)))
	return respData, nil
}
