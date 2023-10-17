// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"os"

	"github.com/osintami/fingerprintz/log"
)

func init() {
	os.Remove("./test/test.log")
	log.InitLogger("./test/", "test.log", "INFO", false)
}
