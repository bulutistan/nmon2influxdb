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

	app := new(application.AppCli).Init(config).MakeCommands().MakeFlags().Ready()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

}
