package main

import (
	"github.com/hashicorp/nomad/plugins/base"
	"github.com/hashicorp/nomad/plugins/drivers"
	"log"
)

func main() {
	log.Println("[INFO] Starting Bastille driver plugin...")
	driver, err := NewBastilleDriver()
	if err != nil {
		log.Fatalf("failed to create driver: %v", err)
	}

	pluginServeOpts := &base.ServeOpts{
		Factory: func() interface{} {
			return driver
		},
	}

	base.Serve(driver.Name(), pluginServeOpts)
}

