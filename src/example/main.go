package main

import (


	"github.com/36625090/involution"
	"github.com/36625090/involution/example/services/account/controller"
	"github.com/36625090/involution/logical"
	"github.com/36625090/involution/option"
	"log"
	"os"
)

func init() {

}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	opts, err := option.NewOptions()
	if err != nil {
		os.Exit(1)
	}

	factories := map[string]logical.Factory{
		"account": controller.Factory,
	}

	inv, err := involution.DefaultInvolution(opts, factories)

	if err != nil {
		log.Fatal(err)
		return
	}

	if err := inv.Start(); err != nil {
		log.Fatal(err)
		return
	}

}
