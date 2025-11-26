package test

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

type Item struct {
	SchoolId      int32
	Scene         string
	ExerciseId    string
	StudentId     int32
	AnswerId      int32
	CorrectStatus int8
}

type StatusCounter struct {
	UnFinishedCount int32
	FinishedCount   int32
	TotalCount      int32
	lock            *sync.Mutex
}

func NewStatusCounter() *StatusCounter {
	return &StatusCounter{
		UnFinishedCount: 0,
		FinishedCount:   0,
		TotalCount:      0,
		lock:            &sync.Mutex{},
	}
}

func (this *StatusCounter) IncreaseCount() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.UnFinishedCount++
	this.TotalCount++
}

func (this *StatusCounter) FinishedOne() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.FinishedCount++
	this.UnFinishedCount--
}

func TestBigMapCalculate(t *testing.T) {
	calculate()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem) // 读取内存统计信息到变量mem

	// 输出不同部分的内存信息
	fmt.Printf("Alloc = %v\n", mem.Alloc)           // 已分配但未释放的字节数
	fmt.Printf("TotalAlloc = %v\n", mem.TotalAlloc) // 从开始运行到现在为止所有对象（包括GC）分配的字节数
	fmt.Printf("Sys = %v\n", mem.Sys)               // 系统调用分配的字节数
	fmt.Printf("NumGC = %v\n", mem.NumGC)           // GC发生次数

}

func calculate() {
	itemCounter := NewItemCounter()
	for i := 0; i < 1000000; i++ {
		schoolId := int32(i / 1000)
		exerciseId := strconv.Itoa(i / 100)
		studentId := int32(i / 10)
		answerId := int32(i)
		go func() {
			fmt.Println(fmt.Sprintf("schoolId:%d, exerciseId:%s, studentId:%d, answerId:%d", schoolId, exerciseId, studentId, answerId))
			itemCounter.AddItem(schoolId, "homework", exerciseId, studentId, answerId)
		}()
	}
	for i := 0; i < 1000; i++ {
		schoolId := int32(i / 1000)
		exerciseId := strconv.Itoa(i / 100)
		studentId := int32(i / 10)
		answerId := int32(i)
		go func() {
			//fmt.Println(fmt.Sprintf("schoolId:%d, exerciseId:%s, studentId:%d, answerId:%d", schoolId, exerciseId, studentId, answerId))
			itemCounter.DeleteItem(schoolId, "homework", exerciseId, studentId, answerId)
		}()
	}
	time.Sleep(1000 * time.Millisecond)
	itemCounter.Print()
	itemCounter.Write()
}

type ItemCounter struct {
	lock           *sync.Mutex
	AllDataMap     *map[string]Item           `json:"-"`
	GlobalSceneMap *map[string]*StatusCounter `json:"globalSceneMap"`
}

func NewItemCounter() ItemCounter {
	lock := sync.Mutex{}
	allDataMap := make(map[string]Item)
	globalSceneMap := make(map[string]*StatusCounter)
	return ItemCounter{
		lock:           &lock,
		AllDataMap:     &allDataMap,
		GlobalSceneMap: &globalSceneMap,
	}
}

func (this *ItemCounter) generateSchoolSceneKey(schoolId int32, scene string) string {
	return fmt.Sprintf("%d_%s", schoolId, scene)
}
func (this *ItemCounter) generateGlobalKey() string {
	return "global"
}

func (this *ItemCounter) generateGlobalSceneKey(scene string) string {
	return fmt.Sprintf("global_%s", scene)
}

func (this *ItemCounter) generateItemKey(schoolId int32, scene string, exerciseId string, studentId int32, answerId int32) string {
	return fmt.Sprintf("%d_%s_%s_%d_%d", schoolId, scene, exerciseId, studentId, answerId)
}

func (this *ItemCounter) calculateGlobalScene(scene string) {
	globalSceneKey := this.generateGlobalSceneKey(scene)
	statusCounter, exist := (*this.GlobalSceneMap)[globalSceneKey]
	if !exist {
		statusCounter = NewStatusCounter()
		(*this.GlobalSceneMap)[globalSceneKey] = statusCounter
	}
	statusCounter.IncreaseCount()
}

func (this *ItemCounter) calculateItem(schoolId int32, scene string, exerciseId string, studentId int32, answerId int32) {
	itemKey := this.generateItemKey(schoolId, scene, exerciseId, studentId, answerId)
	(*this.AllDataMap)[itemKey] = Item{
		SchoolId:      schoolId,
		Scene:         scene,
		ExerciseId:    exerciseId,
		StudentId:     studentId,
		AnswerId:      answerId,
		CorrectStatus: 0,
	}
}

// AddItem
func (this *ItemCounter) AddItem(schoolId int32, scene string, exerciseId string, studentId int32, answerId int32) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.calculateGlobalScene(scene)
	this.calculateItem(schoolId, scene, exerciseId, studentId, answerId)
}

// GetCountItem
func (this *ItemCounter) GetCountItem(scene string) StatusCounter {
	globalSceneKey := this.generateGlobalSceneKey(scene)
	statusCounter, exist := (*this.GlobalSceneMap)[globalSceneKey]
	if !exist {
		statusCounter = NewStatusCounter()
	}
	return *statusCounter
}

// DeleteItem
func (this *ItemCounter) DeleteItem(schoolId int32, scene string, exerciseId string, studentId int32, answerId int32,
) {
	// global scene
	this.lock.Lock()
	defer this.lock.Unlock()
	key := this.generateItemKey(schoolId, scene, exerciseId, studentId, answerId)
	delete(*this.AllDataMap, key)

	globalSceneKey := this.generateGlobalSceneKey(scene)
	statusCounter, exist := (*this.GlobalSceneMap)[globalSceneKey]
	if exist {
		statusCounter.FinishedOne()
	}
}

func (this *ItemCounter) Print() {
	for globalSceneKey, statusCounter := range *this.GlobalSceneMap {
		fmt.Println(fmt.Sprintf("global scene key:%s, status counter:%v", globalSceneKey, *statusCounter))
	}
}

func (this *ItemCounter) Write() {
	jsonBytes, _ := json.Marshal(this)
	fmt.Println(string(jsonBytes))
	err := os.WriteFile("export.json", jsonBytes, 0777)
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
}
