// nmon2influxdb
// author: adejoux@djouxtech.net

package nmon2influxdblib

import (
	"bufio"
	"bytes"
	"fmt"
	influxdbclient2 "github.com/adejoux/nmon2influxdb/influxdbv2/influxdbclient"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/naoina/toml"
	"github.com/urfave/cli/v2"
)

//used for debug
var secretUser = "secretuser"
var secretPassword = "secret"

// Config is the configuration structure used by nmon2influxdb
type Config struct {
	Debug                 bool
	DebugFile             string
	Timezone              string
	InfluxdbUser          string
	InfluxdbPassword      string
	InfluxdbServer        string
	InfluxdbPort          string
	InfluxdbSecure        bool
	InfluxdbSkipCertCheck bool
	InfluxdbDatabase      string
	InfluxdbOrganization  string
	GrafanaUser           string
	GrafanaPassword       string
	GrafanaURL            string `toml:"grafana_URL"`
	GrafanaAccess         string
	GrafanaDatasource     string
	HMCServer             string `toml:"hmc_server"`
	HMCUser               string `toml:"hmc_user"`
	HMCPassword           string `toml:"hmc_password"`
	HMCDatabase           string `toml:"hmc_database"`
	HMCDataRetention      string `toml:"hmc_data_retention"`
	HMCManagedSystem      string `toml:"hmc_managed_system"`
	HMCManagedSystemOnly  bool   `toml:"hmc_managed_system_only"`
	HMCSamples            int    `toml:"hmc_samples"`
	HMCTimeout            int    `toml:"hmc_timeout"`
	ImportSkipDisks       bool
	ImportAllCpus         bool
	ImportBuildDashboard  bool
	ImportForce           bool
	ImportSkipMetrics     string
	ImportLogDatabase     string
	ImportLogRetention    string
	ImportDataRetention   string
	ImportSSHUser         string `toml:"import_ssh_user"`
	ImportSSHKey          string `toml:"import_ssh_key"`
	DashboardWriteFile    bool
	StatsLimit            int
	StatsSort             string
	StatsFilter           string
	StatsFrom             string
	StatsTo               string
	StatsHost             string
	Metric                string `toml:"metric,omitempty"`
	ListFilter            string `toml:",omitempty"`
	ListHost              string `toml:",omitempty"`
	Inputs                Inputs `toml:"input"`
	ConfFile              string
}

// Inputs allows to put multiple input in the configuration file
type Inputs []Input

// Input specify how to apply new filters
type Input struct {
	Measurement string
	Name        string
	Match       string
	Tags        Tags `toml:"tag"`
}

// InitConfig setup initial configuration with sane values
func InitConfig() Config {
	currUser, _ := user.Current()
	home := currUser.HomeDir
	sshKey := filepath.Join(home, "/.ssh/id_rsa")

	return Config{Debug: false,
		Timezone:              "Europe/Paris",
		InfluxdbUser:          "root",
		InfluxdbPassword:      "root",
		InfluxdbServer:        "localhost",
		InfluxdbPort:          "8086",
		InfluxdbDatabase:      "nmon_reports",
		InfluxdbSecure:        false,
		InfluxdbSkipCertCheck: false,
		HMCUser:               "hscroot",
		HMCPassword:           "abc123",
		HMCDatabase:           "nmon2influxdbHMC",
		HMCSamples:            10,
		HMCTimeout:            30,
		GrafanaUser:           "admin",
		GrafanaPassword:       "admin",
		GrafanaURL:            "http://localhost:3000",
		GrafanaAccess:         "direct",
		GrafanaDatasource:     "nmon2influxdb",
		ImportSkipDisks:       false,
		ImportAllCpus:         false,
		ImportBuildDashboard:  false,
		ImportForce:           false,
		ImportLogDatabase:     "nmon2influxdb_log",
		ImportLogRetention:    "2d",
		ImportSSHUser:         currUser.Username,
		ImportSSHKey:          sshKey,
		DashboardWriteFile:    false,
		ImportSkipMetrics:     "JFSINODE|TOP|PCPU",
		StatsLimit:            20,
		StatsSort:             "mean",
		StatsFilter:           "",
		StatsFrom:             "",
		StatsTo:               "",
		StatsHost:             "",
	}
}

