package test

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

type ItemCount struct {
	lock sync.Mutex
	//AllDataMap    map[string]Item           // "%d_%s_%s_%d_%d", schoolId, scene, exerciseId, studentId, answerId
	//SceneDataMap  map[string]*StatusCounter // ["global_%s", scene]*StatusCounter
	GlobalCount StatusCount
}
type StatusCount struct {
	UnFinishedCount int32 // 未完成批改/当前队列中有xx题
	FinishedCount   int32 // 已完成批改
	//TotalCount      int32 // 进题总数
	lock *sync.Mutex
}

func (this *StatusCount) IncreaseCount(count int32) {
	this.lock.Lock()
	this.UnFinishedCount += count
	//this.TotalCount++
	fmt.Println(fmt.Sprintf("CountLog: IncreaseCount : %v", count))
	this.lock.Unlock()
}

// 题目完成批改
func (this *StatusCount) FinishedCorrect(count int32) {
	this.lock.Lock()
	this.FinishedCount += count
	// 防止未批改数小于0
	this.UnFinishedCount -= count
	fmt.Println(fmt.Sprintf("CountLog: FinishedCorrect : %v", count))
	//if this.UnFinishedCount < 0 {
	//	this.UnFinishedCount = 0
	//}
	this.lock.Unlock()
}

func TestCalculate(t *testing.T) {
	calculateTest()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem) // 读取内存统计信息到变量mem

	// 输出不同部分的内存信息
	fmt.Printf("Alloc = %v\n", mem.Alloc)           // 已分配但未释放的字节数
	fmt.Printf("TotalAlloc = %v\n", mem.TotalAlloc) // 从开始运行到现在为止所有对象（包括GC）分配的字节数
	fmt.Printf("Sys = %v\n", mem.Sys)               // 系统调用分配的字节数
	fmt.Printf("NumGC = %v\n", mem.NumGC)           // GC发生次数

}

func calculateTest() {
	itemCount := ItemCount{
		GlobalCount: StatusCount{
			UnFinishedCount: 0,
			FinishedCount:   0,
			lock:            &sync.Mutex{},
		},
	}
	for i := 0; i < 1000; i++ {
		go func(i int32) {
			itemCount.GlobalCount.IncreaseCount(i)
		}(int32(i))
	}
	for i := 0; i < 1000; i++ {
		go func(i int32) {
			itemCount.GlobalCount.FinishedCorrect(i)
		}(int32(i))
	}
	time.Sleep(10000 * time.Millisecond)
	a, _ := json.Marshal(itemCount)

	fmt.Println(string(a))
}
