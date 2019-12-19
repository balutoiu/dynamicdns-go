package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alinbalutoiu/dynamicdns-go/googledomains"
	"github.com/alinbalutoiu/dynamicdns-go/mailinabox"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	GOOGLE_DOMAINS = "googledomains"
	MAIL_IN_A_BOX  = "mailinabox"
)

type DNSClient interface {
	UpdateIP() error
}

var (
	port                  int
	logLevel              string
	supportedDnsProviders = []string{GOOGLE_DOMAINS, MAIL_IN_A_BOX}
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

func main() {
	app := cli.NewApp()
	app.Name = "Dynamic DNS Go"
	app.Usage = "Starts the Dynamic DNS monitoring"
	app.Authors = []*cli.Author{
		&cli.Author{
			Name: "Alin Balutoiu",
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
			Value:       GOOGLE_DOMAINS,
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

func runApp(c *cli.Context) error {
	initLog()
	log.Infof("Log successfully initialized - %v", logLevel)
	log.Infof("DNS provider selected: %v", dnsProvider)
	log.Infof("Config location: %v", configFilePath)
	config, err := getConfig(configFilePath, dnsProvider)
	if err != nil {
		return err
	}

	log.Infof("Configuration: %+v", config)
	var client DNSClient
	switch dnsProvider {
	case GOOGLE_DOMAINS:
		client = googledomains.NewClient(config.(googledomains.Config))
	case MAIL_IN_A_BOX:
		client = mailinabox.NewClient(config.(mailinabox.Config))
	}
	log.Infof("Client initialized")

	log.Infof("Updating DNS with interval: %v", sleepInterval)
	for {
		err = client.UpdateIP()
		if err != nil {
			return err
		}
		log.Infof("Sleeping for %v", sleepInterval)
		time.Sleep(sleepInterval)
	}
	return nil
}