//GetCfgFile returns the current configuration file path
func GetCfgFile(c *cli.Context) string {
	// if configuration file exist in /etc/nmon2influxdb. Stop here.

	if c != nil {
		paramFile := c.String("config_path")
		if IsFile(paramFile) {
			return paramFile
		}
	}

	if IsFile("/etc/nmon2influxdb/nmon2influxdb.cfg") {
		return "/etc/nmon2influxdb/nmon2influxdb.cfg"
	}

	if OnCurrentPath, err := os.Getwd(); err == nil {
		currentCFGFile := fmt.Sprintf("%s/%s", OnCurrentPath, ".nmon2influxdb.cfg")
		if IsFile(currentCFGFile) {
			return currentCFGFile
		}
	}

	currUser, _ := user.Current()
	home := currUser.HomeDir
	return filepath.Join(home, ".nmon2influxdb.cfg")
}

//IsFile returns true if the file doesn't exist
func IsFile(file string) bool {
	stat, err := os.Stat(file)
	if err != nil {
		return false
	}
	if stat.Mode().IsRegular() {
		return true
	}

	return false
}

//BuildCfgFile creates a default configuration file
func (config *Config) BuildCfgFile(cfgfile string) {
	file, err := os.Create(cfgfile)
	CheckError(err)
	defer file.Close()
	writer := bufio.NewWriter(file)
	b, err := toml.Marshal(*config)
	CheckError(err)
	r := bytes.NewReader(b)
	r.WriteTo(writer)
	writer.Flush()
	log.Printf("Generating default configuration file : %s\n", cfgfile)
}

// LoadCfgFile loads current configuration file settings
func (config *Config) LoadCfgFile(c *cli.Context) (cfgfile string) {

	cfgfile = GetCfgFile(c)

	//it would be only if no conf file exists. And it will build a configuration file in the home directory
	if !IsFile(cfgfile) {
		config.BuildCfgFile(cfgfile)
	} else {
		log.Printf("Using configuration file %s\n", cfgfile)
	}

	file, err := os.Open(cfgfile)
	if err != nil {
		log.Printf("Error opening configuration file %s\n", cfgfile)
		return
	}

	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		CheckError(err)
	}

	if err := toml.Unmarshal(buf, &config); err != nil {
		log.Printf("syntax error in configuration file: %s \n", err.Error())
		os.Exit(1)
	}
	return
}

// AddDashboardParams initialize default parameters for dashboard
func (config *Config) AddDashboardParams() {
	dfltConfig := InitConfig()
	dfltConfig.LoadCfgFile(nil)

	config.GrafanaAccess = dfltConfig.GrafanaAccess
	config.GrafanaURL = dfltConfig.GrafanaURL
	config.GrafanaDatasource = dfltConfig.GrafanaDatasource
	config.GrafanaUser = dfltConfig.GrafanaUser
	config.GrafanaPassword = dfltConfig.GrafanaPassword
	config.DashboardWriteFile = dfltConfig.DashboardWriteFile
}

