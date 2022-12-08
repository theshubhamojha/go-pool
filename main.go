package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type poolService struct {
	limit     int
	processor func(...interface{}) interface{}
	counter   chan bool
}

type Pool interface {
	Process(ctx context.Context, payload ...interface{})
}

func NewGoroutinePool(maxRoutines int, process func(...interface{}) interface{}) Pool {
	ch := make(chan bool, maxRoutines)

	return &poolService{
		limit:     maxRoutines,
		processor: process,
		counter:   ch,
	}
}

func (s *poolService) executeFunction(taskId int, payload ...interface{}) {
	now := time.Now()
	s.processor(payload)
	fmt.Printf("program completed, time taken %v, taskId: %d \n", time.Now().Sub(now).Seconds(), taskId)
	<-s.counter
}

func (s *poolService) Process(ctx context.Context, payload ...interface{}) {
	taskId := rand.Intn(500)
	now := time.Now()

	fmt.Printf("program waiting for execution, id: %d \n", taskId)
	s.counter <- true
	fmt.Printf("program assigned go routine for execution: %d, wait_time: %v \n", taskId, time.Now().Sub(now).Seconds())

	go s.executeFunction(taskId, payload)
}

// driver code ahead
// mocking a sample sync function
func sampleFunction(waitTime ...interface{}) (result interface{}) {
	fmt.Println("wait ", waitTime)
	switch duration := waitTime[0].(type) {
	case int:
		fmt.Println("Will sleep now", waitTime)
		time.Sleep(time.Duration(duration) * time.Second)
	default:
		fmt.Println("error typecasting waitTime variable to integer")
		return false
	}

	return true
}

func main() {
	ctx := context.Background()

	// specify total number of routines and the sync function which we want to execute
	p := NewGoroutinePool(1, sampleFunction)

	p.Process(ctx, 8, 2)
	p.Process(ctx, 3, 1)
	//p.Process(ctx, 7)
	//p.Process(ctx, 4)
	//p.Process(ctx, 2)
	//p.Process(ctx, 3)
	//p.Process(ctx, 5)
	//p.Process(ctx, 6)
}
