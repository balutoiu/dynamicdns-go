package googledomains

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

// Doc: https://support.google.com/domains/answer/6147083?hl=en

const googleDomainsURL = "https://domains.google.com/nic/update"
const ipifyURL = "https://api.ipify.org?format=json"

type GoogleDomainsClient struct {
	Username string
	Password string
	Hostname string
	MyIP     string
}

type PublicIP struct {
	IP string `json:"ip"`
}

func NewClient(user, pass, hostname string) *GoogleDomainsClient {
	gdc := &GoogleDomainsClient{
		Username: user,
		Password: pass,
		Hostname: hostname,
	}

	return gdc
}

func getCurrentPublicIP(target interface{}) error {
	resp, err := http.Get(ipifyURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func getDNSIp(domainName string) (*PublicIP, error) {
	ips, err := net.LookupIP(domainName)
	if err != nil {
		return nil, err
	}
	if len(ips) != 1 {
		return nil, fmt.Errorf(fmt.Sprintf("Multiple ips found: %v", ips))
	}
	return &PublicIP{IP: ips[0].String()}, nil
}

func (gdc *GoogleDomainsClient) UpdateIP() error {
	publicIP := &PublicIP{}
	getCurrentPublicIP(publicIP)

	publicIPUpstream, err := getDNSIp(gdc.Hostname)
	if err != nil {
		return err
	}

	if publicIP.IP == publicIPUpstream.IP {
		log.Printf("The IP did not change: %v", publicIP.IP)
		return nil
	} else {
		log.Printf("IP changed: %v != %v", publicIP.IP, publicIPUpstream.IP)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", googleDomainsURL, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(gdc.Username, gdc.Password)
	q := req.URL.Query()
	q.Add("hostname", gdc.Hostname)
	q.Add("myip", publicIP.IP)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if error, ok := ErrorMap[string(body)]; ok {
		return error
	}

	log.Printf("Body: %v", string(body))
	return nil
}
