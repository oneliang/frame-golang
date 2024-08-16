package base

import (
	"database/sql"
	"log"
)

type Query struct {
	logger       *log.Logger
	databaseName string
	db           *sql.DB
}

func NewQuery(logger *log.Logger, databaseName string, db *sql.DB) *Query {
	return &Query{
		logger:       logger,
		databaseName: databaseName,
		db:           db,
	}
}
func (this Query) TableExists(tableName string) bool {
	var count int
	querySql := "SELECT COUNT(0) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?"

	err := this.QueryRow(querySql, []any{this.databaseName, tableName}, []any{&count})
	if err != nil {
		this.logger.Printf("Error checking if table exists: %v", err)
	}
	return count > 0
}

// QueryRow .
func (this Query) QueryRow(querySql string, args []any, result []any) error {
	this.logger.Printf("Query row sql: %s", querySql)
	err := this.db.QueryRow(querySql, args...).Scan(result...)
	if err != nil {
		this.logger.Printf("Error in QueryRow %v", err)
		return err
	}
	return nil
}

// Query .
func (this Query) Query(querySql string, args []any, dataCallback func(rows *sql.Rows)) error {
	this.logger.Printf("Query sql: %s", querySql)
	rows, err := this.db.Query(querySql, args...)
	if err != nil {
		this.logger.Printf("Error in query: %v, sql: %s", err, querySql)
		return err
	}
	for rows.Next() {
		dataCallback(rows)
	}
	return nil
}

// Exec for insert or update
func (this Query) Exec(execSql string, args []any) (sql.Result, error) {
	this.logger.Printf("Exec sql: %s", execSql)
	result, err := this.db.Exec(execSql, args...)
	if err != nil {
		this.logger.Printf("Error in exec. sql: %s, err: %v", execSql, err)
		return nil, err
	}
	return result, nil
}

func (this Query) InsertBatch(sql string, dataList [][]any) error {
	// begin transaction
	tx, err := this.db.Begin()
	if err != nil {
		this.logger.Printf("Error in begin transaction, sql:%s, err:%v", sql, err)
		return err
	}

	//
	stmt, err := tx.Prepare(sql)
	if err != nil {
		this.logger.Printf("Error in prepare sql, sql:%s, err:%v", sql, err)
		return err
	}
	defer func() { _ = stmt.Close() }()

	// insert batch
	for index, value := range dataList {
		if _, stmtErr := stmt.Exec(value...); stmtErr != nil {
			this.logger.Printf("Error in statement exec, sql:%s, data index:%d, err:%v", sql, index, stmtErr)
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				this.logger.Printf("Error in transaction rollback, sql:%s, data index:%d, err:%v", sql, index, rollbackError)
				return rollbackError
			} else {
				return stmtErr
			}
		}
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		this.logger.Printf("Error in transaction commit, sql:%s, err:%v", sql, err)
		return err
	}
	return nil
}
