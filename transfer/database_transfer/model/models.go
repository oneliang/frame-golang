package model

import (
	"fmt"
	"github.com/oneliang/frame-golang/query/base"
	"github.com/oneliang/util-golang/common"
	"github.com/oneliang/util-golang/constants"
	"time"
)

type KeyToColumn struct {
	Key          string
	ToColumnName string //to column name
}

func GetColumnDataMap(columnToKeyList []*base.ColumnData) map[string]*base.ColumnData {
	if columnToKeyList == nil {
		return make(map[string]*base.ColumnData)
	}
	return common.ListToMap[*base.ColumnData, string](columnToKeyList, func(index int, item *base.ColumnData) string {
		return item.ColumnName
	})
}

func GetKeyToColumnMap(columnToKeyList []*KeyToColumn) map[string]*KeyToColumn {
	if columnToKeyList == nil {
		return make(map[string]*KeyToColumn)
	}
	return common.ListToMap[*KeyToColumn, string](columnToKeyList, func(index int, item *KeyToColumn) string {
		return item.Key
	})
}

func GetBeginTime(date string) string {
	return fmt.Sprintf("%s 00:00:00", date)
}

func GetEndTime(date string) string {
	return fmt.Sprintf("%s 23:59:59.999", date)
}

func GetEndTimeWithoutMillisecond(date string) string {
	return fmt.Sprintf("%s 23:59:59", date)
}

var cstSh, _ = time.LoadLocation("Asia/Shanghai")

func GetTimeNext(date string, offset int) string {
	locationTime, _ := time.ParseInLocation(constants.TIME_LAYOUT_YEAR_MONTH_DAY, date, cstSh)
	dayZeroTimeMillis := common.GetDayZeroTimeNext(locationTime.UnixMilli(), offset)
	return time.UnixMilli(dayZeroTimeMillis).Format(constants.TIME_LAYOUT_YEAR_MONTH_DAY)
}
