# logbeat [![Build Status](https://travis-ci.org/xentek/logbeat.svg?branch=master)](https://travis-ci.org/xentek/logbeat) [![GoDoc](https://godoc.org/github.com/xentek/logbeat?status.svg)](https://godoc.org/github.com/xentek/logbeat) [![Go Report Card](https://goreportcard.com/badge/github.com/xentek/logbeat)](https://goreportcard.com/report/github.com/xentek/logbeat) [![Coverage Status](https://coveralls.io/repos/github/xentek/logbeat/badge.svg)](https://coveralls.io/github/xentek/logbeat) [![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/xentek/logbeat/blob/master/LICENSE)

`logbeat` is a [Logrus](https://github.com/sirupsen/logrus) Hook to report errors to [Opbeat](https://opbeat.com/).

---

## Install

    go get https://github.com/xentek/logbeat

## Import

    import "github.com/xentek/logbeat"

## Usage

```go
package main

import (
    "os"

    log "github.com/sirupsen/logrus"
    logbeat "github.com/xentek/logbeat"
)

func init() {
    orgId := os.Getenv("OPBEAT_ORGANIZATION_ID")
    appId := os.Getenv("OPBEAT_APP_ID")
    token := os.Getenv("OPBEAT_SECRET_TOKEN")
    log.AddHook(logbeat.NewOpbeatHook(orgId, appId, token))
}

func main() {
    log.WithField("notify", "opbeat").Error("This error will be sent to Opbeat")
}
```

## Contributing

Refer to our [Contributor's Guide](CONTRIBUTING.md) to learn how you can participate in this project.

## More Info

  - [GoDoc](https://godoc.org/github.com/xentek/logbeat)
  - [Wiki](https://github.com/xentek/logbeat/wiki)
