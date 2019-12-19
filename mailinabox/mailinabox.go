package mailinabox

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alinbalutoiu/dynamicdns-go/utils"

	log "github.com/sirupsen/logrus"
)

type MailInABoxClient struct {
	config         Config
	updateEndpoint string
}

func NewClient(config Config) *MailInABoxClient {
	miabc := &MailInABoxClient{
		config: config,
		updateEndpoint: fmt.Sprintf("%v/admin/dns/custom/%v",
			config.APIUrl, config.Domain),
	}

	return miabc
}

func (miabc *MailInABoxClient) UpdateIP() error {
	ip, err := utils.IPChanged(miabc.config.Domain)
	if err != nil {
		return err
	} else if ip == "" {
		// The ip did not change
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", miabc.updateEndpoint, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(miabc.config.Username, miabc.config.Password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	miabResponse := string(body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("mailinabox API returned an error "+
			"(%v), response: %v", resp.StatusCode, miabResponse)
	}

	log.Infof("Mail-In-A-Box reply: %v", miabResponse)

	// Set to wait for DNS propagation
	utils.SetWaitDNSPropagation(true)

	return nil
}
