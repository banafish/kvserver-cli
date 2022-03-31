package client

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	key = "test"

	TestTypeRead  = "read"
	TestTypeWrite = "write"
)

func SpeedTest(serverIDs []string, threadNum int, testType string, testTime int) {
	var cnt int64
	wg := &sync.WaitGroup{}
	wg.Add(threadNum)
	switch testType {
	case TestTypeRead:
		for i := 0; i < threadNum; i++ {
			go read(clone(serverIDs), wg, &cnt, testTime)
		}
	case TestTypeWrite:
		for i := 0; i < threadNum; i++ {
			go write(clone(serverIDs), wg, &cnt, testTime)
		}
	default:
		fmt.Println("error testType:", testType)
		return
	}
	wg.Wait()
	fmt.Printf("%v %v op/s with %v client\n", testType, cnt/int64(testTime), threadNum)
}

func read(serverIDs []string, wg *sync.WaitGroup, cnt *int64, testTime int) {
	t := time.Now()
	ck := MakeClerk(serverIDs)
	for time.Since(t).Seconds() < float64(testTime) {
		if _, err := ck.Get(key); err != nil {
			fmt.Println(err)
			ck = MakeClerk(clone(serverIDs))
		} else {
			atomic.AddInt64(cnt, 1)
		}
	}
	wg.Done()
}

func write(serverIDs []string, wg *sync.WaitGroup, cnt *int64, testTime int) {
	t := time.Now()
	ck := MakeClerk(serverIDs)
	for time.Since(t).Seconds() < float64(testTime) {
		v := fmt.Sprintf("%v", time.Now().UnixNano())
		if err := ck.Put(key, v); err != nil {
			fmt.Println(err)
			ck = MakeClerk(clone(serverIDs))
		} else {
			atomic.AddInt64(cnt, 1)
		}
	}
	wg.Done()
}

func clone(orig []string) []string {
	tmp := make([]string, len(orig))
	copy(tmp, orig)
	return tmp
}
