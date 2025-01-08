package model

import (
	"errors"
	"github.com/oneliang/frame-golang/query/base"
	"github.com/oneliang/frame-golang/query/mysql"
	"github.com/oneliang/util-golang/common"
)

type TransferConfig struct {
	FromTable *FromTable `xml:"fromTable"`
	ToTable   *ToTable   `xml:"toTable"`
}
type FromTable struct {
	//XMLName   xml.Name `xml:"tables"`
	MainTable  *MainTable    `xml:"mainTable"`
	SlaveTable []*SlaveTable `xml:"slaveTable"`
}

type MainTable struct {
	Name        string    `xml:"name,attr"`
	SequenceKey string    `xml:"sequenceKey,attr"`
	Columns     []*Column `xml:"column"`
}

type SlaveTable struct {
	Name      string    `xml:"name,attr"`
	MergeKeys string    `xml:"mergeKeys,attr"`
	Columns   []*Column `xml:"column"`
}
type Column struct {
	Name    string `xml:"name,attr"`
	DataKey string `xml:"dataKey,attr"`
}

type ToTable struct {
	Name      string      `xml:"name,attr"`
	Comment   string      `xml:"comment,attr"`
	ToColumns []*ToColumn `xml:"toColumn"`
}

type ToColumn struct {
	DataKey      string `xml:"dataKey,attr"`
	Name         string `xml:"name,attr"`
	DatabaseType string `xml:"databaseType,attr"`
	Length       int64  `xml:"length,attr"`
	Precision    int64  `xml:"precision,attr"`
	Nullable     bool   `xml:"nullable,attr"`
	DefaultValue string `xml:"defaultValue,attr"`
	Comment      string `xml:"comment,attr"`
}

func (this TransferConfig) GetFromTable() (*FromTable, error) {
	if this.FromTable == nil {
		return nil, errors.New("FromTable is nil or FromTable.MainTable is nil")
	}
	return this.FromTable, nil
}

func (this TransferConfig) GetMainTable() (*MainTable, error) {
	if this.FromTable == nil || this.FromTable.MainTable == nil {
		return nil, errors.New("FromTable is nil or FromTable.MainTable is nil")
	}
	return this.FromTable.MainTable, nil
}

func (this TransferConfig) GetFromMainTableColumnDataList() []*base.ColumnData {
	if this.FromTable == nil || this.FromTable.MainTable == nil {
		return make([]*base.ColumnData, 0)
	}
	if this.FromTable.MainTable.Columns == nil {
		return make([]*base.ColumnData, 0)
	}

	return common.ListToNewList[*Column, *base.ColumnData](this.FromTable.MainTable.Columns, func(index int, item *Column) *base.ColumnData {
		return &base.ColumnData{
			ColumnName: item.Name,
			DataKey:    item.DataKey,
		}
	})
}

func (this TransferConfig) GetToKeyColumnList() []*KeyToColumn {
	if this.ToTable == nil {
		return make([]*KeyToColumn, 0)
	}
	if this.ToTable.ToColumns == nil {
		return make([]*KeyToColumn, 0)
	}

	return common.ListToNewList[*ToColumn, *KeyToColumn](this.ToTable.ToColumns, func(index int, item *ToColumn) *KeyToColumn {
		return &KeyToColumn{
			Key:          item.DataKey,
			ToColumnName: item.Name,
		}
	})
}

func (this TransferConfig) GetToTableDefinition() *mysql.TableDefinition {
	if this.ToTable == nil {
		return nil
	}
	if this.ToTable.ToColumns == nil {
		return nil
	}
	columnDefinitionList := common.ListToNewList[*ToColumn, *mysql.ColumnDefinition](this.ToTable.ToColumns, func(index int, item *ToColumn) *mysql.ColumnDefinition {
		return &mysql.ColumnDefinition{
			Name:         item.Name,
			DatabaseType: item.DatabaseType,
			Length:       item.Length,
			Precision:    item.Precision,
			Nullable:     item.Nullable,
			DefaultValue: item.DefaultValue,
			Comment:      item.Comment,
		}
	})
	return &mysql.TableDefinition{
		Name:                 this.ToTable.Name,
		Comment:              this.ToTable.Comment,
		ColumnDefinitionList: columnDefinitionList,
	}

}
