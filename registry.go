package main

import (
	"fmt"
	"net/url"
)

type Registry struct {
	Host string
	RegistryHost string
}

func NewRegistry(endpoint string, registryDomain string) (*Registry, error) {
	u, e := url.Parse(endpoint)
	if e != nil {
		return nil, e
	}
	if u.Host == "" {
		u.Host = u.Path
	}
	host := u.Host
	u, e = url.Parse(registryDomain)
	if e != nil {
		return nil, e
	}
	if u.Host == "" {
		u.Host = u.Path
	}
	registryHost := u.Host
	return &Registry{host, registryHost}, nil
}

func (reg *Registry) GetToken(username string, password string, reposName string) (string, error) {
	u := fmt.Sprintf("%s/v1/repositories/%s/images", reg.Host)
	res, e := http.Get(u)
	if e != nil {
		return "", nil
	}
	return "", nil
}
