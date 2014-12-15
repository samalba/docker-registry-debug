package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Registry struct {
	Url          *url.URL
	Host         string
	RegistryHost string
	Logger       func(format string, args ...interface{})
	client       *http.Client
}

type LayerJson struct {
	Id            string
	Parent        string
	DockerVersion string `json:"docker_version"`
	Created       string
	Author        string
	Size          int
}

func NewRegistry(endpoint string, registryDomain string) (*Registry, error) {
	if endpoint == "" && registryDomain == "" {
		return nil, fmt.Errorf("needs a registry endpoint when index endpoint is missing")
	}
	parseEndpoint := func(endpoint string) (*url.URL, error) {
		if endpoint != "" {
			u, e := url.Parse(endpoint)
			if e != nil {
				return nil, e
			}
			if u.Host == "" {
				u.Host = u.Path
			}
			return u, nil
		}
		return nil, nil
	}
	host := ""
	origUrl, e := parseEndpoint(endpoint)
	if e != nil {
		return nil, e
	}
	if origUrl != nil {
		host = fmt.Sprintf("%s://%s", origUrl.Scheme, origUrl.Host)
	}
	registryHost := ""
	if registryDomain != "" {
		u, e := parseEndpoint(registryDomain)
		if e != nil {
			return nil, e
		}
		registryHost = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	}
	client := &http.Client{}
	return &Registry{
		Url:          origUrl,
		Host:         host,
		RegistryHost: registryHost,
		client:       client}, nil
}

func (reg *Registry) log(format string, args ...interface{}) {
	if reg.Logger == nil {
		return
	}
	reg.Logger(format, args...)
}

func (reg *Registry) GetToken(username string, password string, reposName string) (string, error) {
	if reg.Host == "" {
		return "", nil
	}
	u := fmt.Sprintf("%s/v1/repositories/%s/images", reg.Host, reposName)
	req, e := http.NewRequest("GET", u, nil)
	if e != nil {
		return "", e
	}
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	req.Header.Add("X-Docker-Token", "true")
	res, e := reg.client.Do(req)
	if e != nil {
		return "", e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("HTTP Error: %s", res.Status)
	}
	if reg.RegistryHost == "" {
		reg.RegistryHost = fmt.Sprintf("%s://%s", reg.Url.Scheme, res.Header.Get("X-Docker-Endpoints"))
		reg.log("Got registry endpoint from the server: %s", reg.RegistryHost)
	} else {
		reg.log("Registry endpoint overridden to %s", reg.RegistryHost)
	}
	token := res.Header.Get("X-Docker-Token")
	reg.log("Got token: %s", token)
	return token, nil
}

func (reg *Registry) doGet(token string, url string) (*http.Response, error) {
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return nil, e
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	res, e := reg.client.Do(req)
	if e != nil {
		return nil, e
	}
	if res.StatusCode != 200 {
		res.Body.Close()
		return nil, fmt.Errorf("HTTP Error: %s", res.Status)
	}
	return res, nil
}

func (reg *Registry) ReposTags(token string, reposName string) (map[string]string, error) {
	u := fmt.Sprintf("%s/v1/repositories/%s/tags", reg.RegistryHost, reposName)
	res, e := reg.doGet(token, u)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	rawJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	if err := json.Unmarshal(rawJSON, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (reg *Registry) LayerJson(token string, layerId string) (*LayerJson, error) {
	u := fmt.Sprintf("%s/v1/images/%s/json", reg.RegistryHost, layerId)
	res, e := reg.doGet(token, u)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	rawJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := &LayerJson{}
	if err := json.Unmarshal(rawJSON, &result); err != nil {
		return nil, err
	}
	result.Size, _ = strconv.Atoi(res.Header.Get("X-Docker-Size"))
	return result, nil
}

func (reg *Registry) LayerAncestry(token string, layerId string) (*[]string, error) {
	u := fmt.Sprintf("%s/v1/images/%s/ancestry", reg.RegistryHost, layerId)
	res, e := reg.doGet(token, u)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	rawJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := new([]string)
	if err := json.Unmarshal(rawJSON, &result); err != nil {
		return nil, err
	}
	return result, nil
}
