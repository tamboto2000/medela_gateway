package main

import (
	"flag"
	"fmt"

	medelagateway "github.com/tamboto2000/medela_gateway"
)

func main() {
	confFile := flag.String("c", "", "path to config file")
	port := flag.String("p", "8080", "port to listen")
	flag.Parse()

	conf, err := medelagateway.ParseConfigFromFile(*confFile)

	if err != nil {
		fmt.Println("error on reading config file:", err.Error())
		return
	}

	r := medelagateway.InitRouter(conf)
	r.Logger.Fatal(r.Start(fmt.Sprintf(":%s", *port)))
}
