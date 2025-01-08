package database_transfer

import (
	"fmt"
	"github.com/oneliang/frame-golang/datasource"
	"github.com/oneliang/frame-golang/parallel"
	"github.com/oneliang/frame-golang/query/base"
	"github.com/oneliang/frame-golang/query/mysql"
	"github.com/oneliang/frame-golang/transfer/database_transfer/model"
	"github.com/oneliang/frame-golang/transfer/database_transfer/transformer"
	"github.com/oneliang/util-golang/common"
	"io"
	"log"
	"os"
	"testing"
)

const (
	FROM_DSN   = "root:123456@tcp(localhost:3306)/gorm_test?parseTime=True&charset=utf8mb4"
	TO_DSN     = "root:123456@tcp(localhost:3306)/gorm_test?parseTime=True&charset=utf8mb4"
	consoleLog = true
)

type defaultSourceProcessor struct {
	parallel.SourceProcessor[any]
	logger                  *log.Logger
	transferConfig          model.TransferConfig
	fromQuery               *base.Query
	toQuery                 *mysql.MySqlQuery
	mainTableColumnDataList []*base.ColumnData
	mainTableColumnDataMap  map[string]*base.ColumnData
}

func newDefaultSourceProcessor(logger *log.Logger, transferConfig model.TransferConfig, fromQuery *base.Query, toQuery *mysql.MySqlQuery) parallel.SourceProcessor[any] {
	defaultSourceProcessorPointer := &defaultSourceProcessor{
		transferConfig: transferConfig,
		logger:         logger,
		fromQuery:      fromQuery,
		toQuery:        toQuery,
	}

	defaultSourceProcessorPointer.mainTableColumnDataList = transferConfig.GetFromMainTableColumnDataList()
	defaultSourceProcessorPointer.mainTableColumnDataMap = model.GetColumnDataMap(defaultSourceProcessorPointer.mainTableColumnDataList)
	return defaultSourceProcessorPointer
}

func (this *defaultSourceProcessor) getData(lastId int64, pageSize int) (mainData []*map[string]any, slaveDataList []*common.SlaveData, err error) {

	querySql := fmt.Sprintf("SELECT id, name, project_id, modify_time, create_time FROM %s WHERE id > ? ORDER BY id ASC LIMIT ?", "t_module")
	dataList, err := this.fromQuery.QueryForMap(querySql, []any{lastId, pageSize}, this.mainTableColumnDataMap)
	if err != nil {
		this.logger.Printf("err:%v", err)
		return nil, nil, err
	}
	return dataList, nil, nil
}

func (this *defaultSourceProcessor) Process(sourceContext parallel.SourceContext[any]) {

	mainTable, err := this.transferConfig.GetMainTable()
	if err != nil {
		this.logger.Printf("err:%v", err)
		return
	}

	var id int64 = 0
	pageSize := 1
	for {
		dataList, slaveDataList, err := this.getData(id, pageSize)
		if err != nil {
			break
		}
		//jsonBytes, _ := json.Marshal(dataList)
		//fmt.Println(string(jsonBytes))

		if err != nil {
			this.logger.Printf("err:%v", err)
		}

		//check data
		var dataListLength = len(dataList)
		this.logger.Printf("data list len:%d, data > id:%d", dataListLength, id)
		if dataListLength == 0 {
			sourceContext.Collect(nil, parallel.CONTEXT_ACTION_FINISHED)
			break
		} else {
			dataMap := *dataList[dataListLength-1]
			id = dataMap[mainTable.SequenceKey].(int64)
		}

		//merge data
		mergeData := &common.MergerConfig{
			MasterDataList: dataList,
			SlaveDataList:  slaveDataList, //[][]*map[string]any{dataList},
			StaticDataList: []*map[string]any{staticData},
		}
		sourceContext.Collect(mergeData, parallel.CONTEXT_ACTION_NONE)
	}
}

var staticData = &map[string]any{
	"SCHOOL_ID":   125,
	"SCHOOL_NAME": "北海中学",
}

type defaultTransformProcessor struct {
}

