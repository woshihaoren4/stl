package sync

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type Config struct {
	Name  string
	Count int
}

func TestCOWLock(t *testing.T) {
	lock := NewCOWLock(&Config{})
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			for true {
				read := lock.Read()
				if read.Count%100 == 0 {
					fmt.Println(i, "--->", read.Count)
				}
				if read.Count >= 1000 {
					return
				}
				time.Sleep(time.Millisecond)
			}
		}(i)
	}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				lock.Write(func(c *Config) *Config {
					c.Count += 1
					time.Sleep(time.Millisecond)
					return c
				})
			}
		}(i)
	}
	wg.Wait()
}
