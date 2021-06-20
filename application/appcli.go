package application

import (
	"github.com/adejoux/nmon2influxdb/hmc"
	"github.com/adejoux/nmon2influxdb/nmon"
	"github.com/adejoux/nmon2influxdb/nmon2influxdblib"
	"github.com/urfave/cli/v2"
)

type AppCli struct {
	App    *cli.App
	Config nmon2influxdblib.Config
}

func (c AppCli) Init(config nmon2influxdblib.Config) AppCli {
	c.App = cli.NewApp()
	c.Config = config

	c.App.Name = "nmon2influxdb"
	c.App.Usage = "upload NMON stats to InfluxDB database"
	c.App.Version = "2.1.7"
	c.App.Authors = []*cli.Author{{Name: "Alain Dejoux", Email: "adejoux@djouxtech.net"}}

	return c
}

func (c AppCli) MakeCommands() AppCli {

	c.App.Commands = []*cli.Command{
		{
			Name:  "import",
			Usage: "import nmon files",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "skip_metrics",
					Usage:   "skip metrics",
					EnvVars: []string{"NMON2INFLUXDB_SKIP_METRICS"},
				},
				&cli.BoolFlag{
					Name:    "nodisks",
					Aliases: []string{"nd"},
					Usage:   "skip disk metrics",
					EnvVars: []string{"NMON2INFLUXDB_SKIP_DISKS"},
				},
				&cli.BoolFlag{
					Name:    "cpus",
					Aliases: []string{"c"},
					Usage:   "add per cpu metrics",
					EnvVars: []string{"NMON2INFLUXDB_ADD_ALL_CPU"},
				},
				&cli.BoolFlag{
					Name:    "build",
					Aliases: []string{"b"},
					Usage:   "build dashboard",
					EnvVars: []string{"NMON2INFLUXDB_BUILD_DASHBOARD"},
				},
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "force import",
					EnvVars: []string{"NMON2INFLUXDB_FORCE"},
				},
				&cli.StringFlag{
					Name:  "log_database",
					Usage: "influxdb database used to log imports",
					Value: c.Config.ImportLogDatabase,
				},
				&cli.StringFlag{
					Name:  "log_retention",
					Usage: "import log retention",
					Value: c.Config.ImportLogRetention,
				},
			},
			Action: nmon.Import,
		},
		{
			Name:  "dashboard",
			Usage: "generate a dashboard from a nmon file or template",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "generate Grafana dashboard file",
					EnvVars: []string{"NMON2INFLUXDB_DASHBOARD_TO_FILE"},
				},
				&cli.StringFlag{
					Name:  "guser",
					Usage: "grafana user",
					Value: c.Config.GrafanaUser,
				},
				&cli.StringFlag{
					Name:    "gpassword",
					Aliases: []string{"gpass"},
					Usage:   "grafana password",
					Value:   c.Config.GrafanaPassword,
				},
				&cli.StringFlag{
					Name:  "gaccess",
					Usage: "grafana datasource access mode : direct or proxy",
					Value: c.Config.GrafanaAccess,
				},
				&cli.StringFlag{
					Name:  "gurl",
					Usage: "grafana url",
					Value: c.Config.GrafanaURL,
				},
				&cli.StringFlag{
					Name:  "datasource",
					Usage: "grafana datasource",
					Value: c.Config.GrafanaDatasource,
				},
			},
			Action: nmon.Dashboard,
		},
		{
			Name:  "stats",
			Usage: "generate stats from a InfluxDB metric",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "metric",
					Aliases: []string{"m"},
					Usage:   "mandatory metric for stats",
				},
				&cli.StringFlag{
					Name:    "statshost",
					Aliases: []string{"s"},
					Usage:   "host metrics",
					Value:   c.Config.StatsHost,
				},
				&cli.StringFlag{
					Name:    "from",
					Aliases: []string{"f"},
					Usage:   "from date",
					Value:   c.Config.StatsFrom,
				},
				&cli.StringFlag{
					Name:    "to",
					Aliases: []string{"t"},
					Usage:   "to date",
					Value:   c.Config.StatsTo,
				},
				&cli.StringFlag{
					Name:  "sort",
					Usage: "field to perform sort",
					Value: c.Config.StatsSort,
				},
				&cli.IntFlag{
					Name:  "limit,l",
					Usage: "limit the output",
					Value: c.Config.StatsLimit,
				},
				&cli.StringFlag{
					Name:  "filter",
					Usage: "specify a filter on fields",
					Value: c.Config.StatsFilter,
				},
			},
			Action: nmon.Stat,
		},
		{
			Name:  "list",
			Usage: "list InfluxDB metrics or measurements",
			Subcommands: []*cli.Command{
				{
					Name:  "measurement",
					Usage: "list InfluxDB measurements",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "host",
							Usage: "only for specified host",
						},
						&cli.StringFlag{
							Name:    "filter",
							Aliases: []string{"f"},
							Usage:   "filter measurement",
						},
					},
					Action: nmon.ListMeasurement,
				},
			},
		},
		{
			Name:  "hmc",
			Usage: "load hmc data",
			Subcommands: []*cli.Command{
				{
					Name:   "import",
					Usage:  "import HMC PCM data",
					Action: hmc.Import,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "hmc",
							Usage:   "HMC server",
							EnvVars: []string{"NMON2INFLUXDB_HMC_SERVER"},
						},
						&cli.StringFlag{
							Name:  "hmcuser",
							Usage: "HMC user",
							Value: c.Config.HMCUser,
						},
						&cli.StringFlag{
							Name:  "hmcpass",
							Usage: "HMC password",
							Value: c.Config.HMCPassword,
						},
						&cli.StringFlag{
							Name:  "managed_system,m",
							Usage: "only import from this managed system",
							Value: c.Config.HMCManagedSystem,
						},
						&cli.BoolFlag{
							Name:  "managed_system-only,sys-only",
							Usage: "skip partition metrics",
						},
						&cli.IntFlag{
							Name:  "samples",
							Usage: "import latest <value> samples",
							Value: c.Config.HMCSamples,
						},
						&cli.IntFlag{
							Name:  "timeout",
							Usage: "HMC connection timeout",
							Value: c.Config.HMCTimeout,
						},
					},
				},
			},
		},
	}

	return c
}

