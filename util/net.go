package util

import "fmt"

func SprintfAddress(addr string, port int) string {
	return fmt.Sprintf("%s:%d", addr, port)
}
