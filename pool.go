package goworkerpool

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Pool struct {
	concurrency int
}

func PoolFuncCallforResults(concurrency int, TaskInputList []TaskInput, workFunc WorkFunc) ([]TaskResult, error) {
	var (
		lock sync.RWMutex
	)
	ResultList := make([]TaskResult, 0)
	taskLen := len(TaskInputList)
	if taskLen == 0 {
		return nil, ERROR_PARAM_ERROR
	}
	chInputDataList := make(chan TaskInput, taskLen)
	for _, task := range TaskInputList {
		chInputDataList <- task
	}

	chResultList := make(chan TaskResult, taskLen)

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

				//do real job
				result := workFunc(taskData)

				//send back results
				chResultList <- *result
			}
		}()
	}

	//block for getting result
	for i := 0; i < taskLen; i++ {
		ResultList = append(ResultList, <-chResultList)
	}
	return ResultList, nil
}

// 业务逻辑
func WorkCallback(rawTask TaskInput) *TaskResult {
	var (
		taskResult *TaskResult
	)
	outputData := map[string]string{"id": rawTask.Id, "guid": rawTask.Guid}
	fmt.Printf("Doing Work...id=%s,guid=%s\n", rawTask.Id, rawTask.Guid)
	time.Sleep(2 * time.Second)

	taskResult = &TaskResult{
		Id:         rawTask.Id,
		Guid:       rawTask.Guid,
		OutputData: outputData,
		Err:        nil,
		Boolret:    true,
	}
	return taskResult
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
		tasklist = append(tasklist, *input)
	}

	output, _ := PoolFuncCallforResults(*worknum, tasklist, WorkCallback)

	for _, v := range output {
		fmt.Println("result=", v.Id, v.OutputData)
	}
}
