package summary_test

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"regexp"
	"unicode"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gorilla/mux"

	"github.com/FidelityInternational/go-concourse-summary/concourse"
)

func Router(config *summary.Config) *mux.Router {
	server := &summary.Server{Config: config}
	r := server.Start()
	return r
}

func init() {
	http.Handle("/", Router(&summary.Config{}))
}

func stringMinifier(in string) (out string) {
	for _, c := range in {
		if !unicode.IsSpace(c) {
			out = out + string(c)
		}
	}
	return
}

func stripDate(in string) (out string) {
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{4}:\d{2}:\d{2}(&#43;|-)\d{4}`)
	return re.ReplaceAllString(in, "yyyy-mm-ddhh:mm:ss&#43;zzzz")
}

func stripHostPort(in string) (out string) {
	re := regexp.MustCompile(`127\.0\.0\.1:\d{1,6}`)
	return re.ReplaceAllString(in, "127.0.0.1:pppp")
}

var _ = Describe("#SetupConfig", func() {
	var (
		config                                                                    *summary.Config
		err                                                                       error
		refreshInterval, groupsJSON, hostsJSON, skipSSLValidationString, teamName string
	)

	JustBeforeEach(func() {
		config, err = summary.SetupConfig(refreshInterval, groupsJSON, hostsJSON, skipSSLValidationString, teamName)
	})

	AfterEach(func() {
		refreshInterval = ""
		groupsJSON = ""
		hostsJSON = ""
		skipSSLValidationString = ""
		teamName = ""
	})

	Context("when refreshInterval is blank", func() {
		It("returns populated config with the default refresh interval", func() {
			Ω(err).Should(BeNil())
			Ω(config.RefreshInterval).Should(Equal(30))
			Ω(config.Protocol).Should(Equal("https"))
		})
	})

	Context("when refreshInterval is not blank", func() {
		Context("and refreshInterval cannot be converted to an int", func() {
			BeforeEach(func() {
				refreshInterval = "notANumber"
			})

			It("returns an error", func() {
				Ω(err).Should(MatchError(`strconv.Atoi: parsing "notANumber": invalid syntax`))
				Ω(config).Should(Equal(&summary.Config{}))
			})
		})

		Context("and refreshInterval can be converted to an int", func() {
			Context("and refresh interval is less than 1", func() {
				BeforeEach(func() {
					refreshInterval = "-1"
				})

				It("returns populated config with the default refresh interval", func() {
					Ω(err).Should(BeNil())
					Ω(config.RefreshInterval).Should(Equal(30))
					Ω(config.Protocol).Should(Equal("https"))
				})
			})

			Context("and refresh interval is greater than or equal to 1", func() {
				BeforeEach(func() {
					refreshInterval = "1"
				})

				It("returns populated config with the provided refresh interval", func() {
					Ω(err).Should(BeNil())
					Ω(config.RefreshInterval).Should(Equal(1))
					Ω(config.Protocol).Should(Equal("https"))
				})
			})
		})
	})

	Context("when groupsJSON is blank", func() {
		It("returns populated config with the empty CSGroups", func() {
			Ω(err).Should(BeNil())
			Ω(config.CSGroups).Should(Equal(summary.CSGroups{}))
			Ω(config.Protocol).Should(Equal("https"))
		})
	})

	Context("when groupsJSON is not blank", func() {
		Context("and the JSON is invalid", func() {
			BeforeEach(func() {
				groupsJSON = "{]"
			})

			It("returns an error", func() {
				Ω(err).Should(MatchError(`invalid character ']' looking for beginning of object key string`))
				Ω(config).Should(Equal(&summary.Config{}))
			})
		})

		Context("and the JSON is valid", func() {
			BeforeEach(func() {
				groupsJSON = `[{"group": "test", "hosts": []}]`
			})

			It("returns populated config with the CSGroups", func() {
				Ω(err).Should(BeNil())
				Ω(config.CSGroups).Should(Equal(summary.CSGroups{
					{
						Group: "test",
						Hosts: []summary.Host{},
					},
				}))
				Ω(config.Protocol).Should(Equal("https"))
			})
		})
	})

	Context("when hostsJSON is blank", func() {
		It("returns populated config with the empty Hosts", func() {
			Ω(err).Should(BeNil())
			Ω(config.CSGroups).Should(Equal(summary.CSGroups{}))
			Ω(config.Protocol).Should(Equal("https"))
		})
	})

	Context("when hostsJSON is not blank", func() {
		Context("and the JSON is invalid", func() {
			BeforeEach(func() {
				hostsJSON = "[}"
			})

			It("returns an error", func() {
				Ω(err).Should(MatchError(`invalid character '}' looking for beginning of value`))
				Ω(config).Should(Equal(&summary.Config{}))
			})
		})

		Context("and the JSON is valid", func() {
			BeforeEach(func() {
				hostsJSON = `["host1", "host2"]`
			})

			It("returns populated config with the Hosts", func() {
				Ω(err).Should(BeNil())
				Ω(config.Hosts).Should(Equal([]summary.Host{
					{
						FQDN: "host1",
					},
					{
						FQDN: "host2",
					},
				}))
				Ω(config.Protocol).Should(Equal("https"))
			})
		})
	})

	Context("when skipSSLValidationString is blank", func() {
		It("returns populated config with the skipSSLValidation as false", func() {
			Ω(err).Should(BeNil())
			Ω(config.SkipSSLValidation).Should(BeFalse())
		})
	})

	Context("when skipSSLValidationString is not blank", func() {
		Context("and the value of skipSSLValidation string is not true", func() {
			BeforeEach(func() {
				skipSSLValidationString = "test"
			})

			It("returns populated config with the skipSSLValidation as false", func() {
				Ω(err).Should(BeNil())
				Ω(config.SkipSSLValidation).Should(BeFalse())
				Ω(config.Protocol).Should(Equal("https"))
			})
		})

		Context("and the value of skipSSLValidation string is true", func() {
			BeforeEach(func() {
				skipSSLValidationString = "true"
			})

			It("returns populated config with the skipSSLValidation as true", func() {
				Ω(err).Should(BeNil())
				Ω(config.SkipSSLValidation).Should(BeTrue())
				Ω(config.Protocol).Should(Equal("https"))
			})
		})
	})
})

var _ = Describe("config#Index", func() {
	var (
		templates    = template.Must(template.ParseGlob("../templates/*"))
		mockRecorder *httptest.ResponseRecorder
		config       = &summary.Config{
			Templates: templates,
		}
	)

	AfterEach(func() {
		config = &summary.Config{
			Templates: templates,
		}
	})

	JustBeforeEach(func() {
		mockRecorder = httptest.NewRecorder()
		config.Index(mockRecorder, &http.Request{})
	})

	Context("when no hosts or groups are configured", func() {
		It("writes an index page", func() {
			Ω(mockRecorder.Code).Should(Equal(200))
			Ω(stringMinifier(mockRecorder.Body.String())).Should(Equal(stringMinifier(`
<!DOCTYPE html>
<html>
	<head rel="v2">
		<title>Concourse Summary</title>
		<link rel="icon" type="image/png" href="/favicon.png" sizes="32x32">
		<link rel="stylesheet" type="text/css" href="/styles.css">
		<script src="/favico-0.3.10.min.js"></script>
		<script src="/refresh.js"></script>
	</head>
	<body>
		<h1>Concourse Summary</h1>
		<p>Use the URL path to show a summary, eg, '/host/[HOST NAME]'</p>


		<p>This project can be found on <a href="https://github.com/FidelityInternational/go-concourse-summary" target="_blank">Github</a></p>
	</body>
</html>`)))
		})
	})

	Context("when hosts are configured", func() {
		BeforeEach(func() {
			config.Hosts = []summary.Host{
				{
					FQDN: "test1",
				},
				{
					FQDN: "test2",
				},
			}
		})

		It("writes an index page", func() {
			Ω(mockRecorder.Code).Should(Equal(200))
			Ω(stringMinifier(mockRecorder.Body.String())).Should(Equal(stringMinifier(`
<!DOCTYPE html>
<html>
	<head rel="v2">
		<title>Concourse Summary</title>
		<link rel="icon" type="image/png" href="/favicon.png" sizes="32x32">
		<link rel="stylesheet" type="text/css" href="/styles.css">
		<script src="/favico-0.3.10.min.js"></script>
		<script src="/refresh.js"></script>
	</head>
	<body>
		<h1>Concourse Summary</h1>
		<p>Use the URL path to show a summary, eg, '/host/[HOST NAME]'</p>

		<div><a href="/host/test1">
			test1
		</a></div>

		<div><a href="/host/test2">
			test2
		</a></div>


		<p>This project can be found on <a href="https://github.com/FidelityInternational/go-concourse-summary" target="_blank">Github</a></p>
	</body>
</html>`)))
		})
	})

	Context("when groups are configured", func() {
		BeforeEach(func() {
			config.CSGroups = summary.CSGroups{
				{
					Group: "testGroup1",
				},
				{
					Group: "testGroup2",
				},
			}
		})

		It("writes an index page", func() {
			Ω(mockRecorder.Code).Should(Equal(200))
			Ω(stringMinifier(mockRecorder.Body.String())).Should(Equal(stringMinifier(`
<!DOCTYPE html>
<html>
	<head rel="v2">
		<title>Concourse Summary</title>
		<link rel="icon" type="image/png" href="/favicon.png" sizes="32x32">
		<link rel="stylesheet" type="text/css" href="/styles.css">
		<script src="/favico-0.3.10.min.js"></script>
		<script src="/refresh.js"></script>
	</head>
	<body>
		<h1>Concourse Summary</h1>
		<p>Use the URL path to show a summary, eg, '/host/[HOST NAME]'</p>


			<div style="margin-top:2em">Groups</div>

				<div><a href="/group/testGroup1">
					testGroup1
				</a></div>

				<div><a href="/group/testGroup2">
					testGroup2
				</a></div>


		<p>This project can be found on <a href="https://github.com/FidelityInternational/go-concourse-summary" target="_blank">Github</a></p>
	</body>
</html>`)))
		})
	})

	Context("when hosts and groups are configured", func() {
		BeforeEach(func() {
			config.Hosts = []summary.Host{
				{
					FQDN: "test1",
				},
				{
					FQDN: "test2",
				},
			}

			config.CSGroups = summary.CSGroups{
				{
					Group: "testGroup1",
				},
				{
					Group: "testGroup2",
				},
			}

		})

		It("writes an index page", func() {
			Ω(mockRecorder.Code).Should(Equal(200))
			Ω(stringMinifier(mockRecorder.Body.String())).Should(Equal(stringMinifier(`
<!DOCTYPE html>
<html>
	<head rel="v2">
		<title>Concourse Summary</title>
		<link rel="icon" type="image/png" href="/favicon.png" sizes="32x32">
		<link rel="stylesheet" type="text/css" href="/styles.css">
		<script src="/favico-0.3.10.min.js"></script>
		<script src="/refresh.js"></script>
	</head>
	<body>
		<h1>Concourse Summary</h1>
		<p>Use the URL path to show a summary, eg, '/host/[HOST NAME]'</p>

		<div><a href="/host/test1">
			test1
		</a></div>

		<div><a href="/host/test2">
			test2
		</a></div>


			<div style="margin-top:2em">Groups</div>

				<div><a href="/group/testGroup1">
					testGroup1
				</a></div>

				<div><a href="/group/testGroup2">
					testGroup2
				</a></div>


		<p>This project can be found on <a href="https://github.com/FidelityInternational/go-concourse-summary" target="_blank">Github</a></p>
	</body>
</html>`)))
		})
	})
})

var _ = Describe("#HostSummary", func() {
	var (
		templates    = template.Must(template.ParseGlob("../templates/*"))
		mockRecorder *httptest.ResponseRecorder
		config       = &summary.Config{
			Templates: templates,
			Protocol:  "http",
		}
	)

	AfterEach(func() {
		if server != nil {
			teardown()
		}

		config = &summary.Config{
			Templates: templates,
			Protocol:  "http",
		}
	})

	JustBeforeEach(func() {
		mockRecorder = httptest.NewRecorder()

		req, _ := http.NewRequest("GET", fmt.Sprintf("http://example.com/host/%s", Host(server)), nil)
		Router(config).ServeHTTP(mockRecorder, req)
	})

	Context("when concourse returns invalid json", func() {
		BeforeEach(func() {
			mocks := []MockRoute{
				{"GET", "/api/v1/teams/pipelines", "[}", 200, "", nil},
			}
			setupMultiple(mocks)
		})

		It("returns an error", func() {
			Ω(mockRecorder.Code).Should(Equal(500))
			Ω(mockRecorder.Body.String()).Should(MatchRegexp(`Error collecting data from concourse \(127.0.0.1:\d{1,6}\) please refer to logs for more details`))
		})
	})

	Context("when concourse has no pipelines", func() {
		BeforeEach(func() {
			mocks := []MockRoute{
				{"GET", "/api/v1/teams/pipelines", "[]", 200, "", nil},
			}
			setupMultiple(mocks)
		})

		It("returns a formatted blank page", func() {
			Ω(mockRecorder.Code).Should(Equal(200))
			Ω(stripDate(stringMinifier(mockRecorder.Body.String()))).Should(Equal(stripDate(stringMinifier(`
<!DOCTYPE html>
<html>
	<head rel="v2">
		<title>Concourse Summary</title>
		<link rel="icon" type="image/png" href="/favicon.png" sizes="32x32">
		<link rel="stylesheet" type="text/css" href="/styles.css">
		<script>window.refresh_interval =  0 </script>
		<script src="/favico-0.3.10.min.js"></script>
		<script src="/refresh.js"></script>
	</head>
	<body>
		<div class="time">
			2017-09-08 15:17:56 &#43;0100 (<span id="countdown">0</span>)
			<div class="right">
				<a class="github" href="https://github.com/FidelityInternational/go-concourse-summary" target="_blank">&nbsp;</a>
			</div>
		</div>

<div class="scalable">



</div>

	</body>
</html>
				`))))
		})
	})

	Context("and concourse has pipelines", func() {
		BeforeEach(func() {
			mocks := []MockRoute{
				{"GET", "/api/v1/teams/pipelines", pipelinesPayload, 200, "", nil},
				{"GET", "/api/v1/teams/pipelines/test1/jobs", jobsPayload, 200, "", nil},
			}
			setupMultiple(mocks)
		})

		It("returns a page with status etc", func() {
			Ω(mockRecorder.Code).Should(Equal(200))
			Ω(stripHostPort(stripDate(stringMinifier(mockRecorder.Body.String())))).Should(Equal(stripHostPort(stripDate(stringMinifier(`
<!DOCTYPE html>
<html>
	<head rel="v2">
		<title>Concourse Summary</title>
		<link rel="icon" type="image/png" href="/favicon.png" sizes="32x32">
		<link rel="stylesheet" type="text/css" href="/styles.css">
		<script>window.refresh_interval =  0 </script>
		<script src="/favico-0.3.10.min.js"></script>
		<script src="/refresh.js"></script>
	</head>
	<body>
		<div class="time">
			2017-09-08 17:05:56 &#43;0100 (<span id="countdown">0</span>)
			<div class="right">
				<a class="github" href="https://github.com/FidelityInternational/go-concourse-summary" target="_blank">&nbsp;</a>
			</div>
		</div>

<div class="scalable">


	<a href="http://127.0.0.1:49898/test1.url" target="_blank" class="outer">
	<div class="status">
		<div class="paused_job" style="width: 0%;"></div>
		<div class="aborted" style="width: 16%;"></div>
		<div class="errored" style="width: 16%;"></div>
		<div class="failed" style="width: 16%;"></div>
		<div class="succeeded" style="width: 16%;"></div>
	</div>


	<div class="inner">
		<span class="test1"><span>test1</span></span>
		<span class=""><span></span></span>
	</div>
	</a>


</div>

	</body>
</html>`)))))
		})
	})
})

var _ = Describe("#GroupSummary", func() {
	var (
		templates    = template.Must(template.ParseGlob("../templates/*"))
		mockRecorder *httptest.ResponseRecorder
		config       = &summary.Config{
			Templates: templates,
			Protocol:  "http",
		}
	)

	AfterEach(func() {
		if server != nil {
			teardown()
		}

		config = &summary.Config{
			Templates: templates,
			Protocol:  "http",
		}
	})

	JustBeforeEach(func() {
		mockRecorder = httptest.NewRecorder()

		config.CSGroups = []summary.CSGroup{
			{
				Group: "test",
				Hosts: []summary.Host{
					{
						FQDN: Host(server),
					},
				},
			},
		}

		req, _ := http.NewRequest("GET", "http://example.com/group/test", nil)
		Router(config).ServeHTTP(mockRecorder, req)
	})

	Context("when concourse returns invalid json", func() {
		BeforeEach(func() {
			mocks := []MockRoute{
				{"GET", "/api/v1/teams/pipelines", "[}", 200, "", nil},
			}
			setupMultiple(mocks)
		})

		It("returns an error", func() {
			Ω(mockRecorder.Code).Should(Equal(500))
			Ω(mockRecorder.Body.String()).Should(MatchRegexp(`Error collecting data from concourse \(127.0.0.1:\d{1,6}\) please refer to logs for more details`))
		})
	})

	Context("when concourse has no pipelines", func() {
		BeforeEach(func() {
			mocks := []MockRoute{
				{"GET", "/api/v1/teams/pipelines", "[]", 200, "", nil},
			}
			setupMultiple(mocks)
		})

		It("returns a formatted blank page", func() {
			Ω(mockRecorder.Code).Should(Equal(200))
			Ω(stripDate(stripHostPort(stringMinifier(mockRecorder.Body.String())))).Should(Equal(stripHostPort(stripDate(stringMinifier(`
<!DOCTYPE html>
<html>
  <head rel="v2">
    <title>Concourse Summary</title>
    <link rel="icon" type="image/png" href="/favicon.png" sizes="32x32">
    <link rel="stylesheet" type="text/css" href="/styles.css">
    <script>window.refresh_interval =  0 </script>
    <script src="/favico-0.3.10.min.js"></script>
    <script src="/refresh.js"></script>
  </head>
  <body>
    <div class="time">
      2017-09-13 09:38:03 &#43;0100 (<span id="countdown">0</span>)
      <div class="right">
        <a class="github" href="https://github.com/FidelityInternational/go-concourse-summary" target="_blank">&nbsp;</a>
      </div>
    </div>


<div class="group">
  <a href="/host/127.0.0.1:53553">127.0.0.1:53553</a>
  <div>



  </div>
</div>


  </body>
</html>
				`)))))
		})
	})

	Context("and concourse has pipelines", func() {
		BeforeEach(func() {
			mocks := []MockRoute{
				{"GET", "/api/v1/teams/pipelines", pipelinesPayload, 200, "", nil},
				{"GET", "/api/v1/teams/pipelines/test1/jobs", jobsPayload, 200, "", nil},
			}
			setupMultiple(mocks)
		})

		It("returns a page with status etc", func() {
			Ω(mockRecorder.Code).Should(Equal(200))
			Ω(stripHostPort(stripDate(stringMinifier(mockRecorder.Body.String())))).Should(Equal(stripHostPort(stripDate(stringMinifier(`
<!DOCTYPE html>
<html>
  <head rel="v2">
    <title>Concourse Summary</title>
    <link rel="icon" type="image/png" href="/favicon.png" sizes="32x32">
    <link rel="stylesheet" type="text/css" href="/styles.css">
    <script>window.refresh_interval =  0 </script>
    <script src="/favico-0.3.10.min.js"></script>
    <script src="/refresh.js"></script>
  </head>
  <body>
    <div class="time">
      2017-09-13 09:38:03 &#43;0100 (<span id="countdown">0</span>)
      <div class="right">
        <a class="github" href="https://github.com/FidelityInternational/go-concourse-summary" target="_blank">&nbsp;</a>
      </div>
    </div>


<div class="group">
  <a href="/host/127.0.0.1:53555">127.0.0.1:53555</a>
  <div>


  <a href="http://127.0.0.1:53555/test1.url" target="_blank" class="outer">
  <div class="status">
    <div class="paused_job" style="width: 0%;"></div>
    <div class="aborted" style="width: 16%;"></div>
    <div class="errored" style="width: 16%;"></div>
    <div class="failed" style="width: 16%;"></div>
    <div class="succeeded" style="width: 16%;"></div>
  </div>


  <div class="inner">
    <span class="test1"><span>test1</span></span>
    <span class=""><span></span></span>
  </div>
  </a>


  </div>
</div>


  </body>
</html>`)))))
		})
	})
})
