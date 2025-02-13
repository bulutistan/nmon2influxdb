# nmon2influxdb


This application take a nmon file and upload it in a [InfluxDB](influxdb.com) database.
It generates also a dashboard to allow data visualization in [Grafana](http://grafana.org/).
It's working on linux only for now.

## Added Influx v2.0 without Grafana!

### Can't execute query and auto create database on this release.
### Only work with credential that username and password.
### We have don't try with grafana, maybe you will tried.
### Change database property with bucket on configuration.
### Added organization property on configuration.

* influxdb_database="<YOUR_BUCKET_HERE>"
* influxdb_organization="<YOUR_ORGANIZATION_HERE>"

### Added argument for configuration file as relative path.
hmc import --samples SAMPLE_VAL --config_path CUSTOM_PATH_HERE

## Great thanks for this project of Adejoux!


# Demo

A live demo is available at : [demo.nmon2influxdb.org](http://demo.nmon2influxdb.org)

user/password: demo/demo

It's a read only editor user. You can change anything but cannot save it.

# Dashboards available on grafana.com

Multiple ready-to-use dashboards are available on [grafana.com](https://grafana.com/adejoux)

The following dashboards are available:

  * [AIX NMON report](https://grafana.com/dashboards/1555)
  * [AIX NMON Individual disks](https://grafana.com/dashboards/1701)
  * [Power Systems HMC partition view](https://grafana.com/dashboards/1510)
  * [Power Systems HMC system view](https://grafana.com/dashboards/1465)

# Download

Go to my [github repository Releases section](https://github.com/adejoux/nmon2influxdb/releases)

You can build the binary from source. You need to have a working GO environment, see [golang.org installation instructions](https://golang.org/doc/install). Check GOPATH environment variable to be sure it's defined.

~~~
go get -u github.com/adejoux/nmon2influxdb
cd $GOPATH/src/github.com/adejoux/nmon2influxdb
go build
~~~

**[FULL Documentation available here](https://nmon2influxdb.org)**


Copyright
==========

The code is licensed as GNU AGPLv3. See the LICENSE file for the full license.

Copyright (c) 2014 Alain Dejoux <adejoux@djouxtech.net>
