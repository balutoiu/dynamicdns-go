package googledomains

import (
	"io/ioutil"
	"net/http"

	"github.com/alinbalutoiu/dynamicdns-go/utils"

	log "github.com/sirupsen/logrus"
)

// Doc: https://support.google.com/domains/answer/6147083?hl=en

const googleDomainsURL = "https://domains.google.com/nic/update"

type GoogleDomainsClient struct {
	config Config
}

func NewClient(config Config) *GoogleDomainsClient {
	gdc := &GoogleDomainsClient{
		config: config,
	}

	return gdc
}

func (gdc *GoogleDomainsClient) UpdateIP() error {
	ip, err := utils.IPChanged(gdc.config.Domain)
	if err != nil {
		return err
	} else if ip == "" {
		// The ip did not change
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", googleDomainsURL, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(gdc.config.Username, gdc.config.Password)
	q := req.URL.Query()
	q.Add("hostname", gdc.config.Domain)
	q.Add("myip", ip)
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

	if err, ok := ErrorMap[string(body)]; ok {
		return err
	}

	log.Infof("Body: %v", string(body))

	// Set to wait for DNS propagation
	utils.SetWaitDNSPropagation(true)

	return nil
}
