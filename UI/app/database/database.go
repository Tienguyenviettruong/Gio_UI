package database

import (
	"fmt"
	"github.com/dzwvip/oracle"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDBOracle(host string, port int, user string, password string, service string) (bool, error) {
	dataSourceName := fmt.Sprintf("%s/%s@%s:%d/%s", user, password, host, port, service)
	fmt.Println(dataSourceName)
	dbConnection, err := gorm.Open(oracle.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return false, err
	}

	if err := dbConnection.Raw("select 1").Error; err != nil {
		return false, err
	}

	db = dbConnection
	return true, nil
}
