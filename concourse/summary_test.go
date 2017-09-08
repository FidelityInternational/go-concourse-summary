package summary_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"unicode"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/FidelityInternational/go-concourse-summary/concourse"
)

func stringMinifier(in string) (out string) {
	for _, c := range in {
		if !unicode.IsSpace(c) {
			out = out + string(c)
		}
	}
	return
}

var _ = Describe("#SetupConfig", func() {
	var (
		config                                                          *summary.Config
		err                                                             error
		refreshInterval, groupsJSON, hostsJSON, skipSSLValidationString string
	)

	JustBeforeEach(func() {
		config, err = summary.SetupConfig(refreshInterval, groupsJSON, hostsJSON, skipSSLValidationString)
	})

	AfterEach(func() {
		refreshInterval = ""
		groupsJSON = ""
		hostsJSON = ""
		skipSSLValidationString = ""
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
