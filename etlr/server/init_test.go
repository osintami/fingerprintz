package server

import (
	"os"

	"github.com/osintami/fingerprintz/log"
)

func init() {
	os.Remove("./test/test.log")
	log.InitLogger("./test/", "test.log", "INFO", false)
}
