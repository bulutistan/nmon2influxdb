package influxdbclient

import (
	"fmt"
	client2 "github.com/adejoux/nmon2influxdb/influxdbv2/influxdb1clientv2/v2"
	"log"
	"time"
)

// InfluxDBConfig store all configuration parameters
type InfluxDBConfig struct {
	Host            string
	Port            string
	Database        string
	User            string
	Pass            string
	RetentionPolicy string
	Debug           bool
	Secure          bool
	SkipCertCheck   bool
}

// InfluxDB contains the main structures and methods used to parse nmon files and upload data in Influxdb
type InfluxDB struct {
	name        string
	Debug       bool
	count       int64
	points      client2.BatchPoints
	batchConfig client2.BatchPointsConfig
	client      client2.Client
	policy      string
}

// NewInfluxDB initialize a Influx structure
func NewInfluxDB(cfg InfluxDBConfig) (db InfluxDB, err error) {

	db.name = cfg.Database
	db.count = 0

	db.Debug = cfg.Debug

	http_mode := "http"

	if cfg.Secure {
		http_mode = "https"
	}
	influxdbURL := fmt.Sprintf("%s://%s:%s", http_mode, cfg.Host, cfg.Port)

	conf := client2.HTTPConfig{
		Addr:     influxdbURL,
		Username: cfg.User,
		Password: cfg.Pass,
	}

	if cfg.SkipCertCheck {
		conf.InsecureSkipVerify = true
	}

	db.client, err = client2.NewHTTPClient(conf)
	if err != nil {
		return
	}

	_, _, err = db.client.Ping(0)
	if err != nil {
		return
	}

	batchConfig := client2.BatchPointsConfig{Precision: "s", Database: cfg.Database}

	if len(cfg.RetentionPolicy) > 0 {
		batchConfig.RetentionPolicy = cfg.RetentionPolicy
	}

	db.batchConfig = batchConfig
	db.points, _ = client2.NewBatchPoints(batchConfig)
	return
}

// queryDB convenience function to query the database
func (db *InfluxDB) queryDB(cmd string, dbname string) (res []client2.Result, err error) {
	query := client2.NewQuery(cmd, dbname, "")
	if response, err := db.client.Query(query); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		log.Print(err)
	}
	return
}

// query convenience function to query Influxdb
func (db *InfluxDB) query(cmd string) (res []client2.Result, err error) {
	query := client2.NewQuery(cmd, "", "")
	if response, err := db.client.Query(query); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return
}

// CreateDB create the database with the dbname provided
func (db *InfluxDB) CreateDB(dbname string) (res []client2.Result, err error) {
	cmd := fmt.Sprintf("create database %s", dbname)
	res, err = db.queryDB(cmd, dbname)
	return
}

// DropDB drop the database named dbname
func (db *InfluxDB) DropDB(dbname string) (res []client2.Result, err error) {
	cmd := fmt.Sprintf("drop database %s", dbname)
	res, err = db.queryDB(cmd, dbname)
	return
}

// SetRetentionPolicy create a retention policy named retention. It's configured as default policy if defaultPolicy is true
func (db *InfluxDB) SetRetentionPolicy(policy string, retention string, defaultPolicy bool) (res []client2.Result, err error) {
	cmd := fmt.Sprintf("CREATE RETENTION POLICY \"%s\" ON \"%s\" DURATION %s REPLICATION 1", policy, db.name, retention)

	if defaultPolicy {
		cmd += " DEFAULT"
	}

	if db.Debug {
		log.Printf("SetRetentionPolicy: %s\n", cmd)
	}
	res, err = db.query(cmd)
	return
}

// UpdateRetentionPolicy update the rentetion policy
func (db *InfluxDB) UpdateRetentionPolicy(policy string, retention string, defaultPolicy bool) (res []client2.Result, err error) {
	cmd := fmt.Sprintf("ALTER RETENTION POLICY \"%s\" ON \"%s\" DURATION %s REPLICATION 1", policy, db.name, retention)

	if defaultPolicy {
		cmd += " DEFAULT"
	}

	if db.Debug {
		log.Printf("UpdateRetentionPolicy: %s\n", cmd)
	}

	res, err = db.query(cmd)
	return
}

