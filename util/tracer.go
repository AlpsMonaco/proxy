package util

import (
	"fmt"
	"runtime"
	"sync"
)

var mu sync.Mutex

func LogTrace(b []byte) {
	mu.Lock()
	defer mu.Unlock()
	_, fileName, line, _ := runtime.Caller(1)
	fmt.Printf("file:%s\tline:%d\t%v\n", fileName, line, b)
}
