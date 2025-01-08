package transformer

import (
	"github.com/oneliang/frame-golang/parallel"
	"github.com/oneliang/frame-golang/transfer/database_transfer/model"
	"github.com/oneliang/util-golang/common"
	"log"
)

type OutputTransformer struct {
	parallel.TransformProcessor[any, any]
	logger          *log.Logger
	keyToColumnList []*model.KeyToColumn
	columnList      []string
	columnIndexMap  map[string]int
	keyToColumnMap  map[string]*model.KeyToColumn
}

func NewOutputTransformer(logger *log.Logger, keyToColumnList []*model.KeyToColumn) parallel.TransformProcessor[any, any] {
	outputTransformer := &OutputTransformer{
		logger:          logger,
		keyToColumnList: keyToColumnList,
	}
	outputTransformer.columnList = common.ListToNewList[*model.KeyToColumn, string](keyToColumnList, func(index int, item *model.KeyToColumn) string {
		return item.ToColumnName
	})
	outputTransformer.columnIndexMap = common.ListToNewMap[string, string, int](outputTransformer.columnList, func(index int, item string) (string, int) {
		return item, index
	})
	outputTransformer.keyToColumnMap = model.GetKeyToColumnMap(keyToColumnList)
	return outputTransformer
}

func (this OutputTransformer) Process(inputDataList any, transformContext parallel.TransformContext[any]) {
	dataList := inputDataList.([]*map[string]interface{})

	columnListLength := len(this.columnList)

	newDataList := make([][]any, len(dataList))

	for index, dataPointer := range dataList {
		data := *dataPointer

		newDataList[index] = make([]any, columnListLength)

		for key, value := range data {
			keyToColumn, exist := this.keyToColumnMap[key]
			if !exist {
				continue
			}
			columnIndex, exist := this.columnIndexMap[keyToColumn.ToColumnName]
			if !exist {
				continue
			}
			newDataList[index][columnIndex] = value
		}

	}
	transformContext.Collect(newDataList)
}
