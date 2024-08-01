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
