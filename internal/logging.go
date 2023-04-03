package internal

import (
	"fmt"
	"log"
	"os"
)

func NewLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("%-9s ", prefix+":"), log.Ldate|log.Ltime|log.Lmsgprefix)
}
