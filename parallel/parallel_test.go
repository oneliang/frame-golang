package parallel

import (
	"fmt"
	"github.com/oneliang/util-golang/common"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

type defaultSourceProcessor struct {
}

func (this *defaultSourceProcessor) Process(sourceContext SourceContext[any]) {
	for i := 0; i < 3; i++ {
		sourceContext.Collect(fmt.Sprintf("%s_%s", "test_string", common.GenerateZeroString(i, 4)), CONTEXT_ACTION_NONE)
	}
}

type defaultTransformProcessor struct {
}

func (this *defaultTransformProcessor) Process(value any, transformContext TransformContext[any]) {
	//fmt.Println(fmt.Sprintf("goroutine id:%v, value:%s", GetGoroutineId(), value))
	transformContext.Collect(fmt.Sprintf("%s_%s", value, "transform"))
}

type defaultSinkProcessor struct {
	valueQueue chan string
	count      int
}

func newDefaultSinkProcessor() *defaultSinkProcessor {
	sinkProcessor := &defaultSinkProcessor{
		valueQueue: make(chan string),
		count:      0,
	}
	go func() {
		//for {
		//	time.After(10000)
		//	select {
		//	case item, ok := <-sinkProcessor.valueQueue:
		//		if ok {
		//			fmt.Println(fmt.Sprintf("item:%s", item))
		//
		//			sinkProcessor.count++
		//		}
		//	}
		//}
	}()
	return sinkProcessor
}
func (this *defaultSinkProcessor) Sink(value any) {
	fmt.Println(fmt.Sprintf("goroutine id:%v, value:%s", common.GetGoroutineId(), value))
	this.valueQueue <- fmt.Sprintf("%s", value)
	fmt.Println(fmt.Sprintf("add value:%s", value))
}

func TestTransfer(t *testing.T) {
	jobConfiguration := NewJobConfiguration(true, false, 4)
	sourceProcessor := &defaultSourceProcessor{}
	job := NewJob("job", jobConfiguration)
	job.AddSourceProcessor(sourceProcessor)
	transformProcessor := &defaultTransformProcessor{}
	sinkProcessor := newDefaultSinkProcessor()
	job.GenerateFirstJobStep().AddTransformProcessor(transformProcessor).AddSinkProcessor(sinkProcessor)
	job.Execute()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Println(fmt.Sprintf("data count:%d", sinkProcessor.count))
}
