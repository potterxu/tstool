package util

import (
	"log"
	"os"
)

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
		log.Fatal(err)
		os.Exit(code)
	}
}
