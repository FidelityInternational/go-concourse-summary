package summary_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/FidelityInternational/go-concourse-summary/concourse"
)

var _ = Describe("#SetupConfig", func() {
	Context("when refreshInterval is blank", func() {
		It("returns populated config with the default refresh interval", func() {
			config, err := summary.SetupConfig("", "", "", "")
			Ω(err).Should(BeNil())
			Ω(config.RefreshInterval).Should(Equal(30))
		})
	})

	Context("when refreshInterval is not blank", func() {
		Context("and refreshInterval cannot be converted to an int", func() {
			It("returns an error", func() {
				config, err := summary.SetupConfig("notANumber", "", "", "")
				Ω(err).Should(MatchError(`strconv.Atoi: parsing "notANumber": invalid syntax`))
				Ω(config).Should(Equal(&summary.Config{}))
			})
		})

		Context("and refreshInterval can be converted to an int", func() {
			Context("and refresh interval is less than 1", func() {
				It("returns populated config with the default refresh interval", func() {
					config, err := summary.SetupConfig("-1", "", "", "")
					Ω(err).Should(BeNil())
					Ω(config.RefreshInterval).Should(Equal(30))
				})
			})

			Context("and refresh interval is greater than or equal to 1", func() {
				It("returns populated config with the provided refresh interval", func() {
					config, err := summary.SetupConfig("1", "", "", "")
					Ω(err).Should(BeNil())
					Ω(config.RefreshInterval).Should(Equal(1))
				})
			})
		})
	})

	Context("when groupsJSON is blank", func() {
		It("returns populated config with the empty CSGroups", func() {
			config, err := summary.SetupConfig("", "", "", "")
			Ω(err).Should(BeNil())
			Ω(config.CSGroups).Should(Equal(summary.CSGroups{}))
		})
	})

	Context("when groupsJSON is not blank", func() {
		Context("and the JSON is invalid", func() {
			It("returns an error", func() {
				config, err := summary.SetupConfig("", "{]", "", "")
				Ω(err).Should(MatchError(`invalid character ']' looking for beginning of object key string`))
				Ω(config).Should(Equal(&summary.Config{}))
			})
		})

		Context("and the JSON is valid", func() {
			It("returns populated config with the CSGroups", func() {
				config, err := summary.SetupConfig("", `[{"group": "test", "hosts": []}]`, "", "")
				Ω(err).Should(BeNil())
				Ω(config.CSGroups).Should(Equal(summary.CSGroups{
					{
						Group: "test",
						Hosts: []summary.Host{},
					},
				}))
			})
		})
	})

	Context("when hostsJSON is blank", func() {
		It("returns populated config with the empty Hosts", func() {
			config, err := summary.SetupConfig("", "", "", "")
			Ω(err).Should(BeNil())
			Ω(config.CSGroups).Should(Equal(summary.CSGroups{}))
		})
	})

	Context("when hostsJSON is not blank", func() {
		Context("and the JSON is invalid", func() {
			It("returns an error", func() {
				config, err := summary.SetupConfig("", "", "[}", "")
				Ω(err).Should(MatchError(`invalid character '}' looking for beginning of value`))
				Ω(config).Should(Equal(&summary.Config{}))
			})
		})

		Context("and the JSON is valid", func() {
			It("returns populated config with the Hosts", func() {
				config, err := summary.SetupConfig("", "", `["host1", "host2"]`, "")
				Ω(err).Should(BeNil())
				Ω(config.Hosts).Should(Equal([]summary.Host{
					{
						FQDN: "host1",
					},
					{
						FQDN: "host2",
					},
				}))
			})
		})
	})

	Context("when skipSSLValidationString is blank", func() {
		It("returns populated config with the skipSSLValidation as false", func() {
			config, err := summary.SetupConfig("", "", "", "")
			Ω(err).Should(BeNil())
			Ω(config.SkipSSLValidation).Should(BeFalse())
		})
	})

	Context("when skipSSLValidationString is not blank", func() {
		Context("and the value of skipSSLValidation string is not true", func() {
			It("returns populated config with the skipSSLValidation as false", func() {
				config, err := summary.SetupConfig("", "", "", "test")
				Ω(err).Should(BeNil())
				Ω(config.SkipSSLValidation).Should(BeFalse())
			})
		})

		Context("and the value of skipSSLValidation string is true", func() {
			It("returns populated config with the skipSSLValidation as true", func() {
				config, err := summary.SetupConfig("", "", "", "true")
				Ω(err).Should(BeNil())
				Ω(config.SkipSSLValidation).Should(BeTrue())
			})
		})
	})
})
