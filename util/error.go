package util

import (
	"os"

	"github.com/potterxu/tstool/logger"
)

func LogOnErr(err error) {
	if err != nil {
		logger.Logger.Fatal(err)
	}
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ExitOnErr(err error) {
	ExitOnErrWithCode(err, 1)
}

func ExitOnErrWithCode(err error, code int) {
	if err != nil {
		logger.Logger.Fatal(err)
		os.Exit(code)
	}
}
