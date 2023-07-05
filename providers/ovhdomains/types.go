package ovhdomains

type Config struct {
	ApplicationKey    string `yaml:"application_key"`
	ApplicationSecret string `yaml:"application_secret"`
	ConsumerKey       string `yaml:"consumer_key"`
	ZoneName          string `yaml:"zone_name"`
	SubDomain         string `yaml:"sub_domain,omitempty"`
}
