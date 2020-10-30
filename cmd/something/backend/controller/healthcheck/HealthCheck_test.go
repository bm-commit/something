package healthcheck

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func setupServer() *gin.Engine {
	router := gin.Default()
	RegisterRoutes(router)
	return router
}

func TestHealthCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Health Check Suite")
}

var _ = Describe("Server", func() {
	var server *httptest.Server

	BeforeEach(func() {
		// start a test http server
		server = httptest.NewServer(setupServer())
	})

	AfterEach(func() {
		server.Close()
	})

	Context("When GET request is sent to /health-check", func() {
		It("Returns the status OK response", func() {
			resp, err := http.Get(server.URL + "/health-check")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(Equal(`{"status":"ok"}`))
		})
	})
})
