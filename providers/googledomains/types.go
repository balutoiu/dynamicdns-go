package googledomains

import "fmt"

var ErrorMap map[string]error = map[string]error{
	"nohost":   fmt.Errorf("hostname does not exist, or does not have dynamic DNS enabled"),
	"badauth":  fmt.Errorf("username / password combination is not valid for the specified host"),
	"notfqdn":  fmt.Errorf("the supplied hostname is not a valid fully-qualified domain name"),
	"badagent": fmt.Errorf("your dynamic DNS client is making bad requests (ensure the user agent is set in the request)"),
	"abuse":    fmt.Errorf("dynamic DNS access for the hostname has been blocked due to failure to interpret previous responses correctly"),
	"911":      fmt.Errorf("nn error happened on our end, wait 5 minutes and retry"),
}

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Domain   string `yaml:"domain"`
}
