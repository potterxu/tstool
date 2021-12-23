package logger

import (
	"io"
	"log"
	"os"
)

var (
	Logger = log.New(io.MultiWriter(stdout), "[Common]", log.Lshortfile|log.LstdFlags)
)

var (
	stdout = os.Stdout
)
