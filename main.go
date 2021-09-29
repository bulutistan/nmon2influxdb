// nmon2influxdb
// import nmon data in InfluxDB

// author: adejoux@djouxtech.net

package main

import (
	"github.com/adejoux/nmon2influxdb/application"
	"log"
	"os"
	"time"
)

func main() {
	config := new(application.AppConfig).Init()

	for i := 0; i < 10; i++ {
		appCli := new(application.AppCli).Init(config)

		if err := appCli.Ready(os.Args); err != nil {
			log.Println(err)
		} else {
			i--
			log.Printf("waiting restart for %d seconds.", config.HMCTimeout)
			time.Sleep(time.Second * time.Duration(config.HMCTimeout))
		}
	}

}
