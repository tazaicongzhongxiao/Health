package app

import (
	"fmt"
	"runtime/debug"
)

func TryErr() {
	defer func() {
		errs := recover()
		if errs == nil {
			return
		}
		fmt.Println(fmt.Sprintf("%v", errs))
		fmt.Println(string(debug.Stack()))
	}()
}
