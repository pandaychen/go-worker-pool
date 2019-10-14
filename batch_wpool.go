package main

import (
	"errors"
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type TaskInput struct {
	Id   string
	Guid string
	InputData interface{}
}

type TaskResult struct {
	Id         string //回调ID
	Guid       string
	OutputData interface{}
	Err        error
	Boolret    bool
}

var ERROR_PARAM_ERROR = errors.New("Pool len illegal")

func DoFixedSizeWorks(concurrency int, TaskInputList []TaskInput, fun func(task TaskInput) (id, guid string, outputData interface{}, err error, boolret bool)) ([]TaskResult, error) {
	//fmt.Println(TaskInputList)
	ResultList := make([]TaskResult, 0)
	taskLen := len(TaskInputList)
	if taskLen == 0 {
		return nil, ERROR_PARAM_ERROR
	}
	fmt.Println(taskLen)
	chInputDataList := make(chan TaskInput, taskLen)
	for _, task := range TaskInputList {
		//fmt.Println("fff",task)
		chInputDataList <- task
	}

	chResultList := make(chan TaskResult, taskLen)
	var lock sync.RWMutex

	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				lock.Lock()
				if len(chInputDataList) == 0 {
					lock.Unlock()
					break
				}
				taskData := <-chInputDataList
				lock.Unlock()
				fmt.Println("???", taskData)
				id, guid, outputData, err, boolret := fun(taskData)
				fmt.Println(id, guid, outputData, err, boolret)
				chResultList <- TaskResult{
					Id:         id,
					Guid:       guid,
					Err:        err,
					OutputData: outputData, Boolret: boolret}
			}
		}()
	}

	//get result
	for i := 0; i < taskLen; i++ {
		ResultList = append(ResultList, <-chResultList)
	}
	return ResultList, nil
}

func WorkCallback(rawTask TaskInput) (id, guid string, outputData interface{}, err error, ret bool) {
	//DO Batch Jobs
	fmt.Println(rawTask)
	outputData = map[string]string{"id": rawTask.Id, "guid": rawTask.Guid}
	fmt.Printf("Doing Work...id=%s,guid=%s\n", rawTask.Id, rawTask.Guid)
	time.Sleep(10 * time.Second)
	ret = true
	return rawTask.Id, rawTask.Guid, outputData, nil, ret
}

func init() {
	//以时间作为初始化种子
	rand.Seed(time.Now().UnixNano())
}

func main() {
	worknum := flag.Int("worknum", 0, "worker num")
	flag.Parse()

	if *worknum <= 0 {
		fmt.Println("Usage Error!")
		return
	}

	jobsize := 200
	tasklist := make([]TaskInput, 0)
	for i := 0; i < jobsize; i++ {
		u, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
			continue
		}
		input := &TaskInput{
			Id:   strconv.Itoa(i),
			Guid: string(u.String())}
		//fmt.Println(input)
		tasklist = append(tasklist, *input)
	}

	//outoutdata:=make([]RetTask,10000)
	output, _ := DoFixedSizeWorks(*worknum, tasklist, WorkCallback)

	for _, v := range output {
		fmt.Println("result=", v.Id, v.OutputData)
	}
}
