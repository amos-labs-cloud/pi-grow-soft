package cmd

import "github.com/amos-labs-cloud/pi-grow-soft/pkg/controller"

var client *Client

type Client struct {
	controllerService *controller.Service
}

func Initialize(controllerService *controller.Service) {
	client = &Client{controllerService: controllerService}
}
