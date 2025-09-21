package iocketsdk

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	green  = "\033[97;42m"
	yellow = "\033[97;43m"
	red    = "\033[97;41m"
	reset  = "\033[0m"
)

func Perror(args ...any) {
	_, file, line, _ := runtime.Caller(1)
	t := time.Now()
	lineStr := strconv.Itoa(line)
	i := strings.LastIndex(file, "/")
	fmt.Println(append([]any{red+"[IOCKET]" + reset, t.Format(time.TimeOnly), file[i+1:]+":"+lineStr,"|"}, args...)...)
}

func Pwarn(args ...any) {
	_, file, line, _ := runtime.Caller(1)
	t := time.Now()
	lineStr := strconv.Itoa(line)
	i := strings.LastIndex(file, "/")
	fmt.Println(append([]any{yellow+"[IOCKET]" + reset, t.Format(time.TimeOnly), file[i+1:]+":"+lineStr,"|"}, args...)...)

}

func P(args ...any) {
	_, file, line, _ := runtime.Caller(1)
	t := time.Now()
	lineStr := strconv.Itoa(line)
	i := strings.LastIndex(file, "/")
	fmt.Println(append([]any{green+"[IOCKET]" + reset, t.Format(time.TimeOnly), file[i+1:]+":"+lineStr,"|"}, args...)...)
}