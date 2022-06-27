package main

import (
	"context"
	"flag"
	"github.com/broswen/terraform-provider-vex/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
)

var (
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	opts := providerserver.ServeOpts{
		Address: "vex",
		Debug:   debug,
	}
	err := providerserver.Serve(context.Background(), vex.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
