# go-concourse-summary

This is a port of [concourse-summary](https://github.com/dgodd/concourse-summary) to Golang. The aim is for all features of `concourse-summary` to be covered in this port. In its current state all features should have been migrated with the exception of the ability to collapse/ expand groups.

The intention of `concourse-summary` is to show a quick overview of all of your [concourse](https://concourse.ci) pipelines and groups in a single summary page.

It was intended that all configuration would be compatible between `go-concourse-summary` and `concourse-summary`. Unfortunately due to the object model differences on Golang marshaling it has been necessary to rework the format of `HOSTS` and `CS_GROUPS` as detailed below. It is a goal to provide a scripted way of migrating these variables.

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

All configuration is managed using environment variables:

| Variable            | Description                                                                               | Example                                                                                                                                                                                                                                                                    |
| ------------------- | ----------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| HOSTS               | A json array of all concourse hosts that you wish to have a dashboard for                 | ["ci.concourse.ci", "appdog.ci.cf-app.com", "buildpacks.ci.cf-app.com", "diego.ci.cf-app.com", "capi.ci.cf-app.com"]                                                                                                                                                       |
| CS_GROUPS           | A json string of a chosen group name, linking to a host, pipeline and groups in concourse | '[{"group":"test","hosts":[{"fqdn":"buildpacks.ci.cf-app.com","pipelines":[{"groups":["automated-builds","manual-builds"],"name":"binary-builder"},{"name":"brats"}]},{"fqdn":"capi.ci.cf-app.com"},{"fqdn":"diego.ci.cf-app.com","pipelines":[{"name":"greenhouse"}]}]}]' |
| SKIP_SSL_VALIDATION | If set to "true" then SSL Validation will be ignored for all hosts                        | "true"                                                                                                                                                                                                                                                                     |
| REFRESH_INTERVAL    | An integer in seconds for configuring the page refresh interval, defaults to 30           | 10                                                                                                                                                                                                                                                                         |

### Dependency management

This project uses [dep](https://github.com/golang/dep) to manage its dependencies.
