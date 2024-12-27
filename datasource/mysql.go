package datasource

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlDataSource struct {
	dataSourceName string
	Db             *sql.DB
}

func NewMysqlDataSource(dataSourceName string) *MysqlDataSource {
	return &MysqlDataSource{
		dataSourceName: dataSourceName,
	}
}

func (this *MysqlDataSource) Open() error {
	db, err := sql.Open("mysql", this.dataSourceName)
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

func (this *MysqlDataSource) Close() error {
	return this.Db.Close()
}
