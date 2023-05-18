package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/balutoiu/dynamicdns-go/config"
	"github.com/balutoiu/dynamicdns-go/providers"
	"github.com/balutoiu/dynamicdns-go/providers/googledomains"
	"github.com/balutoiu/dynamicdns-go/providers/mailinabox"
	"github.com/balutoiu/dynamicdns-go/providers/ovhdomains"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	supportedDnsProviders = []string{config.GOOGLE_DOMAINS, config.MAIL_IN_A_BOX, config.OVH_DOMAINS}
	logLevel              string
	dnsProvider           string
	configFilePath        string
	sleepInterval         time.Duration
)

func initLog() {
	logrusLogLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v (see --help for more info)", logLevel)
	}

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	log.SetLevel(logrusLogLevel)

	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
}

func runApp(c *cli.Context) error {
	initLog()
	log.Infof("Log successfully initialized - %v", logLevel)
	log.Infof("DNS provider selected: %v", dnsProvider)
	log.Infof("Config location: %v", configFilePath)
	cfg, err := config.GetConfig(configFilePath, dnsProvider)
	if err != nil {
		return err
	}

	var client providers.DNSClient
	switch dnsProvider {
	case config.GOOGLE_DOMAINS:
		client = googledomains.NewClient(cfg.(googledomains.Config))
	case config.MAIL_IN_A_BOX:
		client = mailinabox.NewClient(cfg.(mailinabox.Config))
	case config.OVH_DOMAINS:
		client = ovhdomains.NewClient(cfg.(ovhdomains.Config))
	}
	log.Infof("Client initialized")

	log.Infof("Updating DNS with interval: %v", sleepInterval)
	for {
		if err := client.UpdateIP(); err != nil {
			log.Warnf("Failed to update IP: %v", err)
		}
		log.Debugf("Sleeping for %v", sleepInterval)
		time.Sleep(sleepInterval)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Dynamic DNS Go"
	app.Usage = "Starts the Dynamic DNS monitoring"
	app.Authors = []*cli.Author{
		{
			Name: "Balutoiu",
		},
	}
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "loglevel",
			Usage:       "sets the log level (error, warn, info, debug)",
			Value:       "info",
			EnvVars:     []string{"LOG_LEVEL"},
			Destination: &logLevel,
		},
		&cli.StringFlag{
			Name: "dns-provider",
			Usage: fmt.Sprintf("DNS provider (currently supported: %v)",
				strings.Join(supportedDnsProviders[:], ",")),
			Value:       config.GOOGLE_DOMAINS,
			EnvVars:     []string{"DNS_PROVIDER"},
			Destination: &dnsProvider,
		},
		&cli.StringFlag{
			Name:        "config",
			Usage:       "Path to the config file (format supported: YAML)",
			Value:       "config.yaml",
			EnvVars:     []string{"CONFIG"},
			Destination: &configFilePath,
		},
		&cli.DurationFlag{
			Name:        "sleep-interval",
			Usage:       "Sleep interval between checks (ex: 1s, 1m, 1h)",
			EnvVars:     []string{"SLEEP_INTERVAL"},
			Value:       1 * time.Hour,
			Destination: &sleepInterval,
		},
	}
	app.Action = func(c *cli.Context) error {
		return runApp(c)
	}

	// Handle ctrl+c (https://github.com/urfave/cli/issues/945)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C detected, exiting...")
		os.Exit(0)
	}()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
