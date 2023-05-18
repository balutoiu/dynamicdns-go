package config

import (
	"fmt"
	"os"

	"github.com/balutoiu/dynamicdns-go/providers/googledomains"
	"github.com/balutoiu/dynamicdns-go/providers/mailinabox"
	"github.com/balutoiu/dynamicdns-go/providers/ovhdomains"

	"gopkg.in/yaml.v3"
)

const (
	GOOGLE_DOMAINS = "googledomains"
	MAIL_IN_A_BOX  = "mailinabox"
	OVH_DOMAINS    = "ovhdomains"
)

type Config struct {
	GoogleDomains googledomains.Config `yaml:"googledomains,omitempty"`
	MailInABox    mailinabox.Config    `yaml:"mailinabox,omitempty"`
	OvhDomains    ovhdomains.Config    `yaml:"ovhdomains,omitempty"`
}

func GetConfig(filePath, dnsProvider string) (interface{}, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg *Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	if err := validateConfig(cfg, dnsProvider); err != nil {
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
		return nil, fmt.Errorf("dns provider not supported")
	}
}

func validateConfig(cfg *Config, dnsProvider string) error {
	switch dnsProvider {
	case GOOGLE_DOMAINS:
		if cfg.GoogleDomains.Username == "" {
			return fmt.Errorf("missing username from configuration")
		}
		if cfg.GoogleDomains.Password == "" {
			return fmt.Errorf("missing password from configuration")
		}
		if cfg.GoogleDomains.Domain == "" {
			return fmt.Errorf("missing domain from configuration")
		}
	case MAIL_IN_A_BOX:
		if cfg.MailInABox.Username == "" {
			return fmt.Errorf("missing username from configuration")
		}
		if cfg.MailInABox.Password == "" {
			return fmt.Errorf("missing password from configuration")
		}
		if cfg.MailInABox.Domain == "" {
			return fmt.Errorf("missing domain from configuration")
		}
		if cfg.MailInABox.APIUrl == "" {
			return fmt.Errorf("missing api_url from configuration")
		}
	case OVH_DOMAINS:
		if cfg.OvhDomains.ApplicationKey == "" {
			return fmt.Errorf("missing application_key from configuration")
		}
		if cfg.OvhDomains.ApplicationSecret == "" {
			return fmt.Errorf("missing application_secret from configuration")
		}
		if cfg.OvhDomains.ConsumerKey == "" {
			return fmt.Errorf("missing consumer_key from configuration")
		}
		if cfg.OvhDomains.ZoneName == "" {
			return fmt.Errorf("missing zone_name from configuration")
		}
		if cfg.OvhDomains.SubDomain == "" {
			return fmt.Errorf("missing sub_domain from configuration")
		}
	default:
		return fmt.Errorf("dns provider not supported")
	}
	return nil
}
