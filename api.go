package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type KongClient struct {
	kongAdminUri string
	httpClient   *http.Client
}

type KongAPIList struct {
	Data  []KongAPI
	Total int
}

type KongAPI struct {
	Id            string `json:"id,omitempty"`
	Name          string `json:"name"`
	Created_at    int64  `json:"-"`
	Upstream_url  string `json:"upstream_url,omitempty"`
	Preserve_host bool   `json:"preserve_host"`
	// kong 0.9.x and earlier
	Request_path       string `json:"request_path,omitempty"`
	Request_host       string `json:"request_host,omitempty"`
	Strip_request_path bool   `json:"strip_request_path,omitempty"`
	// kong 0.10.x and later
	Hosts     []string `json:"hosts,omitempty"`
	Uris      []string `json:"uris,omitempty"`
	Strip_uri bool     `json:"strip_uri,omitempty"`
}

const (
	ERR_REQUEST_FAILED   = "Kong Admin API request filed: %v"
	ERR_MARSHAL_FAILED   = "Can't marshal object to JSON: %v"
	ERR_UNMARSHAL_FAILED = "Can't unmarshal JSON returned from Kong Admin API: %v"
)

func NewKongClient(httpClient *http.Client, kongAdminUri string) *KongClient {
	kongClient := KongClient{
		httpClient:   httpClient,
		kongAdminUri: strings.Trim(kongAdminUri, "/ "),
	}
	return &kongClient
}

func (k *KongClient) GetVersion() (string, error) {
	data := make(map[string]interface{})
	body, err := k.getUri(k.kongAdminUri)
	if err != nil {
		return "", fmt.Errorf(ERR_REQUEST_FAILED, err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf(ERR_REQUEST_FAILED, err)
	}
	version := data["version"].(string)
	return version, nil
}

func (k *KongClient) GetAPIs() (KongAPIList, error) {

	apis := KongAPIList{}

	body, err := k.getUri(k.getApisUri())
	if err != nil {
		return apis, fmt.Errorf(ERR_REQUEST_FAILED, err)
	}

	err = json.Unmarshal(body, &apis)
	if err != nil {
		return apis, fmt.Errorf(ERR_UNMARSHAL_FAILED, err)
	}

	return apis, nil
}

func (k *KongClient) GetAPI(apiName string) (KongAPI, error) {

	api := KongAPI{}

	body, err := k.getUri(k.getApiUri(apiName))
	if err != nil {
		return api, fmt.Errorf(ERR_REQUEST_FAILED, err)
	}

	return k.jsonToKongApi(body)

}

func (k *KongClient) AddAPI(api KongAPI) (KongAPI, error) {

	json, err := k.kongApiToJson(api)

	body, err := k.postUri(k.getApisUri(), 201, json)

	if err != nil {
		return KongAPI{}, fmt.Errorf(ERR_REQUEST_FAILED, err)
	}

	return k.jsonToKongApi(body)
}

func (k *KongClient) DeleteAPI(apiName string) error {

	err := k.deleteUri(k.getApiUri(apiName), 204)
	if err != nil {
		return fmt.Errorf(ERR_REQUEST_FAILED, err)
	}
	return nil
}

func (k *KongClient) getUri(uri string) ([]byte, error) {

	b := make([]byte, 0)

	resp, err := http.Get(uri)

	if err != nil {
		return b, err
	}

	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return b, err
	}

	return b, nil
}

func (k *KongClient) postUri(uri string, expectedResponseCode int, data []byte) ([]byte, error) {

	resp, err := k.httpClient.Post(uri, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != expectedResponseCode {
		return nil, fmt.Errorf("%v", resp.Status+" "+string(body[:len(body)]))
	}

	return body, nil
}

func (k *KongClient) deleteUri(uri string, expectedResponseCode int) error {

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}

	resp, err := k.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if resp.StatusCode != expectedResponseCode {
		return fmt.Errorf("%v", resp.Status+" "+string(body[:len(body)]))
	}

	return nil
}

func (k *KongClient) kongApiToJson(api KongAPI) ([]byte, error) {
	body, err := json.Marshal(api)
	if err != nil {
		return nil, fmt.Errorf(ERR_MARSHAL_FAILED, err)
	}
	return body, nil
}

func (k *KongClient) jsonToKongApi(body []byte) (KongAPI, error) {
	api := KongAPI{}

	err := json.Unmarshal(body, &api)
	if err != nil {
		return api, fmt.Errorf(ERR_UNMARSHAL_FAILED, err)
	}

	return api, nil
}

func (k *KongClient) getApisUri() string {
	return k.kongAdminUri + "/apis"
}

func (k *KongClient) getApiUri(apiName string) string {
	return k.getApisUri() + "/" + apiName
}