func (c AppCli) MakeFlags() AppCli {

	c.App.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "server,s",
			Usage: "InfluxDB server and port",
			Value: c.Config.InfluxdbServer,
		},
		&cli.StringFlag{
			Name:  "port,p",
			Usage: "InfluxDB port",
			Value: c.Config.InfluxdbPort,
		},
		&cli.BoolFlag{
			Name:    "secure",
			Usage:   "use ssl for InfluxDB",
			EnvVars: []string{"NMON2INFLUXDB_SECURE"},
		},
		&cli.BoolFlag{
			Name:    "skip_cert_check",
			Usage:   "skip cert check for ssl connzction to InfluxDB",
			EnvVars: []string{"NMON2INFLUXDB_SKIP_CERT_CHECK"},
		},
		&cli.StringFlag{
			Name:  "db,d",
			Usage: "InfluxDB database",
			Value: c.Config.InfluxdbDatabase,
		},
		&cli.StringFlag{
			Name:  "user,u",
			Usage: "InfluxDB administrator user",
			Value: c.Config.InfluxdbUser,
		},
		&cli.StringFlag{
			Name:  "pass",
			Usage: "InfluxDB administrator pass",
			Value: c.Config.InfluxdbPassword,
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "debug mode",
			EnvVars: []string{"NMON2INFLUXDB_DEBUG"},
		},
		&cli.StringFlag{
			Name:  "debug-file",
			Usage: "debug file",
			Value: c.Config.DebugFile,
		},
		&cli.StringFlag{
			Name:  "tz,t",
			Usage: "timezone",
			Value: c.Config.Timezone,
		},
	}

	return c
}

func (c AppCli) Ready() *cli.App {
	return c.App
}
