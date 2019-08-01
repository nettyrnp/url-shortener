package http_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttpService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HttpService Suite")
}