// ParseParameters parse parameter from command line in Config struct
func ParseParameters(c *cli.Context) (config *Config) {
	config = new(Config)
	*config = InitConfig()
	config.LoadCfgFile(c)

	config.ConfFile = c.String("config_path")
	config.HMCSamples = c.Int("samples")

	if config.ConfFile == "" {
		config.Metric = c.String("metric")
		config.StatsHost = c.String("statshost")
		config.StatsFrom = c.String("from")
		config.StatsTo = c.String("to")
		config.StatsLimit = c.Int("limit")
		config.StatsFilter = c.String("filter")
		config.ImportSkipDisks = c.Bool("nodisks")
		if c.IsSet("cpus") {
			config.ImportAllCpus = c.Bool("cpus")
		}
		config.ImportBuildDashboard = c.Bool("build")
		config.ImportSkipMetrics = c.String("skip_metrics")
		config.ImportLogDatabase = c.String("log_database")
		config.ImportLogRetention = c.String("log_retention")
		config.DashboardWriteFile = c.Bool("file")
		config.ListFilter = c.String("filter")
		config.ImportForce = c.Bool("force")
		config.ListHost = c.String("host")
		config.GrafanaUser = c.String("guser")
		config.GrafanaPassword = c.String("gpassword")
		config.GrafanaAccess = c.String("gaccess")
		config.GrafanaURL = c.String("gurl")
		config.GrafanaDatasource = c.String("datasource")
		config.Debug = c.Bool("debug")
		config.DebugFile = c.String("debug-file")
		config.HMCServer = c.String("hmc")
		config.HMCUser = c.String("hmcuser")
		config.HMCPassword = c.String("hmcpass")
		config.HMCManagedSystem = c.String("managed_system")
		config.HMCManagedSystemOnly = c.Bool("managed_system-only")
		config.HMCTimeout = c.Int("timeout")
		config.InfluxdbServer = c.String("server")
		config.InfluxdbUser = c.String("user")
		config.InfluxdbPort = c.String("port")
		config.InfluxdbDatabase = c.String("db")
		config.InfluxdbSecure = c.Bool("secure")
		config.InfluxdbSkipCertCheck = c.Bool("skip_cert_check")
		config.InfluxdbPassword = c.String("pass")
		config.Timezone = c.String("tz")
	}

	if len(config.DebugFile) > 0 {
		//if a debug file is set. Debug is true
		config.Debug = true

		debugFile, err := os.OpenFile(config.DebugFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		// never closing the file here to be able to change the log output in all packages
		// No better solution for now.
		//		defer debugFile.Close()
		log.SetOutput(debugFile)
		log.Printf("NEW NMON2INFLUXDB EXECUTION\n")

	}

	if config.ImportBuildDashboard {
		config.AddDashboardParams()
	}

	return config
}

// ConnectDB connect to the specified influxdb database
func (config *Config) ConnectDB(db, org string) *influxdbclient2.InfluxDB {
	influxdbConfig := influxdbclient2.InfluxDBConfig{
		Host:          config.InfluxdbServer,
		Port:          config.InfluxdbPort,
		Database:      db,
		Organization:  org,
		User:          config.InfluxdbUser,
		Pass:          config.InfluxdbPassword,
		Debug:         config.Debug,
		Secure:        config.InfluxdbSecure,
		SkipCertCheck: config.InfluxdbSkipCertCheck,
	}
	influxdb, err := influxdbclient2.NewInfluxDB(influxdbConfig)
	CheckError(err)

	return &influxdb
}

// GetDB create or get the influxdb database used for nmon data
func (config *Config) GetDB(dbType string) *influxdbclient2.InfluxDB {

	db := config.InfluxdbDatabase
	org := config.InfluxdbOrganization
	retention := config.ImportDataRetention

	if dbType == "hmc" {
		db = config.HMCDatabase
		retention = config.HMCDataRetention
	}

	influxdb := config.ConnectDB(db, org)

	if exist, _ := influxdb.ExistDB(db); exist != true {
		log.Printf("Creating InfluxDB database %s\n", db)
		_, createErr := influxdb.CreateDB(db)
		CheckError(createErr)
	}

	// update default retention policy if ImportDataRetention is set
	if len(retention) > 0 {
		// Get default retention policy name
		policyName, policyErr := influxdb.GetDefaultRetentionPolicy()
		CheckError(policyErr)
		log.Printf("Updating  %s retention policy to keep only the last %s days. Timestamp based.\n", policyName, retention)
		_, err := influxdb.UpdateRetentionPolicy(policyName, retention, true)
		CheckError(err)
	}
	return influxdb
}

// GetLogDB create or get the influxdb database like defined in config
func (config *Config) GetLogDB() *influxdbclient2.InfluxDB {

	influxdb := config.ConnectDB(config.ImportLogDatabase, config.InfluxdbOrganization)

	if exist, _ := influxdb.ExistDB(config.ImportLogDatabase); exist != true {
		_, err := influxdb.CreateDB(config.ImportLogDatabase)
		CheckError(err)
		_, err = influxdb.SetRetentionPolicy("log_retention", config.ImportLogRetention, true)
		CheckError(err)
	} else {
		logPolicyName, logPolicyErr := influxdb.GetDefaultRetentionPolicy()
		CheckError(logPolicyErr)
		_, err := influxdb.UpdateRetentionPolicy(logPolicyName, config.ImportLogRetention, true)
		CheckError(err)
	}
	return influxdb
}

// Sanitized returns a copy of the config struct without the password. Used for debug
func (config *Config) Sanitized() (debugConfig Config) {
	debugConfig = *config
	debugConfig.HMCUser = secretUser
	debugConfig.HMCPassword = secretPassword
	debugConfig.GrafanaUser = secretUser
	debugConfig.GrafanaPassword = secretPassword
	debugConfig.InfluxdbUser = secretUser
	debugConfig.InfluxdbPassword = secretPassword
	return
}
