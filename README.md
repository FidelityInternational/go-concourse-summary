# go-concourse-summary

[![codecov.io](https://codecov.io/github/FidelityInternational/go-concourse-summary/coverage.svg?branch=master)](https://codecov.io/github/FidelityInternational/go-concourse-summary?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/FidelityInternational/go-concourse-summary)](https://goreportcard.com/report/github.com/FidelityInternational/go-concourse-summary)
[![Build Status](https://travis-ci.org/FidelityInternational/go-concourse-summary.svg?branch=master)](https://travis-ci.org/FidelityInternational/go-concourse-summary)

This is a port of [concourse-summary](https://github.com/dgodd/concourse-summary) to Golang. The aim is for all features of `concourse-summary` to be covered in this port. In its current state all features should have been migrated with the exception of the ability to collapse/ expand groups.

The intention of `concourse-summary` is to show a quick overview of all of your [concourse](https://concourse.ci) pipelines and groups in a single summary page.

It was intended that all configuration would be compatible between `go-concourse-summary` and `concourse-summary`. Unfortunately due to the object model differences on Golang marshaling it has been necessary to rework the format of `HOSTS` and `CS_GROUPS` as detailed below, to make this easy we provide some scripts utilizing `jq` to transform the formats for you.

### Converting from concourse-summary

#### HOSTS

```
HOSTS="ci.concourse.ci appdog.ci.cf-app.com buildpacks.ci.cf-app.com diego.ci.cf-app.com capi.ci.cf-app.com" ./scripts/migrate_hosts.sh
```

#### CS_GROUPS

```
CS_GROUPS='{"test":{"buildpacks.ci.cf-app.com":{"binary-builder":["automated-builds","manual-builds"],"brats":null},"diego.ci.cf-app.com":{"greenhouse":null},"capi.ci.cf-app.com":null}}' ./scripts/migrate_cs_groups.sh
```

### Usage

As this app is written in [GoLang](https://golang.org/) it can be run in a number of ways:

#### Using go run

```
go run main.go
```

#### As a binary

```
go build
./go-concourse-summary
```

#### As a CF app

You may want to modify the example `manifest.yml` file prior to running your CF push

```
cf push
```

**Note:** For the purpose of migrations to show all groups for a pipeline you can either run omit `groups` from `CS_GROUPS` entirely, set it as an empty array (`[]`) or set it with a single value of `["all"]`. However if you use `all` and the pipeline has a group of `all` then only that group will be displayed.

All configuration is managed using environment variables:

| Variable            | Description                                                                               | Example                                                                                                                                                                                                                                                                    |
| ------------------- | ----------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| HOSTS               | A json array of all concourse hosts that you wish to have a dashboard for                 | '["ci.concourse.ci", "appdog.ci.cf-app.com", "buildpacks.ci.cf-app.com", "diego.ci.cf-app.com", "capi.ci.cf-app.com"]'                                                                                                                                                     |
| CS_GROUPS           | A json string of a chosen group name, linking to a host, pipeline and groups in concourse | '[{"group":"test","hosts":[{"fqdn":"buildpacks.ci.cf-app.com","pipelines":[{"groups":["automated-builds","manual-builds"],"name":"binary-builder"},{"name":"brats"}]},{"fqdn":"capi.ci.cf-app.com"},{"fqdn":"diego.ci.cf-app.com","pipelines":[{"name":"greenhouse"}]}]}]' |
| SKIP_SSL_VALIDATION | If set to "true" then SSL Validation will be ignored for all hosts                        | "true"                                                                                                                                                                                                                                                                     |
| REFRESH_INTERVAL    | An integer in seconds for configuring the page refresh interval, defaults to 30           | 10                                                                                                                                                                                                                                                                         |

### Dependency management

This project uses [dep](https://github.com/golang/dep) to manage its dependencies.
