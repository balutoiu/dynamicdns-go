package main

import (
	"fmt"
	"os"

	"github.com/alinbalutoiu/dynamicdns-go/googledomains"
	"github.com/alinbalutoiu/dynamicdns-go/mailinabox"
	"github.com/alinbalutoiu/dynamicdns-go/ovhdomains"

	"gopkg.in/yaml.v2"
)

type Config struct {
	GoogleDomains googledomains.Config `yaml:"googledomains,omitempty"`
	MailInABox    mailinabox.Config    `yaml:"mailinabox,omitempty"`
	OvhDomains    ovhdomains.Config    `yaml:"ovhdomains,omitempty"`
}

func getConfig(filePath, dnsProvider string) (interface{}, error) {
	var cfg *Config

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// var cfg googledomains.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	err = validateConfig(cfg, dnsProvider)

	if err != nil {
		return nil, err
	}
	switch dnsProvider {
	case GOOGLE_DOMAINS:
		return cfg.GoogleDomains, nil
	case MAIL_IN_A_BOX:
		return cfg.MailInABox, nil
	case OVH_DOMAINS:
		return cfg.OvhDomains, nil
	default:
		return nil, fmt.Errorf("DNS Provider not supported")
	}
}

func validateConfig(cfg *Config, dnsProvider string) error {
	switch dnsProvider {
	case GOOGLE_DOMAINS:
		if cfg.GoogleDomains.Username == "" {
			return fmt.Errorf("Missing username from configuration")
		}
		if cfg.GoogleDomains.Password == "" {
			return fmt.Errorf("Missing password from configuration")
		}
		if cfg.GoogleDomains.Domain == "" {
			return fmt.Errorf("Missing domain from configuration")
		}
	case MAIL_IN_A_BOX:
		if cfg.MailInABox.Username == "" {
			return fmt.Errorf("Missing username from configuration")
		}
		if cfg.MailInABox.Password == "" {
			return fmt.Errorf("Missing password from configuration")
		}
		if cfg.MailInABox.Domain == "" {
			return fmt.Errorf("Missing domain from configuration")
		}
		if cfg.MailInABox.APIUrl == "" {
			return fmt.Errorf("Missing api_url from configuration")
		}
	case OVH_DOMAINS:
		if cfg.OvhDomains.ApplicationKey == "" {
			return fmt.Errorf("Missing application_key from configuration")
		}
		if cfg.OvhDomains.ApplicationSecret == "" {
			return fmt.Errorf("Missing application_secret from configuration")
		}
		if cfg.OvhDomains.ConsumerKey == "" {
			return fmt.Errorf("Missing consumer_key from configuration")
		}
		if cfg.OvhDomains.ZoneName == "" {
			return fmt.Errorf("Missing zone_name from configuration")
		}
		if cfg.OvhDomains.SubDomain == "" {
			return fmt.Errorf("Missing sub_domain from configuration")
		}
	default:
		return fmt.Errorf("DNS Provider not supported")
	}
	return nil
}
