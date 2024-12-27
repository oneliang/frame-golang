package datasource

import (
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"log"
	"testing"
)

func TestClickHouseDataSource(t *testing.T) {
	dsn := "tcp://default:123456@localhost:8123/test"

	clickHouseDataSource := NewClickHouseDataSource(dsn)
	err := clickHouseDataSource.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = clickHouseDataSource.Close() }()

	rows, err := clickHouseDataSource.Db.Query("show tables")
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

	result, err := clickHouseDataSource.Db.Exec("insert into t_user(id,name) values (?,?)", "2", "李四")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.RowsAffected())
}
