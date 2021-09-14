// nmon2influxdb
// import nmon data in InfluxDB

// author: adejoux@djouxtech.net

package main

import (
	"github.com/adejoux/nmon2influxdb/application"
	"log"
	"os"
)

func main() {
	config := new(application.AppConfig).Init()

	if err := new(application.AppCli).Init(config).Ready(os.Args); err != nil {
		log.Fatalln(err)
	}

}
