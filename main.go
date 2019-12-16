package main

import (
	"log"
	"os"
	"time"

	"github.com/alinbalutoiu/dynamicdns-go/googledomains"
)

func main() {
	username, exists := os.LookupEnv("USERNAME")
	if !exists {
		panic("Could not find username")
	}
	password, exists := os.LookupEnv("PASSWORD")
	if !exists {
		panic("Could not find password")
	}
	domain, exists := os.LookupEnv("DOMAIN")
	if !exists {
		panic("Could not find domain")
	}

	gdc := googledomains.NewClient(username, password, domain)
	for {
		err := gdc.UpdateIP()
		if err != nil {
			panic(err)
		}
		log.Printf("Sleeping for 1 hour...")
		time.Sleep(1 * time.Hour)
	}
}
