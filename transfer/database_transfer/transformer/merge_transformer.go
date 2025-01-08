package transformer

import (
	"github.com/oneliang/frame-golang/parallel"
	"github.com/oneliang/util-golang/common"
	"github.com/oneliang/util-golang/constants"
	"log"
	"strconv"
)

type MergeTransformer struct {
	parallel.TransformProcessor[any, any]
	logger *log.Logger
}

var emptyOutputList = make([]*map[string]any, 0)

func generateKey(data map[string]any, mergeKeys []string) string {
	return common.ListJoinToString[string](mergeKeys, func(index int, item string) string {
		value, exist := data[item]
		if !exist {
			value = ""
		}
		switch value.(type) {
		case int64:
			return strconv.FormatInt(value.(int64), 10)
		}
		return value.(string)
	}, constants.SYMBOL_COMMA)
}

func NewMergeTransformer(logger *log.Logger) parallel.TransformProcessor[any, any] {
	return &MergeTransformer{
		logger: logger,
	}
}

func (this MergeTransformer) Process(needToMergeData any, transformContext parallel.TransformContext[any]) {
	if needToMergeData == nil {
		transformContext.Collect(emptyOutputList)
		return
	}

	mergerConfig := needToMergeData.(*common.MergerConfig)

	dataList := common.Merge(mergerConfig)
	transformContext.Collect(dataList)
	return
}
