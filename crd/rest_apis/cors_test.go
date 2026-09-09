package restapis

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"kubesonde.io/controllers/state"
)

var _ = Describe("withCORS", func() {
	var stateManager *state.StateManager

	BeforeEach(func() {
		stateManager = state.NewStateManager()
	})

	It("adds Access-Control-Allow-Origin to GET responses through the wrapped handler", func() {
		req := httptest.NewRequest("GET", "http://localhost:2709/probes", nil)
		w := httptest.NewRecorder()

		withCORS(GetProbesHandlerWithManager(stateManager)).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(200))
		Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("*"))
		Expect(resp.Header.Get("Access-Control-Allow-Methods")).To(Equal("GET, POST, OPTIONS"))
		Expect(resp.Header.Get("Access-Control-Allow-Headers")).To(Equal("Content-Type"))
	})

	It("answers OPTIONS preflight with CORS headers and no body without hitting the underlying handler", func() {
		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})

		req := httptest.NewRequest("OPTIONS", "http://localhost:2709/probes", nil)
		w := httptest.NewRecorder()

		withCORS(next).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		Expect(called).To(BeFalse())
		Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("*"))
		Expect(resp.Header.Get("Access-Control-Allow-Methods")).To(Equal("GET, POST, OPTIONS"))
		Expect(resp.Header.Get("Access-Control-Allow-Headers")).To(Equal("Content-Type"))
	})
})
