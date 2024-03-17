package ovhdomains

import (
	"fmt"

	"github.com/balutoiu/dynamicdns-go/utils"
	"github.com/ovh/go-ovh/ovh"

	log "github.com/sirupsen/logrus"
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
		log.Debug("IP did not change")
		return nil
	}

	var response []int64

	log.Debug("Getting DNS record ID")
	url := fmt.Sprintf("/domain/zone/%s/record?fieldType=A&subDomain=%s", ovhdc.config.ZoneName, ovhdc.config.SubDomain)
	err = ovhdc.client.Get(url, &response)
	if err != nil {
		return err
	}
	if len(response) == 0 {
		return fmt.Errorf("no DNS record found for %s.%s", ovhdc.config.ZoneName, ovhdc.config.SubDomain)
	}
	if len(response) > 1 {
		return fmt.Errorf("multiple DNS records found %s.%s", ovhdc.config.ZoneName, ovhdc.config.SubDomain)
	}

	log.Debugf("Updating DNS record %d with IP %s", response[0], ip)
	err = ovhdc.client.Put(
		fmt.Sprintf("/domain/zone/%s/record/%d", ovhdc.config.ZoneName, response[0]),
		map[string]string{"target": ip},
		nil)
	if err != nil {
		return err
	}

	log.Debug("Refreshing DNS zone")
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
