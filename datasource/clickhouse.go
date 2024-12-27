package datasource

import (
	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseDataSource struct {
	dataSourceName string
	Db             *sql.DB
}

func NewClickHouseDataSource(dataSourceName string) *ClickHouseDataSource {
	return &ClickHouseDataSource{
		dataSourceName: dataSourceName,
	}
}

func (this *ClickHouseDataSource) Open() error {
	db, err := sql.Open("clickhouse", this.dataSourceName)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	this.Db = db
	return nil
}

func (this *ClickHouseDataSource) Close() error {
	return this.Db.Close()
}
