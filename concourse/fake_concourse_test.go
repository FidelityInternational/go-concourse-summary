package summary_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/onsi/ginkgo"
)

var (
	server *httptest.Server
)

func Host(server *httptest.Server) string {
	split := strings.Split(server.URL, "://")
	return split[1]
}

type MockRoute struct {
	Method      string
	Endpoint    string
	Output      string
	Status      int
	QueryString string
	PostForm    *string
}

func testQueryString(QueryString string, QueryStringExp string) {
	value, _ := url.QueryUnescape(QueryString)

	if QueryStringExp != value {
		defer ginkgo.GinkgoRecover()
		ginkgo.Fail(fmt.Sprintf("Error: Query string '%s' should be equal to '%s'", QueryStringExp, value))
	}
}

func testPostQuery(req *http.Request, postFormBody *string) {
	if postFormBody != nil {
		if body, err := ioutil.ReadAll(req.Body); err != nil {
			defer ginkgo.GinkgoRecover()
			ginkgo.Fail("No request body but expected one")
		} else {
			defer req.Body.Close()
			if strings.TrimSpace(string(body)) != strings.TrimSpace(*postFormBody) {
				defer ginkgo.GinkgoRecover()
				ginkgo.Fail(fmt.Sprintf("Expected POST body (%s) does not equal POST body (%s)", *postFormBody, body))
			}
		}
	}
}

func setup(mock MockRoute) {
	setupMultiple([]MockRoute{mock})
}

func setupMultiple(mockEndpoints []MockRoute) {
	router := mux.NewRouter()

	for _, mock := range mockEndpoints {
		method := mock.Method
		endpoint := mock.Endpoint
		output := mock.Output
		status := mock.Status
		queryString := mock.QueryString
		postFormBody := mock.PostForm

		if method == "POST" {
			router.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
				testQueryString(r.URL.RawQuery, queryString)
				testPostQuery(r, postFormBody)
				w.WriteHeader(status)
				fmt.Fprintf(w, output)
			}).Methods(method)
		} else {
			router.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
				testQueryString(r.URL.RawQuery, queryString)
				w.WriteHeader(status)
				fmt.Fprintf(w, output)
			}).Methods(method)
		}
	}

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer ginkgo.GinkgoRecover()
		ginkgo.Fail(fmt.Sprintf("Route requested but not mocked: %s", r.URL))
	})

	server = httptest.NewServer(router)
}

func teardown() {
	server.Close()
	server = nil
}