func (this *defaultTransformProcessor) Process(value any, transformContext parallel.TransformContext[any]) {
	//fmt.Println(fmt.Sprintf("goroutine id:%v, value:%s", GetGoroutineId(), value))
	transformContext.Collect(fmt.Sprintf("%s_%s", value, "transform"))
}

type defaultSinkProcessor struct {
	parallel.SinkProcessor[any]
	logger         *log.Logger
	transferConfig model.TransferConfig
	toQuery        *mysql.MySqlQuery
	insertSql      string
}

func newDefaultSinkProcessor(logger *log.Logger, transferConfig model.TransferConfig, toQuery *mysql.MySqlQuery) parallel.SinkProcessor[any] {
	sinkProcessor := &defaultSinkProcessor{
		logger:         logger,
		transferConfig: transferConfig,
		toQuery:        toQuery,
	}

	tableDefinition := transferConfig.GetToTableDefinition()
	//initialize local variable
	toQuery.DropTableExists(tableDefinition.Name)
	createTableSql := mysql.GenerateCreateTableSql(tableDefinition)
	_, err := toQuery.Exec(createTableSql, nil)
	if err != nil {
		logger.Printf("err:%v", err)
	}

	keyToColumnList := transferConfig.GetToKeyColumnList()
	columnList := common.ListToNewList[*model.KeyToColumn, string](keyToColumnList, func(index int, item *model.KeyToColumn) string {
		return item.ToColumnName
	})
	sinkProcessor.insertSql = mysql.GenerateInsertSql(tableDefinition.Name, columnList)
	return sinkProcessor
}
func (this *defaultSinkProcessor) Sink(inputDataList any) {
	dataList := inputDataList.([][]any)

	//jsonBytes, _ := json.Marshal(dataList)
	//fmt.Println(string(jsonBytes))

	err := this.toQuery.InsertBatch(this.insertSql, dataList)
	if err != nil {
		this.logger.Printf("err:%v", err)
	}
}

func TestTransfer(t *testing.T) {
	//logger
	var writer io.Writer
	if consoleLog {
		writer = os.Stderr
	} else {
		logFile, err := os.OpenFile("database_transfer_default.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer func() { _ = logFile.Close() }()
		writer = logFile
	}
	logger := log.New(writer, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	//read config
	transferConfig := model.TransferConfig{}
	err := common.LoadXmlToObject("config/transfer.xml", &transferConfig)
	if err != nil {
		logger.Printf("err:%v", err)
		return
	}
	keyToColumnList := transferConfig.GetToKeyColumnList()

	//from database
	fromMysqlDataSource := datasource.NewMysqlDataSource(FROM_DSN)
	err = fromMysqlDataSource.Open()
	if err != nil {
		logger.Fatal(err)
	}
	defer func() { _ = fromMysqlDataSource.Close() }()

	//to database
	toMysqlDataSource := datasource.NewMysqlDataSource(TO_DSN)
	err = toMysqlDataSource.Open()
	if err != nil {
		logger.Fatal(err)
	}
	defer func() { _ = toMysqlDataSource.Close() }()

	fromQuery := base.NewQuery(logger, "from_datasource", fromMysqlDataSource.Db)
	toQuery := mysql.NewMysqlQuery(logger, "to_datasource", toMysqlDataSource.Db)

	//parallel job
	jobConfiguration := parallel.NewJobConfiguration(false, false)

	sourceProcessor := newDefaultSourceProcessor(logger, transferConfig, fromQuery, toQuery)
	job := parallel.NewJob("job", jobConfiguration)
	job.AddSourceProcessor(sourceProcessor)
	//transformProcessor := &defaultTransformProcessor{}
	mergeTransformer := transformer.NewMergeTransformer(logger)
	outputTransformer := transformer.NewOutputTransformer(logger, keyToColumnList)

	sinkProcessor := newDefaultSinkProcessor(logger, transferConfig, toQuery)

	job.GenerateFirstJobStep().
		AddTransformProcessor(mergeTransformer).
		AddTransformProcessor(outputTransformer).
		AddSinkProcessor(sinkProcessor)
	job.Execute()
	//sigChan := make(chan os.Signal, 1)
	//signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	//<-sigChan
}
