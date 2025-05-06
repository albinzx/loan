package mysql

import (
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// DriverMySQL is driver name for mysql
	DriverMySQL = "mysql"
)

// DataSource is mysql data source
type DataSource struct {
	Host      string
	Port      string
	User      string
	Password  string
	Database  string
	ParseTime bool
	Location  string
	Timeout   time.Duration
}

// dsn returns mysql data source name
func (my *DataSource) dsn() string {
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", my.User, my.Password, my.Host, my.Port, my.Database)
	val := url.Values{}

	if my.ParseTime {
		val.Add("parseTime", "1")
	}
	if len(my.Location) > 0 {
		val.Add("loc", my.Location)
	}
	if my.Timeout > 0 {
		val.Add("timeout", my.Timeout.String())
	}

	if len(val) == 0 {
		return connection
	}

	return fmt.Sprintf("%s?%s", connection, val.Encode())
}

// Name returns mysql driver name and data source name
func (my *DataSource) Name() (string, string, error) {
	return DriverMySQL, my.dsn(), nil
}
