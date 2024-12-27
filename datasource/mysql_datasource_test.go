package datasource

import (
	"fmt"
	"log"
	"testing"
)

func TestMysqlDataSource(t *testing.T) {
	dsn := "root:123456@tcp(localhost:3306)/gorm_test?parseTime=True"
	mysqlDataSource := NewMysqlDataSource(dsn)
	err := mysqlDataSource.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = mysqlDataSource.Close() }()

	rows, err := mysqlDataSource.Db.Query("show tables")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var databaseName string
		if err = rows.Scan(&databaseName); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("database: %v\n", databaseName)
	}
}
