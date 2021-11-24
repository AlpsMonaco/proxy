package util

import (
	"encoding/json"
	"fmt"
	"runtime"
)

func ErrorCatch(err error) {
	var i int = 0
	var m = map[string]interface{}{
		"err":    err.Error(),
		"stacks": &[]string{},
	}
	stackspointer := m["stacks"].(*[]string)
	for {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		i++
		*stackspointer = append(*stackspointer, fmt.Sprintf("%s:%d", file, line))
	}

	b, err := json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
