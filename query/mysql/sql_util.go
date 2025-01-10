package mysql

import (
	"fmt"
	"github.com/oneliang/util-golang/common"
	"github.com/oneliang/util-golang/constants"
	"strings"
)

const (
	VARCHAR  = "varchar"
	DATETIME = "datetime"
	INT      = "int"
	UINT     = "uint"
	FLOAT    = "float"
)

type TableDefinition struct {
	Name                 string
	Comment              string
	ColumnDefinitionList []*ColumnDefinition
}
type ColumnDefinition struct {
	Name         string
	DatabaseType string
	Length       int64
	Precision    int64
	Nullable     bool
	DefaultValue string
	Comment      string
}

// GenerateCreateTableSql .
func GenerateCreateTableSql(tableDefinition *TableDefinition) string {
	tableName := tableDefinition.Name
	columnDefinitionList := tableDefinition.ColumnDefinitionList

	if common.StringIsBlank(tableName) || len(columnDefinitionList) <= 0 {
		return constants.STRING_BLANK
	}
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("CREATE TABLE `")
	sqlBuilder.WriteString(tableName)
	sqlBuilder.WriteString("` (")
	sqlBuilder.WriteString("`id` int unsigned NOT NULL AUTO_INCREMENT,")
	for _, columnDefinition := range columnDefinitionList {
		columnName := columnDefinition.Name
		columnNullable := columnDefinition.Nullable
		columnDefaultValue := columnDefinition.DefaultValue
		var columnDefaultValueString string

		var columnType string
		var columnTempDefaultValue string
		switch columnDefinition.DatabaseType {
		case VARCHAR:
			columnTempDefaultValue = fmt.Sprintf("'%s'", columnDefaultValue)
			columnType = fmt.Sprintf("varchar(%d)", columnDefinition.Length)
			break
		case DATETIME:
			columnTempDefaultValue = fmt.Sprintf("'%s'", columnDefaultValue)
			columnType = fmt.Sprintf("datetime(%d)", columnDefinition.Length)
			break
		case INT:
			if common.StringIsBlank(columnDefaultValue) {
				columnTempDefaultValue = "0"
			} else {
				columnTempDefaultValue = fmt.Sprintf("%s", columnDefaultValue)
			}
			columnType = fmt.Sprintf("int(%d)", columnDefinition.Length)
			break
		case UINT:
			if common.StringIsBlank(columnDefaultValue) {
				columnTempDefaultValue = "0"
			} else {
				columnTempDefaultValue = fmt.Sprintf("%s", columnDefaultValue)
			}
			columnType = fmt.Sprintf("int(%d) unsigned", columnDefinition.Length)
			break
		case FLOAT:
			if common.StringIsBlank(columnDefaultValue) {
				columnTempDefaultValue = "0"
			} else {
				columnTempDefaultValue = fmt.Sprintf("%s", columnDefaultValue)
			}
			columnType = FLOAT
			break
		}

		var columnNullableString = constants.STRING_BLANK
		if !columnNullable {
			columnNullableString = "NOT NULL"
			columnDefaultValueString = fmt.Sprintf("DEFAULT %s", columnTempDefaultValue)
		}

		columnComment := columnDefinition.Comment
		sqlBuilder.WriteString(fmt.Sprintf("`%s` %s %s %s COMMENT '%s'", columnName, columnType, columnNullableString, columnDefaultValueString, columnComment))
		sqlBuilder.WriteString(",")
	}
	sqlBuilder.WriteString("PRIMARY KEY (`id`)")
	sqlBuilder.WriteString(" )")
	sqlBuilder.WriteString(fmt.Sprintf(" COMMENT='%s'", tableDefinition.Comment))
	sqlBuilder.WriteString(" ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci")
	return sqlBuilder.String()
}

// GenerateInsertSql .
func GenerateInsertSql(tableName string, columnList []string) string {
	if len(columnList) <= 0 {
		return constants.STRING_BLANK
	}
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("INSERT INTO `")
	sqlBuilder.WriteString(tableName)
	sqlBuilder.WriteString("` ")
	sqlBuilder.WriteString("(")
	sqlBuilder.WriteString(common.ListJoinToString[string](columnList, func(index int, item string) string {
		return item
	}, constants.SYMBOL_COMMA+constants.STRING_SPACE))
	sqlBuilder.WriteString(")")
	sqlBuilder.WriteString(" VALUES ")
	sqlBuilder.WriteString("(")
	sqlBuilder.WriteString(common.ListJoinToString[string](columnList, func(index int, item string) string {
		return constants.SYMBOL_QUESTION_MARK
	}, constants.SYMBOL_COMMA+constants.STRING_SPACE))
	sqlBuilder.WriteString(")")
	return sqlBuilder.String()
}
