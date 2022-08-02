go-fleek
========
[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/mrusme/go-fleek) [![license](http://img.shields.io/badge/license-GPLv3-red.svg?style=flat)](https://raw.githubusercontent.com/mrusme/go-fleek/master/LICENSE)

Tiny Go library for the 
[Fleek API](https://docs.fleek.co/fleek-api/overview/).

## Installation

```sh
go get -u github.com/mrusme/go-fleek
```


## Getting Started


### Querying Sites by Team ID

```go
package main

import (
  "log"
  "github.com/mrusme/go-fleek"
)

func main() {
  f := fleek.New("apiKeyHere")

  sites, err := f.GetSitesByTeamId("my-team")
  if err != nil {
    log.Panic(err)
  }

  for _, site := range sites {
    log.Printf(
      "Site ID: %v\nName: %s\nPlatform: %s\n\n",
      site.Id,
      site.Name,
      site.Platform,
    )
  }
}
```

