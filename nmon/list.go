// Package nmon provides wrapper on nmon reltaed commands
// import nmon report in InfluxDB
// author: adejoux@djouxtech.net
package nmon

import (
	"fmt"
	influxdbclient2 "github.com/adejoux/nmon2influxdb/influxdbv2/influxdbclient"
	"regexp"

	"github.com/adejoux/nmon2influxdb/nmon2influxdblib"
	"github.com/urfave/cli/v2"
	//	"os"
)

//ListMeasurement list all measurements in INFLUXDB database
func ListMeasurement(c *cli.Context) error {
	// parsing parameters
	config := nmon2influxdblib.ParseParameters(c)

	influxdb := config.ConnectDB(config.InfluxdbDatabase, config.InfluxdbOrganization)
	filters := new(influxdbclient2.Filters)

	if len(config.ListHost) > 0 {
		filters.Add("host", config.ListHost, "text")
	}

	measurements, err := influxdb.ListMeasurement(filters)
	nmon2influxdblib.CheckError(err)
	if measurements != nil {
		fmt.Printf("%s\n", measurements.Name)
		for _, value := range measurements.Datas {
			if len(config.ListFilter) == 0 {
				fmt.Printf("%s\n", value)
				continue
			}
			matched, _ := regexp.MatchString(config.ListFilter, value)
			if matched {
				fmt.Printf("%s\n", value)
			}
		}
	}
	return nil
}
