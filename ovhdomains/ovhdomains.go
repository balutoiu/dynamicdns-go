package ovhdomains

import (
	"fmt"

	"github.com/alinbalutoiu/dynamicdns-go/utils"
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
	ip, err := utils.IPChanged(
		fmt.Sprintf("%s.%s", ovhdc.config.SubDomain, ovhdc.config.ZoneName))
	if err != nil {
		return err
	} else if ip == "" {
		// The ip did not change
		return nil
	}

	var response []int64

	err = ovhdc.client.Get(
		fmt.Sprintf("/domain/zone/%s/record?fieldType=A&subDomain=%s", ovhdc.config.ZoneName, ovhdc.config.SubDomain),
		&response)
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

	// Set to wait for DNS propagation
	utils.SetWaitDNSPropagation(true)

	return nil
}
