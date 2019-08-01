package http_service_test

import (
	//"github.com/nettyrnp/url-shortener/http"
	"github.com/nettyrnp/url-shortener/router"
	"io"
	"net/http"
	//"context"
	"net/http/httptest"
	//"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"github.com/nettyrnp/url-shortener/config"
)

func buildRequest(method string, uri string, reader io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, uri, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "text/html; charset=utf-8")
	return request, nil
}

var _ = Describe("Service", func() {
	var (
		//conf  config.HTTPConfig
		//service *http.HTTPService
		//err     error
		//ctx context.Context
		h1 *router.RootHandler
		h2 *router.SvcHandler
	)

	BeforeEach(func() {
		//conf = config.HTTPConfig{
		//		Host:            "127.0.0.1",
		//		Port:            8765,
		//		ShutdownTimeout: time.Duration(10 * time.Second),
		//}
		//ctx = context.Background()
		h1 = &router.RootHandler{Filename: "templates/index.html"}
		h2 = &router.SvcHandler{}
	})

	Describe("Requests to the server", func() {
		var (
			uri      string
			err      error
			request  *http.Request
			recorder *httptest.ResponseRecorder
		)
		BeforeEach(func() {
			recorder = httptest.NewRecorder()
		})
		JustBeforeEach(func() {
			h1.ServeHTTP(recorder, request)
		})

		Describe("getRoot", func() {
			BeforeEach(func() {
				uri = "/"
				//uri = "/v1/login/john"
				request, err = buildRequest("GET", uri, nil)
			})

			Measure("it should getRoot efficiently", func(b Benchmarker) {
				runtime := b.Time("runtime", func() {
					h1.ServeHTTP(recorder, request)
				})
				Ω(runtime.Seconds()).Should(BeNumerically("<", 0.2), "getRoot() shouldn't take too long.")
			}, 50)

			It("should not return an error", func() {
				Ω(err).Should(BeNil())
			})
			It("should return 200 OK", func() {
				Ω(recorder.Code).Should(Equal(http.StatusOK))
			})
			It("should have the correct Content-Type", func() {
				Ω(recorder.HeaderMap.Get("Content-Type")).Should(Equal("text/html; charset=utf-8"))
			})
		})

		Describe("getGeneric", func() {
			BeforeEach(func() {
				uri = "/v1/login/john"
				request, err = buildRequest("GET", uri, nil)
			})

			Measure("it should getGeneric efficiently", func(b Benchmarker) {
				runtime := b.Time("runtime", func() {
					h2.ServeHTTP(recorder, request)
				})
				Ω(runtime.Seconds()).Should(BeNumerically("<", 0.2), "getGeneric() shouldn't take too long.")
			}, 50)

			It("should not return an error", func() {
				Ω(err).Should(BeNil())
			})
			It("should return 200 OK", func() {
				Ω(recorder.Code).Should(Equal(http.StatusOK))
			})
			It("should have the correct Content-Type", func() {
				Ω(recorder.HeaderMap.Get("Content-Type")).Should(Equal("text/html; charset=utf-8"))
			})
		})

	})
})
