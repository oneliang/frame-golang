package mysql

import (
	"database/sql"
	"fmt"
	"github.com/oneliang/frame-golang/query/base"
	"log"
)

type MySqlQuery struct {
	*base.Query
	logger       *log.Logger
	databaseName string
	db           *sql.DB
}

func NewMysqlQuery(logger *log.Logger, databaseName string, db *sql.DB) *MySqlQuery {
	return &MySqlQuery{
		Query:        base.NewQuery(logger, databaseName, db),
		logger:       logger,
		databaseName: databaseName,
		db:           db,
	}
}

func (this MySqlQuery) TableExists(tableName string) bool {
	var count int
	querySql := "SELECT COUNT(0) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?"

	err := this.QueryRow(querySql, []any{this.databaseName, tableName}, []any{&count})
	if err != nil {
		this.logger.Printf("Error checking if table exists: %v", err)
	}
	return count > 0
}

func (this MySqlQuery) DropTableExists(tableName string) bool {
	var count int
	execSql := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)

	_, err := this.Exec(execSql, []any{})
	if err != nil {
		this.logger.Printf("Error checking if table exists: %v", err)
	}
	return count > 0
}
