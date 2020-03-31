package ginkinjectgo

import (
	"github.com/onsi/ginkgo"
)

var curEnv *env = &global

func Describe(text string, body func()) (result bool) {
	curEnv = curEnv.Env()
	result = ginkgo.Describe(text, body)
	curEnv = curEnv.parent
	return
}

func It(text string, body interface{}, timeout ...float64) bool {
	return ginkgo.It(text, curEnv.Inject(body), timeout...)
}
