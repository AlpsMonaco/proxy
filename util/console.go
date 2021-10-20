package util

import (
	"fmt"
	"time"
)

type console struct{}

var Console console

func (c *console) Info(log ...interface{}) {
	c.Format("INFO", fmt.Sprint(log...))
}

func (c *console) Format(logType, log string) {
	fmt.Printf("[%s]\t[%s]\t%s\n", time.Now().Format("2006-01-02 15:04:05"), logType, log)
}

func (c *console) Error(log ...interface{}) {
	c.Format("Error", fmt.Sprint(log...))
}
