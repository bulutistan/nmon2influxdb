package influxdbclient

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

const testdb = "testdb"

const timeformat = "15:04:05,02-Jan-2006"

var cfg = InfluxDBConfig{
	Host:     "192.168.56.101",
	Port:     "8086",
	Database: testdb,
	User:     "root",
	Pass:     "root",
}
var fields = map[string]interface{}{
	"a": 1.0,
	"b": 2.0,
	"c": 3.0,
	"d": 4.0,
}

var fields2 = map[string]interface{}{
	"a": 11.0,
	"b": 12.0,
	"c": 13.0,
	"d": 14.0,
}

var tags = map[string]string{
	"test": "yes",
}

func Test_BadConnect(t *testing.T) {
	badCfg := cfg
	badCfg.Port = "8087"
	_, err := NewInfluxDB(badCfg)
	assert.NotNil(t, err, "We are expecting error and didn't got one")
}

func Test_GoodConnect(t *testing.T) {
	_, err := NewInfluxDB(cfg)
	assert.Nil(t, err, "We are expecting no errors and got one")
}

func Test_CreateDB(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)
	res, err := testDB.CreateDB(testdb)
	assert.Nil(t, err, "We are expecting no error and got one")
	assert.NotNil(t, res, "We are expecting no error and got one")
}

func Test_ShowDB(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)
	res, err := testDB.ShowDB()
	assert.Nil(t, err, "We are expecting no error and got one")
	assert.Contains(t, res, testdb, "We are expecting to retrieve testdb database name")
}

func Test_ExistDB(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)
	res, err := testDB.ExistDB(testdb)
	assert.Nil(t, err, "We are expecting no error and got one")
	assert.Equal(t, res, true)
}

func Test_SetRetentionPolicy(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)
	res, err := testDB.SetRetentionPolicy("testpol", "1d", false)
	assert.Nil(t, err, "We are expecting no error and got one")
	assert.NotNil(t, res, "We are expecting a result and got none")
}

func Test_UpdateRetentionPolicy(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)
	res, err := testDB.UpdateRetentionPolicy("testpol", "10d", false)
	assert.Nil(t, err, "We are expecting no error and got one")
	assert.NotNil(t, res, "We are expecting a result and got none")
}

func Test_GetDefaultRetentionPolicy(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)
	res, err := testDB.GetDefaultRetentionPolicy()
	assert.Nil(t, err, "We are expecting no error and got one")
	assert.Equal(t, res, "autogen")
}

func Test_AddPoint(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)
	ti, _ := ConvertTimeStamp("23:55:28,13-MAY-2015", "Europe/Paris")
	testDB.AddPoint("test", ti, fields, tags)
	assert.Nil(t, err, "We are expecting no error and got one")
	assert.Equal(t, testDB.PointsCount(), int64(1))
}

func Test_WritePoints(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)
	ti, _ := ConvertTimeStamp("23:55:28,13-MAY-2015", "Europe/Paris")
	testDB.AddPoint("test", ti, fields, tags)
	testDB.AddPoint("test2", ti, fields, tags)
	ti, _ = ConvertTimeStamp("23:56:28,13-MAY-2015", "Europe/Paris")
	testDB.AddPoint("test", ti, fields2, tags)
	testDB.AddPoint("test2", ti, fields2, tags)
	err = testDB.WritePoints()
	assert.Nil(t, err, "We are expecting no errors and got one")
}

func Test_ReadPoints(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)
	filters := new(Filters)
	res, err := testDB.ReadPoints("a", filters, "test", "test2", "2015-01-19", "2016-09-19", "")
	assert.Nil(t, err, "We are expecting no errors and got one")
	assert.Equal(t, res[0].Name, "test2")
}

func Test_ReadLastPoint(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)
	filters := new(Filters)
	res, err := testDB.ReadLastPoint("a", filters, "test")
	assert.Nil(t, err, "We are expecting no errors and got one")
	assert.Equal(t, string(res), "")
	//assert.Equal(t, string(res[1].(json.Number)), "11")
}

func Test_ListMeasurement(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)
	filters := new(Filters)
	res, err := testDB.ListMeasurement(filters)
	assert.Nil(t, err, "We are expecting no errors and got one")
	assert.NotNil(t, res, "We are expecting a result and got none")

}

func Test_DropDB(t *testing.T) {
	testDB, err := NewInfluxDB(cfg)

	res, err := testDB.DropDB(testdb)

	assert.Nil(t, err, "We are expecting no error and got one")
	assert.NotNil(t, res, "We are expecting a result and got none")

}

func ConvertTimeStamp(s string, tz string) (time.Time, error) {
	var err error
	var loc *time.Location
	if len(tz) > 0 {
		loc, err = time.LoadLocation(tz)
		if err != nil {
			loc = time.FixedZone("Europe/Paris", 2*60*60)
		}
	} else {
		timezone, _ := time.Now().In(time.Local).Zone()
		loc, err = time.LoadLocation(timezone)
		if err != nil {
			loc = time.FixedZone("Europe/Paris", 2*60*60)
		}
	}

	t, err := time.ParseInLocation(timeformat, s, loc)
	return t, err
}
