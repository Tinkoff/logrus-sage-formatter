# Logrus SAGE Formatter

**Logrus SAGE formatter** is a tiny library which allows you to write logs in SAGE Logs format and uses ENV variables.

## Getting started

```
go get -u github.com/Tinkoff/logrus-sage-formatter.git
```

Use with environment variables.

```go
package main

import (
	log "github.com/sirupsen/logrus"
	sage "github.com/Tinkoff/logrus-sage-formatter"
)

func main() {
	log.SetFormatter(sage.NewFormatter(sage.MetadataFromEnv)) // Will automatically set logger metadata from environment variables

	log.Info("hello world!")
}
```

Use with environment metadata.

```go

package main

import (
	log "github.com/sirupsen/logrus"
	sage "github.com/Tinkoff/logrus-sage-formatter"
)

func main() {
	log.SetFormatter(sage.NewFormatter(sage.Metadata(sage.EnvironmentMetadata{
		Env: "prod",
		Group: "devplatform",
		System: "gitlab-manager",
    })))
	
	log.WithField("extraField", "fieldValue").Info("Hello world") // Will output extraField to .extra
}
```

## Environment variables

| Name | Description | Format | Example | 
| --- | --- | --- | --- |
| SAGE_DC | Application data center | ^[a-z][a-z0-9_-]*$ | m1 |
| SAGE_ENV | Environment | ^[a-z][a-z0-9_-]*$ | test |
| SAGE_GROUP | Access Group | ^[a-z][a-z0-9_-]*$ | default
| SAGE_SYSTEM | System (usually the name of the application) | ^[a-z][a-z0-9_-]*$ | filebeat
| SAGE_INST | The host from which the data are sent | ^[a-zA-Z0-9][a-zA-Z0-9_.-]*$ | ds-sage-cybertruck-prod01.example.ru
| SAGE_TIME_FORMAT | Time Format | yyyy-mm-ddThh:mm:ss.Z | 2006-01-02T15:04:05.000Z0700