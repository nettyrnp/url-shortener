package util

import (
	"github.com/nettyrnp/url-shortener/log"
	"os"
)

// Die kills the failing program.
func Die(err error) {
	logger := log.GetLogger()
	if err == nil {
		return
	}
	logger.Fatal(err.Error())
	os.Exit(1)
}
