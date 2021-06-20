package application

import (
	"github.com/adejoux/nmon2influxdb/nmon2influxdblib"
	"log"
	"os"
)

type AppConfig struct {
}

func (c AppConfig) Init() nmon2influxdblib.Config {
	config := nmon2influxdblib.InitConfig()

	cfgfile := config.LoadCfgFile()
	if len(config.DebugFile) > 0 {
		debugFile, err := os.OpenFile("config.DebugFile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer debugFile.Close()
		log.SetOutput(debugFile)

	}

	log.Printf("Using configuration file %s\n", cfgfile)

	// cannot set values directly for boolean flags
	if config.DashboardWriteFile {
		os.Setenv("NMON2INFLUXDB_DASHBOARD_TO_FILE", "true")
	}

	if config.ImportSkipDisks {
		os.Setenv("NMON2INFLUXDB_SKIP_DISKS", "true")
	}

	if config.ImportAllCpus {
		os.Setenv("NMON2INFLUXDB_ADD_ALL_CPUS", "true")
	}

	if config.ImportBuildDashboard {
		os.Setenv("NMON2INFLUXDB_BUILD_DASHBOARD", "true")
	}

	if config.ImportForce {
		os.Setenv("NMON2INFLUXDB_FORCE", "true")
	}

	if config.InfluxdbSecure {
		os.Setenv("NMON2INFLUXDB_SECURE", "true")
	}
	if config.InfluxdbSkipCertCheck {
		os.Setenv("NMON2INFLUXDB_SKIP_CERT_CHECK", "true")
	}
	if len(config.ImportSkipMetrics) > 0 {
		os.Setenv("NMON2INFLUXDB_SKIP_METRICS", config.ImportSkipMetrics)
	}

	if len(config.HMCServer) > 0 {
		os.Setenv("NMON2INFLUXDB_HMC_SERVER", config.HMCServer)
	}

	if len(config.ImportSkipMetrics) > 0 {
		os.Setenv("NMON2INFLUXDB_HMC_USER", config.HMCServer)
	}

	return config
}
