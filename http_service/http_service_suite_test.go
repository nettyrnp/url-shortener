package http_service_test

//package http_service

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttpService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HttpService Suite")
}
