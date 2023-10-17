package utils

import (
	"os"

	"github.com/osintami/fingerprintz/log"
)

func init() {
	os.Remove("test.log")
	log.InitLogger("./", "test.log", "INFO", false)
}
