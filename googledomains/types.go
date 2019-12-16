package googledomains

import "fmt"

var ErrorMap map[string]error = map[string]error{
	"nohost":  fmt.Errorf("The hostname does not exist, or does not have Dynamic DNS enabled."),
	"badauth": fmt.Errorf("The username / password combination is not valid for the specified host."),
	"notfqdn": fmt.Errorf("The supplied hostname is not a valid fully-qualified domain name."),
	"badagent": fmt.Errorf("	Your Dynamic DNS client is making bad requests. " +
		"Ensure the user agent is set in the request."),
	"abuse": fmt.Errorf("Dynamic DNS access for the hostname has been blocked due to " +
		"failure to interpret previous responses correctly."),
	"911": fmt.Errorf("An error happened on our end. Wait 5 minutes and retry."),
}
