[![Build Status](https://travis-ci.org/comail/cologrus.svg?branch=master)](https://travis-ci.org/comail/cologrus)&nbsp;[![godoc reference](https://godoc.org/comail.io/go/cologrus?status.png)](https://godoc.org/comail.io/go/cologrus)

Package cologrus provides functionality to wrap [Logrus](https://github.com/Sirupsen/logrus) hooks and formatters as ready to use [CoLog](https://github.com/comail/colog) hooks and formatters.

### Example

```go
package main

import (
	"log"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/sentry"

	"comail.io/go/colog"
	"comail.io/go/cologrus"
)

func main() {
	colog.Register()
	colog.ParseFields(true)

	hook, err := logrus_sentry.NewSentryHook("... sentry DNS ...", []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})
	if err != nil {
		panic(err)
	}

	colog.AddHook(cologrus.NewLogrusHook(hook))
	colog.SetFormatter(cologrus.NewLogrusFormatter(new(logrus.TextFormatter)))
	
	log.Println("error: this is bad foo=bar")
}
```