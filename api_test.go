package api

import (
	"errors"
	"log"
	"os"
	"testing"
)

var client *API

var errMissingCredentials = errors.New("Missing SCALEWAY_ORGANIZATION and/or SCALEWAY_TOKEN environment variable")

func buildClient(region string) (*API, error) {
	if os.Getenv("SCALEWAY_ORGANIZATION") == "" || os.Getenv("SCALEWAY_TOKEN") == "" {
		return nil, errMissingCredentials
	}
	return New(os.Getenv("SCALEWAY_ORGANIZATION"), os.Getenv("SCALEWAY_TOKEN"), region)
}

func TestMain(m *testing.M) {
	c, err := buildClient("par1")
	if err != nil && err != errMissingCredentials {
		log.Printf("Unable to create scaleway client")
		os.Exit(1)
	}
	client = c
	code := m.Run()
	os.Exit(code)
}
