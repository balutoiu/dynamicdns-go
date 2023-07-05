package ovhdomains

import (
	"fmt"

	"github.com/balutoiu/dynamicdns-go/utils"
	"github.com/ovh/go-ovh/ovh"
)

type OvhDomainsClient struct {
	config Config
	client *ovh.Client
}

func NewClient(config Config) *OvhDomainsClient {
	ovhClient, err := ovh.NewClient(
		"ovh-eu",
		config.ApplicationKey,
		config.ApplicationSecret,
		config.ConsumerKey)

	if err != nil {
		return nil
	}

	return &OvhDomainsClient{
		config: config,
		client: ovhClient,
	}
}

func (ovhdc *OvhDomainsClient) UpdateIP() error {
	hostname := ovhdc.config.ZoneName
	if ovhdc.config.SubDomain != "" {
		hostname = fmt.Sprintf("%s.%s", ovhdc.config.SubDomain, hostname)
	}
	ip, err := utils.IPChanged(hostname)
	if err != nil {
		return err
	} else if ip == "" {
		// The ip did not change
		return nil
	}

	var response []int64

	url := fmt.Sprintf("/domain/zone/%s/record?fieldType=A", ovhdc.config.ZoneName)
	if ovhdc.config.SubDomain != "" {
		url = fmt.Sprintf("%s&subDomain=%s", url, ovhdc.config.SubDomain)
	}
	err = ovhdc.client.Get(url, &response)
	if err != nil {
		return err
	}

	err = ovhdc.client.Put(
		fmt.Sprintf("/domain/zone/%s/record/%d", ovhdc.config.ZoneName, response[0]),
		map[string]string{"target": ip},
		nil)
	if err != nil {
		return err
	}

	err = ovhdc.client.Post(
		fmt.Sprintf("/domain/zone/%s/refresh", ovhdc.config.ZoneName),
		nil,
		nil)
	if err != nil {
		return err
	}

	// Set to wait for DNS propagation
	utils.SetWaitDNSPropagation(true)

	return nil
}
