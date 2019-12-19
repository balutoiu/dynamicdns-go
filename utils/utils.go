package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const ipifyURL = "https://api.ipify.org?format=json"

var (
	waitDnsPropagation bool
	previousIP         = ""
)

type PublicIP struct {
	IP string `json:"ip"`
}

func GetCurrentPublicIP(target interface{}) error {
	resp, err := http.Get(ipifyURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func GetDNSIp(domainName string) (*PublicIP, error) {
	ips, err := net.LookupIP(domainName)
	if err != nil {
		return nil, err
	}
	ipv4List := []string{}
	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4List = append(ipv4List, ip.String())
		}
	}
	if len(ipv4List) != 1 {
		return nil, fmt.Errorf(fmt.Sprintf("Multiple ipv4 found: %v", ipv4List))
	}
	return &PublicIP{IP: ipv4List[0]}, nil
}

func SetWaitDNSPropagation(state bool) {
	waitDnsPropagation = state
}

// Returns the new IP if it changed
func IPChanged(hostname string) (string, error) {
	publicIP := &PublicIP{}
	err := GetCurrentPublicIP(publicIP)
	if err != nil {
		return "", err
	}
	log.Debugf("Current public IP from %v: %v", ipifyURL, publicIP.IP)

	publicIPUpstream, err := GetDNSIp(hostname)
	if err != nil {
		return "", err
	}
	log.Debugf("Current public IP from DNS lookup: %v", publicIPUpstream.IP)

	if publicIP.IP == publicIPUpstream.IP {
		if waitDnsPropagation {
			// DNS Propagated, should not wait anymore
			log.Infof("DNS is propagated")
			SetWaitDNSPropagation(false)
		}
		log.Infof("The IP did not change: %v", publicIP.IP)
		return "", nil
	} else {
		if previousIP != publicIP.IP {
			SetWaitDNSPropagation(false)
			previousIP = publicIP.IP
		} else if waitDnsPropagation {
			log.Infof("IP already updated, waiting for DNS propagation")
			return "", nil
		}
		log.Infof("IP changed: %v != %v", publicIP.IP, publicIPUpstream.IP)
		return publicIP.IP, nil
	}
}