// GetDefaultRetentionPolicy get the default retention policy name
func (db *InfluxDB) GetDefaultRetentionPolicy() (policyName string, err error) {
	cmd := fmt.Sprintf("SHOW RETENTION POLICIES ON \"%s\"", db.name)

	res, err := db.query(cmd)
	if res == nil {
		return
	}
	for _, policy := range res[0].Series[0].Values {

		// if default policy
		if policy[4] == true {
			policyName = policy[0].(string)
		}
	}
	return
}

// ShowDB returns a rray of dbname
func (db *InfluxDB) ShowDB() (databases []string, err error) {
	cmd := fmt.Sprintf("show databases")
	res, err := db.query(cmd)
	if err != nil {
		return
	}

	if res == nil {
		return
	}

	for _, dbs := range res[0].Series[0].Values {
		for _, db := range dbs {
			if str, ok := db.(string); ok {
				databases = append(databases, str)
			}
		}
	}
	if db.Debug {
		log.Printf("databases: %v\n", databases)
	}
	return
}

// ExistDB returns true if the db exists
func (db *InfluxDB) ExistDB(dbname string) (check bool, err error) {
	dbs, err := db.ShowDB()
	check = false

	if err != nil {
		return
	}

	for _, val := range dbs {
		if dbname == val {
			check = true
			return
		}
	}
	return
}

// AddPoint add standard point
func (db *InfluxDB) AddPoint(measurement string, timestamp time.Time, fields map[string]interface{}, tags map[string]string) {
	point, err := client2.NewPoint(measurement, tags, fields, timestamp)

	if err != nil {
		log.Println("Error: ", err.Error())
	}

	db.points.AddPoint(point)
	db.count++
}

// WritePoints wirte points to database
func (db *InfluxDB) WritePoints() (err error) {
	err = db.client.Write(db.points)
	if err != nil {
		log.Print(err)
	}
	return
}

// PointsCount return the number of influxdb points
func (db *InfluxDB) PointsCount() int64 {
	return db.count
}

// ClearPoints reset the number of points stored
func (db *InfluxDB) ClearPoints() {
	db.count = 0
	db.points, _ = client2.NewBatchPoints(db.batchConfig)
}

// ListMeasurement returns all measurements inside a TextSet
func (db *InfluxDB) ListMeasurement(filters *Filters) (tset *TextSet, err error) {

	query := "SHOW MEASUREMENTS"
	var fQuery FilterQuery
	if len(*filters) > 0 {
		fQuery.AddFilters(filters)
	}
	if len(fQuery.Content) > 0 {
		query += " WHERE " + fQuery.Content
	}
	if db.Debug {
		log.Printf("query: %s\n", query)
	}
	res, err := db.queryDB(query, db.name)
	if err != nil {
		return
	}
	return ConvertToTextSet(res), err

}

// ReadPoints perform a influxdb query on meqsurements and return points inside a dataset
func (db *InfluxDB) ReadPoints(fields string, filters *Filters, groupby string, serie string, from string, to string, function string) (ds []*DataSet, err error) {
	cmd := buildQuery(fields, filters, groupby, serie, from, to, function)
	if db.Debug {
		log.Printf("query: %s\n", cmd)
	}

	res, err := db.queryDB(cmd, db.name)
	if err != nil {
		return
	}
	ds = ConvertToDataSet(res)
	return
}

// ReadLastPoint perform a influxdb query on meqsurements and return the last point as string
func (db *InfluxDB) ReadLastPoint(fields string, filters *Filters, serie string) (result string, err error) {
	var filterQuery FilterQuery
	cmd := fmt.Sprintf("SELECT last(\"%s\") FROM \"%s\"", fields, serie)
	if len(*filters) > 0 {
		filterQuery.AddFilters(filters)
		cmd += fmt.Sprintf(" WHERE %s", filterQuery.Content)
	}

	if db.Debug {
		log.Printf("query: %s\n", cmd)
	}

	res, err := db.queryDB(cmd, db.name)
	if err != nil {
		return
	}

	if len(res) == 0 {
		return
	}

	if len(res[0].Series) == 0 {
		return
	}

	result = res[0].Series[0].Values[0][1].(string)
	return
}
