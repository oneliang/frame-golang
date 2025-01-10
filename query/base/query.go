package base

import (
	"database/sql"
	"errors"
	"log"
)

type Query struct {
	logger       *log.Logger
	databaseName string
	db           *sql.DB
}

type ColumnData struct {
	ColumnName string //to column name
	DataKey    string //
}

func NewQuery(logger *log.Logger, databaseName string, db *sql.DB) *Query {
	return &Query{
		logger:       logger,
		databaseName: databaseName,
		db:           db,
	}
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
func (this Query) Query(querySql string, args []any, dataCallback func(rows *sql.Rows) error) error {
	this.logger.Printf("Query begin, sql: %s", querySql)
	rows, err := this.db.Query(querySql, args...)
	if err != nil {
		this.logger.Printf("Error in query: %v, sql: %s", err, querySql)
		return err
	}
	err = dataCallback(rows)
	this.logger.Printf("Query end, sql: %s", querySql)
	if err != nil {
		return err
	}
	return nil
}

// QueryWithEachRow .
func (this Query) QueryWithEachRow(querySql string, args []any, eachRowCallback func(rows *sql.Rows)) error {
	return this.Query(querySql, args, func(rows *sql.Rows) error {
		for rows.Next() {
			eachRowCallback(rows)
		}
		return nil
	})
}
func (this Query) QueryForMap(querySql string, args []any, columnDataMap map[string]*ColumnData) ([]*map[string]any, error) {
	return this.QueryForMap(querySql, args, columnDataMap)
}

// QueryForMap .
func (this Query) QueryForMapWithRowProcessor(querySql string, args []any, columnDataMap map[string]*ColumnData, rowProcessor func(rowDataMap *map[string]any)) ([]*map[string]any, error) {
	dataList := make([]*map[string]any, 0)
	if columnDataMap == nil || len(columnDataMap) <= 0 {
		return dataList, errors.New("columnDataMap is nil or empty")
	}
	this.logger.Printf("QueryForMap begin, sql: %s", querySql)

	err := this.Query(querySql, args, func(rows *sql.Rows) error {
		columns, err := rows.Columns()
		if err != nil {
			return err
		}

		for rows.Next() {
			rowData := make([]any, len(columns))
			rowDataPointers := make([]any, len(columns))
			for i := range rowData {
				rowDataPointers[i] = &rowData[i]
			}

			if err := rows.Scan(rowDataPointers...); err != nil {
				this.logger.Printf("Error in scan, sql:%s, err:%v", querySql, err)
				return err
			}
			rowDataMap := make(map[string]any)
			for i, column := range columns {
				columnData, exist := columnDataMap[column]
				if !exist {
					continue
				}
				//exist, fix type
				switch rowData[i].(type) {
				case []uint8:
					rowData[i] = string(rowData[i].([]uint8)) //convert to string
					break
				}
				//set to map with key
				rowDataMap[columnData.DataKey] = rowData[i]
			}
			if rowProcessor != nil {
				rowProcessor(&rowDataMap)
			}

			dataList = append(dataList, &rowDataMap)
		}
		return nil
	})

	this.logger.Printf("QueryForMap end[rows:%d], sql: %s", len(dataList), querySql)

	if err != nil {
		return nil, err
	}
	return dataList, nil
}

// Exec for insert or update
func (this Query) Exec(execSql string, args []any) (sql.Result, error) {
	this.logger.Printf("Exec begin, sql: %s", execSql)
	result, err := this.db.Exec(execSql, args...)
	if err != nil {
		this.logger.Printf("Error in exec. sql: %s, err: %v", execSql, err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		this.logger.Printf("Error in RowsAffected() after exec. sql: %s, err: %v", execSql, err)
		return nil, err
	}
	this.logger.Printf("Exec end[%d], sql: %s", rowsAffected, execSql)

	return result, nil
}

func (this Query) InsertBatch(insertSql string, dataList [][]any) error {
	this.logger.Printf("InsertBatch inputDataLength:%d, sql: %s", len(dataList), insertSql)
	// begin transaction
	tx, err := this.db.Begin()
	if err != nil {
		this.logger.Printf("Error in begin transaction, sql:%s, err:%v", insertSql, err)
		return err
	}

	//
	stmt, err := tx.Prepare(insertSql)
	if err != nil {
		this.logger.Printf("Error in prepare sql, sql:%s, err:%v", insertSql, err)
		return err
	}
	defer func() { _ = stmt.Close() }()

	// insert batch
	for index, value := range dataList {
		if _, stmtErr := stmt.Exec(value...); stmtErr != nil {
			this.logger.Printf("Error in statement exec, sql:%s, data index:%d, err:%v", insertSql, index, stmtErr)
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				this.logger.Printf("Error in transaction rollback, sql:%s, data index:%d, err:%v", insertSql, index, rollbackError)
				return rollbackError
			} else {
				return stmtErr
			}
		}
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		this.logger.Printf("Error in transaction commit, sql:%s, err:%v", insertSql, err)
		return err
	}
	return nil
}
