package main

import (
	"flag"
	"fmt"

	medelagateway "github.com/tamboto2000/medela_gateway"
)

func main() {
	confFile := flag.String("c", "", "path to config file")
	flag.Parse()

	conf, err := medelagateway.ParseConfigFromFile(*confFile)

	if err != nil {
		fmt.Println("error on reading config file:", err.Error())
		return
	}

	r := medelagateway.InitRouter(conf)
	r.Logger.Fatal(r.Start(":8080"))
}
