# go-concourse-summary

This project is heavily influenced by [concourse-summary](https://github.com/dgodd/concourse-summary). As `concourse-summary` (in its current state its actually a clone with replacement backend) it intends to show a quick overview of all of your [concourse](https://concourse.ci) pipelines in groups in a single summary page.

`concourse-summary` is a brilliant app and provided exactly what you would expect from a summary view of concourse, unfortunately as it is written in [crystal](https://crystal-lang.org) which is still an alpha language we had many issues of compatability when either `crystal` or `concourse-summary` were updated.

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
