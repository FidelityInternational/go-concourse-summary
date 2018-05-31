package summary

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/concourse/go-concourse/concourse"
)

// Data concourse data structure
type Data struct {
	Pipeline       string
	Group          string
	URL            string `json:"pipeline_url"`
	Running        bool
	Paused         bool
	BrokenResource bool
	Statuses       map[string]int
}

// GroupData a grouping structure for Data
type GroupData struct {
	Host     string
	Statuses []Data
}

func filterData(data []Data, pipelines []Pipeline) []Data {
	var filteredData []Data
	for _, datum := range data {
		if len(pipelines) == 0 || pipelines == nil {
			filteredData = append(filteredData, datum)
		}
		for _, pipeline := range pipelines {
			if datum.Pipeline != pipeline.Name {
				continue
			}
			if len(pipeline.Groups) == 0 || pipeline.Groups == nil {
				filteredData = append(filteredData, datum)
			} else {
				for _, group := range pipeline.Groups {
					if group == "all" && datum.Group == "" {
						filteredData = append(filteredData, datum)
					}
					if group == datum.Group {
						filteredData = append(filteredData, datum)
					}
				}
			}
		}
	}
	return filteredData
}

func getData(host string, config *Config) ([]Data, error) {
	uri := fmt.Sprintf("%s://%s", config.Protocol, host)
	webURI := uri + "/teams/" + config.Team + "/pipelines/"
	httpClient := createHTTPClient(config)
	client := concourse.NewClient(uri, httpClient, false)
	team := client.Team(config.Team)
	pipelines, err := team.ListPipelines()
	if err != nil {
		return []Data{}, err
	}
	data := map[string]Data{}
	for _, pipeline := range pipelines {
		jobs, err := team.ListJobs(pipeline.Name)
		if err != nil {
			return []Data{}, err
		}
		for _, job := range jobs {
			groups := job.Groups
			if len(groups) == 0 {
				groups = []string{""}
			}
			for _, group := range groups {
				key := fmt.Sprintf("%s:%s", pipeline.Name, group)
				datum := data[key]
				if datum.Statuses == nil {
					datum.Statuses = map[string]int{}
					datum.Pipeline = pipeline.Name
					datum.Group = group
					datum.Paused = pipeline.Paused
					if group == "" {
						datum.URL = fmt.Sprintf("%s%s", webURI, pipeline.Name)
					} else {
						datum.URL = fmt.Sprintf("%s%s?groups=%s", webURI, pipeline.URL, group)
					}
				}
				if !datum.Running {
					datum.Running = (job.NextBuild != nil)
				}
				if job.FinishedBuild != nil {
					datum.Statuses[job.FinishedBuild.Status]++
				} else {
					datum.Statuses["pending"]++
				}
				data[key] = datum
			}
		}
	}
	values := make([]Data, 0, len(data))
	for _, value := range data {
		values = append(values, value)
	}

	sort.Sort(byData(values))

	return values, nil
}

// Percent calculate the a percentage value for a particular status from data statuses
func (d Data) Percent(status string) int {
	if len(d.Statuses) == 0 {
		return 0
	}
	return int((float64(d.Statuses[status]) / float64(mapValueSum(d.Statuses))) * 100)
}

func mapValueSum(sourceData map[string]int) int {
	sum := 0
	for i := range sourceData {
		sum += sourceData[i]
	}
	return sum
}

func createHTTPClient(config *Config) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 2,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: config.SkipSSLValidation},
		},
		Timeout: time.Duration(30) * time.Second,
	}

	return client
}

type byData []Data

func (r byData) Len() int {
	return len(r)
}

func (r byData) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r byData) Less(i, j int) bool {
	first := fmt.Sprintf("%s%s", r[i].Pipeline, r[i].Group)
	second := fmt.Sprintf("%s%s", r[j].Pipeline, r[j].Group)

	return first < second
}
